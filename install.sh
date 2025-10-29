#!/bin/bash
set -e

# codezure installer script
# Usage: curl -fsSL https://raw.githubusercontent.com/OlaHulleberg/codezure/main/install.sh | bash

REPO="OlaHulleberg/codezure"
INSTALL_DIR="/usr/local/bin"
FALLBACK_DIR="$HOME/.local/bin"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo "Installing codezure..."

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
    EXT="zip"; BINARY="codezure.exe"
else
    EXT="tar.gz"; BINARY="codezure"
fi

ARCHIVE="codezure_${OS}_${ARCH}.${EXT}"
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
if command -v codezure >/dev/null 2>&1; then
    VERSION=$(codezure manage version 2>&1 || echo "unknown")
    echo ""; echo -e "${GREEN}✓ codezure installed successfully!${NC}"; echo "  Version: $VERSION"; echo ""
    echo "Get started:"
    echo "  codezure manage config"
    echo "  codezure"
else
    echo ""; echo -e "${YELLOW}Note: You may need to restart your shell or update PATH${NC}"; echo "Then verify with: codezure manage version"
fi
# Detect legacy codzure and offer uninstall + config migration
if command -v codzure >/dev/null 2>&1; then
    LEGACY_PATH=$(command -v codzure)
    echo "Found legacy installation: $LEGACY_PATH"
    # Prompt (stdout) and read from tty if available; default to Y otherwise
    RESP="Y"
    if [ -t 0 ] || [ -t 1 ] || [ -t 2 ]; then
        printf "Uninstall legacy codzure and migrate configs to codezure? [Y/n] "
        set +e
        if [ -r /dev/tty ]; then
            read -r RESP < /dev/tty || RESP="Y"
        else
            read -r RESP || RESP="Y"
        fi
        set -e
        [ -z "$RESP" ] && RESP="Y"
    else
        echo "Non-interactive shell detected; defaulting to 'Y' for uninstall + migration."
        RESP="Y"
    fi
    if [ "$RESP" = "Y" ] || [ "$RESP" = "y" ]; then
        echo "Uninstalling codzure..."
        if [ -w "$LEGACY_PATH" ]; then
            rm -f "$LEGACY_PATH" || echo -e "${YELLOW}Warning: failed to remove $LEGACY_PATH${NC}"
        elif command -v sudo >/dev/null 2>&1; then
            sudo rm -f "$LEGACY_PATH" || echo -e "${YELLOW}Warning: failed to remove $LEGACY_PATH${NC}"
        else
            echo -e "${YELLOW}Warning: cannot remove $LEGACY_PATH without sudo; please remove manually${NC}"
        fi
        # Migrate ~/.codzure -> ~/.codezure (non-destructive)
        OLD_DIR="$HOME/.codzure"
        NEW_DIR="$HOME/.codezure"
        if [ -d "$OLD_DIR" ]; then
            echo "Migrating configuration from ~/.codzure to ~/.codezure..."
            mkdir -p "$NEW_DIR/profiles"
            # Copy profiles if present (skip existing)
            if [ -d "$OLD_DIR/profiles" ]; then
                for f in "$OLD_DIR"/profiles/*.json; do
                    [ -e "$f" ] || continue
                    base=$(basename "$f")
                    dst="$NEW_DIR/profiles/$base"
                    if [ ! -e "$dst" ]; then
                        cp "$f" "$dst"
                    fi
                done
            fi
            # Copy current-profile.txt if not present
            if [ -f "$OLD_DIR/current-profile.txt" ] && [ ! -f "$NEW_DIR/current-profile.txt" ]; then
                cp "$OLD_DIR/current-profile.txt" "$NEW_DIR/current-profile.txt"
            fi
            # Copy legacy current.env so app can migrate it on first run
            if [ -f "$OLD_DIR/current.env" ] && [ ! -f "$NEW_DIR/current.env" ]; then
                cp "$OLD_DIR/current.env" "$NEW_DIR/current.env"
            fi
            echo -e "${GREEN}✓ Migration complete${NC}"
        else
            echo "No ~/.codzure directory found; skipping config migration."
        fi
    else
        echo "Skipping uninstall/migration of legacy codzure."
    fi
fi
