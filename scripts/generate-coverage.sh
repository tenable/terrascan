#!/bin/bash

# source: https://github.com/codecov/example-go
# go test can't generate code coverage for multiple packages in one command

 set -e
touch coverage.out
go test -i -race ./cmd/terrascan
for d in $(go list ./... | grep -v vendor | grep -v tests | grep -v integration_test); do
    go test -race -coverprofile=profile.out -covermode=atomic $d
    if [ -f profile.out ]; then
        cat profile.out >> coverage.out
        rm profile.out
    fi
done
