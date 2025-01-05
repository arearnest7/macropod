package main

import (
	pb "app/deployer_pb"
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
	"math"
	"net"
	"os"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"
)

type server struct {
	pb.DeploymentServiceServer
}

type Function struct {
	Registry  string            `json:"registry"`
	Endpoints []string          `json:"endpoints,omitempty"`
	Envs      map[string]string `json:"envs,omitempty"`
	Secrets   map[string]string `json:"secrets,omitempty"`
}

type Workflow struct {
	Name               string              `json:"name,omitempty"`
	Functions          map[string]Function `json:"functions"`
	Pods               [][]string
	Deployments        map[string]map[string]string
	IngressVersion     map[string]int
	LatestVersion      int
	LastUpdated        time.Time
	ReplicaCount       int
	Updating           bool
	InitialPods        []string
	FullyDisaggregated bool
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
	kclient            *kubernetes.Clientset
	nodeCapacityCPU    = make(map[string]float64)
	nodeCapacityMemory = make(map[string]float64)

	workflows          = make(map[string]*Workflow)

	update_threshold   int
	debug              int
	ttl_seconds        int
	macropod_namespace string

	depLock            sync.Mutex
)

func deleteTTL(func_name string, labels string) {
	if debug > 2 {
		fmt.Println("delete TTL")
	}
	depLock.Lock()
	delete(workflows[func_name].Deployments, labels)
	delete(workflows[func_name].IngressVersion, labels)
	depLock.Unlock()
	labels_replica := "workflow_replica=" + labels
	services, err := kclient.CoreV1().Services(macropod_namespace).List(context.Background(), metav1.ListOptions{LabelSelector: labels_replica})
	if err != nil && debug > 0 {
		fmt.Println(err)
	}
	for _, service := range services.Items {
		kclient.CoreV1().Services(macropod_namespace).Delete(context.Background(), service.ObjectMeta.Name, metav1.DeleteOptions{})
	}
	deployments, err := kclient.AppsV1().Deployments(macropod_namespace).List(context.Background(), metav1.ListOptions{LabelSelector: labels_replica})
	if err != nil && debug > 0 {
		fmt.Println(err)
	}
	for _, deployment := range deployments.Items {
		kclient.AppsV1().Deployments(macropod_namespace).Delete(context.Background(), deployment.ObjectMeta.Name, metav1.DeleteOptions{})
	}
	ingresses, err := kclient.NetworkingV1().Ingresses(macropod_namespace).List(context.Background(), metav1.ListOptions{LabelSelector: labels_replica})
	if err != nil && debug > 0 {
		fmt.Println(err)
	}
	for _, ingress := range ingresses.Items {
		kclient.NetworkingV1().Ingresses(macropod_namespace).Delete(context.Background(), ingress.ObjectMeta.Name, metav1.DeleteOptions{})
	}

	for {
		deployments_list, _ := kclient.CoreV1().Pods(macropod_namespace).List(context.Background(), metav1.ListOptions{LabelSelector: labels_replica})
		if deployments_list == nil || len(deployments_list.Items) == 0 {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func createDeployment(func_name string, bypass bool) string {
	if !bypass {
		depLock.Lock()
		if _, exists := workflows[func_name]; !exists {
			log.Printf(" %s is not present", func_name)
			return "0"
		}
		if workflows[func_name].Updating {
			return "1"
		}
		depLock.Unlock()
	}
	depLock.Lock()
	var nodes NodeMetricList
	data, err := kclient.RESTClient().Get().AbsPath("apis/metrics.k8s.io/v1beta1/nodes").Do(context.Background()).Raw()
	if err != nil {
		return "0"
	}
	err = json.Unmarshal(data, &nodes)
	if err != nil {
		return "0"
	}
	start_node := -1
	for i, node := range nodes.Items {
		if start_node != -1 {
			break
		}
		value, exists := node.Metadata.Labels["node-role.kubernetes.io/master"]
		if exists && value == "true" {
			continue
		}
		in_use := false
		for _, deployment := range workflows[func_name].Deployments {
			if in_use && deployment[strings.ToLower(strings.ReplaceAll(workflows[func_name].Pods[0][0], "_", "-"))] == node.Metadata.Name {
				in_use = true
			}
		}
		if !in_use {
			start_node = i
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
	return "0"
}

func nodeReclaim(func_name string) {
	depLock.Lock()
	reclaim_replicas := workflows[func_name].Deployments
	depLock.Unlock()
	for replica_name, _ := range reclaim_replicas {
		depLock.Lock()
		delete(workflows[func_name].Deployments, replica_name)
		delete(workflows[func_name].IngressVersion, replica_name)
		depLock.Unlock()
		labels_replica := "workflow_replica=" + replica_name
		services, err := kclient.CoreV1().Services(macropod_namespace).List(context.Background(), metav1.ListOptions{LabelSelector: labels_replica})
		if err != nil && debug > 0 {
			fmt.Println(err)
		}
		for _, service := range services.Items {
			kclient.CoreV1().Services(macropod_namespace).Delete(context.Background(), service.ObjectMeta.Name, metav1.DeleteOptions{})
		}
		deployments, err := kclient.AppsV1().Deployments(macropod_namespace).List(context.Background(), metav1.ListOptions{LabelSelector: labels_replica})
		if err != nil && debug > 0 {
			fmt.Println(err)
		}
		for _, deployment := range deployments.Items {
			kclient.AppsV1().Deployments(macropod_namespace).Delete(context.Background(), deployment.ObjectMeta.Name, metav1.DeleteOptions{})
		}
		ingresses, err := kclient.NetworkingV1().Ingresses(macropod_namespace).List(context.Background(), metav1.ListOptions{LabelSelector: labels_replica})
		if err != nil && debug > 0 {
			fmt.Println(err)
		}
		for _, ingress := range ingresses.Items {
			kclient.NetworkingV1().Ingresses(macropod_namespace).Delete(context.Background(), ingress.ObjectMeta.Name, metav1.DeleteOptions{})
		}

		for {
			deployments_list, _ := kclient.CoreV1().Pods(macropod_namespace).List(context.Background(), metav1.ListOptions{LabelSelector: labels_replica})
			if deployments_list == nil || len(deployments_list.Items) == 0 {
				break
			}
			time.Sleep(10 * time.Millisecond)
		}

		createDeployment(func_name, true)
	}
	workflows[func_name].Updating = false
}

func updateDeployments(func_name string) string {
	depLock.Lock()
	if _, exists := workflows[func_name]; !exists {
		log.Printf(" %s is not present", func_name)
		return "0"
	}
	if workflows[func_name].Updating {
		return "1"
	}
	if time.Since(workflows[func_name].LastUpdated) < time.Second*time.Duration(update_threshold) {
		return "0"
	}
	if workflows[func_name].FullyDisaggregated {
		return "0"
	}
	workflows[func_name].Updating = true
	depLock.Unlock()

	// get the percentage of the node utilisation
	var nodes NodeMetricList
	data, err := kclient.RESTClient().Get().AbsPath("apis/metrics.k8s.io/v1beta1/nodes").Do(context.Background()).Raw()
	if err != nil {
		return "0"
	}
	err = json.Unmarshal(data, &nodes)
	if err != nil {
		return "0"
	}
	static := os.Getenv("STATIC")
	if static != "" {
		return "0"
	}
	for _, node := range nodes.Items {
		value, exists := node.Metadata.Labels["node-role.kubernetes.io/master"]
		if exists && value == "true" {
			continue
		}
		depLock.Lock()
		node_name := node.Metadata.Name
		usage_mem, _ := memory_raw_to_float(node.Usage.Memory)
		usage_cpu, _ := cpu_raw_to_float(node.Usage.CPU)
		percentage_cpu := (usage_cpu / nodeCapacityCPU[node_name]) * 100
		percentage_mem := (usage_mem / nodeCapacityMemory[node_name]) * 100
		depLock.Unlock()
		if percentage_mem > 70 || percentage_cpu > 70 {
			depLock.Lock()
			workflows[func_name].LatestVersion += 1
			var pods_updated [][]string
			for _, pod := range workflows[func_name].Pods {
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
				workflows[func_name].FullyDisaggregated = true
			}
			workflows[func_name].Pods = pods_updated
			workflows[func_name].LastUpdated = time.Now()
			go nodeReclaim(func_name)
			depLock.Unlock()
			return "2"
		}
	}

	depLock.Lock()
	if _, exists := workflows[func_name]; exists {
		workflows[func_name].Updating = false
	} else {
		log.Printf(" %s is not present", func_name)
	}
	depLock.Unlock()
	return "0"
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

func bfs_initial_pod(pod []string, func_name string, pod_list []string) []string {
	if len(pod_list) == 0 {
		return pod
	}
	entrypoint := pod_list[0]
	if !slices.Contains(pod, entrypoint) {
		pod = append(pod, entrypoint)
	}
	pod_list = pod_list[1:]
	for _, endpoint := range workflows[func_name].Functions[entrypoint].Endpoints {
		if !slices.Contains(pod, endpoint) {
			pod = append(pod, endpoint)
			pod_list = append(pod_list, endpoint)

		}
	}
	return bfs_initial_pod(pod, func_name, pod_list)
}

func createInitialPod(func_name string) {
	var initial_pod []string
	var frontend_func string
	var endpoints []string
	func_endpoint := make(map[string][]string)
	for func_name, function := range workflows[func_name].Functions {
		for _, endpoint := range function.Endpoints {
			func_endpoint[func_name] = append(func_endpoint[func_name], endpoint)
			if !slices.Contains(endpoints, endpoint) {
				if func_name != endpoint {
					endpoints = append(endpoints, endpoint)
				}
			}
		}
	}
	for func_name := range workflows[func_name].Functions {
		if !slices.Contains(endpoints, func_name) {
			frontend_func = func_name
			break
		}

	}
	var pod_list []string
	pod_list = append(pod_list, frontend_func)
	initial_pod = bfs_initial_pod(initial_pod, func_name, pod_list)
	static := os.Getenv("STATIC")
	log.Printf("printing the pods: %v", initial_pod)
	if static == "full-agg" {
		workflows[func_name].Pods = append(workflows[func_name].Pods, initial_pod)
	} else if static == "partial" {
		pods_split := make([][]string, 0)
		pods_split = append(pods_split, initial_pod[0:(len(initial_pod)/2)])
		pods_split = append(pods_split, initial_pod[(len(initial_pod)/2):])
		workflows[func_name].Pods = pods_split
	} else if static == "full-disagg" {
		pods_disagg := make([][]string, 0)
		for _, container := range initial_pod {
			container_pod := make([]string, 0)
			container_pod = append(container_pod, container)
			pods_disagg = append(pods_disagg, container_pod)
		}
		workflows[func_name].Pods = pods_disagg
	} else {
		workflows[func_name].Pods = append(workflows[func_name].Pods, initial_pod)
	}
	workflows[func_name].InitialPods = initial_pod
	workflows[func_name].LatestVersion = 1

}

func createWorkflow(func_name string, func_str string) {
	depLock.Lock()
	workflow := Workflow{}
	json.Unmarshal([]byte(func_str), &workflow)
	_, exists := workflows[func_name]
	if exists {
		fmt.Println("workflow " + func_name + " already exists. If you are updating it please use update instead.")
		return
	}
	workflows[func_name] = &workflow
	workflows[func_name].Deployments = make(map[string]map[string]string)
	workflows[func_name].IngressVersion = make(map[string]int)
	createInitialPod(func_name)
	depLock.Unlock()
}

func updateWorkflow(func_name string, workflow_str string) {
	depLock.Lock()
	workflow := Workflow{}
	json.Unmarshal([]byte(workflow_str), &workflow)
	_, exists := workflows[func_name]
	if exists {
		delete(workflows, func_name)
		label_workflow := "workflow_name=" + func_name
		services, err := kclient.CoreV1().Services(macropod_namespace).List(context.Background(), metav1.ListOptions{LabelSelector: label_workflow})
		if err != nil && debug > 0 {
			fmt.Println(err)
		}
		for _, service := range services.Items {
			kclient.CoreV1().Services(macropod_namespace).Delete(context.Background(), service.ObjectMeta.Name, metav1.DeleteOptions{})
		}
		deployments, err := kclient.AppsV1().Deployments(macropod_namespace).List(context.Background(), metav1.ListOptions{LabelSelector: label_workflow})
		if err != nil && debug > 0 {
			fmt.Println(err)
		}
		for _, deployment := range deployments.Items {
			kclient.AppsV1().Deployments(macropod_namespace).Delete(context.Background(), deployment.ObjectMeta.Name, metav1.DeleteOptions{})
		}
		ingresses, err := kclient.NetworkingV1().Ingresses(macropod_namespace).List(context.Background(), metav1.ListOptions{LabelSelector: label_workflow})
		if err != nil && debug > 0 {
			fmt.Println(err)
		}
		for _, ingress := range ingresses.Items {
			kclient.NetworkingV1().Ingresses(macropod_namespace).Delete(context.Background(), ingress.ObjectMeta.Name, metav1.DeleteOptions{})
		}
		for {
			deployments_list, _ := kclient.CoreV1().Pods(macropod_namespace).List(context.Background(), metav1.ListOptions{LabelSelector: label_workflow})
			if deployments_list == nil || len(deployments_list.Items) == 0 {
				break
			}
			time.Sleep(10 * time.Millisecond)
		}
	}
	workflows[func_name] = &workflow
	workflows[func_name].Deployments = make(map[string]map[string]string)
	workflows[func_name].IngressVersion = make(map[string]int)
	createInitialPod(func_name)
	depLock.Unlock()
}

func deleteWorkflow(func_name string) {
	depLock.Lock()
	_, exists := workflows[func_name]
	if exists {
		delete(workflows, func_name)
		label_workflow := "workflow_name=" + func_name
		services, err := kclient.CoreV1().Services(macropod_namespace).List(context.Background(), metav1.ListOptions{LabelSelector: label_workflow})
		if err != nil && debug > 0 {
			fmt.Println(err)
		}
		for _, service := range services.Items {
			kclient.CoreV1().Services(macropod_namespace).Delete(context.Background(), service.ObjectMeta.Name, metav1.DeleteOptions{})
		}
		deployments, err := kclient.AppsV1().Deployments(macropod_namespace).List(context.Background(), metav1.ListOptions{LabelSelector: label_workflow})
		if err != nil && debug > 0 {
			fmt.Println(err)
		}
		for _, deployment := range deployments.Items {
			kclient.AppsV1().Deployments(macropod_namespace).Delete(context.Background(), deployment.ObjectMeta.Name, metav1.DeleteOptions{})
		}
		ingresses, err := kclient.NetworkingV1().Ingresses(macropod_namespace).List(context.Background(), metav1.ListOptions{LabelSelector: label_workflow})
		if err != nil && debug > 0 {
			fmt.Println(err)
		}
		for _, ingress := range ingresses.Items {
			kclient.NetworkingV1().Ingresses(macropod_namespace).Delete(context.Background(), ingress.ObjectMeta.Name, metav1.DeleteOptions{})
		}
		for {
			deployments_list, _ := kclient.CoreV1().Pods(macropod_namespace).List(context.Background(), metav1.ListOptions{LabelSelector: label_workflow})
			if deployments_list == nil || len(deployments_list.Items) == 0 {
				break
			}
			time.Sleep(10 * time.Millisecond)
		}
	}
	depLock.Unlock()
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
					delete(nodeCapacityCPU, node.Name)
					delete(nodeCapacityMemory, node.Name)
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
			nodeCapacityCPU[node.Name] = cpu_float
			nodeCapacityMemory[node.Name] = mem_float
			depLock.Unlock()
		}
		time.Sleep(30 * time.Second)
	}

}

func (s *server) Deployment(ctx context.Context, req *pb.DeploymentServiceRequest) (*pb.DeploymentServiceReply, error) {
	func_name := req.Name
	request_type := req.FunctionCall
	result := "0"
	switch request_type {
		case "create_workflow":
			createWorkflow(func_name, *req.Workflow)
		case "update_workflow":
			updateWorkflow(func_name, *req.Workflow)
		case "delete_workflow":
			deleteWorkflow(func_name)
		case "update_deployments":
			result = updateDeployments(func_name)
		case "create_deployment":
			result = createDeployment(func_name, false)
		case "ttl_delete":
			deleteTTL(func_name, *req.Workflow)
	}
	return &pb.DeploymentServiceReply{
		Message: result,
	}, nil
}

func main() {
	var err error
	service_port, err := strconv.Atoi(os.Getenv("SERVICE_PORT"))
	if err != nil {
		service_port = 8082
	}
	macropod_namespace = os.Getenv("MACROPOD_NAMESPACE")
	if macropod_namespace == "" {
		macropod_namespace = "macropod-functions"
	}
	debug, err = strconv.Atoi(os.Getenv("DEBUG"))
	if err != nil {
		debug = 0
	}
	update_threshold, err = strconv.Atoi(os.Getenv("UPDATE_THRESHOLD"))
	if err != nil {
		update_threshold = 100
	}
	ttl_seconds, err = strconv.Atoi(os.Getenv("TTL"))
	if err != nil {
		ttl_seconds = 180
	}
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
	go checkNodeStatus()
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", service_port))
	if err != nil {
		fmt.Println("error listener - " + err.Error())
		return
	}
	s := grpc.NewServer()
	pb.RegisterDeploymentServiceServer(s, &server{})
	if err := s.Serve(l); err != nil {
		fmt.Println("failed to serve - " + err.Error())
		return
	}
}
