#!/bin/bash
# Start the Wails development server with hot reload.
# Usage: ./scripts/dev.sh

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

cd "$PROJECT_DIR"

echo "==> Starting Cloud 管理小助手 in development mode..."
echo ""

export PATH=$PATH:/root/go/bin
wails dev
