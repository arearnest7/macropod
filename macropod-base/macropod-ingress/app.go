package main

import (
    "net/http"
    "time"
    "fmt"
    "strconv"
    "context"
    "io/ioutil"
    "math/rand"
    "os"
    "encoding/json"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
    networkingv1 "k8s.io/api/networking/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/watch"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/rest"
    wf_pb "app/wf_pb"
    pb "app/deployer_pb"
)

type Function struct {
    Registry string `json:"registry"`
    Endpoints []string `json:"endpoints,omitempty"`
    Envs map[string]string `json:"envs,omitempty"`
    Secrets map[string]string `json:"secrets,omitempty"`
}

type Workflow struct {
    Name string `json:"name,omitempty"`
    Functions map[string]Function `json:"functions"`
    IngressList map[string]string
    InUse map[string]bool
    LastUsed map[string]time.Time
}

var (
    workflows map[string]Workflow
    kclient *kubernetes.Clientset
    ttl_seconds float64
)

func internal_log(message string) {
    fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.000000 UTC") + ": " + message)
}

func Serve_Help(res http.ResponseWriter, req *http.Request) {
    help_print := "macropod ingress\nplease use the following services to interact/deploy/invoke workflows:\n"
    help_print += "Invoke:\n - path: /invoke/[workflow_name]\n - purpose: invoke workflow at entrypoint\n - payload: JSON body for entrypoint function\n - output: string result\n"
    help_print += "Create:\n - path: /deploy/[workflow_name]\n - purpose: deploy workflow based on configuration\n - payload: macropod JSON configuration\n - output: string confirmation\n"
    help_print += "Update:\n - path: /update/[workflow_name]\n - purpose: update previously deployed workflow\n - payload: macropod JSON configuration\n - output: string confirmation\n"
    help_print += "Delete:\n - path: /delete/[workflow_name]\n - purpose: delete workflow of same name\n - payload: none\n - output: string confirmation\n"
    help_print += "Logs:\n - path: /logs/[workflow_name]\n - purpose: output logs for a workflow\n - payload: none\n - output: CSV\n"
    help_print += "Metrics:\n - path: /metrics\n - purpose: output current metrics\n - payload: none\n - output: CSV\n"
    fmt.Fprintf(res, help_print)
}

