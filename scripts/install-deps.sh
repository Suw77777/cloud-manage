#!/bin/bash
# Install dependencies for Cloud Manage GUI
# Supports: Ubuntu/Debian, Linux Mint, Fedora/RHEL, Arch, openSUSE

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

info()  { echo -e "${GREEN}[INFO]${NC} $1"; }
warn()  { echo -e "${YELLOW}[WARN]${NC} $1"; }
error() { echo -e "${RED}[ERROR]${NC} $1"; }

# Check if running as root
SUDO="sudo"
if [ "$EUID" -eq 0 ]; then
    SUDO=""
fi

# Detect distro family
detect_distro() {
    if [ -f /etc/os-release ]; then
        . /etc/os-release
        echo "$ID"
    elif command -v apt-get &> /dev/null; then
        echo "debian"
    elif command -v dnf &> /dev/null; then
        echo "fedora"
    elif command -v pacman &> /dev/null; then
        echo "arch"
    elif command -v zypper &> /dev/null; then
        echo "opensuse"
    else
        echo "unknown"
    fi
}

# Detect webkit2gtk version available
detect_webkit_version() {
    # Check what's available in the package manager
    if command -v apt-cache &> /dev/null; then
        if apt-cache show libwebkit2gtk-4.1-dev &> /dev/null 2>&1; then
            echo "4.1"
        elif apt-cache show libwebkit2gtk-4.0-dev &> /dev/null 2>&1; then
            echo "4.0"
        else
            echo "none"
        fi
    elif command -v dnf &> /dev/null; then
        if dnf list webkit2gtk4.1-devel &> /dev/null 2>&1; then
            echo "4.1"
        elif dnf list webkit2gtk4.0-devel &> /dev/null 2>&1; then
            echo "4.0"
        else
            echo "none"
        fi
    else
        echo "4.0"  # default fallback
    fi
}

# Install for Debian/Ubuntu/Mint
install_debian() {
    local webkit_ver="$1"
    info "Detected webkit2gtk version: $webkit_ver"

    $SUDO apt-get update -qq

    local packages=(
        build-essential
        pkg-config
        libgtk-3-dev
        libglib2.0-dev
        libgdk-pixbuf2.0-dev
        libpango1.0-dev
        libcairo2-dev
        libatk1.0-dev
        libatk-bridge2.0-dev
    )

    if [ "$webkit_ver" = "4.1" ]; then
        packages+=(libwebkit2gtk-4.1-dev)
        info "Installing webkit2gtk-4.1-dev"
    elif [ "$webkit_ver" = "4.0" ]; then
        packages+=(libwebkit2gtk-4.0-dev)
        info "Installing webkit2gtk-4.0-dev"
    else
        error "No webkit2gtk development package found!"
        echo ""
        echo "Try enabling universe/multiverse repositories:"
        echo "  $SUDO add-apt-repository universe"
        echo "  $SUDO apt-get update"
        echo ""
        echo "Or install manually:"
        echo "  Ubuntu 22.04 / Mint 21: $SUDO apt install libwebkit2gtk-4.0-dev"
        echo "  Ubuntu 24.04 / Mint 22: $SUDO apt install libwebkit2gtk-4.1-dev"
        exit 1
    fi

    info "Installing packages..."
    $SUDO apt-get install -y "${packages[@]}"
}

# Install for Fedora/RHEL
install_fedora() {
    local webkit_ver="$1"
    info "Detected webkit2gtk version: $webkit_ver"

    local packages=(
        gcc
        pkg-config
        gtk3-devel
        glib2-devel
        gdk-pixbuf2-devel
        pango-devel
        cairo-devel
        atk-devel
        at-spi2-atk-devel
    )

    if [ "$webkit_ver" = "4.1" ]; then
        packages+=(webkit2gtk4.1-devel)
    else
        packages+=(webkit2gtk4.0-devel)
    fi

    $SUDO dnf install -y "${packages[@]}"
}

# Install for Arch
install_arch() {
    info "Installing packages..."
    $SUDO pacman -S --noconfirm --needed \
        base-devel \
        pkg-config \
        gtk3 \
        webkit2gtk \
        glib2 \
        gdk-pixbuf2 \
        pango \
        cairo \
        atk
}

# Install for openSUSE
install_opensuse() {
    info "Installing packages..."
    $SUDO zypper install -y \
        gcc \
        pkg-config \
        gtk3-devel \
        webkit2gtk-devel \
        glib2-devel \
        gdk-pixbuf-devel \
        pango-devel \
        cairo-devel \
        atk-devel
}

# Check if dependencies are already satisfied
check_deps() {
    if ! command -v pkg-config &> /dev/null; then
        return 1
    fi
    if ! pkg-config --exists gtk+-3.0 2>/dev/null; then
        return 1
    fi
    if ! pkg-config --exists webkit2gtk-4.0 2>/dev/null && ! pkg-config --exists webkit2gtk-4.1 2>/dev/null; then
        return 1
    fi
    return 0
}

# Main
echo ""
echo "  Cloud 管理小助手 - GUI 依赖安装"
echo "  ================================="
echo ""

# Check if deps already satisfied
if check_deps; then
    info "GUI dependencies are already installed!"
    exit 0
fi

DISTRO=$(detect_distro)
info "Detected distro: $DISTRO"

case "$DISTRO" in
    ubuntu|debian|linuxmint|pop|elementary|zorin)
        WEBKIT_VER=$(detect_webkit_version)
        install_debian "$WEBKIT_VER"
        ;;
    fedora|rhel|centos|rocky|alma)
        WEBKIT_VER=$(detect_webkit_version)
        install_fedora "$WEBKIT_VER"
        ;;
    arch|manjaro|endeavouros)
        install_arch
        ;;
    opensuse|sles)
        install_opensuse
        ;;
    *)
        error "Unsupported distro: $DISTRO"
        echo ""
        echo "Please install these packages manually:"
        echo "  - GTK 3 development files"
        echo "  - WebKit2GTK development files (4.0 or 4.1)"
        echo "  - GLib, GDK-Pixbuf, Pango, Cairo, ATK development files"
        echo ""
        echo "Or use CLI/TUI mode which has no GUI dependencies:"
        echo "  ./cloud-manage --cli ecs list"
        echo "  ./cloud-manage --tui"
        exit 1
        ;;
esac

echo ""
info "Dependencies installed successfully!"
echo ""
echo "You can now run:"
echo "  ./cloud-manage          # Auto-detect mode"
echo "  ./cloud-manage --gui    # Force GUI mode"
echo "  wails dev               # GUI development with hot reload"
echo ""
