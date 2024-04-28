package main

import (
	"context"
	"encoding/json"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

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

type FunctionTimestamp struct {
	data map[string]time.Time
	sync.Mutex
}

var (
	clientset          *kubernetes.Clientset
	err                error
	namesapce_ingress  string
	function_timestamp FunctionTimestamp
	ttl_seconds        int // time in seconds
	version_function   map[string]int
)

func init() {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	version_function = make(map[string]int)
	function_timestamp = FunctionTimestamp{data: make(map[string]time.Time)}
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
	} else {
		return 0, fmt.Errorf("unsupported CPU usage format")
	}
}

// continously check TTL if TTL is above certain seconds delete the namespace
func checkTTL() {
	for {
		//log.Print(function_timestamp.data)
		function_timestamp.Lock()
		currentTime := time.Now()
		for name, timestamp := range function_timestamp.data {
			elapsedTime := currentTime.Sub(timestamp)
			if elapsedTime.Seconds() > float64(ttl_seconds) {
				version := version_function[name]
				versionStr := strconv.Itoa(version)
				namespace := name + "-" + versionStr
				log.Print("Deleting because of TTL")
				DeleteNamespace(namespace)
				delete(function_timestamp.data, name)
			}
		}
		function_timestamp.Unlock()
		time.Sleep(time.Second)
	}
}

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
func getMetrics(clientset *kubernetes.Clientset, pods *PodMetricsList, namespace string) error {
	// result := clientset.RESTClient().
	// 	Get().
	// 	Namespace(namespace).
	// 	AbsPath("apis/metrics.k8s.io/v1beta1/pods").
	// 	Do(context.TODO())
	// data, err := result.Raw()
	// if err != nil {
	// 	return err
	// }
	// err = json.Unmarshal(data, &pods)
	// log.Print(pods)
	// return err

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

// func getMetricsNode(clientset *kubernetes.Clientset, nodes *NodeMetricList) error {
// 	result := clientset.RESTClient().Get().AbsPath("apis/metrics.k8s.io/v1beta1/nodes").Do(context.TODO())
// 	data, err := result.Raw()
// 	if err != nil {
// 		print("error")
// 		return err
// 	}
// 	err = json.Unmarshal(data, &nodes)
// 	return err
// }

func getConfigMapData(namespace, configMapName string) ([]map[string]interface{}, error) {
	cm, err := clientset.CoreV1().ConfigMaps(namespace).Get(context.TODO(), configMapName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	var configData []map[string]interface{}
	if err := json.Unmarshal([]byte(cm.Data["my-config.json"]), &configData); err != nil {
		return nil, err
	}
	log.Print(configData)
	return configData, nil
}

func getConfigMapThreshold(namespace, configMapName string) (float64, float64, float64, float64, error) {
	log.Printf("Fetching  metrics for %s from configMap %s", namespace, configMapName)
	config, err := rest.InClusterConfig()
	if err != nil {
		return 0.0, 0.0, 0.0, 0.0, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return 0.0, 0.0, 0.0, 0.0, err
	}

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

//check for changes in config map configurations provided by uiser -> if the user update sthe existing depkoyment it needs to be changed toooo

func watchConfigMaps() {
	watcher, err := clientset.CoreV1().ConfigMaps(namesapce_ingress).Watch(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Failed to set up watch: %s", err)
	}

	go func() {
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
	}()
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
func MakeDeployment(namespace string, func_name string, ingress bool, update bool) error {
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
		envVariables := configMapData["env"].(map[string]interface{})
		// log.Print(envVariables)
		imageData := configMapData["image"].(map[string]interface{})
		// log.Print(imageData)
		imageName := imageData["image"].(string)
		// log.Print(imageName)
		imagePullPolicy := corev1.PullPolicy(imageData["pullPolicy"].(string))
		// log.Print(imagePullPolicy)
		var env []corev1.EnvVar
		for key, value := range envVariables {
			env = append(env, corev1.EnvVar{
				Name:  key,
				Value: value.(string),
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
								}, {

									Name:  "tcp-dump",
									Image: "docker.io/dockersec/tcpdump",
								},
							},
						},
					},
				},
			}

			// log.Print(deployment)
			if update == false{
				log.Printf("Creating deployment and service %s in %s", name, namespace)
			_, err = clientset.AppsV1().Deployments(namespace).Create(context.Background(), deployment, metav1.CreateOptions{})
			if err != nil {
				panic(err.Error())
			}
		}else{
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
							Port:       servicePort,
							TargetPort: intstr.FromInt(int(containerPort)),
						},
					},
				},
			}
			if update == false{
			_, err = clientset.CoreV1().Services(namespace).Create(context.Background(), service, metav1.CreateOptions{})
			if err != nil {
				panic(err.Error())
			}
		}else{
			_, err = clientset.CoreV1().Services(namespace).Update(context.Background(), service, metav1.UpdateOptions{})
			if err != nil {
				panic(err.Error())
			}
		}

		}
	}

	return nil

}

func DeleteNamespace(namespace string) error {
	log.Print("Deleteing namespace %s\n", namespace)
	err = clientset.CoreV1().Namespaces().Delete(context.TODO(), namespace, metav1.DeleteOptions{})
	return err
}

func metricEvalHandler(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()
	func_name := query.Get("function")
	function_timestamp.Lock()
	defer function_timestamp.Unlock()
	function_timestamp.data[func_name] = time.Now()
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

// this will look into the function name and check if version 0 is available - if yes, move ahead and deploy - this will also keep track of the versions and functions that are being handled
func makeNewFunctionHandler(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()
	func_name := query.Get("function")
	function_timestamp.Lock()
	defer function_timestamp.Unlock()
	function_timestamp.data[func_name] = time.Now()
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

}

func main() {
	go checkTTL()
	go watchConfigMaps()
	http.HandleFunc("/metric_eval", metricEvalHandler)            //after first invokation
	http.HandleFunc("/make_new_function", makeNewFunctionHandler) // this will be triggered when no version of function is deployed
	log.Print("Server listening on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Printf("Failed to start server: %v\n", err)
	}
}
