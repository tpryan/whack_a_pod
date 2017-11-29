# Copyright 2017 Google Inc. All Rights Reserved.

# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at

#     http://www.apache.org/licenses/LICENSE-2.0

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
FROM bitnami/nginx

COPY default.conf /bitnami/nginx/conf
COPY default /app

USER root

RUN apt-get update && apt-get install -y curl
HEALTHCHECK CMD curl --fail http://localhost:8080/healthz/

RUN chmod -R 755 /app
RUN chown -R 1001 /app
RUN chgrp -R 1001 /app

USER 1001
