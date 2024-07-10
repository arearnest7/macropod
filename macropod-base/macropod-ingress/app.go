package main

// to do - first invoke
import (
	pb "app/deployer_pb"
	pp_pb "app/prepuller_pb"
	wf_pb "app/wf_pb"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
	"strings"
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

type WorkflowMetadata struct {
	IngressList map[string]string
	InUse       map[string]bool
	LastUsed    map[string]time.Time
}

var (
	workflows          map[string]Workflow
	deployerInProgress bool
	workflows_metadata map[string]WorkflowMetadata
	kclient            *kubernetes.Clientset
	ttl_seconds        float64
	mutex sync.RWMutex
	retrypolicy        = `{
		"methodConfig": [{
		  "name": [{}],
		  "waitForReady": true,
		  "retryPolicy": {
			  "MaxAttempts": 30,
			  "InitialBackoff": "1s",
			  "MaxBackoff": "10s",
			  "BackoffMultiplier": 2.0,
			  "RetryableStatusCodes": [ "UNAVAILABLE", "UNKNOWN", "DEADLINE_EXCEEDED"]
		  }
		}]}`
)

func internal_log(message string) {
	fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.000000 UTC") + ": " + message)
}

func Serve_Help(res http.ResponseWriter, req *http.Request) {
	help_print := "macropod ingress\nplease use the following services to interact/deploy/invoke workflows:\n"
	help_print += "Invoke:\n - path: /invoke/[workflow_name]\n - purpose: invoke workflow at entrypoint\n - payload: JSON body for entrypoint function\n - output: string result\n"
	help_print += "Create:\n - path: /create/[workflow_name]\n - purpose: deploy workflow based on configuration\n - payload: macropod JSON configuration\n - output: string confirmation\n"
	help_print += "Update:\n - path: /update/[workflow_name]\n - purpose: update previously deployed workflow\n - payload: macropod JSON configuration\n - output: string confirmation\n"
	help_print += "Delete:\n - path: /delete/[workflow_name]\n - purpose: delete workflow of same name\n - payload: none\n - output: string confirmation\n"
	help_print += "Logs:\n - path: /logs/[workflow_name]\n - purpose: output logs for a workflow\n - payload: none\n - output: CSV\n"
	help_print += "Metrics:\n - path: /metrics\n - purpose: output current metrics\n - payload: none\n - output: CSV\n"
	fmt.Fprint(res, help_print)
}

func callDepController(existing bool, wf_name string) string {
	if existing && deployerInProgress {
		internal_log("deployer in progress for existing invoke")
		return "deployer in progress"
	}
	deployerInProgress = true
	internal_log("deployer call start with " + wf_name)
	depControllerAddr := os.Getenv("DEP_CONTROLLER_ADD")
	if depControllerAddr == "" {
		internal_log("DEP_CONTROLLER_ADD environment variable not set")
		return ""
	}
	var type_call string
	if existing {
		type_call = "existing_invoke"
	} else {
		type_call = "new_invoke"
	}
	opts := grpc.WithInsecure()
	cc, err := grpc.Dial(depControllerAddr, opts)
	if err != nil {
		internal_log("failure to dial deployment controller - " + err.Error())
		return ""
	}
	defer cc.Close()
	client := pb.NewDeploymentServiceClient(cc)
	request := &pb.DeploymentServiceRequest{WorkflowName: wf_name, RequestType: type_call}
	resp, err := client.Deployment(context.Background(), request)
	if err != nil {
		internal_log("deployment controller failure with " + type_call + " " + wf_name + "-" + err.Error())
		return ""
	}
	internal_log("deployer call end with " + wf_name)
	deployerInProgress = false
	return resp.Message
}

func Serve_WF_Invoke(res http.ResponseWriter, req *http.Request) {
	internal_log("WF_INVOKE_START " + req.PathValue("wf_name"))
	payload, err := ioutil.ReadAll(req.Body)
	if err != nil {
		internal_log("invoke body - " + err.Error())
		return
	}
	workflow_id := strconv.Itoa(rand.Intn(10000000))
	wf_name := req.PathValue("wf_name")
	invoked := false
	var results string
	request_type := "gg"
	mutex.RLock()
	ingressList := workflows_metadata[wf_name].IngressList
	mutex.RUnlock()
	for namespace, ingress := range ingressList {
		mutex.Lock()
		defer mutex.Unlock()
		if !workflows_metadata[wf_name].InUse[namespace] {
			invoked = true
			workflows_metadata[wf_name].InUse[namespace] = true
			workflows_metadata[wf_name].LastUsed[namespace] = time.Now()
			internal_log("requesting existing replica " + ingress)
			go callDepController(true, req.PathValue("wf_name"))
			channel, err := grpc.Dial(ingress, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*200), grpc.MaxCallSendMsgSize(1024*1024*200)), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultServiceConfig(retrypolicy))
			if err != nil {
				internal_log("failure to dial existing ingress " + ingress + " - " + err.Error())
				workflows_metadata[wf_name].InUse[namespace] = false
				return
			}
			defer channel.Close()
			stub := wf_pb.NewGRPCFunctionClient(channel)
			response, err := stub.GRPCFunctionHandler(context.Background(), &wf_pb.RequestBody{Data: payload, WorkflowId: workflow_id, Depth: 0, Width: 0, RequestType: &request_type})
			if err != nil {
				internal_log("existing invoke response error from " + ingress + err.Error())
				workflows_metadata[wf_name].InUse[namespace] = false
				// log.Print(results)
				http.Error(res, err.Error(), http.StatusBadGateway)
				return
			}
			results = response.GetReply()
			internal_log("response from existing replica " + ingress)
			internal_log("WF_INVOKE_END " + req.PathValue("wf_name"))
			// log.Print(results)
			fmt.Fprint(res, results)
			internal_log("releasing namespace " + namespace)
			workflows_metadata[wf_name].InUse[namespace] = false
			return
		}
	}
	if !invoked {
		internal_log("New Invoke")
		go callDepController(false, wf_name) 
		for {
			request_type := "gg"
			mutex.RLock()
			ingressList := workflows_metadata[wf_name].IngressList
			mutex.RUnlock()
			for namespace, ingress := range ingressList {
				mutex.Lock()
				defer mutex.Unlock()
				if !workflows_metadata[wf_name].InUse[namespace] {
					invoked = true
					workflows_metadata[wf_name].InUse[namespace] = true
					workflows_metadata[wf_name].LastUsed[namespace] = time.Now()
					internal_log("requesting existing replica " + ingress)
					go callDepController(true, req.PathValue("wf_name"))
					channel, err := grpc.Dial(ingress, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*200), grpc.MaxCallSendMsgSize(1024*1024*200)), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultServiceConfig(retrypolicy))
					if err != nil {
						internal_log("failure to dial existing ingress " + ingress + " - " + err.Error())
						workflows_metadata[wf_name].InUse[namespace] = false
						return
					}
					defer channel.Close()
					stub := wf_pb.NewGRPCFunctionClient(channel)
					response, err := stub.GRPCFunctionHandler(context.Background(), &wf_pb.RequestBody{Data: payload, WorkflowId: workflow_id, Depth: 0, Width: 0, RequestType: &request_type})
					if err != nil {
						internal_log("existing invoke response error from " + ingress + err.Error())
						workflows_metadata[wf_name].InUse[namespace] = false
						// log.Print(results)
						http.Error(res, err.Error(), http.StatusBadGateway)
						return
					}
					results = response.GetReply()
					internal_log("response from existing replica " + ingress)
					internal_log("WF_INVOKE_END " + req.PathValue("wf_name"))
					// log.Print(results)
					fmt.Fprint(res, results)
					internal_log("releasing namespace " + namespace)
					workflows_metadata[wf_name].InUse[namespace] = false
					return
				}
			}
		}

	}
}

