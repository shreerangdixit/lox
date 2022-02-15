GOLANGCI_LINT := $(shell command -v golangci-lint 2> /dev/null)
GYCYCLO := $(shell command -v gocyclo 2> /dev/null)

default: build

build:
	@go build .

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
