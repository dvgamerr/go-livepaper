#!/bin/bash

# Get version from VERSION file
VERSION=$(cat VERSION)
echo "Building livepaper version $VERSION"

# Build with version embedded
go build -ldflags "-X main.VERSION=$VERSION" -o livepaper