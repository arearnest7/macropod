package main

import (
    pb "app/macropod_pb"

    "os"
    "fmt"
    "strconv"
    "encoding/json"
    "io/ioutil"
    "strings"
    "time"
    "slices"
    "sync"
    "math"

    appsv1 "k8s.io/api/apps/v1"
    corev1 "k8s.io/api/core/v1"
    networkingv1 "k8s.io/api/networking/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/util/intstr"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/rest"

    "net"
    "net/http"

    "google.golang.org/grpc"
    "golang.org/x/net/context"
)

type DeployerService struct {
    pb.UnimplementedMacroPodDeployerServer
}

type NodeMetricList struct {
    Kind             string              `json:"kind"`
    APIVersion       string              `json:"apiVersion"`
    Metadata         struct {
                     SelfLink string     `json:"selfLink"`
    }                                    `json:"metadata"`
    Items            []NodeMetricsItem   `json:"items"`
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

var (
    kclient                      *kubernetes.Clientset

    nodeCapacityCPU              = make(map[string]float64)
    nodeCapacityMemory           = make(map[string]float64)

    workflows                    = make(map[string]*pb.WorkflowStruct)
    workflow_pods                = make(map[string][][]string)
    workflow_deployments         = make(map[string]map[string]map[string]string)
    workflow_latest_version      = make(map[string]int)
    workflow_last_updated        = make(map[string]time.Time)
    workflow_replica_count       = make(map[string]int)
    workflow_updating            = make(map[string]bool)
    workflow_initial_pods        = make(map[string][]string)
    workflow_fully_disaggregated = make(map[string]bool)

    default_config               *pb.ConfigStruct
    ingress_address              string
    logger_address               string

    dataLock                     sync.Mutex
)

func Debug(message string, debug_level int) {
    if default_config.GetDebug() > int32(debug_level) {
        fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.000000 UTC") + ": " + message + "\n")
    }
}

func bfs_initial_pod(pod []string, workflow_name string, pod_list []string) []string {
    if len(pod_list) == 0 {
        return pod
    }
    entrypoint := pod_list[0]
    if !slices.Contains(pod, entrypoint) {
        pod = append(pod, entrypoint)
    }
    pod_list = pod_list[1:]
    for _, endpoint := range workflows[workflow_name].Functions[entrypoint].Endpoints {
        if !slices.Contains(pod, endpoint) {
            pod = append(pod, endpoint)
            pod_list = append(pod_list, endpoint)

        }
    }
    return bfs_initial_pod(pod, workflow_name, pod_list)
}

func CreateInitialPod(workflow_name string) {
    var initial_pod []string
    var frontend_func string
    var endpoints []string
    func_endpoint := make(map[string][]string)
    for func_name, function := range workflows[workflow_name].Functions {
        for _, endpoint := range function.Endpoints {
            func_endpoint[func_name] = append(func_endpoint[func_name], endpoint)
            if !slices.Contains(endpoints, endpoint) {
                if func_name != endpoint {
                    endpoints = append(endpoints, endpoint)
                }
            }
        }
    }
    for func_name, _ := range workflows[workflow_name].Functions {
        if !slices.Contains(endpoints, func_name) {
            frontend_func = func_name
            break
        }

    }
    var pod_list []string
    pod_list = append(pod_list, frontend_func)
    initial_pod = bfs_initial_pod(initial_pod, workflow_name, pod_list)
    aggregation := default_config.GetAggregation()
    if workflows[workflow_name].GetConfig() != nil && workflows[workflow_name].GetConfig().GetAggregation() != "" {
        aggregation = workflows[workflow_name].GetConfig().GetAggregation()
    }
    if aggregation == "agg" {
        workflow_pods[workflow_name] = make([][]string, 0)
        workflow_pods[workflow_name] = append(workflow_pods[workflow_name], initial_pod)
    } else if aggregation == "disagg" {
        pods_disagg := make([][]string, 0)
        for _, container := range initial_pod {
            container_pod := make([]string, 0)
            container_pod = append(container_pod, container)
            pods_disagg = append(pods_disagg, container_pod)
        }
        workflow_pods[workflow_name] = pods_disagg
    } else {
        workflow_pods[workflow_name] = make([][]string, 0)
        workflow_pods[workflow_name] = append(workflow_pods[workflow_name], initial_pod)
    }
    workflow_initial_pods[workflow_name] = initial_pod
    workflow_latest_version[workflow_name] = 1
}

// return in bytes
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
    } else if strings.HasSuffix(memory_str, "Gi") {
        memory_str = strings.TrimSuffix(memory_str, "Gi")
        memory, err := strconv.ParseFloat(memory_str, 64)
        if err != nil {
            return 0, err
        }
        return memory * 1024 * 1024 * 1024, nil
    } else {
        memory, err := strconv.ParseFloat(memory_str, 64)
        if err != nil {
            return 0, err
        }
        return memory, nil
    }

}

