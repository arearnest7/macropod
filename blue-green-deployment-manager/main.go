package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc"
	pb "main/pb"
	"io"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
	"math"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

type server struct {
	pb.DeploymentServiceServer
}

// // this will also take the ping metrics for latancy in pod to pod communication and

// // we will retrievce 2 things - metrics and then ping results from pod to pod communication and build on that
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

var (
	clientset          *kubernetes.Clientset
	err                error
	namesapce_ingress  string
	function_timestamp map[string]time.Time
	ttl_seconds        int // time in seconds
	version_function   map[string]int
)

func init() {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	version_function = make(map[string]int)
	function_timestamp = make(map[string]time.Time)
	ttl_seconds, _ = strconv.Atoi(os.Getenv("TTL"))
	namesapce_ingress = os.Getenv("NAMESPACE_INGRESS")
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
}

func convertCPUUsage(cpuUsage string) (float64, error) {

	if cpuUsage == "0" {
		return 0, nil
	}
	if strings.HasSuffix(cpuUsage, "n") {
		cpuUsage = strings.TrimSuffix(cpuUsage, "n")
		cpu, err := strconv.ParseFloat(cpuUsage, 64)
		if err != nil {
			return 0, err
		}

		return cpu / 1000000, nil
	} else if strings.HasSuffix(cpuUsage, "m") {
		cpuUsage = strings.TrimSuffix(cpuUsage, "m")
		cpu, err := strconv.ParseFloat(cpuUsage, 64)
		if err != nil {
			return 0, err
		}

		return cpu, nil
	} else if strings.HasSuffix(cpuUsage, "u") {

		cpuUsage = strings.TrimSuffix(cpuUsage, "m")
		cpu, err := strconv.ParseFloat(cpuUsage, 64)
		if err != nil {
			return 0, err
		}

		return cpu / 1000, nil

	} else {
		return 0, fmt.Errorf("unsupported CPU usage format %s", cpuUsage)
	}
}

// continously check TTL if TTL is above certain seconds delete the namespace
func checkTTL() {
	for {
		//log.Print(function_timestamp.data)
		currentTime := time.Now()
		for name, timestamp := range function_timestamp {
			elapsedTime := currentTime.Sub(timestamp)
			if elapsedTime.Seconds() > float64(ttl_seconds) {
				version := version_function[name]
				versionStr := strconv.Itoa(version)
				namespace := name + "-" + versionStr
				log.Printf("Deleting %s because of TTL", namespace)
				DeleteNamespace(namespace)
				delete(function_timestamp, name)
				version += 1
				versionStr = strconv.Itoa(version)
				// delete both blue  and green version
				namespace = name + "-" + versionStr
				_, exists := clientset.CoreV1().Namespaces().Get(context.Background(), namespace, metav1.GetOptions{})
				if exists == nil {
					log.Printf("Deleting %s because of TTL", namespace)
					DeleteNamespace(namespace)
				}
			}
		}
		time.Sleep(time.Second)
	}
}

// convert memory given by metric server to float for comaprisions
func convertMemoryUsage(memoryUsage string) (float64, error) {
	if memoryUsage == "0" {
		return 0, nil
	}
	if strings.HasSuffix(memoryUsage, "Ki") {
		memoryUsage = strings.TrimSuffix(memoryUsage, "Ki")
		memory, err := strconv.ParseFloat(memoryUsage, 64)
		if err != nil {
			return 0, err
		}
		return memory * 1024, nil
	} else if strings.HasSuffix(memoryUsage, "Mi") {
		memoryUsage = strings.TrimSuffix(memoryUsage, "Mi")
		memory, err := strconv.ParseFloat(memoryUsage, 64)
		if err != nil {
			return 0, err
		}
		return memory * 1024 * 1024, nil
	} else {
		return 0, fmt.Errorf("unsupported memory usage format")
	}
}

func getMetricsNodes(node *NodeMetricList) error {

	result := clientset.RESTClient().Get().AbsPath("apis/metrics.k8s.io/v1beta1/nodes").Do(context.TODO())
	data, err := result.Raw()
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &node)
	return err
}

