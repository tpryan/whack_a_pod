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

build: creds
	docker build -t gcr.io/$(PROJECT)/game "$(BASEDIR)/../containers/."
	gcloud docker -- push gcr.io/$(PROJECT)/game

deploy: creds
	kubectl apply -f "$(BASEDIR)/game-service.yaml"
	kubectl apply -f "$(BASEDIR)/game-deployment.yaml"	

clean: clean.deployment clean.service 

clean.deployment: env
	-kubectl delete -f "$(BASEDIR)/game-deployment.yaml"
	
clean.service: env
	-kubectl delete -f "$(BASEDIR)/game-service.yaml"

retry: clean build deploy

config:
	@echo Creating GAME Yaml files based on samples and setting in your Makefile.properties
	@cp "$(BASEDIR)/game-deployment-sample.yaml" "$(BASEDIR)/game-deployment.yaml"
	$(call rewritefile,"$(BASEDIR)/game-deployment.yaml",%PROJECT%,$(PROJECT))
	@cp "$(BASEDIR)/game-service-sample.yaml" "$(BASEDIR)/game-service.yaml"
	$(call rewritefile,"$(BASEDIR)/game-service.yaml",%GAMEIP%,$(GAMEIP))
	@cp "$(BASEDIR)/../containers/default/assets/js/config.sample.js" "$(BASEDIR)/../containers/default/assets/js/config.js"
	$(call rewritefile,"$(BASEDIR)/../containers/default/assets/js/config.js",%APIHOST%,$(APIHOST))
	$(call rewritefile,"$(BASEDIR)/../containers/default/assets/js/config.js",%ADMINHOST%,$(ADMINHOST))

define rewritefile
	@sed s/$(2)/$(3)/g <""$(1)"" >"$(BASEDIR)/.temp"
	@cp "$(BASEDIR)/.temp" $(1)
	@rm "$(BASEDIR)/.temp" 
endef