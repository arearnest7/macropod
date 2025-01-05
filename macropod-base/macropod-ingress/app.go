package main

import (
	pb "app/deployer_pb"
	wf_pb "app/wf_pb"
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"io/ioutil"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	toolsWatch "k8s.io/client-go/tools/watch"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Function struct {
	Registry  string            `json:"registry"`
	Endpoints []string          `json:"endpoints,omitempty"`
	Envs	  map[string]string `json:"envs,omitempty"`
	Secrets   map[string]string `json:"secrets,omitempty"`
}

type Workflow struct {
	Name	  string              `json:"name,omitempty"`
	Functions map[string]Function `json:"functions"`
}

var (
	workflows                    = make(map[string]Workflow)
	workflow_target_concurrency  = make(map[string]int)
	workflow_invocations_current = make(map[string]int)
	workflow_invocations_total   = make(map[string]int)
	service_target               = make(map[string][]string)
	service_count                = make(map[string]int)
	service_timestamp            = make(map[string]time.Time)
	service_channel              = make(map[string][]*grpc.ClientConn)
	service_stub                 = make(map[string][]wf_pb.GRPCFunctionClient)

	ttl_seconds                  int
	debug                        int
	deployer_add                 string
	macropod_namespace           string

	dataLock                     sync.Mutex
	timerLock                    sync.Mutex
	depLock                      sync.Mutex

	retrypolicy                  = `{
		"methodConfig": [{
		"name": [{}],
		"waitForReady": true,
		"retryThrottling": {
			"maxTokens": 100,
			"tokenRatio": 0.1
		  },
		"retryPolicy": {
			"MaxAttempts": 3,
			"InitialBackoff": "1s",
			"MaxBackoff": "10s",
			"BackoffMultiplier": 2.0,
			"RetryableStatusCodes": ["UNAVAILABLE", "UNKNOWN"]
		}
	}]}`
)

func ifPodsAreRunning(workflow_replica string, namespace string) bool {
	label_replica := "workflow_replica=" + workflow_replica
	config, err := rest.InClusterConfig()
	if err != nil {
		if debug > 0 {
			fmt.Println("Failed to get in-cluster config: " + err.Error())
		}
		return false
	}
	k, err := kubernetes.NewForConfig(config)
	if err != nil {
		if debug > 0 {
			fmt.Println("Failed to create k: " + err.Error())
		}
	}
	pods, err := k.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{LabelSelector: label_replica})
	if err != nil {
		if debug > 0 {
			fmt.Println(err)
		}
		return false
	}
	for _, pod := range pods.Items {
		if pod.Status.Phase != "Running" {
			return false
		}
		for _, container_status := range pod.Status.ContainerStatuses {
			if container_status.State.Running == nil {
				return false
			}
		}
	}
	return true
}

func callDepController(type_call string, func_name string, payload string) error {
	opts := grpc.WithInsecure()
	cc, err := grpc.Dial(deployer_add, opts)
	if err != nil && debug > 0 {
		fmt.Println(err)
	}
	defer cc.Close()
	client := pb.NewDeploymentServiceClient(cc)
	request := &pb.DeploymentServiceRequest{Name: func_name, FunctionCall: type_call}
	switch type_call {
		case "create_workflow":
			request.Workflow = &payload
		case "update_workflow":
			request.Workflow = &payload
		case "delete_workflow":
			if debug > 2 {
				fmt.Println("calling deployer delete")
			}
		case "create_deployment":
			rn, _ := strconv.Atoi(payload)
			request.ReplicaNumber = int32(rn)
		case "update_deployments":
			rn, _ := strconv.Atoi(payload)
			request.ReplicaNumber = int32(rn)
		case "ttl_delete":
			request.Workflow = &payload
	}
	depLock.Lock()
	resp, err := client.Deployment(context.Background(), request)
	depLock.Unlock()
	if err != nil && debug > 0 {
		fmt.Println(err)
		return err
	}
	switch resp_code := strings.Split(resp.Message, "."); resp_code[0] {
		case "0": // status ok
			if debug > 3 {
				fmt.Println("deployer controller return ok")
			}
		case "1": // already deploying
			if debug > 2 {
				fmt.Println("deployer already modifying workflow %s", func_name)
			}
		case "2": // target concurrency decrease
			if debug > 2 {
				fmt.Println("workflow hit threshold, scaling down target concurrency")
				workflow_target_concurrency[func_name] = workflow_invocations_current[func_name] / len(service_target[func_name]) / 2
			}
		default:
			if debug > 2 {
				fmt.Println("Unknown response code " + resp_code[0])
			}
	}

	return nil
}

