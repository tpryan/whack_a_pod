package main

import (
	"bytes"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"k8s.io/api/apps/v1beta1"
	"k8s.io/api/core/v1"
)

var (
	client              *http.Client
	pool                *x509.CertPool
	token               = ""
	errItemNotExist     = fmt.Errorf("Item does not exist")
	errItemAlreadyExist = fmt.Errorf("Item already exists")
)

const (
	tokenFile = "/var/run/secrets/kubernetes.io/serviceaccount/token"
	certFile  = "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"
	root      = "https://kubernetes"
	namespace = "default"
	selector  = "app=api"
)

func tokenFromDisk() (string, error) {
	b, err := ioutil.ReadFile(tokenFile)
	if err != nil {
		return "", fmt.Errorf("could not retrieve kubernetes token from disk: %v", err)
	}
	return string(b), nil
}

func certsFromDisk() ([]byte, error) {
	b, err := ioutil.ReadFile(certFile)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve kubernetes token from disk: %v", err)
	}
	return b, nil
}

func queryK8sAPI(url, method string, data []byte) ([]byte, int, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("could not create HTTP request: %v", err)
	}
	req.Header.Add("Authorization", "Bearer "+token)

	if method == http.MethodPost {
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Content-Length", strconv.Itoa(len(string(data))))
	}

	if method == http.MethodPatch {
		req.Header.Set("Content-Type", "application/merge-patch+json")
		req.Header.Set("Content-Length", strconv.Itoa(len(string(data))))
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("could not execute HTTP request: %v", err)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("could not read HTTP request: %v", err)
	}
	return b, resp.StatusCode, nil
}

func listAllPods() (*v1.PodList, error) {
	url := root + "/api/v1/pods?labelSelector=" + selector

	b, _, err := queryK8sAPI(url, "GET", nil)
	if err != nil {
		return nil, fmt.Errorf("I can't even with the HTTP: %v", err)
	}

	var pods v1.PodList
	if err := json.Unmarshal(b, &pods); err != nil {
		return nil, fmt.Errorf("could not convert HTTP request to Pod data structure: %v", err)
	}

	return &pods, nil
}

func listAllNodes() (*v1.NodeList, error) {
	url := root + "/api/v1/nodes"

	b, _, err := queryK8sAPI(url, "GET", nil)
	if err != nil {
		return nil, fmt.Errorf("I can't even with the HTTP: %v", err)
	}

	var node v1.NodeList
	if err := json.Unmarshal(b, &node); err != nil {
		return nil, fmt.Errorf("could not convert HTTP request to Pod data structure: %v", err)
	}

	return &node, nil
}

func pod(podname string) (*v1.Pod, error) {
	url := root + podname

	b, _, err := queryK8sAPI(url, "GET", nil)
	if err != nil {
		return nil, fmt.Errorf("I can't even with the HTTP: %v", err)
	}

	var pod v1.Pod
	if err := json.Unmarshal(b, &pod); err != nil {
		return nil, fmt.Errorf("could not convert HTTP request to Pod data structure: %v", err)
	}

	return &pod, nil

}

func deletePod(podname string) (*v1.Pod, error) {
	url := root + podname

	b, status, err := queryK8sAPI(url, "DELETE", nil)
	if err != nil {
		return nil, fmt.Errorf("I can't even with the HTTP: %v", err)
	}

	if status == http.StatusNotFound {
		return nil, errItemNotExist
	}

	var pod v1.Pod
	if err := json.Unmarshal(b, &pod); err != nil {
		return nil, fmt.Errorf("could not convert HTTP request to Pod data structure: %v", err)
	}

	return &pod, nil

}

func deletePods(pod, node string) (*v1.PodList, error) {
	url := root
	if len(pod) > 0 {
		url += pod
	} else {
		url += "/api/v1/namespaces/" + namespace + "/pods" + "?labelSelector=" + selector
	}
	if len(node) > 0 {
		fs := "&fieldSelector=spec.nodeName=" + node
		url += fs
	}

	b, status, err := queryK8sAPI(url, "DELETE", nil)
	if err != nil {
		return nil, fmt.Errorf("I can't even with the HTTP: %v", err)
	}

	if status == http.StatusNotFound {

		return nil, errItemNotExist
	}

	var podlist v1.PodList
	if err := json.Unmarshal(b, &podlist); err != nil {
		return nil, fmt.Errorf("could not convert HTTP request to Pod data structure: %v", err)
	}

	return &podlist, nil

}

func toggleNode(nodename string, inactive bool) (*v1.Node, error) {
	url := root + "/api/v1/nodes/" + nodename

	j := fmt.Sprintf("{\"spec\": {\"unschedulable\": %t}}", inactive)
	b, status, err := queryK8sAPI(url, "PATCH", []byte(j))
	if err != nil {
		return nil, fmt.Errorf("I can't even with the HTTP: %v", err)
	}

	if status == http.StatusNotFound {
		return nil, errItemNotExist
	}

	var node v1.Node
	if err := json.Unmarshal(b, &node); err != nil {
		return nil, fmt.Errorf("could not convert HTTP request to deployment data structure: %v", err)
	}

	return &node, nil
}

