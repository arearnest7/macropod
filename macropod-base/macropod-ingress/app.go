package main

import (
    "net/http"
    "fmt"
    "math/rand"
    "bytes"
    "context"
    "time"
    "strconv"
    "strings"
    "io"
    "io/ioutil"
    "encoding/json"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"

    apiv1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/rest"
    "k8s.io/client-go/kubernetes"
    appsv1 "k8s.io/api/apps/v1"

    pb "app/app_pb"
)

type server struct {
    pb.UnimplementedGRPCFunctionServer
}

type Context struct {
    Request []byte
    WorkflowId string
    Depth int
    Width int
    RequestType string
    InvokeType string
    IsJson bool
}

type WFConf struct {
    workflow_name string
    func_names []string
    registries []string
    endpoints [][]string
    envs []map[string]string
    secrets []map[string]string
}

type container struct {
    func_name string
    registry string
    endpoints []string
    envs map[string]string
    secrets map[string]string
}

type pod struct {
    pod_name string
    containers map[int]*container
}

type deployment struct {
    name string
    workflow *workflow
    pods map[int]*pod
    config *appsv1.Deployment
    ttl time.Time
}

type workflow struct {
    name string
    containers map[string]*container
    pattern [][]string
    deployments map[int]*deployment
    entrypoint map[int]string
    in_use map[int]bool
}

var config *rest.Config
var client *kubernetes.Clientset
var wt map[string]*workflow

func Help(res http.ResponseWriter, req *http.Request) {
    help_print := "macropod ingress\nplease use the following services to interact/deploy/invoke workflows:\n"
    help_print += "Invoke:\n - path: /invoke/[workflow_name]\n - purpose: invoke workflow at entrypoint\n - payload: JSON body for entrypoint function\n - output: string result\n"
    help_print += "Deploy:\n - path: /deploy\n - purpose: deploy workflow based on configuration\n - payload: macropod YAML configuration\n - output: string confirmation\n"
    help_print += "Update:\n - path: /update\n - purpose: update previously deployed workflow\n - payload: macropod YAML configuration\n - output: string confirmation\n"
    help_print += "Delete:\n - path: /delete\n - purpose: delete workflow of same name\n - payload: workflow name string\n - output: string confirmation\n"
    help_print += "Logs:\n - path: /logs\n - purpose: output logs for a workflow\n - payload: workflow name string\n - output: CSV\n"
    help_print += "Metrics:\n - path: /metrics\n - purpose: output current metrics\n - payload: None\n - output: CSV\n"
    fmt.Fprintf(res, help_print)
}

func Metrics(res http.ResponseWriter, req *http.Request) {
    //TODO
    fmt.Fprintf(res, "metrics")
}

func Logs(res http.ResponseWriter, req *http.Request) {
    logs_arr := make(map[string]string)
    payload, _ := ioutil.ReadAll(req.Body)
    wf := wt[string(payload)]
    for func_name, _ := range wf.containers {
        logs_arr[func_name] = ""
    }
    for _, deployment := range wf.deployments {
        for _, pod := range deployment.pods {
            for _, container := range pod.containers {
                logOpts := apiv1.PodLogOptions{Container: container.func_name}
                req := client.CoreV1().Pods("default").GetLogs(pod.pod_name, &logOpts)
                log, _ := req.Stream(context.TODO())
                defer log.Close()
                b := new(bytes.Buffer)
                io.Copy(b, log)
                logs_arr[container.func_name] = b.String()
            }
        }
    }
    result := ""
    for func_name, log := range logs_arr {
        result += func_name + "\n" + log + "\n\n"
    }
    fmt.Fprintf(res, result)
}

func WF_Realloc() (string) {
    return ""
}

func Dep_Create(wf *workflow) (int, string) {
    return 0, ""
}

func Dep_Delete(workflow_name string, dep_idx int) {
    wf := wt[workflow_name]
    deletePolicy := metav1.DeletePropagationForeground
    client.AppsV1().Deployments("default").Delete(context.TODO(), wf.deployments[dep_idx].name, metav1.DeleteOptions{PropagationPolicy: &deletePolicy})
    delete(wf.deployments, dep_idx)
}

func Dep_Reclaim() {
    for true {
        time.Sleep(1 * time.Second)
        for workflow_name, workflow := range wt {
            for i, deployment := range workflow.deployments {
                if !workflow.in_use[i] {
                    if time.Now().Before(deployment.ttl) {
                        Dep_Delete(workflow_name, i)
                    }
                }
            }
        }
    }
}

