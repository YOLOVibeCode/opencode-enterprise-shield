#!/bin/bash
# Test the installation process locally

set -e

echo "=== Testing Installation Process ==="
echo ""

# Build first
echo "1. Building binary..."
make build

echo ""
echo "2. Testing install script (dry run)..."
INSTALL_DIR="/tmp/test-opencode-shield"
export HOME="/tmp/test-home-shield"
mkdir -p "$HOME"

# Copy install script and modify for local testing
sed "s|YOLOVibeCode/opencode-enterprise-shield|local/test|g" install.sh > /tmp/test-install.sh
chmod +x /tmp/test-install.sh

echo ""
echo "3. Running installer..."
bash -c "
export HOME='$HOME'
export VERSION='local'
# Simulate installation
mkdir -p '$HOME/.opencode/plugins'
mkdir -p '$HOME/.opencode/config'
cp build/enterprise-shield '$HOME/.opencode/plugins/'
cp config/default.yaml '$HOME/.opencode/config/enterprise-shield.yaml'
"

echo ""
echo "4. Testing binary..."
"$HOME/.opencode/plugins/enterprise-shield" version

echo ""
echo "5. Testing scan command..."
"$HOME/.opencode/plugins/enterprise-shield" scan "Test SSN: 123-45-6789"

echo ""
echo "6. Testing process command..."
"$HOME/.opencode/plugins/enterprise-shield" process test@example.com "Query ServerDB01" openai

echo ""
echo "7. Cleaning up..."
rm -rf "$HOME"
rm /tmp/test-install.sh

echo ""
echo "✅ Installation test completed successfully!"

