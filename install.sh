#!/usr/bin/env bash

set -e

# === Konfiguration ===
REPO="maxischmaxi/calendar-export"
APP_NAME="calendar-export"
INSTALL_DIR="$HOME/.config/$APP_NAME"
BINARY_NAME="$APP_NAME"
GITHUB_API="https://api.github.com/repos/$REPO/releases/latest"

# === Plattform und Architektur erkennen ===
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# ARM Macs haben uname -m = arm64 oder aarch64
if [[ "$ARCH" == "x86_64" ]]; then
    ARCH="amd64"
elif [[ "$ARCH" == "arm64" || "$ARCH" == "aarch64" ]]; then
    ARCH="arm64"
else
    echo "Unsupported architecture: $ARCH"
    exit 1
fi

# Windows Detection (für Git Bash oder WSL)
if [[ "$OS" == "mingw"* || "$OS" == "msys"* || "$OS" == "cygwin"* ]]; then
    OS="windows"
    INSTALL_DIR="$APPDATA\\$APP_NAME"
    BINARY_NAME="$APP_NAME.exe"
fi

# === Letztes Release abfragen ===
echo "📦 Lade neuestes Release von $REPO..."

DOWNLOAD_URL=$(curl -s "$GITHUB_API" | grep "browser_download_url" | grep "$OS-$ARCH" | cut -d '"' -f 4)

if [ -z "$DOWNLOAD_URL" ]; then
    echo "❌ Keine passende Datei für $OS-$ARCH gefunden."
    exit 1
fi

TMP_DIR=$(mktemp -d)
cd "$TMP_DIR"

echo "⬇️  Herunterladen: $DOWNLOAD_URL"
curl -L "$DOWNLOAD_URL" -o "$APP_NAME.tar.gz"

echo "📂 Entpacken..."
tar -xzf "$APP_NAME.tar.gz"

# === Installation ===
mkdir -p "$INSTALL_DIR"
mv "$BINARY_NAME" "$INSTALL_DIR/"

chmod +x "$INSTALL_DIR/$BINARY_NAME"

# === Symlink setzen (nur für Unix-Systeme) ===
if [[ "$OS" != "windows" ]]; then
    echo "🔗 Erstelle Symlink in ~/.local/bin..."
    mkdir -p "$HOME/.local/bin"
    ln -sf "$INSTALL_DIR/$BINARY_NAME" "$HOME/.local/bin/$APP_NAME"
    echo "✅ Fertig! Stelle sicher, dass ~/.local/bin in deinem PATH ist."
else
    echo "✅ Installation abgeschlossen. Die Datei befindet sich unter:"
    echo "$INSTALL_DIR\\$BINARY_NAME"
    echo "Füge diesen Ordner manuell zu deinem PATH hinzu, falls nötig."
fi

# Aufräumen
cd ~
rm -rf "$TMP_DIR"

