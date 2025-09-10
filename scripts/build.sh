#!/usr/bin/env sh

# get version from latest tag
VERSION="$(git describe --tags)"

# create local build
go build -ldflags "-X main.version=$VERSION" -o depgraph depgraph.go
