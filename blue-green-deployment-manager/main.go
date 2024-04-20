package main

import (
	"context"
	"encoding/json"
	"fmt"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
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

var (
	clientset *kubernetes.Clientset
	err       error
)

func init() {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

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
	result := clientset.RESTClient().
		Get().
		Namespace(namespace).
		AbsPath("apis/metrics.k8s.io/v1beta1/pods").
		Do(context.TODO())
	data, err := result.Raw()
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &pods)
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

func InstallChart(dir, release, namespace string, configMapName string) error {
	log.Printf("Installing releaseName %s in namespace %s", release, namespace)
	chart, err := loader.Load(dir)
	if err != nil {
		return err
	}
	configDataArray, err := getConfigMapData(namespace, configMapName)
	if err != nil {
		return err
	}

	os.Setenv("HELM_DRIVER", "secrets")
	for _, configData := range configDataArray {
		log.Print(configData)
		actionConfig := new(action.Configuration)
		if err := actionConfig.Init(
			&genericclioptions.ConfigFlags{
				Namespace: &namespace,
			},
			namespace,
			os.Getenv("HELM_DRIVER"),
			log.Printf,
		); err != nil {
			return err
		}
		client := action.NewInstall(actionConfig)
		client.ReleaseName = release + configData["name"].(string)
		client.Namespace = namespace

		if _, err := client.Run(chart, configData); err != nil {
			return err
		}
	}

	return nil
}

func DeleteChart(releaseName, namespace string, configMapName string) error {
	actionConfig := new(action.Configuration)
	log.Printf("Deleting releaseName %s in namespace %s", releaseName, namespace)
	if err := actionConfig.Init(
		&genericclioptions.ConfigFlags{
			Namespace: &namespace,
		},
		namespace,
		os.Getenv("HELM_DRIVER"),
		log.Printf,
	); err != nil {
		return err
	}
	configDataArray, err := getConfigMapData(namespace, configMapName)
	if err != nil {
		return err
	}
	for _, configData := range configDataArray {
		log.Print(configData)
		actionConfig := new(action.Configuration)
		if err := actionConfig.Init(
			&genericclioptions.ConfigFlags{
				Namespace: &namespace,
			},
			namespace,
			os.Getenv("HELM_DRIVER"),
			log.Printf,
		); err != nil {
			return err
		}
		client := action.NewUninstall(actionConfig)
		releaseName = releaseName + configData["name"].(string)
		_, err := client.Run(releaseName)
		if err != nil {
			return err
		}

		log.Print("Successfully deleted release: %s\n", releaseName)
	}

	return nil
}

// this will evaluate the metrics and make deciions
var version_function = make(map[string]int)

func metricEvalHandler(w http.ResponseWriter, r *http.Request) {
	
	query := r.URL.Query()
	func_name := query.Get("function")
	log.Printf("Evaluation metrics of %s", func_name)
	version := version_function[func_name]
	versionStr := strconv.Itoa(version)
	namespace := func_name + versionStr
	cpu_threshold1, cpu_threshold2, mm_threshold1, mm_threshold2, _ := getConfigMapThreshold(namespace, func_name)
	var cpu_usage, mm_usage float64
	for {
		cpu_usage, mm_usage, err = getPodMetricsAndChanges(namespace)
		if err == nil {
			break
		}
	}

	if cpu_usage > cpu_threshold1 || mm_usage > mm_threshold1{
		version_upate := version +1
		log.Printf("Threshold 1 of %s reached........deploying version %d", func_name, version_upate)
		version_updateStr := strconv.Itoa(version_upate)
		namespace_update := func_name +version_updateStr
		releaseNameDeploy := func_name +version_updateStr+"_deployment"
		InstallChart("./deployment-charts", releaseNameDeploy, namespace_update, func_name)
	} 

	if cpu_usage > cpu_threshold2 || mm_usage > mm_threshold2{
		version_upate := version +1
		version_updateStr := strconv.Itoa(version_upate)
		log.Printf("Threshold 2 of %s reached........shifting requests to %d", func_name, version_upate)
		namespace_update := func_name +version_updateStr
		releaseNameIngress := namespace_update+"_ingress"
		InstallChart("./ingress-charts", releaseNameIngress, namespace_update, func_name)
		releaseNameDeploy := namespace+"_deployment"
		releaseNameIngress = namespace+"_ingress"
		// delete the older versions 
		DeleteChart(releaseNameDeploy, namespace_update, func_name)
		DeleteChart(releaseNameIngress, namespace_update, func_name)
		version_function[func_name] = version_upate
	} 

}

// this will look into the function name and check if version 0 is available - if yes, move ahead and deploy - this will also keep track of the versions and functions that are being handled
func makeNewFunctionHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	func_name := query.Get("function")
	log.Printf("Deploying new function.....%s", func_name)
	version := 0
	versionStr := strconv.Itoa(version)
	namespace := func_name + versionStr
	releaseNameDeploy := namespace+"_deployment"
	InstallChart("./deployment-charts", releaseNameDeploy, namespace, func_name)
	releaseNameIngress := namespace+"_ingress"
	InstallChart("./ingress-charts", releaseNameIngress, namespace, func_name)
	version_function[func_name] = version
	log.Printf("Deployed initial version of %s", func_name)

}

func main() {
	http.HandleFunc("/metric_eval", metricEvalHandler)
	http.HandleFunc("/make_new_function", makeNewFunctionHandler)
	log.Print("Server listening on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Printf("Failed to start server: %v\n", err)
	}
}
