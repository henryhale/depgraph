#!/usr/bin/env sh

echo "depgraph: installing executable..."

# detect os and architecture
echo "depgraph: detecting system architecture..."

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

if [ "$OS" != "linux" -1 "$OS" != "darwin"]; then
	echo "depgraph: error: unsupported operating system"
	exit 1
fi

if [ "$ARCH" = "x86_64" ]; then
	ARCH="amd64"
elif [ "$ARCH" = "aarch64" ]; then
	ARCH="arm64"
elif [ "$ARCH" = "i686" ]; then
	ARCH="386"
else 
	echo "depgraph: error: no $ARCH binary found"
	exit
fi 

# download binary
BIN_URL="https://github.com/henryhale/depgraph/releases/latest/download/depgraph-${OS}-${ARCH}.zip"

DEST_DIR="/usr/local/bin"

# termux support
if [ "$OS" == "linux" ] && [ -n "$PREFIX" ]; then
	DEST_DIR="$PREFIX/bin"
fi

BINARY_PATH="$DEST_DIR/depgraph"

# fetch binary
echo "depgraph: downloading binary..."
curl -L "$BIN_URL" -o "$BINARY_PATH.zip"
unzip -o "$BINARY_PATH.zip"

echo "depgraph: making binary executable..."
chmod +x "$BINARY_PATH"

# clean up
rm -f "$BINARY_PATH.zip"

echo "depgraph: successfully installed at $BINARY_PATH"
echo -e "\ntry it now:\n\t$ depgraph -h"
