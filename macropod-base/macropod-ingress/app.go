package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"log"
	wf_pb "app/wf_pb"
	pb "app/deployer_pb"
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

type Function struct {
	Registry  string            `json:"registry"`
	Endpoints []string          `json:"endpoints,omitempty"`
	Envs      map[string]string `json:"envs,omitempty"`
	Secrets   map[string]string `json:"secrets,omitempty"`
}

type Workflow struct {
	Name      string              `json:"name,omitempty"`
	Functions map[string]Function `json:"functions"`
}

var (
	workflows                   = make(map[string]Workflow)
	kclient                     *kubernetes.Clientset
	hostTargets                 = make(map[string][]string)
	serviceCount                = make(map[string]int)
	serviceTimeStamp            = make(map[string]time.Time)
	runningDeploymentController = make(map[string]bool) // this can be used for lock mechanism
	ttl_seconds                 int // time in seconds
	max_concurrency             int
	countLock                   sync.Mutex
)

func internal_log(message string) {
	fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.000000 UTC") + ": " + message)
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
		rn := int32(replicaNumber)
		request := &pb.DeploymentServiceRequest{Name: func_name, FunctionCall: type_call, ReplicaNumber: rn}
		resp, err := client.Deployment(context.Background(), request)
		if err != nil {
			log.Fatal(err)
			return err
		}
		fmt.Printf("Receive response => %s ", resp.Message)

	}
	return nil
}

func checkTTL() {
	internal_log("Checking TTL")
	for {
		currentTime := time.Now()
		for name, timestamp := range serviceTimeStamp {
			elapsedTime := currentTime.Sub(timestamp)
			if elapsedTime.Seconds() > float64(ttl_seconds) {
				service_name := strings.Split(name, ".")[0]
				log.Print("deleting because of TTL %s", service_name)
				log.Printf("Deleting service and deployment of %s because of TTL\n", service_name)
				kclient.CoreV1().Services("").Delete(context.TODO(), service_name, metav1.DeleteOptions{})
				kclient.AppsV1().Deployments("").Delete(context.TODO(), service_name, metav1.DeleteOptions{})
				kclient.NetworkingV1().Ingresses("").Delete(context.TODO(), service_name, metav1.DeleteOptions{})
			}
			time.Sleep(time.Second)
		}
	}

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

func watchIngress(kclient *kubernetes.Clientset) {
        watcher, err := kclient.NetworkingV1().Ingresses("").Watch(context.TODO(), metav1.ListOptions{})
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

func Serve_Help(res http.ResponseWriter, req *http.Request) {
	help_print := "macropod ingress\nplease use the following services to interact/deploy/invoke workflows:\n"
	help_print += "Invoke:\n - path: /invoke/[workflow_name]\n - purpose: invoke workflow at entrypoint\n - payload: JSON body for entrypoint function\n - output: string result\n"
	help_print += "Create:\n - path: /create/[workflow_name]\n - purpose: deploy workflow based on configuration\n - payload: macropod JSON configuration\n - output: string confirmation\n"
	help_print += "Update:\n - path: /update/[workflow_name]\n - purpose: update previously deployed workflow\n - payload: macropod JSON configuration\n - output: string confirmation\n"
	help_print += "Delete:\n - path: /delete/[workflow_name]\n - purpose: delete workflow of same name\n - payload: none\n - output: string confirmation\n"
	fmt.Fprint(res, help_print)
}

func Serve_WF_Invoke(res http.ResponseWriter, req *http.Request) {
        func_name := req.PathValue("func_name")
        results := ""
        internal_log("function name: " + func_name)
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
        internal_log("forwarding request to " + target)

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
        stub := wf_pb.NewGRPCFunctionClient(channel)
        ctx, cancel := context.WithTimeout(context.Background(), time.Second)
        defer cancel()
        request_type := "gg"
        response, _ := stub.GRPCFunctionHandler(ctx, &wf_pb.RequestBody{Data: payload, WorkflowId: workflow_id, Depth: 0, Width: 0, RequestType: &request_type})
        results = response.GetReply()
        countLock.Lock()
	serviceCount[target]--
        countLock.Unlock()
        fmt.Fprint(res, results)
}




func Serve_WF_Create(res http.ResponseWriter, req *http.Request) {
	internal_log("WF_CREATE_START " + req.PathValue("func_name"))
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		internal_log("create body - " + err.Error())
	}
	body_u := Workflow{}
	json.Unmarshal(body, &body_u)
	defer req.Body.Close()
	opts := grpc.WithInsecure()
	cc, err := grpc.Dial(os.Getenv("DEP_CONTROLLER_ADD"), opts)
	if err != nil {
		internal_log("create grpc - " + err.Error())
	}
	defer cc.Close()
	workflow := string(body)
	client := pb.NewDeploymentServiceClient(cc)
	request := &pb.DeploymentServiceRequest{Name: req.PathValue("func_name"), FunctionCall: "create", Workflow: &workflow}
	internal_log("requesting create for " + req.PathValue("func_name"))
	_, err = client.Deployment(context.Background(), request)
	internal_log("returned create for " + req.PathValue("func_name"))
	if err != nil {
		internal_log("create return - " + err.Error())
	}
	workflows[req.PathValue("func_name")] = body_u
	internal_log("WF_CREATE_END " + req.PathValue("func_name"))
	fmt.Fprintf(res, "Workflow created successfully. Invoke your workflow with /invoke/"+req.PathValue("func_name"))
}

func Serve_WF_Update(res http.ResponseWriter, req *http.Request) {
	internal_log("WF_UPDATE_START " + req.PathValue("func_name"))
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		internal_log("update body - " + err.Error())
	}
	body_u := Workflow{}
	json.Unmarshal(body, &body_u)
	defer req.Body.Close()
	opts := grpc.WithInsecure()
	cc, err := grpc.Dial(os.Getenv("DEP_CONTROLLER_ADD"), opts)
	if err != nil {
		internal_log("update grpc - " + err.Error())
	}
	defer cc.Close()
	workflow := string(body)
	client := pb.NewDeploymentServiceClient(cc)
	request := &pb.DeploymentServiceRequest{Name: req.PathValue("func_name"), FunctionCall: "update", Workflow: &workflow}
	internal_log("requesting update for " + req.PathValue("func_name"))
	_, err = client.Deployment(context.Background(), request)
	internal_log("returned update for " + req.PathValue("func_name"))
	if err != nil {
		internal_log("update return - " + err.Error())
	}
	workflows[req.PathValue("func_name")] = body_u
	internal_log("WF_UPDATE_END " + req.PathValue("func_name"))
	fmt.Fprintf(res, "Workflow \""+req.PathValue("func_name")+"\" has been updated successfully.")
}

