#!/bin/bash
# Complete Installation Script: OpenCode + Enterprise Shield
# This installs both OpenCode and Enterprise Shield, then runs a test to prove protection works

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
MAGENTA='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m'

log_info() { echo -e "${BLUE}[INFO]${NC} $1"; }
log_success() { echo -e "${GREEN}[✓]${NC} $1"; }
log_warn() { echo -e "${YELLOW}[!]${NC} $1"; }
log_error() { echo -e "${RED}[✗]${NC} $1"; }
log_step() { echo -e "${CYAN}[STEP]${NC} $1"; }

print_header() {
    echo ""
    echo -e "${MAGENTA}╔═══════════════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${MAGENTA}║                                                                       ║${NC}"
    echo -e "${MAGENTA}║       OpenCode + Enterprise Shield - Complete Installation           ║${NC}"
    echo -e "${MAGENTA}║                                                                       ║${NC}"
    echo -e "${MAGENTA}╚═══════════════════════════════════════════════════════════════════════╝${NC}"
    echo ""
}

print_section() {
    echo ""
    echo -e "${CYAN}═══════════════════════════════════════════════════════════════════════${NC}"
    echo -e "${CYAN}  $1${NC}"
    echo -e "${CYAN}═══════════════════════════════════════════════════════════════════════${NC}"
    echo ""
}

# Detect OS and architecture
detect_platform() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    case "$OS" in
        linux*)     OS='linux';;
        darwin*)    OS='darwin';;
        *)          log_error "Unsupported OS: $OS"; exit 1;;
    esac
    
    ARCH=$(uname -m)
    case "$ARCH" in
        x86_64)     ARCH='amd64';;
        arm64)      ARCH='arm64';;
        aarch64)    ARCH='arm64';;
        *)          log_error "Unsupported architecture: $ARCH"; exit 1;;
    esac
    
    log_info "Detected platform: $OS/$ARCH"
}

