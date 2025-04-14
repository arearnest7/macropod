package main

import (
    pb "app/macropod_pb"

    "os"
    "encoding/json"
    "fmt"
    "strconv"

    "github.com/soheilhy/cmux"
    "net"
    "net/http"

    "google.golang.org/grpc"
    "golang.org/x/net/context"
)

type IngressService struct {
    pb.UnimplementedMacroPodIngressServer
}

var (
    manifest         = make(pb.MacroPodManifest)

    service_channel  = make(map[string][]*grpc.ClientConn)
    service_stub     = make(map[string][]pb.MacroPodFunctionClient)

    deployer_channel *grpc.ClientConn
    deployer_stub    pb.MacroPodDeployerClient
    deployer_add     string

    dataLock         sync.Mutex

    retrypolicy      = `{
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
    dataLock.Lock()
    for deployer_channel == nil || deployer_channel.GetState() != connectivity.Ready {
        deployer_channel, _ = grpc.Dial(deployer_add, grpc.WithInsecure())
        deployer_stub = pb.NewDeploymentServiceClient(deployer_channel)
        time.Sleep(10 * time.Millisecond)
    }
    dataLock.Unlock()
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
    resp, err := deployer_stub.Deployment(context.Background(), request)
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
                fmt.Println("deployer already modifying workflow " + func_name)
            }
        case "2": // target concurrency decrease
            if debug > 2 {
                fmt.Println("workflow hit threshold, scaling down target concurrency")
                dataLock.Lock()
                workflow_target_concurrency[func_name] = workflow_invocations_current[func_name] / len(service_target[func_name]) / 2
                dataLock.Unlock()
            }
        default:
            if debug > 2 {
                fmt.Println("Unknown response code " + resp_code[0])
            }
    }

    return nil//TODO
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
                    for {
                        go callDepController("ttl_delete", func_name, labels)
                        time.Sleep(100 * time.Millisecond)
                        deployments_list, _ := k.CoreV1().Pods(macropod_namespace).List(context.Background(), metav1.ListOptions{LabelSelector: labels})
                        if deployments_list == nil || len(deployments_list.Items) == 0 {
                            break
                        }
                    }
                }
            }
            time.Sleep(time.Second)
        }
    }//TODO

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
    }//TODO
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

    }//TODO

}

func deployerSync() {
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

    }//TODO
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
    return target, triggered//TODO
}

func Ingress_Config(config pb.ConfigStruct) (error) {
    
    return nil //TODO implement in macropod-ingress
}

func Ingress_WorkflowInvoke(ingress_request pb.IngressRequest) (string, error) {
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
    dataLock.Lock()
    invocations_current := strconv.Itoa(workflow_invocations_current[func_name])
    dataLock.Unlock()
    go callDepController("update_deployments", func_name, invocations_current)
    request_type := "gg"
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
    return "", nil //TODO implement in macropod-ingress
}

func Ingress_FunctionInvoke(function_request pb.FunctionRequest) (string, error) {
    
    return "", nil //TODO implement in macropod-ingress
}

func Ingress_CreateWorkflow(workflow pb.WorkflowStruct) (error) {
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
    return nil //TODO implement in macropod-ingress
}

func Ingress_UpdateWorkflow(workflow pb.WorkflowStruct) (error) {
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
    return nil //TODO implement in macropod-ingress
}

func Ingress_DeleteWorkflow(ingress_request pb.IngressRequest) (error) {
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
    return nil //TODO implement in macropod-ingress
}

func Ingress_Eval(workflow pb.WorkflowStruct) (string, error) {
    
    return "", nil //TODO implement in macropod-ingress
}

func Ingress_EvalMetrics(ingress_eval_request pb.IngressEvalRequest) (string, error) {
    
    return "", nil //TODO implement in macropod-ingress
}

func Ingress_EvalLatency(ingress_eval_request pb.IngressEvalRequest) (string, error) {
    
    return "", nil //TODO implement in macropod-ingress
}

func Ingress_EvalSummary(ingress_eval_request pb.IngressEvalRequest) (string, error) {
    
    return "", nil //TODO implement in macropod-ingress
}

func (s *IngressService) Config(ctx context.Context, in *pb.ConfigStruct) (*pb.MacroPodReply, error) {
    err := Ingress_Config(*in)
    res := pb.MacroPodReply{}
    return &res, err
}

func (s *IngressService) WorkflowInvoke(ctx context.Context, in *pb.IngressRequest) (*pb.MacroPodReply, error) {
    results, err := Ingress_WorkflowInvoke(*in)
    res := pb.MacroPodReply{Reply: results}
    return &res, err
}

func (s *IngressService) FunctionInvoke(ctx context.Context, in *pb.FunctionRequest) (*pb.MacroPodReply, error) {
    results, err := Ingress_FunctionInvoke(*in)
    res := pb.MacroPodReply{Reply: results}
    return &res, err
}

func (s *IngressService) CreateWorkflow(ctx context.Context, in *pb.WorkflowStruct) (*pb.MacroPodReply, error) {
    err := Ingress_CreateWorkflow(*in)
    res := pb.MacroPodReply{}
    return &res, err
}

func (s *IngressService) UpdateWorkflow(ctx context.Context, in *pb.WorkflowStruct) (*pb.MacroPodReply, error) {
    err := Ingress_UpdateWorkflow(*in)
    res := pb.MacroPodReply{}
    return &res, err
}

func (s *IngressService) DeleteWorkflow(ctx context.Context, in *pb.IngressRequest) (*pb.MacroPodReply, error) {
    err := Ingress_DeleteWorkflow(*in)
    res := pb.MacroPodReply{}
    return &res, err
}

func (s *IngressService) Eval(ctx context.Context, in *pb.WorkflowStruct) (*pb.MacroPodReply, error) {
    id, err := Ingress_Eval(*in)
    res := pb.MacroPodReply{Reply: id}
    return &res, err
}

func (s *IngressService) EvalMetrics(ctx context.Context, in *pb.IngressEvalRequest) (*pb.MacroPodReply, error) {
    metrics, err := Ingress_EvalMetrics(*in)
    res := pb.MacroPodReply{Reply: metrics}
    return &res, err
}

func (s *IngressService) EvalLatency(ctx context.Context, in *pb.IngressEvalRequest) (*pb.MacroPodReply, error) {
    latency, err := Ingress_EvalLatency(*in)
    res := pb.MacroPodReply{Reply: latency}
    return &res, err
}

func (s *IngressService) EvalSummary(ctx context.Context, in *pb.IngressEvalRequest) (*pb.MacroPodReply, error) {
    summary, err := Ingress_EvalSummary(*in)
    res := pb.MacroPodReply{Reply: summary}
    return &res, err
}

func Serve_Ingress_Help(res http.ResponseWriter, req *http.Request) {
    help_print := "MacroPod Ingress\nplease use the following paths:\n"
    help_print += "Config:\n - path: /config\n - purpose: \n - payload: \n - output: \n"
    help_print += "Workflow Invoke:\n - path: /workflow/invoke/{wf_name}\n - purpose: \n - payload: \n - output: \n"
    help_print += "Function Invoke:\n - path: /function/invoke/{func_name}\n - purpose: \n - payload: \n - output: \n"
    help_print += "Workflow Create:\n - path: /workflow/create/{wf_name}\n - purpose: \n - payload: \n - output: \n"
    help_print += "Workflow Update:\n - path: /workflow/update/{wf_name}\n - purpose: \n - payload: \n - output: \n"
    help_print += "Workflow Delete:\n - path: /workflow/delete/{wf_name}\n - purpose: \n - payload: \n - output: \n"
    help_print += "Eval Start:\n - path: /eval/start\n - purpose: \n - payload: \n - output: \n"
    help_print += "Eval Metrics:\n - path: /eval/metrics/{id}\n - purpose: \n - payload: \n - output: \n"
    help_print += "Eval Latency:\n - path: /eval/latency/{id}\n - purpose: \n - payload: \n - output: \n"
    help_print += "Eval Summary:\n - path: /eval/summary/{id}\n - purpose: \n - payload: \n - output: \n"
    fmt.Fprint(res, help_print)
}

func Serve_Ingress_Config(res http.ResponseWriter, req *http.Request) {
    body, _ := ioutil.ReadAll(req.Body)
    config := pb.ConfigStruct{}
    json.Unmarshal(body, &config)
    defer req.Body.Close()
    err := Ingress_Config(config)
    if err != nil {
        fmt.Fprint(res, err)
        return
    }
    fmt.Fprint(res, "MacroPod config has been set")
}

func Serve_Ingress_WorkflowInvoke(res http.ResponseWriter, req *http.Request) {
    body, _ := ioutil.ReadAll(req.Body)
    ingress_request := pb.IngressRequest{Workflow: req.PathValue("wf_name"), Payload: body}
    defer req.Body.Close()
    results, err := Ingress_WorkflowInvoke(ingress_request)
    if err != nil {
        fmt.Fprint(res, err)
        return
    }
    fmt.Print(res, results)
}

func Serve_Ingress_FunctionInvoke(res http.ResponseWriter, req *http.Request) {
    body, _ := ioutil.ReadAll(req.Body)
    function_request := pb.FunctionRequest{}
    json.Unmarshal(body, &payload)
    defer req.Body.Close()
    function_request.Function = req.PathValue("func_name")
    results, err := Ingress_FunctionInvoke(function_request)
    if err != nil {
        fmt.Fprint(res, err)
        return
    }
    fmt.Fprint(res, results)
}

func Serve_Ingress_CreateWorkflow(res http.ResponseWriter, req *http.Request) {
    body, _ := ioutil.ReadAll(req.Body)
    workflow := pb.WorkflowStruct{}
    json.Unmarshal(body, &workflow)
    defer req.Body.Close()
    err := Ingress_CreateWorkflow(workflow)
    if err != nil {
        fmt.Fprint(res, err)
        return
    }
    fmt.Fprint(res, "Workflow created successfully. Invoke your workflow with /workflow/invoke/"+workflow.GetWorkflow()+"\n")
}

func Serve_Ingress_UpdateWorkflow(res http.ResponseWriter, req *http.Request) {
    body, _ := ioutil.ReadAll(req.Body)
    workflow := pb.WorkflowStruct{}
    json.Unmarshal(body, &workflow)
    defer req.Body.Close()
    err := Ingress_UpdateWorkflow(workflow)
    if err != nil {
        fmt.Fprint(res, err)
        return
    }
    fmt.Fprint(res, "Workflow \""+workflow.GetWorkflow()+"\" has been updated successfully.\n")
}

func Serve_Ingress_DeleteWorkflow(res http.ResponseWriter, req *http.Request) {
    ingress_request := pb.IngressRequest{Workflow: req.PathValue("wf_name")}
    err := Ingress_DeleteWorkflow(ingress_request)
    if err != nil {
        fmt.Fprint(res, err)
        return
    }
    fmt.Fprint(res, "Workflow \""+workflow.GetWorkflow()+"\" has been deleted successfully.\n")
}

func Serve_Ingress_Eval(res http.ResponseWriter, req *http.Request) {
    body, _ := ioutil.ReadAll(req.Body)
    workflow := pb.WorkflowStruct{}
    json.Unmarshal(body, &workflow)
    defer req.Body.Close()
    id, err := Ingress_Eval(workflow)
    if err != nil {
        fmt.Fprint(res, err)
        return
    }
    fmt.Fprint(res, id)
}

func Serve_Ingress_EvalMetrics(res http.ResponseWriter, req *http.Request) {
    id := int(req.PathValue("id"))
    metrics, err := Ingress_EvalMetrics(pb.IngressEvalRequest{ID: id})
    if err != nil {
        fmt.Fprint(res, err)
        return
    }
    fmt.Fprint(res, metrics)
}

func Serve_Ingress_EvalLatency(res http.ResponseWriter, req *http.Request) {
    id := int(req.PathValue("id"))
    latency, err := Ingress_EvalLatency(pb.IngressEvalRequest{ID: id})
    if err != nil {
        fmt.Fprint(res, err)
        return
    }
    fmt.Fprint(res, latency)
}

func Serve_Ingress_EvalSummary(res http.ResponseWriter, req *http.Request) {
    id := int(req.PathValue("id"))
    summary, err := Ingress_EvalSummary(pb.IngressEvalRequest{ID: id})
    if err != nil {
        fmt.Fprint(res, err)
        return
    }
    fmt.Fprint(res, summary)
}

func main() {
    service_port := os.Getenv("SERVICE_PORT")
    if service_port == "" {
        service_port = "8081"
    }
    deployer_add = os.Getenv("DEP_CONTROLLER_ADD")
    if deployer_add == "" {
        deployer_add = "127.0.0.1:8082"
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
    default_config = pb.ConfigStruct{Namespace: namespace, TTL: ttl, UpdateThreshold: update_threshold, Deployment: deployment, Communication: communication, Aggregation: aggregation, TargetConcurrency: target_concurrency, Debug: debug}

    l, err := net.Listen("tcp", ":" + service_port)
    if err != nil {
        fmt.Println("listener not operational")
    }
    m := cmux.New(l)
    l_g := m.Match(cmux.HTTP2HeaderField("content-type", "application/grpc"))
    l_h := m.Match(cmux.HTTP1Fast())

    s_g := grpc.NewServer()
    m_h := http.NewServeMux()
    pb.RegisterMacroPodIngressServer(s_g, &IngressService{})
    m_h.HandleFunc("/", Serve_Ingress_Help)
    m_h.HandleFunc("/config", Serve_Ingress_Config)
    m_h.HandleFunc("/workflow/invoke/{wf_name}", Serve_Ingress_WorkflowInvoke)
    m_h.HandleFunc("/function/invoke/{func_name}", Serve_Ingress_FunctionInvoke)
    m_h.HandleFunc("/workflow/create/{wf_name}", Serve_Ingress_CreateWorkflow)
    m_h.HandleFunc("/workflow/update/{wf_name}", Serve_Ingress_UpdateWorkflow)
    m_h.HandleFunc("/workflow/delete/{wf_name}", Serve_Ingress_DeleteWorkflow)
    m_h.HandleFunc("/eval/start", Serve_Ingress_Eval)
    m_h.HandleFunc("/eval/metrics/{id}", Serve_Ingress_EvalMetrics)
    m_h.HandleFunc("/eval/latency/{id}", Serve_Ingress_EvalLatency)
    m_h.HandleFunc("/eval/summary/{id}", Serve_Ingress_EvalSummary)

    go watchTTL()
    go deployerSync()

    go s_g.Serve(l_g)
    s_h := &http.Server{Handler: m_h}
    go s_h.Serve(l_h)
}
