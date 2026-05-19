#!/bin/bash
# Install dependencies for Cloud Manage GUI

set -e

echo "Installing Cloud Manage dependencies..."

# Detect package manager
if command -v apt-get &> /dev/null; then
    # Debian/Ubuntu
    sudo apt-get update
    sudo apt-get install -y \
        libgtk-3-0 \
        libwebkit2gtk-4.0-37 \
        libglib2.0-0 \
        libgdk-pixbuf2.0-0 \
        libpango-1.0-0 \
        libcairo2 \
        libatk1.0-0 \
        libatk-bridge2.0-0 \
        libgtk-3-dev \
        libwebkit2gtk-4.0-dev
elif command -v dnf &> /dev/null; then
    # Fedora/RHEL
    sudo dnf install -y \
        gtk3 \
        webkit2gtk4.0 \
        glib2 \
        gdk-pixbuf2 \
        pango \
        cairo \
        atk \
        at-spi2-atk
elif command -v pacman &> /dev/null; then
    # Arch Linux
    sudo pacman -S --noconfirm \
        gtk3 \
        webkit2gtk \
        glib2 \
        gdk-pixbuf2 \
        pango \
        cairo \
        atk
else
    echo "Unsupported package manager. Please install manually:"
    echo "  - libgtk-3-0"
    echo "  - libwebkit2gtk-4.0-37"
    exit 1
fi

echo "Dependencies installed successfully!"