func watchTTL() {
	if debug > 3 {
		fmt.Println("watch TTL")
	}
	config, err := rest.InClusterConfig()
	if err != nil && debug > 0 {
		fmt.Println("Failed to get in-cluster config: " + err.Error())
	}
	k, err := kubernetes.NewForConfig(config)
	if err != nil && debug > 0 {
		fmt.Println("Failed to create k: " + err.Error())
	}
	for {
		currentTime := time.Now()
		for name, timestamp := range service_timestamp {
			elapsedTime := currentTime.Sub(timestamp)
			if elapsedTime.Seconds() > float64(ttl_seconds) {
				service_name := strings.Split(name, ".")[0]
				if debug > 1 {
					log.Printf("Deleting service and deployment of %s because of TTL\n", service_name)
				}
				service, exists := k.CoreV1().Services(macropod_namespace).Get(context.Background(), service_name, metav1.GetOptions{})
				if exists == nil {
					dataLock.Lock()
					func_name := service.Labels["workflow_name"]
					labels := service.Labels["workflow_replica"]
					dataLock.Unlock()
					labels_replica := "workflow_replica=" + labels
					services, err := k.CoreV1().Services(macropod_namespace).List(context.Background(), metav1.ListOptions{LabelSelector: labels_replica})
					if err != nil && debug > 0 {
						fmt.Println(err)
					}
					for _, service := range services.Items {
						dataLock.Lock()
						delete(service_count, service.Name)
						delete(service_timestamp, service.Name)
						delete(service_channel, service.Name)
						delete(service_stub, service.Name)
						dataLock.Unlock()
					}
					go callDepController("ttl_delete", func_name, labels)
				}
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
				dataLock.Lock()
				func_name := ingress.Labels["workflow_name"]
				hostname := fmt.Sprintf("%s.%s.svc.cluster.local:%d", serviceName, namespace, port)
				replica_name := ingress.Labels["workflow_replica"]
				dataLock.Unlock()
				for !ifPodsAreRunning(replica_name, namespace) {}
				if debug > 2 {
					fmt.Println("service found: " + hostname)
				}
				channel, _ := grpc.Dial(hostname, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*200), grpc.MaxCallSendMsgSize(1024*1024*200)), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultServiceConfig(retrypolicy))
				stub := wf_pb.NewGRPCFunctionClient(channel)
				for channel.GetState() != connectivity.Ready {
					time.Sleep(10 * time.Millisecond)
				}
				dataLock.Lock()
				service_count[hostname] = 0
				service_timestamp[hostname] = time.Now()
				service_target[func_name] = append(service_target[func_name], hostname)
				service_channel[hostname] = append(service_channel[hostname], channel)
				service_stub[hostname] = append(service_stub[hostname], stub)
				dataLock.Unlock()
			}
		}
	}
}

