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

deploy: env
	cd "$(BASEDIR)/apps/api/kubernetes/" && $(MAKE) deploy
	cd "$(BASEDIR)/apps/game/kubernetes/" && $(MAKE) deploy
	cd "$(BASEDIR)/apps/admin/kubernetes/" && $(MAKE) deploy
	cd "$(BASEDIR)/apps/ingress/" && $(MAKE) deploy
	

clean: env
	cd "$(BASEDIR)/apps/api/kubernetes/" && $(MAKE) clean
	cd "$(BASEDIR)/apps/game/kubernetes/" && $(MAKE) clean
	cd "$(BASEDIR)/apps/admin/kubernetes/" && $(MAKE) clean	
	cd "$(BASEDIR)/apps/ingress/" && $(MAKE) clean

build: env
	cd "$(BASEDIR)/apps/api/kubernetes/" && $(MAKE) build
	cd "$(BASEDIR)/apps/game/kubernetes/" && $(MAKE) build
	cd "$(BASEDIR)/apps/admin/kubernetes/" && $(MAKE) build		

config: env
	@cd "$(BASEDIR)/apps/ingress/" && $(MAKE) config

test: env
	cd "$(BASEDIR)/apps/api/kubernetes/" && $(MAKE) test
	cd "$(BASEDIR)/apps/admin/kubernetes/" && $(MAKE) test	