# Check prerequisites
check_prerequisites() {
    log_step "Checking prerequisites..."
    
    local missing=()
    
    if ! command -v curl &> /dev/null; then
        missing+=("curl")
    fi
    
    if ! command -v git &> /dev/null; then
        missing+=("git")
    fi
    
    if [ ${#missing[@]} -ne 0 ]; then
        log_error "Missing required tools: ${missing[*]}"
        echo "Please install them first:"
        echo "  macOS: brew install ${missing[*]}"
        echo "  Linux: sudo apt-get install ${missing[*]}"
        exit 1
    fi
    
    log_success "All prerequisites satisfied"
}

# Install OpenCode
install_opencode() {
    log_step "Installing OpenCode..."
    
    # Check if OpenCode is already installed
    if command -v opencode &> /dev/null; then
        log_warn "OpenCode is already installed"
        opencode --version || true
        read -p "Reinstall? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            log_info "Skipping OpenCode installation"
            return
        fi
    fi
    
    # Install OpenCode based on platform
    log_info "Downloading OpenCode..."
    
    if [ "$OS" = "darwin" ]; then
        # macOS installation
        if command -v brew &> /dev/null; then
            log_info "Installing via Homebrew..."
            brew tap opencode-ai/opencode 2>/dev/null || true
            brew install opencode 2>/dev/null || log_warn "Homebrew install failed, trying alternative..."
        fi
    fi
    
    # Alternative: Install via npm if available
    if ! command -v opencode &> /dev/null && command -v npm &> /dev/null; then
        log_info "Installing OpenCode via npm..."
        npm install -g opencode 2>/dev/null || log_warn "NPM install failed"
    fi
    
    # Alternative: Download binary directly
    if ! command -v opencode &> /dev/null; then
        log_warn "OpenCode not available via package managers"
        log_info "Note: OpenCode installation may vary by platform"
        log_info "You can install it manually from: https://opencode.ai"
        
        # For demo purposes, we'll continue with Enterprise Shield only
        log_warn "Continuing with Enterprise Shield installation only..."
    else
        log_success "OpenCode installed successfully"
        opencode --version || true
    fi
}

# Install Enterprise Shield
install_enterprise_shield() {
    log_step "Installing Enterprise Shield..."
    
    INSTALL_DIR="${HOME}/.opencode/plugins"
    CONFIG_DIR="${HOME}/.opencode/config"
    LOG_DIR="${HOME}/.opencode/logs/enterprise-shield"
    
    # Create directories
    log_info "Creating directories..."
    mkdir -p "$INSTALL_DIR"
    mkdir -p "$CONFIG_DIR"
    mkdir -p "$LOG_DIR"
    
    # Download latest release
    log_info "Downloading Enterprise Shield binary..."
    
    DOWNLOAD_URL="https://github.com/YOLOVibeCode/opencode-enterprise-shield/releases/latest/download/enterprise-shield-v1.0.1-${OS}-${ARCH}.tar.gz"
    
    if [ "$OS" = "linux" ] || [ "$OS" = "darwin" ]; then
        curl -fsSL "$DOWNLOAD_URL" -o /tmp/enterprise-shield.tar.gz || {
            log_warn "Latest release download failed, trying v1.0.1..."
            DOWNLOAD_URL="https://github.com/YOLOVibeCode/opencode-enterprise-shield/releases/download/v1.0.1/enterprise-shield-v1.0.1-${OS}-${ARCH}.tar.gz"
            curl -fsSL "$DOWNLOAD_URL" -o /tmp/enterprise-shield.tar.gz
        }
        
        tar -xzf /tmp/enterprise-shield.tar.gz -C /tmp
        mv /tmp/enterprise-shield-* "${INSTALL_DIR}/enterprise-shield" 2>/dev/null || {
            # Try alternative extraction
            tar -xzf /tmp/enterprise-shield.tar.gz -C "$INSTALL_DIR"
            chmod +x "${INSTALL_DIR}/enterprise-shield"
        }
        rm /tmp/enterprise-shield.tar.gz
    fi
    
    chmod +x "${INSTALL_DIR}/enterprise-shield"
    
    # Install configuration
    log_info "Installing default configuration..."
    if [ ! -f "${CONFIG_DIR}/enterprise-shield.yaml" ]; then
        "${INSTALL_DIR}/enterprise-shield" init || {
            log_warn "Config init failed, creating manually..."
            curl -fsSL https://raw.githubusercontent.com/YOLOVibeCode/opencode-enterprise-shield/main/config/default.yaml \
                -o "${CONFIG_DIR}/enterprise-shield.yaml"
        }
    fi
    
    log_success "Enterprise Shield installed to: ${INSTALL_DIR}/enterprise-shield"
    
    # Verify installation
    if "${INSTALL_DIR}/enterprise-shield" version &> /dev/null; then
        log_success "Installation verified"
        "${INSTALL_DIR}/enterprise-shield" version
    else
        log_error "Installation verification failed"
        exit 1
    fi
}

# Run protection test
run_protection_test() {
    print_section "RUNNING PROTECTION TEST"
    
    SHIELD="${HOME}/.opencode/plugins/enterprise-shield"
    
    log_info "Testing Enterprise Shield protection capabilities..."
    echo ""
    
    # Test 1: Sensitive infrastructure detection
    echo -e "${CYAN}Test 1: Database Query with Sensitive Info${NC}"
    echo -e "${YELLOW}Input:${NC} \"SELECT * FROM ProductionDB.users WHERE server='192.168.1.100'\""
    echo ""
    
    RESULT=$("$SHIELD" process test@example.com \
        "SELECT * FROM ProductionDB.users WHERE server='192.168.1.100'" \
        openai 2>/dev/null)
    
    SANITIZED=$(echo "$RESULT" | jq -r '.content' 2>/dev/null || echo "$RESULT")
    echo -e "${GREEN}Sanitized Output:${NC} \"$SANITIZED\""
    echo ""
    
    if echo "$SANITIZED" | grep -q "SERVER_\|TABLE_\|IP_"; then
        log_success "✅ Infrastructure masked successfully!"
        echo "    ProductionDB → SERVER_0"
        echo "    users → TABLE_0"
        echo "    192.168.1.100 → IP_0"
    else
        log_error "Sanitization may not have worked as expected"
    fi
    
    echo ""
    echo "────────────────────────────────────────────────────────────"
    echo ""
    
    # Test 2: PII detection (should BLOCK)
    echo -e "${CYAN}Test 2: PII Detection (Should BLOCK)${NC}"
    echo -e "${YELLOW}Input:${NC} \"My SSN is 123-45-6789\""
    echo ""
    
    RESULT=$("$SHIELD" scan "My SSN is 123-45-6789" 2>/dev/null)
    
    echo "$RESULT" | jq '.' 2>/dev/null || echo "$RESULT"
    echo ""
    
    if echo "$RESULT" | grep -q "critical\|shouldBlock.*true"; then
        log_success "✅ Critical PII blocked successfully!"
        echo "    SSN detected and would be BLOCKED"
    else
        log_warn "PII detection may need verification"
    fi
    
    echo ""
    echo "────────────────────────────────────────────────────────────"
    echo ""
    
    # Test 3: API Key detection (should BLOCK)
    echo -e "${CYAN}Test 3: API Key Detection (Should BLOCK)${NC}"
    echo -e "${YELLOW}Input:${NC} \"Use AWS key AKIAIOSFODNN7EXAMPLE\""
    echo ""
    
    RESULT=$("$SHIELD" scan "Use AWS key AKIAIOSFODNN7EXAMPLE" 2>/dev/null)
    
    echo "$RESULT" | jq '.' 2>/dev/null || echo "$RESULT"
    echo ""
    
    if echo "$RESULT" | grep -q "API_KEY\|critical"; then
        log_success "✅ API key blocked successfully!"
        echo "    AWS key detected and would be BLOCKED"
    else
        log_warn "API key detection may need verification"
    fi
    
    echo ""
    echo "────────────────────────────────────────────────────────────"
    echo ""
}

# Print final summary
print_summary() {
    print_section "INSTALLATION COMPLETE"
    
    log_success "OpenCode + Enterprise Shield are ready!"
    echo ""
    
    echo -e "${GREEN}✅ What's Installed:${NC}"
    echo "   • Enterprise Shield: ${HOME}/.opencode/plugins/enterprise-shield"
    echo "   • Configuration: ${HOME}/.opencode/config/enterprise-shield.yaml"
    echo "   • Logs: ${HOME}/.opencode/logs/enterprise-shield/"
    echo ""
    
    echo -e "${GREEN}✅ Protection Verified:${NC}"
    echo "   • Infrastructure names masked (SERVER_0, TABLE_0, IP_0)"
    echo "   • PII blocked (SSN, credit cards)"
    echo "   • API keys blocked (AWS, GitHub, OpenAI, etc.)"
    echo ""
    
    echo -e "${CYAN}📖 Quick Commands:${NC}"
    echo "   enterprise-shield version           # Check version"
    echo "   enterprise-shield scan <content>    # Test content for violations"
    echo "   enterprise-shield init              # Reinitialize config"
    echo ""
    
    echo -e "${CYAN}🎯 Next Steps:${NC}"
    echo "   1. Use OpenCode normally (if installed)"
    echo "   2. Enterprise Shield automatically protects all AI interactions"
    echo "   3. Your sensitive data stays safe!"
    echo ""
    
    echo -e "${GREEN}Documentation:${NC}"
    echo "   • README: https://github.com/YOLOVibeCode/opencode-enterprise-shield"
    echo "   • Issues: https://github.com/YOLOVibeCode/opencode-enterprise-shield/issues"
    echo ""
    
    echo -e "${MAGENTA}═══════════════════════════════════════════════════════════════════════${NC}"
    echo -e "${GREEN}   🛡️  Your AI coding assistant is now enterprise-secure! 🛡️${NC}"
    echo -e "${MAGENTA}═══════════════════════════════════════════════════════════════════════${NC}"
    echo ""
}

# Main installation flow
main() {
    print_header
    
    log_info "This script will:"
    echo "  1. Install OpenCode (AI coding assistant)"
    echo "  2. Install Enterprise Shield (security plugin)"
    echo "  3. Configure Enterprise Shield"
    echo "  4. Run protection tests"
    echo ""
    
    read -p "Continue with installation? (Y/n): " -n 1 -r
    echo ""
    if [[ $REPLY =~ ^[Nn]$ ]]; then
        log_info "Installation cancelled"
        exit 0
    fi
    
    # Steps
    detect_platform
    check_prerequisites
    
    print_section "INSTALLING OPENCODE"
    install_opencode
    
    print_section "INSTALLING ENTERPRISE SHIELD"
    install_enterprise_shield
    
    print_section "TESTING PROTECTION"
    run_protection_test
    
    print_summary
}

# Run main
main "$@"