func deleteHostTargets(ingress *networkingv1.Ingress) {
	func_name := ""
	namespace := ""
	hostname_deleted := ""
	for _, rule := range ingress.Spec.Rules {
		if rule.HTTP != nil {
			for _, path := range rule.HTTP.Paths {
				serviceName := path.Backend.Service.Name
				namespace = ingress.Namespace
				port := path.Backend.Service.Port.Number
				dataLock.Lock()
				func_name = ingress.Labels["workflow_name"]
				dataLock.Unlock()
				hostname_deleted = fmt.Sprintf("%s.%s.svc.cluster.local:%d", serviceName, namespace, port)
			}
		}
	}
	for i, val := range service_target[func_name] {
		if val == hostname_deleted {
			dataLock.Lock()
			service_target[func_name] = append(service_target[func_name][:i], service_target[func_name][i+1:]...)
			dataLock.Unlock()
			if debug > 2 {
				log.Printf("deleting %v", hostname_deleted)
			}
			if _, exists := service_count[hostname_deleted]; exists {
				dataLock.Lock()
				delete(service_count, hostname_deleted)
				delete(service_timestamp, hostname_deleted)
				delete(service_channel, hostname_deleted)
				delete(service_stub, hostname_deleted)
				dataLock.Unlock()
			}
			break
		}

	}

}

func watchIngress() {
	config, err := rest.InClusterConfig()
	if err != nil && debug > 0 {
		fmt.Println("Failed to get in-cluster config: " + err.Error())
	}
	k, err := kubernetes.NewForConfig(config)
	if err != nil && debug > 0{
		fmt.Println("Failed to create k: " + err.Error())
	}
	//reference: https://blog.mimacom.com/k8s-watch-resources/
	watchFunc := func(options metav1.ListOptions) (watch.Interface, error) {
		return k.NetworkingV1().Ingresses("").Watch(context.Background(), metav1.ListOptions{})
	}
	watcher, _ := toolsWatch.NewRetryWatcher("1", &cache.ListWatch{WatchFunc: watchFunc})
	for event := range watcher.ResultChan() {
		ingress, ok := event.Object.(*networkingv1.Ingress)
		if !ok {
			continue
		}
		switch event.Type {
		case watch.Added, watch.Modified:
			updateHostTargets(ingress)
		case watch.Deleted:
			deleteHostTargets(ingress)
		}

	}
}

func Serve_Help(res http.ResponseWriter, req *http.Request) {
	help_print := "macropod ingress\nplease use the following services to interact/deploy/invoke workflows:\n"
	help_print += "Invoke:\n - path: /invoke/[workflow_name]\n - purpose: invoke workflow at entrypoint\n - payload: JSON body for entrypoint function\n - output: string result\n"
	help_print += "Create:\n - path: /create/[workflow_name]\n - purpose: deploy workflow based on configuration\n - payload: macropod JSON configuration\n - output: string confirmation\n"
	help_print += "Update:\n - path: /update/[workflow_name]\n - purpose: update previously deployed workflow\n - payload: macropod JSON configuration\n - output: string confirmation\n"
	help_print += "Delete:\n - path: /delete/[workflow_name]\n - purpose: delete workflow of same name\n - payload: none\n - output: string confirmation\n"
	fmt.Fprint(res, help_print)
}

func getTargets(triggered bool, func_name string, target string) (string, bool) {
	dataLock.Lock()
	if len(service_target[func_name]) > 0 {
		idx := workflow_invocations_total[func_name] % len(service_target[func_name])
		service := service_target[func_name][idx]
		target_concurrency, target_concurrency_set := workflow_target_concurrency[func_name]
		channel, service_channel_exists := service_channel[service]
		if service_channel_exists && channel[0].GetState() == connectivity.Ready {
			if !target_concurrency_set || service_count[service] < target_concurrency {
				target = service
				service_count[service]++
				service_timestamp[service] = time.Now()
			}
		}
	}
	if !triggered && len(service_target[func_name]) < workflow_invocations_current[func_name]{
		triggered = true
		go callDepController("create_deployment", func_name, strconv.Itoa(len(service_target[func_name])))
	}
	dataLock.Unlock()
	return target, triggered
}

