#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# build
mkdir -p bin
echo "Building cmd/opensesame → bin/opensesame-hub"
go build -o bin/opensesame-hub cmd/opensesame/main.go

# cd into bin/ to run binaries then back to the original directory
pushd bin > /dev/null
echo "Starting application from $(pwd)"
./opensesame-hub
popd > /dev/null