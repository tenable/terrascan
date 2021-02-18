#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

go test -v -coverpkg=./pkg/... -coverprofile=coverage.out ./pkg/...
go tool cover -func coverage.out
