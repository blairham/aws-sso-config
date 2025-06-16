# Project configuration
BINARY_NAME=aws-sso-config
MAIN_PACKAGE=.
BUILD_DIR=bin
DIST_DIR=dist
COVERAGE_DIR=coverage

# Build information
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT?=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME?=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Go configuration (for development builds only)
GOOS?=$(shell go env GOOS)
GOARCH?=$(shell go env GOARCH)
CGO_ENABLED?=0

# Development build flags (simpler for faster iteration)
DEV_LDFLAGS=-ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.buildTime=$(BUILD_TIME)"

# Default target
.DEFAULT_GOAL := help

# Phony targets
.PHONY: help clean deps tidy fmt vet lint test test-race test-coverage \
        build build-dev build-local install run dev check \
        release snapshot docker docker-only docker-clean pre-commit ci security goreleaser-check \
        tag tag-unsigned tag-lightweight tag-check tag-delete tag-verify tag-list tag tag-check tag-delete

## help: Show this help message
help:
	@echo "Available targets:"
	@awk '/^##/ { \
		sub(/^## /, "", $$0); \
		split($$0, arr, ": "); \
		printf "  %-15s %s\n", arr[1], arr[2] \
	}' $(MAKEFILE_LIST)

## clean: Remove build artifacts and temporary files
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR) $(DIST_DIR) $(COVERAGE_DIR)
	@rm -rf build/
	@rm -f $(BINARY_NAME)
	@rm -f coverage.out profile.out cpu.prof mem.prof block.prof mutex.prof
	@rm -f *.test *.cover *.log *.tmp
	@rm -f aws-config.yaml aws-config.json aws-config.toml
	@find . -name "*.test" -type f -delete 2>/dev/null || true
	@find . -name "*.prof" -type f -delete 2>/dev/null || true
	@find . -name "onfig=" -type d -exec rm -rf {} + 2>/dev/null || true
	@find . -name "*onfig=" -type d -exec rm -rf {} + 2>/dev/null || true
	@find . -name "*Test*" -type d -path "*/command/*" -exec rm -rf {} + 2>/dev/null || true
	@go clean -cache -testcache -modcache
	@echo "Clean complete"

## deps: Download and verify dependencies
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod verify

## tidy: Clean up dependencies
tidy:
	@echo "Tidying dependencies..."
	@go mod tidy

## fmt: Format Go code
fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@gofumpt -l -w . 2>/dev/null || true

## vet: Run go vet
vet:
	@echo "Running go vet..."
	@go vet ./...

## lint: Run golangci-lint
lint:
	@echo "Running golangci-lint..."
	@golangci-lint run

## test: Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

## test-race: Run tests with race detection
test-race:
	@echo "Running tests with race detection..."
	@go test -v -race ./...

## test-coverage: Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@mkdir -p $(COVERAGE_DIR)
	@go test -v -coverprofile=$(COVERAGE_DIR)/coverage.out ./...
	@go tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	@echo "Coverage report: $(COVERAGE_DIR)/coverage.html"

## goreleaser-check: Check GoReleaser configuration
goreleaser-check:
	@echo "Checking GoReleaser configuration..."
	@goreleaser check

## build: Build using GoReleaser (recommended for production)
build: clean goreleaser-check
	@echo "Building with GoReleaser..."
	@goreleaser build --clean --snapshot --single-target
	@echo "Build complete, artifacts in $(DIST_DIR)/"

## build-dev: Quick development build using go build
build-dev:
	@echo "Building $(BINARY_NAME) for development ($(GOOS)/$(GOARCH))..."
	@mkdir -p $(BUILD_DIR)
	@CGO_ENABLED=$(CGO_ENABLED) go build $(DEV_LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PACKAGE)
	@echo "Development build complete: $(BUILD_DIR)/$(BINARY_NAME)"

## build-local: Alias for build-dev (backward compatibility)
build-local: build-dev

