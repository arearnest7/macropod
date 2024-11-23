package main

import (
	pb "app/deployer_pb"
	wf_pb "app/wf_pb"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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
	Registry	string			`json:"registry"`
	Endpoints	[]string		`json:"endpoints,omitempty"`
	Envs		map[string]string	`json:"envs,omitempty"`
	Secrets		map[string]string	`json:"secrets,omitempty"`
}

type Workflow struct {
	Name		string			`json:"name,omitempty"`
	Functions	map[string]Function	`json:"functions"`
}

var (
	workflows				= make(map[string]Workflow)
	kclient					*kubernetes.Clientset
	hostTargets				= make(map[string][]string)
	standbyTargets				= make(map[string]string)
	serviceCount				= make(map[string]int)
	serviceTimeStamp			= make(map[string]time.Time)
	runningDeploymentController		= make(map[string]bool) // this can be used for lock mechanism
	ttl_seconds				int
	max_concurrency			 	int
	countLock				sync.Mutex
	macropod_namespace			string
	workflow_invocations			= make(map[string]int)
	replicaCount				= make(map[string]int)
	workflow_deployments			= make(map[string]int)
	retrypolicy				= `{
							"methodConfig": [{
							"name": [{}],
							"waitForReady": true,
							"retryPolicy": {
								"MaxAttempts": 30,
								"InitialBackoff": "1s",
								"MaxBackoff": "10s",
								"BackoffMultiplier": 2.0,
								"RetryableStatusCodes": [ "UNAVAILABLE", "UNKNOWN"]
							}
						}]}`
)

func internal_log(message string) {
	fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.000000 UTC") + ": " + message)
}

func ifPodsAreRunning(workflow_replica string) bool {
	label_replica := "workflow_replica=" + workflow_replica
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("Failed to get in-cluster config: %s", err)
	}
	k, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed to create kclient: %s", err)
	}
	pods, err := k.CoreV1().Pods(macropod_namespace).List(context.Background(), metav1.ListOptions{LabelSelector: label_replica})
	if err != nil {
		fmt.Println(err)
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

func callDepController(func_found bool, func_name string, replicaNumber int) error {
	depControllerAddr := os.Getenv("DEP_CONTROLLER_ADD")
	if depControllerAddr == "" {
		fmt.Println("DEP_CONTROLLER_ADD environment variable not set")
		return nil
	}
	var type_call string
	if func_found {
		type_call = "existing_invoke"
		internal_log("getting metrics from deployment controller")
	} else {
		type_call = "new_invoke"
		internal_log("invoking new replica")
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
	if type_call == "existing_invoke" {
		countLock.Lock()
		max_concurrency, _ = strconv.Atoi(resp.Message)
		countLock.Unlock()
	}
	fmt.Printf("Receive response => %s ", resp.Message)
	return nil
}

func checkTTL() {
	internal_log("Checking TTL")
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("Failed to get in-cluster config: %s", err)
	}
	k, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed to create kclient: %s", err)
	}
	for {
		currentTime := time.Now()
		for name, timestamp := range serviceTimeStamp {
			elapsedTime := currentTime.Sub(timestamp)
			if elapsedTime.Seconds() > float64(ttl_seconds) {
				service_name := strings.Split(name, ".")[0]
				log.Printf("deleting because of TTL %s", service_name)
				log.Printf("Deleting service and deployment of %s because of TTL\n", service_name)
				k.CoreV1().Services("macropod-functions").Delete(context.TODO(), service_name, metav1.DeleteOptions{})
				log.Print("deleted service")
				k.AppsV1().Deployments("macropod-functions").Delete(context.TODO(), service_name, metav1.DeleteOptions{})
				log.Print("deleted deployment")
				k.NetworkingV1().Ingresses("macropod-functions").Delete(context.TODO(), service_name, metav1.DeleteOptions{})
				log.Print("deleted ingress")
			}
			time.Sleep(time.Second)
		}
	}

}

