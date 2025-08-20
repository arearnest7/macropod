package main

import (
	pb "app/macropod_pb"
	structpb "google.golang.org/protobuf/types/known/structpb"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	toolsWatch "k8s.io/client-go/tools/watch"

	"net"
	"net/http"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

type IngressService struct {
	pb.UnimplementedMacroPodIngressServer
}

var (
	workflows                    = make(map[string]*pb.WorkflowStruct)
	workflow_target_concurrency  = make(map[string]int32)
	workflow_invocations_current = make(map[string]int)
	workflow_invocations_total   = make(map[string]int)
	entry_function               = make(map[string]string)
	workflow_functions_created   = make(map[string]int)
	service_target               = make(map[string][]string)
	service_channel              = make(map[string]*grpc.ClientConn)
	service_stub                 = make(map[string]pb.MacroPodFunctionClient)
	service_count                = make(map[string]int)
	service_timestamp            = make(map[string]time.Time)
	service_ttl                  = make(map[string]int32)
	service_downstream_channel   = make(map[string]map[string]*grpc.ClientConn)
        service_downstream_stub      = make(map[string]map[string]pb.MacroPodFunctionClient)

	deployer_channel        *grpc.ClientConn
	deployer_stub           pb.MacroPodDeployerClient
	deployer_update_stub    pb.MacroPodDeployerClient
	deployer_update_channel *grpc.ClientConn

	default_config   *pb.ConfigStruct
	deployer_address string
	updateRun        bool
	dataLock         sync.Mutex
	createRun        bool

	retrypolicy = `{
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

func Debug(message string, debug_level int) {
	if default_config.GetDebug() >= int32(debug_level) {
		fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.000000 UTC") + ": " + message + "\n")
	}
}

func Deployer_Check() {
	var err error
	for deployer_channel == nil {
		deployer_channel, _ = grpc.Dial(deployer_address, grpc.WithInsecure())
		deployer_stub = pb.NewMacroPodDeployerClient(deployer_channel)
		Debug("attempting rebuild deployer stub...", 5)
		time.Sleep(10 * time.Millisecond)
	}
	for deployer_channel.GetState() != connectivity.Ready {
		if deployer_channel.GetState() == connectivity.Connecting {
			Debug("waiting for deployer stub to be finish connecting", 5)
			time.Sleep(10 * time.Millisecond)
		} else if deployer_channel.GetState() == connectivity.TransientFailure || deployer_channel.GetState() == connectivity.Shutdown {
			deployer_channel.Close()
			deployer_channel, err = grpc.Dial(deployer_address, grpc.WithInsecure())
			if err != nil {
				Debug(err.Error(), 0)
			}
			deployer_stub = pb.NewMacroPodDeployerClient(deployer_channel)
			Debug("attempting rebuild deployer stub due to transient failure or shutdown...", 5)
			time.Sleep(10 * time.Millisecond)
		} else {
			break
		}
	}
}

// Create a separate channel for each area of high load in the application (as per the documentation)
func Deployer_Update_Check() {
	var err error
	for deployer_update_channel == nil {
		deployer_update_channel, _ = grpc.Dial(deployer_address, grpc.WithInsecure())
		deployer_update_stub = pb.NewMacroPodDeployerClient(deployer_channel)
		Debug("attempting rebuild deployer update stub...", 5)
		time.Sleep(10 * time.Millisecond)
	}

	for deployer_update_channel.GetState() != connectivity.Ready {
		if deployer_update_channel.GetState() == connectivity.Connecting {
			Debug("waiting for deployer update stub to be finish connecting", 5)
			time.Sleep(10 * time.Millisecond)
		} else if deployer_update_channel.GetState() == connectivity.TransientFailure || deployer_update_channel.GetState() == connectivity.Shutdown {
			deployer_update_channel.Close()
			deployer_update_channel, err = grpc.Dial(deployer_address, grpc.WithInsecure())
			if err != nil {
				Debug(err.Error(), 0)
			}
			deployer_update_stub = pb.NewMacroPodDeployerClient(deployer_update_channel)
			Debug("attempting rebuild deployer update stub due to transient failure or shutdown...", 5)
			time.Sleep(10 * time.Millisecond)
		} else {
			break
		}
	}
}

func IfPodsAreRunning(workflow_replica string, namespace string) bool {
	label_replica := "workflow_replica=" + workflow_replica
	config, err := rest.InClusterConfig()
	if err != nil {
		Debug("Failed to get in-cluster config: "+err.Error(), 0)
		return false
	}
	k, err := kubernetes.NewForConfig(config)
	if err != nil {
		Debug("Failed to create k: "+err.Error(), 0)
	}
	pods, err := k.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{LabelSelector: label_replica})
	if err != nil {
		Debug(err.Error(), 0)
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

func WatchTTL() {
	Debug("watch TTL", 3)
	config, err := rest.InClusterConfig()
	if err != nil {
		Debug("Failed to get in-cluster config: "+err.Error(), 0)
	}
	k, err := kubernetes.NewForConfig(config)
	if err != nil {
		Debug("Failed to create k: "+err.Error(), 0)
	}
	for {
		currentTime := time.Now()
		dataLock.Lock()
		var service_timestamp_check []string
		for name, _ := range service_timestamp {
			service_timestamp_check = append(service_timestamp_check, name)
		}
		dataLock.Unlock()
		var deleted_services []string
		for _, name := range service_timestamp_check {
			dataLock.Lock()
			timestamp := service_timestamp[name]
			elapsedTime := currentTime.Sub(timestamp)
			cnt := service_count[name]
			ttl := float64(service_ttl[name])
			dataLock.Unlock()
			if elapsedTime.Seconds() > ttl && cnt == 0 {
				service_name := strings.Split(name, ".")[0]
				namespace := strings.Split(name, ".")[1]
				Debug("Deleting service and deployment of "+service_name+" because of TTL: "+currentTime.UTC().Format("2006-01-02 15:04:05.000000 UTC")+" - "+timestamp.UTC().Format("2006-01-02 15:04:05.000000 UTC")+" > "+strconv.Itoa(int(ttl)), 1)
				service, err := k.CoreV1().Services(namespace).Get(context.Background(), service_name, metav1.GetOptions{})
				if service != nil && err == nil {
					dataLock.Lock()
					workflow_name := service.Labels["workflow_name"]
					labels := service.Labels["workflow_replica"]
					dataLock.Unlock()
					//labels_replica := "workflow_replica=" + labels
					dataLock.Lock()
					for i, val := range service_target[workflow_name] {
						if val == name {
							service_target[workflow_name] = append(service_target[workflow_name][:i], service_target[workflow_name][i+1:]...)
							delete(service_count, name)
							deleted_services = append(deleted_services, name)
							_, channel_exists := service_channel[name]
							if channel_exists {
								service_channel[name].Close()
							}
							_, exists := workflow_functions_created[workflow_name]
							if exists {
								workflow_functions_created[workflow_name]--
							}
							delete(service_channel, name)
							delete(service_stub, name)
							delete(service_ttl, name)
						}
					}
					dataLock.Unlock()
					for {
						if labels == "" {
							break
						}
						ttl_request := pb.MacroPodRequest{Workflow: &workflow_name, Target: &labels}
						Debug("deleting resources for"+labels, 5)
						go deployer_stub.TTLDelete(context.Background(), &ttl_request)
						time.Sleep(100 * time.Millisecond)
						deployments_list, err := k.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{LabelSelector: labels})
						if err != nil || deployments_list == nil || len(deployments_list.Items) == 0 {
							break
						}
					}
				} else {
					dataLock.Lock()
					for workflow_name_val, _ := range service_target {
						for i, val := range service_target[workflow_name_val] {
							if val == name {
								service_target[workflow_name_val] = append(service_target[workflow_name_val][:i], service_target[workflow_name_val][i+1:]...)
								delete(service_count, name)
								deleted_services = append(deleted_services, name)
								_, channel_exists := service_channel[name]
								if channel_exists {
									service_channel[name].Close()
								}
								delete(service_channel, name)
								delete(service_stub, name)
								delete(service_ttl, name)
							}
						}
						dataLock.Unlock()
					}
				}
			}
		}
		dataLock.Lock()
		for _, sname := range deleted_services {
			delete(service_timestamp, sname)
		}
		dataLock.Unlock()
		time.Sleep(time.Second)
	}
}

func UpdateHostTargets(ingress *networkingv1.Ingress) {
	for _, rule := range ingress.Spec.Rules {
		if rule.HTTP != nil {
			for _, path := range rule.HTTP.Paths {
				serviceName := path.Backend.Service.Name
				namespace := ingress.Namespace
				port := path.Backend.Service.Port.Number
				dataLock.Lock()
				workflow_name := ingress.Labels["workflow_name"]
				hostname := fmt.Sprintf("%s.%s.svc.cluster.local:%d", serviceName, namespace, port)
				replica_name := ingress.Labels["workflow_replica"]
				dataLock.Unlock()
				for !IfPodsAreRunning(replica_name, namespace) {
				}
				//Debug("service found: " + hostname, 2)
				channel, _ := grpc.Dial(hostname, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*200), grpc.MaxCallSendMsgSize(1024*1024*200)), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultServiceConfig(retrypolicy))
				stub := pb.NewMacroPodFunctionClient(channel)
				for channel.GetState() != connectivity.Ready {
					time.Sleep(10 * time.Millisecond)
				}
				dataLock.Lock()

				_, exists := workflows[workflow_name]
				if exists {
					Debug("service found: "+hostname, 2)
					service_count[hostname] = 0
					service_timestamp[hostname] = time.Now()
					service_target[workflow_name] = append(service_target[workflow_name], hostname)
					service_channel[hostname] = channel
					service_stub[hostname] = stub
					service_ttl[hostname] = default_config.GetTTL()
					if workflows[workflow_name].GetConfig() != nil && workflows[workflow_name].GetConfig().GetTTL() != 0 {
						service_ttl[hostname] = workflows[workflow_name].GetConfig().GetTTL()
					}
				}
				//fmt.Printf("INGRESS: service target %v\n", service_target)
				dataLock.Unlock()
			}
		}
	}
}

func DeleteHostTargets(ingress *networkingv1.Ingress) {
	for _, rule := range ingress.Spec.Rules {
		if rule.HTTP != nil {
			for _, path := range rule.HTTP.Paths {
				serviceName := path.Backend.Service.Name
				namespace := ingress.Namespace
				port := path.Backend.Service.Port.Number
				workflow_name := ingress.Labels["workflow_name"]
				hostname_deleted := fmt.Sprintf("%s.%s.svc.cluster.local:%d", serviceName, namespace, port)
				dataLock.Lock()
				for i, val := range service_target[workflow_name] {
					if val == hostname_deleted {
						service_target[workflow_name] = append(service_target[workflow_name][:i], service_target[workflow_name][i+1:]...)
						Debug("deleting "+hostname_deleted, 2)
						if _, exists := service_count[hostname_deleted]; exists {
							delete(service_count, hostname_deleted)
							delete(service_timestamp, hostname_deleted)
							_, channel_exists := service_channel[hostname_deleted]
							if channel_exists {
								service_channel[hostname_deleted].Close()
							}
							_, exists := workflow_functions_created[workflow_name]
							if exists {
								workflow_functions_created[workflow_name]--
							}
							delete(service_channel, hostname_deleted)
							delete(service_stub, hostname_deleted)
							delete(service_ttl, hostname_deleted)
						}
						break
					}
				}
				dataLock.Unlock()
			}
		}
	}
}

func WatchIngress() {
	config, err := rest.InClusterConfig()
	if err != nil {
		Debug("Failed to get in-cluster config: "+err.Error(), 0)
	}
	k, err := kubernetes.NewForConfig(config)
	if err != nil {
		Debug("Failed to create k: "+err.Error(), 0)
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
			UpdateHostTargets(ingress)
		case watch.Deleted:
			DeleteHostTargets(ingress)
		}
	}
}

func runCreateFunction(workflow_name_for_create string) {
	Debug("create request to deployer", 6)
	dataLock.Lock()
	inv_cnt := workflow_invocations_current[workflow_name_for_create]
	dep_cnt := workflow_functions_created[workflow_name_for_create]
	dataLock.Unlock()
	fmt.Printf("number of invocations: %d, number of deployments: %d\n", inv_cnt, dep_cnt)
	create_deployment_request := pb.MacroPodRequest{Workflow: &workflow_name_for_create}
	response, _ := deployer_stub.CreateDeployment(context.Background(), &create_deployment_request)
	Debug("response from create deployment: "+response.GetReply(), 5)
	if response.GetReply() == "0" {
		dataLock.Lock()
		workflow_functions_created[workflow_name_for_create]++
		dataLock.Unlock()
	}
	Debug("set  createRun to false", 6)
	return
}
func GetTarget(triggered bool, workflow_name string, target string) (pb.MacroPodFunctionClient, string, bool) {
	if len(service_target[workflow_name]) > 0 {
		idx := workflow_invocations_total[workflow_name] % len(service_target[workflow_name])
		service := service_target[workflow_name][idx]
		target_concurrency, target_concurrency_set := workflow_target_concurrency[workflow_name]
		channel, service_channel_exists := service_channel[service]
		if service_channel_exists && channel.GetState() == connectivity.Ready {
			if !target_concurrency_set || int(target_concurrency) == -1 || service_count[service] < int(target_concurrency) {
				target = service
				service_count[service]++
				service_timestamp[service] = time.Now()
			}
		}
	}
	if !triggered && workflow_functions_created[workflow_name] < workflow_invocations_current[workflow_name] {
		triggered = true
		//create_deployment_request := pb.MacroPodRequest{Workflow: &workflow_name}
		//Debug("sending request to create deployments\n", 0)
		go runCreateFunction(workflow_name)
	}
	stub_to_return := service_stub[target]

	return stub_to_return, target, triggered
}

func Serve_Config(request *pb.ConfigStruct) string {
	dataLock.Lock()
	default_config = request
	if default_config.GetNamespace() == "" {
		namespace := "macropod-functions"
		default_config.Namespace = &namespace
	}
	if default_config.GetTTL() == 0 {
		ttl := int32(180)
		default_config.TTL = &ttl
	}
	if default_config.GetDeployment() == "" {
		deployment := "macropod"
		default_config.Deployment = &deployment
	}
	if default_config.GetCommunication() == "" {
		communication := "direct"
		default_config.Communication = &communication
	}
	if default_config.GetAggregation() == "" {
		aggregation := "agg"
		default_config.Aggregation = &aggregation
	}
	if default_config.GetTargetConcurrency() == 0 {
		target := int32(-1)
		default_config.TargetConcurrency = &target
	}
	dataLock.Unlock()
	_, err := deployer_stub.Config(context.Background(), request)
	if err != nil {
		Debug(err.Error(), 0)
		return ""
	}
	dataLock.Lock()
	config_txt := ""
	config_txt += "Namespace: " + default_config.GetNamespace() + "\n"
	config_txt += "TTL: " + string(default_config.GetTTL()) + "\n"
	config_txt += "Deployment: " + default_config.GetDeployment() + "\n"
	config_txt += "Communication: " + default_config.GetCommunication() + "\n"
	config_txt += "Aggregation: " + default_config.GetAggregation() + "\n"
	config_txt += "TargetConcurrency: " + string(default_config.GetTargetConcurrency()) + "\n"
	config_txt += "Debug: " + string(default_config.GetDebug()) + "\n"
	dataLock.Unlock()
	return config_txt
}
func runUpdateFunction(update_request pb.MacroPodRequest) {
	Debug("update request to deployer", 6)
	updateRun = true
	deployer_update_stub.UpdateDeployments(context.Background(), &update_request)
	updateRun = false
	Debug("set updateRun to false", 6)
	return
}
func Serve_WorkflowInvoke(request *pb.MacroPodRequest) (string, int32) {
	workflow_name := request.GetWorkflow()
	function_name := entry_function[workflow_name]
	request.Function = &function_name
	var stub pb.MacroPodFunctionClient
	triggered := false
	var invocations_current string
	invoked := false
	target := ""
	for target == "" {
		dataLock.Lock()
		if !invoked {
			workflow_invocations_current[workflow_name]++
			invocations_current = strconv.Itoa(workflow_invocations_current[workflow_name])
			workflow_invocations_total[workflow_name]++
			invoked = true
		}
		stub, target, triggered = GetTarget(triggered, workflow_name, target)
		dataLock.Unlock()
	}
	Debug("target set to  "+target, 6)
	workflow_id := strconv.Itoa(rand.Intn(100000))
	dw := int32(0)
	request.WorkflowID = &workflow_id
	request.Depth = &dw
	request.Width = &dw
	status := int32(0)
	var response *pb.MacroPodReply
	update_request := pb.MacroPodRequest{Workflow: &workflow_name, WorkflowID: &invocations_current}
	response, err := stub.Invoke(context.Background(), request)
	if err != nil {
		Debug("error while request "+err.Error(), 3)
	}
	status = response.GetCode()
	dataLock.Lock()
	if !updateRun {
		go runUpdateFunction(update_request)
	}
	service_count[target]--
	workflow_invocations_current[workflow_name]--
	dataLock.Unlock()
	Debug("workflow invoke request was served by: "+target, 6)
	return response.GetReply(), status
}

func Serve_FunctionInvoke(request *pb.MacroPodRequest) (string, int32) {
	Debug("got invoke request", 3)
        dataLock.Lock()
        stub, exists := service_downstream_stub[request.GetWorkflow()][request.GetTarget()]
        if !exists {
            channel, _ := grpc.Dial(request.GetTarget(), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*200), grpc.MaxCallSendMsgSize(1024*1024*200)), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultServiceConfig(retrypolicy))
	    stub = pb.NewMacroPodFunctionClient(channel)
            service_downstream_channel[request.GetWorkflow()][request.GetTarget()] = channel
            service_downstream_stub[request.GetWorkflow()][request.GetTarget()] = stub
        }
        dataLock.Unlock()
	response, _ := stub.Invoke(context.Background(), request)
	status := response.GetCode()
	Debug("function invoke request was served by: "+request.GetTarget(), 6)
	return response.GetReply(), status
}

func Serve_CreateWorkflow(request *pb.WorkflowStruct) string {
	if request.GetName() == "" {
		return "Workflow definition is malformed\n"
	}
	Debug("WF_CREATE "+request.GetName(), 2)
	_, exists := workflows[request.GetName()]
	if exists {
		Debug("WF_CREATE_END_EXISTS "+request.GetName(), 2)
		return "Workflow already exists\n"
	}
	dataLock.Lock()
	fmt.Printf("deployer channel status : %v\n", deployer_channel.GetState())
	workflows[request.GetName()] = request
	workflow_target_concurrency[request.GetName()] = default_config.GetTargetConcurrency()
	if request.GetConfig() != nil && request.GetConfig().GetTargetConcurrency() != 0 {
		workflow_target_concurrency[request.GetName()] = request.GetConfig().GetTargetConcurrency()
	}
        service_downstream_channel[request.GetName()] = make(map[string]*grpc.ClientConn)
        service_downstream_stub[request.GetName()] = make(map[string]pb.MacroPodFunctionClient)
	workflow_invocations_current[request.GetName()] = 0
	workflow_functions_created[request.GetName()] = 0
	workflow_invocations_total[request.GetName()] = 0
	Deployer_Check()
	Deployer_Update_Check()
	dataLock.Unlock()
	results, _ := deployer_stub.CreateWorkflow(context.Background(), request)
	entrypoint := results.GetReply()
	dataLock.Lock()
	entry_function[request.GetName()] = entrypoint
	dataLock.Unlock()
	Debug("WF_CREATE_END "+request.GetName(), 2)
	return "Workflow created successfully. Invoke your workflow with /workflow/invoke/" + request.GetName() + "\n"
}

func Serve_UpdateWorkflow(request *pb.WorkflowStruct) string {
	if request.GetName() == "" {
		return "Workflow definition is malformed\n"
	}
	Debug("WF_UPDATE "+request.GetName(), 2)
	delete_request := pb.MacroPodRequest{Workflow: &request.Name}
	Serve_DeleteWorkflow(&delete_request)
	Serve_CreateWorkflow(request)
	Debug("WF_UPDATE_END"+request.GetName(), 2)
	return "Workflow " + request.GetName() + " has been updated.\n"
}

func Serve_DeleteWorkflow(request *pb.MacroPodRequest) string {
	if request.GetWorkflow() == "" {
		return "Workflow name is missing\n"
	}
	Debug("WF_DELETE "+request.GetWorkflow(), 2)
	label_workflow := "workflow_name=" + request.GetWorkflow()
	fmt.Printf("delete request for %s", label_workflow)
	config, err := rest.InClusterConfig()
	if err != nil {
		Debug("Failed to get in-cluster config: "+err.Error(), 0)
	}
	k, err := kubernetes.NewForConfig(config)
	if err != nil {
		Debug("Failed to create k: "+err.Error(), 0)
	}
	namespace := default_config.GetNamespace()
	dataLock.Lock()
	if workflows[request.GetWorkflow()].GetConfig() != nil && workflows[request.GetWorkflow()].GetConfig().GetNamespace() != "" {
		namespace = workflows[request.GetWorkflow()].GetConfig().GetNamespace()
	}
	dataLock.Unlock()
	ingresses, err := k.NetworkingV1().Ingresses(namespace).List(context.Background(), metav1.ListOptions{LabelSelector: label_workflow})
	if err != nil {
		Debug(err.Error(), 0)
	}
	for _, ingress := range ingresses.Items {
		for _, rule := range ingress.Spec.Rules {
			if rule.HTTP != nil {
				for _, path := range rule.HTTP.Paths {
					serviceName := path.Backend.Service.Name
					namespace := ingress.Namespace
					port := path.Backend.Service.Port.Number

					hostname_deleted := fmt.Sprintf("%s.%s.svc.cluster.local:%d", serviceName, namespace, port)
					Debug("DELETE: "+hostname_deleted+" is deleted\n", 3)
					dataLock.Lock()
					delete(service_count, hostname_deleted)
					delete(service_timestamp, hostname_deleted)
					_, channel_exists := service_channel[hostname_deleted]
					if channel_exists {
						service_channel[hostname_deleted].Close()
					}
					delete(service_channel, hostname_deleted)
					delete(service_stub, hostname_deleted)
					delete(service_ttl, hostname_deleted)
					dataLock.Unlock()
				}

			}
		}

	}
	response, err := deployer_stub.DeleteWorkflow(context.Background(), request)
	if err != nil {
		Debug("error in response of delete "+err.Error(), 5)
	}
	fmt.Printf("response: %v", response)
	dataLock.Lock()
	delete(service_target, request.GetWorkflow())
	delete(workflows, request.GetWorkflow())
	delete(workflow_target_concurrency, request.GetWorkflow())
	delete(workflow_invocations_current, request.GetWorkflow())
	delete(workflow_functions_created, request.GetWorkflow())
	delete(workflow_invocations_total, request.GetWorkflow())
	delete(entry_function, request.GetWorkflow())
        for _, channel := range service_downstream_channel[request.GetWorkflow()] {
            if channel != nil {
                channel.Close()
            }
        }
        delete(service_downstream_channel, request.GetWorkflow())
        delete(service_downstream_stub, request.GetWorkflow())
	Debug("deleting and recreating the deployer channel", 5)
	deployer_channel.Close()
	deployer_update_channel.Close()
	Deployer_Check()
	Deployer_Update_Check()
	fmt.Printf("service target %v\n service_count :%v\n", service_target, service_count)
	dataLock.Unlock()
	Debug("WF_DELETE_END "+request.GetWorkflow(), 2)
	return "Workflow \"" + request.GetWorkflow() + "\" has been successfully deleted.\n"

}

func (s *IngressService) Config(ctx context.Context, req *pb.ConfigStruct) (*pb.MacroPodReply, error) {
	reply := Serve_Config(req)
	results := pb.MacroPodReply{Reply: &reply}
	return &results, nil
}

func (s *IngressService) WorkflowInvoke(ctx context.Context, req *pb.MacroPodRequest) (*pb.MacroPodReply, error) {
	function := entry_function[req.GetWorkflow()]
	req.Function = &function
	reply, code := Serve_WorkflowInvoke(req)
	results := pb.MacroPodReply{Reply: &reply, Code: &code}
	return &results, nil
}

func (s *IngressService) FunctionInvoke(ctx context.Context, req *pb.MacroPodRequest) (*pb.MacroPodReply, error) {
	reply, code := Serve_FunctionInvoke(req)
	results := pb.MacroPodReply{Reply: &reply, Code: &code}
	return &results, nil
}

func (s *IngressService) CreateWorkflow(ctx context.Context, req *pb.WorkflowStruct) (*pb.MacroPodReply, error) {
	reply := Serve_CreateWorkflow(req)
	results := pb.MacroPodReply{Reply: &reply}
	return &results, nil
}

func (s *IngressService) UpdateWorkflow(ctx context.Context, req *pb.WorkflowStruct) (*pb.MacroPodReply, error) {
	reply := Serve_UpdateWorkflow(req)
	results := pb.MacroPodReply{Reply: &reply}
	return &results, nil
}

func (s *IngressService) DeleteWorkflow(ctx context.Context, req *pb.MacroPodRequest) (*pb.MacroPodReply, error) {
	reply := Serve_DeleteWorkflow(req)
	results := pb.MacroPodReply{Reply: &reply}
	return &results, nil
}

func HTTP_Help(res http.ResponseWriter, req *http.Request) {
	help_print := "TODO\n"
	fmt.Fprint(res, help_print)
}

func HTTP_Config(res http.ResponseWriter, req *http.Request) {
	body, _ := ioutil.ReadAll(req.Body)
	request := pb.ConfigStruct{}
	json.Unmarshal(body, &request)
	Serve_Config(&request)
}

func HTTP_WorkflowInvoke(res http.ResponseWriter, req *http.Request) {
	workflow := req.PathValue("workflow")
	function := entry_function[workflow]
	request := pb.MacroPodRequest{Workflow: &workflow, Function: &function}
	body, _ := ioutil.ReadAll(req.Body)
	switch content := req.Header.Get("Content-Type"); content {
	case "text/plain":
		t := string(body)
		request.Text = &t
	case "application/json":
		j_i := make(map[string]interface{})
		json.Unmarshal(body, &j_i)
		j, _ := structpb.NewStruct(j_i)
		request.JSON = j
	case "application/octet-stream":
		request.Data = body
	default:
		t := string(body)
		request.Text = &t
	}
	results, _ := Serve_WorkflowInvoke(&request)
	fmt.Fprint(res, results)
}

func HTTP_FunctionInvoke(res http.ResponseWriter, req *http.Request) {
	body, _ := ioutil.ReadAll(req.Body)
	request := pb.MacroPodRequest{}
	json.Unmarshal(body, &request)
	results, _ := Serve_FunctionInvoke(&request)
	fmt.Fprint(res, results)
}

func HTTP_CreateWorkflow(res http.ResponseWriter, req *http.Request) {
	body, _ := ioutil.ReadAll(req.Body)
	request := pb.WorkflowStruct{}
	json.Unmarshal(body, &request)
	fmt.Printf("\nBody : %s\n request : %v\n", body, request)
	results := Serve_CreateWorkflow(&request)
	fmt.Fprint(res, results)
}

func HTTP_UpdateWorkflow(res http.ResponseWriter, req *http.Request) {
	body, _ := ioutil.ReadAll(req.Body)
	request := pb.WorkflowStruct{}
	json.Unmarshal(body, &request)
	results := Serve_UpdateWorkflow(&request)
	fmt.Fprint(res, results)
}

func HTTP_DeleteWorkflow(res http.ResponseWriter, req *http.Request) {
	workflow := req.PathValue("workflow")
	request := pb.MacroPodRequest{Workflow: &workflow}
	results := Serve_DeleteWorkflow(&request)
	fmt.Fprint(res, results)
}

func main() {
	service_port := os.Getenv("SERVICE_PORT")
	if service_port == "" {
		service_port = "8000"
	}
	http_port := os.Getenv("HTTP_PORT")
	if http_port == "" {
		http_port = "9000"
	}
	deployer_address = os.Getenv("DEPLOYER_ADDRESS")
	if deployer_address == "" {
		deployer_address = "127.0.0.1:8002"
	}
	namespace := os.Getenv("NAMESPACE")
	if namespace == "" {
		namespace = "macropod-functions"
	}
	ttl, err := strconv.Atoi(os.Getenv("TTL"))
	if err != nil {
		ttl = 180
	}
	deployment := os.Getenv("DEPLOYMENT")
	if deployment == "" {
		deployment = "macropod"
	}
	communication := os.Getenv("COMMUNICATION")
	if communication == "" {
		communication = "direct"
	}
	aggregation := os.Getenv("AGGREGATION")
	if aggregation == "" {
		aggregation = "agg"
	}
	target_concurrency, err := strconv.Atoi(os.Getenv("TARGET_CONCURRENCY"))
	if err != nil {
		target_concurrency = -1
	}
	debug, err := strconv.Atoi(os.Getenv("DEBUG"))
	if err != nil {
		debug = 0
	}
	ttl_i := int32(ttl)
	target_concurrency_i := int32(target_concurrency)
	debug_i := int32(debug)
	c := pb.ConfigStruct{Namespace: &namespace, TTL: &ttl_i, Deployment: &deployment, Communication: &communication, Aggregation: &aggregation, TargetConcurrency: &target_concurrency_i, Debug: &debug_i}
	default_config = &c
	l, _ := net.Listen("tcp", ":"+service_port)
	s := grpc.NewServer(grpc.MaxSendMsgSize(1024*1024*200), grpc.MaxRecvMsgSize(1024*1024*200))
	pb.RegisterMacroPodIngressServer(s, &IngressService{})
	updateRun = false
	createRun = false
	Deployer_Check()
	Deployer_Update_Check()
	go WatchTTL()
	go WatchIngress()
	go s.Serve(l)

	h := http.NewServeMux()
	h.HandleFunc("/", HTTP_Help)
	h.HandleFunc("/config", HTTP_Config)
	h.HandleFunc("/workflow/invoke/{workflow}", HTTP_WorkflowInvoke)
	h.HandleFunc("/function/invoke", HTTP_FunctionInvoke)
	h.HandleFunc("/workflow/create", HTTP_CreateWorkflow)
	h.HandleFunc("/workflow/update", HTTP_UpdateWorkflow)
	h.HandleFunc("/workflow/delete/{workflow}", HTTP_DeleteWorkflow)
	http.ListenAndServe(":"+http_port, h)
}

