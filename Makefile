GOLANGCI_LINT := $(shell command -v golangci-lint 2> /dev/null)

build:
	@go build .

fmt:
	@go fmt ./...

test:
	@go test -v ./...

lint: lint.deps
	@golangci-lint run

lint.fix: lint.deps
	@golangci-lint run --fix

lint.deps:
ifndef GOLANGCI_LINT
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.44.0
endif

.PHONY: build fmt test lint lint.fix lint.deps
