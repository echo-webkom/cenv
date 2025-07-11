#!/usr/bin/env bash

REPO="echo-webkom/cenv"

if [ "$OS" = "Windows_NT" ]; then
    target="windows-amd64"
    extension=".zip"
else
    case $(uname -sm) in
    "Darwin x86_64")
        target="darwin-amd64"
        extension=".tar.gz"
        ;;
    "Darwin arm64")
        target="darwin-arm64"
        extension=".tar.gz"
        ;;
    "Linux x86_64")
        target="linux-amd64"
        extension=".tar.gz"
        ;;
    "Linux aarch64")
        target="linux-arm64"
        extension=".tar.gz"
        ;;
    *) target="unknown" ;;
    esac
fi

if [ "$target" = "unknown" ]; then
    echo "Error: Unsupported OS or architecture."
    exit 1
fi

if [ "$extension" == ".zip" ]; then
    tool="unzip"
else
    tool="tar"
fi

latest_release=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$latest_release" ]; then
    echo "Error: Unable to fetch the latest release from $REPO."
    exit 1
fi

bin_dir="$HOME/.local/bin"
mkdir -p "$bin_dir"

bins=("cenv")

for bin in "${bins[@]}"; do
    binary_name="${bin}-${latest_release}-${target}${extension}"
    download_url="https://github.com/$REPO/releases/download/$latest_release/$binary_name"
    archive_path="$bin_dir/${bin}${extension}"
    exe="$bin_dir/$bin"

    echo "Downloading $bin from $download_url..."

    curl --fail --location --progress-bar --output "$archive_path" "$download_url"

    if [ $? -ne 0 ]; then
        echo "Error: Failed to download $bin from $download_url."
        exit 1
    fi

    if [ "$tool" == "unzip" ]; then
        if ! command -v unzip &>/dev/null; then
            echo "Error: 'unzip' command is required to extract the binary."
            exit 1
        fi

        echo "Unzipping $bin..."
        unzip -o "$archive_path" -d "$bin_dir"
    else
        if ! command -v tar &>/dev/null; then
            echo "Error: 'tar' command is required to extract the binary."
            exit 1
        fi

        echo "Extracting $bin..."
        tar -xzf "$archive_path" -C "$bin_dir"
    fi

    if [ ! -f "$exe" ]; then
        echo "Error: Failed to extract $bin from the archive."
        exit 1
    fi

    chmod +x "$exe"
    rm "$archive_path"
done

if [ "$OS" != "Windows_NT" ]; then
    man_dir="$HOME/.local/share/man/man1"
    mkdir -p "$man_dir"

    man_url="https://raw.githubusercontent.com/$REPO/main/cenv.1"
    curl --fail --location --silent --output "$man_dir/cenv.1" "$man_url"

    if [ $? -ne 0 ]; then
        echo "Warning: Failed to install man page"
    fi
fi

echo "Installation completed successfully!"
echo "Run 'cenv --help' to get started"