func Serve_WF_Invoke(res http.ResponseWriter, req *http.Request) {
	func_name := req.PathValue("func_name")
	dataLock.Lock()
	workflow_invocations_current[func_name]++
	workflow_invocations_total[func_name]++
	dataLock.Unlock()
	target := ""
	triggered := false
	look_arget := time.Now()
	for target == "" {
		target, triggered = getTargets(triggered, func_name, target)
	}
	payload, _ := ioutil.ReadAll(req.Body)
	workflow_id := strconv.Itoa(rand.Intn(100000))
	status := int32(0)
	var response *wf_pb.ResponseBody
	go callDepController("update_deployments", func_name, strconv.Itoa(workflow_invocations_current[func_name]))
	request_type := "gg"
	dataLock.Lock()

	dataLock.Unlock()
	stub := service_stub[target][0]
	start_time := time.Now()
	response, err := stub.GRPCFunctionHandler(context.Background(), &wf_pb.RequestBody{Data: payload, WorkflowId: workflow_id, Depth: 0, Width: 0, RequestType: &request_type})
	end_time := time.Now()
	status = response.GetCode()
	if status == 200 {
		if debug > 2 {
			fmt.Println("request was served by: " + target + " in " + end_time.Sub(start_time).String() + " time to find target " + start_time.Sub(look_arget).String())
		}
		fmt.Fprint(res, response.GetReply())

	} else {
		if debug > 2 {
			log.Printf("Non 200 code %s: %d, error : %s in : %s time ti find target : %s", target, response.GetCode(), err.Error(), end_time.Sub(start_time).String(), start_time.Sub(look_arget).String())
		}
		http.Error(res, err.Error(), http.StatusBadGateway)
	}
	dataLock.Lock()
	service_count[target]--
	workflow_invocations_current[func_name]--
	dataLock.Unlock()
}

func Serve_WF_Create(res http.ResponseWriter, req *http.Request) {
	if debug > 2 {
		fmt.Println("WF_CREATE_START " + req.PathValue("func_name"))
	}
	func_name := req.PathValue("func_name")
	_, exists := workflows[func_name]
	if exists {
		if debug > 2 {
			fmt.Fprintf(res, "Workflow " +func_name+" already exists\n")
		}
		return
	}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil && debug > 0{
		fmt.Println("create body - " + err.Error())
	}
	body_u := Workflow{}
	json.Unmarshal(body, &body_u)
	defer req.Body.Close()
	workflow := string(body)
	callDepController("create_workflow", func_name, workflow)
	dataLock.Lock()
	workflows[req.PathValue("func_name")] = body_u
	workflow_invocations_current[req.PathValue("func_name")] = 0
	workflow_invocations_total[req.PathValue("func_name")] = 0
	dataLock.Unlock()
	if debug > 2 {
		fmt.Println("WF_CREATE_END " + func_name)
	}
	fmt.Fprintf(res, "Workflow created successfully. Invoke your workflow with /invoke/"+func_name+"\n")
}