func dialPrePuller(body string) {
	internal_log("Prepuller")
	ip := "127.0.0.1"
	prepuller := ip + ":5003"
	opts := grpc.WithInsecure()
	cc, err := grpc.Dial(prepuller, opts)
	if err != nil {
		internal_log("create grpc - " + err.Error())
	}
	defer cc.Close()
	client := pp_pb.NewPrepullerServiceClient(cc)
	request := &pp_pb.PrepullerServiceRequest{Data: &body}
	r, err := client.Prepuller(context.Background(), request)
	if err != nil {
		internal_log("create return - " + err.Error())
	}
	internal_log(r.Message)

}

func Serve_WF_Create(res http.ResponseWriter, req *http.Request) {
	internal_log("WF_CREATE_START " + req.PathValue("wf_name"))
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
	request := &pb.DeploymentServiceRequest{WorkflowName: req.PathValue("wf_name"), RequestType: "create", Data: &workflow}
	dialPrePuller(workflow)
	internal_log("requesting create for " + req.PathValue("wf_name"))
	_, err = client.Deployment(context.Background(), request)
	internal_log("returned create for " + req.PathValue("wf_name"))
	if err != nil {
		internal_log("create return - " + err.Error())
	}
	workflows[req.PathValue("wf_name")] = body_u
	internal_log("WF_CREATE_END " + req.PathValue("wf_name"))
	fmt.Fprintf(res, "Workflow created successfully. Invoke your workflow with /invoke/"+req.PathValue("wf_name"))
}

