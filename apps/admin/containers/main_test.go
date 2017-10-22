// Copyright 2017 Google Inc. All Rights Reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//     http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"k8s.io/api/apps/v1beta2"
	"k8s.io/api/core/v1"
)

func TestHealth(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	req, err := http.NewRequest("GET", "/healthz", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(health)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `ok`
	if rr.Body.String() != expected {
		t.Errorf("unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestHandlePods(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	client = new(MockClient)
	req, err := http.NewRequest("GET", "/admin/k8s/pods/get", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleAPI(handlePods))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("wrong status code: got %v want %v", status, http.StatusOK)
	}

	var podlist v1.PodList
	if err := json.Unmarshal([]byte(rr.Body.String()), &podlist); err != nil {
		t.Errorf("error turning response to podlist: %v", err)
		t.FailNow()
	}

	if len(podlist.Items) < 12 {
		t.Errorf("podlist.Items: got %d want %d", len(podlist.Items), 12)
	}

}

func TestHandlePodDeleteExisting(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	client = new(MockClient)
	req, err := http.NewRequest("GET", "/admin/k8s/pod/delete?pod="+podExistsSelfLink, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleAPI(handlePodDelete))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("wrong status code: got %v want %v", status, http.StatusOK)
	}

	var pod v1.Pod
	if err := json.Unmarshal([]byte(rr.Body.String()), &pod); err != nil {
		t.Errorf("error turning response to pod: %v", err)
		t.FailNow()
	}

	if pod.SelfLink != podExistsSelfLink {
		t.Errorf("podlist.SelfLink: got %s want %s", pod.SelfLink, podExistsSelfLink)
	}

}

func TestHandlePodDeleteNonExisting(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	client = new(MockClient)
	req, err := http.NewRequest("GET", "/admin/k8s/pod/delete?pod=dsadasdasdsa", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleAPI(handlePodDelete))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusAccepted {
		t.Errorf("wrong status code: got %v want %v", status, http.StatusAccepted)
	}

}

func TestHandlePodsDeleteAll(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	client = new(MockClient)
	req, err := http.NewRequest("GET", "/admin/k8s/pods/delete", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleAPI(handlePodsDelete))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("wrong status code: got %v want %v", status, http.StatusOK)
	}

	var podlist v1.PodList
	if err := json.Unmarshal([]byte(rr.Body.String()), &podlist); err != nil {
		t.Errorf("error turning response to pod: %v", err)
		t.FailNow()
	}

	if len(podlist.Items) != 12 {
		t.Errorf("len(podlist.Items): got %d want %d", len(podlist.Items), 12)
	}

}

func TestHandleDeploymentCreate(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	client = new(MockClient)
	os.Setenv("APIIMAGE", "gcr.io/carnivaldemos/api")
	req, err := http.NewRequest("GET", "/admin/k8s/deployment/create", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleAPI(handleDeploymentCreate))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("wrong status code: got %v want %v", status, http.StatusOK)
	}

	var d v1beta2.Deployment
	if err := json.Unmarshal([]byte(rr.Body.String()), &d); err != nil {
		t.Errorf("error turning response to deployment: %v", err)
		t.FailNow()
	}

	if d.Name != deploymentName {
		t.Errorf("Deployment Name: got %s want %s", d.Name, deploymentName)
	}
	os.Unsetenv("APIIMAGE")
}

func TestHandleDeploymentDelete(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	client = new(MockClient)
	req, err := http.NewRequest("GET", "/admin/k8s/deployment/delete", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleAPI(handleDeploymentDelete))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("wrong status code: got %v want %v", status, http.StatusOK)
	}

	var podlist v1.PodList
	if err := json.Unmarshal([]byte(rr.Body.String()), &podlist); err != nil {
		t.Errorf("error turning response to pod: %v", err)
		t.FailNow()
	}

	if len(podlist.Items) < 12 {
		t.Errorf("podlist.Items: got %d want %d", len(podlist.Items), 12)
	}
}

func TestHandleDeploymentCreateNoEnvSet(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	client = new(MockClient)
	req, err := http.NewRequest("GET", "/admin/k8s/deployment/create", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleAPI(handleDeploymentCreate))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}

}

func TestHandleNodes(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	client = new(MockClient)
	req, err := http.NewRequest("GET", "/admin/k8s/nodes/get", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleAPI(handleNodes))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("wrong status code: got %v want %v", status, http.StatusOK)
	}

	var nodeList v1.NodeList
	if err := json.Unmarshal([]byte(rr.Body.String()), &nodeList); err != nil {
		t.Errorf("error turning response to podlist: %v", err)
		t.FailNow()
	}

	if len(nodeList.Items) < 2 {
		t.Errorf("nodelist.Items: got %d want %d", len(nodeList.Items), 2)
	}

}

func TestHandleNodeDrain(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	client = new(MockClient)
	req, err := http.NewRequest("GET", "/admin/k8s/node/drain?node="+nodeName, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleAPI(handleNodeDrain))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("wrong status code: got %v want %v", status, http.StatusOK)
	}

	var node v1.Node
	if err := json.Unmarshal([]byte(rr.Body.String()), &node); err != nil {
		t.Errorf("error turning response to podlist: %v", err)
		t.FailNow()
	}

	if !node.Spec.Unschedulable {
		t.Errorf("node.Spec.Unschedulable: got %t want %t", node.Spec.Unschedulable, true)
	}

}

func TestHandleNodeUncordon(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	client = new(MockClient)
	req, err := http.NewRequest("GET", "/admin/k8s/node/uncordon?node="+nodeName, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleAPI(handleNodeUncordon))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("wrong status code: got %v want %v", status, http.StatusOK)
	}

	var node v1.Node
	if err := json.Unmarshal([]byte(rr.Body.String()), &node); err != nil {
		t.Errorf("error turning response to podlist: %v", err)
		t.FailNow()
	}

	if node.Spec.Unschedulable {
		t.Errorf("node.Spec.Unschedulable: got %t want %t", node.Spec.Unschedulable, false)
	}

}
