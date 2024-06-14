package main

import (
    "time"
    "fmt"
    "bytes"
    "context"
    "encoding/json"
    "io"
    "math"
    "net"
    "os"
    "strconv"
    "strings"
    "google.golang.org/grpc"
    appsv1 "k8s.io/api/apps/v1"
    corev1 "k8s.io/api/core/v1"
    networkingv1 "k8s.io/api/networking/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/util/intstr"
    "k8s.io/apimachinery/pkg/watch"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/rest"
    pb "app/deployer_pb"
)

type server struct {
    pb.DeploymentServiceServer
}

type PodMetricsList struct {
    Kind string `json:"kind"`
    APIVersion string `json:"apiVersion"`
    Metadata MetadataPod `json:"metadata"`
    Items []PodMetricsItem `json:"items"`
}

type MetadataPod struct {
    Name string `json:"name,omitempty"`
    Namespace string `json:"namespace,omitempty"`
    CreationTimestamp time.Time `json:"creationTimestamp,omitempty"`
    Labels map[string]string `json:"labels,omitempty"`
}

type PodMetricsItem struct {
    Metadata MetadataPod `json:"metadata"`
    Timestamp time.Time `json:"timestamp"`
    Window string `json:"window"`
    Containers []ContainerMetric `json:"containers"`
}

type ContainerMetric struct {
    Name string `json:"name"`
    Usage Usage `json:"usage"`
}

type NodeMetricList struct {
    Kind string `json:"kind"`
    APIVersion string `json:"apiVersion"`
    Metadata struct {
        SelfLink string `json:"selfLink"`
    } `json:"metadata"`
    Items []NodeMetricsItem `json:"items"`
}

type NodeMetricsItem struct {
    Metadata Metadata `json:"metadata"`
    Timestamp time.Time `json:"timestamp"`
    Window string `json:"window"`
    Usage Usage `json:"usage"`
}
type Metadata struct {
    Name string `json:"name,omitempty"`
    CreationTimestamp time.Time `json:"creationTimestamp,omitempty"`
    Labels map[string]string `json:"labels,omitempty"`
}

type Usage struct {
    CPU string `json:"cpu"`
    Memory string `json:"memory"`
}

type ContainerMetrics struct {
    Name string `json:"name"`
    CPU string `json:"cpu"`
    Memory string `json:"memory"`
}

type PodMetrics struct {
    Name string `json:"name"`
    Namespace string `json:"namespace"`
    Containers []ContainerMetrics `json:"containers"`
}

type NodeMetrics struct {
    Name string `json:"name"`
    CPU string `json:"cpu"`
    Memory string `json:"memory"`
}

type Metrics struct {
    Pods []PodMetrics `json:"pods"`
    Nodes []NodeMetrics `json:"nodes"`
}

type Function struct {
    Registry string `json:"registry"`
    Endpoints []string `json:"endpoints,omitempty"`
    Envs map[string]string `json:"envs,omitempty"`
    Secrets map[string]string `json:"secrets,omitempty"`
}

type Workflow struct {
    Name string `json:"name,omitempty"`
    Functions map[string]Function `json:"functions"`
    Pods [][]string
    IngressVersion map[string]int
    LatestVersion int
    LastUpdated time.Time
    NextReplicaIndex int
    Updating bool
    FullyDisaggregated bool
}

var (
    kclient *kubernetes.Clientset
    namespace_ingress string
    workflows map[string]Workflow
    cpu_threshold_1 float64
    cpu_threshold_2 float64
    mem_threshold_1 float64
    mem_threshold_2 float64
    update_threshold int
)

func internal_log(message string) {
    fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.000000 UTC") + ": " + message)
}

