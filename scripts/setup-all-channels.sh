#!/bin/bash
# Complete Setup Script for All Distribution Channels

set -e

GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

log() { echo -e "${BLUE}[SETUP]${NC} $1"; }
success() { echo -e "${GREEN}[SUCCESS]${NC} $1"; }
warn() { echo -e "${YELLOW}[INFO]${NC} $1"; }

echo ""
echo "╔════════════════════════════════════════════════════════════╗"
echo "║                                                            ║"
echo "║  Complete Distribution Setup for Enterprise Shield        ║"
echo "║                                                            ║"
echo "╚════════════════════════════════════════════════════════════╝"
echo ""

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Already working
success "✅ Already working (no setup needed):"
echo "  1. GitHub Releases"
echo "  2. Install Script (curl | bash)"
echo "  3. Docker (GHCR)"
echo "  4. Go Install"
echo ""

# To be set up
warn "⏸️  To be configured:"
echo "  5. Homebrew (5 minutes)"
echo "  6. NPM (3 minutes)"
echo ""

read -p "Set up Homebrew? (Y/n): " -n 1 -r
echo ""
if [[ ! $REPLY =~ ^[Nn]$ ]]; then
    log "Setting up Homebrew tap..."
    bash "${SCRIPT_DIR}/setup-homebrew.sh"
    echo ""
else
    warn "Skipping Homebrew setup"
    echo ""
fi

read -p "Set up NPM? (Y/n): " -n 1 -r
echo ""
if [[ ! $REPLY =~ ^[Nn]$ ]]; then
    log "Setting up NPM publishing..."
    bash "${SCRIPT_DIR}/setup-npm.sh"
    echo ""
else
    warn "Skipping NPM setup"
    echo ""
fi

echo ""
echo "╔════════════════════════════════════════════════════════════╗"
echo "║                                                            ║"
echo "║              ✅ Setup Complete!                            ║"
echo "║                                                            ║"
echo "╚════════════════════════════════════════════════════════════╝"
echo ""

success "Distribution channels configured!"
echo ""
echo "Users can now install via:"
echo "  • curl -sSL https://raw.githubusercontent.com/YOLOVibeCode/enterprise-shield/main/install.sh | bash"
echo "  • go install github.com/YOLOVibeCode/opencode-enterprise-shield/cmd/plugin@latest"
echo "  • docker pull ghcr.io/yolovibeCode/opencode-enterprise-shield:latest"
echo "  • brew tap YOLOVibeCode/opencode-enterprise-shield && brew install enterprise-shield"
echo "  • npm install -g @yolovibeCode/opencode-enterprise-shield"
echo ""
echo "Next release will auto-publish to all channels! 🚀"
echo ""

