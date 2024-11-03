# Installs cenv to /usr/local/bin

REPO="echo-webkom/cenv"
TARGET_DIR="temp_cenv"

OS=$(uname | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

if [ "$ARCH" == "x86_64" ]; then
    ARCH="amd64"
elif [ "$ARCH" == "aarch64" ]; then
    ARCH="arm64"
fi

echo "[LOG] Getting latest release"

# Get all release binary urls
URL=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | \
    grep -oP '"browser_download_url": "\K(.*'"$OS-$ARCH"'.*\.tar\.gz)')

# Remove duplicate results
URL=$(echo $URL | cut -d' ' -f1)

if [ -z "$URL" ]; then
    echo "[ERROR] Could not find a release for OS=$OS and ARCH=$ARCH."
    exit 1
fi

echo "[LOG] Downloading cenv..."
curl -s -L -o cenv.tar.gz "$URL"

echo "[LOG] Creating target dir"
mkdir -p "$TARGET_DIR"

echo "[LOG] Unpacking"
tar -xzf cenv.tar.gz -C "$TARGET_DIR"

echo "[LOG] Downloaded and unpacked to $TARGET_DIR/"

echo "[LOG] Copying to usr/local/bin"
sudo cp $TARGET_DIR/cenv /usr/local/bin

echo "[LOG] Cleanup"
rm cenv.tar.gz
rm -r $TARGET_DIR

echo "[LOG] Done"