func manageDeployment(wf_name string, ns string) (string, error) {
    workflow := workflows[wf_name]
    namespace := ns
    update := true
    if namespace == "" {
        namespace = wf_name + "-" + strconv.Itoa(workflow.NextReplicaIndex)
        workflow.NextReplicaIndex += 1
        update = false
    }
    if !update {
        namespace_object := &corev1.Namespace{
            ObjectMeta: metav1.ObjectMeta{
                Name: namespace,
            },
        }
        _, err := kclient.CoreV1().Namespaces().Create(context.Background(), namespace_object, metav1.CreateOptions{})
        if err != nil {
            internal_log("namespace " + namespace + " unable to be created - " + err.Error())
            return "", err
        }
        internal_log("namespace " + namespace + " has been created")
    }
    labels_ingress := map[string]string{
        "workflow_name": wf_name,
        "app": namespace,
    }
    pathType := networkingv1.PathTypePrefix
    ingress := &networkingv1.Ingress{
        ObjectMeta: metav1.ObjectMeta{
            Name: namespace,
            Namespace: namespace,
            Labels: labels_ingress,
        },
        Spec: networkingv1.IngressSpec{
            Rules: []networkingv1.IngressRule{
                {
                    Host: wf_name + "." + namespace + ".macropod",
                    IngressRuleValue: networkingv1.IngressRuleValue{
                        HTTP: &networkingv1.HTTPIngressRuleValue{
                            Paths: []networkingv1.HTTPIngressPath{
                                {
                                    Path: "/",
                                    PathType: &pathType,
                                    Backend: networkingv1.IngressBackend{
                                        Service: &networkingv1.IngressServiceBackend{
                                            Name: workflow.Pods[0][0],
                                            Port: networkingv1.ServiceBackendPort{
                                                Number: 5000,
                                            },
                                        },
                                    },
                                },
                            },
                        },
                    },
                },
            },
        },
    }
    if update {
        _, err := kclient.NetworkingV1().Ingresses(namespace).Update(context.Background(), ingress, metav1.UpdateOptions{})
        if err != nil {
            internal_log("unable to update ingress for " + namespace + " - " + err.Error())
            return "", err
        }
    } else {
        _, err := kclient.NetworkingV1().Ingresses(namespace).Create(context.Background(), ingress, metav1.CreateOptions{})
        if err != nil {
            internal_log("unable to create ingress for " + namespace + " - " + err.Error())
            return "", err
        }
    }
    for _, pod := range workflow.Pods {
        service := &corev1.Service{
            ObjectMeta: metav1.ObjectMeta{
                Name: pod[0],
                Namespace: namespace,
            },
            Spec: corev1.ServiceSpec{
                Selector: map[string]string{
                    "app": pod[0],
                },
                Ports: []corev1.ServicePort{
                    {
                        Name: pod[0],
                        Port: 5000,
                        TargetPort: intstr.FromInt(5000),
                    },
                },
            },
        }
        node_name := cpu_node_sort()
        labels := map[string]string{
            "workflow_name": wf_name,
            "app": pod[0],
        }
        replicaCount := int32(1)
        deployment := &appsv1.Deployment{
            ObjectMeta: metav1.ObjectMeta{
                Name: pod[0],
                Labels: labels,
            },
            Spec: appsv1.DeploymentSpec{
                Replicas: &replicaCount,
                Selector: &metav1.LabelSelector{
                    MatchLabels: labels,
                },
                Template: corev1.PodTemplateSpec{
                    ObjectMeta: metav1.ObjectMeta{
                        Labels: labels,
                    },
                    Spec: corev1.PodSpec{
                        NodeSelector: map[string]string{
                            "kubernetes.io/hostname": node_name,
                        },
                        Containers: make([]corev1.Container, 0),
                    },
                },
            },
        }
        func_port := 5000
        for i, container := range pod {
            function := workflow.Functions[container]
            registry := function.Registry
            var env []corev1.EnvVar
            for name, value := range function.Envs {
                env = append(env, corev1.EnvVar{Name: name, Value: value,})
            }
            env = append(env, corev1.EnvVar{Name: "SERVICE_TYPE", Value: "GRPC",})
            env = append(env, corev1.EnvVar{Name: "GRPC_THREAD", Value: "10",})
            func_port_s := strconv.Itoa(func_port)
            env = append(env, corev1.EnvVar{Name: "FUNC_PORT", Value: func_port_s,})
            for _, endpoint := range function.Endpoints {
                in_pod := false
                idx := 0
                for j, c := range pod {
                    if endpoint == c {
                        in_pod = true
                        idx = j
                        break
                    }
                }
                endpoint_upper := strings.ToUpper(endpoint)
                endpoint_name := strings.ReplaceAll(endpoint_upper, "-", "_")
                var service_name string
                if in_pod {
                    endpoint_port := strconv.Itoa(5000 + idx)
                    service_name = "127.0.0.1:" + endpoint_port
                } else {
                    service_name = endpoint + "." + namespace + ".svc.cluster.local"
                }
                env = append(env, corev1.EnvVar{Name: endpoint_name, Value: service_name,})
            }
            container_port := int32(func_port)
            func_port += 1
            imagePullPolicy := corev1.PullPolicy("Always")
            deployment.Spec.Template.Spec.Containers[i] = corev1.Container{
                Name: container,
                Image: registry,
                ImagePullPolicy: imagePullPolicy,
                Ports: []corev1.ContainerPort{
                    {
                        ContainerPort: container_port,
                    },
                },
                Env: env,
            }
        }
        _, err := kclient.CoreV1().Services(namespace).Update(context.Background(), service, metav1.UpdateOptions{})
        if err != nil {
            internal_log("unable to update service " + pod[0] + " for " + namespace + " - " + err.Error())
            return "", err
        }
        _, err = kclient.AppsV1().Deployments(namespace).Update(context.Background(), deployment, metav1.UpdateOptions{})
        if err != nil {
            internal_log("unable to update deployment " + pod[0] + " for " + namespace + " - " + err.Error())
            return "", err
        }
    }
    workflow.IngressVersion[namespace] = workflow.LatestVersion
    return wf_name + "." + namespace + ".macropod", nil
}

