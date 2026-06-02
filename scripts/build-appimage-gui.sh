#!/bin/bash
# Build GUI AppImage for Cloud Manage
# Usage: ./scripts/build-appimage-gui.sh

set -e

APP_NAME="cloud-manage"
APP_VERSION="0.1.0"
ARCH="x86_64"

echo "Building ${APP_NAME} GUI AppImage..."

# Build GUI binary with wails
echo "Building GUI binary..."
export PATH=/usr/local/go/bin:~/go/bin:$PATH
wails build --skipbindings

# Create AppDir structure
APPDIR="build/AppDir-gui"
rm -rf "${APPDIR}"
mkdir -p "${APPDIR}/usr/bin"
mkdir -p "${APPDIR}/usr/lib"
mkdir -p "${APPDIR}/usr/share/applications"
mkdir -p "${APPDIR}/usr/share/icons/hicolor/256x256/apps"

# Copy binary
cp "build/bin/${APP_NAME}" "${APPDIR}/usr/bin/"

# Copy required shared libraries
echo "Copying shared libraries..."
LIBS=(
    "libwebkit2gtk-4.0.so.37"
    "libjavascriptcoregtk-4.0.so.18"
    "libgtk-3.so.0"
    "libgdk-3.so.0"
    "libgdk_pixbuf-2.0.so.0"
    "libpango-1.0.so.0"
    "libpangocairo-1.0.so.0"
    "libcairo.so.2"
    "libcairo-gobject.so.2"
    "libatk-1.0.so.0"
    "libatk-bridge-2.0.so.0"
    "libglib-2.0.so.0"
    "libgio-2.0.so.0"
    "libgobject-2.0.so.0"
    "libgmodule-2.0.so.0"
    "libfontconfig.so.1"
    "libfreetype.so.6"
    "libX11.so.6"
    "libXcomposite.so.1"
    "libXdamage.so.1"
    "libXext.so.6"
    "libXfixes.so.3"
    "libXi.so.6"
    "libXrandr.so.2"
    "libXrender.so.1"
    "libwayland-client.so.0"
    "libwayland-server.so.0"
    "libwayland-egl.so.1"
    "libwayland-cursor.so.0"
    "libEGL.so.1"
    "libGL.so.1"
    "libGLESv2.so.2"
    "libxkbcommon.so.0"
    "libepoxy.so.0"
    "libpcre2-8.so.0"
    "libffi.so.8"
    "libresolv.so.2"
    "libselinux.so.1"
    "libthai.so.0"
    "libfribidi.so.0"
    "libharfbuzz.so.0"
    "libpng16.so.16"
    "libjpeg.so.8"
    "libtiff.so.5"
    "libxml2.so.2"
    "libxslt.so.1"
    "libicui18n.so.70"
    "libicuuc.so.70"
    "libicudata.so.70"
    "libz.so.1"
    "liblzma.so.5"
    "libsqlite3.so.0"
    "libsecret-1.so.0"
    "libtasn1.so.6"
    "libgnutls.so.30"
    "libhogweed.so.6"
    "libnettle.so.8"
    "libgmp.so.10"
    "libidn2.so.0"
    "libunistring.so.2"
    "libp11-kit.so.0"
)

for lib in "${LIBS[@]}"; do
    lib_path=$(find /usr/lib/x86_64-linux-gnu /lib/x86_64-linux-gnu -name "${lib}*" 2>/dev/null | head -1)
    if [ -n "${lib_path}" ] && [ -f "${lib_path}" ]; then
        cp -L "${lib_path}" "${APPDIR}/usr/lib/" 2>/dev/null || true
    fi
done

# Copy additional libs from ldd output
echo "Copying additional dependencies..."
ldd "build/bin/${APP_NAME}" 2>/dev/null | grep "=>" | awk '{print $3}' | while read lib; do
    if [ -f "${lib}" ] && [ ! -f "${APPDIR}/usr/lib/$(basename ${lib})" ]; then
        cp -L "${lib}" "${APPDIR}/usr/lib/" 2>/dev/null || true
    fi
done

# Create desktop file
cat > "${APPDIR}/usr/share/applications/${APP_NAME}.desktop" << EOF
[Desktop Entry]
Name=Cloud Manage
Comment=阿里云资源管理工具
Exec=${APP_NAME}
Icon=${APP_NAME}
Terminal=false
Type=Application
Categories=Utility;
EOF

# Create icon
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

# Create AppRun with library path
cat > "${APPDIR}/AppRun" << 'EOF'
#!/bin/bash
SELF=$(readlink -f "$0")
HERE=${SELF%/*}
export LD_LIBRARY_PATH="${HERE}/usr/lib/:${LD_LIBRARY_PATH}"
export GTK_PATH="${HERE}/usr/lib/gtk-3.0"
export GDK_PIXBUF_MODULE_FILE="${HERE}/usr/lib/gdk-pixbuf-2.0/2.10.0/loaders.cache"
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
OUTPUT="build/${APP_NAME}-gui-${APP_VERSION}-${ARCH}.AppImage"
ARCH=${ARCH} ./${APPIMAGETOOL} "${APPDIR}" "${OUTPUT}" 2>/dev/null || \
ARCH=${ARCH} ./${APPIMAGETOOL} --no-appstream "${APPDIR}" "${OUTPUT}"

chmod +x "${OUTPUT}"

echo ""
echo "GUI AppImage created: ${OUTPUT}"
echo "Size: $(du -h ${OUTPUT} | cut -f1)"
echo ""
echo "To run: ./${OUTPUT}"