func Serve_WF_Delete(res http.ResponseWriter, req *http.Request) {
	internal_log("WF_DELETE_START " + req.PathValue("func_name"))
	opts := grpc.WithInsecure()
	cc, err := grpc.Dial(os.Getenv("DEP_CONTROLLER_ADD"), opts)
	if err != nil {
		internal_log("delete grpc - " + err.Error())
	}
	defer cc.Close()
	client := pb.NewDeploymentServiceClient(cc)
	request := &pb.DeploymentServiceRequest{Name: req.PathValue("func_name"), FunctionCall: "delete"}
	internal_log("requesting delete for " + req.PathValue("func_name"))
	_, err = client.Deployment(context.Background(), request)
	internal_log("returned delete for " + req.PathValue("func_name"))
	if err != nil {
		internal_log("delete return - " + err.Error())
	}
	delete(workflows, req.PathValue("func_name"))
	internal_log("WF_DELETE_END " + req.PathValue("func_name"))
	fmt.Fprintf(res, "Workflow \""+req.PathValue("func_name")+"\" has been deleted successfully.")
}

func Serve_Logs(res http.ResponseWriter, req *http.Request) {
	internal_log("LOGS_START " + req.PathValue("func_name"))
	opts := grpc.WithInsecure()
	cc, err := grpc.Dial(os.Getenv("DEP_CONTROLLER_ADD"), opts)
	if err != nil {
		internal_log("logs grpc - " + err.Error())
	}
	defer cc.Close()
	client := pb.NewDeploymentServiceClient(cc)
	request := &pb.DeploymentServiceRequest{Name: req.PathValue("func_name"), FunctionCall: "logs"}
	internal_log("requesting logs for " + req.PathValue("func_name"))
	response, err := client.Deployment(context.Background(), request)
	internal_log("returned logs for " + req.PathValue("func_name"))
	if err != nil {
		internal_log("logs return - " + err.Error())
	}
	internal_log("LOGS_END " + req.PathValue("func_name"))
	fmt.Fprintf(res, response.Message)
}

func Serve_Metrics(res http.ResponseWriter, req *http.Request) {
	internal_log("METRICS_START")
	opts := grpc.WithInsecure()
	cc, err := grpc.Dial(os.Getenv("DEP_CONTROLLER_ADD"), opts)
	if err != nil {
		internal_log("metrics grpc - " + err.Error())
	}
	defer cc.Close()
	client := pb.NewDeploymentServiceClient(cc)
	request := &pb.DeploymentServiceRequest{Name: "", FunctionCall: "metrics"}
	internal_log("requesting metrics")
	response, err := client.Deployment(context.Background(), request)
	internal_log("returned metrics")
	if err != nil {
		internal_log("metrics return - " + err.Error())
	}
	internal_log("METRICS_END")
	fmt.Fprintf(res, response.Message)
}




func main() {
	log.Print("Ingress controller started")
	go checkTTL()
	ttl_seconds, _ = strconv.Atoi(os.Getenv("TTL"))
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("Failed to get in-cluster config: %s", err)
	}

	kclient, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed to create kclient: %s", err)
	}
	workflows = make(map[string]Workflow)
	max_concurrency = 3
	log.Print("watch ingress")
	watchIngress(kclient)
	h := http.NewServeMux()
	h.HandleFunc("/", Serve_Help)
	h.HandleFunc("/invoke/{func_name}", Serve_WF_Invoke)
	h.HandleFunc("/create/{func_name}", Serve_WF_Create)
	h.HandleFunc("/update/{func_name}", Serve_WF_Update)
	h.HandleFunc("/delete/{func_name}", Serve_WF_Delete)
	http.ListenAndServe(":"+os.Getenv("SERVICE_PORT"), h)
}
