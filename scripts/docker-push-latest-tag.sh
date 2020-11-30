#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null)
DOCKER_REPO="accurics/terrascan"
DOCKERFILE="./build/Dockerfile"
LATEST_TAG=$(git describe --abbrev=0 --tags)
LATEST_TAG_SHORT=$(echo "${LATEST_TAG//v}")

# PS: It is a prerequisite to execute 'docker login' before running this script
docker tag ${DOCKER_REPO}:${GIT_COMMIT} ${DOCKER_REPO}:${LATEST_TAG_SHORT}
docker push ${DOCKER_REPO}:${LATEST_TAG_SHORT}