//TODO
func updateHostTargets(ingress *networkingv1.Ingress) {
	for _, rule := range ingress.Spec.Rules {
		if rule.HTTP != nil {
			for _, path := range rule.HTTP.Paths {
				serviceName := path.Backend.Service.Name
				namespace := ingress.Namespace
				port := path.Backend.Service.Port.Number
				func_name := ingress.Labels["workflow_name"]
				replica_name := ingress.Labels["workflow_replica"]
				for !ifPodsAreRunning(replica_name) {
				}
				hostname := fmt.Sprintf("%s.%s.svc.cluster.local:%d", serviceName, namespace, port)
				internal_log("service found: " + hostname)
				if func_name != "" {
					countLock.Lock()
					serviceCount[hostname] = 0
					hostTargets[func_name] = append(hostTargets[func_name], hostname)
					countLock.Unlock()
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
				func_name = ingress.Labels["workflow_name"]
				hostname_deleted = fmt.Sprintf("%s.%s.svc.cluster.local:%d", serviceName, namespace, port)
			}
		}
	}
	log.Printf("deleting %s\n", hostname_deleted)
	log.Printf("hosttargets %v",hostTargets[func_name])
	for i, val := range hostTargets[func_name] {
		log.Printf("valuessssss: %s", val)
		if val == hostname_deleted {
			countLock.Lock()
			hostTargets[func_name] = append(hostTargets[func_name][:i], hostTargets[func_name][i+1:]...)
			log.Printf("deletingggggggg%v", hostTargets)
			delete(serviceCount, hostname_deleted)
			delete(serviceTimeStamp, hostname_deleted)
			workflow_deployments[func_name]--
			countLock.Unlock()
			break
		}
	}


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
				log.Printf("Deleted host targets")
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

//TODO
func Serve_WF_Invoke(res http.ResponseWriter, req *http.Request) {
	func_name := req.PathValue("func_name")
	internal_log("function name: " + func_name)
	countLock.Lock()
	workflow_invocations[func_name]++
	countLock.Unlock()
	target := ""
	triggered := false
	_, exists := runningDeploymentController[func_name]
	if !exists {
		runningDeploymentController[func_name] = false
	}
	for target == "" {
		countLock.Lock()
		for idx := range len(hostTargets[func_name]) {
			service := hostTargets[func_name][len(hostTargets[func_name])-1-idx]
			if serviceCount[service] < max_concurrency {
				target = service
				serviceCount[service]++
				serviceTimeStamp[service] = time.Now()
				break
			}
		}
		if target == "" && !triggered && !runningDeploymentController[func_name] && workflow_invocations[func_name] > max_concurrency * workflow_deployments[func_name] {
			runningDeploymentController[func_name] = true
			triggered = true
			replicaCount[func_name] = replicaCount[func_name] + 1
			go callDepController(false, func_name, replicaCount[func_name])
			workflow_deployments[func_name]++
			runningDeploymentController[func_name] = false
		}
		countLock.Unlock()

	}
	fmt.Println(target)
	internal_log("forwarding request to " + target)
	opts := grpc.WithInsecure()
	cc, err := grpc.Dial(target, opts)
	if err != nil {
		log.Fatal(err)
	}
	defer cc.Close()
	go callDepController(true, func_name, max_concurrency)
	payload, _ := ioutil.ReadAll(req.Body)
	workflow_id := strconv.Itoa(rand.Intn(100000))
	status := int32(0)
	var response *wf_pb.ResponseBody
	for status == 0 {
		channel, _ := grpc.Dial(target, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*200), grpc.MaxCallSendMsgSize(1024*1024*200)), grpc.WithTransportCredentials(insecure.NewCredentials()),grpc.WithDefaultServiceConfig(retrypolicy))
		defer channel.Close()
		stub := wf_pb.NewGRPCFunctionClient(channel)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		request_type := "gg"
		response, _ = stub.GRPCFunctionHandler(ctx, &wf_pb.RequestBody{Data: payload, WorkflowId: workflow_id, Depth: 0, Width: 0, RequestType: &request_type})
		status = response.GetCode()
	}
	countLock.Lock()
	serviceCount[target]--
	internal_log("decreasing the count for" + target)
	workflow_invocations[func_name]--
	countLock.Unlock()
	if status == 200 {
		fmt.Fprint(res, response)
	} else {
		http.Error(res, "Non 200 status code", http.StatusBadRequest)
	}
}

//TODO
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
	runningDeploymentController[req.PathValue("func_name")] = false
	workflow_invocations[req.PathValue("func_name")] = 0
	workflow_deployments[req.PathValue("func_name")] = 0
	internal_log("WF_CREATE_END " + req.PathValue("func_name"))
	fmt.Fprintf(res, "Workflow created successfully. Invoke your workflow with /invoke/"+req.PathValue("func_name")+"\n")
}