func Serve_WF_Update(res http.ResponseWriter, req *http.Request) {
	if debug > 2 {
		fmt.Println("WF_UPDATE_START " + req.PathValue("func_name"))
	}
	label_workflow := "workflow_name=" + req.PathValue("func_name")
	config, err := rest.InClusterConfig()
	if err != nil && debug > 0 {
		fmt.Println("Failed to get in-cluster config: " + err.Error())
	}
	k, err := kubernetes.NewForConfig(config)
	if err != nil && debug > 0 {
		fmt.Println("Failed to create k: " + err.Error())
	}
	services, err := k.CoreV1().Services(macropod_namespace).List(context.Background(), metav1.ListOptions{LabelSelector: label_workflow})
	if err != nil && debug > 0 {
		fmt.Println(err)
	}
	for _, service := range services.Items {
		dataLock.Lock()
		delete(service_target, service.Name)
		delete(service_count, service.Name)
		delete(service_timestamp, service.Name)
		delete(service_channel, service.Name)
		delete(service_stub, service.Name)
		dataLock.Unlock()
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil && debug > 0 {
		fmt.Println("update body - " + err.Error())
	}
	body_u := Workflow{}
	json.Unmarshal(body, &body_u)
	defer req.Body.Close()
	workflow := string(body)
	callDepController("update_workflow", req.PathValue("func_name"), workflow)
	dataLock.Lock()
	workflows[req.PathValue("func_name")] = body_u
	delete(workflow_target_concurrency, req.PathValue("func_name"))
	workflow_invocations_current[req.PathValue("func_name")] = 0
	workflow_invocations_total[req.PathValue("func_name")] = 0
	dataLock.Unlock()
	if debug > 2 {
		fmt.Println("WF_UPDATE_END " + req.PathValue("func_name"))
	}
	fmt.Fprintf(res, "Workflow \""+req.PathValue("func_name")+"\" has been updated successfully.\n")
}

func Serve_WF_Delete(res http.ResponseWriter, req *http.Request) {
	if debug > 2 {
		fmt.Println("WF_DELETE_START " + req.PathValue("func_name"))
	}
	label_workflow := "workflow_name=" + req.PathValue("func_name")
	config, err := rest.InClusterConfig()
	if err != nil && debug > 0 {
		fmt.Println("Failed to get in-cluster config: " + err.Error())
	}
	k, err := kubernetes.NewForConfig(config)
	if err != nil && debug > 0 {
		fmt.Println("Failed to create k: " + err.Error())
	}
	services, err := k.CoreV1().Services(macropod_namespace).List(context.Background(), metav1.ListOptions{LabelSelector: label_workflow})
	if err != nil && debug > 0 {
		fmt.Println(err)
	}
	for _, service := range services.Items {
		dataLock.Lock()
		delete(service_target, service.Name)
		delete(service_count, service.Name)
		delete(service_timestamp, service.Name)
		delete(service_channel, service.Name)
		delete(service_stub, service.Name)
		dataLock.Unlock()
	}

	callDepController("delete_workflow", req.PathValue("func_name"), "")

	dataLock.Lock()
	delete(workflows, req.PathValue("func_name"))
	delete(workflow_target_concurrency, req.PathValue("func_name"))
	delete(workflow_invocations_current, req.PathValue("func_name"))
	delete(workflow_invocations_total, req.PathValue("func_name"))
	dataLock.Unlock()
	if debug > 2 {
		fmt.Println("WF_DELETE_END " + req.PathValue("func_name"))
	}
	fmt.Fprintf(res, "Workflow \""+req.PathValue("func_name")+"\" has been deleted successfully.\n")
}

func main() {
	var err error
	service_port := os.Getenv("SERVICE_PORT")
	if service_port == "" {
		service_port = "8081"
	}
	deployer_add = os.Getenv("DEP_CONTROLLER_ADD")
	if deployer_add == "" {
		deployer_add = "127.0.0.1:8082"
	}
	macropod_namespace = os.Getenv("MACROPOD_NAMESPACE")
	if macropod_namespace == "" {
		macropod_namespace = "macropod-functions"
	}
	ttl_seconds, err = strconv.Atoi(os.Getenv("TTL"))
	if err != nil {
		ttl_seconds = 180
	}
	debug, err = strconv.Atoi(os.Getenv("DEBUG"))
	if err != nil {
		debug = 0
	}
	go watchTTL()
	go watchIngress()
	h := http.NewServeMux()
	h.HandleFunc("/", Serve_Help)
	h.HandleFunc("/invoke/{func_name}", Serve_WF_Invoke)
	h.HandleFunc("/create/{func_name}", Serve_WF_Create)
	h.HandleFunc("/update/{func_name}", Serve_WF_Update)
	h.HandleFunc("/delete/{func_name}", Serve_WF_Delete)
	http.ListenAndServe(":" + service_port, h)
}
