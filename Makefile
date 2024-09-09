GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null)
BUILD_FLAGS := -v -ldflags "-w -s -X google.golang.org/protobuf/reflect/protoregistry.conflictPolicy=ignore"

BUILD_DIR = ./bin
BINARY_NAME = terrascan

# default
default: help


# please keep the commands in lexicographical order
help:
	@echo "usage: make [command]\ncommands:"
	@echo "build\n\tbuild terrascan binary"
	@echo "cicd\n\tsimulate CI/CD pipeline locally"
	@echo "clean\n\tclean up build"
	@echo "docker-build\n\tbuild terrascan docker image"
	@echo "docker-push\n\tpush terrascan docker image"
	@echo "docker-push-latest\n\tpush terrascan docker image with latest tag"
	@echo "docker-push-latest-tag\n\tpush terrascan docker image with latest release tag"
	@echo "docker-atlantis-build\n\tbuild terrascan_atlantis docker image"
	@echo "docker-atlantis-push\n\tpush terrascan_atlantis docker image"
	@echo "docker-atlantis-push-latest\n\tpush terrascan_atlantis docker image with latest tag"
	@echo "docker-atlantis-push-latest-tag\n\tpush terrascan_atlantis docker image with latest release tag"
	@echo "gofmt\n\tvalidate gofmt"
	@echo "golint\n\tvalidate golint"
	@echo "gomodverify\n\tverify go modules"
	@echo "govet\n\tvalidate govet"
	@echo "staticcheck\n\trun static code analysis"
	@echo "test\n\texecute unit and integration tests"
	@echo "unit-tests\n\texecute unit tests"
	@echo "e2e-tests\n\texecute e2e tests"
	@echo "e2e-admission-control-tests\n\texecute e2e admission control tests"
	@echo "e2e-vulnerability-tests\n\texecute e2e vulnerability tests"
	@echo "validate\n\trun all validations"

# build terrascan binary
build: clean
	@mkdir -p $(BUILD_DIR) > /dev/null
	@export GO111MODULE=on
	go build ${BUILD_FLAGS} -o ${BUILD_DIR}/${BINARY_NAME} cmd/terrascan/main.go
	@echo "binary created at ${BUILD_DIR}/${BINARY_NAME}"


# clean build
clean:
	@rm -rf $(BUILD_DIR)


# run all cicd steps
cicd: validate build test docker-build


# run all unit and integration tests
test: unit-tests e2e-tests


# run all validation tests
validate: gofmt govet golint gomodverify staticcheck


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
	GOLANG_PROTOBUF_REGISTRATION_CONFLICT=ignore ./scripts/generate-coverage.sh

# run e2e tests
e2e-tests: build
	./scripts/run-e2e.sh

# run e2e validating webhook
e2e-admission-control-tests: build
	./scripts/e2e-admission-control.sh

# run e2e vulnerability tests
e2e-vulnerability-tests: build
	./scripts/e2e-vulnerability.sh

# install kind
install-kind:
	./scripts/install-kind.sh

# build terrascan docker image
docker-build:
	./scripts/docker-build.sh

# build and push latest terrascan docker image
docker-build-push-latest:
	./scripts/docker-build.sh latest

# build and push release tag terrascan docker image
docker-build-push-latest-tag:
	./scripts/docker-build.sh tag


# push terrascan docker image
docker-push:
	./scripts/docker-push.sh

# push latest terrascan docker image
docker-push-latest:
	./scripts/docker-push-latest.sh

# push release tag terrascan docker image
docker-push-latest-tag:
	./scripts/docker-push-latest-tag.sh

# build terrascan_atlantis docker image
atlantis-docker-build:
	./scripts/atlantis/docker-build.sh

# push terrascan_atlantis docker image
atlantis-docker-push:
	./scripts/atlantis/docker-push.sh

# push latest terrascan_atlantis docker image
atlantis-docker-push-latest:
	./scripts/atlantis/docker-push-latest.sh

# push release tag terrascan_atlantis docker image
atlantis-docker-push-latest-tag:
	./scripts/atlantis/docker-push-latest-tag.sh
