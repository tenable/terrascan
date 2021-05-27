#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

export TERRASCAN_BIN_PATH=${PWD}/bin/terrascan
export DEFAULT_CHART_PATH=${PWD}/deploy/helm-charts

go test -p 1 -v ./test/...