#!/bin/bash
# Build GUI using Docker (for Ubuntu 22.04 compatibility).
# Usage: ./scripts/build-gui-docker.sh

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

cd "$PROJECT_DIR"

VERSION=$(grep 'const version' main.go | cut -d'"' -f2)
if [ -z "$VERSION" ]; then
    VERSION="v0.0.9"
fi

OUTPUT_DIR="release"
mkdir -p "$OUTPUT_DIR"

echo "==> Building Cloud GUI $VERSION using Docker..."
echo ""

# Build Docker image
echo "  Building Docker image..."
docker build -t cloud-manage-builder .

# Run build inside container
echo "  Running wails build inside container..."
docker run --rm \
    -v "$(pwd):/app" \
    -v "$HOME/go/pkg/mod:/root/go/pkg/mod" \
    cloud-manage-builder \
    bash -c "wails build && cp build/bin/cloud-manage /app/release/cloud-manage-linux-amd64"

echo ""
echo "==> Build complete!"
echo ""
ls -lh "${OUTPUT_DIR}/cloud-manage-linux-amd64"
