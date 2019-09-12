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

include Makefile.properties

deploy: env creds
	cd "$(BASEDIR)/apps/api/kubernetes/" && $(MAKE) deploy
	cd "$(BASEDIR)/apps/game/kubernetes/" && $(MAKE) deploy
	cd "$(BASEDIR)/apps/admin/kubernetes/" && $(MAKE) deploy
	cd "$(BASEDIR)/apps/ingress/" && $(MAKE) deploy

reset.safe: env creds
	cd "$(BASEDIR)/apps/api/kubernetes/" && $(MAKE) reset.safe
	cd "$(BASEDIR)/apps/game/kubernetes/" && $(MAKE) reset.safe
	cd "$(BASEDIR)/apps/admin/kubernetes/" && $(MAKE) reset.safe

deploy.minikube: creds.minikube
	cd "$(BASEDIR)/apps/api/kubernetes/" && $(MAKE) deploy.minikube
	cd "$(BASEDIR)/apps/game/kubernetes/" && $(MAKE) deploy.minikube
	cd "$(BASEDIR)/apps/admin/kubernetes/" && $(MAKE) deploy.minikube
	cd "$(BASEDIR)/apps/ingress/" && $(MAKE) deploy.minikube

deploy.minikube.dockerhub: creds.minikube
	minikube addons enable ingress
	cd "$(BASEDIR)/apps/api/kubernetes/" && $(MAKE) deploy.minikube.dockerhub
	cd "$(BASEDIR)/apps/game/kubernetes/" && $(MAKE) deploy.minikube.dockerhub
	cd "$(BASEDIR)/apps/admin/kubernetes/" && $(MAKE) deploy.minikube.dockerhub
	cd "$(BASEDIR)/apps/ingress/" && $(MAKE) deploy.minikube
	@printf -- "*** DONE ***\n"
	@printf -- "add the following line to your /etc/hosts file:\n\n"
	@printf -- "$(shell minikube ip) minikube.wap\n\n"

deploy.generic: 
	cd "$(BASEDIR)/apps/api/kubernetes/" && $(MAKE) deploy.generic
	cd "$(BASEDIR)/apps/game/kubernetes/" && $(MAKE) deploy.generic
	cd "$(BASEDIR)/apps/admin/kubernetes/" && $(MAKE) deploy.generic
	cd "$(BASEDIR)/apps/ingress/" && $(MAKE) deploy.generic

clean: env creds
	cd "$(BASEDIR)/apps/api/kubernetes/" && $(MAKE) clean
	cd "$(BASEDIR)/apps/game/kubernetes/" && $(MAKE) clean
	cd "$(BASEDIR)/apps/admin/kubernetes/" && $(MAKE) clean	
	cd "$(BASEDIR)/apps/ingress/" && $(MAKE) clean

clean.generic: 
	cd "$(BASEDIR)/apps/api/kubernetes/" && $(MAKE) clean.generic
	cd "$(BASEDIR)/apps/game/kubernetes/" && $(MAKE) clean.generic
	cd "$(BASEDIR)/apps/admin/kubernetes/" && $(MAKE) clean.generic
	cd "$(BASEDIR)/apps/ingress/" && $(MAKE) clean.generic

clean.minikube: 
	cd "$(BASEDIR)/apps/api/kubernetes/" && $(MAKE) clean.minikube
	cd "$(BASEDIR)/apps/game/kubernetes/" && $(MAKE) clean.minikube
	cd "$(BASEDIR)/apps/admin/kubernetes/" && $(MAKE) clean.minikube	
	cd "$(BASEDIR)/apps/ingress/" && $(MAKE) clean.minikube

clean.minikube.dockerhub: 
	cd "$(BASEDIR)/apps/api/kubernetes/" && $(MAKE) clean.minikube
	cd "$(BASEDIR)/apps/game/kubernetes/" && $(MAKE) clean.minikube
	cd "$(BASEDIR)/apps/admin/kubernetes/" && $(MAKE) clean.minikube.dockerhub
	cd "$(BASEDIR)/apps/ingress/" && $(MAKE) clean.minikube

build: env creds
	cd "$(BASEDIR)/apps/api/kubernetes/" && $(MAKE) build
	cd "$(BASEDIR)/apps/game/kubernetes/" && $(MAKE) build
	cd "$(BASEDIR)/apps/admin/kubernetes/" && $(MAKE) build

build.dockerhub:
	cd "$(BASEDIR)/apps/api/kubernetes/" && $(MAKE) build.dockerhub
	cd "$(BASEDIR)/apps/game/kubernetes/" && $(MAKE) build.dockerhub
	cd "$(BASEDIR)/apps/admin/kubernetes/" && $(MAKE) build.dockerhub

build.generic: 
	cd "$(BASEDIR)/apps/api/kubernetes/" && $(MAKE) build.generic
	cd "$(BASEDIR)/apps/game/kubernetes/" && $(MAKE) build.generic
	cd "$(BASEDIR)/apps/admin/kubernetes/" && $(MAKE) build.generic

config: env creds
	@cd "$(BASEDIR)/apps/ingress/" && $(MAKE) config

test: 
	cd "$(BASEDIR)/apps/api/kubernetes/" && $(MAKE) test
	cd "$(BASEDIR)/apps/admin/kubernetes/" && $(MAKE) test	
