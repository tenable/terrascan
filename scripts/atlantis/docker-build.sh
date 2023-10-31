#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null)
DOCKER_REPO="docker-terrascan-local.artifactory.eng.tenable.com/terrascan_atlantis"
DIR="./integrations/atlantis"

docker build -t ${DOCKER_REPO}:${GIT_COMMIT} ${DIR}
