# GitFlow TUI Makefile

# Variables
BINARY_NAME=gitflow-tui
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS=-ldflags "-s -w -X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.Commit=$(COMMIT)"

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Directories
BUILD_DIR=build
DIST_DIR=dist
CMD_DIR=cmd/gitflow-tui

# Platforms
PLATFORMS=darwin/amd64 darwin/arm64 linux/amd64 linux/arm64 windows/amd64

# Colors
GREEN=\033[0;32m
BLUE=\033[0;34m
ORANGE=\033[0;33m
NC=\033[0m # No Color

.PHONY: all build clean test lint install uninstall deps build-all package release

all: deps build

## Build the binary for current platform
build:
	@echo "$(BLUE)Building $(BINARY_NAME)...$(NC)"
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./$(CMD_DIR)
	@echo "$(GREEN)✓ Built: $(BUILD_DIR)/$(BINARY_NAME)$(NC)"

## Build for all platforms
build-all:
	@echo "$(BLUE)Building for all platforms...$(NC)"
	@mkdir -p $(DIST_DIR)
	@for platform in $(PLATFORMS); do \
		GOOS=$$(echo $$platform | cut -d/ -f1); \
		GOARCH=$$(echo $$platform | cut -d/ -f2); \
		OUTPUT=$(DIST_DIR)/$(BINARY_NAME)-$$GOOS-$$GOARCH; \
		if [ "$$GOOS" = "windows" ]; then OUTPUT="$$OUTPUT.exe"; fi; \
		echo "Building for $$GOOS/$$GOARCH..."; \
		GOOS=$$GOOS GOARCH=$$GOARCH $(GOBUILD) $(LDFLAGS) -o $$OUTPUT ./$(CMD_DIR); \
	done
	@echo "$(GREEN)✓ Built all platforms$(NC)"

## Clean build artifacts
clean:
	@echo "$(ORANGE)Cleaning...$(NC)"
	$(GOCLEAN)
	@rm -rf $(BUILD_DIR) $(DIST_DIR)
	@echo "$(GREEN)✓ Cleaned$(NC)"

## Run tests
test:
	@echo "$(BLUE)Running tests...$(NC)"
	$(GOTEST) -v ./...
	@echo "$(GREEN)✓ Tests passed$(NC)"

## Run tests with coverage
test-coverage:
	@echo "$(BLUE)Running tests with coverage...$(NC)"
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)✓ Coverage report generated: coverage.html$(NC)"

## Run linter
lint:
	@echo "$(BLUE)Running linter...$(NC)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "$(ORANGE)golangci-lint not installed, using go vet$(NC)"; \
		go vet ./...; \
	fi
	@echo "$(GREEN)✓ Linting passed$(NC)"

## Install dependencies
deps:
	@echo "$(BLUE)Installing dependencies...$(NC)"
	$(GOMOD) download
	$(GOMOD) tidy
	@echo "$(GREEN)✓ Dependencies installed$(NC)"

## Install binary to system
install: build
	@echo "$(BLUE)Installing $(BINARY_NAME)...$(NC)"
	@install -d $(DESTDIR)/usr/local/bin
	@install -m 755 $(BUILD_DIR)/$(BINARY_NAME) $(DESTDIR)/usr/local/bin/
	@echo "$(GREEN)✓ Installed to /usr/local/bin/$(BINARY_NAME)$(NC)"

## Uninstall binary
uninstall:
	@echo "$(ORANGE)Uninstalling $(BINARY_NAME)...$(NC)"
	@rm -f /usr/local/bin/$(BINARY_NAME)
	@echo "$(GREEN)✓ Uninstalled$(NC)"

## Package for distribution
package: build-all
	@echo "$(BLUE)Packaging...$(NC)"
	@mkdir -p $(DIST_DIR)/packages
	@for platform in $(PLATFORMS); do \
		GOOS=$$(echo $$platform | cut -d/ -f1); \
		GOARCH=$$(echo $$platform | cut -d/ -f2); \
		BINARY=$(BINARY_NAME)-$$GOOS-$$GOARCH; \
		if [ "$$GOOS" = "windows" ]; then BINARY="$$BINARY.exe"; fi; \
		ARCHIVE=$(DIST_DIR)/packages/$(BINARY_NAME)-$(VERSION)-$$GOOS-$$GOARCH.tar.gz; \
		if [ "$$GOOS" = "windows" ]; then ARCHIVE="$${ARCHIVE%.tar.gz}.zip"; fi; \
		echo "Packaging $$BINARY..."; \
		if [ "$$GOOS" = "windows" ]; then \
			cd $(DIST_DIR) && zip -q ../$$ARCHIVE $$BINARY && cd ..; \
		else \
			tar -czf $$ARCHIVE -C $(DIST_DIR) $$BINARY; \
		fi; \
	done
	@echo "$(GREEN)✓ Packaged$(NC)"

## Create release
release: clean test lint package
	@echo "$(GREEN)✓ Release $(VERSION) ready$(NC)"

## Run the application
run: build
	@echo "$(BLUE)Running $(BINARY_NAME)...$(NC)"
	@./$(BUILD_DIR)/$(BINARY_NAME)

## Run with debug mode
debug: build
	@echo "$(BLUE)Running in debug mode...$(NC)"
	@DEBUG=1 ./$(BUILD_DIR)/$(BINARY_NAME)

## Show help
help:
	@echo "GitFlow TUI Makefile"
	@echo ""
	@echo "Usage:"
	@echo "  make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(BLUE)%-15s$(NC) %s\n", $$1, $$2}'

## Development helpers

fmt:
	@echo "$(BLUE)Formatting code...$(NC)"
	@gofmt -w .
	@echo "$(GREEN)✓ Formatted$(NC)"

vet:
	@echo "$(BLUE)Running go vet...$(NC)"
	@go vet ./...
	@echo "$(GREEN)✓ Vet passed$(NC)"

check: fmt vet lint test
	@echo "$(GREEN)✓ All checks passed$(NC)"

## Install development tools
install-tools:
	@echo "$(BLUE)Installing development tools...$(NC)"
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "$(GREEN)✓ Tools installed$(NC)"

## Generate assets (banners, etc.)
generate-assets:
	@echo "$(BLUE)Generating assets...$(NC)"
	@mkdir -p assets
	@echo "$(GREEN)✓ Assets generated$(NC)"

## Build Neovim plugin
build-nvim:
	@echo "$(BLUE)Building Neovim plugin...$(NC)"
	@echo "$(GREEN)✓ Neovim plugin ready$(NC)"

## Build VSCode extension
build-vscode:
	@echo "$(BLUE)Building VSCode extension...$(NC)"
	@cd editors/vscode && npm install && npm run compile
	@echo "$(GREEN)✓ VSCode extension built$(NC)"

## Package VSCode extension
package-vscode: build-vscode
	@echo "$(BLUE)Packaging VSCode extension...$(NC)"
	@cd editors/vscode && npx vsce package
	@echo "$(GREEN)✓ VSCode extension packaged$(NC)"

## Install Neovim plugin locally
install-nvim:
	@echo "$(BLUE)Installing Neovim plugin...$(NC)"
	@mkdir -p ~/.config/nvim/lua/gitflow
	@cp -r editors/nvim/lua/gitflow/* ~/.config/nvim/lua/gitflow/
	@echo "$(GREEN)✓ Neovim plugin installed$(NC)"

## Install VSCode extension locally
install-vscode: package-vscode
	@echo "$(BLUE)Installing VSCode extension...$(NC)"
	@code --install-extension editors/vscode/gitflow-tui-*.vsix
	@echo "$(GREEN)✓ VSCode extension installed$(NC)"

.DEFAULT_GOAL := help