func Serve_WF_Update(res http.ResponseWriter, req *http.Request) {
	internal_log("WF_UPDATE_START " + req.PathValue("wf_name"))
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
	request := &pb.DeploymentServiceRequest{WorkflowName: req.PathValue("wf_name"), RequestType: "update", Data: &workflow}
	dialPrePuller(workflow)
	internal_log("requesting update for " + req.PathValue("wf_name"))
	_, err = client.Deployment(context.Background(), request)
	internal_log("returned update for " + req.PathValue("wf_name"))
	if err != nil {
		internal_log("update return - " + err.Error())
	}
	workflows[req.PathValue("wf_name")] = body_u
	internal_log("WF_UPDATE_END " + req.PathValue("wf_name"))
	fmt.Fprintf(res, "Workflow \""+req.PathValue("wf_name")+"\" has been updated successfully.")
}

func Serve_WF_Delete(res http.ResponseWriter, req *http.Request) {
	internal_log("WF_DELETE_START " + req.PathValue("wf_name"))
	opts := grpc.WithInsecure()
	cc, err := grpc.Dial(os.Getenv("DEP_CONTROLLER_ADD"), opts)
	if err != nil {
		internal_log("delete grpc - " + err.Error())
	}
	defer cc.Close()
	client := pb.NewDeploymentServiceClient(cc)
	request := &pb.DeploymentServiceRequest{WorkflowName: req.PathValue("wf_name"), RequestType: "delete"}
	internal_log("requesting delete for " + req.PathValue("wf_name"))
	_, err = client.Deployment(context.Background(), request)
	internal_log("returned delete for " + req.PathValue("wf_name"))
	if err != nil {
		internal_log("delete return - " + err.Error())
	}
	delete(workflows, req.PathValue("wf_name"))
	wf_name := req.PathValue("wf_name")
	namespaces, _ := kclient.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	for _, ns := range namespaces.Items {

		if strings.Contains(ns.Name, wf_name) {
			kclient.CoreV1().Namespaces().Delete(context.Background(), ns.Name, metav1.DeleteOptions{})
		}
	}

	internal_log("WF_DELETE_END " + req.PathValue("wf_name"))
	fmt.Fprintf(res, "Workflow \""+req.PathValue("wf_name")+"\" has been deleted successfully.")
}

func Serve_Logs(res http.ResponseWriter, req *http.Request) {
	internal_log("LOGS_START " + req.PathValue("wf_name"))
	opts := grpc.WithInsecure()
	cc, err := grpc.Dial(os.Getenv("DEP_CONTROLLER_ADD"), opts)
	if err != nil {
		internal_log("logs grpc - " + err.Error())
	}
	defer cc.Close()
	client := pb.NewDeploymentServiceClient(cc)
	request := &pb.DeploymentServiceRequest{WorkflowName: req.PathValue("wf_name"), RequestType: "logs"}
	internal_log("requesting logs for " + req.PathValue("wf_name"))
	response, err := client.Deployment(context.Background(), request)
	internal_log("returned logs for " + req.PathValue("wf_name"))
	if err != nil {
		internal_log("logs return - " + err.Error())
	}
	internal_log("LOGS_END " + req.PathValue("wf_name"))
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
	request := &pb.DeploymentServiceRequest{WorkflowName: "", RequestType: "metrics"}
	internal_log("requesting metrics")
	response, err := client.Deployment(context.Background(), request)
	internal_log("returned metrics")
	if err != nil {
		internal_log("metrics return - " + err.Error())
	}
	internal_log("METRICS_END")
	fmt.Fprintf(res, response.Message)
}

func deleteWorkflowIngress(ingress *networkingv1.Ingress) {
	wf_name := ingress.Labels["workflow_name"]
	mutex.RLock()
	_, ok := workflows_metadata[wf_name]
	mutex.RUnlock()
	namespace := ingress.Labels["replica_namespace"]
	internal_log("deleting workflow ingress for " + namespace)
	if ok {
		mutex.Lock()
		delete(workflows_metadata[wf_name].IngressList, namespace)
		delete(workflows_metadata[wf_name].InUse, namespace)
		delete(workflows_metadata[wf_name].LastUsed, namespace)
		mutex.Unlock()
	}
	internal_log("deleted workflow ingress for " + namespace)
}

