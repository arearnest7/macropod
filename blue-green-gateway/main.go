package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var (
	hostTargets                 = make(map[string]string)
	runningDeploymentController = make(map[string]bool) // this can be used for lock mechanism
)

// TODO - if deploymnet controller for a fucntion is already running dont run it again

type baseHandle struct{}

func (h *baseHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// host := r.Host
	// hostWithoutPort, _, err := net.SplitHostPort(host)
	// if err != nil {
	//     hostWithoutPort = host
	// }

	pathParts := strings.Split(r.URL.Path, "/")
	func_name := pathParts[1]
	r.URL.Path = ""
	target, ok := hostTargets[func_name]
	if !ok {
		http.Error(w, "Target not found making the function workflow", http.StatusNotFound)
		runningDeploymentController[func_name] = false
		callDepController(false, func_name)

		return
	}
	remoteUrl, err := url.Parse(target)
	if err != nil {
		log.Fatal("target parse fail:", err)
	}
	go callDepController(true, func_name)
	log.Print("Forwarding request to", remoteUrl)
	proxy := httputil.NewSingleHostReverseProxy(remoteUrl)
	proxy.ServeHTTP(w, r)
}

func callDepController(func_found bool, func_name string) error {
	if runningDeploymentController[func_name] == false {
		runningDeploymentController[func_name] = true
		depControllerAddr := os.Getenv("DEP_CONTROLLER_ADD")
		if depControllerAddr == "" {
			fmt.Println("DEP_CONTROLLER_ADD environment variable not set")
			runningDeploymentController[func_name] = false
			return nil
		}
		depControllerAddr = "http://" + depControllerAddr //schema support
		var path string
		if func_found {
			path = "/metric_eval"
		} else {
			path = "/make_new_function"
		}
		targetURL := depControllerAddr + path + "?function=" + func_name
		_, err := http.Get(targetURL)
		if err != nil {
			fmt.Printf("Failed to make HTTP request: %v\n", err)
			runningDeploymentController[func_name] = false
			return err
		}
		runningDeploymentController[func_name] = false
	} else {

		log.Printf("Deployment Controller for function this %s is already in progress", func_name)
	}
	return nil
}

func main() {
	log.Print("Ingress controller started")

	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("Failed to get in-cluster config: %s", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed to create clientset: %s", err)
	}

	watchIngress(clientset)
	h := &baseHandle{}
	http.Handle("/", h)
	server := &http.Server{
		Addr:    ":8080", // reverse proxy
		Handler: h,
	}
	log.Fatal(server.ListenAndServe())
}

func watchIngress(clientset *kubernetes.Clientset) {
	watcher, err := clientset.NetworkingV1().Ingresses("").Watch(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Failed to set up watch: %s", err)
	}

	go func() {
		for event := range watcher.ResultChan() {
			ingress, ok := event.Object.(*networkingv1.Ingress)
			if !ok {
				log.Printf("Expected Ingress type, got %T", event.Object)
				continue
			}

			switch event.Type {
			case watch.Added, watch.Modified:
				log.Printf("Updated host targets: %+v", hostTargets)
				updateHostTargets(ingress)
			case watch.Deleted:
				deleteHostTargets(ingress)
			}

		}
	}()
}

func updateHostTargets(ingress *networkingv1.Ingress) {
	for _, rule := range ingress.Spec.Rules {
		if rule.HTTP != nil {
			for _, path := range rule.HTTP.Paths {
				serviceName := path.Backend.Service.Name
				namespace := ingress.Namespace
				port := path.Backend.Service.Port.Number
				func_name := ingress.Labels["function_name"]
				hostname := fmt.Sprintf("http://%s.%s.svc.cluster.local:%d", serviceName, namespace, port)
				hostTargets[func_name] = hostname
			}
		}
	}
}

func deleteHostTargets(ingress *networkingv1.Ingress) {
	func_name := ""
	hostname_deleted := ""
	for _, rule := range ingress.Spec.Rules {
		if rule.HTTP != nil {
			for _, path := range rule.HTTP.Paths {
				serviceName := path.Backend.Service.Name
				namespace := ingress.Namespace
				port := path.Backend.Service.Port.Number
				func_name = ingress.Labels["function_name"]
				hostname_deleted = fmt.Sprintf("http://%s.%s.svc.cluster.local:%d", serviceName, namespace, port)
			}
		}
	}

	// delete entry only when deleted ingress is same as present ingress
	if hostTargets[func_name] == hostname_deleted {
		log.Printf("Deleting ingress for %s", func_name)
		delete(hostTargets, func_name)
	}

}
