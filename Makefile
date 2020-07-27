GITCOMMIT := $(shell git rev-parse --short HEAD 2>/dev/null)
BUILD_FLAGS := -v -ldflags "-w -s"

BUILD_DIR = ./bin
BINARY_NAME = terrascan


# default
default: help


# help
help:
	@echo "usage: make [command]\ncommands:"
	@echo "build\n\tbuild terrascan binary"
	@echo "cicd\n\tsimulate CI/CD pipeline locally"
	@echo "clean\n\tclean up build"
	@echo "gofmt\n\tvalidate gofmt"
	@echo "golint\n\tvalidate golint"
	@echo "gomodverify\n\tverify go modules"
	@echo "govet\n\tvalidate govet"
	@echo "staticcheck\n\trun static code analysis"
	@echo "test\n\texecute unit and integration tests"
	@echo "unit-tests\n\texecute unit tests"
	@echo "validate\n\trun all validations"

# build terrascan binary
build: clean
	@mkdir -p $(BUILD_DIR) > /dev/null
	go build ${BUILD_FLAGS} -o ${BUILD_DIR}/${BINARY_NAME} cmd/terrascan/main.go
	@echo "binary created at ${BUILD_DIR}/${BINARY_NAME}"


# clean build 
clean:
	@rm -rf $(BUILD_DIR)


# run all cicd steps
cicd: validate build test


# run all unit and integration tests
test: unit-tests


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


# static code analysis
staticcheck:
	./scripts/staticcheck.sh


# run unit tests
unit-tests:
	./scripts/generate-coverage.sh
