#!/bin/bash
# Update Homebrew formula with new version and checksums

set -e

VERSION=$1
if [ -z "$VERSION" ]; then
    echo "Usage: $0 <version>"
    echo "Example: $0 v1.0.1"
    exit 1
fi

FORMULA_FILE="Formula/enterprise-shield.rb"
REPO="YOLOVibeCode/opencode-enterprise-shield"

echo "Updating Homebrew formula for version ${VERSION}..."

# Download checksums from GitHub release
echo "Downloading checksums..."
CHECKSUMS_URL="https://github.com/${REPO}/releases/download/${VERSION}/checksums.txt"
curl -fsSL "$CHECKSUMS_URL" -o /tmp/checksums.txt

# Extract SHA256 for each platform
SHA_DARWIN_AMD64=$(grep "darwin-amd64.tar.gz" /tmp/checksums.txt | awk '{print $1}')
SHA_DARWIN_ARM64=$(grep "darwin-arm64.tar.gz" /tmp/checksums.txt | awk '{print $1}')
SHA_LINUX_AMD64=$(grep "linux-amd64.tar.gz" /tmp/checksums.txt | awk '{print $1}')
SHA_LINUX_ARM64=$(grep "linux-arm64.tar.gz" /tmp/checksums.txt | awk '{print $1}')

echo "Checksums:"
echo "  darwin-amd64: ${SHA_DARWIN_AMD64}"
echo "  darwin-arm64: ${SHA_DARWIN_ARM64}"
echo "  linux-amd64:  ${SHA_LINUX_AMD64}"
echo "  linux-arm64:  ${SHA_LINUX_ARM64}"

# Update formula
echo "Updating formula..."
sed -i.bak "s/version \".*\"/version \"${VERSION#v}\"/" "$FORMULA_FILE"
sed -i.bak "s|REPLACE_WITH_ACTUAL_SHA256_AMD64|${SHA_DARWIN_AMD64}|" "$FORMULA_FILE"
sed -i.bak "s|REPLACE_WITH_ACTUAL_SHA256_ARM64|${SHA_DARWIN_ARM64}|" "$FORMULA_FILE"
sed -i.bak "s|REPLACE_WITH_ACTUAL_SHA256_LINUX_AMD64|${SHA_LINUX_AMD64}|" "$FORMULA_FILE"
sed -i.bak "s|REPLACE_WITH_ACTUAL_SHA256_LINUX_ARM64|${SHA_LINUX_ARM64}|" "$FORMULA_FILE"

# Update version in URLs
sed -i.bak "s|v[0-9]\+\.[0-9]\+\.[0-9]\+|${VERSION}|g" "$FORMULA_FILE"

rm "${FORMULA_FILE}.bak"

echo "✅ Formula updated successfully!"
echo ""
echo "Next steps:"
echo "  1. Review changes: git diff ${FORMULA_FILE}"
echo "  2. Test formula: brew install --build-from-source ${FORMULA_FILE}"
echo "  3. Commit: git add ${FORMULA_FILE} && git commit -m 'Update to ${VERSION}'"
echo "  4. Push to tap repository"

