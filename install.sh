#!/bin/bash
# Enterprise Shield Plugin - Installation Script
# Usage: curl -sSL https://raw.githubusercontent.com/yourorg/enterprise-shield/main/install.sh | bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
REPO="yourorg/opencode-enterprise-shield"
GITHUB_URL="https://github.com/${REPO}"
BINARY_NAME="enterprise-shield"
INSTALL_DIR="${HOME}/.opencode/plugins"
CONFIG_DIR="${HOME}/.opencode/config"
LOG_DIR="${HOME}/.opencode/logs/enterprise-shield"

# Version (default to latest)
VERSION="${VERSION:-latest}"

# Functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

detect_os() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    case "$OS" in
        linux*)     OS='linux';;
        darwin*)    OS='darwin';;
        mingw*)     OS='windows';;
        msys*)      OS='windows';;
        *)          log_error "Unsupported OS: $OS"; exit 1;;
    esac
    echo "$OS"
}

detect_arch() {
    ARCH=$(uname -m)
    case "$ARCH" in
        x86_64)     ARCH='amd64';;
        amd64)      ARCH='amd64';;
        arm64)      ARCH='arm64';;
        aarch64)    ARCH='arm64';;
        *)          log_error "Unsupported architecture: $ARCH"; exit 1;;
    esac
    echo "$ARCH"
}

