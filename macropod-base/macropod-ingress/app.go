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
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"
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
	workflows   = make(map[string]Workflow)
	hostTargets = make(map[string][]string)
	//standbyTargets              = make(map[string]string)
	serviceCount                    = make(map[string]int)
	serviceTimeStamp                = make(map[string]time.Time)
	runningDeploymentController     = make(map[string]bool) // this can be used for lock mechanism
	ttl_seconds                     int
	deploymentStarted               bool
	overall_ttl                     int
	max_concurrency                 int
	dataLock                        sync.Mutex
	timerLock                       sync.Mutex
	depLock                         sync.Mutex
	replicaRunning                  []string
	deleteIngressList []*networkingv1.Ingress
	workflow_invocations            = make(map[string]int)
	serviceCreationTime             = make(map[string]time.Time)
	replicaCount                    = make(map[string]int)
	workflow_deployments            = make(map[string]int)
	overallWorkflowTTL              = make(map[string]time.Time)
	channelOfService                = make(map[string][]*grpc.ClientConn)
	stubsOfChannel                  = make(map[string][]wf_pb.GRPCFunctionClient)
	concurrencyList                 = make(map[string]int)
	currentReplicasThatCanBeHandled = make(map[string]int)
	invocationsCount                = make(map[string]int)
	debug                           int
	maxConcurrencyStanby            int
	services_list                   []string
	retrypolicy                     = `{
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

type NodeMetricList struct {
	Kind       string `json:"kind"`
	APIVersion string `json:"apiVersion"`
	Metadata   struct {
		SelfLink string `json:"selfLink"`
	} `json:"metadata"`
	Items []NodeMetricsItem `json:"items"`
}

type NodeMetricsItem struct {
	Metadata  Metadata  `json:"metadata"`
	Timestamp time.Time `json:"timestamp"`
	Window    string    `json:"window"`
	Usage     Usage     `json:"usage"`
}
type Metadata struct {
	Name              string            `json:"name,omitempty"`
	CreationTimestamp time.Time         `json:"creationTimestamp,omitempty"`
	Labels            map[string]string `json:"labels,omitempty"`
}

type Usage struct {
	CPU    string `json:"cpu"`
	Memory string `json:"memory"`
}

func internal_log(message string) {
	fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.000000 UTC") + ": " + message)
}

func ifPodsAreRunning(workflow_replica string, namespace string) bool {
	label_replica := "workflow_replica=" + workflow_replica
	config, err := rest.InClusterConfig()
	if err != nil {
		fmt.Println("Failed to get in-cluster config: " + err.Error())
		return false
	}
	k, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Println("Failed to create kclient: " + err.Error())
	}
	pods, err := k.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{LabelSelector: label_replica})
	if err != nil {
		log.Print(err)
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

func deleteDeployments(label_to_delete string, func_name string) {
	config, err := rest.InClusterConfig()
	if err != nil {
		fmt.Println("Failed to get in-cluster config: " + err.Error())
	}
	k, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Println("Failed to create kclient: " + err.Error())
	}
	for _, ingress := range deleteIngressList {
		service_name := ""
		label_replica := ingress.Labels["workflow_replica"]
		for _, rule := range ingress.Spec.Rules {
			if rule.HTTP != nil {
				for _, path := range rule.HTTP.Paths {
					serviceName := path.Backend.Service.Name
					services_list = append(services_list, serviceName)
					namespace := ingress.Namespace
					port := path.Backend.Service.Port.Number
					service_name = fmt.Sprintf("%s.%s.svc.cluster.local:%d", serviceName, namespace, port)
					log.Printf("deleting service %s", service_name)
					dataLock.Lock()
					delete(serviceTimeStamp, service_name)
					delete(stubsOfChannel, service_name)
					for i, val := range hostTargets[func_name] {
						if val == service_name {
							hostTargets[func_name] = append(hostTargets[func_name][:i], hostTargets[func_name][i+1:]...)
							break
						}
					}
					dataLock.Unlock()

				}
			}
		}
		delete_labels := "workflow_replica=" + label_replica
		services, err := k.CoreV1().Services("macropod-functions").List(context.Background(), metav1.ListOptions{LabelSelector: delete_labels})
		if err != nil {
			fmt.Println(err)
		}
		for _, service := range services.Items {
			k.CoreV1().Services("macropod-functions").Delete(context.Background(), service.ObjectMeta.Name, metav1.DeleteOptions{})
		}
		deployments, err := k.AppsV1().Deployments("macropod-functions").List(context.Background(), metav1.ListOptions{LabelSelector: delete_labels})
		if err != nil {
			log.Print(err)
		}
		for _, deployment := range deployments.Items {
			k.AppsV1().Deployments("macropod-functions").Delete(context.Background(), deployment.ObjectMeta.Name, metav1.DeleteOptions{})
		}
		ingresses, err := k.NetworkingV1().Ingresses("macropod-functions").List(context.Background(), metav1.ListOptions{LabelSelector: delete_labels})
		if err != nil {
			log.Print(err)
		}
		for _, ingress := range ingresses.Items {
			k.NetworkingV1().Ingresses("macropod-functions").Delete(context.Background(), ingress.ObjectMeta.Name, metav1.DeleteOptions{})
		}

		for {
			deployments_list, _ := k.CoreV1().Pods("macropod-functions").List(context.Background(), metav1.ListOptions{LabelSelector: delete_labels})
			if deployments_list == nil || len(deployments_list.Items) == 0 {
				fmt.Println("Deployment does not exist")
				break
			}
			time.Sleep(10 * time.Millisecond)
		}
		dataLock.Lock()
		currentReplicasThatCanBeHandled[func_name] -= concurrencyList[service_name]
		delete(concurrencyList, service_name)
		dataLock.Unlock()
		for !deploymentStarted{
			
		}
	}

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
	} else {
		type_call = "new_invoke"
		dataLock.Lock()
		deploymentStarted = true
		dataLock.Unlock()
		internal_log("invoking new replica")
	}
	opts := grpc.WithInsecure()
	cc, err := grpc.Dial(depControllerAddr, opts)
	if err != nil {
		fmt.Println(err)
	}
	defer cc.Close()
	client := pb.NewDeploymentServiceClient(cc)
	rn := int32(replicaNumber)
	request := &pb.DeploymentServiceRequest{Name: func_name, FunctionCall: type_call, ReplicaNumber: rn}
	depLock.Lock()
	// fmt.Println("depLock1")
	resp, err := client.Deployment(context.Background(), request)
	depLock.Unlock()
	if err != nil {
		fmt.Println(err)
		return err
	}
	if func_found {
		config, err := rest.InClusterConfig()
		if err != nil {
			fmt.Println("Failed to get in-cluster config: " + err.Error())
		}
		k, err := kubernetes.NewForConfig(config)
		if err != nil {
			fmt.Println("Failed to create kclient: " + err.Error())
		}
		labels_to_delete := resp.Message
		if labels_to_delete != "" {
			dataLock.Lock()
			fmt.Print("deleting.....")
			max_concurrency = 2 * max_concurrency
			dataLock.Unlock()
			log.Printf("Increasing the concurrency to %d", max_concurrency)
			ingresses, _ := k.NetworkingV1().Ingresses("macropod-functions").List(context.Background(), metav1.ListOptions{LabelSelector: labels_to_delete})
			for _, ingress := range ingresses.Items{
				deleteIngressList = append(deleteIngressList, &ingress)
			}
			go deleteDeployments(labels_to_delete, func_name)
		}

	}
	if !func_found {
		if debug > 2 {
			log.Printf("waiting for a new replica %s to be up ", resp.Message)
		}
		replicaName := func_name + "-" + strconv.Itoa(replicaNumber)
		dataLock.Lock()
		for !slices.Contains(replicaRunning, replicaName) {
			dataLock.Unlock()
			time.Sleep(10 * time.Millisecond)
			dataLock.Lock()
		}
		deploymentStarted = false
		dataLock.Unlock()
	}

	return nil
}

func checkTTL() {
	internal_log("Checking TTL")
	config, err := rest.InClusterConfig()
	if err != nil {
		fmt.Println("Failed to get in-cluster config: " + err.Error())
	}
	k, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Println("Failed to create kclient: " + err.Error())
	}
	for {
		currentTime := time.Now()
		for name, timestamp := range serviceTimeStamp {
			elapsedTime := currentTime.Sub(timestamp)
			if elapsedTime.Seconds() > float64(ttl_seconds) {
				service_name := strings.Split(name, ".")[0]
				if debug > 1 {
					log.Printf("deleting because of TTL %s", service_name)
					log.Printf("Deleting service and deployment of %s because of TTL\n", service_name)
				}
				service, exists := k.CoreV1().Services("macropod-functions").Get(context.Background(), service_name, metav1.GetOptions{})
				if exists == nil {
					dataLock.Lock()
					// fmt.Println("dataLock5")
					labels := service.Labels["workflow_replica"]
					dataLock.Unlock()
					labels_replica := "workflow_replica=" + labels
					services, err := k.CoreV1().Services("macropod-functions").List(context.Background(), metav1.ListOptions{LabelSelector: labels_replica})
					if err != nil {
						fmt.Println(err)
					}
					for _, service := range services.Items {
						dataLock.Lock()
						// fmt.Println("dataLock6")
						delete(serviceTimeStamp, service.Name)
						delete(channelOfService, service.Name)
						delete(stubsOfChannel, service.Name)
						dataLock.Unlock()
						fmt.Println(service.ObjectMeta.Name)
						k.CoreV1().Services("macropod-functions").Delete(context.Background(), service.ObjectMeta.Name, metav1.DeleteOptions{})
					}
					deployments, err := k.AppsV1().Deployments("macropod-functions").List(context.Background(), metav1.ListOptions{LabelSelector: labels_replica})
					if err != nil {
						fmt.Println(err)
					}
					for _, deployment := range deployments.Items {
						fmt.Println(deployment.ObjectMeta.Name)
						k.AppsV1().Deployments("macropod-functions").Delete(context.Background(), deployment.ObjectMeta.Name, metav1.DeleteOptions{})
					}
					ingresses, err := k.NetworkingV1().Ingresses("macropod-functions").List(context.Background(), metav1.ListOptions{LabelSelector: labels_replica})
					if err != nil {
						fmt.Println(err)
					}
					for _, ingress := range ingresses.Items {
						fmt.Println(ingress.ObjectMeta.Name)
						k.NetworkingV1().Ingresses("macropod-functions").Delete(context.Background(), ingress.ObjectMeta.Name, metav1.DeleteOptions{})
					}
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
				// fmt.Println("dataLock7")
				func_name := ingress.Labels["workflow_name"]
				dataLock.Unlock()
				hostname := fmt.Sprintf("%s.%s.svc.cluster.local:%d", serviceName, namespace, port)
				dataLock.Lock()
				// fmt.Println("dataLock8")
				replica_name := ingress.Labels["workflow_replica"]
				dataLock.Unlock()
				for !ifPodsAreRunning(replica_name, namespace) {
				}
				internal_log("service found: " + hostname)
				channel, _ := grpc.Dial(hostname, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*200), grpc.MaxCallSendMsgSize(1024*1024*200)), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultServiceConfig(retrypolicy))
				stub := wf_pb.NewGRPCFunctionClient(channel)
				for channel.GetState() != connectivity.Ready {
					time.Sleep(10 * time.Millisecond)
				}
				dataLock.Lock()
				serviceCount[hostname] = 0
				replicaRunning = append(replicaRunning, replica_name)
				serviceTimeStamp[hostname] = time.Now()
				serviceCreationTime[hostname] = time.Now()
				hostTargets[func_name] = append(hostTargets[func_name], hostname)
				concurrencyList[hostname] = max_concurrency
				channelOfService[hostname] = append(channelOfService[hostname], channel)
				stubsOfChannel[hostname] = append(stubsOfChannel[hostname], stub)
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
				// fmt.Println("dataLock11")
				func_name = ingress.Labels["workflow_name"]
				dataLock.Unlock()
				hostname_deleted = fmt.Sprintf("%s.%s.svc.cluster.local:%d", serviceName, namespace, port)
			}
		}
	}
	log.Printf("deleting %s\n", hostname_deleted)

	log.Printf("hosttargets %v", hostTargets[func_name])
	for i, val := range hostTargets[func_name] {
		if val == hostname_deleted {
			dataLock.Lock()
			// fmt.Println("dataLock13")
			hostTargets[func_name] = append(hostTargets[func_name][:i], hostTargets[func_name][i+1:]...)
			dataLock.Unlock()
			if debug > 2 {
				log.Printf("deleting%v", hostname_deleted)
			}
			if _, exists := serviceCount[hostname_deleted]; exists {
				dataLock.Lock()
				// fmt.Println("dataLock14")
				delete(serviceCount, hostname_deleted)
				delete(serviceTimeStamp, hostname_deleted)
				delete(channelOfService, hostname_deleted)
				delete(stubsOfChannel, hostname_deleted)
				workflow_deployments[func_name]--
				dataLock.Unlock()
			}
			break
		}

	}

}

func watchIngress(kclient *kubernetes.Clientset) {
	//reference: https://blog.mimacom.com/k8s-watch-resources/
	watchFunc := func(options metav1.ListOptions) (watch.Interface, error) {
		return kclient.NetworkingV1().Ingresses("").Watch(context.Background(), metav1.ListOptions{})
	}
	watcher, _ := toolsWatch.NewRetryWatcher("1", &cache.ListWatch{WatchFunc: watchFunc})
	for event := range watcher.ResultChan() {
		ingress, ok := event.Object.(*networkingv1.Ingress)
		if !ok {
			continue
		}
		switch event.Type {
		case watch.Added, watch.Modified:
			log.Printf("Found a new ingress")
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
	if len(hostTargets[func_name]) > 0 {
		idx := invocationsCount[func_name] % len(hostTargets[func_name])
		service := hostTargets[func_name][idx]
		if serviceCount[service] < concurrencyList[service] && len(channelOfService[service]) > 0 && channelOfService[service][0].GetState() == connectivity.Ready {
			target = service
			serviceCount[service]++
			serviceTimeStamp[service] = time.Now()
		}
	}
	if target == "" && !triggered && !runningDeploymentController[func_name] && workflow_invocations[func_name] > currentReplicasThatCanBeHandled[func_name] {
		runningDeploymentController[func_name] = true
		triggered = true
		replicaCount[func_name] = replicaCount[func_name] + 1
		currentReplicasThatCanBeHandled[func_name] += max_concurrency
		go callDepController(false, func_name, replicaCount[func_name])
		workflow_deployments[func_name]++
		runningDeploymentController[func_name] = false
	}
	dataLock.Unlock()
	return target, triggered
}

func Serve_WF_Invoke(res http.ResponseWriter, req *http.Request) {
	func_name := req.PathValue("func_name")
	dataLock.Lock()
	// fmt.Println("dataLock23")
	workflow_invocations[func_name]++
	invocationsCount[func_name]++
	_, exists := runningDeploymentController[func_name]
	if !exists {
		runningDeploymentController[func_name] = false
	}
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
	go callDepController(true, func_name, max_concurrency)
	request_type := "gg"
	dataLock.Lock()

	dataLock.Unlock()
	stub := stubsOfChannel[target][0]
	start_time := time.Now()
	response, err := stub.GRPCFunctionHandler(context.Background(), &wf_pb.RequestBody{Data: payload, WorkflowId: workflow_id, Depth: 0, Width: 0, RequestType: &request_type})
	end_time := time.Now()
	status = response.GetCode()
	if status == 200 {
		fmt.Println("request was served by: " + target + " in " + end_time.Sub(start_time).String() + " time to find target " + start_time.Sub(look_arget).String())
		fmt.Fprint(res, response.GetReply())

	} else {
		log.Printf("Non 200 code %s: %d, error : %s in : %s time ti find target : %s", target, response.GetCode(), err.Error(), end_time.Sub(start_time).String(), start_time.Sub(look_arget).String())
		http.Error(res, err.Error(), http.StatusBadGateway)
	}
	dataLock.Lock()
	serviceCount[target]--
	workflow_invocations[func_name]--
	dataLock.Unlock()
}

func Serve_WF_Create(res http.ResponseWriter, req *http.Request) {
	internal_log("WF_CREATE_START " + req.PathValue("func_name"))
	func_name := req.PathValue("func_name")
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
	max_concurrency = 5 // we will need to fix this at some point - AE
	client := pb.NewDeploymentServiceClient(cc)
	request := &pb.DeploymentServiceRequest{Name: func_name, FunctionCall: "create", Workflow: &workflow}
	internal_log("requesting create for " + func_name)
	_, err = client.Deployment(context.Background(), request)
	internal_log("returned create for " + func_name)
	if err != nil {
		internal_log("create return - " + err.Error())
	}
	dataLock.Lock()
	// fmt.Println("dataLock26")
	workflows[func_name] = body_u
	runningDeploymentController[func_name] = false
	workflow_invocations[func_name] = 0
	workflow_deployments[func_name] = 0
	invocationsCount[func_name] = 0
	overallWorkflowTTL[func_name] = time.Now()
	currentReplicasThatCanBeHandled[func_name] = 0
	dataLock.Unlock()
	internal_log("WF_CREATE_END " + func_name)
	fmt.Fprintf(res, "Workflow created successfully. Invoke your workflow with /invoke/"+func_name+"\n")
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
	_, _ = client.Deployment(context.Background(), request)
	dataLock.Lock()
	workflows[req.PathValue("func_name")] = body_u
	runningDeploymentController[req.PathValue("func_name")] = false
	workflow_invocations[req.PathValue("func_name")] = 0
	workflow_deployments[req.PathValue("func_name")] = 0
	invocationsCount[req.PathValue("func_name")] = 0
	currentReplicasThatCanBeHandled[req.PathValue("func_name")] = 0
	dataLock.Unlock()
	internal_log("WF_UPDATE_END " + req.PathValue("func_name"))
	fmt.Fprintf(res, "Workflow \""+req.PathValue("func_name")+"\" has been updated successfully.\n")
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
	label_workflow := "workflow_name=" + req.PathValue("func_name")
	config, err := rest.InClusterConfig()
	if err != nil {
		fmt.Println("Failed to get in-cluster config: " + err.Error())
	}
	k, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Println("Failed to create kclient: " + err.Error())
	}
	services, err := k.CoreV1().Services("macropod-functions").List(context.Background(), metav1.ListOptions{LabelSelector: label_workflow})
	if err != nil {
		fmt.Println(err)
	}
	for _, service := range services.Items {
		dataLock.Lock()
		fmt.Println("deleting service: " + service.Name)
		// fmt.Println("dataLock28")
		delete(serviceTimeStamp, service.Name)
		delete(channelOfService, service.Name)
		delete(stubsOfChannel, service.Name)
		dataLock.Unlock()
		fmt.Println(service.ObjectMeta.Name)
		k.CoreV1().Services("macropod-functions").Delete(context.Background(), service.ObjectMeta.Name, metav1.DeleteOptions{})
	}
	deployments, err := k.AppsV1().Deployments("macropod-functions").List(context.Background(), metav1.ListOptions{LabelSelector: label_workflow})
	if err != nil {
		fmt.Println(err)
	}
	for _, deployment := range deployments.Items {
		fmt.Println(deployment.ObjectMeta.Name)
		k.AppsV1().Deployments("macropod-functions").Delete(context.Background(), deployment.ObjectMeta.Name, metav1.DeleteOptions{})
	}
	ingresses, err := k.NetworkingV1().Ingresses("macropod-functions").List(context.Background(), metav1.ListOptions{LabelSelector: label_workflow})
	if err != nil {
		fmt.Println(err)
	}
	for _, ingress := range ingresses.Items {
		fmt.Println(ingress.ObjectMeta.Name)
		k.NetworkingV1().Ingresses("macropod-functions").Delete(context.Background(), ingress.ObjectMeta.Name, metav1.DeleteOptions{})
	}

	for {
		deployments_list, _ := k.CoreV1().Pods("macropod-functions").List(context.Background(), metav1.ListOptions{LabelSelector: label_workflow})
		if deployments_list == nil || len(deployments_list.Items) == 0 {
			log.Printf("Pods for function %s are deleted", req.PathValue("func_name"))
			break
		}
		time.Sleep(10 * time.Millisecond) // let all the deployments be deleted before new ones
	}

	dataLock.Lock()
	// fmt.Println("dataLock30")
	delete(workflows, req.PathValue("func_name"))
	delete(runningDeploymentController, req.PathValue("func_name"))
	delete(workflow_invocations, req.PathValue("func_name"))
	delete(workflow_deployments, req.PathValue("func_name"))
	delete(replicaCount, req.PathValue("func_name"))
	delete(overallWorkflowTTL, req.PathValue("func_name"))
	delete(hostTargets, req.PathValue("func_name"))
	dataLock.Unlock()
	internal_log("WF_DELETE_END " + req.PathValue("func_name"))
	fmt.Fprintf(res, "Workflow \""+req.PathValue("func_name")+"\" has been deleted successfully.\n")
}

func main() {
	log.Print("Ingress controller started")
	ttl_seconds, _ = strconv.Atoi(os.Getenv("TTL"))
	overall_ttl = ttl_seconds * 60 //VD : need to discuss the value of this
	config, err := rest.InClusterConfig()
	if err != nil {
		fmt.Println("Failed to get in-cluster config: " + err.Error())
	}
	debug, err = strconv.Atoi(os.Getenv("DEBUG"))
	if err != nil {
		debug = 0
	}
	kclient, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Println("Failed to create kclient: " + err.Error())
	}
	go checkTTL()
	//go checkOverallTTL() // disabled for the moment; seems buggy - AE
	workflows = make(map[string]Workflow)
	max_concurrency = 5
	maxConcurrencyStanby = 6
	go watchIngress(kclient)
	h := http.NewServeMux()
	h.HandleFunc("/", Serve_Help)
	h.HandleFunc("/invoke/{func_name}", Serve_WF_Invoke)
	h.HandleFunc("/create/{func_name}", Serve_WF_Create)
	h.HandleFunc("/update/{func_name}", Serve_WF_Update)
	h.HandleFunc("/delete/{func_name}", Serve_WF_Delete)
	http.ListenAndServe(":"+os.Getenv("SERVICE_PORT"), h)
}
