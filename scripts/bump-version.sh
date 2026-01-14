#!/bin/bash
# Bump version script - for manual version bumps

set -e

BUMP_TYPE=${1:-patch}

# Read current version
CURRENT_VERSION=$(cat VERSION)
IFS='.' read -r MAJOR MINOR PATCH <<< "$CURRENT_VERSION"

# Bump based on type
case "$BUMP_TYPE" in
    major)
        MAJOR=$((MAJOR + 1))
        MINOR=0
        PATCH=0
        ;;
    minor)
        MINOR=$((MINOR + 1))
        PATCH=0
        ;;
    patch)
        PATCH=$((PATCH + 1))
        ;;
    *)
        echo "Usage: $0 [major|minor|patch]"
        echo "Current version: $CURRENT_VERSION"
        exit 1
        ;;
esac

NEW_VERSION="${MAJOR}.${MINOR}.${PATCH}"

echo "Bumping version: ${CURRENT_VERSION} → ${NEW_VERSION}"
echo "$NEW_VERSION" > VERSION

echo "Version bumped to ${NEW_VERSION}"
echo ""
echo "Next steps:"
echo "  git add VERSION"
echo "  git commit -m 'chore: Bump version to ${NEW_VERSION}'"
echo "  git push origin main"
echo ""
echo "Or create a release:"
echo "  git tag -a v${NEW_VERSION} -m 'Release v${NEW_VERSION}'"
echo "  git push origin v${NEW_VERSION}"

