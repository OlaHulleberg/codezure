#!/bin/bash
set -e

# codzure installer script
# Usage: curl -fsSL https://raw.githubusercontent.com/OlaHulleberg/codzure/main/install.sh | bash

REPO="OlaHulleberg/codzure"
INSTALL_DIR="/usr/local/bin"
FALLBACK_DIR="$HOME/.local/bin"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo "Installing codzure..."

OS="$(uname -s)"
ARCH="$(uname -m)"

case "$OS" in
    Linux*)     OS="linux";;
    Darwin*)    OS="darwin";;
    MINGW*|MSYS*|CYGWIN*) OS="windows";;
    *) echo -e "${RED}Error: Unsupported operating system: $OS${NC}"; exit 1;;
esac

case "$ARCH" in
    x86_64|amd64)   ARCH="amd64";;
    aarch64|arm64)  ARCH="arm64";;
    *) echo -e "${RED}Error: Unsupported architecture: $ARCH${NC}"; exit 1;;
esac

echo "Detected: $OS/$ARCH"

if [ "$OS" = "windows" ]; then
    EXT="zip"; BINARY="codzure.exe"
else
    EXT="tar.gz"; BINARY="codzure"
fi

ARCHIVE="codzure_${OS}_${ARCH}.${EXT}"
echo "Fetching latest release..."
REL_URL="https://api.github.com/repos/$REPO/releases/latest"
DL_URL=$(curl -s "$REL_URL" | grep "browser_download_url.*$ARCHIVE" | cut -d '"' -f 4)
if [ -z "$DL_URL" ]; then
    echo -e "${RED}Error: Could not find release for $OS/$ARCH${NC}";
    echo "Please check https://github.com/$REPO/releases"; exit 1
fi

TMP_DIR=$(mktemp -d); trap "rm -rf $TMP_DIR" EXIT; cd "$TMP_DIR"
echo "Downloading $ARCHIVE..."
curl -fsSL "$DL_URL" -o "$ARCHIVE"
echo "Extracting..."
if [ "$EXT" = "zip" ]; then unzip -q "$ARCHIVE"; else tar -xzf "$ARCHIVE"; fi

if [ ! -f "$BINARY" ]; then echo -e "${RED}Error: Binary not found in archive${NC}"; exit 1; fi
chmod +x "$BINARY"

INSTALLED=false; FINAL_DIR=""
if [ -w "$INSTALL_DIR" ]; then
    mv "$BINARY" "$INSTALL_DIR/"; INSTALLED=true; FINAL_DIR="$INSTALL_DIR"; echo -e "${GREEN}Installed to $INSTALL_DIR/$BINARY${NC}"
elif command -v sudo >/dev/null 2>&1; then
    echo "Installing to $INSTALL_DIR requires sudo..."
    if sudo mv "$BINARY" "$INSTALL_DIR/"; then INSTALLED=true; FINAL_DIR="$INSTALL_DIR"; echo -e "${GREEN}Installed to $INSTALL_DIR/$BINARY${NC}"; fi
fi

if [ "$INSTALLED" = false ]; then
    echo -e "${YELLOW}Installing to $FALLBACK_DIR (user directory)${NC}";
    mkdir -p "$FALLBACK_DIR"; mv "$BINARY" "$FALLBACK_DIR/"; FINAL_DIR="$FALLBACK_DIR";
    echo -e "${GREEN}Installed to $FALLBACK_DIR/$BINARY${NC}";
    if [[ ":$PATH:" != *":$FALLBACK_DIR:"* ]]; then
        echo -e "${YELLOW}Warning: $FALLBACK_DIR is not in your PATH${NC}";
        echo "Add this line to your shell rc:"
        echo "  export PATH=\"\$PATH:$FALLBACK_DIR\""
    fi
fi

cd - > /dev/null
if command -v codzure >/dev/null 2>&1; then
    VERSION=$(codzure manage version 2>&1 || echo "unknown")
    echo ""; echo -e "${GREEN}✓ codzure installed successfully!${NC}"; echo "  Version: $VERSION"; echo ""
    echo "Get started:"
    echo "  codzure manage config"
    echo "  codzure"
else
    echo ""; echo -e "${YELLOW}Note: You may need to restart your shell or update PATH${NC}"; echo "Then verify with: codzure manage version"
fi