## build-all: Build for all platforms using GoReleaser
build-all: clean goreleaser-check
	@echo "Building for all platforms with GoReleaser..."
	@goreleaser build --clean --snapshot
	@echo "Multi-platform build complete, artifacts in $(DIST_DIR)/"

## go-install: Install with go install and version info (recommended)
go-install:
	@echo "Installing $(BINARY_NAME) with go install..."
	@echo "This will install from the latest tagged version or main branch"
	@go install github.com/blairham/aws-sso-config@latest

## go-install-dev: Install development version with local changes
go-install-dev:
	@echo "Installing $(BINARY_NAME) from local source with version info..."
	@go install $(DEV_LDFLAGS) $(MAIN_PACKAGE)
	@echo "Development installation complete!"
	@if [ -n "$(shell go env GOBIN)" ]; then \
		echo "Installed to: $(shell go env GOBIN)/$(BINARY_NAME)"; \
		echo "Make sure $(shell go env GOBIN) is in your PATH"; \
	else \
		echo "Installed to: $(shell go env GOPATH)/bin/$(BINARY_NAME)"; \
		echo "Make sure $(shell go env GOPATH)/bin is in your PATH"; \
	fi

## check-install: Verify installation and show version
check-install:
	@echo "Checking $(BINARY_NAME) installation..."
	@if command -v $(BINARY_NAME) >/dev/null 2>&1; then \
		echo "✓ $(BINARY_NAME) is installed"; \
		$(BINARY_NAME) --version; \
	else \
		echo "✗ $(BINARY_NAME) is not found in PATH"; \
		if [ -n "$(shell go env GOBIN)" ]; then \
			echo "Make sure $(shell go env GOBIN) is in your PATH"; \
			echo "Binary should be at: $(shell go env GOBIN)/$(BINARY_NAME)"; \
			if [ -f "$(shell go env GOBIN)/$(BINARY_NAME)" ]; then \
				echo "✓ Binary exists at expected location"; \
			else \
				echo "✗ Binary not found at expected location"; \
			fi; \
		else \
			echo "Make sure $(shell go env GOPATH)/bin is in your PATH"; \
			echo "Binary should be at: $(shell go env GOPATH)/bin/$(BINARY_NAME)"; \
			if [ -f "$(shell go env GOPATH)/bin/$(BINARY_NAME)" ]; then \
				echo "✓ Binary exists at expected location"; \
			else \
				echo "✗ Binary not found at expected location"; \
			fi; \
		fi; \
		echo "Current PATH: $$PATH"; \
	fi

## install: Install the binary to GOPATH/bin
install:
	@echo "Installing $(BINARY_NAME)..."
	@go install $(DEV_LDFLAGS) $(MAIN_PACKAGE)

## run: Build and run the application (development build)
run: build-dev
	@echo "Running $(BINARY_NAME)..."
	@$(BUILD_DIR)/$(BINARY_NAME)

## dev: Run the application in development mode (no build)
dev:
	@echo "Running in development mode..."
	@go run $(DEV_LDFLAGS) $(MAIN_PACKAGE)

## check: Run all checks (fmt, vet, lint, test)
check: fmt vet lint test
	@echo "All checks passed!"

## pre-commit: Run pre-commit checks
pre-commit: tidy fmt vet lint test-race goreleaser-check
	@echo "Pre-commit checks complete!"

## ci: Run GitHub Actions CI pipeline locally using act
ci:
	@echo "Running GitHub Actions CI pipeline locally..."
	@echo "Note: This requires Docker to be running and 'act' to be installed"
	@act --container-architecture linux/amd64 push

## ci-job: Run a specific CI job locally (usage: make ci-job JOB=test)
ci-job:
	@echo "Running CI job: $(or $(JOB),test)"
	@act --container-architecture linux/amd64 --job $(or $(JOB),test) push

