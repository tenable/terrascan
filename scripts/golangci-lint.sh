#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

export GO111MODULE=on
export PATH=$PATH:$(go env GOPATH)/bin
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
golangci-lint run ./...