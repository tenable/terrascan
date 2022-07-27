#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

find_files() {
  find . -not \( \
      \( \
        -wholename '*/vendor/*' \
      \) -prune \
    \) -name '*.go'
}


export GO111MODULE=on
export PATH=$PATH:$(go env GOPATH)/bin
go get -d golang.org/x/lint/golint
go install golang.org/x/lint/golint

bad_files=$(find_files | xargs -I@ bash -c "$GOPATH/bin/golint @")
if [[ -n "${bad_files}" ]]; then
  echo "${bad_files}"
  exit 1
fi
