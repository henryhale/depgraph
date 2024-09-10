#!/usr/bin/env sh

# detect os and architecture
echo "detecting system architecture..."

OS=$(uname | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

if [ "$ARCH" == "x86_64" ]; then
	ARCH="amd64"
elif [ "$ARCH" == "aarch64" ]; then
	ARCH="arm64"
elif [ "$ARCH" == "i686" ]; then
	ARCH="386"
fi

# download binary
BIN_URL="https://github.com/henryhale/depgraph/releases/latest/download/depgraph-${OS}-${ARCH}"

DEST_DIR="/usr/local/bin"

# termux support
if [ "$OS" == "linux" ] && [ -n "$PREFIX" ]; then
	DEST_DIR="$PREFIX/bin"
fi

# fetch
echo "downloading binary..."

curl -L "$BIN_URL" -o "$DEST_DIR/depgraph"
chmod +x "$DEST_DIR/depgraph"

echo "successfully installed depgraph at $DEST_DIR/depgraph"
echo -e "\ntry it now:\n\t$ depgraph -h"
