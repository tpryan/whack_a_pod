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
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	"k8s.io/api/apps/v1beta1"
	"k8s.io/api/core/v1"
)

var (
	deploymentPath    = "/apis/extensions/v1beta1/namespaces/default/deployments"
	deploymentName    = "api-deployment"
	nodePath          = "/api/v1/nodes"
	nodeName          = "gke-whack-a-pod-default-pool-8deaa3a5-b9p7"
	podDeletePath     = "/api/v1/namespaces/default/pods"
	podExistsName     = "api-deployment-1435701907-xx9lm"
	podExistsSelfLink = "/api/v1/namespaces/default/pods/api-deployment-1435701907-xx9lm"
	podListPath       = "/api/v1/pods"
)

func TestCreateDeploymentNoEnv(t *testing.T) {
	client = new(MockClient)
	_, err := createDeployment()
	if err == nil {
		t.Errorf("Expected error due to no ENV variable")
	}
}

func TestCreateDeployment(t *testing.T) {
	client = new(MockClient)
	os.Setenv("APIIMAGE", "gcr.io/carnivaldemos/api")
	b, err := createDeployment()
	if err != nil {
		t.Errorf("Error getting deployment : %v", err)
	}
	var dep v1beta1.Deployment
	if err := json.Unmarshal(b, &dep); err != nil {
		t.Errorf("error turning response to deployment: %v", err)
		t.FailNow()
	}

	if dep.Name != "api-deployment" {
		t.Errorf("createDeployment().Name got: %s, wanted: %s", dep.Name, "api-deplyment")
	}
	os.Unsetenv("APIIMAGE")
}

func TestDeleteDeployment(t *testing.T) {
	client = new(MockClient)
	_, err := deleteDeployment("api-deployment")
	if err != nil {
		t.Errorf("Error getting deployment : %v", err)
	}
}

func TestToggleNode(t *testing.T) {

	cases := []struct {
		node     string
		inactive bool
	}{
		{"gke-whack-a-pod-default-pool-8deaa3a5-b9p7", true},
		{"gke-whack-a-pod-default-pool-8deaa3a5-b9p7", false},
	}

	client = new(MockClient)

	for _, c := range cases {
		b, err := toggleNode(c.node, c.inactive)
		if err != nil {
			t.Errorf("Error getting deployment : %v", err)
		}

		var node v1.Node
		if err := json.Unmarshal(b, &node); err != nil {
			t.Errorf("error turning response to node: %v", err)
			t.FailNow()
		}

		if node.Spec.Unschedulable != c.inactive {
			t.Errorf("toggleNode().Spec.Unschedulable got %t wanted %t", node.Spec.Unschedulable, c.inactive)
		}

	}

}

func TestListNodes(t *testing.T) {

	client = new(MockClient)

	b, err := listNodes()
	if err != nil {
		t.Errorf("Error getting deployment : %v", err)
	}

	var nodes v1.NodeList
	if err := json.Unmarshal(b, &nodes); err != nil {
		t.Errorf("error turning response to nodelist: %v", err)
		t.FailNow()
	}

	if len(nodes.Items) != 2 {
		t.Errorf("Wanted 2 nodes got : %d", len(nodes.Items))
	}

}

func TestListPods(t *testing.T) {

	client = new(MockClient)

	b, err := listPods()
	if err != nil {
		t.Errorf("Error getting deployment : %v", err)
	}

	var pods v1.PodList
	if err := json.Unmarshal(b, &pods); err != nil {
		t.Errorf("error turning response to podlist: %v", err)
		t.FailNow()
	}

	if len(pods.Items) != 12 {
		t.Errorf("listPods(): got %d wanted : %d", len(pods.Items), 12)
	}

}

func TestDeletePods(t *testing.T) {

	client = new(MockClient)

	b, err := deletePods("")
	if err != nil {
		t.Errorf("Error getting deployment : %v", err)
	}

	var pods v1.PodList
	if err := json.Unmarshal(b, &pods); err != nil {
		t.Errorf("error turning response to podlist: %v", err)
		t.FailNow()
	}

	if len(pods.Items) != 12 {
		t.Errorf("Wanted 12 pods got : %d", len(pods.Items))
	}

}

func TestDescribePod(t *testing.T) {

	client = new(MockClient)
	self := "/api/v1/namespaces/default/pods/api-deployment-1435701907-xx9lm"

	b, err := describePod(self)
	if err != nil {
		t.Errorf("Error getting pod : %v", err)
	}

	var pod v1.Pod
	if err := json.Unmarshal(b, &pod); err != nil {
		t.Errorf("error turning response to podlist: %v", err)
		t.FailNow()
	}

	if pod.ObjectMeta.SelfLink != self {
		t.Errorf("Selflink should be %s got : %s", pod.ObjectMeta.SelfLink, self)
	}

}

func TestDeletePod(t *testing.T) {

	client = new(MockClient)
	self := "/api/v1/namespaces/default/pods/api-deployment-1435701907-xx9lm"

	b, err := deletePod(self)
	if err != nil {
		t.Errorf("Unexpected error getting pod, %v", err)
	}

	var pod v1.Pod
	if err := json.Unmarshal(b, &pod); err != nil {
		t.Errorf("error turning response to podlist: %v", err)
		t.FailNow()
	}

	if pod.ObjectMeta.SelfLink != self {
		t.Errorf("Selflink should be %s got : %s", pod.ObjectMeta.SelfLink, self)
	}

}

func TestDeletePodDoesNotExist(t *testing.T) {
	client = new(MockClient)
	self := "/api/v1/namespaces/default/pods/api-deployment-1435701907-FALSE"

	_, err := deletePod(self)
	if err != errItemNotExist {
		t.Errorf("Unexpected error gerring pod, got %v wanted: %v", err, errItemNotExist)
	}

}

type MockClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {

	r := http.Response{}
	r.StatusCode = http.StatusOK
	r.Request = req
	filename := ""

	switch {

	case req.URL.Path == deploymentPath:
		filename = "testdata/deployment_create.json"

	case req.URL.Path == deploymentPath+"/"+deploymentName:
		filename = "testdata/deployment_delete.json"

	case req.URL.Path == nodePath:
		filename = "testdata/nodes.json"

	case req.URL.Path == podDeletePath+"/"+podExistsName:
		filename = "testdata/pod.json"

	case req.URL.Path == podListPath && req.FormValue("labelSelector") == selector:
		filename = "testdata/pods.json"

	case req.URL.Path == podDeletePath && req.FormValue("labelSelector") == selector && req.FormValue("fieldSelector") == "spec.nodeName="+nodeName:
		filename = "testdata/pods_onenode_only.json"

	case req.URL.Path == podDeletePath && req.FormValue("labelSelector") == selector:
		filename = "testdata/pods.json"

	case req.URL.Path == nodePath+"/"+nodeName:
		b, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return nil, fmt.Errorf("could not read request data : %v", err)
		}

		filename = "testdata/node_schedulable.json"
		if strings.Index(string(b), "true") >= 0 {
			filename = "testdata/node_unschedulable.json"
		}
	default:
		filename = "testdata/empty.json"
		r.StatusCode = http.StatusNotFound
	}

	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("could not get cached file (%s) for test: %v", filename, err)
	}
	reader := bufio.NewReader(f)
	r.Body = ioutil.NopCloser(reader)

	return &r, nil

}
