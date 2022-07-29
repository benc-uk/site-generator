# Common variables
VERSION ?= 1.0.1
BUILD_INFO ?= Manual build

.PHONY: help build run lint lint-fix
.DEFAULT_GOAL := help

help: ## 💬 This help message :)
	@figlet $@ || true
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

lint: ## 🌟 Lint & format, will not fix but sets exit code on error
	@figlet $@ || true
	golangci-lint run ./...

lint-fix: ## 🔍 Lint & format, will try to fix errors and modify code
	@figlet $@ || true
	golangci-lint run ./... --fix

build: ## 🔨 Run a local build, placing binary in ./bin/site-generator
	@figlet $@ || true
	go build -o bin/site-generator -ldflags "-X main.version=$(VERSION) -X 'main.buildInfo=$(BUILD_INFO)'" ./...

run: ## 🚀 Run application, using sample folder as input
	@figlet $@ || true
	go run main.go -s sample