func callDepController(existing bool, wf_name string) (string) {
    fmt.Printf("deployer call end with " + wf_name)
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
    fmt.Printf("deployer call end with " + wf_name)
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
    workflow := workflows[req.PathValue("wf_name")]
    var response *wf_pb.ResponseBody
    var invoked bool
    request_type := "gg"
    for namespace, ingress := range workflow.IngressList {
        if !workflow.InUse[namespace] {
            invoked = true
            workflow.InUse[namespace] = true
            internal_log("requesting existing replica " + ingress)
            go callDepController(true, req.PathValue("wf_name"))
            channel, err := grpc.Dial(ingress, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*200), grpc.MaxCallSendMsgSize(1024*1024*200)), grpc.WithTransportCredentials(insecure.NewCredentials()))
            if err != nil {
                internal_log("failure to dial existing ingress " + ingress + " - " + err.Error())
                return
            }
            defer channel.Close()
            stub := wf_pb.NewGRPCFunctionClient(channel)
            ctx, cancel := context.WithTimeout(context.Background(), time.Second)
            defer cancel()
            response, err = stub.GRPCFunctionHandler(ctx, &wf_pb.RequestBody{Data: payload, WorkflowId: workflow_id, Depth: 0, Width: 0, RequestType: &request_type})
            if err != nil {
                internal_log("response failure from existing ingress " + ingress + " - " + err.Error())
                return
            }
            workflow.InUse[namespace] = false
            internal_log("response from existing replica " + ingress)
        }
    }
    if !invoked {
        ingress := callDepController(false, req.PathValue("wf_name"))
        channel, err := grpc.Dial(ingress, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*200), grpc.MaxCallSendMsgSize(1024*1024*200)), grpc.WithTransportCredentials(insecure.NewCredentials()))
        if err != nil {
            internal_log("failure to dial new ingress " + ingress + " - " + err.Error())
            return
        }
        defer channel.Close()
        stub := wf_pb.NewGRPCFunctionClient(channel)
        ctx, cancel := context.WithTimeout(context.Background(), time.Second)
        defer cancel()
        response, err = stub.GRPCFunctionHandler(ctx, &wf_pb.RequestBody{Data: payload, WorkflowId: workflow_id, Depth: 0, Width: 0, RequestType: &request_type})
        if err != nil {
            internal_log("response failure from new ingress " + ingress + " - " + err.Error())
            return
        }
        internal_log("response from new replica " + ingress)
    }
    internal_log("WF_INVOKE_END " + req.PathValue("wf_name"))
    fmt.Fprintf(res, response.GetReply())
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
    internal_log("requesting create for " + req.PathValue("wf_name"))
    _, err = client.Deployment(context.Background(), request)
    internal_log("returned create for " + req.PathValue("wf_name"))
    if err != nil {
        internal_log("create return - " + err.Error())
    }
    workflows[req.PathValue("wf_name")] = body_u
    internal_log("WF_CREATE_END " + req.PathValue("wf_name"))
    fmt.Fprintf(res, "Workflow created successfully. Invoke your workflow with /invoke/" + req.PathValue("wf_name"))
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
    internal_log("requesting update for " + req.PathValue("wf_name"))
    _, err = client.Deployment(context.Background(), request)
    internal_log("returned update for " + req.PathValue("wf_name"))
    if err != nil {
        internal_log("update return - " + err.Error())
    }
    workflows[req.PathValue("wf_name")] = body_u
    internal_log("WF_UPDATE_END " + req.PathValue("wf_name"))
    fmt.Fprintf(res, "Workflow \"" + req.PathValue("wf_name") + "\" has been updated successfully.")
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
    internal_log("WF_DELETE_END " + req.PathValue("wf_name"))
    fmt.Fprintf(res, "Workflow \"" + req.PathValue("wf_name") + "\" has been deleted successfully.")
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
    _, ok := workflows[wf_name]
    namespace := ingress.Labels["replica_namespace"]
    internal_log("deleting workflow ingress for " + namespace)
    if ok {
        delete(workflows[wf_name].IngressList, namespace)
        delete(workflows[wf_name].InUse, namespace)
        delete(workflows[wf_name].LastUsed, namespace)
    }
    internal_log("deleted workflow ingress for " + namespace)
}

func updateWorkflowIngress(ingress *networkingv1.Ingress) {
    wf_name := ingress.Labels["workflow_name"]
    namespace := ""
    for _, rule := range ingress.Spec.Rules {
        if rule.HTTP != nil {
            for _, path := range rule.HTTP.Paths {
                serviceName := path.Backend.Service.Name
                namespace = ingress.Namespace
                port := path.Backend.Service.Port.Number
                hostname := fmt.Sprintf("%s.%s.svc.cluster.local:%d", serviceName, namespace, port)
                internal_log("updating ingress for " + namespace)
                workflows[wf_name].IngressList[namespace] = hostname
                workflows[wf_name].InUse[namespace] = false
                workflows[wf_name].LastUsed[namespace] = time.Now()
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
        for _, workflow := range workflows {
            for namespace, _ := range workflow.IngressList {
                if !workflow.InUse[namespace] && currentTime.Sub(workflow.LastUsed[namespace]).Seconds() > ttl_seconds {
                    internal_log("TTL exceeded... deleting ingress " + namespace)
                    kclient.CoreV1().Namespaces().Delete(context.TODO(), namespace, metav1.DeleteOptions{})
                    delete(workflow.IngressList, namespace)
                    delete(workflow.InUse, namespace)
                    delete(workflow.LastUsed, namespace)
                    internal_log("TTL exceeded... deleted ingress " + namespace)
                }
            }
        }
        time.Sleep(time.Second)
    }
    internal_log("WATCH_TTL_END")
}

func main() {
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
    // with future services add them above this line in form shown
    http.ListenAndServe(":" + os.Getenv("SERVICE_PORT"), h)
}
