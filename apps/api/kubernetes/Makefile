# Copyright 2017 Google Inc. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
BASEDIR = $(shell pwd)

include ../../../Makefile.properties

app: build deploy

reset: clean.deployment app
	say "app refresh complete"

reset.safe: env creds
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o "$(BASEDIR)/../containers/main" "$(BASEDIR)/../containers/main.go"
	gcloud container builds submit "$(BASEDIR)/../containers/." --tag=$(DOCKERREPO)/api
	-kubectl delete deployment api-deployment
	-kubectl run api-deployment --image=$(DOCKERREPO)/api --replicas=12 --port=8080 --labels=app=api 
	say "app refresh complete"		

main:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o "$(BASEDIR)/../containers/main" "$(BASEDIR)/../containers/main.go"

build: env main
	gcloud container builds submit "$(BASEDIR)/../containers/." --tag=$(DOCKERREPO)/api

build.dockerhub: main
	docker build "$(BASEDIR)/../containers/." --tag=$(DOCKERREPO):api
	docker push $(DOCKERREPO):api

build.generic: main
	docker build "$(BASEDIR)/../containers/." --tag=$(DOCKERREPO)/whackapod-api
	docker push $(DOCKERREPO)/whackapod-api

deploy: env creds deployment service

deploy.minikube: deployment service

deploy.minikube.dockerhub:
	kubectl run api-deployment --image=$(DOCKERREPO):api --replicas=12 --port=8080 --labels=app=api
	kubectl expose deployment api-deployment --name=api --target-port=8080  --type=NodePort --labels="app=api"

deploy.generic: 
	kubectl run api-deployment --image=$(DOCKERREPO)/whackapod-api --replicas=12 --port=8080 --labels=app=api 
	kubectl expose deployment api-deployment --name=api --target-port=8080  --type=NodePort --labels="app=api"

test: 
	cd 	"$(BASEDIR)/../containers" && go test

deployment:
	kubectl run api-deployment --image=$(DOCKERREPO)/api --replicas=12 --port=8080 --labels=app=api 

service:
	kubectl expose deployment api-deployment --name=api --target-port=8080  --type=NodePort --labels="app=api"

clean: env creds clean.deployment clean.service 
clean.generic: clean.deployment clean.service 

clean.minikube: clean.deployment clean.service 

clean.deployment: 
	-kubectl delete deployment api-deployment
	
clean.service: 
	-kubectl delete service api

retry: clean build deploy

config: env 
	@echo Creating API Yaml files based on samples and setting in your Makefile.properties
	@cp "$(BASEDIR)/api-deployment.sample.yaml" "$(BASEDIR)/api-deployment.yaml"
	$(call rewritefile,"$(BASEDIR)/api-deployment.yaml",%PROJECT%,$(PROJECT))

define rewritefile
	@sed s/$(2)/$(3)/g <""$(1)"" >"$(BASEDIR)/.temp"
	@cp "$(BASEDIR)/.temp" $(1)
	@rm "$(BASEDIR)/.temp" 
endef
