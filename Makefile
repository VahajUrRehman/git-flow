# GitFlow TUI - Universal Build Makefile
# Builds for all platforms: Windows, macOS, Linux (including WSL)

# Binary name
BINARY_NAME=gitflow-tui

# Version info (auto-detect or use 'dev')
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Build flags
LDFLAGS=-ldflags "-s -w -X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.Commit=$(COMMIT)"

# Directories
BUILD_DIR=build
DIST_DIR=dist
RELEASE_DIR=release
SCRIPTS_DIR=scripts

# Platforms to build for
PLATFORMS=windows/amd64 windows/arm64 darwin/amd64 darwin/arm64 linux/amd64 linux/arm64

# Colors for output
GREEN=\033[0;32m
BLUE=\033[0;34m
ORANGE=\033[0;33m
RED=\033[0;31m
NC=\033[0m

.PHONY: all build build-all build-windows build-macos build-linux build-wsl clean test lint fmt install uninstall \
        package package-all package-windows package-macos package-linux \
        release release-check release-local \
        install-scripts update-formulas checksums

# Default target
all: deps build-all package-all checksums

## Dependencies
deps:
	@echo "$(BLUE)Installing dependencies...$(NC)"
	go mod download
	go mod tidy
	@echo "$(GREEN)✓ Dependencies ready$(NC)"

## Build for current platform
build:
	@echo "$(BLUE)Building for current platform...$(NC)"
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/gitflow-tui
	@echo "$(GREEN)✓ Built: $(BUILD_DIR)/$(BINARY_NAME)$(NC)"

## Build for all platforms
build-all: build-windows build-macos build-linux
	@echo "$(GREEN)✓ All platforms built$(NC)"

## Build for Windows
build-windows:
	@echo "$(BLUE)Building for Windows...$(NC)"
	@mkdir -p $(DIST_DIR)/windows
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(DIST_DIR)/windows/$(BINARY_NAME)-windows-amd64.exe ./cmd/gitflow-tui
	GOOS=windows GOARCH=arm64 go build $(LDFLAGS) -o $(DIST_DIR)/windows/$(BINARY_NAME)-windows-arm64.exe ./cmd/gitflow-tui
	@echo "$(GREEN)✓ Windows builds complete$(NC)"

## Build for macOS
build-macos:
	@echo "$(BLUE)Building for macOS...$(NC)"
	@mkdir -p $(DIST_DIR)/macos
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(DIST_DIR)/macos/$(BINARY_NAME)-darwin-amd64 ./cmd/gitflow-tui
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(DIST_DIR)/macos/$(BINARY_NAME)-darwin-arm64 ./cmd/gitflow-tui
	@echo "$(GREEN)✓ macOS builds complete$(NC)"

## Build for Linux
build-linux:
	@echo "$(BLUE)Building for Linux...$(NC)"
	@mkdir -p $(DIST_DIR)/linux
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(DIST_DIR)/linux/$(BINARY_NAME)-linux-amd64 ./cmd/gitflow-tui
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(DIST_DIR)/linux/$(BINARY_NAME)-linux-arm64 ./cmd/gitflow-tui
	@echo "$(GREEN)✓ Linux builds complete$(NC)"

## Build specifically for WSL (Linux AMD64)
build-wsl: build-linux
	@echo "$(BLUE)WSL build ready at: $(DIST_DIR)/linux/$(BINARY_NAME)-linux-amd64$(NC)"
	@echo "$(GREEN)✓ Copy to WSL: cp $(DIST_DIR)/linux/$(BINARY_NAME)-linux-amd64 /usr/local/bin/$(BINARY_NAME)$(NC)"

## Create all packages
package-all: package-windows package-macos package-linux
	@echo "$(GREEN)✓ All packages created$(NC)"

## Package Windows builds
package-windows: build-windows
	@echo "$(BLUE)Packaging Windows builds...$(NC)"
	@mkdir -p $(RELEASE_DIR)
	cd $(DIST_DIR)/windows && \
		zip -q ../../$(RELEASE_DIR)/$(BINARY_NAME)-$(VERSION)-windows-amd64.zip $(BINARY_NAME)-windows-amd64.exe && \
		zip -q ../../$(RELEASE_DIR)/$(BINARY_NAME)-$(VERSION)-windows-arm64.zip $(BINARY_NAME)-windows-arm64.exe
	@echo "$(GREEN)✓ Windows packages created$(NC)"

## Package macOS builds
package-macos: build-macos
	@echo "$(BLUE)Packaging macOS builds...$(NC)"
	@mkdir -p $(RELEASE_DIR)
	cd $(DIST_DIR)/macos && \
		tar -czf ../../$(RELEASE_DIR)/$(BINARY_NAME)-$(VERSION)-darwin-amd64.tar.gz $(BINARY_NAME)-darwin-amd64 && \
		tar -czf ../../$(RELEASE_DIR)/$(BINARY_NAME)-$(VERSION)-darwin-arm64.tar.gz $(BINARY_NAME)-darwin-arm64
	@echo "$(GREEN)✓ macOS packages created$(NC)"

## Package Linux builds
package-linux: build-linux
	@echo "$(BLUE)Packaging Linux builds...$(NC)"
	@mkdir -p $(RELEASE_DIR)
	cd $(DIST_DIR)/linux && \
		tar -czf ../../$(RELEASE_DIR)/$(BINARY_NAME)-$(VERSION)-linux-amd64.tar.gz $(BINARY_NAME)-linux-amd64 && \
		tar -czf ../../$(RELEASE_DIR)/$(BINARY_NAME)-$(VERSION)-linux-arm64.tar.gz $(BINARY_NAME)-linux-arm64
	@echo "$(GREEN)✓ Linux packages created$(NC)"

## Generate SHA256 checksums
checksums: package-all
	@echo "$(BLUE)Generating checksums...$(NC)"
	@cd $(RELEASE_DIR) && sha256sum * > checksums.txt
	@echo "$(GREEN)✓ Checksums saved to $(RELEASE_DIR)/checksums.txt$(NC)"

## Update install scripts with version
update-scripts:
	@echo "$(BLUE)Updating install scripts...$(NC)"
	@sed -i.bak 's/VERSION_PLACEHOLDER/$(VERSION)/g' $(SCRIPTS_DIR)/homebrew-formula.rb && rm $(SCRIPTS_DIR)/homebrew-formula.rb.bak
	@sed -i.bak 's/VERSION_PLACEHOLDER/$(VERSION)/g' $(SCRIPTS_DIR)/scoop-manifest.json && rm $(SCRIPTS_DIR)/scoop-manifest.json.bak
	@echo "$(GREEN)✓ Scripts updated$(NC)"

## Copy install scripts to release
copy-scripts: update-scripts
	@echo "$(BLUE)Copying install scripts...$(NC)"
	@cp $(SCRIPTS_DIR)/homebrew-formula.rb $(RELEASE_DIR)/
	@cp $(SCRIPTS_DIR)/scoop-manifest.json $(RELEASE_DIR)/
	@cp install.sh $(RELEASE_DIR)/
	@echo "$(GREEN)✓ Scripts copied$(NC)"

## Create complete release
release: clean deps build-all package-all checksums copy-scripts
	@echo "$(GREEN)========================================$(NC)"
	@echo "$(GREEN)✓ Release $(VERSION) ready!$(NC)"
	@echo "$(GREEN)========================================$(NC)"
	@echo ""
	@echo "Files in $(RELEASE_DIR)/:"
	@ls -lh $(RELEASE_DIR)/
	@echo ""
	@echo "$(BLUE)Next steps:$(NC)"
	@echo "  1. Test binaries: $(DIST_DIR)/<platform>/"
	@echo "  2. Upload packages from $(RELEASE_DIR)/"
	@echo "  3. Update Homebrew formula with SHA256 hashes"

## Create local release (without git requirements)
release-local: deps build-all package-all
	@echo "$(GREEN)✓ Local release built$(NC)"
	@echo "Files in $(DIST_DIR)/:"
	@ls -lh $(DIST_DIR)/*/

## Run tests
test:
	@echo "$(BLUE)Running tests...$(NC)"
	go test -v ./...
	@echo "$(GREEN)✓ Tests passed$(NC)"

## Run linter
lint:
	@echo "$(BLUE)Running linter...$(NC)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		go vet ./...; \
	fi
	@echo "$(GREEN)✓ Linting passed$(NC)"

## Format code
fmt:
	@echo "$(BLUE)Formatting code...$(NC)"
	gofmt -w .
	@echo "$(GREEN)✓ Code formatted$(NC)"

## Clean build artifacts
clean:
	@echo "$(ORANGE)Cleaning build artifacts...$(NC)"
	@rm -rf $(BUILD_DIR) $(DIST_DIR) $(RELEASE_DIR)
	@go clean
	@echo "$(GREEN)✓ Cleaned$(NC)"

## Install locally (current platform only)
install: build
	@echo "$(BLUE)Installing to /usr/local/bin...$(NC)"
	@cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
	@chmod +x /usr/local/bin/$(BINARY_NAME)
	@echo "$(GREEN)✓ Installed to /usr/local/bin/$(BINARY_NAME)$(NC)"

## Uninstall
uninstall:
	@echo "$(ORANGE)Uninstalling...$(NC)"
	@rm -f /usr/local/bin/$(BINARY_NAME)
	@echo "$(GREEN)✓ Uninstalled$(NC)"

## Install to WSL (from Windows)
install-wsl: build-wsl
	@echo "$(BLUE)Installing to WSL...$(NC)"
	@cp $(DIST_DIR)/linux/$(BINARY_NAME)-linux-amd64 /mnt/c/temp/$(BINARY_NAME)
	@wsl -e sudo cp /mnt/c/temp/$(BINARY_NAME) /usr/local/bin/
	@wsl -e sudo chmod +x /usr/local/bin/$(BINARY_NAME)
	@echo "$(GREEN)✓ Installed to WSL$(NC)"

## Show help
help:
	@echo "GitFlow TUI Build System"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Main Targets:"
	@echo "  $(GREEN)make all$(NC)          - Full build for all platforms + packages"
	@echo "  $(GREEN)make build$(NC)        - Build for current platform only"
	@echo "  $(GREEN)make build-all$(NC)    - Build for all platforms"
	@echo "  $(GREEN)make release$(NC)      - Create complete release with packages"
	@echo ""
	@echo "Platform-Specific:"
	@echo "  $(BLUE)make build-windows$(NC) - Build Windows binaries (amd64, arm64)"
	@echo "  $(BLUE)make build-macos$(NC)   - Build macOS binaries (Intel, ARM)"
	@echo "  $(BLUE)make build-linux$(NC)   - Build Linux binaries (amd64, arm64)"
	@echo "  $(BLUE)make build-wsl$(NC)     - Build for WSL (Linux amd64)"
	@echo ""
	@echo "Packaging:"
	@echo "  $(BLUE)make package-all$(NC)  - Create zip/tar.gz for all platforms"
	@echo "  $(BLUE)make package-windows$(NC) - Create Windows zip files"
	@echo "  $(BLUE)make package-macos$(NC)   - Create macOS tar.gz files"
	@echo "  $(BLUE)make package-linux$(NC)   - Create Linux tar.gz files"
	@echo ""
	@echo "Development:"
	@echo "  $(BLUE)make test$(NC)         - Run tests"
	@echo "  $(BLUE)make lint$(NC)         - Run linter"
	@echo "  $(BLUE)make fmt$(NC)          - Format code"
	@echo "  $(BLUE)make clean$(NC)        - Clean all build artifacts"
	@echo ""
	@echo "Installation:"
	@echo "  $(BLUE)make install$(NC)      - Install locally (current platform)"
	@echo "  $(BLUE)make uninstall$(NC)    - Remove local installation"
	@echo "  $(BLUE)make install-wsl$(NC)  - Install to WSL from Windows"
	@echo ""
	@echo "Current Version: $(ORANGE)$(VERSION)$(NC)"