func updateDeployments(wf_name string) {
    if workflows[wf_name].Updating {
        return
    }
    workflow := workflows[wf_name]
    workflow.Updating = true
    var cpu_total float64
    var memory_total float64
    for _, wf := range workflows {
        for namespace, _ := range wf.IngressVersion {
            cpu, memory := getNamespaceMetrics(namespace)
            cpu_total += cpu
            memory_total += memory
        }
    }
    if cpu_total > cpu_threshold_1 || memory_total > mem_threshold_1 {
        internal_log("threshold 1 reached - " + wf_name)
        if time.Now().Sub(workflow.LastUpdated) > time.Second * time.Duration(update_threshold) && !workflow.FullyDisaggregated {
            workflow.LatestVersion += 1
            internal_log("workflow " + wf_name + " updated to version " + strconv.Itoa(workflow.LatestVersion))
            var pods_updated [][]string
            for _, pod := range workflow.Pods {
                if len(pod) > 1 {
                    idx := int(math.Floor(float64(len(pod)) / 2))
                    pods_updated = append(pods_updated, pod[:idx])
                    pods_updated = append(pods_updated, pod[idx+1:])
                } else {
                    pods_updated = append(pods_updated, pod)
                }
            }
            workflow.Pods = pods_updated
            pod_2_or_more := false
            for _, pod := range pods_updated {
                if len(pod) > 1 {
                    pod_2_or_more = true
                    break
                }
            }
            if !pod_2_or_more {
                internal_log(wf_name + " has been fully disaggregated")
                workflow.FullyDisaggregated = true
            }
            workflow.LastUpdated = time.Now()
        }
        if cpu_total > cpu_threshold_2 || memory_total > mem_threshold_2 {
            internal_log("threshold 2 reached - " + wf_name)
            for namespace, version := range workflow.IngressVersion {
                if version < workflow.LatestVersion {
                    go manageDeployment(wf_name, namespace)
                }
            }
        }
    }
    workflow.Updating = false
}

func memory_raw_to_float(memory_str string) (float64, error) {
    if memory_str == "0" {
        return 0, nil
    } else if strings.HasSuffix(memory_str, "Ki") {
        memory_str = strings.TrimSuffix(memory_str, "Ki")
        memory, err := strconv.ParseFloat(memory_str, 64)
        if err != nil {
            return 0, err
        }
        return memory * 1024, nil
    } else if strings.HasSuffix(memory_str, "Mi") {
        memory_str = strings.TrimSuffix(memory_str, "Mi")
        memory, err := strconv.ParseFloat(memory_str, 64)
        if err != nil {
            return 0, err
        }
        return memory * 1024 * 1024, nil
    }
    return 0, fmt.Errorf("unsupported memory usage format")
}

func cpu_raw_to_float(cpu_str string) (float64, error) {
    if cpu_str == "0" {
        return 0, nil
    } else if strings.HasSuffix(cpu_str, "n") {
        cpu_str = strings.TrimSuffix(cpu_str, "n")
        cpu, err := strconv.ParseFloat(cpu_str, 64)
        if err != nil {
            return 0, err
        }
        return cpu / 1000000, nil
    } else if strings.HasSuffix(cpu_str, "m") {
        cpu_str = strings.TrimSuffix(cpu_str, "m")
        cpu, err := strconv.ParseFloat(cpu_str, 64)
        if err != nil {
            return 0, err
        }
        return cpu, nil
    } else if strings.HasSuffix(cpu_str, "u") {
        cpu_str = strings.TrimSuffix(cpu_str, "m")
        cpu, err := strconv.ParseFloat(cpu_str, 64)
        if err != nil {
            return 0, err
        }
        return cpu / 1000, nil
    }
    return 0, fmt.Errorf("unsupported CPU usage format %s", cpu_str)
}

