package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"os"
	networkingv1 "k8s.io/api/networking/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	hostTargets = make(map[string]string)
)

type baseHandle struct{}

func (h *baseHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// host := r.Host
	// hostWithoutPort, _, err := net.SplitHostPort(host)
    // if err != nil {
    //     hostWithoutPort = host
    // }

	pathParts := strings.Split(r.URL.Path, "/")
	func_name := pathParts[1]

    target, ok := hostTargets[func_name]
    if !ok {
        http.Error(w, "Target not found making the function workflow", http.StatusNotFound)
		callDepController(false, func_name)

	return
   }
	log.Print("Forwarding request to",target)
	remoteUrl, err := url.Parse(target)
	if err != nil {
		log.Fatal("target parse fail:", err)
	}
	callDepController(true, func_name)
	proxy := httputil.NewSingleHostReverseProxy(remoteUrl)
	proxy.ServeHTTP(w, r)
}

func callDepController(func_found bool, func_name string) {
    depControllerAddr := os.Getenv("DEP_CONTROLLER_ADD")
    if depControllerAddr == "" {
        fmt.Println("DEP_CONTROLLER_ADD environment variable not set")
        return
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
        return
    }
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
		Addr:    ":8080",
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
                updateHostTargets(ingress)
            case watch.Deleted:
                deleteHostTargets(ingress)
            }
            log.Printf("Updated host targets: %+v", hostTargets)
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
				//host := rule.Host
				func_name:= ingress.Labels["function_name"]
				hostname := fmt.Sprintf("http://%s.%s.svc.cluster.local:%d", serviceName, namespace, port)
				hostTargets[func_name] = hostname
			}
		}
	}
}

func deleteHostTargets(ingress *networkingv1.Ingress) {
	for _, rule := range ingress.Spec.Rules {
		delete(hostTargets, rule.Host)
	}
}
