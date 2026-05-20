#!/bin/bash
# ManPaSik Flutter app launche
set -euo pipefail

echo "[ManPaSik] Starting Flutter app..."

# Use absolute paths for stable execution.
PROJECT_ROOT="$HOME/Manpasik"
FLUTTER_APP_DIR="$PROJECT_ROOT/frontend/flutter-app"

# Flutter SDK path (WSL native install).
FLUTTER_BIN="$HOME/flutter/bin/flutter"

if [ ! -f "$FLUTTER_BIN" ]; then
    echo "[WARN] $FLUTTER_BIN not found. Trying flutter from PATH."
    if ! command -v flutter >/dev/null 2>&1; then
        echo "[ERROR] flutter command not found. Please verify SDK install."
        exit 1
    else
        FLUTTER_BIN="flutter"
    fi
fi

if [ -d "$FLUTTER_APP_DIR" ]; then
    echo "[INFO] cd $FLUTTER_APP_DIR"
    cd "$FLUTTER_APP_DIR"
else
    echo "[ERROR] Missing directory: $FLUTTER_APP_DIR"
    exit 1
fi

if [ ! -f "pubspec.yaml" ]; then
    echo "[ERROR] pubspec.yaml not found."
    exit 1
fi

configure_wsl_graphics_fallback() {
    # Default on WSL: software rendering fallback for stability.
    # Disable with: MANPASIK_WSL_SOFTWARE_RENDER=0 ./run_app.sh
    if ! grep -qiE "(microsoft|wsl)" /proc/version 2>/dev/null; then
        return 0
    fi

    local enable_fallback="${MANPASIK_WSL_SOFTWARE_RENDER:-1}"
    if [ "$enable_fallback" != "1" ]; then
        echo "  [INFO] WSL software fallback disabled (MANPASIK_WSL_SOFTWARE_RENDER=0)"
        return 0
    fi

    export LIBGL_ALWAYS_SOFTWARE="${LIBGL_ALWAYS_SOFTWARE:-1}"
    export GALLIUM_DRIVER="${GALLIUM_DRIVER:-llvmpipe}"
    export MESA_LOADER_DRIVER_OVERRIDE="${MESA_LOADER_DRIVER_OVERRIDE:-llvmpipe}"
    echo "  [INFO] WSL graphics fallback enabled (LIBGL_ALWAYS_SOFTWARE=1, llvmpipe)"
}

check_linux_prerequisites() {
    if ! command -v xclip >/dev/null 2>&1; then
        echo "  [WARN] xclip not installed: CEF clipboard can be limited."
        echo "         Install: sudo apt-get update && sudo apt-get install -y xclip"
    fi
}

ensure_webview_cef_linux_compat_sources() {
    local linux_dir runner_dir file src dst link_target

    linux_dir="$FLUTTER_APP_DIR/linux"
    runner_dir="$linux_dir/runner"

    if [ ! -d "$runner_dir" ]; then
        return 0
    fi

    for file in main.cc my_application.cc; do
        src="$runner_dir/$file"
        dst="$linux_dir/$file"
        link_target="runner/$file"

        if [ ! -f "$src" ]; then
            continue
        fi

        if [ -L "$dst" ]; then
            if [ "$(readlink "$dst")" = "$link_target" ]; then
                continue
            fi
            rm -f "$dst"
        elif [ -e "$dst" ]; then
            if cmp -s "$src" "$dst"; then
                echo "  [INFO] webview_cef compat source already present: linux/$file"
            else
                cp -f "$src" "$dst"
                echo "  [INFO] Synced webview_cef compat source: linux/$file (runner -> linux)"
            fi
            continue
        fi

        ln -s "$link_target" "$dst"
        echo "  [INFO] Created webview_cef compat link: linux/$file -> $link_target"
    done
}

prepare_webview_cef_linux() {
    local pub_cache pkg_dir arch url version extracted_name prebuilt cef_di

    pub_cache="$HOME/.pub-cache/hosted/pub.dev"
    pkg_dir=$(ls -d "$pub_cache"/webview_cef-* 2>/dev/null | sort -V | tail -n 1 || true)

    if [ -z "$pkg_dir" ]; then
        echo "  [INFO] webview_cef package not found. Skipping CEF bootstrap."
        return 0
    fi

    arch=$(uname -m)
    if [ "$arch" = "aarch64" ]; then
        url="https://cef-builds.spotifycdn.com/cef_binary_130.1.2%2Bg48f3ef6%2Bchromium-130.0.6723.44_linuxarm64.tar.bz2"
        version="cef_binary_130.1.2%2Bg48f3ef6%2Bchromium-130.0.6723.44_linuxarm64.tar.bz2"
        extracted_name="cef_binary_130.1.2+g48f3ef6+chromium-130.0.6723.44_linuxarm64"
    else
        url="https://cef-builds.spotifycdn.com/cef_binary_130.1.2%2Bg48f3ef6%2Bchromium-130.0.6723.44_linux64.tar.bz2"
        version="cef_binary_130.1.2%2Bg48f3ef6%2Bchromium-130.0.6723.44_linux64.tar.bz2"
        extracted_name="cef_binary_130.1.2+g48f3ef6+chromium-130.0.6723.44_linux64"
    fi

    prebuilt="$pkg_dir/linux/prebuilt.zip"
    cef_dir="$pkg_dir/third/cef"

    if [ -f "$cef_dir/cmake/FindCEF.cmake" ] && [ -f "$cef_dir/Release/libcef.so" ]; then
        echo "  [INFO] webview_cef CEF binaries ready"
        return 0
    fi

    if ! command -v curl >/dev/null 2>&1; then
        echo "[ERROR] curl is required for automatic CEF bootstrap."
        echo "        Install: sudo apt install curl"
        exit 1
    fi

    echo "  [INFO] Preparing webview_cef CEF binaries (first run may be large download)."
    mkdir -p "$pkg_dir/linux" "$pkg_dir/third"

    if [ ! -s "$prebuilt" ]; then
        curl -L --fail --retry 3 -o "$prebuilt" "$url"
    fi

    rm -rf "$cef_dir" "$pkg_dir/third/$extracted_name"
    tar -xf "$prebuilt" -C "$pkg_dir/third"

    if [ ! -d "$pkg_dir/third/$extracted_name" ]; then
        echo "[ERROR] Failed to locate extracted CEF directory."
        exit 1
    fi

    mv "$pkg_dir/third/$extracted_name" "$cef_dir"
    printf "%s" "$version" > "$cef_dir/version.txt"
    echo "  [INFO] webview_cef CEF binaries ready"
}

echo "[INFO] Launching Flutter Linux desktop app..."
echo "----------------------------------------------------------------"
echo "  [INFO] Syncing dependencies..."
"$FLUTTER_BIN" pub get >/dev/null
echo "  [INFO] Starting build..."
echo "  [INFO] Initial build can take 1-3 minutes."
echo "  [INFO] GUI window will open after build completes."
configure_wsl_graphics_fallback
check_linux_prerequisites
ensure_webview_cef_linux_compat_sources
prepare_webview_cef_linux
echo "----------------------------------------------------------------"

"$FLUTTER_BIN" run -d linux

echo "----------------------------------------------------------------"
echo "[INFO] App run finished."
echo "----------------------------------------------------------------"
