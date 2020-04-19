SHELL := /bin/bash

GO =                 GOFLAGS=-mod=vendor GOPRIVATE=github.com/kanbara go
BUILD_DIR            = bin
BINARY               = $(BUILD_DIR)/lis
PKGFILES               = $(shell find pkg -name '*.go' -type f)
PACKAGES              := $(shell $(GO) list ./...| grep -v node_modules)
DEPENDENCIES          := $(shell find ./vendor -type f)
SYSTEM                := $(shell uname -s | tr A-Z a-z)_$(shell uname -m | sed "s/x86_64/amd64/")
GOOS				?= darwin
PROGAM               = lis
GOLANGCI_LINT_VERSION = 1.24.0

$(BUILD_DIR)/%: %/*.go $(PKGFILES) $(DEPENDENCIES)
	@echo Building $@
	env GOOS=$(GOOS) $(GO) build -ldflags="-s -w -X main.BuildTime=$(BUILD_TIME)" -o $@ ./$(notdir $@)

.PHONY: build
build: $(BINARY)

.PHONY: lint
lint: bin/golangci-lint
	./bin/golangci-lint --modules-download-mode vendor run $(LINT_TARGETS)

.PHONY: test
test:
	$(GO) test -cover -race $(PACKAGES)

bin/golangci-lint:
	mkdir -p bin
	curl -sSLf \
		https://github.com/golangci/golangci-lint/releases/download/v$(GOLANGCI_LINT_VERSION)/golangci-lint-$(GOLANGCI_LINT_VERSION)-$(shell echo $(SYSTEM) | tr '_' '-').tar.gz \
		| tar xzOf - golangci-lint-$(GOLANGCI_LINT_VERSION)-$(shell echo $(SYSTEM) | tr '_' '-')/golangci-lint > bin/golangci-lint && chmod +x bin/golangci-lint