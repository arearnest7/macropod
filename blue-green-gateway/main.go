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
	hostTargets                 = make(map[string]string)
	runningDeploymentController = make(map[string]bool) // this can be used for lock mechanism
	function_timestamp          = make(map[string]time.Time)
	clientset                   *kubernetes.Clientset
	ttl_seconds                 int // time in seconds
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
	if func_name == "logs" {
		func_name_for_logs := pathParts[2]
		depControllerAddr := os.Getenv("DEP_CONTROLLER_ADD")
		function_timestamp[func_name_for_logs] = time.Now()
		log.Print(function_timestamp)
		depControllerAddr = depControllerAddr //grpc
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
		function_timestamp[func_name] = time.Now()
		log.Print(function_timestamp)
		req.URL.Path = ""
		target, ok := hostTargets[func_name]
		if !ok {
			http.Error(res, "Target not found making the function workflow", http.StatusNotFound)
			runningDeploymentController[func_name] = false
			callDepController(false, func_name)

			return
		}
		// remoteUrl, err := url.Parse(target)
		// if err != nil {
		// 	log.Fatal("target parse fail:", err)
		// }
		go callDepController(true, func_name)
		// log.Print("Forwarding request to", remoteUrl)
		// proxy := httputil.NewSingleHostReverseProxy(remoteUrl)
		// proxy.ServeHTTP(res, req)

		opts := grpc.WithInsecure()
		cc, err := grpc.Dial(target, opts)
		if err != nil {
			log.Fatal(err)
		}
		defer cc.Close()

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
		fmt.Fprintf(res, results)
		return
	}
}

func callDepController(func_found bool, func_name string) error {
	if !runningDeploymentController[func_name] {
		runningDeploymentController[func_name] = true
		depControllerAddr := os.Getenv("DEP_CONTROLLER_ADD")
		if depControllerAddr == "" {
			fmt.Println("DEP_CONTROLLER_ADD environment variable not set")
			runningDeploymentController[func_name] = false
			return nil
		}
		//depControllerAddr = "http://" + depControllerAddr
		depControllerAddr = depControllerAddr //grpc
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
		request := &pb.DeploymentServiceRequest{Name: func_name, FunctionCall: type_call}

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

func fetchUpdatedTimeStamp() map[string]time.Time {
	//log.Print(function_timestamp)
	return function_timestamp
}
func checkTTL() {
	for {
		currentTime := time.Now()
		function_timestamp_updated := fetchUpdatedTimeStamp()
		//log.Print(function_timestamp_updated)
		for name, timestamp := range function_timestamp_updated {
			elapsedTime := currentTime.Sub(timestamp)
			if elapsedTime.Seconds() > float64(ttl_seconds) {
				log.Printf("Deleting resources of function %s\n", name)
				namespaces, _ := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
				//log.Print(namespaces)
				for _, ns := range namespaces.Items {
					namespace := ns.Name
					if strings.Contains(namespace, name) {
						log.Printf("Namespace deletion %s", namespace)
						clientset.CoreV1().Namespaces().Delete(context.TODO(), namespace, metav1.DeleteOptions{})
					}
				}
				delete(function_timestamp, name)
			}

		}
		time.Sleep(time.Second)
	}

}

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

	watchIngress(clientset)
	h := &baseHandle{}
	http.Handle("/", h)
	server := &http.Server{
		Addr:    ":8080", // reverse proxy
		Handler: h,
	}
	log.Fatal(server.ListenAndServe())
}

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

func updateHostTargets(ingress *networkingv1.Ingress) {
	for _, rule := range ingress.Spec.Rules {
		if rule.HTTP != nil {
			for _, path := range rule.HTTP.Paths {
				serviceName := path.Backend.Service.Name
				namespace := ingress.Namespace
				port := path.Backend.Service.Port.Number
				func_name := ingress.Labels["function_name"]
				hostname := fmt.Sprintf("%s.%s.svc.cluster.local:%d", serviceName, namespace, port)
				hostTargets[func_name] = hostname
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
				hostname_deleted = fmt.Sprintf("http://%s.%s.svc.cluster.local:%d", serviceName, namespace, port)
			}
		}
	}

	// delete entry only when deleted ingress is same as present ingress
	if hostTargets[func_name] == hostname_deleted {
		log.Printf("Deleting ingress for %s", func_name)
		delete(hostTargets, func_name)
	}

}