//TODO
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

	label_workflow := "workflow_name=" + req.PathValue("func_name")
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("Failed to get in-cluster config: %s", err)
	}
	k, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed to create kclient: %s", err)
	}
	services, err := k.CoreV1().Services("macropod-functions").List(context.TODO(), metav1.ListOptions{LabelSelector: label_workflow})
	if err != nil {
		fmt.Println(err)
	}
	for _, service := range services.Items {
		delete(serviceTimeStamp, service.Name)
		fmt.Println(service.ObjectMeta.Name)
		k.CoreV1().Services("macropod-functions").Delete(context.TODO(), service.ObjectMeta.Name, metav1.DeleteOptions{})
	}
	log.Print("deleted service for update")
	deployments, err := k.AppsV1().Deployments("macropod-functions").List(context.TODO(), metav1.ListOptions{LabelSelector: label_workflow})
	if err != nil {
		fmt.Println(err)
	}
	for _, deployment := range deployments.Items {
		fmt.Println(deployment.ObjectMeta.Name)
		k.AppsV1().Deployments("macropod-functions").Delete(context.TODO(), deployment.ObjectMeta.Name, metav1.DeleteOptions{})
	}
	log.Print("deleted deployment for update")
	ingresses, err := k.NetworkingV1().Ingresses("macropod-functions").List(context.TODO(), metav1.ListOptions{LabelSelector: label_workflow})
	if err != nil {
		fmt.Println(err)
	}
	for _, ingress := range ingresses.Items {
		fmt.Println(ingress.ObjectMeta.Name)
		k.NetworkingV1().Ingresses("macropod-functions").Delete(context.TODO(), ingress.ObjectMeta.Name, metav1.DeleteOptions{})
	}
	log.Print("deleted ingress for update")

	workflows[req.PathValue("func_name")] = body_u
	runningDeploymentController[req.PathValue("func_name")] = false
	workflow_invocations[req.PathValue("func_name")] = 0
	workflow_deployments[req.PathValue("func_name")] = 0
	internal_log("WF_UPDATE_END " + req.PathValue("func_name"))
	fmt.Fprintf(res, "Workflow \""+req.PathValue("func_name")+"\" has been updated successfully.\n")
}

//TODO
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
	label_workflow := "workflow_name=" + req.PathValue("func_name")
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("Failed to get in-cluster config: %s", err)
	}
	k, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed to create kclient: %s", err)
	}
	services, err := k.CoreV1().Services("macropod-functions").List(context.TODO(), metav1.ListOptions{LabelSelector: label_workflow})
	if err != nil {
		fmt.Println(err)
	}
	for _, service := range services.Items {
		delete(serviceTimeStamp, service.Name)
		fmt.Println(service.ObjectMeta.Name)
		k.CoreV1().Services("macropod-functions").Delete(context.TODO(), service.ObjectMeta.Name, metav1.DeleteOptions{})
	}
	log.Print("deleted service for delete")
	deployments, err := k.AppsV1().Deployments("macropod-functions").List(context.TODO(), metav1.ListOptions{LabelSelector: label_workflow})
	if err != nil {
		fmt.Println(err)
	}
	for _, deployment := range deployments.Items {
		fmt.Println(deployment.ObjectMeta.Name)
		k.AppsV1().Deployments("macropod-functions").Delete(context.TODO(), deployment.ObjectMeta.Name, metav1.DeleteOptions{})
	}
	log.Print("deleted deployment for delete")
	ingresses, err := k.NetworkingV1().Ingresses("macropod-functions").List(context.TODO(), metav1.ListOptions{LabelSelector: label_workflow})
	if err != nil {
		fmt.Println(err)
	}
	for _, ingress := range ingresses.Items {
		fmt.Println(ingress.ObjectMeta.Name)
		k.NetworkingV1().Ingresses("macropod-functions").Delete(context.TODO(), ingress.ObjectMeta.Name, metav1.DeleteOptions{})
	}
	log.Print("deleted ingress for delete")

	delete(workflows, req.PathValue("func_name"))
	delete(runningDeploymentController, req.PathValue("func_name"))
	delete(workflow_invocations, req.PathValue("func_name"))
	delete(workflow_deployments, req.PathValue("func_name"))
	internal_log("WF_DELETE_END " + req.PathValue("func_name"))
	fmt.Fprintf(res, "Workflow \""+req.PathValue("func_name")+"\" has been deleted successfully.\n")
}

func main() {
	log.Print("Ingress controller started")
	ttl_seconds, _ = strconv.Atoi(os.Getenv("TTL"))
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("Failed to get in-cluster config: %s", err)
	}
	kclient, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed to create kclient: %s", err)
	}
	go checkTTL()
	workflows = make(map[string]Workflow)
	max_concurrency = 5
	macropod_namespace = "macropod-functions"
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
