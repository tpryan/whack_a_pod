package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"k8s.io/api/apps/v1beta1"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var (
	k         *kubernetes.Clientset
	namespace = "default"
)

const defaultGracePeriod = 3

func main() {
	log.Printf("starting whack a pod admin api")
	http.ListenAndServe(":8080", handler())
}

func handler() http.Handler {

	var err error
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Printf("could not initialize InClustterConfig %v", err)
	}
	// creates the clientset
	k, err = kubernetes.NewForConfig(config)
	if err != nil {
		log.Printf("could not initialize NewForConfig %v", err)
	}

	if err != nil {
		log.Printf("could not initialize connection to Kubernetes cluster %v", err)
	}

	r := http.NewServeMux()
	r.HandleFunc("/", health)
	r.HandleFunc("/healthz", health)
	r.HandleFunc("/admin/healthz", health)
	r.HandleFunc("/admin/k8s/getpods", handleAPI(handlePods))
	r.HandleFunc("/admin/k8s/getnodes", handleAPI(handleNodes))
	r.HandleFunc("/admin/k8s/deletepod", handleAPI(handlePodDelete))
	r.HandleFunc("/admin/k8s/deleteallpods", handleAPI(handlePodsDelete))
	r.HandleFunc("/admin/k8s/drain", handleAPI(handleNodeDrain))
	r.HandleFunc("/admin/k8s/uncordon", handleAPI(handleNodeUncordon))
	r.HandleFunc("/admin/k8s/deletedeploy", handleAPI(deleteDeployment))
	r.HandleFunc("/admin/k8s/createdeploy", handleAPI(createDeployment))

	return r
}

func health(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "ok")
}

type apiHandler func(http.ResponseWriter, *http.Request) ([]byte, error)

func handleAPI(h apiHandler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	})
}

func (h apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b, err := h(w, r)
	if err != nil {
		sendJSON(w, fmt.Sprintf("{error:\"%s\"}", err), http.StatusInternalServerError)
	}
	sendJSON(w, string(b), http.StatusOK)
	log.Printf("Request %s", r.URL)
}

func handlePods(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	pods, err := listAllPods()
	if err != nil {
		return nil, fmt.Errorf("could not retrieve k8s api: %v", err)
	}

	return sendBytes(pods)
}

func listAllPods() (*v1.PodList, error) {
	l := metav1.ListOptions{LabelSelector: "app=api"}
	pods, err := k.CoreV1().Pods(namespace).List(l)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve k8s pods: %v", err)
	}

	return pods, nil
}

func handlePodDelete(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	in := strings.Split(r.FormValue("pod"), "/")
	podname := in[len(in)-1]

	pod, err := k.CoreV1().Pods(namespace).Get(podname, metav1.GetOptions{})

	if err != nil {
		return nil, fmt.Errorf("could not get k8s object: %v", err)
	}

	d := metav1.DeleteOptions{GracePeriodSeconds: newInt64(defaultGracePeriod)}
	if err := k.CoreV1().Pods(namespace).Delete(podname, &d); err != nil {
		return nil, fmt.Errorf("could not delete k8s object: %v", err)
	}

	return sendBytes(pod)
}

func handlePodsDelete(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	d := metav1.DeleteOptions{GracePeriodSeconds: newInt64(defaultGracePeriod)}
	l := metav1.ListOptions{LabelSelector: "app=api"}

	err := k.CoreV1().Pods(namespace).DeleteCollection(&d, l)
	if err != nil {
		return nil, fmt.Errorf("could not delete k8s pods: %v", err)
	}

	pods, err := k.CoreV1().Pods(namespace).List(l)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve k8s pods: %v", err)
	}

	return sendBytes(pods)
}

func newInt32(x int) *int32 {
	r := int32(x)
	return &r
}

func newInt64(x int) *int64 {
	r := int64(x)
	return &r
}

