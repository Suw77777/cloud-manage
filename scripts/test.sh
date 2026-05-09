#!/bin/bash
# Run all Go tests for the Cloud 管理小助手 project.
# Usage: ./scripts/test.sh

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

cd "$PROJECT_DIR"

echo "==> Running Go tests..."
echo ""

GO111MODULE=on go test ./... -v

echo ""
echo "==> All tests passed."
