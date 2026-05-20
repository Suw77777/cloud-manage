#!/bin/bash
# scripts/test.sh - Run tests with coverage

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

cd "$PROJECT_DIR"

echo "Running unit tests..."
go test ./... -v -coverprofile=coverage.out -covermode=atomic

echo ""
echo "Coverage report:"
go tool cover -func=coverage.out

echo ""
echo "Generating HTML coverage report..."
go tool cover -html=coverage.out -o coverage.html

echo ""
echo "Total coverage:"
go tool cover -func=coverage.out | tail -1
