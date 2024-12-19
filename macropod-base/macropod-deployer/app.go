package main

import (
	pb "app/deployer_pb"
	"context"
	"encoding/json"
	"errors"
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
	IngressVersion     map[string]int
	LatestVersion      int
	LastUpdated        time.Time
	Updating           bool
	InitialPods        []string
	FullyDisaggregated bool
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

var (
	kclient          *kubernetes.Clientset
	workflows        map[string]*Workflow
	versionFunction  map[string]int
	update_threshold int
	//deploymentRunning bool
	nodeCapacityCPU    map[string]float64
	nodeCapacityMemory map[string]float64
	//countLock     sync.Mutex
	updateDeployment sync.Mutex
	standbyNodeMap   map[string]string
)

func internal_log(message string) {
	fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.000000 UTC") + ": " + message)
}

func nodesAreStable(func_name string) bool {
	var nodes NodeMetricList
	data, err := kclient.RESTClient().Get().AbsPath("apis/metrics.k8s.io/v1beta1/nodes").Do(context.TODO()).Raw()
	if err != nil {
		log.Print("issue")
		return true
	}
	err = json.Unmarshal(data, &nodes)
	if err != nil {
		log.Print("issue")
		return true
	}
	for _, node := range nodes.Items {
		value, exists := node.Metadata.Labels["node-role.kubernetes.io/master"]
		if exists && value == "true" {
			continue
		}
		if node.Metadata.Name == standbyNodeMap[func_name] { //1 workflow assumption we can skip all the standby nodes in future
			continue
		}
		node_name := node.Metadata.Name
		usage_mem, _ := memory_raw_to_float(node.Usage.Memory)
		usage_cpu, _ := cpu_raw_to_float(node.Usage.CPU)
		percentage_cpu := (usage_cpu / nodeCapacityCPU[node_name]) * 100
		percentage_mem := (usage_mem / nodeCapacityMemory[node_name]) * 100
		if percentage_mem > 70 || percentage_cpu > 70 {
			log.Printf("Node %s is still not stable", node_name)
			return false
		}
	}
	return true
}

func getNodes() string {
	var nodes NodeMetricList
	data, err := kclient.RESTClient().Get().AbsPath("apis/metrics.k8s.io/v1beta1/nodes").Do(context.TODO()).Raw()
	if err != nil {
		internal_log("unable to retrieve metrics from nodes API - " + err.Error())
		return ""
	}
	_ = json.Unmarshal(data, &nodes)
	for _, node := range nodes.Items {
		value, exists := node.Metadata.Labels["node-role.kubernetes.io/master"]
		if exists && value == "true" {
			continue
		}
		pods, _ := kclient.CoreV1().Pods("macropod-functions").List(context.Background(), metav1.ListOptions{FieldSelector: "spec.nodeName=" + node.Metadata.Name})
		log.Printf("%d", len(pods.Items))
		if len(pods.Items) == 0 {
			return node.Metadata.Name
		}

	}
	return ""
}
func createStandByDeployment(func_name string, node_name string) (string, error) {
	var update_deployments []appsv1.Deployment
	namespace := "standby-functions"
	labels_ingress := map[string]string{
		"workflow_name":    func_name,
		"workflow_replica": func_name,
	}
	pathType := networkingv1.PathTypePrefix
	var iteratePods [][]string
	iteratePods = append(iteratePods, workflows[func_name].InitialPods) // standby is always fully reduced
	service_name_ingress := strings.ToLower(strings.ReplaceAll(iteratePods[0][0], "_", "-"))
	for _, pod := range iteratePods {
		pod_name := strings.ToLower(strings.ReplaceAll(pod[0], "_", "-"))
		service := &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      pod_name,
				Namespace: namespace,
			},
			Spec: corev1.ServiceSpec{
				Selector: map[string]string{
					"app": pod_name,
				},
				Ports: []corev1.ServicePort{},
			},
		}

		labels := map[string]string{
			"workflow_name":    func_name,
			"app":              pod_name,
			"workflow_replica": func_name,
			"version":          strconv.Itoa(workflows[func_name].LatestVersion),
		}
		replicaCount := int32(1)
		deployment := &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:   pod_name,
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
						Tolerations: []corev1.Toleration{
							{
								Key:      "workflow_standby",
								Value:    func_name,
								Effect:   "NoSchedule",
								Operator: corev1.TolerationOperator("Equal"),
							},
						},
						Containers: make([]corev1.Container, len(pod)),
					},
				},
			},
		}

		for i, container := range pod {
			container_name := strings.ToLower(strings.ReplaceAll(container, "_", "-"))
			func_port := 5000 + slices.Index(workflows[func_name].InitialPods, container)
			function := workflows[func_name].Functions[container]
			registry := function.Registry
			var env []corev1.EnvVar
			for name, value := range function.Envs {
				env = append(env, corev1.EnvVar{Name: name, Value: value})
			}
			env = append(env, corev1.EnvVar{Name: "SERVICE_TYPE", Value: "GRPC"})
			env = append(env, corev1.EnvVar{Name: "GRPC_THREAD", Value: strconv.Itoa(10)}) //TODO
			func_port_s := strconv.Itoa(func_port)
			env = append(env, corev1.EnvVar{Name: "FUNC_PORT", Value: func_port_s})
			log.Print(function.Endpoints)
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
					service_name = strings.ReplaceAll(strings.ToLower(endpoint_name), "_", "-") + "." + namespace + ".svc.cluster.local:" + endpoint_port
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
		_, exists := kclient.AppsV1().Deployments(namespace).Get(context.Background(), strings.ToLower(pod[0]), metav1.GetOptions{})
		if exists != nil {
			internal_log("Creating a new deployment " + deployment.Name)
			_, err := kclient.AppsV1().Deployments(namespace).Create(context.Background(), deployment, metav1.CreateOptions{})
			if err != nil {
				internal_log("unable to create deployment " + strings.ToLower(pod[0]) + " for " + namespace + " - " + err.Error())
				return "", err

			}

		} else {
			internal_log("Updating the existing deployment " + deployment.Name)
			update_deployments = append(update_deployments, *deployment)
		}
	}

	for _, dp := range update_deployments {
		internal_log("deploying existing deployment " + dp.Name)
		log.Print(dp)
		kclient.AppsV1().Deployments(namespace).Update(context.Background(), &dp, metav1.UpdateOptions{})
	}

	for _, pod := range workflows[func_name].Pods {
		pod_name := strings.ToLower(strings.ReplaceAll(pod[0], "_", "-"))
		labels := map[string]string{
			"workflow_name":    func_name,
			"app":              pod_name,
			"workflow_replica": func_name,
			"version":          strconv.Itoa(workflows[func_name].LatestVersion),
		}
		service := &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      pod_name,
				Namespace: namespace,
				Labels:    labels,
			},
			Spec: corev1.ServiceSpec{
				Selector: map[string]string{
					"app": pod_name,
				},
				Ports: []corev1.ServicePort{},
			},
		}
		for _, container := range pod {
			container_name := strings.ToLower(strings.ReplaceAll(container, "_", "-"))
			container_port := int32(5000 + slices.Index(workflows[func_name].InitialPods, container))
			service.Spec.Ports = append(service.Spec.Ports, corev1.ServicePort{
				Name:       container_name,
				Port:       container_port,
				TargetPort: intstr.FromInt(int(container_port)),
			})
		}

		_, exists := kclient.CoreV1().Services(namespace).Get(context.Background(), pod_name, metav1.GetOptions{})
		if exists != nil {
			_, err := kclient.CoreV1().Services(namespace).Create(context.Background(), service, metav1.CreateOptions{})
			if err != nil {
				internal_log("unable to create service " + strings.ToLower(pod[0]) + " for " + namespace + " - " + err.Error())
				return "", err
			}
		} else {
			_, err := kclient.CoreV1().Services(namespace).Update(context.Background(), service, metav1.UpdateOptions{})
			if err != nil {
				internal_log("unable to update service " + strings.ToLower(pod[0]) + " for " + namespace + " - " + err.Error())
				return "", err
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
					Host: func_name + "." + namespace + ".macropod",
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
	_, exists := kclient.NetworkingV1().Ingresses(namespace).Get(context.Background(), ingress.Name, metav1.GetOptions{})
	if exists != nil {
		_, err := kclient.NetworkingV1().Ingresses(namespace).Create(context.Background(), ingress, metav1.CreateOptions{})
		if err != nil {
			internal_log("unable to create ingress ")
			return "", err
		}
	} else {
		_, err := kclient.NetworkingV1().Ingresses(namespace).Update(context.Background(), ingress, metav1.UpdateOptions{})
		if err != nil {
			internal_log("unable to update service")
			return "", err
		}
	}
	return service_name_ingress + "." + namespace + ".svc.cluster.local:5000", nil
}

func manageDeployment(func_name string, replicaNumber string) (string, error) {
	log.Printf("deploying the workflow of version %d", workflows[func_name].LatestVersion)
	var update_deployments []appsv1.Deployment
	namespace := "macropod-functions"
	labels_ingress := map[string]string{
		"workflow_name":    func_name,
		"workflow_replica": func_name + "-" + replicaNumber,
		"version":          strconv.Itoa(workflows[func_name].LatestVersion),
	}
	updateDeployment.Lock()
	if _, exists := workflows[func_name]; exists {
		log.Printf(" %s is present", func_name)
	} else {
		log.Printf(" %s is not present", func_name)
		updateDeployment.Unlock()
		return "", errors.New("Workflow does not exist")
	}
	label_workflow := "workflow_name=" + func_name
	label_version := "version=" + strconv.Itoa(workflows[func_name].LatestVersion-1)
	labels_to_check := label_workflow + "," + label_version
	log.Printf("checking if older versions exist %s", labels_to_check)

	for {
		deployments_list, _ := kclient.CoreV1().Pods("macropod-functions").List(context.Background(), metav1.ListOptions{LabelSelector: labels_to_check})
		if deployments_list == nil || len(deployments_list.Items) == 0 {
			fmt.Println("Deployment does not exist")
			break
		}
		time.Sleep(10 * time.Millisecond) // let all the deployments be deleted before new ones
	}

	static := os.Getenv("STATIC")
	if static == "" {
		for !nodesAreStable(func_name) {
			fmt.Println("waiting for nodes to be stable")
			time.Sleep(1*time.Second)
		}
	}

	pathType := networkingv1.PathTypePrefix
	service_name_ingress := strings.ToLower(strings.ReplaceAll(workflows[func_name].Pods[0][0], "_", "-")) + "-" + replicaNumber
	for _, pod := range workflows[func_name].Pods {
		pod_name := strings.ToLower(strings.ReplaceAll(pod[0], "_", "-"))
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
			log.Print(function.Endpoints)
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
					service_name = strings.ReplaceAll(strings.ToLower(endpoint_service), "_", "-") + "-" + replicaNumber + "." + namespace + ".svc.cluster.local:" + endpoint_port
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
		// check if deployment with name ecists, if does not make a new one else update the existing lets start with creating new ones and then update the existing ones
		internal_log("Looking for " + pod[0] + "_" + replicaNumber)
		_, exists := kclient.AppsV1().Deployments(namespace).Get(context.Background(), strings.ToLower(pod[0])+"-"+replicaNumber, metav1.GetOptions{})
		if exists != nil {
			internal_log("Creating a new deployment " + deployment.Name)
			_, err := kclient.AppsV1().Deployments(namespace).Create(context.Background(), deployment, metav1.CreateOptions{})
			if err != nil {
				internal_log("unable to create deployment " + strings.ToLower(pod[0]) + " for " + namespace + " - " + err.Error())
				updateDeployment.Unlock()
				return "", err

			}

		} else {
			internal_log("Updating the existing deployment " + deployment.Name)
			update_deployments = append(update_deployments, *deployment)
		}
	}

	for _, dp := range update_deployments {
		internal_log("deploying existing deployment " + dp.Name)
		kclient.AppsV1().Deployments(namespace).Update(context.Background(), &dp, metav1.UpdateOptions{})
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
			container_port := int32(5000 + slices.Index(workflows[func_name].InitialPods, container))
			service.Spec.Ports = append(service.Spec.Ports, corev1.ServicePort{
				Name:       container_name,
				Port:       container_port,
				TargetPort: intstr.FromInt(int(container_port)),
			})
		}

		_, exists := kclient.CoreV1().Services(namespace).Get(context.Background(), pod_name+"-"+replicaNumber, metav1.GetOptions{})
		if exists != nil {
			_, err := kclient.CoreV1().Services(namespace).Create(context.Background(), service, metav1.CreateOptions{})
			if err != nil {
				internal_log("unable to create service " + strings.ToLower(pod[0]) + "-" + replicaNumber + " for " + namespace + " - " + err.Error())
				updateDeployment.Unlock()
				return "", err
			}
		} else {
			_, err := kclient.CoreV1().Services(namespace).Update(context.Background(), service, metav1.UpdateOptions{})
			if err != nil {
				internal_log("unable to update service " + strings.ToLower(pod[0]) + "-" + replicaNumber + " for " + namespace + " - " + err.Error())
				updateDeployment.Unlock()
				return "", err
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
					Host: func_name + "." + replicaNumber + "." + namespace + ".macropod",
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

	_, err := kclient.NetworkingV1().Ingresses(namespace).Create(context.Background(), ingress, metav1.CreateOptions{})
	if err != nil {
		internal_log("unable to create ingress for " + namespace + " - " + err.Error())
		updateDeployment.Unlock()
		return "", err
	}
	updateDeployment.Unlock()
	return service_name_ingress + "." + namespace + ".svc.cluster.local:5000", nil
}

// TODO
func updateDeployments(func_name string, max_concurrency int) string {
	var ingresses_deleted string
	if _, exists := workflows[func_name]; exists {
	} else {
		log.Printf(" %s is not present", func_name)
		return ingresses_deleted
	}

	if workflows[func_name].Updating {
		return ingresses_deleted
	}

	if time.Since(workflows[func_name].LastUpdated) < time.Second*time.Duration(update_threshold) {
		return ingresses_deleted
	}

	if workflows[func_name].FullyDisaggregated {
		return ingresses_deleted
	}

	workflows[func_name].Updating = true

	// get the percentage of the node utilisation
	var nodes NodeMetricList
	data, err := kclient.RESTClient().Get().AbsPath("apis/metrics.k8s.io/v1beta1/nodes").Do(context.TODO()).Raw()
	if err != nil {
		return ingresses_deleted
	}
	err = json.Unmarshal(data, &nodes)
	if err != nil {
		return ingresses_deleted
	}
	labels_to_check := ""
	static := os.Getenv("STATIC")
	if static != "" {
		return ""
	}
	for _, node := range nodes.Items {
		value, exists := node.Metadata.Labels["node-role.kubernetes.io/master"]
		if exists && value == "true" {
			continue
		}
		if node.Metadata.Name == standbyNodeMap[func_name] { //1 workflow assumption we can skip all the standby nodes in future
			continue
		}
		node_name := node.Metadata.Name
		usage_mem, _ := memory_raw_to_float(node.Usage.Memory)
		usage_cpu, _ := cpu_raw_to_float(node.Usage.CPU)
		percentage_cpu := (usage_cpu / nodeCapacityCPU[node_name]) * 100
		percentage_mem := (usage_mem / nodeCapacityMemory[node_name]) * 100
		if percentage_mem > 70 || percentage_cpu > 70 {
			updateDeployment.Lock()
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
			max_concurrency *= 2
			workflows[func_name].Updating = false
			label_workflow := "workflow_name=" + func_name
			label_version := "version=" + strconv.Itoa(workflows[func_name].LatestVersion-1)
			labels_to_check = label_workflow + "," + label_version
			updateDeployment.Unlock()
			return labels_to_check //ingress controller will delete the resources based on the usaage
		}
	}

	if _, exists := workflows[func_name]; exists {
		workflows[func_name].Updating = false
	} else {
		log.Printf(" %s is not present", func_name)
		return ingresses_deleted
	}
	return labels_to_check
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

func createInitialPod(func_name string) string {
	standbyNode := getNodes()
	standbyNodeMap[func_name] = standbyNode
	node, _ := kclient.CoreV1().Nodes().Get(context.Background(), standbyNode, metav1.GetOptions{})
	node.Spec.Taints = append(node.Spec.Taints, corev1.Taint{
		Key:    "workflow_standby",
		Value:  func_name,
		Effect: "NoSchedule",
	})
	kclient.CoreV1().Nodes().Update(context.Background(), node, metav1.UpdateOptions{})
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
		fmt.Println("original: ")
		workflows[func_name].Pods = pods_disagg
	} else {
		workflows[func_name].Pods = append(workflows[func_name].Pods, initial_pod)
	}
	log.Printf("we are going to deploy %v", workflows[func_name].Pods)
	workflows[func_name].InitialPods = initial_pod
	workflows[func_name].LatestVersion = 1
	createStandByDeployment(func_name, standbyNode)
	return standbyNode
}

func createWorkflow(func_name string, func_str string) string {
	workflow := Workflow{}
	json.Unmarshal([]byte(func_str), &workflow)
	_, exists := workflows[func_name]
	if exists {
		internal_log("workflow " + func_name + " already exists. If you are updating it please use update instead.")
		return ""
	}
	workflows[func_name] = &workflow
	standby := createInitialPod(func_name)
	return standby
}

func updateWorkflow(func_name string, workflow_str string) string {
	workflow := Workflow{}
	json.Unmarshal([]byte(workflow_str), &workflow)
	_, exists := workflows[func_name]
	if exists {
		delete(workflows, func_name)
		delete(standbyNodeMap, func_name)
	}
	workflows[func_name] = &workflow
	standby := createInitialPod(func_name)
	return standby
}

func deleteWorkflow(func_name string) {
	updateDeployment.Lock()
	_, exists := workflows[func_name]
	node_name := standbyNodeMap[func_name]
	node, _ := kclient.CoreV1().Nodes().Get(context.Background(), node_name, metav1.GetOptions{})
	taint := corev1.Taint{
		Key:    "workflow_standby",
		Value:  func_name,
		Effect: "NoSchedule",
	}
	index := slices.Index(node.Spec.Taints, taint)
	if index != -1 {
		node.Spec.Taints = append(node.Spec.Taints[:index], node.Spec.Taints[index+1:]...)
	}
	log.Printf("Removing tain from %s in %s", func_name, node_name)
	kclient.CoreV1().Nodes().Update(context.Background(), node, metav1.UpdateOptions{})
	if exists {
		delete(workflows, func_name)
		delete(standbyNodeMap, func_name)
	}
	updateDeployment.Unlock()
}

func updateExistingIngress(func_name string, current_concurrency int) string {
	ingress_deleted := updateDeployments(func_name, current_concurrency)
	return ingress_deleted
}

func createNewIngress(func_name string, rn int) string {
	_, exist := workflows[func_name]
	if !exist {
		internal_log("unable to create new ingress for " + func_name + " - workflow does not exist")
		return ""
	}
	replicaNumber := strconv.Itoa(rn)
	internal_log("deploying replica number " + replicaNumber + " for workflow " + func_name)
	ingress, err := manageDeployment(func_name, replicaNumber)
	if err != nil {
		internal_log("Failed to deploy new ingress - " + err.Error())
		return ""
	}
	return ingress
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
					delete(nodeCapacityCPU, node.Name)
					delete(nodeCapacityMemory, node.Name)
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
			nodeCapacityCPU[node.Name] = cpu_float
			nodeCapacityMemory[node.Name] = mem_float
			// log.Printf("%v", nodeCapacityCPU)
			// log.Printf("%v", nodeCapacityMemory)
		}
		time.Sleep(30 * time.Minute)
	}

}

func (s *server) Deployment(ctx context.Context, req *pb.DeploymentServiceRequest) (*pb.DeploymentServiceReply, error) {
	func_name := req.Name
	request_type := req.FunctionCall
	replicaNumber := req.ReplicaNumber
	var result string
	if request_type == "create" {
		result = createWorkflow(func_name, *req.Workflow)
	} else if request_type == "update" {
		result = updateWorkflow(func_name, *req.Workflow)
	} else if request_type == "delete" {
		deleteWorkflow(func_name)
	} else if request_type == "existing_invoke" {
		result = updateExistingIngress(func_name, int(replicaNumber)) //replicaNumber here is the currentc_concurrency
	} else if request_type == "new_invoke" {
		result = createNewIngress(func_name, int(replicaNumber))
	}
	return &pb.DeploymentServiceReply{
		Message: result,
	}, nil
}

func main() {
	internal_log("Deployer Started")
	// deploymentRunning = false
	update_threshold, _ = strconv.Atoi(os.Getenv("UPDATE_THRESHOLD"))
	workflows = make(map[string]*Workflow)
	standbyNodeMap = make(map[string]string)
	versionFunction = make(map[string]int)
	nodeCapacityCPU = make(map[string]float64)
	nodeCapacityMemory = make(map[string]float64)
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
	go checkNodeStatus()
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
