GITCOMMIT := $(shell git rev-parse --short HEAD 2>/dev/null)
BUILD_FLAGS := -v -ldflags "-w -s"

BUILD_DIR = ./bin
BINARY_NAME = terrascan


# default
default: build


# build terrascan binary
build: clean
	@mkdir -p $(BUILD_DIR) > /dev/null
	go build ${BUILD_FLAGS} -o ${BUILD_DIR}/${BINARY_NAME} cmd/terrascan/main.go
	@echo "binary created at ${BUILD_DIR}/${BINARY_NAME}"


# clean build 
clean:
	@rm -rf $(BUILD_DIR)


# run all validation tests
validate: gofmt govet golint gomodverify


# gofmt validation
gofmt:
	./scripts/validate-gofmt.sh


# golint validation
golint:
	./scripts/validate-golint.sh


# govet validation
govet:
	./scripts/validate-govet.sh


# go mod validation
gomodverify:
	go mod verify


# run unit tests
unit-tests:
	./scripts/generate-coverage.sh
