package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io/ioutil"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
	app_pb "main/app_pb"
	pb "main/pb"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	hostTargets                = make(map[string][]string)
	serviceCount                = make(map[string]int)
	serviceTimeStamp            = make(map[string]time.Time)
	runningDeploymentController = make(map[string]bool) // this can be used for lock mechanism
	clientset                   *kubernetes.Clientset
	ttl_seconds                 int // time in seconds
	max_concurrency             int
	countLock                   sync.Mutex
)

// TODO - if deploymnet controller for a fucntion is already running dont run it again
func init() {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	ttl_seconds, _ = strconv.Atoi(os.Getenv("TTL"))
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
}

type baseHandle struct{}

func (h *baseHandle) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	// host := req.Host
	// hostWithoutPort, _, err := net.SplitHostPort(host)
	// if err != nil {
	//     hostWithoutPort = host
	// }

	pathParts := strings.Split(req.URL.Path, "/")
	func_name := pathParts[1]
	fmt.Printf("function name: %s", func_name)
	if func_name == "logs" {
		func_name_for_logs := pathParts[2]
		depControllerAddr := os.Getenv("DEP_CONTROLLER_ADD")
		for _, service := range hostTargets[func_name_for_logs] {
			serviceTimeStamp[service] = time.Now()
		}
		type_call := "logs"
		opts := grpc.WithInsecure()
		cc, err := grpc.Dial(depControllerAddr, opts)
		if err != nil {
			log.Fatal(err)
		}
		defer cc.Close()
		client := pb.NewDeploymentServiceClient(cc)
		request := &pb.DeploymentServiceRequest{Name: func_name_for_logs, FunctionCall: type_call}

		resp, err := client.Deployment(context.Background(), request)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(res, resp.Message)
		return
	} else {
		target := ""
		index := 0
		triggered := false
		for target == "" {
			countLock.Lock()
			for i, service := range hostTargets[func_name] {
				if serviceCount[service] < max_concurrency {
					target = service
					serviceCount[service] ++
					index = i
					serviceTimeStamp[service] = time.Now()
				}
			}
			if target == "" && !triggered  {
				triggered = true
				go callDepController(false, func_name,len(hostTargets[func_name]))
			} 
			countLock.Unlock()
		}
		fmt.Printf("forwarding request to %s", target)

		opts := grpc.WithInsecure()
		cc, err := grpc.Dial(target, opts)
		if err != nil {
			log.Fatal(err)
		}
		defer cc.Close()
		go callDepController(true, func_name, index)
		/*    optional bytes data = 1;
		string workflow_id = 2;
		int32 depth = 3;
		int32 width = 4;
		optional string request_type = 5;
		optional string pv_path = 6;*/
		payload, _ := ioutil.ReadAll(req.Body)
		workflow_id := strconv.Itoa(rand.Intn(100000))
		channel, _ := grpc.Dial(target, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*200), grpc.MaxCallSendMsgSize(1024*1024*200)), grpc.WithTransportCredentials(insecure.NewCredentials()))
		defer channel.Close()
		stub := app_pb.NewGRPCFunctionClient(channel)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		request_type := "gg"
		response, _ := stub.GRPCFunctionHandler(ctx, &app_pb.RequestBody{Data: payload, WorkflowId: workflow_id, Depth: 0, Width: 0, RequestType: &request_type})
		results := response.GetReply()
		countLock.Lock()
			serviceCount[target] --
		countLock.Unlock()
		fmt.Fprintf(res, results)
		return
	}
}

func callDepController(func_found bool, func_name string, replicaNumber int) error {
	if !runningDeploymentController[func_name] {
		runningDeploymentController[func_name] = true
		depControllerAddr := os.Getenv("DEP_CONTROLLER_ADD")
		if depControllerAddr == "" {
			fmt.Println("DEP_CONTROLLER_ADD environment variable not set")
			runningDeploymentController[func_name] = false
			return nil
		}
		//depControllerAddr = "http://" + depControllerAddr
		var type_call string
		if func_found {
			type_call = "existing_invoke"
		} else {
			type_call = "new_invoke"
		}
		opts := grpc.WithInsecure()
		cc, err := grpc.Dial(depControllerAddr, opts)
		if err != nil {
			log.Fatal(err)
		}
		defer cc.Close()
		client := pb.NewDeploymentServiceClient(cc)
		request := &pb.DeploymentServiceRequest{Name: func_name, FunctionCall: type_call, ReplicaNumber: int32(replicaNumber)}
		resp, err := client.Deployment(context.Background(), request)
		if err != nil {
			log.Fatal(err)
			return err
		}
		fmt.Printf("Receive response => %s ", resp.Message)

	}
	return nil
}

