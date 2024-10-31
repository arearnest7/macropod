package main

import (
	pb "app/deployer_pb"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
	"google.golang.org/grpc"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type server struct {
	pb.DeploymentServiceServer
}

type PodMetricsList struct {
	Kind       string           `json:"kind"`
	APIVersion string           `json:"apiVersion"`
	Metadata   MetadataPod      `json:"metadata"`
	Items      []PodMetricsItem `json:"items"`
}

type MetadataPod struct {
	Name              string            `json:"name,omitempty"`
	Namespace         string            `json:"namespace,omitempty"`
	CreationTimestamp time.Time         `json:"creationTimestamp,omitempty"`
	Labels            map[string]string `json:"labels,omitempty"`
}

type PodMetricsItem struct {
	Metadata   MetadataPod       `json:"metadata"`
	Timestamp  time.Time         `json:"timestamp"`
	Window     string            `json:"window"`
	Containers []ContainerMetric `json:"containers"`
}

type ContainerMetric struct {
	Name  string `json:"name"`
	Usage Usage  `json:"usage"`
}

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

type ContainerMetrics struct {
	Name   string `json:"name"`
	CPU    string `json:"cpu"`
	Memory string `json:"memory"`
}

type PodMetrics struct {
	Name       string             `json:"name"`
	Namespace  string             `json:"namespace"`
	Containers []ContainerMetrics `json:"containers"`
}

type NodeMetrics struct {
	Name   string `json:"name"`
	CPU    string `json:"cpu"`
	Memory string `json:"memory"`
}

type Metrics struct {
	Pods  []PodMetrics  `json:"pods"`
	Nodes []NodeMetrics `json:"nodes"`
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
	IngressVersion     map[string]int
	LatestVersion      int
	LastUpdated        time.Time
	Updating           bool
	InitialPods        []string
	FullyDisaggregated bool
}

var (
	kclient          *kubernetes.Clientset
	workflows        map[string]*Workflow
	cpu_threshold_1  float64
	cpu_threshold_2  float64
	mem_threshold_1  float64
	mem_threshold_2  float64
	update_threshold int
	ready_deployment []string
	isSorting        bool
	node_sort        string
	nodes_list []string
	deploymentRunning bool
)

func internal_log(message string) {
	fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.000000 UTC") + ": " + message)
}

func ifPodsAreRunning(deployment_array []string) bool {
	for _, d_name := range deployment_array {
		internal_log("Checking " + d_name)
		for {
			if slices.Contains(ready_deployment, d_name) {
				break
			}
		}

	}

	return true
}

// func watchDeployments() {

// 	internal_log("WATCH_WORKFLOW_START")
// 	watcher, err := kclient.CoreV1().Pods("").Watch(context.TODO(), metav1.ListOptions{
// 		LabelSelector: "workflow_name",
// 	})
// 	internal_log("teststttt")
// 	if err != nil {
// 		internal_log("Failed to set up watcher - " + err.Error())
// 	}
// 	for event := range watcher.ResultChan() {
// 		pod, ok := event.Object.(*corev1.Pod)
// 		if !ok {
// 			internal_log("Invalid Deployment Event")
// 			continue
// 		}
// 		internal_log(pod.Labels["app-version"])
// 		internal_log(pod.Labels["app"])
// 		entry_name := pod.Labels["app-version"] + pod.Labels["app"]
// 		if pod.Status.Phase == "Running" {
// 			if !slices.Contains(ready_deployment, entry_name) {
// 				ready_deployment = append(ready_deployment, entry_name)
// 			}
// 		}else{
// 			if slices.Contains(ready_deployment, entry_name) {
// 				index := slices.Index(ready_deployment,entry_name)
// 				ready_deployment[index] =""
// 			}			
// 		}
// 		log.Print(ready_deployment)
// 	}
// 	internal_log("WATCH_WORKFLOW_END")
// }
func manageDeployment(wf_name string, replicaNumber string) (string, error) {
	//log.Print(workflows)
	deploymentRunning = true
	var update_deployments []appsv1.Deployment
	var created_deployment []string
	var updated_deployment_names []string
	namespace := "macropod-functions"
	update := true
	if workflows[wf_name].IngressVersion == nil {
		workflows[wf_name].IngressVersion = make(map[string]int)
	}
	labels_ingress := map[string]string{
		"workflow_name":     wf_name,
	}
	// log.Print(workflows[wf_name].Pods)
	pathType := networkingv1.PathTypePrefix
	service_name_ingress := workflows[wf_name].Pods[0][0] + "-" + replicaNumber
	for _, pod := range workflows[wf_name].Pods {
		service := &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      pod[0]+ "-"+ replicaNumber,
				Namespace: namespace,
			},
			Spec: corev1.ServiceSpec{
				Selector: map[string]string{
					"app": pod[0]+ "-"+ replicaNumber,
				},
				Ports: []corev1.ServicePort{},
			},
		}

		labels := map[string]string{
			"workflow_name": wf_name,
			"app":           pod[0]+ "-"+ replicaNumber,
		}
		replicaCount := int32(1)
		log.Print(pod[0])
		deployment := &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name: pod[0]+ "-"+ replicaNumber,
			},
			Spec: appsv1.DeploymentSpec{
				Replicas: &replicaCount,
				Selector: &metav1.LabelSelector{
					MatchLabels: labels,
				},
				Template: corev1.PodTemplateSpec{
					Spec: corev1.PodSpec{
						Containers: make([]corev1.Container, len(pod)),
					},
				},
			},
		}
		for i, container := range pod {
			func_port := 5000 + slices.Index(workflows[wf_name].InitialPods, container)
			function := workflows[wf_name].Functions[container]
			registry := function.Registry
			var env []corev1.EnvVar
			for name, value := range function.Envs {
				env = append(env, corev1.EnvVar{Name: name, Value: value})
			}
			env = append(env, corev1.EnvVar{Name: "SERVICE_TYPE", Value: "GRPC"})
			env = append(env, corev1.EnvVar{Name: "GRPC_THREAD", Value: "10"})
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
				endpoint_port := strconv.Itoa(5000 + slices.Index(workflows[wf_name].InitialPods, endpoint))
				var service_name string
				if in_pod {
					service_name = "127.0.0.1:" + endpoint_port // structuring because we are fixating on the port number
				} else {
					service_name = endpoint + "." + namespace + ".svc.cluster.local:" + endpoint_port
				}
				env = append(env, corev1.EnvVar{Name: endpoint_name, Value: service_name})
			}
			container_port := int32(5000 + slices.Index(workflows[wf_name].InitialPods, container))
			imagePullPolicy := corev1.PullPolicy("IfNotPresent")
			deployment.Spec.Template.Spec.Containers[i] = corev1.Container{
				Name:            container,
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
				Name:       container,
				Port:       container_port,
				TargetPort: intstr.FromInt(int(container_port)),
			})
		}
		// check if deployment with name ecists, if does not make a new one else update the existing lets start with creating new ones and then update the existing ones
		internal_log("Looking for " + pod[0]+ "-"+ replicaNumber)
		_, exists := kclient.AppsV1().Deployments(namespace).Get(context.Background(), pod[0]+ "-"+ replicaNumber, metav1.GetOptions{})
		if exists != nil {
			internal_log("Creating a new deployment " + deployment.Name)
			_, err := kclient.AppsV1().Deployments(namespace).Create(context.Background(), deployment, metav1.CreateOptions{})
			if err != nil {
				internal_log("unable to create deployment " + pod[0] + " for " + namespace + " - " + err.Error())
				deploymentRunning = false
				return "", err

			}
			//entry_name := pod.Namespace+pod.Labels["app-version"]+pod.Labels["app"]
			entry_name := deployment.Spec.Template.ObjectMeta.Labels["app"]
			created_deployment = append(created_deployment, entry_name)
		} else {
			internal_log("Updating the existing deployment " + deployment.Name)
			entry_name :=  deployment.Spec.Template.ObjectMeta.Labels["app"]
			update_deployments = append(update_deployments, *deployment)
			updated_deployment_names = append(updated_deployment_names, entry_name)
		}
	}
	if !ifPodsAreRunning(created_deployment) {
		deploymentRunning = false
		return "", nil
	}
	for _, dp := range update_deployments {
		internal_log("deploying existing deployment " + dp.Name)
		kclient.AppsV1().Deployments(namespace).Update(context.Background(), &dp, metav1.UpdateOptions{})
	}
	if !ifPodsAreRunning(updated_deployment_names) {
		deploymentRunning = false
		return "", nil
	}
	for _, pod := range workflows[wf_name].Pods {
		service := &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      pod[0],
				Namespace: namespace,
			},
			Spec: corev1.ServiceSpec{
				Selector: map[string]string{
					"app": pod[0],
				},
				Ports: []corev1.ServicePort{},
			},
		}
		for _, container := range pod {
			container_port := int32(5000 + slices.Index(workflows[wf_name].InitialPods, container))
			service.Spec.Ports = append(service.Spec.Ports, corev1.ServicePort{
				Name:       container,
				Port:       container_port,
				TargetPort: intstr.FromInt(int(container_port)),
			})
		}

		_, exists := kclient.CoreV1().Services(namespace).Get(context.Background(), pod[0], metav1.GetOptions{})
		if exists != nil {
			_, err := kclient.CoreV1().Services(namespace).Create(context.Background(), service, metav1.CreateOptions{})
			if err != nil {
				internal_log("unable to create service " + pod[0] + " for " + namespace + " - " + err.Error())
				deploymentRunning = false
				return "", err
			}
		} else {
			_, err := kclient.CoreV1().Services(namespace).Update(context.Background(), service, metav1.UpdateOptions{})
			if err != nil {
				internal_log("unable to update service " + pod[0] + " for " + namespace + " - " + err.Error())
				deploymentRunning = false
				return "", err
			}
		}

	}

	ingress := &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      namespace,
			Namespace: namespace,
			Labels:    labels_ingress,
		},
		Spec: networkingv1.IngressSpec{
			Rules: []networkingv1.IngressRule{
				{
					Host: wf_name + "." + namespace + ".macropod",
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
	if update {
		_, err := kclient.NetworkingV1().Ingresses(namespace).Update(context.Background(), ingress, metav1.UpdateOptions{})
		if err != nil {
			internal_log("unable to update ingress for " + namespace + " - " + err.Error())
			deploymentRunning = false
			return "", err
		}
	} else {
		_, err := kclient.NetworkingV1().Ingresses(namespace).Create(context.Background(), ingress, metav1.CreateOptions{})
		if err != nil {
			internal_log("unable to create ingress for " + namespace + " - " + err.Error())
			deploymentRunning = false
			return "", err
		}
	}
	deploymentRunning = false
	return service_name_ingress + "." + namespace + ".svc.cluster.local:5000", nil
}

func updateDeployments(wf_name string) {
	if workflows[wf_name].Updating {
		internal_log("Already updating..........")
		return
	}

	for ns, version := range workflows[wf_name].IngressVersion {
		internal_log("version running in " + ns + " is " + strconv.Itoa(version))
	}
	workflows[wf_name].Updating = true
	// var cpu_total float64
	// var memory_total float64

	// // this is across all clusters, Question? Why are we thinking its because of the invoked ? We should check namepsaces belonging  to wf_name only right ?
	// for _, wf := range workflows {
	// 	for namespace := range wf.IngressVersion {
	// 		cpu, memory := getNamespaceMetrics(namespace)
	// 		cpu_total += cpu
	// 		memory_total += memory
	// 	}
	// }
	// internal_log("Cpu is: " + strconv.Itoa(int(cpu_total)))
	// internal_log("Memory is : " + strconv.Itoa(int(memory_total)))
	// if cpu_total > cpu_threshold_1 || memory_total > mem_threshold_1{
	// 	internal_log("threshold 1 reached - " + wf_name)
	// 	internal_log("update threshold is " + strconv.Itoa(update_threshold) + " seconds")
	// 	if time.Since(workflows[wf_name].LastUpdated) > time.Second*time.Duration(update_threshold) && !workflows[wf_name].FullyDisaggregated {
	// 		workflows[wf_name].LatestVersion += 1
	// 		internal_log("workflow " + wf_name + " updated to version " + strconv.Itoa(workflows[wf_name].LatestVersion))
	// 		var pods_updated [][]string
	// 		for _, pod := range workflows[wf_name].Pods {
	// 			if len(pod) > 1 {
	// 				idx := int(math.Floor(float64(len(pod)) / 2))
	// 				pods_updated = append(pods_updated, pod[:idx])
	// 				pods_updated = append(pods_updated, pod[idx:])
	// 			} else {
	// 				pods_updated = append(pods_updated, pod)
	// 			}
	// 		}
	// 		workflows[wf_name].Pods = pods_updated
	// 		log.Print(workflows[wf_name].Pods)
	// 		pod_2_or_more := false
	// 		for _, pod := range pods_updated {
	// 			if len(pod) > 1 {
	// 				pod_2_or_more = true
	// 				break
	// 			}
	// 		}
	// 		if !pod_2_or_more {
	// 			internal_log(wf_name + " has been fully disaggregated")
	// 			workflows[wf_name].FullyDisaggregated = true
	// 		}
	// 		workflows[wf_name].LastUpdated = time.Now()
	// 	}
	// 	if cpu_total > cpu_threshold_2 || memory_total > mem_threshold_2 {
	// 		internal_log("threshold 2 reached - " + wf_name)
	// 		for namespace, version := range workflows[wf_name].IngressVersion {
	// 			if version < workflows[wf_name].LatestVersion {
	// 				log.Printf("curret version is")
	// 				log.Print(version)
	// 				log.Print(workflows[wf_name].LatestVersion)
	// 				workflows[wf_name].IngressVersion[namespace] = workflows[wf_name].LatestVersion
	// 				go manageDeployment(wf_name, namespace)
	// 			}
	// 		}
	// 	}
	// }
	workflows[wf_name].Updating = false
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
	return math.Inf(1), fmt.Errorf("unsupported memory usage format")
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
	return math.Inf(1), fmt.Errorf("unsupported CPU usage format %s", cpu_str)
}

func cpu_node_sort() string {


	for {
		if !isSorting{
				isSorting = true
				node_sort = nodes_list[0]
				nodes_list = nodes_list[1:]
				nodes_list = append(nodes_list, node_sort)
			break
		}


	}
	isSorting = false
	return node_sort
	// we dont want race conditions
	// if isSorting {
	// 	internal_log("sorting in process")
	// 	for {
	// 		if node_sort != "" {
	// 			break
	// 		}
	// 	}
	// 	return node_sort
	// }
	// isSorting = true
	// internal_log("SORT_CPU_START")
	// var nodes NodeMetricList
	// data, err := kclient.RESTClient().Get().AbsPath("apis/metrics.k8s.io/v1beta1/nodes").Do(context.TODO()).Raw()
	// if err != nil {
	// 	internal_log("unable to retrieve metrics from nodes API - " + err.Error())
	// 	return ""
	// }
	// err = json.Unmarshal(data, &nodes)
	// if err != nil {
	// 	internal_log("unable to unmarshal metrics from nodes API - " + err.Error())
	// 	return ""
	// }
	// node_name := ""
	// var node_usage_minimum float64 = math.Inf(1)
	// for _, item := range nodes.Items {
	// 	cpu_current, err := memory_raw_to_float(item.Usage.Memory)
	// 	if err != nil {
	// 		internal_log("unable to convert cpu to float - " + err.Error())
	// 		return ""
	// 	}
	// 	mem_usage, _ := memory_raw_to_float(item.Usage.Memory)
	// 	internal_log("memory usage of "+ item.Metadata.Name+ " is "+strconv.Itoa(int(mem_usage)))
	// 	if cpu_current < node_usage_minimum {
	// 		node_usage_minimum = cpu_current
	// 		node_name = item.Metadata.Name
	// 	}
	// }
	// internal_log("SORT_CPU_END")
	// node_sort = node_name
	// isSorting = false
	// return node_name
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
	// log.Printf("\nendpoints of %s:", entrypoint)
	// log.Print(workflows[wf_name].Functions[entrypoint].Endpoints)
	for _, endpoint := range workflows[wf_name].Functions[entrypoint].Endpoints {
		if !slices.Contains(pod, endpoint) {
			pod = append(pod, endpoint)
			pod_list = append(pod_list, endpoint)
			// log.Print(pod_list)
			// log.Print(pod)

		}
	}
	return bfs_initial_pod(pod, wf_name, pod_list)
}
func createInitialPod(wf_name string) {
	var initial_pod []string

	var frontend_func string
	var endpoints []string
	func_endpoint := make(map[string][]string)
	for func_name, function := range workflows[wf_name].Functions {
		for _, endpoint := range function.Endpoints {
			func_endpoint[func_name] = append(func_endpoint[func_name], endpoint)
			if !slices.Contains(endpoints, endpoint) {
				if func_name != endpoint {
					endpoints = append(endpoints, endpoint)
				}
			}
		}
	}

	for func_name := range workflows[wf_name].Functions {
		if !slices.Contains(endpoints, func_name) {
			frontend_func = func_name
			break
		}

	}
	var pod_list []string
	pod_list = append(pod_list, frontend_func)
	initial_pod = bfs_initial_pod(initial_pod, wf_name, pod_list)
	workflows[wf_name].Pods = append(workflows[wf_name].Pods, initial_pod)
	workflows[wf_name].InitialPods = initial_pod
	log.Print(len(initial_pod))
	log.Print(workflows[wf_name].InitialPods)
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
	workflows[wf_name] = &workflow
	createInitialPod(wf_name)
	internal_log("CREATE_WORKFLOW_END - " + wf_name)
}


//to do 
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
	workflows[wf_name] = &workflow
	createInitialPod(wf_name)
	internal_log("UPDATE_WORKFLOW_END - " + wf_name)
}

func deleteWorkflow(wf_name string) {
	internal_log("DELETE_WORKFLOW_START - " + wf_name)
	_, exists := workflows[wf_name]
	if exists {
		internal_log("workflow " + wf_name + " exists")
		
	} else {
		internal_log("workflow " + wf_name + " does not exist")
	}
	internal_log("DELETE_WORKFLOW_END - " + wf_name)
}


//to-do 
func updateExistingIngress(wf_name string) {
	internal_log("UPDATE_EXISTING_START - " + wf_name)
	updateDeployments(wf_name)
	internal_log("UPDATE_EXISTING_END - " + wf_name)
}

func createNewIngress(wf_name string, rn int) string {
	internal_log("CREATE_INGRESS_START - " + wf_name)
	_, exist := workflows[wf_name]
	if !exist {
		internal_log("unable to create new ingress for " + wf_name + " - workflow does not exist")
		return ""
	}
	replicaNumber := strconv.Itoa(rn)
	internal_log("deploying replica number " + replicaNumber + " for workflow "+wf_name)
	ingress, err := manageDeployment(wf_name, replicaNumber)
	if err != nil {
		internal_log("Failed to deploy new ingress - " + err.Error())
		return ""
	}
	internal_log("CREATE_INGRESS_END - " + wf_name)
	return ingress
}


func (s *server) Deployment(ctx context.Context, req *pb.DeploymentServiceRequest) (*pb.DeploymentServiceReply, error) {
	wf_name := req.Name
	request_type := req.FunctionCall
	replicaNumber := req.ReplicaNumber
	var result string
	if request_type == "create" {
		internal_log("create workflow request start - " + wf_name)
		createWorkflow(wf_name, *req.Workflow)
		internal_log("create workflow request end - " + wf_name)
	} else if request_type == "update" {
		internal_log("update workflow request start - " + wf_name)
		updateWorkflow(wf_name, *req.Workflow)
		internal_log("update workflow request end - " + wf_name)
	} else if request_type == "delete" {
		internal_log("delete workflow request start - " + wf_name)
		deleteWorkflow(wf_name)
		internal_log("delete workflow request end - " + wf_name)
	} else if request_type == "existing_invoke" {
		internal_log("existing invoke request start - " + wf_name)
		updateExistingIngress(wf_name)
		internal_log("existing invoke request end - " + wf_name)
	} else if request_type == "new_invoke" {
		internal_log("new invoke request start - " + wf_name)
		result = createNewIngress(wf_name, int(replicaNumber))
		internal_log("new invoke request end - " + wf_name)
	}
	return &pb.DeploymentServiceReply{
		Message: fmt.Sprintf("%s", result),
	}, nil
}

func main() {
	internal_log("Ingress Controller Started")
	isSorting = false
	node_sort = ""
	deploymentRunning = false
	workflows = make(map[string]*Workflow)
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