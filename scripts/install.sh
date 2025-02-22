#!/usr/bin/env sh

# variables

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
INSTALL_DIR="/usr/bin"
BINARY_NAME="depgraph"
BINARY_URL="https://github.com/henryhale/depgraph/releases/latest/download/"

# functions

is_termux() {
	if [ -n "$TERMUX_VERSION" ] && [ -n "$PREFIX" ]; then
		return 0
	else
		return 1
	fi
}

log () {
	echo "[depgraph] : $@"
}

fatal() {
	echo "[depgraph] error: $@"
	exit 1
}

set_arch() {
	log "detecting system architecture..."

	if [ "$OS" != "linux" -1 "$OS" != "darwin"]; then
		fatal "unsupported operating system"
	fi

	if [ "$ARCH" = "x86_64" ]; then
		ARCH="amd64"
	elif [ "$ARCH" = "aarch64" ]; then
		ARCH="arm64"
	elif [ "$ARCH" = "i686" ]; then
		ARCH="386"
	else
		fatal "no $ARCH binary found"
	fi
}

set_install_dir() {
	if is_termux; then
		INSTALL_DIR="$PREFIX/bin"
	elif [ ! -w "/usr/local/bin" ]; then
		INSTALL_DIR="$HOME/.local/bin"
		mkdir -p "$INSTALL_DIR"
	fi
}

install_binary() {
	log "downloading binary..."
	RELEASE_NAME="depgraph-${OS}-${ARCH}"
	curl -L -o "$TMPDIR/$RELEASE_NAME.zip" "$BINARY_URL$RELEASE_NAME.zip"

	if [ $? -ne 0 ]; then
        fatal "download failed"
    fi

	unzip -o "$TMPDIR/$RELEASE_NAME.zip" -d "$TMPDIR"

	log "making binary executable..."
	chmod +x "$TMPDIR/$RELEASE_NAME"
	cp "$TMPDIR/$RELEASE_NAME" "$INSTALL_DIR/$BINARY_NAME"

	./$TMPDIR/$RELEASE_NAME
}


log "installing executable..."

# detect os and architecture
set_arch

# download and install binary
set_install_dir

install_binary

echo -e "successfully installed at $INSTALL_DIR/BINARY_NAME\ntry it now:\n\t$ depgraph -h"