// store cpu in cores
func cpu_raw_to_float(cpu_str string) (float64, error) {
    if cpu_str == "0" {
        return 0, nil
    } else if strings.HasSuffix(cpu_str, "n") {
        cpu_str = strings.TrimSuffix(cpu_str, "n")
        cpu, err := strconv.ParseFloat(cpu_str, 64)
        if err != nil {
            return 0, err
        }
        return cpu / 1000000000, nil
    } else if strings.HasSuffix(cpu_str, "m") {
        cpu_str = strings.TrimSuffix(cpu_str, "m")
        cpu, err := strconv.ParseFloat(cpu_str, 64)
        if err != nil {
            return 0, err
        }
        return cpu / 1000, nil
    } else if strings.HasSuffix(cpu_str, "u") {
        cpu_str = strings.TrimSuffix(cpu_str, "u")
        cpu, err := strconv.ParseFloat(cpu_str, 64)
        if err != nil {
            return 0, err
        }
        return cpu / 1000000, nil
    } else {
        cpu, err := strconv.ParseFloat(cpu_str, 64)
        if err != nil {
            return 0, err
        }
        return cpu, nil
    }

}

func CheckNodeStatus() {
    for {
        nodes, err := kclient.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
        if err != nil {
            panic(err)
        }
        for _, node := range nodes.Items {

            notReady := false
            for _, condition := range node.Status.Conditions {
                if condition.Type == "Ready" && condition.Status != "True" {
                    dataLock.Lock()
                    delete(nodeCapacityCPU, node.Name)
                    delete(nodeCapacityMemory, node.Name)
                    dataLock.Unlock()
                    notReady = true
                }
            }
            if notReady {
                continue
            }

            cpu_float, err := cpu_raw_to_float(node.Status.Capacity.Cpu().String())
            if err != nil {
                panic(err)
            }
            mem_float, err := memory_raw_to_float(node.Status.Capacity.Memory().String())
            if err != nil {
                panic(err)
            }
            dataLock.Lock()
            nodeCapacityCPU[node.Name] = cpu_float
            nodeCapacityMemory[node.Name] = mem_float
            dataLock.Unlock()
        }
        time.Sleep(30 * time.Second)
    }
}

func NodeReclaim(workflow_name string) {
    dataLock.Lock()
    reclaim_replicas := workflow_deployments[workflow_name]
    namespace := default_config.GetNamespace()
    if workflows[workflow_name].GetConfig() != nil && workflows[workflow_name].GetConfig().GetNamespace() != "" {
        namespace = workflows[workflow_name].GetConfig().GetNamespace()
    }
    dataLock.Unlock()
    for replica_name, _ := range reclaim_replicas {
        dataLock.Lock()
        delete(workflow_deployments[workflow_name], replica_name)
        dataLock.Unlock()
        labels_replica := "workflow_replica=" + replica_name
        services, err := kclient.CoreV1().Services(namespace).List(context.Background(), metav1.ListOptions{LabelSelector: labels_replica})
        if err != nil {
            Debug(err.Error(), 0)
        }
        for _, service := range services.Items {
            kclient.CoreV1().Services(namespace).Delete(context.Background(), service.ObjectMeta.Name, metav1.DeleteOptions{})
        }
        deployments, err := kclient.AppsV1().Deployments(namespace).List(context.Background(), metav1.ListOptions{LabelSelector: labels_replica})
        if err != nil {
            Debug(err.Error(), 0)
        }
        for _, deployment := range deployments.Items {
            kclient.AppsV1().Deployments(namespace).Delete(context.Background(), deployment.ObjectMeta.Name, metav1.DeleteOptions{})
        }
        ingresses, err := kclient.NetworkingV1().Ingresses(namespace).List(context.Background(), metav1.ListOptions{LabelSelector: labels_replica})
        if err != nil {
            Debug(err.Error(), 0)
        }
        for _, ingress := range ingresses.Items {
            kclient.NetworkingV1().Ingresses(namespace).Delete(context.Background(), ingress.ObjectMeta.Name, metav1.DeleteOptions{})
        }

        for {
            deployments_list, _ := kclient.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{LabelSelector: labels_replica})
            if deployments_list == nil || len(deployments_list.Items) == 0 {
                break
            }
            time.Sleep(10 * time.Millisecond)
        }
        request := pb.MacroPodRequest{Workflow: &workflow_name}
        Serve_CreateDeployment(&request, true)
    }
    dataLock.Lock()
    workflow_updating[workflow_name] = false
    dataLock.Unlock()
}

