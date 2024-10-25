package main

import (
	"context"
	"fmt"
	"io/ioutil"
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

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var (
	hostTargets                 = make(map[string][]string)
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
	pathParts := strings.Split(req.URL.Path, "/")
	func_name := pathParts[1]
	results := ""
	fmt.Printf("function name: %s", func_name)
	target := ""
	triggered := false
	for target == "" {
		countLock.Lock()
		for _, service := range hostTargets[func_name] {
			if serviceCount[service] < max_concurrency {
				target = service
				serviceCount[service]++
				serviceTimeStamp[service] = time.Now()
			}
		}
		countLock.Unlock()
		if target == "" && !triggered {
			triggered = true
			go callDepController(false, func_name, len(hostTargets[func_name]))
			log.Print(target)
		}

	}
	log.Printf("forwarding request to %s", target)

	opts := grpc.WithInsecure()
	cc, err := grpc.Dial(target, opts)
	if err != nil {
		log.Fatal(err)
	}
	defer cc.Close()
	payload, _ := ioutil.ReadAll(req.Body)
	workflow_id := strconv.Itoa(rand.Intn(100000))
	channel, _ := grpc.Dial(target, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*200), grpc.MaxCallSendMsgSize(1024*1024*200)), grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer channel.Close()
	stub := app_pb.NewGRPCFunctionClient(channel)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	request_type := "gg"
	response, _ := stub.GRPCFunctionHandler(ctx, &app_pb.RequestBody{Data: payload, WorkflowId: workflow_id, Depth: 0, Width: 0, RequestType: &request_type})
	results = response.GetReply()
	countLock.Lock()
	serviceCount[target]--
	countLock.Unlock()
	fmt.Fprint(res, results)
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
				service_name := strings.Split(name, ".")[0]
				log.Print("deleting because of TTL %s", service_name)
				log.Printf("Deleting service and deployment of %s because of TTL\n", service_name)
				clientset.CoreV1().Services("").Delete(context.TODO(), service_name, metav1.DeleteOptions{})
				clientset.AppsV1().Deployments("").Delete(context.TODO(), service_name, metav1.DeleteOptions{})
				clientset.NetworkingV1().Ingresses("").Delete(context.TODO(), service_name, metav1.DeleteOptions{})
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
	max_concurrency = 3
	log.Print("watch ingresssss")
	watchIngress(clientset)
	h := &baseHandle{}
	http.Handle("/", h)
	server := &http.Server{
		Addr:    ":8081", // reverse proxy
		Handler: h,
	}
	log.Fatal(server.ListenAndServe())
}

// done
func watchIngress(clientset *kubernetes.Clientset) {
	watcher, err := clientset.NetworkingV1().Ingresses("").Watch(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Failed to set up watch: %s", err)
	}

	go func() {
		for event := range watcher.ResultChan() {
			ingress, ok := event.Object.(*networkingv1.Ingress)
			if !ok {
				log.Printf("Expected Ingress type, got %T", event.Object)
				continue
			}

			switch event.Type {
			case watch.Added, watch.Modified:
				log.Printf("Updated host targets: %+v", hostTargets)
				updateHostTargets(ingress)
			case watch.Deleted:
				deleteHostTargets(ingress)
			}

		}
	}()
}

// done
func updateHostTargets(ingress *networkingv1.Ingress) {
	for _, rule := range ingress.Spec.Rules {
		if rule.HTTP != nil {
			for _, path := range rule.HTTP.Paths {
				serviceName := path.Backend.Service.Name
				namespace := ingress.Namespace
				port := path.Backend.Service.Port.Number
				func_name := ingress.Labels["function_name"]
				hostname := fmt.Sprintf("%s.%s.svc.cluster.local:%d", serviceName, namespace, port)
				if func_name != "" {
					countLock.Lock()
					serviceCount[hostname] = 0
					hostTargets[func_name] = append(hostTargets[func_name], hostname)
					fmt.Print(hostname)
					countLock.Unlock()
					log.Print("hostname found : %v ", hostTargets)
				}
			}
		}
	}
}

func deleteHostTargets(ingress *networkingv1.Ingress) {
	func_name := ""
	hostname_deleted := ""
	for _, rule := range ingress.Spec.Rules {
		if rule.HTTP != nil {
			for _, path := range rule.HTTP.Paths {
				serviceName := path.Backend.Service.Name
				namespace := ingress.Namespace
				port := path.Backend.Service.Port.Number
				func_name = ingress.Labels["function_name"]
				hostname_deleted = fmt.Sprintf("%s.%s.svc.cluster.local:%d", serviceName, namespace, port)
			}
		}
	}
	log.Print("deleting %s\n", hostname_deleted)
	for i, val := range hostTargets[func_name] {
		if val == hostname_deleted {
			countLock.Lock()
			hostTargets[func_name] = append(hostTargets[func_name][:i], hostTargets[func_name][i+1:]...)
			countLock.Unlock()
			break
		}
	}
	delete(serviceCount, hostname_deleted)
	delete(serviceTimeStamp, hostname_deleted)

}
