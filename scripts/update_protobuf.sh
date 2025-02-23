#!/bin/bash

set -e  # Exit on any error

echo "Updating common protobuf..."
GOPROXY=direct go get -u github.com/ngdangkietswe/swe-protobuf-shared/generated/common
go mod tidy
go mod vendor
echo "Successfully updated common protobuf!"