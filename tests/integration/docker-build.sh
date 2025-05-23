#!/bin/sh

#
# Copyright Red Hat
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

#set the docker alias if necessary
. ../../setenv.sh

# Devfile Registry library build arguments
LIBRARY_REPO=${LIBRARY_REPO:-"https://github.com/devfile/registry-support"}
LIBRARY_REF=${LIBRARY_REF:-"main"}

# LICENSE build arguments
LICENSE_REPO=${LICENSE_REPO:-"devfile/registry-support"}
LICENSE_REF=${LICENSE_REF:-"main"}

# Build the container image
docker build -t devfile-registry-integration \
    --build-arg LIBRARY_REPO=${LIBRARY_REPO} \
    --build-arg LIBRARY_REF=${LIBRARY_REF} \
    --build-arg LICENSE_REPO=${LICENSE_REPO} \
    --build-arg LICENSE_REF=${LICENSE_REF} ./
