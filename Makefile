# A Self-Documenting Makefile: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html

GO_SOURCE_FILES = $(shell find . -type f -name "*.go" -not -name "bindata.go" -not -path "./vendor/*")
GO_PACKAGES = $(shell go list ./... | grep -v /vendor/)

.PHONY: install check test fmt csfix envcheck help
.DEFAULT_GOAL := help

install: ## Install dependencies
	@glide install

check: test fmt ## Run tests and linters

test: ## Run unit tests
	@go test ${GO_PACKAGES}

watch-test: ## Watch for file changes and run tests
	reflex -t 2s -d none -r '\.go$$' -- go test ${GO_PACKAGES}

fmt: ## Check that all source files follow the Coding Style
	@gofmt -l ${GO_SOURCE_FILES} | read something && echo "Code differs from gofmt's style" 1>&2 && exit 1 || true

csfix: ## Fix Coding Standard violations
	@gofmt -l -w -s ${GO_SOURCE_FILES}

envcheck: ## Check environment for all the necessary requirements
	$(call executable_check,Go,go)
	$(call executable_check,Glide,glide)
	$(call executable_check,Reflex,reflex)

define executable_check
    @printf "\033[36m%-30s\033[0m %s\n" "$(1)" `if which $(2) > /dev/null 2>&1; then echo "\033[0;32m✓\033[0m"; else echo "\033[0;31m✗\033[0m"; fi`
endef

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
