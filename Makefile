GOLANGCI_LINT := $(shell command -v golangci-lint 2> /dev/null)
GYCYCLO := $(shell command -v gocyclo 2> /dev/null)
VERSION := $(shell git describe --always --dirty)
BUILD_DATE := $(shell date)
BUILD_OS := $(shell uname -s)
BUILD_HOST := $(shell uname -n)
BUILD_ARCH := $(shell uname -mp)
BUILD_KERNEL_VERSION := $(shell uname -r)

default: build

build:
	@go build -ldflags="\
	                    -X 'github.com/shreerangdixit/lox/build.version=$(VERSION)' \
	                    -X 'github.com/shreerangdixit/lox/build.date=$(BUILD_DATE)' \
	                    -X 'github.com/shreerangdixit/lox/build.os=$(BUILD_OS)' \
	                    -X 'github.com/shreerangdixit/lox/build.host=$(BUILD_HOST)' \
	                    -X 'github.com/shreerangdixit/lox/build.arch=$(BUILD_ARCH)' \
	                    -X 'github.com/shreerangdixit/lox/build.kernelVersion=$(BUILD_KERNEL_VERSION)' \
	                    " \
	                    .

fmt:
	@go fmt ./...

test:
	@go test -v ./...

lint: lint.deps
# TODO: Fixme
	@gocyclo -over 50 .
	@go vet ./...
	@golangci-lint run

lint.fix: lint.deps
	@golangci-lint run --fix

lint.deps:
ifndef GOLANGCI_LINT
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.44.0
endif

ifndef GYCYCLO
	@go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
endif

.PHONY: build fmt test lint lint.fix lint.deps