func deleteAllReplicaSets() (*v1.Pod, error) {
	url := root + "/apis/extensions/v1beta1/namespaces/" + namespace + "/replicasets" + "?labelSelector=" + selector

	b, status, err := queryK8sAPI(url, "DELETE", nil)
	if err != nil {
		return nil, fmt.Errorf("I can't even with the HTTP: %v", err)
	}

	if status == http.StatusNotFound {

		return nil, errItemNotExist
	}

	var pod v1.Pod
	if err := json.Unmarshal(b, &pod); err != nil {
		return nil, fmt.Errorf("could not convert HTTP request to Pod data structure: %v", err)
	}

	return &pod, nil

}

func deleteDeployment() (*v1beta1.Deployment, error) {
	selflink := "/apis/extensions/v1beta1/namespaces/" + namespace + "/deployments/api-deployment"
	url := root + selflink

	b, status, err := queryK8sAPI(url, "DELETE", nil)
	if err != nil {
		return nil, fmt.Errorf("I can't even with the HTTP: %v", err)
	}

	if status == http.StatusNotFound {
		return nil, errItemNotExist
	}
	log.Printf("response %s", string(b))
	var d v1beta1.Deployment
	if err := json.Unmarshal(b, &d); err != nil {
		return nil, fmt.Errorf("could not convert HTTP request to deployment data structure: %v", err)
	}

	return &d, nil

}

type minimumDeployment struct {
	APIVersion string `json:"apiVersion,omitempty"`
	Kind       string `json:"kind,omitempty"`
	Metadata   struct {
		Name string `json:"name,omitempty"`
	} `json:"metadata,omitempty"`
	Spec struct {
		Replicas int `json:"replicas,omitempty"`
		Selector struct {
			MatchLabels map[string]string `json:"matchLabels,omitempty"`
		} `json:"selector,omitempty"`
		Strategy struct {
			Type string `json:"type,omitempty"`
		} `json:"strategy,omitempty"`
		Template struct {
			Metadata struct {
				Labels map[string]string `json:"labels,omitempty"`
			} `json:"metadata,omitempty"`
			Spec struct {
				Containers []minimumContainer `json:"containers,omitempty"`
			} `json:"spec,omitempty"`
		} `json:"template,omitempty"`
	} `json:"spec,omitempty"`
}

type minimumContainer struct {
	Image           string        `json:"image,omitempty"`
	ImagePullPolicy string        `json:"imagePullPolicy,omitempty"`
	Name            string        `json:"name,omitempty"`
	Ports           []minimumPort `json:"ports,omitempty"`
}

type minimumPort struct {
	ContainerPort int    `json:"containerPort,omitempty"`
	Name          string `json:"name,omitempty"`
	Protocol      string `json:"protocol,omitempty"`
}

func imageFromEnv() (string, error) {
	i := os.Getenv("APIIMAGE")
	if len(i) == 0 {
		return "", fmt.Errorf("env var APIIMAGE not set")
	}
	return i, nil
}

func createDeployment() (*v1beta1.Deployment, error) {
	selflink := "/apis/extensions/v1beta1/namespaces/" + namespace + "/deployments/"
	url := root + selflink

	image, err := imageFromEnv()
	if err != nil {
		return nil, fmt.Errorf("could not get the name of the container image from env : %v", err)
	}

	var d minimumDeployment
	d.APIVersion = "extensions/v1beta1"
	d.Kind = "Deployment"
	d.Metadata.Name = "api-deployment"
	d.Spec.Replicas = 12
	d.Spec.Selector.MatchLabels = map[string]string{"app": "api", "visualize": "true"}
	d.Spec.Strategy.Type = "RollingUpdate"
	d.Spec.Template.Metadata.Labels = map[string]string{"app": "api", "visualize": "true"}
	d.Spec.Template.Spec.Containers = []minimumContainer{
		minimumContainer{
			Name:            "api",
			Image:           image,
			ImagePullPolicy: "Always",
			Ports: []minimumPort{
				minimumPort{
					ContainerPort: 8080,
					Name:          "http",
					Protocol:      "TCP",
				},
			},
		},
	}

	dbytes, err := json.Marshal(d)
	if err != nil {
		return nil, fmt.Errorf("could not convert deployment to json: %v", err)
	}

	b, status, err := queryK8sAPI(url, "POST", dbytes)
	if err != nil {
		return nil, fmt.Errorf("I can't even with the HTTP: %v", err)
	}
	log.Printf("response: %d %s", status, string(b))

	if status == http.StatusNotFound {
		return nil, errItemNotExist
	}

	if status == http.StatusConflict {
		return nil, errItemAlreadyExist
	}
	var dep v1beta1.Deployment
	if err := json.Unmarshal(b, &dep); err != nil {
		return nil, fmt.Errorf("could not convert HTTP request to deployment data structure: %v", err)
	}

	return &dep, nil

}