// 	targetURL := depControllerAddr + path + "?function=" + func_name
// 	_, err := http.Get(targetURL)
// 	if err != nil {
// 		fmt.Printf("Failed to make HTTP request: %v\n", err)
// 		runningDeploymentController[func_name] = false
// 		return err
// 	}
// 	runningDeploymentController[func_name] = false
// } else {

// 	log.Printf("Deployment Controller for function this %s is already in progress", func_name)
// }
// return nil

func checkTTL() {
	for {
		currentTime := time.Now()
		for name, timestamp := range serviceTimeStamp {
			elapsedTime := currentTime.Sub(timestamp)
			if elapsedTime.Seconds() > float64(ttl_seconds) {
				service_name := strings.Split(name, ".")[0] // this will give the service name which is same as deployment
				log.Printf("Deleting service and deployment of %s because of TTL\n", service_name)
				clientset.CoreV1().Services("").Delete(context.TODO(), service_name, metav1.DeleteOptions{})
				clientset.AppsV1().Deployments("").Delete(context.TODO(), service_name, metav1.DeleteOptions{})

			}
			time.Sleep(time.Second)
		}
	}

}

// done
func main() {
	log.Print("Ingress controller started")
	go checkTTL()
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("Failed to get in-cluster config: %s", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed to create clientset: %s", err)
	}
	max_concurrency, _ = strconv.Atoi(os.Getenv("MAX_CONCURRENCY"))
	watchService(clientset)
	h := &baseHandle{}
	http.Handle("/", h)
	server := &http.Server{
		Addr:    ":8080", // reverse proxy
		Handler: h,
	}
	log.Fatal(server.ListenAndServe())
}

// done
func watchService(clientset *kubernetes.Clientset) {
	watcher, err := clientset.CoreV1().Services("").Watch(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Failed to set up watch: %s", err)
	}

	go func() {
		for event := range watcher.ResultChan() {
			services, ok := event.Object.(*corev1.Service)
			if !ok {
				log.Printf("Expected Ingress type, got %T", event.Object)
				continue
			}

			switch event.Type {
			case watch.Added, watch.Modified:
				updateHostTargets(services)
			case watch.Deleted:
				deleteHostTargets(services)
			}

		}
	}()
}

// done
func updateHostTargets(service *corev1.Service) {
	for _, port := range service.Spec.Ports {
		serviceName := service.Name
		namespace := service.Namespace
		func_name := service.Labels["function_name"]
		hostname := fmt.Sprintf("%s.%s.svc.cluster.local:%d", serviceName, namespace, port.Port)
		countLock.Lock()
		serviceCount[hostname] = 0
		hostTargets[func_name] = append(hostTargets[func_name], hostname)
		countLock.Unlock()
	}
	countLock.Lock()
	log.Printf("Updated host targets: %+v", hostTargets)
	log.Printf("Service Count: %+v", serviceCount)
	countLock.Unlock()

}

func deleteHostTargets(service *corev1.Service) {
	func_name := ""
	hostname_deleted := ""

	for _, port := range service.Spec.Ports {
		serviceName := service.Name
		namespace := service.Namespace
		func_name = service.Labels["function_name"]
		hostname_deleted = fmt.Sprintf("%s.%s.svc.cluster.local:%d", serviceName, namespace, port.Port)
	}

	for i, val := range hostTargets[func_name] {
		if val == hostname_deleted {
			log.Printf("found the hostname in %d", i)
			countLock.Lock()
			hostTargets[func_name] = append(hostTargets[func_name][:i], hostTargets[func_name][i+1:]...)
			countLock.Unlock()
			break
		}
	}
	delete(serviceCount, hostname_deleted)
	delete(serviceTimeStamp, hostname_deleted)
	// log.Printf("deleted host targets: %+v", hostTargets)
	// log.Printf("delted Service Count: %+v", serviceCount)

}