func nodeCPUSort() string {

	var nodes NodeMetricList
	err := getMetricsNodes(&nodes)
	if err != nil {

		return ""

	}

	node_name := ""
	var node_usage_minimum float64 = math.Inf(1)

	for _, item := range nodes.Items {
		cpu_current, _ := convertCPUUsage(item.Usage.CPU)
		if cpu_current < node_usage_minimum {
			node_usage_minimum = cpu_current
			node_name = item.Metadata.Name

		}
	}
	return node_name

}

func getMetrics(clientset *kubernetes.Clientset, pods *PodMetricsList, namespace string) error {
	result := clientset.RESTClient().Get().AbsPath("apis/metrics.k8s.io/v1beta1/pods").Do(context.TODO())
	data, err := result.Raw()
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &pods)
	//log.Print(pods)
	return err
}

func getPodMetricsAndChanges(namespace string) (float64, float64, error) {
	_, exists := clientset.CoreV1().Namespaces().Get(context.Background(), namespace, metav1.GetOptions{})
	if exists != nil {
		log.Printf("%s worflow version does not exist", namespace)
		return 0.0, 0.0, nil
	}
	var podMetricsList PodMetricsList
	err := getMetrics(clientset, &podMetricsList, namespace)
	if err != nil {
		return 0, 0, err
	}
	var nsCPUUsage float64 = 0.0
	var nsMemUsage float64 = 0.0
	for _, item := range podMetricsList.Items {
		podName := item.Metadata.Name
		podNamespace := item.Metadata.Namespace
		if podNamespace != namespace {
			continue
		}
		var podCPUUsage float64
		var podMemUsage float64
		for _, container := range item.Containers {
			mem, err := convertMemoryUsage(container.Usage.Memory)
			if err != nil {
				fmt.Printf("Error converting CPU usage for pod %s: %v\n", podName, err)
				return 0.0, 0.0, err
			}
			cpu, err := convertCPUUsage(container.Usage.CPU)
			if err != nil {
				fmt.Printf("Error converting CPU usage for pod %s: %v\n", podName, err)
				return 0.0, 0.0, err
			}
			podCPUUsage += cpu
			podMemUsage += mem
		}

		nsCPUUsage += podCPUUsage
		nsMemUsage += podMemUsage

	}

	log.Printf("Namespace : %s\n", namespace)
	log.Printf("CPU Usage:%v\n", nsCPUUsage)
	log.Printf("Memory Usage:%v\n", nsMemUsage)
	return nsCPUUsage, nsMemUsage, nil
}

