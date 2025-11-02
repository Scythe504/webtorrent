#!/usr/bin/env bash
set -e

REPO="scythe504/fluxstream"
INSTALL_DIR="$HOME/.local/bin"
mkdir -p "$INSTALL_DIR"

echo "Detecting latest version..."
LATEST_TAG=$(curl -s https://api.github.com/repos/$REPO/releases/latest | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/')
VERSION="${LATEST_TAG#v}"  # remove leading 'v' if present
echo "Latest version: $LATEST_TAG"

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
  x86_64) ARCH="amd64" ;;
  aarch64) ARCH="arm64" ;;
esac

FILE="fluxstream_${VERSION}_${OS}_${ARCH}.tar.gz"
URL="https://github.com/${REPO}/releases/download/${LATEST_TAG}/${FILE}"

echo "⬇️ Downloading $FILE ..."
curl -L -o /tmp/"$FILE" "$URL"

echo "📦 Extracting to $INSTALL_DIR ..."
tar -xzf /tmp/"$FILE" -C "$INSTALL_DIR"

chmod +x "$INSTALL_DIR/fluxstream"

if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
  echo "Adding $INSTALL_DIR to PATH (append to ~/.bashrc or ~/.zshrc)"
  echo "export PATH=\"$INSTALL_DIR:\$PATH\"" >> ~/.bashrc
fi

echo "✅ FluxStream installed successfully!"
"$INSTALL_DIR/fluxstream" --version || true
