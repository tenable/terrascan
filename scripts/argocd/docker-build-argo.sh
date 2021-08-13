#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null)
DOCKER_REPO="accurics/terrascan_argocd"
DIR="."
FILE="./build/argo-Dockerfile"

docker build -t ${DOCKER_REPO}:${GIT_COMMIT} -f ${FILE} ${DIR}