## ci-local: Run basic CI checks without Docker (legacy)
ci-local: deps check test-coverage build
	@echo "Local CI pipeline complete!"

## security: Run security checks
security:
	@echo "Running security checks..."
	@govulncheck ./... 2>/dev/null || echo "govulncheck not installed, skipping vulnerability check"

## release: Create a release using goreleaser
release: check
	@echo "Creating release..."
	@goreleaser release --clean

## snapshot: Create a snapshot build using goreleaser
snapshot: check
	@echo "Creating snapshot build..."
	@goreleaser build --snapshot --clean

## docker: Build Docker image using GoReleaser (if Dockerfile exists)
docker:
	@if [ -f Dockerfile ]; then \
		echo "Building Docker image with GoReleaser..."; \
		goreleaser build --single-target --snapshot --clean && \
		goreleaser release --snapshot --clean --skip=publish; \
	else \
		echo "No Dockerfile found, skipping Docker build"; \
	fi

## docker-build: Build Docker image with version info using GoReleaser
docker-build:
	@echo "Building Docker image $(BINARY_NAME):$(VERSION) with GoReleaser..."
	@goreleaser build --single-target --snapshot --clean
	@goreleaser release --snapshot --clean --skip=publish
	@echo "Docker image built successfully!"
	@docker images $(BINARY_NAME)

## docker-run: Run the Docker container interactively
docker-run:
	@echo "Running $(BINARY_NAME) in Docker container..."
	@docker run --rm -it \
		-v $(HOME)/.aws:/home/appuser/.aws:ro \
		-v $(PWD):/workspace \
		$(BINARY_NAME):$(VERSION) $(ARGS)

## docker-run-config: Run with config directory mounted
docker-run-config:
	@echo "Running $(BINARY_NAME) with config directory mounted..."
	@docker run --rm -it \
		-v $(HOME)/.aws:/home/appuser/.aws:ro \
		-v $(HOME)/.awsssoconfig:/home/appuser/.awsssoconfig \
		-v $(PWD):/workspace \
		$(BINARY_NAME):$(VERSION) $(ARGS)

## docker-shell: Run an interactive shell in the container (debug mode)
docker-shell:
	@echo "Starting debug shell in $(BINARY_NAME) container..."
	@docker run --rm -it \
		-v $(HOME)/.aws:/home/appuser/.aws:ro \
		-v $(PWD):/workspace \
		--entrypoint="" \
		--user root \
		alpine:latest \
		sh -c "apk add --no-cache curl && exec sh"

## docker-size: Show Docker image size
docker-size:
	@echo "Docker image sizes:"
	@docker images $(BINARY_NAME) --format "table {{.Repository}}\t{{.Tag}}\t{{.Size}}"

## docker-inspect: Inspect the Docker image
docker-inspect:
	@echo "Inspecting Docker image $(BINARY_NAME):$(VERSION)..."
	@docker inspect $(BINARY_NAME):$(VERSION)

## docker-test: Test the Docker image
docker-test: docker-build
	@echo "Testing Docker image..."
	@docker run --rm $(BINARY_NAME):$(VERSION) --version
	@docker run --rm $(BINARY_NAME):$(VERSION) --help
	@echo "Docker image test complete!"

## docker-only: Build only Docker images with GoReleaser (no binaries)
docker-only:
	@echo "Building Docker images only with GoReleaser..."
	@goreleaser release --snapshot --clean --skip=publish --skip=validate

## docker-clean: Remove Docker images created by this project
docker-clean:
	@echo "Cleaning up Docker images for $(BINARY_NAME)..."
	@if docker images -q $(BINARY_NAME) 2>/dev/null | grep -q .; then \
		echo "Removing $(BINARY_NAME) images..."; \
		docker rmi $$(docker images -q $(BINARY_NAME)) --force; \
		echo "Docker images for $(BINARY_NAME) removed successfully!"; \
	else \
		echo "No $(BINARY_NAME) images found to remove."; \
	fi
	@echo "Cleaning up dangling images..."
	@if docker images -f "dangling=true" -q 2>/dev/null | grep -q .; then \
		docker rmi $$(docker images -f "dangling=true" -q) --force; \
		echo "Dangling images removed."; \
	else \
		echo "No dangling images found."; \
	fi

