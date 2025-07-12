#!/usr/bin/env bash
set -euo pipefail

# 1) Move to the project root (where this script lives)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# 2) Build the binary into bin/
mkdir -p bin
echo "Building cmd/opensesame → bin/opensesame-hub"
go build -o bin/opensesame-hub cmd/opensesame/main.go

# 3) Copy migrations/ so migrate can find them in bin/
echo "Copying migrations/ → bin/migrations/"
rm -rf bin/migrations
cp -r migrations bin/

# 4) cd into bin/ so our working‐dir is bin/
pushd bin > /dev/null

# 5) Run the binary (it will create app.db here)
echo "Starting application from $(pwd)"
./opensesame-hub

# 6) back out
popd > /dev/null