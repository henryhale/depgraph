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
	if [ "$OS" != "linux" ] && [ "$OS" != "darwin"]; then
		fatal "unsupported operating system"
	fi
}

detect_arch() {
	log "detecting system architecture..."

	if [ "$OS" != "linux" ] && [ "$OS" != "darwin"]; then
		fatal "unsupported operating system"
	fi

	ARCH=$(uname -m)

	case "$ARCH" in
        x86_64|amd64)
            ARCH="x86_64"
            ;;
        i686|i386)
            ARCH="i386"
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
	RELEASE_NAME="depgraph_${OS}_${ARCH}"
	RELEASE_TAR="$RELEASE_NAME.tar"
	RELEASE_GZ="$RELEASE_NAME.tar.gz"

	TMP_DIR="/tmp"
	if is_termux; then
		TMP_DIR="$TMPDIR"
	fi
	
	curl -L -o "$TMP_DIR/$RELEASE_GZ" "$BINARY_URL$RELEASE_NAME.tar.gz"

	check_exec_error

	gunzip -f "$TMP_DIR/$RELEASE_GZ"

	check_exec_error

	cd "$TMP_DIR" && tar -xvf "$RELEASE_TAR"

	check_exec_error

	log "making binary executable..."
	chmod +x "$TMP_DIR/$RELEASE_NAME"

	check_exec_error

	cp "$TMP_DIR/$RELEASE_NAME" "$INSTALL_DIR/$BINARY_NAME"

	check_exec_error

	log "cleaning up..."
	rm -f "$TMP_DIR/$BINARY_NAME*"
}


log "installing executable..."

# detect os and architecture
detect_os
detect_arch

# download and install binary
set_install_dir

install_binary

log "successfully installed at $INSTALL_DIR/$BINARY_NAME"
log "try it now: $ depgraph -h"
