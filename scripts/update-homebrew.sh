#!/bin/bash

# Script to update Homebrew formula with new version and checksums
# Usage: ./scripts/update-homebrew.sh v1.0.0

set -e

if [ $# -eq 0 ]; then
    echo "Usage: $0 <version>"
    echo "Example: $0 v1.0.0"
    exit 1
fi

VERSION=$1
VERSION_NO_V=${VERSION#v}  # Remove 'v' prefix

echo "Updating Homebrew formula for version $VERSION"

# Download checksums from GitHub release
CHECKSUMS_URL="https://github.com/deployaja/deployaja-cli/releases/download/$VERSION/checksums.txt"
echo "Downloading checksums from $CHECKSUMS_URL"

# Create temp file for checksums
TEMP_CHECKSUMS=$(mktemp)
curl -sL "$CHECKSUMS_URL" > "$TEMP_CHECKSUMS"

# Extract checksums for each platform
SHA256_DARWIN_ARM64=$(grep "aja-darwin-arm64.tar.gz" "$TEMP_CHECKSUMS" | cut -d' ' -f1)
SHA256_DARWIN_AMD64=$(grep "aja-darwin-amd64.tar.gz" "$TEMP_CHECKSUMS" | cut -d' ' -f1)
SHA256_LINUX_ARM64=$(grep "aja-linux-arm64.tar.gz" "$TEMP_CHECKSUMS" | cut -d' ' -f1)
SHA256_LINUX_AMD64=$(grep "aja-linux-amd64.tar.gz" "$TEMP_CHECKSUMS" | cut -d' ' -f1)

echo "SHA256 checksums:"
echo "  Darwin ARM64: $SHA256_DARWIN_ARM64"
echo "  Darwin AMD64: $SHA256_DARWIN_AMD64"
echo "  Linux ARM64:  $SHA256_LINUX_ARM64"
echo "  Linux AMD64:  $SHA256_LINUX_AMD64"

# Update the Homebrew formula
FORMULA_FILE="Formula/aja.rb"

# Update version
sed -i.bak "s/version \".*\"/version \"$VERSION_NO_V\"/" "$FORMULA_FILE"

# Update checksums
sed -i.bak "s/PLACEHOLDER_SHA256_ARM64/$SHA256_DARWIN_ARM64/" "$FORMULA_FILE"
sed -i.bak "s/PLACEHOLDER_SHA256_AMD64/$SHA256_DARWIN_AMD64/" "$FORMULA_FILE"
sed -i.bak "s/PLACEHOLDER_SHA256_LINUX_ARM64/$SHA256_LINUX_ARM64/" "$FORMULA_FILE"
sed -i.bak "s/PLACEHOLDER_SHA256_LINUX_AMD64/$SHA256_LINUX_AMD64/" "$FORMULA_FILE"

# Remove backup file
rm "$FORMULA_FILE.bak"

# Clean up
rm "$TEMP_CHECKSUMS"

echo "✅ Updated $FORMULA_FILE with version $VERSION_NO_V and checksums"
echo ""
echo "Next steps:"
echo "1. Review the changes in $FORMULA_FILE"
echo "2. Test the formula: brew install --build-from-source $FORMULA_FILE"
echo "3. Commit and push to your homebrew-tap repository"
echo "4. Or submit to homebrew-core if ready for wider distribution" 