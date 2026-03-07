#!/usr/bin/env bash

REPO="echo-webkom/cenv"

if [ "$OS" = "Windows_NT" ]; then
    target="x86_64-pc-windows-gnu"
    exe_suffix=".exe"
else
    case $(uname -sm) in
    "Darwin x86_64") target="x86_64-apple-darwin" ;;
    "Darwin arm64")  target="aarch64-apple-darwin" ;;
    "Linux x86_64")  target="x86_64-unknown-linux-gnu" ;;
    *)
        echo "Error: Unsupported OS or architecture."
        exit 1
        ;;
    esac
    exe_suffix=""
fi

if ! command -v tar &>/dev/null; then
    echo "Error: 'tar' is required to extract the binary."
    exit 1
fi

latest_release=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$latest_release" ]; then
    echo "Error: Unable to fetch the latest release from $REPO."
    exit 1
fi

bin_dir="$HOME/.local/bin"
mkdir -p "$bin_dir"

archive="cenv-${latest_release}-${target}.tar.gz"
download_url="https://github.com/$REPO/releases/download/$latest_release/$archive"
archive_path="$bin_dir/$archive"
exe="$bin_dir/cenv${exe_suffix}"

echo "Downloading cenv from $download_url..."

curl --fail --location --progress-bar --output "$archive_path" "$download_url"

if [ $? -ne 0 ]; then
    echo "Error: Failed to download cenv from $download_url."
    exit 1
fi

echo "Extracting cenv..."
tar -xzf "$archive_path" -C "$bin_dir"
rm "$archive_path"

if [ ! -f "$exe" ]; then
    echo "Error: Failed to extract cenv from the archive."
    exit 1
fi

chmod +x "$exe"

echo "Installation completed successfully!"
echo "Run 'cenv --help' to get started"

