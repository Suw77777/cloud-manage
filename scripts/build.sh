#!/bin/bash
# Build the Cloud 管理小助手 application.
# Usage: ./scripts/build.sh

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

cd "$PROJECT_DIR"

echo "==> Building Cloud 管理小助手..."
echo ""

export PATH=$PATH:/root/go/bin
wails build

echo ""
echo "==> Build complete. Output: build/bin/"