func Serve_Config(request *pb.ConfigStruct) (string) {
    dataLock.Lock()
    default_config = request
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

func Serve_CreateWorkflow(request *pb.WorkflowStruct) (string) {
    dataLock.Lock()
    _, exists := workflows[request.GetName()]
    if exists {
        dataLock.Unlock()
        return ""
    }
    workflows[request.GetName()] = request
    workflow_deployments[request.GetName()] = make(map[string]map[string]string)
    CreateInitialPod(request.GetName())
    dataLock.Unlock()
    return workflow_pods[request.GetName()][0][0]
}

func Serve_UpdateWorkflow(request *pb.WorkflowStruct) (string) {
    delete_request := pb.MacroPodRequest{Workflow: &request.Name}
    Serve_DeleteWorkflow(&delete_request)
    function := Serve_CreateWorkflow(request)
    return function
}

func Serve_DeleteWorkflow(request *pb.MacroPodRequest) (string) {
    dataLock.Lock()
    _, exists := workflows[request.GetWorkflow()]
    dataLock.Unlock()
    if exists {
        Debug("deleting workflow " + request.GetWorkflow(), 2)
        label_workflow := "workflow_name=" + request.GetWorkflow()
        namespace := default_config.GetNamespace()
        if workflows[request.GetWorkflow()].GetConfig() != nil && workflows[request.GetWorkflow()].GetConfig().GetNamespace() != "" {
            namespace = workflows[request.GetWorkflow()].GetConfig().GetNamespace()
        }
        services, err := kclient.CoreV1().Services(namespace).List(context.Background(), metav1.ListOptions{LabelSelector: label_workflow})
        if err != nil {
            Debug(err.Error(), 0)
        }
        for _, service := range services.Items {
            kclient.CoreV1().Services(namespace).Delete(context.Background(), service.ObjectMeta.Name, metav1.DeleteOptions{})
        }
        deployments, err := kclient.AppsV1().Deployments(namespace).List(context.Background(), metav1.ListOptions{LabelSelector: label_workflow})
        if err != nil {
            Debug(err.Error(), 0)
        }
        for _, deployment := range deployments.Items {
            kclient.AppsV1().Deployments(namespace).Delete(context.Background(), deployment.ObjectMeta.Name, metav1.DeleteOptions{})
        }
        ingresses, err := kclient.NetworkingV1().Ingresses(namespace).List(context.Background(), metav1.ListOptions{LabelSelector: label_workflow})
        if err != nil {
            Debug(err.Error(), 0)
        }
        for _, ingress := range ingresses.Items {
            kclient.NetworkingV1().Ingresses(namespace).Delete(context.Background(), ingress.ObjectMeta.Name, metav1.DeleteOptions{})
        }
        for {
            Debug("waiting for " + request.GetWorkflow() + " to delete", 4)
            deployments_list, _ := kclient.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{LabelSelector: label_workflow})
            if deployments_list == nil || len(deployments_list.Items) == 0 {
                break
            }
            time.Sleep(10 * time.Millisecond)
        }
        dataLock.Lock()
        delete(workflows, request.GetWorkflow())
        delete(workflow_pods, request.GetWorkflow())
        delete(workflow_deployments, request.GetWorkflow())
        delete(workflow_latest_version, request.GetWorkflow())
        delete(workflow_last_updated, request.GetWorkflow())
        delete(workflow_replica_count, request.GetWorkflow())
        delete(workflow_updating, request.GetWorkflow())
        delete(workflow_initial_pods, request.GetWorkflow())
        delete(workflow_fully_disaggregated, request.GetWorkflow())
        dataLock.Unlock()
    }
    return "deleted workflow " + request.GetWorkflow() + "\n"
}