func updateWorkflowIngress(ingress *networkingv1.Ingress) {
	wf_name := ingress.Labels["workflow_name"]
	namespace := ""
	if workflows_metadata[wf_name].IngressList == nil {
		workflows_metadata[wf_name] = WorkflowMetadata{
			IngressList: make(map[string]string),
			LastUsed:    make(map[string]time.Time),
			InUse:       make(map[string]bool),
		}
	}
	mutex.Lock()
	defer mutex.Unlock()
	for _, rule := range ingress.Spec.Rules {
		if rule.HTTP != nil {
			for _, path := range rule.HTTP.Paths {
				serviceName := path.Backend.Service.Name
				namespace = ingress.Namespace
				port := path.Backend.Service.Port.Number
				hostname := fmt.Sprintf("%s.%s.svc.cluster.local:%d", serviceName, namespace, port)
				internal_log("updating ingress for " + namespace)
				workflows_metadata[wf_name].IngressList[namespace] = hostname
				workflows_metadata[wf_name].InUse[namespace] = false
				workflows_metadata[wf_name].LastUsed[namespace] = time.Now()
				internal_log("updated ingress for " + namespace)
				return
			}
		}
	}
}

func watchWorkflows() {
	internal_log("WATCH_WORKFLOW_START")
	watcher, err := kclient.NetworkingV1().Ingresses("").Watch(context.TODO(), metav1.ListOptions{})
	if err != nil {
		internal_log("Failed to set up watcher - " + err.Error())
	}
	for event := range watcher.ResultChan() {
		ingress, ok := event.Object.(*networkingv1.Ingress)
		if !ok {
			internal_log("Invalid Ingress Event")
			continue
		}
		switch event.Type {
		case watch.Added, watch.Modified:
			updateWorkflowIngress(ingress)
		case watch.Deleted:
			deleteWorkflowIngress(ingress)
		}
	}
	internal_log("WATCH_WORKFLOW_END")
}

func watchTTL() {
	internal_log("WATCH_TTL_START")
	for {
		currentTime := time.Now()
		mutex.Lock()
		for _, workflow_metadata := range workflows_metadata {
			for namespace := range workflow_metadata.IngressList {
				if !workflow_metadata.InUse[namespace] && currentTime.Sub(workflow_metadata.LastUsed[namespace]).Seconds() > ttl_seconds {
					internal_log("TTL exceeded... deleting ingress " + namespace)
					kclient.CoreV1().Namespaces().Delete(context.TODO(), namespace, metav1.DeleteOptions{})
					delete(workflow_metadata.IngressList, namespace)
					delete(workflow_metadata.InUse, namespace)
					delete(workflow_metadata.LastUsed, namespace)
					internal_log("TTL exceeded... deleted ingress " + namespace)
				}
			}
		}
		mutex.Unlock()
		time.Sleep(time.Second)
	}
	// internal_log("WATCH_TTL_END") -> unreachable code
}

func main() {
	deployerInProgress = false
	workflows = make(map[string]Workflow)
	workflows_metadata = make(map[string]WorkflowMetadata)

	internal_log("Ingress Controller Started")
	config, err := rest.InClusterConfig()
	if err != nil {
		internal_log("config - " + err.Error())
		return
	}
	kclient, err = kubernetes.NewForConfig(config)
	if err != nil {
		internal_log("client - " + err.Error())
		return
	}
	ttl_seconds, err = strconv.ParseFloat(os.Getenv("TTL"), 64)
	if err != nil {
		internal_log("ttl - " + err.Error())
		return
	}
	go watchWorkflows()
	go watchTTL()
	h := http.NewServeMux()
	h.HandleFunc("/", Serve_Help)
	h.HandleFunc("/invoke/{wf_name}", Serve_WF_Invoke)
	h.HandleFunc("/create/{wf_name}", Serve_WF_Create)
	h.HandleFunc("/update/{wf_name}", Serve_WF_Update)
	h.HandleFunc("/delete/{wf_name}", Serve_WF_Delete)
	h.HandleFunc("/logs/{wf_name}", Serve_Logs)
	h.HandleFunc("/metrics", Serve_Metrics)
	http.ListenAndServe(":"+os.Getenv("SERVICE_PORT"), h)
}
