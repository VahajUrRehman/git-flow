#!/bin/bash

# GitFlow TUI Installation Script
# Usage: curl -sSL https://raw.githubusercontent.com/gitflow/tui/main/install.sh | bash

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
REPO="gitflow/tui"
BINARY_NAME="gitflow-tui"
INSTALL_DIR="/usr/local/bin"

# Detect OS and architecture
detect_platform() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)
    
    case "$ARCH" in
        x86_64)
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
    
    case "$OS" in
        linux|darwin)
            ;;
        mingw*|msys*|cygwin*)
            OS="windows"
            BINARY_NAME="${BINARY_NAME}.exe"
            ;;
        *)
            echo -e "${RED}Unsupported operating system: $OS${NC}"
            exit 1
            ;;
    esac
    
    PLATFORM="${OS}-${ARCH}"
}

# Get latest release version
get_latest_version() {
    echo -e "${BLUE}Fetching latest version...${NC}"
    
    if command -v curl &> /dev/null; then
        VERSION=$(curl -s "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    elif command -v wget &> /dev/null; then
        VERSION=$(wget -qO- "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    else
        echo -e "${RED}curl or wget is required${NC}"
        exit 1
    fi
    
    if [ -z "$VERSION" ]; then
        echo -e "${RED}Could not determine latest version${NC}"
        exit 1
    fi
    
    echo -e "${GREEN}Latest version: $VERSION${NC}"
}

# Download binary
download_binary() {
    DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${BINARY_NAME}-${PLATFORM}"
    
    if [ "$OS" = "windows" ]; then
        DOWNLOAD_URL="${DOWNLOAD_URL}.exe"
    fi
    
    TEMP_DIR=$(mktemp -d)
    TEMP_FILE="${TEMP_DIR}/${BINARY_NAME}"
    
    echo -e "${BLUE}Downloading ${BINARY_NAME} ${VERSION} for ${PLATFORM}...${NC}"
    
    if command -v curl &> /dev/null; then
        curl -sSL "$DOWNLOAD_URL" -o "$TEMP_FILE"
    elif command -v wget &> /dev/null; then
        wget -q "$DOWNLOAD_URL" -O "$TEMP_FILE"
    fi
    
    if [ ! -f "$TEMP_FILE" ]; then
        echo -e "${RED}Download failed${NC}"
        exit 1
    fi
    
    chmod +x "$TEMP_FILE"
}

# Install binary
install_binary() {
    echo -e "${BLUE}Installing to ${INSTALL_DIR}...${NC}"
    
    if [ -w "$INSTALL_DIR" ]; then
        mv "$TEMP_FILE" "${INSTALL_DIR}/${BINARY_NAME}"
    else
        echo -e "${YELLOW}Requesting sudo access to install to ${INSTALL_DIR}${NC}"
        sudo mv "$TEMP_FILE" "${INSTALL_DIR}/${BINARY_NAME}"
    fi
    
    rm -rf "$TEMP_DIR"
    
    echo -e "${GREEN}✓ ${BINARY_NAME} installed successfully!${NC}"
}

# Verify installation
verify_installation() {
    if command -v "$BINARY_NAME" &> /dev/null; then
        echo -e "${GREEN}✓ Installation verified${NC}"
        echo ""
        echo -e "${BLUE}Version:${NC}"
        "$BINARY_NAME" --version
        echo ""
        echo -e "${BLUE}Usage:${NC}"
        echo "  $BINARY_NAME           # Launch TUI"
        echo "  $BINARY_NAME --help    # Show help"
    else
        echo -e "${RED}Installation verification failed${NC}"
        echo -e "${YELLOW}Make sure ${INSTALL_DIR} is in your PATH${NC}"
        exit 1
    fi
}

# Print banner
print_banner() {
    echo -e "${GREEN}"
    echo "╔═══════════════════════════════════════════════════════════════╗"
    echo "║                                                               ║"
    echo "║   🌿 GitFlow TUI Installer                                   ║"
    echo "║   Complete Git Management Terminal UI                        ║"
    echo "║                                                               ║"
    echo "╚═══════════════════════════════════════════════════════════════╝"
    echo -e "${NC}"
    echo ""
}

# Main installation
main() {
    print_banner
    
    detect_platform
    get_latest_version
    download_binary
    install_binary
    verify_installation
    
    echo ""
    echo -e "${GREEN}Installation complete! 🎉${NC}"
    echo ""
    echo -e "${BLUE}Next steps:${NC}"
    echo "  1. Navigate to a git repository"
    echo "  2. Run: gitflow-tui"
    echo "  3. Press '?' for help"
    echo ""
    echo -e "${BLUE}Documentation:${NC} https://github.com/${REPO}"
}

# Run main function
main
