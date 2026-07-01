#!/usr/bin/env sh

go vet ./

gofmt -w ./

golangci-lint run