func WF_Deploy(res http.ResponseWriter, req *http.Request) {
    var workflow_conf WFConf
    err := json.NewDecoder(req.Body).Decode(&workflow_conf)
    if err != nil {
        fmt.Fprintf(res, "Invalid Workflow JSON, please resubmit with proper formatting\n")
    } else {
        var wf workflow
        wf.name = workflow_conf.workflow_name
        for i, func_name := range workflow_conf.func_names {
            var c container
            c.func_name = func_name
            c.registry = workflow_conf.registries[i]
            c.endpoints = workflow_conf.endpoints[i]
            for name, env := range workflow_conf.envs[i] {
                c.envs[name] = env
            }
            for name, secret := range workflow_conf.secrets[i] {
                c.secrets[name] = secret
            }
            wf.containers[func_name] = &c
        }
        Dep_Create(&wf)
        wt[wf.name] = &wf
        fmt.Fprintf(res, "Deployment was successful. Invoke your workflow with /invoke/" + wf.name + "\n")
    }
}

func WF_Delete(res http.ResponseWriter, req *http.Request) {
    payload, _ := ioutil.ReadAll(req.Body)
    wf_name := string(payload)
    wf := wt[wf_name]
    for key, _ := range wf.deployments {
        Dep_Delete(wf_name, key)
    }
    for key, _ := range wf.containers {
        delete(wf.containers, key)
    }
    delete(wt, wf_name)
    fmt.Fprintf(res, wf_name + " has been deleted.\n")
}

func WF_Update(res http.ResponseWriter, req *http.Request) {
    var workflow_conf WFConf
    err := json.NewDecoder(req.Body).Decode(&workflow_conf)
    if err != nil {
        fmt.Fprintf(res, "Invalid Workflow JSON, please resubmit with proper formatting\n")
    } else {
        wf := wt[workflow_conf.workflow_name]
        for key, _ := range wf.deployments {
            Dep_Delete(workflow_conf.workflow_name, key)
        }
        for key, _ := range wf.containers {
            delete(wf.containers, key)
        }
        delete(wt, workflow_conf.workflow_name)
        wf = new(workflow)
        wf.name = workflow_conf.workflow_name
        for i, func_name := range workflow_conf.func_names {
            c := new(container)
            c.func_name = func_name
            c.registry = workflow_conf.registries[i]
            c.endpoints = workflow_conf.endpoints[i]
            for name, env := range workflow_conf.envs[i] {
                c.envs[name] = env
            }
            for name, secret := range workflow_conf.secrets[i] {
                c.secrets[name] = secret
            }
            wf.containers[func_name] = c
        }
        Dep_Create(wf)
        wt[wf.name] = wf
        fmt.Fprintf(res, "Update was successful. Invoke your workflow with /invoke/" + wf.name + "\n")
    }
}

func WF_Invoke(res http.ResponseWriter, req *http.Request) {
    wf_name := strings.Split(req.URL.Path, "/")[2]
    wf := wt[wf_name]
    idx := -1
    entrypoint := ""
    not_used_found:
    for i, v := range wf.in_use {
        if v {
            idx = i
            entrypoint = wf.entrypoint[i]
            wf.in_use[i] = true
            break not_used_found
        }
    }
    if idx == -1 {
        idx, entrypoint = Dep_Create(wf)
    }
    payload, _ := ioutil.ReadAll(req.Body)
    workflow_id := strconv.Itoa(rand.Intn(10000000))
    channel, _ := grpc.Dial(entrypoint, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*200), grpc.MaxCallSendMsgSize(1024*1024*200)), grpc.WithTransportCredentials(insecure.NewCredentials()))
    defer channel.Close()
    stub := pb.NewGRPCFunctionClient(channel)
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()
    request_type := "gg"
    response, _ := stub.GRPCFunctionHandler(ctx, &pb.RequestBody{Data: payload, WorkflowId: workflow_id, Depth: 0, Width: 0, RequestType: &request_type})
    results := response.GetReply()
    wf.in_use[idx] = false
    wf.deployments[idx].ttl = time.Now().Add(time.Minute * 30)
    fmt.Fprintf(res, results)
}

func main() {
    wt = make(map[string]*workflow)
    config, _ = rest.InClusterConfig()
    client, _ = kubernetes.NewForConfig(config)
    go Dep_Reclaim()
    http.HandleFunc("/", Help)
    http.HandleFunc("/invoke", WF_Invoke)
    http.HandleFunc("/deploy", WF_Deploy)
    http.HandleFunc("/update", WF_Update)
    http.HandleFunc("/delete", WF_Delete)
    http.HandleFunc("/logs", Logs)
    http.HandleFunc("/metrics", Metrics)
    http.ListenAndServe(":80", nil)
}
