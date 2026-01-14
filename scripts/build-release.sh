#!/bin/bash
# Build script for creating multi-platform release binaries

set -e

VERSION=${1:-dev}
OUTPUT_DIR="dist"

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

log() {
    echo -e "${BLUE}[BUILD]${NC} $1"
}

success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

# Create output directory
mkdir -p "$OUTPUT_DIR"
rm -rf "${OUTPUT_DIR:?}/"*

log "Building Enterprise Shield v${VERSION}"

# Build for multiple platforms
PLATFORMS=(
    "linux/amd64"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
)

for PLATFORM in "${PLATFORMS[@]}"; do
    IFS='/' read -r GOOS GOARCH <<< "$PLATFORM"
    
    OUTPUT_NAME="enterprise-shield-${VERSION}-${GOOS}-${GOARCH}"
    if [ "$GOOS" = "windows" ]; then
        OUTPUT_NAME="${OUTPUT_NAME}.exe"
    fi
    
    log "Building ${GOOS}/${GOARCH}..."
    
    env GOOS="$GOOS" GOARCH="$GOARCH" CGO_ENABLED=0 go build \
        -ldflags "-X main.version=${VERSION} -X main.buildTime=$(date -u '+%Y-%m-%d_%H:%M:%S') -w -s" \
        -o "${OUTPUT_DIR}/${OUTPUT_NAME}" \
        ./cmd/plugin
    
    # Create archive
    pushd "$OUTPUT_DIR" > /dev/null
    
    if [ "$GOOS" = "windows" ]; then
        zip "${OUTPUT_NAME}.zip" "$OUTPUT_NAME"
        sha256sum "${OUTPUT_NAME}.zip" > "${OUTPUT_NAME}.zip.sha256"
        rm "$OUTPUT_NAME"
        success "Created ${OUTPUT_NAME}.zip"
    else
        tar -czf "${OUTPUT_NAME}.tar.gz" "$OUTPUT_NAME"
        sha256sum "${OUTPUT_NAME}.tar.gz" > "${OUTPUT_NAME}.tar.gz.sha256"
        rm "$OUTPUT_NAME"
        success "Created ${OUTPUT_NAME}.tar.gz"
    fi
    
    popd > /dev/null
done

log "Generating checksums file..."
cat "${OUTPUT_DIR}"/*.sha256 > "${OUTPUT_DIR}/checksums.txt"

log "Build complete! Artifacts in ${OUTPUT_DIR}/"
ls -lh "$OUTPUT_DIR"

success "All binaries built successfully for version ${VERSION}"

