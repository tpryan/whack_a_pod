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

include ../Makefile.properties
NODEDISKSIZE="200"
NODECOUNT="2"
NODETYPE="n1-standard-4"

# Builds the entire infrastructure you need to run this set of demos. 
build: env build.cluster ips

# Requests a GKE cluster
build.cluster: 
	-gcloud container clusters create $(CLUSTER) --num-nodes $(NODECOUNT) \
	--disk-size=$(NODEDISKSIZE) --machine-type=$(NODETYPE) --zone=$(ZONE)

# Runs all the clean commands. 
clean: env clean.cluster clean.repo

# deletes the cluster. 
clean.cluster:
	gcloud container clusters delete $(CLUSTER) -q

# Gets rid of all of the images in GKR
clean.repo:
	-gsutil -m rm -r gs://artifacts.$(PROJECT).appspot.com

ips:
	gcloud compute --project $(PROJECT) addresses create "api" --region $(REGION)
	gcloud compute --project $(PROJECT) addresses create "admin" --region $(REGION)
	gcloud compute --project $(PROJECT) addresses create "game" --region $(REGION)			