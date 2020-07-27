#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail


export GO111MODULE=on
go get honnef.co/go/tools/cmd/staticcheck
staticcheck -f stylish ./...
