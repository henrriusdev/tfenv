#!/bin/bash

# Detect OS and Architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

if [ "$ARCH" = "x86_64" ]; then
    ARCH="amd64"
elif [ "$ARCH" = "aarch64" ] || [ "$ARCH" = "arm64" ]; then
    ARCH="arm64"
else
    echo "❌ Unsupported architecture: $ARCH"
    exit 1
fi

# Set binary name based on OS
if [ "$OS" = "darwin" ]; then
    BINARY_NAME="tfenv-darwin-$ARCH"
elif [ "$OS" = "linux" ]; then
    BINARY_NAME="tfenv-linux-$ARCH"
else
    echo "❌ Unsupported operating system: $OS"
    exit 1
fi

# Download the correct binary
DOWNLOAD_URL="https://github.com/henrriusdev/tfenv/releases/latest/download/$BINARY_NAME"
DEST_PATH="/usr/local/bin/tfenv"

echo "🔽 Downloading $BINARY_NAME from $DOWNLOAD_URL..."
curl -L -o "$DEST_PATH" "$DOWNLOAD_URL"

# Make it executable
chmod +x "$DEST_PATH"

echo "✅ Installation complete! Run 'tfenv' to use the CLI."
