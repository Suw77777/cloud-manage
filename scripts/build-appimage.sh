#!/bin/bash
# Build AppImage for Cloud Manage
# Usage: ./scripts/build-appimage.sh

set -e

APP_NAME="cloud-manage"
APP_VERSION="0.0.12"
ARCH="x86_64"

echo "Building ${APP_NAME} AppImage..."

# Build the application first
echo "Building application..."
wails build

# Create AppDir structure
APPDIR="build/AppDir"
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
    "libatk-1.0.so.0"
    "libatk-bridge-2.0.so.0"
    "libglib-2.0.so.0"
    "libgio-2.0.so.0"
    "libgobject-2.0.so.0"
)

for lib in "${LIBS[@]}"; do
    lib_path=$(ldconfig -p 2>/dev/null | grep "${lib}" | head -1 | awk '{print $NF}')
    if [ -n "${lib_path}" ] && [ -f "${lib_path}" ]; then
        cp -L "${lib_path}" "${APPDIR}/usr/lib/"
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

# Copy icon if exists
if [ -f "frontend/src/assets/images/logo.png" ]; then
    cp "frontend/src/assets/images/logo.png" "${APPDIR}/usr/share/icons/hicolor/256x256/apps/${APP_NAME}.png"
else
    # Create a simple placeholder icon
    convert -size 256x256 xc:blue -fill white -gravity center -pointsize 40 -annotate 0 "CM" \
        "${APPDIR}/usr/share/icons/hicolor/256x256/apps/${APP_NAME}.png" 2>/dev/null || true
fi

# Copy desktop file to AppDir root
cp "${APPDIR}/usr/share/applications/${APP_NAME}.desktop" "${APPDIR}/"
# Copy icon to AppDir root
cp "${APPDIR}/usr/share/icons/hicolor/256x256/apps/${APP_NAME}.png" "${APPDIR}/" 2>/dev/null || true

# Create AppRun
cat > "${APPDIR}/AppRun" << 'EOF'
#!/bin/bash
SELF=$(readlink -f "$0")
HERE=${SELF%/*}
export PATH="${HERE}/usr/bin/:${PATH}"
export LD_LIBRARY_PATH="${HERE}/usr/lib/:${LD_LIBRARY_PATH}"
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
