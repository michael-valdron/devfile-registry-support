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
ARG BASE_IMAGE=localhost/devfile-index-base

FROM registry.access.redhat.com/ubi8/go-toolset:1.21 AS builder

# Set user as root
USER root

# Automatically set when --platform flag is set, will default to amd64 if no platform is given
ARG TARGETARCH=amd64

# Install yq
ENV YQ_VERSION=v4.44.1
RUN curl -sL -O https://github.com/mikefarah/yq/releases/download/${YQ_VERSION}/yq_linux_${TARGETARCH} -o /usr/local/bin/yq && mv ./yq_linux_${TARGETARCH} /usr/local/bin/yq && chmod +x /usr/local/bin/yq

# Copy the registry build tools
COPY build-tools /build-tools
# Copy the index generator source
COPY index/ /index
COPY tests/registry /registry

# Download all the offline parent devfiles
RUN bash /build-tools/dl_parent_devfiles.sh

# Download offline starter projects
RUN bash /build-tools/dl_starter_projects.sh go-starter community

# Run the registry build tools
RUN bash /build-tools/build.sh /registry /build

# Set user as non-root
USER 1001

FROM ${BASE_IMAGE}

COPY --from=builder /build /registry