func getConfigMapData(namespace, configMapName string) ([]map[string]interface{}, error) {
	cm, err := clientset.CoreV1().ConfigMaps(namespace).Get(context.TODO(), configMapName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	var configData []map[string]interface{}
	if err := json.Unmarshal([]byte(cm.Data["my-config.json"]), &configData); err != nil {
		return nil, err
	}
	//log.Print(configData)
	return configData, nil
}

func getConfigMapThreshold(namespace, configMapName string) (float64, float64, float64, float64, error) {
	cm, err := clientset.CoreV1().ConfigMaps(namespace).Get(context.TODO(), configMapName, metav1.GetOptions{})
	if err != nil {
		return 0.0, 0.0, 0.0, 0.0, err
	}

	cpuThreshold1Str := cm.Labels["cpu_threshold_1"]
	cpuThreshold2Str := cm.Labels["cpu_threshold_2"]
	mmThreshold1Str := cm.Labels["mm_threshold_1"]
	mmThreshold2Str := cm.Labels["mm_threshold_2"]

	cpu_threshold1, _ := strconv.ParseFloat(cpuThreshold1Str, 64)

	cpu_threshuld2, _ := strconv.ParseFloat(cpuThreshold2Str, 64)

	mm_threshold1, _ := strconv.ParseFloat(mmThreshold1Str, 64)

	mm_threshold2, _ := strconv.ParseFloat(mmThreshold2Str, 64)

	return cpu_threshold1, cpu_threshuld2, mm_threshold1, mm_threshold2, nil

}

func getKind(configMapName string) string {

	cm, _ := clientset.CoreV1().ConfigMaps(namesapce_ingress).Get(context.TODO(), configMapName, metav1.GetOptions{})

	if kind, ok := cm.Labels["kind"]; ok {
		return kind
	}
	return ""

}

// check for changes in config map configurations provided by uiser -> if the user update sthe existing depkoyment it needs to be changed toooo
// keep running continuously
func watchConfigMaps() {
	for {
		watcher, err := clientset.CoreV1().ConfigMaps(namesapce_ingress).Watch(context.TODO(), metav1.ListOptions{})
		if err != nil {
			log.Fatalf("Failed to set up watch: %s", err)
		}

		for event := range watcher.ResultChan() {
			configmap, ok := event.Object.(*corev1.ConfigMap)
			if !ok {
				log.Printf("Expected Ingress type, got %T", event.Object)
				continue
			}

			switch event.Type {
			case watch.Modified:
				log.Printf("Updated host targets: %+v", configmap)
				updateDeployment(configmap)
			}

		}
	}
}

func updateDeployment(configMap *corev1.ConfigMap) {
	namespace := configMap.Name
	func_name := configMap.Labels["function_name"]
	log.Printf("ConfigMap for %s updates", namespace)
	_, exists := clientset.CoreV1().Namespaces().Get(context.Background(), namespace, metav1.GetOptions{})
	if exists != nil {
		//does not exist so ignore
		return
	} else {
		log.Printf("ConfigMap for %s updates", namespace)
		MakeDeployment(namespace, func_name, false, true)
		MakeDeployment(namespace, func_name, true, true)

	}

}

func MakeDeploymentSinglePod(kind string, namespace string, func_name string, ingress bool, update bool) error {

	log.Print("Deploying Single pod")

	configMapName := namespace

	configDataArray, err := getConfigMapData(namesapce_ingress, configMapName)
	if err != nil {
		log.Printf("Falied to delete deployment: %v\n", err)
	}

	// first lets make deplouyment and then service

	namespace_object := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
	}

	_, exists := clientset.CoreV1().Namespaces().Get(context.Background(), namespace, metav1.GetOptions{})
	if exists != nil {
		_, err = clientset.CoreV1().Namespaces().Create(context.Background(), namespace_object, metav1.CreateOptions{})
		if err != nil {
			panic(err.Error())
		}
		log.Print("Namespace created successfuly")
	}
	if ingress == false {
		// deployment := appsv1.Deployment{
		// 	ObjectMeta: metav1.ObjectMeta{
		// 		Name: namespace,
		// 	},
		// 	Spec: appsv1.DeploymentSpec{
		// 		Selector: &metav1.LabelSelector{
		// 			MatchLabels: map[string]string{"app": namespace},
		// 		},
		// 		Template: corev1.PodTemplateSpec{
		// 			ObjectMeta: metav1.ObjectMeta{
		// 				Labels: map[string]string{"app": namespace},
		// 			},
		// 			Spec: corev1.PodSpec{
		// 				Containers: make([]corev1.Container, len(configDataArray)),
		// 			},
		// 		},
		// 	},
		// }
		service := &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      namespace,
				Namespace: namespace,
			},
			Spec: corev1.ServiceSpec{
				Selector: map[string]string{
					"app": namespace,
				},
				Ports: []corev1.ServicePort{},
			},
		}
		// if kind == "mmap" {
		// 	volume := corev1.Volume{
		// 		Name: "macropod-pv",
		// 		VolumeSource: corev1.VolumeSource{
		// 			EmptyDir: &corev1.EmptyDirVolumeSource{
		// 				Medium: corev1.StorageMediumMemory,
		// 			},
		// 		},
		// 	}

		// 	if deployment.Spec.Template.Spec.Volumes == nil {
		// 		deployment.Spec.Template.Spec.Volumes = make([]corev1.Volume, 0)
		// 	}

		// 	deployment.Spec.Template.Spec.Volumes = append(deployment.Spec.Template.Spec.Volumes, volume)

		// }

		// log.Print(configDataArray)

		for _, configMapData := range configDataArray {
			name := configMapData["name"].(string)
			servicePort := int32(configMapData["service"].(map[string]interface{})["port"].(float64))

			containerPort := int32(configMapData["service"].(map[string]interface{})["targetPort"].(float64))
			service.Spec.Ports = append(service.Spec.Ports, corev1.ServicePort{
				Name:       name,
				Port:       servicePort,
				TargetPort: intstr.FromInt(int(containerPort)),
			})

		}

		deployment := appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name: namespace,
			},
			Spec: appsv1.DeploymentSpec{
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{"app": namespace},
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{"app": namespace},
					},
					Spec: corev1.PodSpec{
						Containers: make([]corev1.Container, len(configDataArray)),
					},
				},
			},
		}

		if kind == "mmap" {
			volume := corev1.Volume{
				Name: "macropod-pv",
				VolumeSource: corev1.VolumeSource{
					EmptyDir: &corev1.EmptyDirVolumeSource{
						Medium: corev1.StorageMediumMemory,
					},
				},
			}

			if deployment.Spec.Template.Spec.Volumes == nil {
				deployment.Spec.Template.Spec.Volumes = make([]corev1.Volume, 0)
			}

			deployment.Spec.Template.Spec.Volumes = append(deployment.Spec.Template.Spec.Volumes, volume)

		}

		// log.Print(configDataArray)

		for i, configMapData := range configDataArray {
			name := configMapData["name"].(string)
			//log.Print(name)
			replicaCount := int32(configMapData["replicaCount"].(float64))

			envVariables, _ := configMapData["env"].([]interface{})

			imageData := configMapData["image"].(map[string]interface{})
			imageName := imageData["image"].(string)
			endpoints, ok := configMapData["endpoints"].(string)
			if !ok {
				endpoints = ""
			}
			containerPort := int32(configMapData["service"].(map[string]interface{})["targetPort"].(float64))

			//log.Printf("EndpointsList %s", endpoints)
			imagePullPolicy := corev1.PullPolicy(imageData["pullPolicy"].(string))

			var env []corev1.EnvVar
			for _, item := range envVariables {
				envData, _ := item.(map[string]interface{})

				name, _ := envData["name"].(string)
				value, _ := envData["value"].(string)

				env = append(env, corev1.EnvVar{
					Name:  name,
					Value: value,
				})
			}

			endpointList := strings.Split(endpoints, ",")
			if endpoints != "" {
				for _, endpoint := range endpointList {
					name_key := strings.ToUpper(endpoint)
					port := ""
					for _, port_svc := range service.Spec.Ports {
						if port_svc.Name == endpoint {
							port = strconv.Itoa(int(port_svc.Port))
							break
						}
					}
					service_name := "127.0.0.1:" + port
					final_name := strings.ReplaceAll(name_key, "-", "_")
					env = append(env, corev1.EnvVar{
						Name:  final_name,
						Value: service_name,
					})
				}
			}
			deployment.Spec.Replicas = &replicaCount
			deployment.Spec.Template.Spec.Containers[i] = corev1.Container{
				Name:            name,
				Image:           imageName,
				ImagePullPolicy: imagePullPolicy,
				Ports: []corev1.ContainerPort{
					{
						ContainerPort: containerPort,
					},
				},
				Env: env,
			}

			if kind == "mmap" {
				volumeMount := corev1.VolumeMount{
					Name:      "macropod-pv",
					MountPath: "/macropod-pv",
				}
				deployment.Spec.Template.Spec.Containers[i].VolumeMounts = append(
					deployment.Spec.Template.Spec.Containers[i].VolumeMounts, volumeMount)
			}
			command, ok := configMapData["command"].(string)
			if !ok {
				command = ""
			}
			if command != "" {
				commandList := strings.Split(command, ",")
				//log.Print(commandList)
				deployment.Spec.Template.Spec.Containers[i].Command = commandList

			}
			args, ok := configMapData["args"].(string)
			if !ok {
				args = ""
			}
			if args != "" {
				argsList := strings.Split(args, ",")
				//log.Print(argsList)
				deployment.Spec.Template.Spec.Containers[i].Args = argsList

			}

		}

		node_name := nodeCPUSort()
		if node_name != "" {
			log.Printf("Assigning node %s", node_name)
			deployment.Spec.Template.Spec.NodeSelector = map[string]string{
				"kubernetes.io/hostname": node_name,
			}
		}

		// log.Print("Deploying...............")
		// log.Print(deployment)
		// log.Print("..............................")

		if update == false {
			log.Printf("Creating deployment and service %s in %s", namespace, namespace)
			_, err = clientset.AppsV1().Deployments(namespace).Create(context.Background(), &deployment, metav1.CreateOptions{})
			if err != nil {
				panic(err.Error())
			}
		} else {
			log.Print("Updating deployment and service %s in %s", namespace, namespace)
			_, err = clientset.AppsV1().Deployments(namespace).Update(context.Background(), &deployment, metav1.UpdateOptions{})
			if err != nil {
				panic(err.Error())
			}
		}
		if update == false {
			_, err = clientset.CoreV1().Services(namespace).Create(context.Background(), service, metav1.CreateOptions{})
			if err != nil {
				panic(err.Error())
			}
		} else {
			_, err = clientset.CoreV1().Services(namespace).Update(context.Background(), service, metav1.UpdateOptions{})
			if err != nil {
				panic(err.Error())
			}
		}
	}
	for _, configMapData := range configDataArray {

		name := configMapData["name"].(string)

		labels := map[string]string{
			"function_name": func_name,
			"app":           name,
		}

		// log.Print(labels)
		//containerPort := int32(configMapData["service"].(map[string]interface{})["targetPort"].(float64))
		// log.Print(containerPort)
		servicePort := int32(configMapData["service"].(map[string]interface{})["port"].(float64))
		// log.Print(servicePort)

		pathType := networkingv1.PathTypePrefix
		// create namespace first

		// create deployment and service if ingress is false else make ingress resource

		if ingress {
			if _, ok := configMapData["ingress"]; ok {

				hostName := configMapData["ingress"].(map[string]interface{})["host"].(string)
				ingress := &networkingv1.Ingress{
					ObjectMeta: metav1.ObjectMeta{
						Name:      name,
						Namespace: namespace,
						Labels:    labels,
					},
					Spec: networkingv1.IngressSpec{
						Rules: []networkingv1.IngressRule{
							{
								Host: hostName,
								IngressRuleValue: networkingv1.IngressRuleValue{
									HTTP: &networkingv1.HTTPIngressRuleValue{
										Paths: []networkingv1.HTTPIngressPath{
											{
												Path:     "/",
												PathType: &pathType,
												Backend: networkingv1.IngressBackend{
													Service: &networkingv1.IngressServiceBackend{
														Name: namespace,
														Port: networkingv1.ServiceBackendPort{
															Number: servicePort,
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
				if update == false {
					log.Printf("Creating ingress for %s in %s", name, namespace)
					_, err = clientset.NetworkingV1().Ingresses(namespace).Create(context.Background(), ingress, metav1.CreateOptions{})
					if err != nil {
						panic(err.Error())
					}
				} else {
					log.Printf("Updating ingress for %s in %s", name, namespace)
					_, err = clientset.NetworkingV1().Ingresses(namespace).Update(context.Background(), ingress, metav1.UpdateOptions{})
					if err != nil {
						panic(err.Error())
					}
				}

			}
		}

	}

	return nil

}

func MakeDeploymentMultiPod(kind string, namespace string, func_name string, ingress bool, update bool) error {
	configMapName := namespace

	configDataArray, err := getConfigMapData(namesapce_ingress, configMapName)
	if err != nil {
		log.Printf("Falied to delete deployment: %v\n", err)
	}
	//log.Print(configDataArray)
	for _, configMapData := range configDataArray {

		name := configMapData["name"].(string)
		// log.Print(name)
		replicaCount := int32(configMapData["replicaCount"].(float64))
		// log.Print(replicaCount)
		envVariables, _ := configMapData["env"].([]interface{})
		// log.Print(envVariables)
		imageData := configMapData["image"].(map[string]interface{})
		// log.Print(imageData)
		imageName := imageData["image"].(string)
		endpoints, ok := configMapData["endpoints"].(string)
		if !ok {
			endpoints = ""
		}
		log.Printf("EndpointsList %s", endpoints)
		// log.Print(imageName)
		imagePullPolicy := corev1.PullPolicy(imageData["pullPolicy"].(string))
		// log.Print(imagePullPolicy)
		var env []corev1.EnvVar
		for _, item := range envVariables {
			envData, _ := item.(map[string]interface{})

			name, _ := envData["name"].(string)
			value, _ := envData["value"].(string)

			env = append(env, corev1.EnvVar{
				Name:  name,
				Value: value,
			})
		}
		// log.Print(env)
		labels := map[string]string{
			"function_name": func_name,
			"app":           name,
		}

		// log.Print(labels)
		containerPort := int32(configMapData["service"].(map[string]interface{})["targetPort"].(float64))
		// log.Print(containerPort)
		servicePort := int32(configMapData["service"].(map[string]interface{})["port"].(float64))
		// log.Print(servicePort)

		pathType := networkingv1.PathTypePrefix
		// create namespace first
		endpointList := strings.Split(endpoints, ",")

		if endpoints != "" {
			for _, endpoint := range endpointList {
				service_name := endpoint + "." + namespace + "." + "svc.cluster.local"
				name_key := strings.ToUpper(endpoint)
				final_name := strings.ReplaceAll(name_key, "-", "_")
				env = append(env, corev1.EnvVar{
					Name:  final_name,
					Value: service_name,
				})
			}
		}
		namespace_object := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: namespace,
			},
		}

		_, exists := clientset.CoreV1().Namespaces().Get(context.Background(), namespace, metav1.GetOptions{})
		if exists != nil {
			_, err = clientset.CoreV1().Namespaces().Create(context.Background(), namespace_object, metav1.CreateOptions{})
			if err != nil {
				panic(err.Error())
			}
			log.Print("Namespace created successfuly")
		}

		// create deployment and service if ingress is false else make ingress resource

		if ingress {
			if _, ok := configMapData["ingress"]; ok {

				hostName := configMapData["ingress"].(map[string]interface{})["host"].(string)
				ingress := &networkingv1.Ingress{
					ObjectMeta: metav1.ObjectMeta{
						Name:      name,
						Namespace: namespace,
						Labels:    labels,
					},
					Spec: networkingv1.IngressSpec{
						Rules: []networkingv1.IngressRule{
							{
								Host: hostName,
								IngressRuleValue: networkingv1.IngressRuleValue{
									HTTP: &networkingv1.HTTPIngressRuleValue{
										Paths: []networkingv1.HTTPIngressPath{
											{
												Path:     "/",
												PathType: &pathType,
												Backend: networkingv1.IngressBackend{
													Service: &networkingv1.IngressServiceBackend{
														Name: name,
														Port: networkingv1.ServiceBackendPort{
															Number: servicePort,
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
				if update == false {
					log.Printf("Creating ingress for %s in %s", name, namespace)
					_, err = clientset.NetworkingV1().Ingresses(namespace).Create(context.Background(), ingress, metav1.CreateOptions{})
					if err != nil {
						panic(err.Error())
					}
				} else {
					log.Printf("Updating ingress for %s in %s", name, namespace)
					_, err = clientset.NetworkingV1().Ingresses(namespace).Update(context.Background(), ingress, metav1.UpdateOptions{})
					if err != nil {
						panic(err.Error())
					}
				}

			}
		} else {

			log.Print("Creating deployment and service %s in %s", name, namespace)
			node_name := nodeCPUSort()
			deployment := &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:   name,
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
							Containers: []corev1.Container{
								{
									Name:            name,
									Image:           imageName,
									ImagePullPolicy: imagePullPolicy,
									Ports: []corev1.ContainerPort{
										{
											ContainerPort: containerPort,
										},
									},
									Env: env,
								},
							},
						},
					},
				},
			}
			command, ok := configMapData["command"].(string)
			if !ok {
				command = ""
			}
			if command != "" {
				commandList := strings.Split(command, ",")
				deployment.Spec.Template.Spec.Containers[0].Command = commandList

			}
			args, ok := configMapData["args"].(string)
			if !ok {
				args = ""
			}
			if args != "" {
				argsList := strings.Split(args, ",")
				deployment.Spec.Template.Spec.Containers[0].Args = argsList

			}

			if node_name != "" {
				log.Printf("Assigning node %s", node_name)
				deployment.Spec.Template.Spec.NodeSelector = map[string]string{
					"kubernetes.io/hostname": node_name,
				}
			}

			if update == false {
				log.Printf("Creating deployment and service %s in %s", name, namespace)
				_, err = clientset.AppsV1().Deployments(namespace).Create(context.Background(), deployment, metav1.CreateOptions{})
				if err != nil {
					panic(err.Error())
				}
			} else {
				log.Print("Updating deployment and service %s in %s", name, namespace)
				_, err = clientset.AppsV1().Deployments(namespace).Update(context.Background(), deployment, metav1.UpdateOptions{})
				if err != nil {
					panic(err.Error())
				}
			}

			// deploy a service as well

			service := &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
				},
				Spec: corev1.ServiceSpec{
					Selector: map[string]string{
						"app": name,
					},
					Ports: []corev1.ServicePort{
						{
							Name:       "main-port",
							Port:       servicePort,
							TargetPort: intstr.FromInt(int(containerPort)),
						},
					},
				},
			}
			if update == false {
				_, err = clientset.CoreV1().Services(namespace).Create(context.Background(), service, metav1.CreateOptions{})
				if err != nil {
					panic(err.Error())
				}
			} else {
				_, err = clientset.CoreV1().Services(namespace).Update(context.Background(), service, metav1.UpdateOptions{})
				if err != nil {
					panic(err.Error())
				}
			}

		}

	}

	return nil

}
func MakeDeployment(namespace string, func_name string, ingress bool, update bool) error {

	configMapName := namespace
	kind := getKind(configMapName)

	if kind == "mmap" || kind == "single-pod" {
		log.Print("Single Pod")
		err := MakeDeploymentSinglePod(kind, namespace, func_name, ingress, update)
		return err

	} else {
		log.Print("Multi-Pod")
		err := MakeDeploymentMultiPod(kind, namespace, func_name, ingress, update)
		return err

	}

}

func DeleteNamespace(namespace string) error {
	log.Print("Deleteing namespace %s\n", namespace)
	err = clientset.CoreV1().Namespaces().Delete(context.TODO(), namespace, metav1.DeleteOptions{})
	return err
}

func metricEvalHandler(func_name string) string {
	function_timestamp[func_name] = time.Now()
	log.Printf("Evaluation metrics of %s", func_name)
	version := version_function[func_name]
	versionStr := strconv.Itoa(version)
	namespace_existing := func_name + "-" + versionStr
	// we have config maps stores in same namesapce as ingress and with name <func_name>-<version>
	cpu_threshold1, cpu_threshold2, mm_threshold1, mm_threshold2, _ := getConfigMapThreshold(namesapce_ingress, namespace_existing)
	var cpu_usage, mm_usage float64
	for {
		cpu_usage, mm_usage, err = getPodMetricsAndChanges(namespace_existing)
		if err == nil {
			break
		}
	}
	if cpu_usage > cpu_threshold1 || mm_usage > mm_threshold1 {
		log.Print("Threshold 1 reached")
		version_update := version + 1
		namespace_update := func_name + "-" + strconv.Itoa(version_update)
		_, exists := clientset.CoreV1().Namespaces().Get(context.Background(), namespace_update, metav1.GetOptions{})
		if exists != nil {
			//namespace is not there
			err := MakeDeployment(namespace_update, func_name, false, false)
			if err != nil {
				log.Printf("Falied to install deployment: %v\n", err)
			}
		}
		if cpu_usage > cpu_threshold2 || mm_usage > mm_threshold2 {
			log.Print("Threshold 2 reached")

			// lets check if the deploymnets are up then only sfift ingress
			deployments, err := clientset.AppsV1().Deployments(namespace_update).List(context.TODO(), metav1.ListOptions{})
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error listing deployments: %v\n", err)
			}
			//log.Print(deployments)
			allRunning := true
			for _, deployment := range deployments.Items {
				if deployment.Status.ReadyReplicas == 0 {
					fmt.Printf("Deployment %s is not running\n", deployment.Name)
					allRunning = false
					break
				}
			}

			if allRunning {

				err = MakeDeployment(namespace_update, func_name, true, false)
				if err != nil {
					log.Printf("Falied to install ingress: %v\n", err)
				}
				err = DeleteNamespace(namespace_existing)
				if err != nil {
					log.Printf("Falied to delete older version: %v\n", err)
				}
				version_function[func_name] = version_update
			}

		}
	}

	return "Evaluated metrics"

}

// this will look into the function name and check if version 0 is available - if yes, move ahead and deploy - this will also keep track of the versions and functions that are being handled
func makeNewFunctionHandler(func_name string) string {
	function_timestamp[func_name] = time.Now()
	log.Printf("Deploying new function.....%s", func_name)
	version := 0
	versionStr := strconv.Itoa(version)
	namespace := func_name + "-" + versionStr
	err := MakeDeployment(namespace, func_name, false, false)
	if err != nil {
		log.Printf("Falied to install deployment: %v\n", err)
	}
	err = MakeDeployment(namespace, func_name, true, false)
	if err != nil {
		log.Printf("Falied to install ingress: %v\n", err)
	}
	version_function[func_name] = version
	log.Printf("Deployed initial version of %s", func_name)
	return "Deployed"

}
func getLogs(func_name string) string {
	log.Printf("getting logs of %s\n", func_name)
	logs_arr := make(map[string]string)
	version := version_function[func_name]
	versionStr := strconv.Itoa(version)
	namespace := func_name + "-" + versionStr
	pods, _ := clientset.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{})
	for _, pod := range pods.Items {
		for _, container_name := range pod.Spec.Containers {
			// log.Print(container_name.Name)
			podLogOpts := corev1.PodLogOptions{
				Container: container_name.Name,
			}
			req := clientset.CoreV1().Pods(namespace).GetLogs(pod.Name, &podLogOpts)
			log.Print("req success")
			logs, err := req.Stream(context.TODO())
			if err != nil {
				log.Print("error in opening stream")
			}
			defer logs.Close()
			b := new(bytes.Buffer)
			io.Copy(b, logs)
			// log.Print("\n-----------------------------------------------------\n")
			// log.Print(pod.Name)
			// log.Print(b.String())
			// log.Print("\n-----------------------------------------------------\n")
			logs_arr[container_name.Name] = b.String()
		}
	}
	result := ""
	for container_name, logs := range logs_arr {
		result += container_name + "\n" + logs + "\n"
	}
	log.Print("\n-----------------------------------------------------\n")
	log.Print(result)
	log.Print("\n-----------------------------------------------------\n")
	return result

}

func (s *server) Deployment(ctx context.Context, req *pb.DeploymentServiceRequest) (*pb.DeploymentServiceReply, error) {
	func_name := req.Name
	var result string
	if req.FunctionCall == "logs" {
		result = getLogs(func_name)
	}
	if req.FunctionCall == "new_invoke" {
		result = makeNewFunctionHandler(func_name)
	}
	if req.FunctionCall == "existing_invoke" {
		result = metricEvalHandler(func_name)
	}
	return &pb.DeploymentServiceReply{
		Message: fmt.Sprintf("%s", result),
	}, nil
}

func main() {
	go checkTTL()
	go watchConfigMaps()
	//http
	// http.HandleFunc("/metric_eval", metricEvalHandler) //after first invokation
	// http.HandleFunc("/make_new_function", makeNewFunctionHandler)
	// http.HandleFunc("/get_logs", getLogs)
	// log.Print("Server listening on :8080...")
	//grpc
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	s := grpc.NewServer()
	pb.RegisterDeploymentServiceServer(s, &server{})
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
