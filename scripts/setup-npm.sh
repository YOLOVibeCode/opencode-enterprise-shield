#!/bin/bash
# NPM Publishing Setup Script

set -e

GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

log() { echo -e "${BLUE}[SETUP]${NC} $1"; }
success() { echo -e "${GREEN}[SUCCESS]${NC} $1"; }
warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
error() { echo -e "${RED}[ERROR]${NC} $1"; }

ORG="YOLOVibeCode"
REPO="opencode-enterprise-shield"
PACKAGE_NAME="@yolovibeCode/opencode-enterprise-shield"

echo ""
echo "╔════════════════════════════════════════════════════════════╗"
echo "║   NPM Publishing Setup for Enterprise Shield              ║"
echo "╚════════════════════════════════════════════════════════════╝"
echo ""

# Step 1: Check if logged in to NPM
log "Checking NPM login status..."
if npm whoami &>/dev/null; then
    NPM_USER=$(npm whoami)
    success "Logged in to NPM as: $NPM_USER"
else
    warn "Not logged in to NPM"
    echo ""
    echo "Please run: npm login"
    echo "Then run this script again."
    exit 1
fi

# Step 2: Check if package name is available
log "Checking if package name is available..."
if npm view "$PACKAGE_NAME" version &>/dev/null; then
    warn "Package $PACKAGE_NAME already exists!"
    echo ""
    echo "Options:"
    echo "1. Use a different name (update package.json)"
    echo "2. If you own this package, you can publish updates"
    echo "3. Use unscoped name: opencode-enterprise-shield"
    echo ""
    read -p "Continue anyway? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 0
    fi
else
    success "Package name is available: $PACKAGE_NAME"
fi

# Step 3: Instructions for getting NPM token
echo ""
log "To enable automated NPM publishing, you need an automation token:"
echo ""
echo "════════════════════════════════════════════════════════════"
echo "  STEP-BY-STEP: Get NPM Automation Token"
echo "════════════════════════════════════════════════════════════"
echo ""
echo "  1. Open in browser:"
echo "     https://www.npmjs.com/settings/${NPM_USER}/tokens"
echo ""
echo "  2. Click 'Generate New Token' → 'Classic Token'"
echo ""
echo "  3. Select token type: 'Automation' (important!)"
echo ""
echo "  4. Click 'Generate Token'"
echo ""
echo "  5. Copy the token (starts with 'npm_...')"
echo ""
echo "════════════════════════════════════════════════════════════"
echo ""
read -p "Press Enter when you have copied your NPM token..."
echo ""

# Step 4: Add token to GitHub
log "Adding NPM_TOKEN to GitHub repository..."
echo ""
echo "Paste your NPM token (it will be hidden):"
read -s NPM_TOKEN
echo ""

if [ -z "$NPM_TOKEN" ]; then
    error "No token provided"
    exit 1
fi

# Set the secret using GitHub CLI
log "Setting GitHub secret..."
echo "$NPM_TOKEN" | gh secret set NPM_TOKEN --repo "${ORG}/${REPO}"

success "✅ NPM_TOKEN secret added to GitHub!"

# Step 5: Test the workflow
echo ""
log "Testing NPM publish workflow..."
echo ""
read -p "Trigger test publish of v1.0.0 to NPM? (y/N): " -n 1 -r
echo ""

if [[ $REPLY =~ ^[Yy]$ ]]; then
    log "Triggering NPM publish workflow..."
    gh workflow run publish-npm.yml \
        --repo "${ORG}/${REPO}" \
        -f version=v1.0.0
    
    echo ""
    success "Workflow triggered!"
    echo ""
    echo "Monitor progress:"
    echo "  gh run watch --repo ${ORG}/${REPO}"
    echo ""
    echo "Or visit:"
    echo "  https://github.com/${ORG}/${REPO}/actions"
    echo ""
    echo "After completion, verify:"
    echo "  npm view ${PACKAGE_NAME} version"
    echo ""
fi

echo ""
echo "════════════════════════════════════════════════════════════"
echo ""
success "✅ NPM setup complete!"
echo ""
echo "Future releases will automatically publish to:"
echo "  https://www.npmjs.com/package/${PACKAGE_NAME}"
echo ""
echo "Users can install with:"
echo "  npm install -g ${PACKAGE_NAME}"
echo ""
echo "════════════════════════════════════════════════════════════"
echo ""

