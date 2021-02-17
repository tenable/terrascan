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

if [ -z $DOCKER_REPO ]; then
    DOCKER_REPO="accurics/terrascan"
fi

if [ -z $TAG ]; then
    TAG=${GIT_COMMIT}
fi

# PS: make sure you've logged in into the docker. if not, `docker login`
docker push ${DOCKER_REPO}:${TAG}
