#!/usr/bin/env sh

# variables

OS="unknown"
ARCH="unknown"
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

detect_os() {
	log "detecting system os..."
if is_termux; then
		OS="linux"
		return
	fi
	OS=$(uname -s | tr '[:upper:]' '[:lower:]')
	if [ "$OS" != "linux" -1 "$OS" != "darwin"]; then
		fatal "unsupported operating system"
	fi
}

detect_arch() {
	log "detecting system architecture..."

	if [ "$OS" != "linux" -1 "$OS" != "darwin"]; then
		fatal "unsupported operating system"
	fi

	ARCH=$(uname -m)

	case "$ARCH" in
        x86_64|amd64)
            ARCH="amd64"
            ;;
        i686|i386)
            ARCH="386"
            ;;
        aarch64|arm64)
            ARCH="arm64"
            ;;
        armv7*|armv8*|arm)
            ARCH="arm"
            ;;
        *)
			fatal "no $ARCH binary found"
            ;;
    esac
}

log() {
	echo "[depgraph] : $@"
}

fatal() {
	echo "[depgraph] error: $@"
	exit 1
}

check_exec_error() {
	if [ "$?" -ne 0 ]; then
		exit 1
	fi
}

set_install_dir() {
	if is_termux; then
		INSTALL_DIR="$PREFIX/bin"
	elif [ ! -w "/usr/local/bin" ]; then
		INSTALL_DIR="$HOME/.local/bin"
		mkdir -p "$INSTALL_DIR"
		check_exec_error
	fi
}

install_binary() {
	log "downloading binary..."
	RELEASE_NAME="depgraph-${OS}-${ARCH}"
	curl -L -o "$TMPDIR/$RELEASE_NAME.zip" "$BINARY_URL$RELEASE_NAME.zip"

	check_exec_error

	unzip -o "$TMPDIR/$RELEASE_NAME.zip" -d "$TMPDIR"

	check_exec_error

	log "making binary executable..."
	chmod +x "$TMPDIR/$RELEASE_NAME"

	check_exec_error

	cp "$TMPDIR/$RELEASE_NAME" "$INSTALL_DIR/$BINARY_NAME"

	check_exec_error
}


log "installing executable..."

# detect os and architecture
detect_os
detect_arch

# download and install binary
set_install_dir

install_binary

echo -e "successfully installed at $INSTALL_DIR/$BINARY_NAME\ntry it now:\n\t$ depgraph -h"