func cpu_node_sort() (string) {
    internal_log("SORT_CPU_START")
    var nodes NodeMetricList
    data, err := kclient.RESTClient().Get().AbsPath("apis/metrics.k8s.io/v1beta1/nodes").Do(context.TODO()).Raw()
    if err != nil {
        internal_log("unable to retrieve metrics from nodes API - " + err.Error())
        return ""
    }
    err = json.Unmarshal(data, &nodes)
    if err != nil {
        internal_log("unable to unmarshal metrics from nodes API - " + err.Error())
        return ""
    }
    node_name := ""
    var node_usage_minimum float64 = math.Inf(1)
    for _, item := range nodes.Items {
        cpu_current, err := cpu_raw_to_float(item.Usage.CPU)
        if err != nil {
            internal_log("unable to convert cpu to float - " + err.Error())
            return ""
        }
        if cpu_current < node_usage_minimum {
            node_usage_minimum = cpu_current
            node_name = item.Metadata.Name
        }
    }
    internal_log("SORT_CPU_END")
    return node_name
}

func getNamespaceMetrics(namespace string) (float64, float64) {
    internal_log("GET_NAMESPACE_METRICS_START " + namespace)
    _, exists := kclient.CoreV1().Namespaces().Get(context.Background(), namespace, metav1.GetOptions{})
    if exists != nil {
        internal_log("namespace " + namespace + " does not exists. Failed to get pods for metrics.")
        return 0, 0
    }
    var podMetricsList PodMetricsList
    data, err := kclient.RESTClient().Get().AbsPath("apis/metrics.k8s.io/v1beta1/pods").Do(context.TODO()).Raw()
    if err != nil {
        internal_log("unable to retrieve metrics from API - " + err.Error())
        return 0, 0
    }
    err = json.Unmarshal(data, &podMetricsList)
    if err != nil {
        internal_log("unable to unmarshal metrics from API - " + err.Error())
        return 0, 0
    }
    cpu_ns := 0.0
    memory_ns := 0.0
    for _, item := range podMetricsList.Items {
        podName := item.Metadata.Name
        podNamespace := item.Metadata.Namespace
        if podNamespace != namespace {
            continue
        }
        var cpu_p float64
        var memory_p float64
        for _, container := range item.Containers {
            cpu, err := cpu_raw_to_float(container.Usage.CPU)
            if err != nil {
                internal_log("Error converting memory usage for pod " + podName + " - " + err.Error())
                return 0, 0
            }
            memory, err := memory_raw_to_float(container.Usage.Memory)
            if err != nil {
                internal_log("Error converting CPU usage for pod " + podName + " - " + err.Error())
                return 0, 0
            }
            cpu_p += cpu
            memory_p += memory
        }
        cpu_ns += cpu_p
        memory_ns += memory_p
    }
    internal_log("GET_NAMESPACE_METRICS_END " + namespace)
    return cpu_ns, memory_ns
}

func dfs_initial_pod(entrypoint string, pod []string, visited map[string]bool, wf_name string) {
    for _, endpoint := range workflows[wf_name].Functions[entrypoint].Endpoints {
        _, exists := visited[endpoint]
        if !exists {
            pod = append(pod, endpoint)
            visited[endpoint] = true
            dfs_initial_pod(endpoint, pod, visited, wf_name)
        }
    }
}

func createInitialPod(wf_name string) {
    workflow := workflows[wf_name]
    var endpoints map[string]string
    for _, function := range workflow.Functions {
        for _, endpoint := range function.Endpoints {
            _, exists := endpoints[endpoint]
            if !exists {
                endpoints[endpoint] = endpoint
            }
        }
    }
    var initial_pod []string
    var visited map[string]bool
    for func_name, _ := range workflow.Functions {
        _, exists := endpoints[func_name]
        if !exists {
            initial_pod = append(initial_pod, func_name)
            visited[func_name] = true
            break
        }
    }
    dfs_initial_pod(initial_pod[0], initial_pod, visited, wf_name)
    workflow.Pods = append(workflow.Pods, initial_pod)
}

func createWorkflow(wf_name string, workflow_str string) {
    internal_log("CREATE_WORKFLOW_START - " + wf_name)
    workflow := Workflow{}
    json.Unmarshal([]byte(workflow_str), &workflow)
    _, exists := workflows[wf_name]
    if exists {
        internal_log("workflow " + wf_name + " already exists. If you are updating it please use update instead.")
        return
    }
    workflows[wf_name] = workflow
    createInitialPod(wf_name)
    internal_log("CREATE_WORKFLOW_END - " + wf_name)
}

