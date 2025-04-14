package main

import (
    pb "app/macropod_pb"

    "os"
    "encoding/json"
    "fmt"
    "strconv"
    "context"
    "log"
    "math"
    "slices"
    "strings"
    "sync"
    "time"

    "github.com/soheilhy/cmux"
    "net"
    "net/http"

    "google.golang.org/grpc"
    "golang.org/x/net/context"

    appsv1 "k8s.io/api/apps/v1"
    corev1 "k8s.io/api/core/v1"
    networkingv1 "k8s.io/api/networking/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/util/intstr"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/rest"
)

type DeployerService struct {
    pb.UnimplementedMacroPodDeployerServer
}

var (
    kclient        *kubernetes.Clientset

    manifest       = make(pb.MacroPodManifest)

    depLock        sync.Mutex
)

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

func bfs_initial_pod(pod []string, wf_name string, pod_list []string) []string {
    if len(pod_list) == 0 {
        return pod
    }
    entrypoint := pod_list[0]
    if !slices.Contains(pod, entrypoint) {
        pod = append(pod, entrypoint)
    }
    pod_list = pod_list[1:]
    for _, endpoint := range manifest.Workflows[wf_name].Workflow.Functions[entrypoint].Endpoints {
        if !slices.Contains(pod, endpoint) {
            pod = append(pod, endpoint)
            pod_list = append(pod_list, endpoint)

        }
    }
    return bfs_initial_pod(pod, wf_name, pod_list)
}

func createInitialPod(wf_name string) {
    var initial_pod []string
    var frontend_func string
    var endpoints []string
    func_endpoint := make(map[string][]string)
    for func_name, function := range manifest.Workflows[wf_name].Workflow.Functions {
        for _, endpoint := range function.Endpoints {
            func_endpoint[func_name] = append(func_endpoint[func_name], endpoint)
            if !slices.Contains(endpoints, endpoint) {
                if func_name != endpoint {
                    endpoints = append(endpoints, endpoint)
                }
            }
        }
    }
    for func_name, _ := range manifest.Workflows[wf_name].Workflow.Functions {
        if !slices.Contains(endpoints, func_name) {
            frontend_func = func_name
            break
        }
    }
    var pod_list []string
    pod_list = append(pod_list, frontend_func)
    initial_pod = bfs_initial_pod(initial_pod, wf_name, pod_list)
    aggregation := manifest.DefaultConfig.Aggregation
    if manifest.Workflows[wf_name].Workflow.Config.Aggregation != "" {
        aggregation = manifest.Workflows[wf_name].Workflow.Config.Aggregation
    }
    deployment_config := pb.DeploymentMetadata{ID: -1, Entrypoint: frontend_func}
    switch aggregation {
        case "agg":
            pod := pb.PodMetadata{Name: frontend_func}
            for _, func_name := range initial_pod {
                pod.Ingresses[func_name] = ""
            }
            deployment_config.Pods[frontend_func] = pod
        case "disagg":
            for _, func_name := range initial_pod {
                deployment_config.Pods[func_name] = pb.PodMetadata{Name: func_name}
                deployment_config.Pods[func_name].Ingresses[func_name] = ""
            }
        default:
            pod := pb.PodMetadata{Name: frontend_func}
            for _, func_name := range initial_pod {
                pod.Ingresses[func_name] = ""
            }
            deployment_config.Pods[frontend_func] = pod
    }
    manifest.Workflows[wf_name].DeploymentConfig = deployment_config
    manifest.Workflows[wf_name].LatestVersion = 1
}

