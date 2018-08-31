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
	gcloud container builds submit "$(BASEDIR)/../containers/." --tag=$(DOCKERREPO)/game
	-kubectl delete deployment game-deployment
	kubectl run game-deployment --image=$(DOCKERREPO)/game --replicas=4 --port=8080 --labels=app=game 
	say "app refresh complete"	

build: env
	gcloud container builds submit "$(BASEDIR)/../containers/." --tag=$(DOCKERREPO)/game

build.dockerhub:
	docker build "$(BASEDIR)/../containers/." --tag=$(DOCKERREPO):game
	docker push $(DOCKERREPO):game

build.generic: 
	docker build "$(BASEDIR)/../containers/." --tag=$(DOCKERREPO)/whackapod-game
	docker push $(DOCKERREPO)/whackapod-game

deploy: env creds deployment service

deploy.minikube: deployment service

deploy.minikube.dockerhub:
	kubectl run game-deployment --image=$(DOCKERREPO):game --replicas=4 --port=8080 --labels=app=game 
	kubectl expose deployment game-deployment --name=game --target-port=8080  --type=NodePort --labels="app=game"

deploy.generic: 
	kubectl run game-deployment --image=$(DOCKERREPO)/whackapod-game --replicas=4 --port=8080 --labels=app=game 
	kubectl expose deployment game-deployment --name=game --target-port=8080  --type=NodePort --labels="app=game"

deployment:
	kubectl run game-deployment --image=$(DOCKERREPO)/game --replicas=4 --port=8080 --labels=app=game 

service:
	kubectl expose deployment game-deployment --name=game --target-port=8080  --type=NodePort --labels="app=game"

clean: env creds clean.deployment clean.service 

clean.generic: clean.deployment clean.service

clean.minikube: clean.deployment clean.service 

clean.deployment: 
	-kubectl delete deployment game-deployment
	
clean.service: 
	-kubectl delete service game


retry: clean build deploy

config: env
