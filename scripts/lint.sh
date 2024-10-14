#!/usr/bin/env sh

gofmt -w .

golangci-lint run