func Serve_UpdateDeployments(request *pb.MacroPodRequest) (string) {
    dataLock.Lock()
    if _, exists := workflows[request.GetWorkflow()]; !exists {
        Debug("workflow is not present", 3)
        dataLock.Unlock()
        return "0"
    }
    if workflow_updating[request.GetWorkflow()] {
        dataLock.Unlock()
        return "1"
    }
    ttl := default_config.GetTTL()
    if workflows[request.GetWorkflow()].GetConfig() != nil && workflows[request.GetWorkflow()].GetConfig().GetTTL() != 0 {
        ttl = workflows[request.GetWorkflow()].GetConfig().GetTTL()
    }
    if time.Since(workflow_last_updated[request.GetWorkflow()]) < time.Second*time.Duration(ttl) {
        dataLock.Unlock()
        return "0"
    }
    aggregation := default_config.GetAggregation()
    if workflows[request.GetWorkflow()].GetConfig() != nil && workflows[request.GetWorkflow()].GetConfig().GetAggregation() != "" {
        aggregation = workflows[request.GetWorkflow()].GetConfig().GetAggregation()
    }
    dataLock.Unlock()
    switch aggregation {
        case "dynamic":
            dataLock.Lock()
            if workflow_fully_disaggregated[request.GetWorkflow()] {
                dataLock.Unlock()
                return "0"
            }
            workflow_updating[request.GetWorkflow()] = true
            dataLock.Unlock()
            // get the percentage of the node utilisation
            var nodes NodeMetricList
            data, err := kclient.RESTClient().Get().AbsPath("apis/metrics.k8s.io/v1beta1/nodes").Do(context.Background()).Raw()
            if err != nil {
                dataLock.Unlock()
                return "0"
            }
            err = json.Unmarshal(data, &nodes)
            if err != nil {
                dataLock.Unlock()
                return "0"
            }
            for _, node := range nodes.Items {
                value, exists := node.Metadata.Labels["node-role.kubernetes.io/master"]
                if exists && value == "true" {
                    continue
                }
                dataLock.Lock()
                node_name := node.Metadata.Name
                usage_mem, _ := memory_raw_to_float(node.Usage.Memory)
                usage_cpu, _ := cpu_raw_to_float(node.Usage.CPU)
                percentage_cpu := (usage_cpu / nodeCapacityCPU[node_name]) * 100
                percentage_mem := (usage_mem / nodeCapacityMemory[node_name]) * 100
                dataLock.Unlock()
                if percentage_mem > 70 || percentage_cpu > 70 {
                    dataLock.Lock()
                    workflow_latest_version[request.GetWorkflow()] += 1
                    var pods_updated [][]string
                    for _, pod := range workflow_pods[request.GetWorkflow()] {
                        if len(pod) > 1 {
                            idx := int(math.Floor(float64(len(pod)) / 2))
                            pods_updated = append(pods_updated, pod[:idx])
                            pods_updated = append(pods_updated, pod[idx:])
                        } else {
                            pods_updated = append(pods_updated, pod)
                        }
                    }

                    pod_2_or_more := false
                    for _, pod := range pods_updated {
                        if len(pod) > 1 {
                            pod_2_or_more = true
                            break
                        }
                    }
                    if !pod_2_or_more {
                        workflow_fully_disaggregated[request.GetWorkflow()] = true
                    }
                    workflow_pods[request.GetWorkflow()] = pods_updated
                    workflow_last_updated[request.GetWorkflow()] = time.Now()
                    go NodeReclaim(request.GetWorkflow())
                    dataLock.Unlock()
                    return "2"
                }
            }

            dataLock.Lock()
            if _, exists := workflows[request.GetWorkflow()]; exists {
                workflow_updating[request.GetWorkflow()] = false
            }
            dataLock.Unlock()
            return "0"
        default:
            return "0"
    }
}

