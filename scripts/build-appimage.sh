#!/bin/bash
# Build AppImage for Cloud Manage CLI
# Usage: ./scripts/build-appimage.sh

set -e

APP_NAME="cloud-manage"
APP_VERSION="0.1.0"
ARCH="x86_64"

echo "Building ${APP_NAME} AppImage..."

# Build CLI binary
echo "Building CLI binary..."
CGO_ENABLED=0 go build -o "build/bin/${APP_NAME}" .

# Create AppDir structure
APPDIR="build/AppDir"
rm -rf "${APPDIR}"
mkdir -p "${APPDIR}/usr/bin"
mkdir -p "${APPDIR}/usr/share/applications"
mkdir -p "${APPDIR}/usr/share/icons/hicolor/256x256/apps"

# Copy binary
cp "build/bin/${APP_NAME}" "${APPDIR}/usr/bin/"

# Create desktop file
cat > "${APPDIR}/usr/share/applications/${APP_NAME}.desktop" << EOF
[Desktop Entry]
Name=Cloud Manage
Comment=阿里云资源管理工具
Exec=${APP_NAME}
Icon=${APP_NAME}
Terminal=true
Type=Application
Categories=Utility;
EOF

# Create icon directory and copy icon
mkdir -p "${APPDIR}/usr/share/icons/hicolor/256x256/apps"

# Create a simple PNG icon using python
python3 -c "
import struct, zlib
def create_png():
    width, height = 256, 256
    ihdr_data = struct.pack('>IIBBBBB', width, height, 8, 2, 0, 0, 0)
    ihdr_crc = zlib.crc32(b'IHDR' + ihdr_data) & 0xffffffff
    ihdr = struct.pack('>I', 13) + b'IHDR' + ihdr_data + struct.pack('>I', ihdr_crc)
    raw_data = b''
    for y in range(height):
        raw_data += b'\x00'
        for x in range(width):
            raw_data += b'\x18\x90\xff'
    compressed = zlib.compress(raw_data)
    idat_crc = zlib.crc32(b'IDAT' + compressed) & 0xffffffff
    idat = struct.pack('>I', len(compressed)) + b'IDAT' + compressed + struct.pack('>I', idat_crc)
    iend_crc = zlib.crc32(b'IEND') & 0xffffffff
    iend = struct.pack('>I', 0) + b'IEND' + struct.pack('>I', iend_crc)
    return b'\x89PNG\r\n\x1a\n' + ihdr + idat + iend
with open('${APPDIR}/usr/share/icons/hicolor/256x256/apps/${APP_NAME}.png', 'wb') as f:
    f.write(create_png())
"

# Copy desktop file and icon to AppDir root
cp "${APPDIR}/usr/share/applications/${APP_NAME}.desktop" "${APPDIR}/"
cp "${APPDIR}/usr/share/icons/hicolor/256x256/apps/${APP_NAME}.png" "${APPDIR}/"

# Create AppRun
cat > "${APPDIR}/AppRun" << 'EOF'
#!/bin/bash
SELF=$(readlink -f "$0")
HERE=${SELF%/*}
exec "${HERE}/usr/bin/cloud-manage" "$@"
EOF
chmod +x "${APPDIR}/AppRun"

# Download appimagetool if not exists
APPIMAGETOOL="build/appimagetool"
if [ ! -f "${APPIMAGETOOL}" ]; then
    echo "Downloading appimagetool..."
    curl -fsSL "https://github.com/AppImage/AppImageKit/releases/download/continuous/appimagetool-${ARCH}.AppImage" -o "${APPIMAGETOOL}"
    chmod +x "${APPIMAGETOOL}"
fi

# Build AppImage
echo "Creating AppImage..."
OUTPUT="build/${APP_NAME}-${APP_VERSION}-${ARCH}.AppImage"
ARCH=${ARCH} ./${APPIMAGETOOL} "${APPDIR}" "${OUTPUT}" 2>/dev/null || \
ARCH=${ARCH} ./${APPIMAGETOOL} --no-appstream "${APPDIR}" "${OUTPUT}"

chmod +x "${OUTPUT}"

echo ""
echo "AppImage created: ${OUTPUT}"
echo "Size: $(du -h ${OUTPUT} | cut -f1)"
echo ""
echo "To run: ./${OUTPUT}"
