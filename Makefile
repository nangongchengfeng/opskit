# OpsKit Makefile

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "v0.1.0-dev")
BUILD_TIME ?= $(shell date -u '+%Y-%m-%d_%H:%M:%S')
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

GO ?= go
GOFLAGS ?= -v
LDFLAGS := -s -w
LDFLAGS += -X main.Version=$(VERSION)
LDFLAGS += -X main.BuildTime=$(BUILD_TIME)
LDFLAGS += -X main.Commit=$(COMMIT)

BINARY_NAME ?= opskit
BUILD_DIR ?= bin
DIST_DIR ?= dist

.PHONY: all
all: build

.PHONY: help
help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) }' $(MAKEFILE_LIST)

##@ Development

.PHONY: download-tools
download-tools: ## Download third-party tool binaries
	chmod +x ./scripts/download-tools.sh
	./scripts/download-tools.sh

.PHONY: build
build: download-tools ## Build opskit for current platform
	CGO_ENABLED=0 $(GO) build $(GOFLAGS) -tags=embed_assets -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/opskit

.PHONY: build-direct
build-direct: ## Build opskit for current platform (skip download-tools)
	CGO_ENABLED=0 $(GO) build $(GOFLAGS) -tags=embed_assets -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/opskit

.PHONY: build-linux-amd64
build-linux-amd64: download-tools ## Build for Linux amd64
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build $(GOFLAGS) -tags=embed_assets -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/opskit

.PHONY: build-linux-amd64-direct
build-linux-amd64-direct: ## Build for Linux amd64 (skip download-tools)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build $(GOFLAGS) -tags=embed_assets -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/opskit

.PHONY: build-linux-arm64
build-linux-arm64: download-tools ## Build for Linux arm64
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(GO) build $(GOFLAGS) -tags=embed_assets -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 ./cmd/opskit

.PHONY: build-linux-arm64-direct
build-linux-arm64-direct: ## Build for Linux arm64 (skip download-tools)
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(GO) build $(GOFLAGS) -tags=embed_assets -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 ./cmd/opskit

.PHONY: build-all
build-all: download-tools ## Build for all Linux platforms
	$(MAKE) build-linux-amd64
	$(MAKE) build-linux-arm64

.PHONY: build-debug
build-debug: download-tools ## Build without stripping (for debugging)
	CGO_ENABLED=0 $(GO) build $(GOFLAGS) -tags=embed_assets -ldflags="-X main.Version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/opskit

.PHONY: run
run: build ## Build and run opskit
	$(BUILD_DIR)/$(BINARY_NAME) --help

.PHONY: test
test: ## Run tests
	$(GO) test -v ./...

.PHONY: test-coverage
test-coverage: ## Run tests with coverage
	$(GO) test -v -coverprofile=coverage.txt ./...
	$(GO) tool cover -html=coverage.txt -o coverage.html

.PHONY: lint
lint: ## Run linter
	if command -v golint &> /dev/null; then golint ./...; fi
	if command -v golangci-lint &> /dev/null; then golangci-lint run; fi

.PHONY: fmt
fmt: ## Format Go code
	$(GO) fmt ./...

.PHONY: vet
vet: ## Run go vet
	$(GO) vet ./...

##@ Release

.PHONY: build-snapshot
build-snapshot: download-tools ## Build snapshot release (requires goreleaser)
	goreleaser build --snapshot --clean

.PHONY: release
release: download-tools ## Create a release (requires goreleaser and GITHUB_TOKEN)
	goreleaser release --clean

.PHONY: release-snapshot
release-snapshot: download-tools ## Create a snapshot release (requires goreleaser)
	goreleaser release --snapshot --skip-publish --clean

##@ Cleanup

.PHONY: clean
clean: ## Remove build artifacts
	rm -rf $(BUILD_DIR) $(DIST_DIR) coverage.txt coverage.html

.PHONY: clean-cache
clean-cache: ## Remove opskit cache directory
	rm -rf ~/.cache/opskit /tmp/.opskit-bin-*

.PHONY: distclean
distclean: clean clean-cache ## Clean everything
	rm -rf assets/linux-amd64/* assets/linux-arm64/*
	touch assets/linux-amd64/.gitkeep assets/linux-arm64/.gitkeep