func Serve_CreateDeployment(request *pb.MacroPodRequest, bypass bool) (string) {
    if !bypass {
        dataLock.Lock()
        if _, exists := workflows[request.GetWorkflow()]; !exists {
            dataLock.Unlock()
            return "0"
        }
        if workflow_updating[request.GetWorkflow()] {
            dataLock.Unlock()
            return "1"
        }
        dataLock.Unlock()
    }
    namespace := default_config.GetNamespace()
    if workflows[request.GetWorkflow()].GetConfig() != nil && workflows[request.GetWorkflow()].GetConfig().GetNamespace() != "" {
        namespace = workflows[request.GetWorkflow()].GetConfig().GetNamespace()
    }
    deployment := default_config.GetDeployment()
    if workflows[request.GetWorkflow()].GetConfig() != nil && workflows[request.GetWorkflow()].GetConfig().GetDeployment() != "" {
        deployment = workflows[request.GetWorkflow()].GetConfig().GetDeployment()
    }
    communication := default_config.GetCommunication()
    if workflows[request.GetWorkflow()].GetConfig() != nil && workflows[request.GetWorkflow()].GetConfig().GetCommunication() != "" {
        communication = workflows[request.GetWorkflow()].GetConfig().GetCommunication()
    }
    switch deployment {
        case "macropod":
            dataLock.Lock()
            var nodes NodeMetricList
            data, err := kclient.RESTClient().Get().AbsPath("apis/metrics.k8s.io/v1beta1/nodes").Do(context.Background()).Raw()
            if err != nil {
                dataLock.Unlock()
                return "0"
            }
            err = json.Unmarshal(data, &nodes)
            if err != nil {
                dataLock.Unlock()
                return "0"
            }
            start_node := -1
            for i, node := range nodes.Items {
                value, exists := node.Metadata.Labels["node-role.kubernetes.io/master"]
                if exists && value == "true" {
                    continue
                }
                in_use := false
                for _, deployment := range workflow_deployments[request.GetWorkflow()] {
                    if deployment[strings.ToLower(strings.ReplaceAll(workflow_pods[request.GetWorkflow()][0][0], "_", "-"))] == node.Metadata.Name {
                        in_use = true
                        break
                    }
                }
                if !in_use {
                    start_node = i
                    break
                }
            }
            if start_node == -1 {
                dataLock.Unlock()
                return "0"
            }
            replicaNumber := strconv.Itoa(workflow_replica_count[request.GetWorkflow()])
            workflow_replica_count[request.GetWorkflow()]++
            var update_deployments []appsv1.Deployment
            labels_ingress := map[string]string{
                "workflow_name":    request.GetWorkflow(),
                "workflow_replica": request.GetWorkflow() + "-" + replicaNumber,
                "version":          strconv.Itoa(workflow_latest_version[request.GetWorkflow()]),
            }
            pathType := networkingv1.PathTypePrefix
            service_name_ingress := strings.ToLower(strings.ReplaceAll(workflow_pods[request.GetWorkflow()][0][0], "_", "-")) + "-" + replicaNumber
            workflow_deployments[request.GetWorkflow()][request.GetWorkflow() + "-" + replicaNumber] = make(map[string]string)
            node_index := start_node
            for _, pod := range workflow_pods[request.GetWorkflow()] {
                node_sel := ""
                for node_sel == "" {
                    value, exists := nodes.Items[node_index].Metadata.Labels["node-role.kubernetes.io/master"]
                    if exists && value == "true" {
                        node_index = (node_index + 1) % len(nodes.Items)
                        continue
                    } else {
                        node_sel = nodes.Items[node_index].Metadata.Name
                        node_index = (node_index + 1) % len(nodes.Items)
                        break
                    }
                }
                pod_name := strings.ToLower(strings.ReplaceAll(pod[0], "_", "-"))
                workflow_deployments[request.GetWorkflow()][request.GetWorkflow() + "-" + replicaNumber][pod_name] = node_sel
                service := &corev1.Service{
                    ObjectMeta: metav1.ObjectMeta{
                        Name:      pod_name + "-" + replicaNumber,
                        Namespace: namespace,
                    },
                    Spec: corev1.ServiceSpec{
                        Selector: map[string]string{
                            "app": pod_name + "-" + replicaNumber,
                        },
                        Ports: []corev1.ServicePort{},
                    },
                }
                labels := map[string]string{
                    "workflow_name":    request.GetWorkflow(),
                    "app":              pod_name + "-" + replicaNumber,
                    "workflow_replica": request.GetWorkflow() + "-" + replicaNumber,
                    "version":          strconv.Itoa(workflow_latest_version[request.GetWorkflow()]),
                }
                replicaCount := int32(1)
                deployment := &appsv1.Deployment{
                    ObjectMeta: metav1.ObjectMeta{
                        Name:   pod_name + "-" + replicaNumber,
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
                                NodeName: node_sel,
                                Containers: make([]corev1.Container, len(pod)),
                            },
                        },
                    },
                }
                for i, container := range pod {
                    container_name := strings.ToLower(strings.ReplaceAll(container, "_", "-")) + "-" + replicaNumber
                    service_port := 5000 + slices.Index(workflow_initial_pods[request.GetWorkflow()], container)
                    http_port := 6000 + slices.Index(workflow_initial_pods[request.GetWorkflow()], container)
                    function := workflows[request.GetWorkflow()].Functions[container]
                    registry := function.GetRegistry()
                    var env []corev1.EnvVar
                    for name, value := range function.Envs {
                        env = append(env, corev1.EnvVar{Name: name, Value: value})
                    }
                    service_port_s := strconv.Itoa(service_port)
                    env = append(env, corev1.EnvVar{Name: "SERVICE_PORT", Value: service_port_s})
                    http_port_s := strconv.Itoa(http_port)
                    env = append(env, corev1.EnvVar{Name: "HTTP_PORT", Value: http_port_s})
                    env = append(env, corev1.EnvVar{Name: "WORKFLOW", Value: request.GetWorkflow()})
                    env = append(env, corev1.EnvVar{Name: "FUNCTION", Value: container})
                    env = append(env, corev1.EnvVar{Name: "INGRESS", Value: ingress_address})
                    env = append(env, corev1.EnvVar{Name: "LOGGER", Value: logger_address})
                    env = append(env, corev1.EnvVar{Name: "COMM_TYPE", Value: communication})
                    for _, endpoint := range function.Endpoints {
                        in_pod := false
                        for _, c := range pod {
                            if endpoint == c {
                                in_pod = true
                                break
                            }
                        }
                        endpoint_upper := strings.ToUpper(endpoint)
                        endpoint_name := strings.ReplaceAll(endpoint_upper, "-", "_")
                        endpoint_port := strconv.Itoa(5000 + slices.Index(workflow_initial_pods[request.GetWorkflow()], endpoint))
                        var service_name string
                        if in_pod {
                            service_name = "127.0.0.1:" + endpoint_port // structuring because we are fixating on the port number
                        } else {
                            endpoint_service := ""
                            for _, pods_list := range workflow_pods[request.GetWorkflow()] {
                                if slices.Contains(pods_list, endpoint_name) {
                                    endpoint_service = pods_list[0]
                                }
                            }
                            service_name = strings.ReplaceAll(strings.ToLower(endpoint_service), "_", "-") + "-" + replicaNumber + "." + namespace + ".svc.cluster.local:" + endpoint_port
                        }
                        env = append(env, corev1.EnvVar{Name: endpoint_name, Value: service_name})
                    }
                    container_port := int32(5000 + slices.Index(workflow_initial_pods[request.GetWorkflow()], container))
                    container_port_http := int32(6000 + slices.Index(workflow_initial_pods[request.GetWorkflow()], container))
                    imagePullPolicy := corev1.PullPolicy("IfNotPresent")
                    deployment.Spec.Template.Spec.Containers[i] = corev1.Container{
                        Name:            container_name,
                        Image:           registry,
                        ImagePullPolicy: imagePullPolicy,
                        Ports: []corev1.ContainerPort{
                            {
                                ContainerPort: container_port,
                            },
                            {
                                ContainerPort: container_port_http,
                            },
                        },
                        Env: env,
                    }
                    service.Spec.Ports = append(service.Spec.Ports, corev1.ServicePort{
                        Name:       container_name,
                        Port:       container_port,
                        TargetPort: intstr.FromInt(int(container_port)),
                    })
                    service.Spec.Ports = append(service.Spec.Ports, corev1.ServicePort{
                        Name:       container_name + "-http",
                        Port:       container_port_http,
                        TargetPort: intstr.FromInt(int(container_port_http)),
                    })
                }
                // check if deployment with name exists, if does not make a new one else update the existing lets start with creating new ones and then update the existing ones
                Debug("Looking for " + pod[0] + "_" + replicaNumber, 0)
                _, exists := kclient.AppsV1().Deployments(namespace).Get(context.Background(), strings.ToLower(pod[0])+"-"+replicaNumber, metav1.GetOptions{})
                if exists != nil {
                    Debug("Creating a new deployment " + deployment.Name, 0)
                    _, err := kclient.AppsV1().Deployments(namespace).Create(context.Background(), deployment, metav1.CreateOptions{})
                    if err != nil {
                        Debug("unable to create deployment " + strings.ToLower(pod[0]) + " for " + namespace + " - " + err.Error(), 0)
                        dataLock.Unlock()
                        return "0"
                    }
                } else {
                    Debug("Updating the existing deployment " + deployment.Name, 0)
                    update_deployments = append(update_deployments, *deployment)
                }
            }
            for _, dp := range update_deployments {
                Debug("deploying existing deployment " + dp.Name, 0)
                kclient.AppsV1().Deployments(namespace).Update(context.Background(), &dp, metav1.UpdateOptions{})
            }
            for _, pod := range workflow_pods[request.GetWorkflow()] {
                pod_name := strings.ToLower(strings.ReplaceAll(pod[0], "_", "-"))
                labels := map[string]string{
                    "workflow_name":    request.GetWorkflow(),
                    "app":              pod_name + "-" + replicaNumber,
                    "workflow_replica": request.GetWorkflow() + "-" + replicaNumber,
                    "version":          strconv.Itoa(workflow_latest_version[request.GetWorkflow()]),
                }
                service := &corev1.Service{
                    ObjectMeta: metav1.ObjectMeta{
                        Name:      pod_name + "-" + replicaNumber,
                        Namespace: namespace,
                        Labels:    labels,
                    },
                    Spec: corev1.ServiceSpec{
                        Selector: map[string]string{
                            "app": pod_name + "-" + replicaNumber,
                        },
                        Ports: []corev1.ServicePort{},
                    },
                }
                for _, container := range pod {
                    container_name := strings.ToLower(strings.ReplaceAll(container, "_", "-")) + "-" + replicaNumber
                    container_port := int32(5000 + slices.Index(workflow_initial_pods[request.GetWorkflow()], container))
                    container_port_http := int32(6000 + slices.Index(workflow_initial_pods[request.GetWorkflow()], container))
                    service.Spec.Ports = append(service.Spec.Ports, corev1.ServicePort{
                        Name:       container_name,
                        Port:       container_port,
                        TargetPort: intstr.FromInt(int(container_port)),
                    })
                    service.Spec.Ports = append(service.Spec.Ports, corev1.ServicePort{
                        Name:       container_name + "http",
                        Port:       container_port_http,
                        TargetPort: intstr.FromInt(int(container_port_http)),
                    })
                }
                _, exists := kclient.CoreV1().Services(namespace).Get(context.Background(), pod_name+"-"+replicaNumber, metav1.GetOptions{})
                if exists != nil {
                    _, err := kclient.CoreV1().Services(namespace).Create(context.Background(), service, metav1.CreateOptions{})
                    if err != nil {
                        Debug("unable to create service " + strings.ToLower(pod[0]) + "-" + replicaNumber + " for " + namespace + " - " + err.Error(), 0)
                        dataLock.Unlock()
                        return "0"
                    }
                } else {
                    _, err := kclient.CoreV1().Services(namespace).Update(context.Background(), service, metav1.UpdateOptions{})
                    if err != nil {
                        Debug("unable to update service " + strings.ToLower(pod[0]) + "-" + replicaNumber + " for " + namespace + " - " + err.Error(), 0)
                        dataLock.Unlock()
                        return "0"
                    }
                }
            }
            ingress := &networkingv1.Ingress{
                ObjectMeta: metav1.ObjectMeta{
                    Name:      service_name_ingress,
                    Namespace: namespace,
                    Labels:    labels_ingress,
                },
                Spec: networkingv1.IngressSpec{
                    Rules: []networkingv1.IngressRule{
                        {
                            Host: request.GetWorkflow() + "." + replicaNumber + "." + namespace + ".macropod",
                            IngressRuleValue: networkingv1.IngressRuleValue{
                                HTTP: &networkingv1.HTTPIngressRuleValue{
                                    Paths: []networkingv1.HTTPIngressPath{
                                        {
                                            Path:     "/",
                                            PathType: &pathType,
                                            Backend: networkingv1.IngressBackend{
                                                Service: &networkingv1.IngressServiceBackend{
                                                    Name: service_name_ingress,
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

            _, err = kclient.NetworkingV1().Ingresses(namespace).Create(context.Background(), ingress, metav1.CreateOptions{})
            if err != nil {
                Debug("unable to create ingress for " + namespace + " - " + err.Error(), 0)
                dataLock.Unlock()
                return "0"
            }
            dataLock.Unlock()
            return "0"
        default:
            return "0"
    }
}

func Serve_TTLDelete(request *pb.MacroPodRequest) (string) {
    labels := request.GetTarget()
    Debug("delete TTL " + labels, 2)
    dataLock.Lock()
    delete(workflow_deployments[request.GetWorkflow()], labels)
    namespace := default_config.GetNamespace()
    if workflows[request.GetWorkflow()].GetConfig() != nil && workflows[request.GetWorkflow()].GetConfig().GetNamespace() != "" {
        namespace = workflows[request.GetWorkflow()].GetConfig().GetNamespace()
    }
    dataLock.Unlock()
    labels_replica := "workflow_replica=" + labels
    services, err := kclient.CoreV1().Services(namespace).List(context.Background(), metav1.ListOptions{LabelSelector: labels_replica})
    if err != nil {
        Debug(err.Error(), 0)
    }
    for _, service := range services.Items {
        kclient.CoreV1().Services(namespace).Delete(context.Background(), service.ObjectMeta.Name, metav1.DeleteOptions{})
    }
    deployments, err := kclient.AppsV1().Deployments(namespace).List(context.Background(), metav1.ListOptions{LabelSelector: labels_replica})
    if err != nil {
        Debug(err.Error(), 0)
    }
    for _, deployment := range deployments.Items {
        kclient.AppsV1().Deployments(namespace).Delete(context.Background(), deployment.ObjectMeta.Name, metav1.DeleteOptions{})
    }
    ingresses, err := kclient.NetworkingV1().Ingresses(namespace).List(context.Background(), metav1.ListOptions{LabelSelector: labels_replica})
    if err != nil {
        Debug(err.Error(), 0)
    }
    for _, ingress := range ingresses.Items {
        kclient.NetworkingV1().Ingresses(namespace).Delete(context.Background(), ingress.ObjectMeta.Name, metav1.DeleteOptions{})
    }

    for {
        deployments_list, _ := kclient.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{LabelSelector: labels_replica})
        if deployments_list == nil || len(deployments_list.Items) == 0 {
            break
        }
        time.Sleep(10 * time.Millisecond)
    }
    return ""
}

func (s *DeployerService) Config(ctx context.Context, req *pb.ConfigStruct) (*pb.MacroPodReply, error) {
    reply := Serve_Config(req)
    results := pb.MacroPodReply{Reply: &reply}
    return &results, nil
}

func (s *DeployerService) CreateWorkflow(ctx context.Context, req *pb.WorkflowStruct) (*pb.MacroPodReply, error) {
    reply := Serve_CreateWorkflow(req)
    results := pb.MacroPodReply{Reply: &reply}
    return &results, nil
}

func (s *DeployerService) UpdateWorkflow(ctx context.Context, req *pb.WorkflowStruct) (*pb.MacroPodReply, error) {
    reply := Serve_UpdateWorkflow(req)
    results := pb.MacroPodReply{Reply: &reply}
    return &results, nil
}

func (s *DeployerService) DeleteWorkflow(ctx context.Context, req *pb.MacroPodRequest) (*pb.MacroPodReply, error) {
    reply := Serve_DeleteWorkflow(req)
    results := pb.MacroPodReply{Reply: &reply}
    return &results, nil
}

func (s *DeployerService) UpdateDeployments(ctx context.Context, req *pb.MacroPodRequest) (*pb.MacroPodReply, error) {
    reply := Serve_UpdateDeployments(req)
    results := pb.MacroPodReply{Reply: &reply}
    return &results, nil
}

func (s *DeployerService) CreateDeployment(ctx context.Context, req *pb.MacroPodRequest) (*pb.MacroPodReply, error) {
    reply := Serve_CreateDeployment(req, false)
    results := pb.MacroPodReply{Reply: &reply}
    return &results, nil
}

func (s *DeployerService) TTLDelete(ctx context.Context, req *pb.MacroPodRequest) (*pb.MacroPodReply, error) {
    reply := Serve_TTLDelete(req)
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

func HTTP_CreateWorkflow(res http.ResponseWriter, req *http.Request) {
    body, _ := ioutil.ReadAll(req.Body)
    request := pb.WorkflowStruct{}
    json.Unmarshal(body, &request)
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

func HTTP_UpdateDeployments(res http.ResponseWriter, req *http.Request) {
    workflow := req.PathValue("workflow")
    request := pb.MacroPodRequest{Workflow: &workflow}
    results := Serve_UpdateDeployments(&request)
    fmt.Fprint(res, results)
}

func HTTP_CreateDeployment(res http.ResponseWriter, req *http.Request) {
    workflow := req.PathValue("workflow")
    request := pb.MacroPodRequest{Workflow: &workflow}
    results := Serve_CreateDeployment(&request, false)
    fmt.Fprint(res, results)
}

func HTTP_TTLDelete(res http.ResponseWriter, req *http.Request) {
    body, _ := ioutil.ReadAll(req.Body)
    request := pb.MacroPodRequest{}
    json.Unmarshal(body, &request)
    results := Serve_TTLDelete(&request)
    fmt.Fprint(res, results)
}

func main() {
    config, err := rest.InClusterConfig()
    if err != nil {
        fmt.Println("error config - " + err.Error())
        return
    }
    kclient, err = kubernetes.NewForConfig(config)
    if err != nil {
        fmt.Println("error kclient - " + err.Error())
        return
    }
    service_port := os.Getenv("SERVICE_PORT")
    if service_port == "" {
        service_port = "8000"
    }
    http_port := os.Getenv("HTTP_PORT")
    if http_port == "" {
        http_port = "9000"
    }
    ingress_address = os.Getenv("INGRESS_ADDRESS")
    if ingress_address == "" {
        ingress_address = "127.0.0.1:8001"
    }
    logger_address = os.Getenv("LOGGER_ADDRESS")
    if logger_address == "" {
        logger_address = "127.0.0.1:8003"
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
    l, _ := net.Listen("tcp", ":" + service_port)
    s := grpc.NewServer(grpc.MaxSendMsgSize(1024*1024*200), grpc.MaxRecvMsgSize(1024*1024*200))
    pb.RegisterMacroPodDeployerServer(s, &DeployerService{})

    go CheckNodeStatus()
    go s.Serve(l)

    h := http.NewServeMux()
    h.HandleFunc("/", HTTP_Help)
    h.HandleFunc("/config", HTTP_Config)
    h.HandleFunc("/workflow/create", HTTP_CreateWorkflow)
    h.HandleFunc("/workflow/update", HTTP_UpdateWorkflow)
    h.HandleFunc("/workflow/delete/{workflow}", HTTP_DeleteWorkflow)
    h.HandleFunc("/deployment/update/{workflow}", HTTP_UpdateDeployments)
    h.HandleFunc("/deployment/create/{workflow}", HTTP_CreateDeployment)
    h.HandleFunc("/ttl", HTTP_TTLDelete)
    http.ListenAndServe(":" + http_port, h)
}
