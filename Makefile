# Enterprise Shield Plugin Makefile

.PHONY: all build test clean install lint fmt help release

# Build variables
BINARY_NAME=enterprise-shield
BUILD_DIR=./build
DIST_DIR=./dist
CMD_DIR=./cmd/plugin
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME) -w -s"

# Go variables
GOFILES=$(shell find . -name '*.go' -type f)
GOPACKAGES=$(shell go list ./...)

# Release platforms
PLATFORMS=linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64

all: build

## build: Build the plugin binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)
	@echo "Built: $(BUILD_DIR)/$(BINARY_NAME)"

## test: Run all tests
test:
	@echo "Running tests..."
	go test -v ./...

## test-cover: Run tests with coverage
test-cover:
	@echo "Running tests with coverage..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

## lint: Run linters
lint:
	@echo "Running linters..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed, running go vet..."; \
		go vet ./...; \
	fi

## fmt: Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...
	@if command -v goimports >/dev/null 2>&1; then \
		goimports -w .; \
	fi

## clean: Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html
	go clean

## install: Install the plugin to OpenCode directory
install: build
	@echo "Installing to OpenCode plugins directory..."
	@mkdir -p ~/.opencode/plugins
	cp $(BUILD_DIR)/$(BINARY_NAME) ~/.opencode/plugins/
	@echo "Installed: ~/.opencode/plugins/$(BINARY_NAME)"
	@echo ""
	@echo "Don't forget to copy the config:"
	@echo "  cp config/default.yaml ~/.opencode/config/enterprise-shield.yaml"

## install-config: Install default configuration
install-config:
	@echo "Installing default configuration..."
	@mkdir -p ~/.opencode/config
	cp config/default.yaml ~/.opencode/config/enterprise-shield.yaml
	@echo "Installed: ~/.opencode/config/enterprise-shield.yaml"

## deps: Download dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

## run: Build and run the plugin
run: build
	$(BUILD_DIR)/$(BINARY_NAME) serve

## demo: Run a quick demo
demo: build
	@echo "=== Enterprise Shield Demo ==="
	@echo ""
	@echo "1. Scanning for SSN..."
	$(BUILD_DIR)/$(BINARY_NAME) scan "My SSN is 123-45-6789"
	@echo ""
	@echo "2. Processing a request..."
	$(BUILD_DIR)/$(BINARY_NAME) process demo@example.com "Query ServerDB01.users_prod from 192.168.1.100" openai
	@echo ""
	@echo "=== Demo Complete ==="

## release: Build release binaries for all platforms
release:
	@echo "Building release binaries for version: $(VERSION)"
	@mkdir -p $(DIST_DIR)
	@./scripts/build-release.sh $(VERSION)

## release-test: Test the release binaries
release-test: release
	@echo "Testing release binaries..."
	@for file in $(DIST_DIR)/*.tar.gz; do \
		echo "Testing $$file..."; \
		tar -tzf "$$file" > /dev/null && echo "  ✓ Archive OK" || echo "  ✗ Archive FAILED"; \
	done

## docker-build: Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t yourorg/enterprise-shield:$(VERSION) .
	docker tag yourorg/enterprise-shield:$(VERSION) yourorg/enterprise-shield:latest

## docker-push: Push Docker image to registry
docker-push: docker-build
	docker push yourorg/enterprise-shield:$(VERSION)
	docker push yourorg/enterprise-shield:latest

## tag: Create a new version tag
tag:
	@if [ -z "$(V)" ]; then \
		echo "Usage: make tag V=1.0.0"; \
		exit 1; \
	fi
	@echo "Creating tag v$(V)..."
	git tag -a "v$(V)" -m "Release v$(V)"
	@echo "Tag created. Push with: git push origin v$(V)"

## changelog: Generate changelog from git commits
changelog:
	@echo "Generating changelog..."
	@git log --pretty=format:"- %s (%h)" $(shell git describe --tags --abbrev=0)..HEAD

## help: Show this help
help:
	@echo "Enterprise Shield Plugin - Makefile Commands"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/  /'
	@echo ""
	@echo "Examples:"
	@echo "  make build            # Build the plugin"
	@echo "  make test             # Run tests"
	@echo "  make install          # Install to ~/.opencode/plugins/"
	@echo "  make release          # Build release binaries"
	@echo "  make tag V=1.0.0      # Create version tag"
	@echo "  make demo             # Run a quick demo"

