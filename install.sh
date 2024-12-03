#!/usr/bin/env bash

REPO="echo-webkom/cenv"

if [ "$OS" = "Windows_NT" ]; then
    target="windows-amd64"
else
    case $(uname -sm) in
        "Darwin x86_64") target="darwin-amd64" ;;
        "Darwin arm64") target="darwin-arm64" ;;
        "Linux x86_64") target="linux-amd64" ;;
        "Linux aarch64") target="linux-arm" ;;
        *) target="unknown" ;;
    esac
fi

if [ "$target" = "unknown" ]; then
    echo "Error: Unsupported OS or architecture."
    exit 1
fi

latest_release=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$latest_release" ]; then
    echo "Error: Unable to fetch the latest release from $REPO."
    exit 1
fi

bin_dir="$HOME/.local/bin"
mkdir -p "$bin_dir"

binary_name="cenv-${latest_release}-${target}.tar.gz"
download_url="https://github.com/$REPO/releases/download/$latest_release/$binary_name"
archive_path="$bin_dir/cenv.tar.gz"
exe="$bin_dir/cenv"
install_exe="$bin_dir/cenv-install"

curl --fail --location --progress-bar --output "$archive_path" "$download_url"

if [ $? -ne 0 ]; then
    echo "Error: Failed to download from $download_url."
    exit 1
fi

if ! command -v tar &> /dev/null; then
    echo "Error: 'tar' command is required to extract the binary."
    exit 1
fi

tar -xzf "$archive_path" -C "$bin_dir"

if [ ! -f "$exe" ]; then
    echo "Error: Failed to extract the binary from the archive."
    exit 1
fi

chmod +x "$exe"
chmod +x "$install_exe"

rm "$archive_path"

echo "Installation completed successfully!"
echo "Run 'cenv --help' to get started."
