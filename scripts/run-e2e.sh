#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

export TERRASCAN_BIN_PATH=${PWD}/bin/terrascan

go test -p 1 -v $(go list ./test/e2e/... | grep -v /vulnerability) 