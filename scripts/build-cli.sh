#!/bin/bash
# Build CLI for multiple platforms.
# Usage: ./scripts/build-cli.sh

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

cd "$PROJECT_DIR"

VERSION=$(grep 'const Version' internal/consts/version.go | cut -d'"' -f2)
if [ -z "$VERSION" ]; then
    VERSION="v0.2.0"
fi

OUTPUT_DIR="release"
mkdir -p "$OUTPUT_DIR"

echo "==> Building Cloud CLI $VERSION for multiple platforms..."
echo ""

# Build matrix
PLATFORMS=(
    "linux/amd64"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
)

for PLATFORM in "${PLATFORMS[@]}"; do
    GOOS="${PLATFORM%/*}"
    GOARCH="${PLATFORM#*/}"

    OUTPUT_NAME="cloud-cli-${GOOS}-${GOARCH}"
    if [ "$GOOS" = "windows" ]; then
        OUTPUT_NAME="${OUTPUT_NAME}.exe"
    fi

    echo "  Building ${OUTPUT_NAME}..."
    env CGO_ENABLED=0 GOOS="$GOOS" GOARCH="$GOARCH" go build \
        -trimpath \
        -ldflags "-s -w" \
        -o "${OUTPUT_DIR}/${OUTPUT_NAME}" \
        .
done

# Generate checksums
echo ""
echo "==> Generating checksums..."
cd "$OUTPUT_DIR"
sha256sum cloud-cli-* > checksums.txt 2>/dev/null || shasum -a 256 cloud-cli-* > checksums.txt
cd "$PROJECT_DIR"

echo ""
echo "==> Build complete! Output in ${OUTPUT_DIR}/"
echo ""
ls -lh "${OUTPUT_DIR}"/cloud-cli-*
echo ""
echo "Checksums:"
cat "${OUTPUT_DIR}/checksums.txt"