func updateWorkflow(wf_name string, workflow_str string) {
    internal_log("UPDATE_WORKFLOW_START - " + wf_name)
    workflow := Workflow{}
    json.Unmarshal([]byte(workflow_str), &workflow)
    existing_workflow, exists := workflows[wf_name]
    if exists {
        for namespace, _ := range existing_workflow.IngressVersion {
            kclient.CoreV1().Namespaces().Delete(context.TODO(), namespace, metav1.DeleteOptions{})
        }
        delete(workflows, wf_name)
    }
    workflows[wf_name] = workflow
    createInitialPod(wf_name)
    internal_log("UPDATE_WORKFLOW_END - " + wf_name)
}

func deleteWorkflow(wf_name string) {
    internal_log("DELETE_WORKFLOW_START - " + wf_name)
    existing_workflow, exists := workflows[wf_name]
    if exists {
        for namespace, _ := range existing_workflow.IngressVersion {
            kclient.CoreV1().Namespaces().Delete(context.TODO(), namespace, metav1.DeleteOptions{})
        }
        delete(workflows, wf_name)
    } else {
        internal_log("workflow " + wf_name + " does not exist. Please check the spelling of your workflow name.")
        return
    }
    internal_log("DELETE_WORKFLOW_END - " + wf_name)
}

func getLogs(wf_name string) (string) {
    internal_log("GET_LOGS_START - " + wf_name)
    logs_arr := make(map[string]string)
    workflow, exists := workflows[wf_name]
    if !exists {
        internal_log("workflow + " + wf_name + " does not exists for logs. Please check workflow name spelling.")
        return ""
    }
    for namespace, _ := range workflow.IngressVersion {
        pods, err := kclient.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{})
        if err != nil {
            internal_log("namespace " + namespace + " does not exists. Failed to get pods for logs.")
        }
        for _, pod := range pods.Items {
            for _, func_name := range pod.Spec.Containers {
                podLogOpts := corev1.PodLogOptions {
                    Container: func_name.Name,
                }
                req := kclient.CoreV1().Pods(namespace).GetLogs(pod.Name, &podLogOpts)
                logs, err := req.Stream(context.TODO())
                if err != nil {
                    internal_log("error in opening stream for " + namespace + " " + pod.Name + " " + func_name.Name)
                }
                defer logs.Close()
                b := new(bytes.Buffer)
                io.Copy(b, logs)
                _, exists := logs_arr[func_name.Name]
                if exists {
                    logs_arr[func_name.Name] += b.String()
                } else {
                    logs_arr[func_name.Name] = b.String()
                }
            }
        }
    }
    result := ""
    for func_name, logs := range logs_arr {
        result += func_name + "\n" + logs + "\n"
    }
    internal_log("GET_LOGS_END - " + wf_name)
    return result
}

func getMetrics() (string) {
    internal_log("GET_METRICS_START")
    var cpu_total float64
    var memory_total float64
    for _, workflow := range workflows {
        for namespace, _ := range workflow.IngressVersion {
            cpu, memory := getNamespaceMetrics(namespace)
            cpu_total += cpu
            memory_total += memory
        }
    }
    internal_log("GET_METRICS_END")
    return strconv.FormatFloat(cpu_total, 'f', -1, 64) + ", " + strconv.FormatFloat(memory_total, 'f', -1, 64)
}

func updateExistingIngress(wf_name string) {
    internal_log("UPDATE_EXISTING_START - " + wf_name)
    updateDeployments(wf_name)
    internal_log("UPDATE_EXISTING_END - " + wf_name)
}

func createNewIngress(wf_name string) (string) {
    internal_log("CREATE_INGRESS_START - " + wf_name)
    _, exist := workflows[wf_name]
    if !exist {
        internal_log("unable to create new ingress for " + wf_name + " - workflow does not exist")
        return ""
    }
    ingress, err := manageDeployment(wf_name, "")
    if err != nil {
        internal_log("Failed to deploy new ingress - " + err.Error())
        return ""
    }
    internal_log("CREATE_INGRESS_END - " + wf_name)
    return ingress
}

