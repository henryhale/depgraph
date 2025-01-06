#!/usr/bin/env sh

CI="$1"

# get version from latest tag
VERSION="$(git describe --tags)"

# create local build
localbuild () {
    go build -ldflags "-X main.version=$VERSION" -o depgraph depgraph.go
}

# create builds in ci - github actions
cibuild () {
    GOOS="$1"
    GOARCH="$2"
    EXT=""

    BINARY_NAME="depgraph-$GOOS-$GOARCH$EXT"

    if [ "$GOOS" = "windows" ]; then EXT=".exe"; fi

    go build -ldflags "-X main.version=$VERSION" -o "$BINARY_NAME" depgraph.go

    zip -m "$BINARY_NAME.zip" "$BINARY_NAME"
}

if [ "$CI" = "ci" ]; then
    echo "building for linux..."
    cibuild linux amd64
    cibuild linux arm64
    cibuild linux 386

    echo "building for darwin..."
    cibuild darwin amd64
    cibuild darwin arm64

    echo "building for windows..."
    cibuild windows amd64
    cibuild windows arm64
else
    localbuild
fi
