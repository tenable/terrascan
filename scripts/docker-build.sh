#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null)
DOCKER_REPO="tenable/terrascan"
DOCKERFILE="./build/Dockerfile"

docker buildx create --platform linux/amd64,linux/arm64 --name terrascan-builder --use

docker buildx build -t ${DOCKER_REPO}:${GIT_COMMIT} -f ${DOCKERFILE} . --load

docker buildx rm terrascan-builder