func imageFromEnv() (string, error) {
	i := os.Getenv("APIIMAGE")
	if len(i) == 0 {
		return "", fmt.Errorf("env var APIIMAGE not set")
	}
	return i, nil
}

func createDeployment(w http.ResponseWriter, r *http.Request) ([]byte, error) {

	name := "api-deployment"
	image, err := imageFromEnv()
	if err != nil {
		return nil, fmt.Errorf("could not get the name of the container image from env : %v", err)
	}

	var d v1beta1.Deployment
	d.Name = name
	d.Spec.Replicas = newInt32(12)
	d.Spec.Strategy.Type = v1beta1.RollingUpdateDeploymentStrategyType
	d.Spec.Template.ObjectMeta.SetLabels(map[string]string{"app": "api"})
	d.Spec.Template.Spec.Containers = []v1.Container{
		v1.Container{
			Name:            "api",
			Image:           image,
			ImagePullPolicy: "Always",
			Ports: []v1.ContainerPort{
				v1.ContainerPort{
					ContainerPort: 8080,
					Name:          "http",
					Protocol:      "TCP",
				},
			},
		},
	}

	dep, err := k.AppsV1beta1().Deployments(namespace).Create(&d)
	if err != nil {
		return nil, fmt.Errorf("could not create k8s deployment: %v", err)
	}

	return sendBytes(dep)
}

func deleteDeployment(w http.ResponseWriter, r *http.Request) ([]byte, error) {

	name := "api-deployment"

	d := metav1.DeleteOptions{GracePeriodSeconds: newInt64(defaultGracePeriod)}

	if err := k.AppsV1beta1().Deployments(namespace).Delete(name, &d); err != nil {
		return nil, fmt.Errorf("could not delete k8s deployment: %v", err)
	}

	rs, err := k.ExtensionsV1beta1().ReplicaSets(namespace).List(metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("could not locate any k8s replica set: %v", err)
	}

	for _, set := range rs.Items {
		if strings.Index(set.Name, name) < 0 {
			continue
		}
		if err := k.ExtensionsV1beta1().ReplicaSets(namespace).Delete(set.Name, &d); err != nil {
			return nil, fmt.Errorf("could not delete k8s replica set: %v", err)
		}
	}

	return handlePodsDelete(w, r)
}

func handleNodes(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	nodes, err := k.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("could not get list of k8s nodes: %v", err)
	}

	return sendBytes(nodes)
}

func handleNodeDrain(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	nodename := r.FormValue("node")

	node, err := setNodeInactive(nodename, true)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve k8s node info: %v", err)
	}

	pods, err := listAllPods()
	if err != nil {
		return nil, fmt.Errorf("could not retrieve k8s pods: %v", err)
	}

	for _, pod := range pods.Items {
		if pod.Spec.NodeName != nodename {
			continue
		}
		if err := k.CoreV1().Pods(namespace).Delete(pod.Name, &metav1.DeleteOptions{}); err != nil {
			return nil, fmt.Errorf("could not delete k8s object: %v", err)
		}
	}

	return sendBytes(node)
}

func handleNodeUncordon(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	nodename := r.FormValue("node")

	node, err := setNodeInactive(nodename, false)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve k8s node info: %v", err)
	}

	return sendBytes(node)
}

func setNodeInactive(nodename string, inactive bool) (*v1.Node, error) {
	node, err := k.CoreV1().Nodes().Get(nodename, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("could not retrieve k8s node info: %v", err)
	}
	node.Spec.Unschedulable = inactive
	node, err = k.CoreV1().Nodes().Update(node)
	if err != nil {
		return nil, fmt.Errorf("could not change schedulable attributes for k8s node %s: %v", nodename, err)
	}
	return node, nil
}

func sendJSON(w http.ResponseWriter, content string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	fmt.Fprint(w, content)
}

func sendBytes(v interface{}) ([]byte, error) {
	result, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshall data: %v", err)
	}

	return result, nil
}
