#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

go test -v -coverpkg=./... -coverprofile=coverage.out ./...
go tool cover -func coverage.out