func checkNodeStatus() {
    for {
        nodes, err := kclient.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
        if err != nil {
            panic(err)
        }
        for _, node := range nodes.Items {

            notReady := false
            for _, condition := range node.Status.Conditions {
                if condition.Type == "Ready" && condition.Status != "True" {
                    depLock.Lock()
                    delete(manifest.NodeMetrics, node.Name)
                    depLock.Unlock()
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
            depLock.Lock()
            manifest.NodeMetrics[node.Name].CPU = cpu_float
            manifest.NodeMetrics[node.Name].Memory = mem_float
            depLock.Unlock()
        }
        time.Sleep(30 * time.Second)
    }

}

func Deployer_Config(config pb.ConfigStruct) (error) {
    depLock.Lock()
    manifest.DefaultConfig = config
    if manifest.DefaultConfig.Namespace == "" {
        manifest.DefaultConfig.Namespace = "macropod-functions"
    }
    if manifest.DefaultConfig.TTL == 0 {
        manifest.DefaultConfig.TTL = 180
    }
    if manifest.DefaultConfig.Deployment == "" {
        manifest.DefaultConfig.Deployment = "macropod"
    }
    if manifest.DefaultConfig.Communication == "" {
        manifest.DefaultConfig.Communication = "direct"
    }
    if manifest.DefaultConfig.Aggregation == "" {
        manifest.DefaultConfig.Aggregation = "agg"
    }
    if manifest.DefaultConfig.TargetConcurrency == 0 {
        manifest.DefaultConfig.TargetConcurrency = -1
    }
    depLock.Unlock()
    return nil
}

func Deployer_CreateWorkflow(workflow pb.WorkflowStruct) (error) {
    depLock.Lock()
    _, exists := manifest.Workflows[workflow.WorkflowName]
    if exists {
        depLock.Unlock()
        return
    }
    manifest.Workflows[workflow.WorkflowName] = pb.WorkflowMetadata{Workflow: workflow}
    createInitialPod(workflow.WorkflowName)
    depLock.Unlock()
    return nil
}

func Deployer_UpdateWorkflow(workflow pb.WorkflowStruct) (error) {
    depLock.Lock()
    _, exists := manifest.Workflows[workflow.WorkflowName]
    if !exists {
        depLock.Unlock()
        Deployer_CreateWorkflow(workflow)
        return nil
    }
    manifest.Workflows[workflow.WorkflowName].Workflow = workflow
    createInitialPod(workflow.WorkflowName)
    depLock.Unlock()
    return nil
}

func Deployer_DeleteWorkflow(deployer_request pb.DeployerRequest) (error) {
    depLock.Lock()
    wf_name := deployer_request.WorkflowName
    _, exists := workflows[wf_name]
    if exists {
        if debug > 2 {
            fmt.Println("deleting workflow " + wf_name)
        }
        label_workflow := "workflow_name=" + wf_name
        namespace := manifest.Workflows[wf_name].Config.Namespace
        if namespace == "" {
            namespace = manifest.DefaultConfig.Namespace
        }
        services, err := kclient.CoreV1().Services(namespace).List(context.Background(), metav1.ListOptions{LabelSelector: label_workflow})
        if err != nil && debug > 0 {
            fmt.Println(err)
        }
        for _, service := range services.Items {
            kclient.CoreV1().Services(namespace).Delete(context.Background(), service.ObjectMeta.Name, metav1.DeleteOptions{})
        }
        deployments, err := kclient.AppsV1().Deployments(namespace).List(context.Background(), metav1.ListOptions{LabelSelector: label_workflow})
        if err != nil && debug > 0 {
            fmt.Println(err)
        }
        for _, deployment := range deployments.Items {
            kclient.AppsV1().Deployments(namespace).Delete(context.Background(), deployment.ObjectMeta.Name, metav1.DeleteOptions{})
        }
        ingresses, err := kclient.NetworkingV1().Ingresses(namespace).List(context.Background(), metav1.ListOptions{LabelSelector: label_workflow})
        if err != nil && debug > 0 {
            fmt.Println(err)
        }
        for _, ingress := range ingresses.Items {
            kclient.NetworkingV1().Ingresses(namespace).Delete(context.Background(), ingress.ObjectMeta.Name, metav1.DeleteOptions{})
        }
        for {
            if debug > 4 {
                fmt.Println("waiting for " + wf_name + " to delete")
            }
            deployments_list, _ := kclient.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{LabelSelector: label_workflow})
            if deployments_list == nil || len(deployments_list.Items) == 0 {
                break
            }
            time.Sleep(10 * time.Millisecond)
        }
        delete(manifest.Workflows, wf_name)
        delete(manifest.WorkflowIngresses, wf_name)
        for func_name_ingress, func_ingress := range manifest.Workflows.FunctionIngresses {
            if func_ingress.WorkflowName == wf_name {
                delete(manifest.FunctionIngresses, func_name_ingress)
            }
        }
    }
    depLock.Unlock()
    return nil
}

func Deployer_UpdateDeployments(deployer_request pb.DeployerRequest) (string, error) {
    return "", nil //implement if macropod-dynamic is needed
}

func Deployer_CreateDeployment(deployer_request pb.DeployerRequest) (string, error) {
    wf_name := deployer_request.WorkflowName
    if !bypass {
        depLock.Lock()
        if _, exists := workflows[func_name]; !exists {
            log.Printf(" %s is not present", func_name)
            depLock.Unlock()
            return "0"
        }
        if workflows[func_name].Updating {
            depLock.Unlock()
            return "1"
        }
        depLock.Unlock()
    }
    depLock.Lock()
    var nodes NodeMetricList
    data, err := kclient.RESTClient().Get().AbsPath("apis/metrics.k8s.io/v1beta1/nodes").Do(context.Background()).Raw()
    if err != nil {
        depLock.Unlock()
        return "0"
    }
    err = json.Unmarshal(data, &nodes)
    if err != nil {
        depLock.Unlock()
        return "0"
    }
    start_node := -1
    for i, node := range nodes.Items {
        value, exists := node.Metadata.Labels["node-role.kubernetes.io/master"]
        if exists && value == "true" {
            continue
        }
        in_use := false
        for _, deployment := range workflows[func_name].Deployments {
            if deployment[strings.ToLower(strings.ReplaceAll(workflows[func_name].Pods[0][0], "_", "-"))] == node.Metadata.Name {
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
        depLock.Unlock()
        return "0"
    }
    replicaNumber := strconv.Itoa(workflows[func_name].ReplicaCount)
    workflows[func_name].ReplicaCount++
    var update_deployments []appsv1.Deployment
    labels_ingress := map[string]string{
        "workflow_name":    func_name,
        "workflow_replica": func_name + "-" + replicaNumber,
        "version":          strconv.Itoa(workflows[func_name].LatestVersion),
    }
    pathType := networkingv1.PathTypePrefix
    service_name_ingress := strings.ToLower(strings.ReplaceAll(workflows[func_name].Pods[0][0], "_", "-")) + "-" + replicaNumber
    workflows[func_name].Deployments[func_name + "-" + replicaNumber] = make(map[string]string)
    node_index := start_node
    for _, pod := range workflows[func_name].Pods {
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
        workflows[func_name].Deployments[func_name + "-" + replicaNumber][pod_name] = node_sel
        service := &corev1.Service{
            ObjectMeta: metav1.ObjectMeta{
                Name:      pod_name + "-" + replicaNumber,
                Namespace: macropod_namespace,
            },
            Spec: corev1.ServiceSpec{
                Selector: map[string]string{
                    "app": pod_name + "-" + replicaNumber,
                },
                Ports: []corev1.ServicePort{},
            },
        }

        labels := map[string]string{
            "workflow_name":    func_name,
            "app":              pod_name + "-" + replicaNumber,
            "workflow_replica": func_name + "-" + replicaNumber,
            "version":          strconv.Itoa(workflows[func_name].LatestVersion),
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
            func_port := 5000 + slices.Index(workflows[func_name].InitialPods, container)
            function := workflows[func_name].Functions[container]
            registry := function.Registry
            var env []corev1.EnvVar
            for name, value := range function.Envs {
                env = append(env, corev1.EnvVar{Name: name, Value: value})
            }
            env = append(env, corev1.EnvVar{Name: "SERVICE_TYPE", Value: "GRPC"})
            env = append(env, corev1.EnvVar{Name: "GRPC_THREAD", Value: strconv.Itoa(10)})
            func_port_s := strconv.Itoa(func_port)
            env = append(env, corev1.EnvVar{Name: "FUNC_PORT", Value: func_port_s})
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
                endpoint_port := strconv.Itoa(5000 + slices.Index(workflows[func_name].InitialPods, endpoint))
                var service_name string
                if in_pod {
                    service_name = "127.0.0.1:" + endpoint_port // structuring because we are fixating on the port number
                } else {
                    endpoint_service := ""
                    for _, pods_list := range workflows[func_name].Pods {
                        if slices.Contains(pods_list, endpoint_name) {
                            endpoint_service = pods_list[0]
                        }
                    }
                    service_name = strings.ReplaceAll(strings.ToLower(endpoint_service), "_", "-") + "-" + replicaNumber + "." + macropod_namespace + ".svc.cluster.local:" + endpoint_port
                }
                env = append(env, corev1.EnvVar{Name: endpoint_name, Value: service_name})
            }
            container_port := int32(5000 + slices.Index(workflows[func_name].InitialPods, container))
            imagePullPolicy := corev1.PullPolicy("IfNotPresent")
            deployment.Spec.Template.Spec.Containers[i] = corev1.Container{
                Name:            container_name,
                Image:           registry,
                ImagePullPolicy: imagePullPolicy,
                Ports: []corev1.ContainerPort{
                    {
                        ContainerPort: container_port,
                    },
                },
                Env: env,
            }
            service.Spec.Ports = append(service.Spec.Ports, corev1.ServicePort{
                Name:       container_name,
                Port:       container_port,
                TargetPort: intstr.FromInt(int(container_port)),
            })
        }
        // check if deployment with name exists, if does not make a new one else update the existing lets start with creating new ones and then update the existing ones
        fmt.Println("Looking for " + pod[0] + "_" + replicaNumber)
        _, exists := kclient.AppsV1().Deployments(macropod_namespace).Get(context.Background(), strings.ToLower(pod[0])+"-"+replicaNumber, metav1.GetOptions{})
        if exists != nil {
            fmt.Println("Creating a new deployment " + deployment.Name)
            _, err := kclient.AppsV1().Deployments(macropod_namespace).Create(context.Background(), deployment, metav1.CreateOptions{})
            if err != nil {
                fmt.Println("unable to create deployment " + strings.ToLower(pod[0]) + " for " + macropod_namespace + " - " + err.Error())
                depLock.Unlock()
                return "0"
            }

        } else {
            fmt.Println("Updating the existing deployment " + deployment.Name)
            update_deployments = append(update_deployments, *deployment)
        }
    }

    for _, dp := range update_deployments {
        fmt.Println("deploying existing deployment " + dp.Name)
        kclient.AppsV1().Deployments(macropod_namespace).Update(context.Background(), &dp, metav1.UpdateOptions{})
    }

    for _, pod := range workflows[func_name].Pods {
        pod_name := strings.ToLower(strings.ReplaceAll(pod[0], "_", "-"))
        labels := map[string]string{
            "workflow_name":    func_name,
            "app":              pod_name + "-" + replicaNumber,
            "workflow_replica": func_name + "-" + replicaNumber,
            "version":          strconv.Itoa(workflows[func_name].LatestVersion),
        }
        service := &corev1.Service{
            ObjectMeta: metav1.ObjectMeta{
                Name:      pod_name + "-" + replicaNumber,
                Namespace: macropod_namespace,
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
            container_port := int32(5000 + slices.Index(workflows[func_name].InitialPods, container))
            service.Spec.Ports = append(service.Spec.Ports, corev1.ServicePort{
                Name:       container_name,
                Port:       container_port,
                TargetPort: intstr.FromInt(int(container_port)),
            })
        }

        _, exists := kclient.CoreV1().Services(macropod_namespace).Get(context.Background(), pod_name+"-"+replicaNumber, metav1.GetOptions{})
        if exists != nil {
            _, err := kclient.CoreV1().Services(macropod_namespace).Create(context.Background(), service, metav1.CreateOptions{})
            if err != nil {
                fmt.Println("unable to create service " + strings.ToLower(pod[0]) + "-" + replicaNumber + " for " + macropod_namespace + " - " + err.Error())
                depLock.Unlock()
                return "0"
            }
        } else {
            _, err := kclient.CoreV1().Services(macropod_namespace).Update(context.Background(), service, metav1.UpdateOptions{})
            if err != nil {
                fmt.Println("unable to update service " + strings.ToLower(pod[0]) + "-" + replicaNumber + " for " + macropod_namespace + " - " + err.Error())
                depLock.Unlock()
                return "0"
            }
        }

    }

    ingress := &networkingv1.Ingress{
        ObjectMeta: metav1.ObjectMeta{
            Name:      service_name_ingress,
            Namespace: macropod_namespace,
            Labels:    labels_ingress,
        },
        Spec: networkingv1.IngressSpec{
            Rules: []networkingv1.IngressRule{
                {
                    Host: func_name + "." + replicaNumber + "." + macropod_namespace + ".macropod",
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

    _, err = kclient.NetworkingV1().Ingresses(macropod_namespace).Create(context.Background(), ingress, metav1.CreateOptions{})
    if err != nil {
        fmt.Println("unable to create ingress for " + macropod_namespace + " - " + err.Error())
        depLock.Unlock()
        return "0"
    }
    depLock.Unlock()
    return "", nil //TODO implement in macropod-deployer
}

func Deployer_TTLDelete(deployer_request pb.DeployerRequest) (string, error) {
    if debug > 2 {
        fmt.Println("delete TTL " + labels)
    }
    depLock.Lock()
    wf_name := deployer_request.WorkflowName
    namespace := manifest.Workflows[wf_name].Config.Namespace
    delete(manifest.Workflows[wf_name].Deployments, labels)
    for idx, label := range manifest.WorkflowIngresses[wf_name].Labels {
        if label == deployer_request.Label {
            delete(manifest.WorkflowIngresses[wf_name].Labels, idx)
            delete(manifest.WorkflowIngresses[wf_name].Ingresses, idx)
        }
    }
    for f_idx, ingress := range manifest.FunctionIngresses {
        for idx, label := range manifest.FunctionIngresses[f_idx].Labels {
            if label == deployer_request.Label {
                delete(manifest.FunctionIngresses[f_idx].Labels, idx)
                delete(manifest.FunctionIngresses[f_idx].Ingresses, idx)
            }
        }
    }
    depLock.Unlock()
    labels_replica := "workflow_replica=" + labels
    services, err := kclient.CoreV1().Services(namespace).List(context.Background(), metav1.ListOptions{LabelSelector: labels_replica})
    if err != nil && debug > 0 {
        fmt.Println(err)
    }
    for _, service := range services.Items {
        kclient.CoreV1().Services(namespace).Delete(context.Background(), service.ObjectMeta.Name, metav1.DeleteOptions{})
    }
    deployments, err := kclient.AppsV1().Deployments(namespace).List(context.Background(), metav1.ListOptions{LabelSelector: labels_replica})
    if err != nil && debug > 0 {
        fmt.Println(err)
    }
    for _, deployment := range deployments.Items {
        kclient.AppsV1().Deployments(namespace).Delete(context.Background(), deployment.ObjectMeta.Name, metav1.DeleteOptions{})
    }
    ingresses, err := kclient.NetworkingV1().Ingresses(namespace).List(context.Background(), metav1.ListOptions{LabelSelector: labels_replica})
    if err != nil && debug > 0 {
        fmt.Println(err)
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
    return "", err
}

func Deployer_UpdateManifest(manifest pb.MacroPodManifest) (error) {
    
    return nil //TODO implement in macropod-deployer
}

func (s *DeployerService) Config(ctx context.Context, in *pb.ConfigStruct) (*pb.MacroPodReply, error) {
    err := Deployer_Config(*in)
    res := pb.MacroPodReply{}
    return &res, err
}

func (s *DeployerService) CreateWorkflow(ctx context.Context, in *pb.WorkflowStruct) (*pb.MacroPodReply, error) {
    err := Deployer_CreateWorkflow(*in)
    res := pb.MacroPodReply{}
    return &res, err
}

func (s *DeployerService) UpdateWorkflow(ctx context.Context, in *pb.WorkflowStruct) (*pb.MacroPodReply, error) {
    err := Deployer_UpdateWorkflow(*in)
    res := pb.MacroPodReply{}
    return &res, err
}

func (s *DeployerService) DeleteWorkflow(ctx context.Context, in *pb.DeployerRequest) (*pb.MacroPodReply, error) {
    err := Deployer_DeleteWorkflow(*in)
    res := pb.MacroPodReply{}
    return &res, err
}

func (s *DeployerService) UpdateDeployments(ctx context.Context, in *pb.DeployerRequest) (*pb.MacroPodReply, error) {
    code, err := Deployer_UpdateDeployments(*in)
    c, _ := strconv.Atoi(code)
    res := pb.MacroPodReply{Code: int32(c)}
    return &res, err
}

func (s *DeployerService) CreateDeployment(ctx context.Context, in *pb.DeployerRequest) (*pb.MacroPodReply, error) {
    code, err := Deployer_CreateDeployment(*in)
    c, _ := strconv.Atoi(code)
    res := pb.MacroPodReply{Code: int32(c)}
    return &res, err
}

func (s *DeployerService) TTLDelete(ctx context.Context, in *pb.DeployerRequest) (*pb.MacroPodReply, error) {
    err := Deployer_TTLDelete(*in)
    res := pb.MacroPodReply{}
    return &res, err
}

func (s *DeployerService) UpdateManifest(ctx context.Context, in *pb.MacroPodManifest) (*pb.MacroPodReply, error) {
    err := Deployer_UpdateDeployerManifest(*in)
    res := &manifest
    return &res, err
}

func Serve_Deployer_Help(res http.ResponseWriter, req *http.Request) {
    help_print := "MacroPod Deployer\nplease use the following paths:\n"
    help_print := "Config:\n - path: /config\n - purpose: \n - payload: \n - output: \n"
    help_print += "Workflow Create:\n - path: /workflow/create/{wf_name}\n - purpose: \n - payload: \n - output: \n"
    help_print += "Workflow Update:\n - path: /workflow/update/{wf_name}\n - purpose: \n - payload: \n - output: \n"
    help_print += "Workflow Delete:\n - path: /workflow/delete/{wf_name}\n - purpose: \n - payload: \n - output: \n"
    help_print += "Deployment Update:\n - path: /deployment/update/{wf_name}\n - purpose: \n - payload: \n - output: \n"
    help_print += "Deployment Create:\n - path: /deployment/create/{wf_name}\n - purpose: \n - payload: \n - output: \n"
    help_print += "TTL Delete:\n - path: /workflow/ttl/{wf_name}\n - purpose: \n - payload: \n - output: \n"
    fmt.Fprint(res, help_print)
}

func Serve_Deployer_Config(res http.ResponseWriter, req *http.Request) {
    body, _ := ioutil.ReadAll(req.Body)
    config := pb.ConfigStruct{}
    json.Unmarshal(body, &config)
    defer req.Body.Close()
    err := Deployer_Config(config)
    if err != nil {
        fmt.Fprint(res, err)
        return
    }
    fmt.Fprint(res, "MacroPod config has been set")
}

func Serve_Deployer_CreateWorkflow(res http.ResponseWriter, req *http.Request) {
    body, _ := ioutil.ReadAll(req.Body)
    workflow := pb.WorkflowStruct{}
    json.Unmarshal(body, &workflow)
    defer req.Body.Close()
    err := Deployer_CreateWorkflow(workflow)
    if err != nil {
        fmt.Fprint(res, err)
        return
    }
    fmt.Fprint(res, "")
}

func Serve_Deployer_UpdateWorkflow(res http.ResponseWriter, req *http.Request) {
    body, _ := ioutil.ReadAll(req.Body)
    workflow := pb.WorkflowStruct{}
    json.Unmarshal(body, &workflow)
    defer req.Body.Close()
    err := Deployer_UpdateWorkflow(workflow)
    if err != nil {
        fmt.Fprint(res, err)
        return
    }
    fmt.Fprint(res, "")
}

func Serve_Deployer_DeleteWorkflow(res http.ResponseWriter, req *http.Request) {
    deployer_request := pb.DeployerRequest{WorkflowName: req.PathValue("wf_name")}
    err := Deployer_DeleteWorkflow(deployer_request)
    if err != nil {
        fmt.Fprint(res, err)
        return
    }
    fmt.Fprint(res, "")
}

func Serve_Deployer_UpdateDeployments(res http.ResponseWriter, req *http.Request) {
    deployer_request := pb.DeployerRequest{WorkflowName: req.PathValue("wf_name")}
    code, err := Deployer_UpdateDeployments(deployer_request)
    if err != nil {
        fmt.Fprint(res, err)
        return
    }
    fmt.Fprint(res, code)
}

func Serve_Deployer_CreateDeployment(res http.ResponseWriter, req *http.Request) {
    deployer_request := pb.DeployerRequest{WorkflowName: req.PathValue("wf_name")}
    code, err := Deployer_CreateDeployment(deployer_request)
    if err != nil {
        fmt.Fprint(res, err)
        return
    }
    fmt.Fprint(res, code)
}

func Serve_Deployer_TTLDelete(res http.ResponseWriter, req *http.Request) {
    body, _ := ioutil.ReadAll(req.Body)
    deployer_request := pb.DeployerRequest{}
    json.Unmarshal(body, &deployer_request)
    defer req.Body.Close()
    deployer_request.WorkflowName = req.PathValue("wf_name")
    err := Deployer_TTLDelete(deployer_request)
    if err != nil {
        fmt.Fprint(res, err)
        return
    }
    fmt.Fprint(res, "ttl delete")
}

func main() {
    service_port := os.Getenv("SERVICE_PORT")
    if service_port == "" {
        service_port = "8082"
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
    pb.RegisterMacroPodDeployerServer(s_g, &DeployerService{})
    m_h.HandleFunc("/", Serve_Deployer_Help)
    m_h.HandleFunc("/config", Serve_Deployer_Config)
    m_h.HandleFunc("/workflow/create/{wf_name}", Serve_Deployer_CreateWorkflow)
    m_h.HandleFunc("/workflow/update/{wf_name}", Serve_Deployer_UpdateWorkflow)
    m_h.HandleFunc("/workflow/delete/{wf_name}", Serve_Deployer_DeleteWorkflow)
    m_h.HandleFunc("/deployment/update/{wf_name}", Serve_Deployer_UpdateDeployments)
    m_h.HandleFunc("/deployment/create/{wf_name}", Serve_Deployer_CreateDeployment)
    m_h.HandleFunc("/workflow/ttl/{wf_name}", Serve_Deployer_TTLDelete)

    go checkNodeStatus()

    go s_g.Serve(l_g)
    s_h := &http.Server{Handler: m_h}
    go s_h.Serve(l_h)
}
