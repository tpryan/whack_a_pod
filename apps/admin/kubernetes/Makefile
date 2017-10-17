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

reset: clean.deployment app
	say "App Reset"

reset.safe: 
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o "$(BASEDIR)/../containers/main" "$(BASEDIR)/../containers/main.go"
	gcloud container builds submit "$(BASEDIR)/../containers/." --tag=$(DOCKERREPO)/admin
	-kubectl delete deployment admin-deployment
	-kubectl run admin-deployment --image=$(DOCKERREPO)/admin --replicas=1 --port=8080 --labels=app=admin --env="APIIMAGE=$(DOCKERREPO)/api"
	say "app refresh complete"		

app: build deploy

main:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o "$(BASEDIR)/../containers/main" "$(BASEDIR)/../containers/main.go"

build: env creds main
	gcloud container builds submit "$(BASEDIR)/../containers/." --tag=$(DOCKERREPO)/admin

deploy: env creds deployment service
	
	

deployment:
	kubectl run admin-deployment --image=$(DOCKERREPO)/admin --replicas=1 --port=8080 --labels=app=admin --env="APIIMAGE=$(DOCKERREPO)/api"

service:
	kubectl expose deployment admin-deployment --name=admin --target-port=8080  --type=NodePort --labels="app=admin"

clean: clean.deployment clean.service 

clean.deployment: env creds
	-kubectl delete deployment admin-deployment
	
clean.service: env creds
	-kubectl delete service admin

retry: clean build deploy

config: env 
	echo "No custom config needed for admin"
	