package main

import (
	pb "app/prepuller_pb"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"time"
	"google.golang.org/grpc"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var (
	kclient *kubernetes.Clientset
)

type PodPhase string

// These are the valid statuses of pods.
// const (
// 	PodPending PodPhase = "Pending"

// 	PodRunning PodPhase = "Running"

// 	PodSucceeded PodPhase = "Succeeded"
// 	PodFailed    PodPhase = "Failed"

// 	PodUnknown PodPhase = "Unknown"
// )

// type PodStatus struct {
// 	Phase PodPhase `json:"phase,omitempty" protobuf:"bytes,1,opt,name=phase,casttype=PodPhase"`
// 	Data json.RawMessage `json:"data"`
// }

// type Pod struct {
// 	Status PodStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`

// }
var (

	node_namespace map[string]string
)
type server struct {
	pb.PrepullerServiceServer
}

type Function struct {
	Registry  string            `json:"registry"`
	Endpoints []string          `json:"endpoints,omitempty"`
	Envs      map[string]string `json:"envs,omitempty"`
	Secrets   map[string]string `json:"secrets,omitempty"`
}

type Workflow struct {
	Name      string              `json:"name,omitempty"`
	Functions map[string]Function `json:"functions"`
	Pods      [][]string
}

func internal_log(message string) {
	fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.000000 UTC") + ": " + message)
}

func checkPodsAndDelete(namespace string) {
	internal_log("Checking  " + namespace)
	pods, _ := kclient.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	for _, pod := range pods.Items {
		internal_log(pod.Name)
		for {
			pod_updated, _ := kclient.CoreV1().Pods(namespace).Get(context.TODO(), pod.Name, metav1.GetOptions{})
			pod_phase := string(pod_updated.Status.Phase)
			if pod_updated.Status.Phase == "Running" {
				internal_log("pod " + pod.Name + "is" + pod_phase + "in " + namespace + "is ready")
				break
			}
			time.Sleep(10 * time.Second)
		}
	}
	internal_log("Deleting  " + namespace)
	kclient.CoreV1().Namespaces().Delete(context.TODO(), namespace, metav1.DeleteOptions{})

}
func prepullImage(workflow_str string) {
	workflow := Workflow{}
	json.Unmarshal([]byte(workflow_str), &workflow)

	nodes, _ := kclient.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	node_namespace := make(map[string]string)
	for _, node := range nodes.Items {
		namespace := "prepuller" + strconv.Itoa(rand.Intn(10000000))
		internal_log("Deploying in " + node.Name)
		namespace_object := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: namespace,
			},
		}

		_, err := kclient.CoreV1().Namespaces().Create(context.Background(), namespace_object, metav1.CreateOptions{})
		if err != nil {
			internal_log("namespace " + namespace + " unable to be created - " + err.Error())
			return
		}
		node_namespace[node.Name] = namespace
		for func_name, functions := range workflow.Functions {
			image_to_pull := functions.Registry
			imagePullPolicy := corev1.PullPolicy("Always")
			deployment := &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name: func_name,
				},
				Spec: appsv1.DeploymentSpec{
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							"app": func_name,
						},
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								"app": func_name,
							},
						},
						Spec: corev1.PodSpec{
							NodeSelector: map[string]string{
								"kubernetes.io/hostname": node.Name,
							},
							Containers: []corev1.Container{
								{
									Name:            func_name,
									Image:           image_to_pull,
									ImagePullPolicy: imagePullPolicy,
									Ports: []corev1.ContainerPort{
										{
											ContainerPort: 5000,
										},
									},
									Env: []corev1.EnvVar{
										{
											Name: "SERVICE_TYPE", Value: "GRPC",
										},
										{
											Name: "GRPC_THREAD", Value: "10",
										},
										{
											Name: "FUNC_PORT", Value: "5000",
										},
									},
								},
							},
						},
					},
				},
			}
			_, err := kclient.AppsV1().Deployments(namespace).Create(context.Background(), deployment, metav1.CreateOptions{})
			if err != nil {
				internal_log("unable to create deployment " + func_name + " for " + namespace + " - " + err.Error())

			}
			
			

		}

	}

	for _, namespace := range node_namespace {
		checkPodsAndDelete(namespace)
	}

}

func (s *server) Prepuller(ctx context.Context, req *pb.PrepullerServiceRequest) (*pb.PrepullerServiceReply, error) {
	internal_log("prepullinnngggg")
	prepullImage(*req.Data)
	return &pb.PrepullerServiceReply{
		Message: "pre-pullling over",
	}, nil
}

func main() {
	config, err := rest.InClusterConfig()
	if err != nil {
		internal_log("config - " + err.Error())
		return
	}
	kclient, err = kubernetes.NewForConfig(config)
	if err != nil {
		internal_log("client - " + err.Error())
		return
	}
	internal_log("macropod pre-loader is ready")
	port := ":5003"
	l, err := net.Listen("tcp", port)
	if err != nil {
		internal_log("error listener - " + err.Error())
		return
	}
	s := grpc.NewServer()
	pb.RegisterPrepullerServiceServer(s, &server{})
	internal_log("Registered")
	if err := s.Serve(l); err != nil {
		internal_log("failed to serve - " + err.Error())
		return
	}
}
