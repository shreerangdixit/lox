GOLANGCI_LINT := $(shell command -v golangci-lint 2> /dev/null)
GYCYCLO := $(shell command -v gocyclo 2> /dev/null)
VERSION := $(shell git describe --always --dirty)
BUILD_DATE := $(shell date)
BUILD_OS := $(shell uname -s)
BUILD_HOST := $(shell uname -n)
BUILD_ARCH := $(shell uname -mp)
BUILD_KERNEL_VERSION := $(shell uname -r)
BUILD_FLAGS := "-X 'github.com/shreerangdixit/redes/build.version=$(VERSION)' \
	            -X 'github.com/shreerangdixit/redes/build.date=$(BUILD_DATE)' \
	            -X 'github.com/shreerangdixit/redes/build.os=$(BUILD_OS)' \
	            -X 'github.com/shreerangdixit/redes/build.host=$(BUILD_HOST)' \
	            -X 'github.com/shreerangdixit/redes/build.arch=$(BUILD_ARCH)' \
	            -X 'github.com/shreerangdixit/redes/build.kernelVersion=$(BUILD_KERNEL_VERSION)'"
REDES_TEST_FILES := $(sort $(shell find ./tests -type f -name '*.rds' -print))

default: build

build:
	@go build -ldflags=$(BUILD_FLAGS) .

fmt:
	@go fmt ./...
	@gofmt -s -w .

test:
	@go test -v ./...

test.rds:
	@make
	@for file in $(REDES_TEST_FILES); do \
		set -e ; \
		./redes $$file; \
	done

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

.PHONY: build fmt test test.redes lint lint.fix lint.deps
