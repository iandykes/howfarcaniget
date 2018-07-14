#!/bin/bash

# Builds application to ./output/

rm -rf output
mkdir output

# Copy templates and static content
cp -r ./templates ./output/templates
cp -r ./static ./output/static

# Build date is UTC as ISO 8601 format
# Expects VERSION and COMMIT environment variables to be set
BUILD_DATE="$(date -Is -u)"
CGO_ENABLED=0
GOOS=linux

go build \
    -a -installsuffix nocgo \
    -ldflags "-X 'main.version=${VERSION}' -X 'main.buildDate=${BUILD_DATE}' -X 'main.commitHash=${COMMIT}'" \
    -o output/app .
