#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

DOCKER_REPO="docker-terrascan-local.artifactory.eng.tenable.com/terrascan"
DOCKERFILE="./build/Dockerfile"

if [ $# -eq 0 ]; then
  LABEL=$(git rev-parse --short HEAD 2>/dev/null)
elif [ $# -eq 1 ]; then
  case "$1" in
    latest)
      LABEL="latest"
      ;;
    tag)
      LATEST_TAG=$(git describe --abbrev=0 --tags)
      LABEL=$(echo "${LATEST_TAG//v}")
      ;;
      *)
  esac
fi

if [ "${LABEL-false}" = "false" ]; then
  echo "Usage:"
  echo "  $0          ->  label is the git commit"
  echo "  $0  tag     ->  label is the latest tag"
  echo "  $0  latest  ->  label is 'latest'"
  exit 1
fi

declare -a PLATFORM
if [ "${MULTIPLATFORM-false}" = "true" ]; then
  OUTPUT_TYPE="--push"
  PLATFORM=("--platform" "linux/amd64,linux/arm64")
else
  OUTPUT_TYPE="--load"
fi

docker buildx create "${PLATFORM[@]}" --name terrascan-builder --use

docker buildx build --provenance=false "${OUTPUT_TYPE}" "${PLATFORM[@]}" -t "${DOCKER_REPO}:${LABEL}" -f "${DOCKERFILE}" .

echo "${LABEL}" > dockerhub-image-label.txt

docker buildx rm terrascan-builder
