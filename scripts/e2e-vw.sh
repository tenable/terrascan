#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null)
DOCKER_REPO="accurics/terrascan"

docker tag ${DOCKER_REPO}:${GIT_COMMIT} ${DOCKER_REPO}:validating-webhook
export REPO_NAME=$DOCKER_REPO

go test -p 1 -v ./test/e2e/validating-webhook/...