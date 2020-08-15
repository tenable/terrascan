#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null)
DOCKER_REPO="accurics/terrascan"
DOCKERFILE="./build/Dockerfile"
LATEST_TAG="latest"

# PS: It is a prerequisite to execute 'docker login' before running this script
docker tag ${DOCKER_REPO}:${GIT_COMMIT} ${DOCKER_REPO}:${LATEST_TAG}
docker push ${DOCKER_REPO}:${LATEST_TAG}
