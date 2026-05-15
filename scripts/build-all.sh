#!/bin/bash
# Build all releases: CLI (all platforms) + GUI (Docker).
# Usage: ./scripts/build-all.sh [--cli-only] [--gui-only]

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

cd "$PROJECT_DIR"

BUILD_CLI=true
BUILD_GUI=true

# Parse args
for arg in "$@"; do
    case $arg in
        --cli-only) BUILD_GUI=false ;;
        --gui-only) BUILD_CLI=false ;;
    esac
done

echo "=========================================="
echo "  Cloud 管理小助手 - Multi-Platform Build"
echo "=========================================="
echo ""

# Build CLI
if [ "$BUILD_CLI" = true ]; then
    echo "[1/2] Building CLI for all platforms..."
    ./scripts/build-cli.sh
    echo ""
fi

# Build GUI
if [ "$BUILD_GUI" = true ]; then
    echo "[2/2] Building GUI using Docker..."
    if command -v docker &> /dev/null; then
        ./scripts/build-gui-docker.sh
    else
        echo "  WARNING: Docker not found. Skipping GUI build."
        echo "  Install Docker or run: ./scripts/build-gui-docker.sh"
    fi
    echo ""
fi

echo "=========================================="
echo "  Build Summary"
echo "=========================================="
echo ""
ls -lh release/
echo ""
echo "Done!"
