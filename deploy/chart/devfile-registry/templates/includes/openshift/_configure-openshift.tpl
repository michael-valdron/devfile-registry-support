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
{{ define "devfileregistry.openshift.configure" }}
- name: configure-openshift
  image: docker.io/alpine/k8s:1.32.1@sha256:11e7706b73a219147a1c7cd2ddd8fb27768200ce87207b2808007f9a52d5c119
  workingDir: /tmp
  command:
    - /bin/sh
    - -c
    - |
      set -o errexit
      set -o nounset
      set -o pipefail
      
      ROUTE=''
      ROUTE_TIMEOUT_SEC=25
      DEPLOY_TIMEOUT_SEC=60
      START_TIME=$(date +%s)

      echo -n "* Waiting for Devfile Registry route: "
      while [ -z "${ROUTE}" ]; do
          ROUTE=$(kubectl get route {{ template "devfileregistry.fullname" . }} -o jsonpath='{.spec.host}' 2>/dev/null)
          new_time=$(date +%s)
          elapsed_time=$((new_time - START_TIME))
          if [ $elapsed_time -ge $ROUTE_TIMEOUT_SEC ]; then
              echo "TIMEOUT: Domain not found."
              exit 1
          fi
          echo -n "_"
          sleep 1
      done
      echo "OK"
      echo Domain found: "$ROUTE"

      {{ $gitUrl := .Values.chartUrlOverride | default (index .Chart.Sources 0) }}
      echo -n "* Cloning helm chart from '{{ $gitUrl }}': "
      REGISTRY_SUPPORT_REPO=$(mktemp -d)
      CHART=${REGISTRY_SUPPORT_REPO}/deploy/chart/devfile-registry
      {{ if hasKey .Values "chartRefOverride" }}
      git clone {{ $gitUrl }} -b {{ .Values.chartRefOverride }} ${REGISTRY_SUPPORT_REPO} >/dev/null 2>&1
      {{ else }}
      git clone {{ $gitUrl }} ${REGISTRY_SUPPORT_REPO} >/dev/null 2>&1
      {{ end }}
      echo "OK"

      echo -n "* Waiting for Devfile Registry deployment: "
      kubectl wait --for=condition=Available deployment/{{ template "devfileregistry.fullname" . }} --timeout=${DEPLOY_TIMEOUT_SEC}s >/dev/null 2>&1
      echo "OK"

      # Run upgrade with the new variable to set fqdn of the Registry Viewer
      echo -n "* Update deployment with route: "
      helm upgrade {{ .Release.Name }} ${CHART} --version {{ .Release.Revision }} --reuse-values --set global.route.domain=$ROUTE >/dev/null 2>&1
      if [ $? -ne 0 ]; then
        echo "problem updating the helm release"
        exit 1
      fi
      echo "OK"
{{ end }}
