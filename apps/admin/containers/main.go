package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

const defaultGracePeriod = 3

func main() {
	log.Printf("starting whack a pod admin api")
	var err error

	token, err = tokenFromDisk()
	if err != nil {
		log.Printf("could not get token from file system")
	}

	certs, err := certsFromDisk()
	if err != nil {
		log.Printf("could not get token from file system")
	}

	pool = x509.NewCertPool()
	pool.AppendCertsFromPEM(certs)
	client = &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{RootCAs: pool}}}

	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Addr:         ":8080",
		Handler:      handler(),
	}

	srv.ListenAndServe()
}

func handler() http.Handler {

	r := http.NewServeMux()
	r.HandleFunc("/", health)
	r.HandleFunc("/healthz", health)
	r.HandleFunc("/admin/healthz", health)
	r.HandleFunc("/healthz/", health)
	r.HandleFunc("/admin/healthz/", health)

	r.HandleFunc("/admin/k8s/pods/get", handleAPI(handlePods))
	r.HandleFunc("/admin/k8s/nodes/get", handleAPI(handleNodes))
	r.HandleFunc("/admin/k8s/pod/delete", handleAPI(handlePodDelete))
	r.HandleFunc("/admin/k8s/pods/delete", handleAPI(handlePodsDelete))
	r.HandleFunc("/admin/k8s/node/drain", handleAPI(handleNodeDrain))
	r.HandleFunc("/admin/k8s/node/uncordon", handleAPI(handleNodeUncordon))
	r.HandleFunc("/admin/k8s/deployment/delete", handleAPI(handleDeploymentDelete))
	r.HandleFunc("/admin/k8s/deployment/create", handleAPI(handleDeploymentCreate))
	r.HandleFunc("/admin/k8s/pods/get/", handleAPI(handlePods))
	r.HandleFunc("/admin/k8s/nodes/get/", handleAPI(handleNodes))
	r.HandleFunc("/admin/k8s/pod/delete/", handleAPI(handlePodDelete))
	r.HandleFunc("/admin/k8s/pods/delete/", handleAPI(handlePodsDelete))
	r.HandleFunc("/admin/k8s/node/drain/", handleAPI(handleNodeDrain))
	r.HandleFunc("/admin/k8s/node/uncordon/", handleAPI(handleNodeUncordon))
	r.HandleFunc("/admin/k8s/deployment/delete/", handleAPI(handleDeploymentDelete))
	r.HandleFunc("/admin/k8s/deployment/create/", handleAPI(handleDeploymentCreate))
	return r
}

func health(w http.ResponseWriter, r *http.Request) {
	r.Close = true
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "ok")
}

type apiHandler func(http.ResponseWriter, *http.Request) ([]byte, error)

func handleAPI(h apiHandler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	})
}

func (h apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.Close = true
	w.Header().Add("Access-Control-Allow-Origin", "*")
	b, err := h(w, r)
	status := http.StatusOK
	if err != nil {
		status = http.StatusInternalServerError
		if err == errItemNotExist {
			status = http.StatusAccepted
		}

		if err == errItemAlreadyExist {
			status = http.StatusAccepted
		}

		sendJSON(w, fmt.Sprintf("{\"error\":\"%v\"}", err), status)
		log.Printf("%s %d %s", r.Method, status, r.URL)
		log.Printf("Error %v", err)
		return
	}
	sendJSON(w, string(b), status)
	log.Printf("%s %d %s", r.Method, status, r.URL)
}

func handlePods(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	pods, err := listAllPods()
	if err != nil {
		return nil, fmt.Errorf("could not retrieve k8s api: %v", err)
	}

	return sendBytes(pods)
}

func handlePodDelete(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	pod, err := deletePod(r.FormValue("pod"))
	if err != nil {
		if err == errItemNotExist {
			return nil, errItemNotExist
		}

		return nil, fmt.Errorf("could not delete k8s object: %v", err)
	}

	return sendBytes(pod)
}

func handlePodsDelete(w http.ResponseWriter, r *http.Request) ([]byte, error) {

	pods, err := deletePods("", "")
	if err != nil && err != errItemNotExist {
		return nil, fmt.Errorf("could not delete k8s pods: %v", err)
	}
	return sendBytes(pods)
}

func handleDeploymentCreate(w http.ResponseWriter, r *http.Request) ([]byte, error) {

	d, err := createDeployment()
	if err != nil {
		return nil, fmt.Errorf("could not create k8s deployment: %v", err)
	}

	return sendBytes(d)
}

func handleDeploymentDelete(w http.ResponseWriter, r *http.Request) ([]byte, error) {

	_, err := deleteDeployment()
	if err != nil && err != errItemNotExist {
		return nil, fmt.Errorf("could not delete k8s deployment: %v", err)
	}

	_, err = deleteAllReplicaSets()
	if err != nil && err != errItemNotExist {
		return nil, fmt.Errorf("could not delete k8s replica set: %v", err)
	}

	return handlePodsDelete(w, r)
}

func handleNodes(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	nodes, err := listAllNodes()
	if err != nil {
		return nil, fmt.Errorf("could not get list of k8s nodes: %v", err)
	}

	return sendBytes(nodes)
}

func handleNodeDrain(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	nodename := r.FormValue("node")

	node, err := toggleNode(nodename, true)
	if err != nil && err != errItemNotExist {
		return nil, fmt.Errorf("could not retrieve k8s node info: %v", err)
	}

	_, err = deletePods("", nodename)
	if err != nil {
		return nil, fmt.Errorf("could not remove all pods on node: %v", err)
	}

	return sendBytes(node)
}

func handleNodeUncordon(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	nodename := r.FormValue("node")

	node, err := toggleNode(nodename, false)
	if err != nil && err != errItemNotExist {
		return nil, fmt.Errorf("could not retrieve k8s node info: %v", err)
	}

	return sendBytes(node)
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