func updateDeletedNamespace(namespace string) {
    internal_log("deleting namespace - " + namespace)
    for _, workflow := range workflows {
        for ns, _ := range workflow.IngressVersion {
            if namespace == ns {
                delete(workflow.IngressVersion, ns)
                internal_log("deleted namespace - " + namespace)
                return
            }
        }
    }
    internal_log("namespace not found - " + namespace)
}

func watchNamespaces() {
    internal_log("WATCH_NAMESPACE_START")
    for {
        watcher, err := kclient.CoreV1().Namespaces().Watch(context.TODO(), metav1.ListOptions{})
        if err != nil {
            internal_log("Failed to set up watch - " + err.Error())
        }
        for event := range watcher.ResultChan() {
            namespace, ok := event.Object.(*corev1.Namespace)
            if !ok {
                continue
            }
            switch event.Type {
                case watch.Deleted:
                    updateDeletedNamespace(namespace.Name)
            }
        }
    }
    internal_log("WATCH_NAMESPACE_END")
}

func (s *server) Deployment(ctx context.Context, req *pb.DeploymentServiceRequest) (*pb.DeploymentServiceReply, error) {
    wf_name := req.WorkflowName
    request_type := req.RequestType
    var result string
    if request_type == "create" {
        internal_log("create workflow request start - " + wf_name)
        createWorkflow(wf_name, *req.Data)
        internal_log("create workflow request end - " + wf_name)
    } else if request_type == "update" {
        internal_log("update workflow request start - " + wf_name)
        updateWorkflow(wf_name, *req.Data)
        internal_log("update workflow request end - " + wf_name)
    } else if request_type == "delete" {
        internal_log("delete workflow request start - " + wf_name)
        deleteWorkflow(wf_name)
        internal_log("delete workflow request end - " + wf_name)
    } else if request_type == "logs" {
        internal_log("logs request start - " + wf_name)
        result = getLogs(wf_name)
        internal_log("logs request end - " + wf_name)
    } else if request_type == "metrics" {
        internal_log("metrics request start - " + wf_name)
        result = getMetrics()
        internal_log("metrics request end - " + wf_name)
    } else if request_type == "existing_invoke" {
        internal_log("existing invoke request start - " + wf_name)
        updateExistingIngress(wf_name)
        internal_log("existing invoke request end - " + wf_name)
    } else if request_type == "new_invoke" {
        internal_log("new invoke request start - " + wf_name)
        result = createNewIngress(wf_name)
        internal_log("new invoke request end - " + wf_name)
    }
    return &pb.DeploymentServiceReply{
        Message: fmt.Sprintf("%s", result),
    }, nil
}

func main() {
    internal_log("Ingress Controller Started")
    workflows = make(map[string]Workflow, 0)
    config, err := rest.InClusterConfig()
    if err != nil {
        internal_log("error config - " + err.Error())
        return
    }
    kclient, err = kubernetes.NewForConfig(config)
    if err != nil {
        internal_log("error kclient - " + err.Error())
        return
    }
    namespace_ingress = os.Getenv("NAMESPACE_INGRESS")
    cpu_threshold_1, err = strconv.ParseFloat(os.Getenv("CPU_THRESHOLD_1"), 64)
    if err != nil {
        internal_log("error cpu_threshold_1 - " + err.Error())
        return
    }
    cpu_threshold_2, err = strconv.ParseFloat(os.Getenv("CPU_THRESHOLD_2"), 64)
    if err != nil {
        internal_log("error cpu_threshold_2 - " + err.Error())
        return
    }
    mem_threshold_1, err = strconv.ParseFloat(os.Getenv("MEM_THRESHOLD_1"), 64)
    if err != nil {
        internal_log("error mem_threshold_1 - " + err.Error())
        return
    }
    mem_threshold_2, err = strconv.ParseFloat(os.Getenv("MEM_THRESHOLD_2"), 64)
    if err != nil {
        internal_log("error mem_threshold_2 - " + err.Error())
        return
    }
    update_threshold, err = strconv.Atoi(os.Getenv("UPDATE_THRESHOLD"))
    if err != nil {
        internal_log("error update_threshold - " + err.Error())
        return
    }
    go watchNamespaces()
    port, err := strconv.Atoi(os.Getenv("SERVICE_PORT"))
    if err != nil {
        internal_log("error port - " + err.Error())
        return
    }
    l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
    if err != nil {
        internal_log("error listener - " + err.Error())
        return
    }
    s := grpc.NewServer()
    pb.RegisterDeploymentServiceServer(s, &server{})
    if err := s.Serve(l); err != nil {
        internal_log("failed to serve - " + err.Error())
        return
    }
}
