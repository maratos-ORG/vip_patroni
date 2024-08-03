.DEFAULT_GOAL := help

.PHONY: help
help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

ifeq ($(GOPATH),)
GOPATH := ~/go
endif

PROJECT_NAME = cmd
RELEASE ?= dev
BINARY_NAME ?= bin/vip_patroni
BUILD_TIME ?= $(shell date '+%Y-%m-%d_%H:%M:%S')
COMMIT_HASH = $(shell git rev-parse --short HEAD)

.PHONY: all
all: linter test build ## Run linter, tests and build a package

test: ## Run tests
	go test -race -coverprofile=coverage.out ./$(PROJECT_NAME)/...

linter: ## Apply linter
	golangci-lint run -c ./.golangci.yml --timeout 3m ./internal/...

clean: ## Clean before build
	@go clean ./...

build: clean ## Build package
	@go build -ldflags "-s -w -X 'main.build=${RELEASE}'" -o ${BINARY_NAME} ./$(PROJECT_NAME)/...