#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null)
DOCKER_REPO="tenable/terrascan"

# PS: It is a prerequisite to execute 'docker login' before running this script
docker push ${DOCKER_REPO}:${GIT_COMMIT}
