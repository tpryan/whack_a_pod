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
	gcloud container builds submit "$(BASEDIR)/../containers/." --tag=gcr.io/$(PROJECT)/admin
	-kubectl delete -f "$(BASEDIR)/admin-deployment.yaml"	
	-kubectl apply -f "$(BASEDIR)/admin-deployment.yaml"
	say "app refresh complete"		

app: build deploy

build: creds
	gcloud container builds submit "$(BASEDIR)/../containers/." --tag=gcr.io/$(PROJECT)/admin

deploy: creds
	kubectl apply -f "$(BASEDIR)/admin-service.yaml"
	kubectl apply -f "$(BASEDIR)/admin-deployment.yaml"	

clean: clean.deployment clean.service 

clean.deployment: env
	-kubectl delete -f "$(BASEDIR)/admin-deployment.yaml"
	
clean.service: env
	-kubectl delete -f "$(BASEDIR)/admin-service.yaml"

retry: clean build deploy

config:
	@echo Creating ADMIN Yaml files based on samples and setting in your Makefile.properties
	@cp "$(BASEDIR)/admin-deployment.sample.yaml" "$(BASEDIR)/admin-deployment.yaml"
	$(call rewritefile,"$(BASEDIR)/admin-deployment.yaml",%PROJECT%,$(PROJECT))
	@cp "$(BASEDIR)/deployment.json" "$(BASEDIR)/../containers/default/api/k8s/createdeploy/deployment.json"
	$(call rewritefile,"$(BASEDIR)/../containers/default/api/k8s/createdeploy/deployment.json",%PROJECT%,$(PROJECT))

define rewritefile
	@sed s/$(2)/$(3)/g <""$(1)"" >"$(BASEDIR)/.temp"
	@cp "$(BASEDIR)/.temp" $(1)
	@rm "$(BASEDIR)/.temp" 
endef