## tag: Create and push a signed annotated git tag (use: make tag VERSION=v1.0.0)
tag:
	@if [ -z "$(VERSION)" ]; then \
		echo "Error: VERSION is required. Use: make tag VERSION=v1.0.0"; \
		exit 1; \
	fi
	@echo "Creating signed annotated tag $(VERSION)..."
	@echo "Using GPG key: $$(git config --get user.signingkey)"
	@git tag -a $(VERSION) -m "Release $(VERSION)"
	@git push origin $(VERSION)
	@echo "✅ Signed tag $(VERSION) created and pushed successfully!"
	@echo "🚀 Release workflow will start automatically."
	@echo "📊 View progress: https://github.com/blairham/aws-sso-config/actions"
	@echo "🔍 Verify signature: make tag-verify VERSION=$(VERSION)"

## tag-unsigned: Create and push an unsigned annotated tag (fallback)
tag-unsigned:
	@if [ -z "$(VERSION)" ]; then \
		echo "Error: VERSION is required. Use: make tag-unsigned VERSION=v1.0.0"; \
		exit 1; \
	fi
	@echo "Creating annotated tag $(VERSION) (unsigned)..."
	@git tag -a $(VERSION) -m "Release $(VERSION)"
	@git push origin $(VERSION)
	@echo "Tag $(VERSION) created and pushed. Release workflow will start automatically."
	@echo "View release progress at: https://github.com/blairham/aws-sso-config/actions"

## tag-lightweight: Create and push a lightweight tag (not recommended for releases)
tag-lightweight:
	@if [ -z "$(VERSION)" ]; then \
		echo "Error: VERSION is required. Use: make tag-lightweight VERSION=v1.0.0"; \
		exit 1; \
	fi
	@echo "Creating lightweight tag $(VERSION)..."
	@git tag $(VERSION)
	@git push origin $(VERSION)
	@echo "Lightweight tag $(VERSION) created and pushed."

## tag-check: Check if current commit is ready for tagging
tag-check:
	@echo "Checking repository status for tagging..."
	@git status --porcelain | grep -q . && echo "Error: Repository has uncommitted changes" && exit 1 || echo "✓ Repository is clean"
	@git diff-index --quiet HEAD -- && echo "✓ No unstaged changes" || (echo "Error: Repository has unstaged changes" && exit 1)
	@echo "✓ Repository is ready for tagging"

## tag-delete: Delete a tag locally and remotely (use: make tag-delete VERSION=v1.0.0)
tag-delete:
	@if [ -z "$(VERSION)" ]; then \
		echo "Error: VERSION is required. Use: make tag-delete VERSION=v1.0.0"; \
		exit 1; \
	fi
	@echo "Deleting tag $(VERSION)..."
	@git tag -d $(VERSION) || true
	@git push origin :refs/tags/$(VERSION) || true
	@echo "Tag $(VERSION) deleted"

## tag-verify: Verify the signature of a tag (use: make tag-verify VERSION=v1.0.0)
tag-verify:
	@if [ -z "$(VERSION)" ]; then \
		echo "Error: VERSION is required. Use: make tag-verify VERSION=v1.0.0"; \
		exit 1; \
	fi
	@echo "Verifying signature of tag $(VERSION)..."
	@git tag -v $(VERSION)

## tag-list: List all tags with verification status
tag-list:
	@echo "Available tags:"
	@git tag -l --sort=-version:refname | head -10

# Version information
version:
	@echo "Version: $(VERSION)"
	@echo "Commit: $(COMMIT)"
	@echo "Build Time: $(BUILD_TIME)"
