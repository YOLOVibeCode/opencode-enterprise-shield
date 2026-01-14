#!/bin/bash
# Automated Homebrew Tap Setup Script

set -e

GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

log() { echo -e "${BLUE}[SETUP]${NC} $1"; }
success() { echo -e "${GREEN}[SUCCESS]${NC} $1"; }
warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }

ORG="YOLOVibeCode"
TAP_REPO="homebrew-opencode-enterprise-shield"
MAIN_REPO_DIR="/Users/admin/Dev/YOLOProjects/opencode-enterprise-shield"

echo ""
echo "╔════════════════════════════════════════════════════════════╗"
echo "║   Homebrew Tap Setup for Enterprise Shield                ║"
echo "╚════════════════════════════════════════════════════════════╝"
echo ""

# Check if tap repo already exists
log "Checking if tap repository exists..."
if gh repo view "${ORG}/${TAP_REPO}" &>/dev/null; then
    warn "Repository ${ORG}/${TAP_REPO} already exists!"
    read -p "Delete and recreate? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        log "Deleting existing repository..."
        gh repo delete "${ORG}/${TAP_REPO}" --yes
    else
        log "Using existing repository"
    fi
fi

# Create tap repository
if ! gh repo view "${ORG}/${TAP_REPO}" &>/dev/null; then
    log "Creating Homebrew tap repository..."
    gh repo create "${ORG}/${TAP_REPO}" \
        --public \
        --description "Homebrew tap for Enterprise Shield - OpenCode security plugin" \
        --homepage "https://github.com/${ORG}/opencode-enterprise-shield"
    
    success "Repository created: https://github.com/${ORG}/${TAP_REPO}"
fi

# Clone tap repository
TEMP_DIR="/tmp/${TAP_REPO}-setup"
rm -rf "$TEMP_DIR"
mkdir -p "$TEMP_DIR"

log "Cloning tap repository..."
git clone "https://github.com/${ORG}/${TAP_REPO}.git" "$TEMP_DIR"
cd "$TEMP_DIR"

# Create Formula directory
log "Setting up Formula directory..."
mkdir -p Formula

# Copy formula from main repository
log "Copying formula..."
cp "${MAIN_REPO_DIR}/Formula/enterprise-shield.rb" Formula/

# Create README
log "Creating README..."
cat > README.md << 'README'
# Homebrew Tap for Enterprise Shield

Enterprise-grade security plugin for OpenCode AI assistants.

## Installation

```bash
brew tap YOLOVibeCode/opencode-enterprise-shield
brew install enterprise-shield
```

## Updating

```bash
brew update
brew upgrade enterprise-shield
```

## Features

- Automatically sanitizes sensitive data before sending to AI
- Blocks PII (SSN, credit cards, API keys)
- Session-based mapping with encryption
- Audit logging for compliance
- Works with OpenAI, Anthropic, Azure OpenAI, Google AI

## Documentation

See main repository: https://github.com/YOLOVibeCode/opencode-enterprise-shield

## Uninstalling

```bash
brew uninstall enterprise-shield
brew untap YOLOVibeCode/opencode-enterprise-shield
```

## License

MIT
README

# Commit and push
log "Committing and pushing to GitHub..."
git add .
git commit -m "Initial Homebrew formula for Enterprise Shield v1.0.0"
git push origin main

success "✅ Homebrew tap setup complete!"
echo ""
echo "════════════════════════════════════════════════════════════"
echo ""
success "Users can now install with:"
echo "  brew tap ${ORG}/opencode-enterprise-shield"
echo "  brew install enterprise-shield"
echo ""
success "Formula will auto-update on future releases!"
echo ""
echo "View tap at: https://github.com/${ORG}/${TAP_REPO}"
echo ""
echo "════════════════════════════════════════════════════════════"
echo ""

# Clean up
rm -rf "$TEMP_DIR"