get_download_url() {
    local os=$1
    local arch=$2
    local version=$3
    
    if [ "$version" = "latest" ]; then
        # Get latest stable release version
        version=$(curl -s "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
        if [ -z "$version" ]; then
            log_warn "Could not fetch latest version, trying dev release"
            version="dev"
        fi
    fi
    
    if [ "$version" = "dev" ]; then
        # Get latest commit SHA for dev builds
        local sha=$(curl -s "https://api.github.com/repos/${REPO}/commits/main" | grep '"sha":' | head -1 | sed -E 's/.*"([^"]+)".*/\1/' | cut -c1-7)
        if [ -n "$sha" ]; then
            version="dev-${sha}"
        fi
    fi
    
    local filename="${BINARY_NAME}-${version}-${os}-${arch}"
    if [ "$os" = "windows" ]; then
        filename="${filename}.exe"
    fi
    
    # Determine which release tag to use
    local release_tag="$version"
    if [[ "$version" == dev-* ]]; then
        release_tag="dev"
    fi
    
    echo "https://github.com/${REPO}/releases/download/${release_tag}/${filename}"
}

check_dependencies() {
    log_info "Checking dependencies..."
    
    local missing_deps=()
    
    if ! command -v curl &> /dev/null; then
        missing_deps+=("curl")
    fi
    
    if ! command -v tar &> /dev/null && ! command -v unzip &> /dev/null; then
        missing_deps+=("tar or unzip")
    fi
    
    if [ ${#missing_deps[@]} -ne 0 ]; then
        log_error "Missing required dependencies: ${missing_deps[*]}"
        exit 1
    fi
    
    log_success "All dependencies satisfied"
}

download_binary() {
    local url=$1
    local output=$2
    
    log_info "Downloading from: $url"
    
    if ! curl -fsSL "$url" -o "$output"; then
        log_error "Failed to download binary"
        
        # If stable release failed, try dev release
        if [[ "$VERSION" != "dev" ]]; then
            log_warn "Stable release not found, trying development build..."
            VERSION="dev"
            url=$(get_download_url "$OS" "$ARCH" "$VERSION")
            log_info "Downloading from: $url"
            
            if ! curl -fsSL "$url" -o "$output"; then
                log_error "Development build also failed"
                log_info "You can build from source instead:"
                log_info "  git clone ${GITHUB_URL}"
                log_info "  cd opencode-enterprise-shield"
                log_info "  make install"
                exit 1
            fi
        else
            exit 1
        fi
    fi
    
    chmod +x "$output"
    log_success "Binary downloaded successfully"
}

create_directories() {
    log_info "Creating directories..."
    
    mkdir -p "$INSTALL_DIR"
    mkdir -p "$CONFIG_DIR"
    mkdir -p "$LOG_DIR"
    
    log_success "Directories created"
}

install_config() {
    local config_file="${CONFIG_DIR}/enterprise-shield.yaml"
    
    if [ -f "$config_file" ]; then
        log_warn "Configuration already exists at: $config_file"
        read -p "Overwrite? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            log_info "Keeping existing configuration"
            return
        fi
    fi
    
    log_info "Installing default configuration..."
    
    # Download default config or create inline
    cat > "$config_file" << 'EOF'
# Enterprise Shield Configuration
enabled: true

session:
  ttl: "8h"
  maxMappings: 10000
  encryption: true

rules:
  - ruleId: "custom_server"
    name: "Custom Server Names"
    pattern: '\b(srv|app|db|web)-[a-z]+-\d{2,4}\b'
    prefix: "SERVER"
    severity: "medium"
    enabled: true
    order: 100

compliance:
  blockOnCritical: true
  detectors:
    - type: "ssn"
      enabled: true
      severity: "critical"
    - type: "credit_card"
      enabled: true
      severity: "critical"
      validateLuhn: true
    - type: "api_key"
      enabled: true
      severity: "critical"

policy:
  defaultAccessLevel: "sanitized_only"
  requireAuth: true

audit:
  enabled: true
  logPath: "~/.opencode/logs/enterprise-shield"
  signEntries: true
  retentionDays: 365
EOF
    
    log_success "Configuration installed at: $config_file"
}

verify_installation() {
    local binary_path="$1"
    
    log_info "Verifying installation..."
    
    if [ ! -f "$binary_path" ]; then
        log_error "Binary not found at: $binary_path"
        exit 1
    fi
    
    if [ ! -x "$binary_path" ]; then
        log_error "Binary is not executable"
        exit 1
    fi
    
    # Test the binary
    if ! "$binary_path" version &> /dev/null; then
        log_error "Binary test failed"
        exit 1
    fi
    
    log_success "Installation verified"
}

add_to_path() {
    local install_dir="$1"
    
    # Check if already in PATH
    if echo "$PATH" | grep -q "$install_dir"; then
        return
    fi
    
    log_info "Adding to PATH..."
    
    local shell_rc=""
    if [ -f "$HOME/.zshrc" ]; then
        shell_rc="$HOME/.zshrc"
    elif [ -f "$HOME/.bashrc" ]; then
        shell_rc="$HOME/.bashrc"
    elif [ -f "$HOME/.bash_profile" ]; then
        shell_rc="$HOME/.bash_profile"
    fi
    
    if [ -n "$shell_rc" ]; then
        echo "" >> "$shell_rc"
        echo "# Enterprise Shield" >> "$shell_rc"
        echo "export PATH=\"\$PATH:$install_dir\"" >> "$shell_rc"
        log_success "Added to PATH in $shell_rc"
        log_warn "Please run: source $shell_rc"
    else
        log_warn "Could not determine shell config file"
        log_info "Please add manually: export PATH=\"\$PATH:$install_dir\""
    fi
}

print_success_message() {
    local binary_path="$1"
    
    echo ""
    echo -e "${GREEN}╔════════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${GREEN}║                                                                ║${NC}"
    echo -e "${GREEN}║         🛡️  Enterprise Shield Installed Successfully!  🛡️       ║${NC}"
    echo -e "${GREEN}║                                                                ║${NC}"
    echo -e "${GREEN}╚════════════════════════════════════════════════════════════════╝${NC}"
    echo ""
    log_success "Binary installed at: $binary_path"
    log_success "Config installed at: ${CONFIG_DIR}/enterprise-shield.yaml"
    log_success "Logs will be written to: $LOG_DIR"
    echo ""
    echo -e "${BLUE}Quick Start:${NC}"
    echo "  1. Test the installation:"
    echo "     $binary_path version"
    echo ""
    echo "  2. Scan content for violations:"
    echo "     $binary_path scan \"Your content here\""
    echo ""
    echo "  3. Process a request:"
    echo "     $binary_path process user@example.com \"Query content\" openai"
    echo ""
    echo -e "${BLUE}Documentation:${NC}"
    echo "  ${GITHUB_URL}"
    echo ""
    echo -e "${YELLOW}Note:${NC} OpenCode integration is automatic if you have OpenCode installed."
    echo ""
}

# Main installation flow
main() {
    echo ""
    echo -e "${BLUE}╔════════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${BLUE}║                                                                ║${NC}"
    echo -e "${BLUE}║              Enterprise Shield Plugin Installer                ║${NC}"
    echo -e "${BLUE}║                    for OpenCode                                ║${NC}"
    echo -e "${BLUE}║                                                                ║${NC}"
    echo -e "${BLUE}╚════════════════════════════════════════════════════════════════╝${NC}"
    echo ""
    
    # Detect system
    OS=$(detect_os)
    ARCH=$(detect_arch)
    log_info "Detected: $OS/$ARCH"
    
    # Check dependencies
    check_dependencies
    
    # Create directories
    create_directories
    
    # Get download URL
    DOWNLOAD_URL=$(get_download_url "$OS" "$ARCH" "$VERSION")
    BINARY_PATH="${INSTALL_DIR}/${BINARY_NAME}"
    
    # Download binary
    download_binary "$DOWNLOAD_URL" "$BINARY_PATH"
    
    # Install configuration
    install_config
    
    # Verify installation
    verify_installation "$BINARY_PATH"
    
    # Add to PATH (optional)
    add_to_path "$INSTALL_DIR"
    
    # Success message
    print_success_message "$BINARY_PATH"
    
    exit 0
}

# Run main function
main "$@"

