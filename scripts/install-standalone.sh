#!/bin/sh
set -e

REPO=danlafeir/devctl-timecard
BINARY=timecard
INSTALL_DIR=~/.local/bin

# Detect OS
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
case "$OS" in
  linux) OS=linux ;;
  darwin) OS=darwin ;;
  *) echo "Unsupported OS: $OS"; exit 1 ;;
esac

# Detect ARCH
ARCH=$(uname -m)
case "$ARCH" in
  x86_64|amd64) ARCH=amd64 ;;
  arm64|aarch64) ARCH=arm64 ;;
  *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

# Find the standalone binary for this OS/ARCH
# Standalone builds use the pattern: timecard-<os>-<arch>
FILENAME="timecard-$OS-$ARCH"
URL="https://raw.githubusercontent.com/$REPO/main/bin/$FILENAME"

TMP=$(mktemp)
echo "Downloading $URL ..."
curl -sSLfL "$URL" -o "$TMP"
chmod +x "$TMP"

# Move to install dir (may require sudo)
echo "Installing to $INSTALL_DIR/$BINARY ..."
mv "$TMP" "$INSTALL_DIR/$BINARY"

if command -v $BINARY >/dev/null 2>&1; then
  echo "$BINARY installed successfully!"
  echo "Ensure $INSTALL_DIR is in your PATH."
else
  echo "Install failed: $BINARY not found in PATH." >&2
  exit 1
fi
