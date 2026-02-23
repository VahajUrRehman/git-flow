#!/bin/bash
# GitFlow TUI Installation Script
# Usage: curl -sSL https://raw.githubusercontent.com/gitflow/tui/main/install.sh | bash

set -e

REPO="gitflow/tui"
BINARY_NAME="gitflow-tui"
INSTALL_DIR="/usr/local/bin"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Detect OS and architecture
detect_platform() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)
    
    case "$OS" in
        linux)
            PLATFORM="linux"
            ;;
        darwin)
            PLATFORM="darwin"
            ;;
        mingw*|msys*|cygwin*)
            PLATFORM="windows"
            INSTALL_DIR="$HOME/bin"
            ;;
        *)
            echo -e "${RED}Unsupported OS: $OS${NC}"
            exit 1
            ;;
    esac
    
    case "$ARCH" in
        x86_64|amd64)
            ARCH="amd64"
            ;;
        arm64|aarch64)
            ARCH="arm64"
            ;;
        *)
            echo -e "${RED}Unsupported architecture: $ARCH${NC}"
            exit 1
            ;;
    esac
    
    if [ "$PLATFORM" = "windows" ]; then
        BINARY_NAME="${BINARY_NAME}.exe"
    fi
}

# Get latest release version
get_latest_version() {
    VERSION=$(curl -s "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    if [ -z "$VERSION" ]; then
        echo -e "${RED}Failed to get latest version${NC}"
        exit 1
    fi
    echo -e "${BLUE}Latest version: $VERSION${NC}"
}

# Download and install
download_and_install() {
    DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${BINARY_NAME%-*}-${PLATFORM}-${ARCH}"
    
    if [ "$PLATFORM" = "windows" ]; then
        DOWNLOAD_URL="${DOWNLOAD_URL}.exe"
    fi
    
    TMP_DIR=$(mktemp -d)
    TMP_FILE="$TMP_DIR/$BINARY_NAME"
    
    echo -e "${BLUE}Downloading from $DOWNLOAD_URL...${NC}"
    curl -sSL "$DOWNLOAD_URL" -o "$TMP_FILE"
    
    if [ ! -f "$TMP_FILE" ]; then
        echo -e "${RED}Download failed${NC}"
        exit 1
    fi
    
    chmod +x "$TMP_FILE"
    
    # Create install directory if needed
    if [ ! -d "$INSTALL_DIR" ]; then
        mkdir -p "$INSTALL_DIR"
    fi
    
    # Check if we need sudo
    if [ -w "$INSTALL_DIR" ]; then
        mv "$TMP_FILE" "$INSTALL_DIR/$BINARY_NAME"
    else
        echo -e "${BLUE}Requesting sudo access to install to $INSTALL_DIR...${NC}"
        sudo mv "$TMP_FILE" "$INSTALL_DIR/$BINARY_NAME"
    fi
    
    rm -rf "$TMP_DIR"
    
    echo -e "${GREEN}✓ Installed $BINARY_NAME to $INSTALL_DIR${NC}"
}

# Verify installation
verify() {
    if command -v "$INSTALL_DIR/$BINARY_NAME" &> /dev/null; then
        VERSION_OUTPUT=$($INSTALL_DIR/$BINARY_NAME --version 2>&1 || true)
        echo -e "${GREEN}✓ Installation successful!${NC}"
        echo -e "${BLUE}$VERSION_OUTPUT${NC}"
    else
        echo -e "${RED}Installation may have failed. Please check $INSTALL_DIR is in your PATH.${NC}"
        exit 1
    fi
}

# Main
main() {
    echo -e "${BLUE}Installing GitFlow TUI...${NC}"
    
    detect_platform
    get_latest_version
    download_and_install
    verify
    
    echo ""
    echo -e "${GREEN}Installation complete! Run 'gitflow-tui' to start.${NC}"
}

main "$@"
