#!/bin/bash

# Builds application to ./output/

OUTPUT_DIR=$1
# Build date is UTC as ISO 8601 format
# Expects VERSION and COMMIT environment variables to be set
BUILD_DATE="$(date -Is -u)"
CGO_ENABLED=0
GOOS=linux

rm -rf ${OUTPUT_DIR}
mkdir ${OUTPUT_DIR}

# Copy templates and static content
cp -r ./templates ${OUTPUT_DIR}/templates
cp -r ./static ${OUTPUT_DIR}/static

go build \
    -a -installsuffix nocgo \
    -ldflags "-X 'main.version=${VERSION}' -X 'main.buildDate=${BUILD_DATE}' -X 'main.commitHash=${COMMIT}'" \
    -o ${OUTPUT_DIR}/app .
