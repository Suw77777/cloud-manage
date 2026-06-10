#!/bin/bash
# Build the Cloud 管理小助手 application.
# Usage: ./scripts/build.sh [--gui]

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

cd "$PROJECT_DIR"

export PATH=/usr/local/go/bin:~/go/bin:$PATH

echo "==> Building Cloud 管理小助手..."
echo ""

if [ "$1" = "--gui" ]; then
    echo "Building with GUI support (wails)..."
    wails build --skipbindings
    echo ""
    echo "==> Build complete. Output: build/bin/cloud-manage"
else
    echo "Building CLI/TUI version..."
    go build -trimpath -ldflags "-s -w" -o cloud-manage .
    echo ""
    echo "==> Build complete. Output: cloud-manage"
fi

echo ""
echo "Usage:"
echo "  ./cloud-manage          # Auto-detect mode (TUI if no display)"
echo "  ./cloud-manage --tui    # Force TUI mode"
echo "  ./cloud-manage --gui    # Force GUI mode (requires display)"
echo "  ./cloud-manage --cli    # Force CLI mode"
echo "  ./cloud-manage ecs list # CLI command"
