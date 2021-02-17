#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

if [ $# -eq 2 ]; then
    DOCKER_REPO=$1
    TAG=$2
fi

if [ $# -eq 1 ]; then
    DOCKER_REPO=$1
fi

GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null)
DOCKERFILE="./build/Dockerfile"

if [ -z $DOCKER_REPO ]; then
    DOCKER_REPO="accurics/terrascan"
fi

if [ -z $TAG ]; then
    TAG=${GIT_COMMIT}
fi

echo "creating docker image : $DOCKER_REPO/$TAG"
docker build -t ${DOCKER_REPO}:${TAG} -f ${DOCKERFILE} .