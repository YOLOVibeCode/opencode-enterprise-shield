#!/bin/bash
# Enterprise Shield Plugin - Uninstallation Script

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

BINARY_NAME="enterprise-shield"
INSTALL_DIR="${HOME}/.opencode/plugins"
CONFIG_DIR="${HOME}/.opencode/config"
LOG_DIR="${HOME}/.opencode/logs/enterprise-shield"

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

confirm() {
    read -p "$1 (y/N): " -n 1 -r
    echo
    [[ $REPLY =~ ^[Yy]$ ]]
}

main() {
    echo ""
    echo -e "${YELLOW}╔════════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${YELLOW}║         Enterprise Shield Plugin Uninstaller                   ║${NC}"
    echo -e "${YELLOW}╚════════════════════════════════════════════════════════════════╝${NC}"
    echo ""
    
    log_warn "This will remove Enterprise Shield from your system."
    echo ""
    
    if ! confirm "Are you sure you want to continue?"; then
        log_info "Uninstallation cancelled"
        exit 0
    fi
    
    # Remove binary
    if [ -f "${INSTALL_DIR}/${BINARY_NAME}" ]; then
        log_info "Removing binary..."
        rm -f "${INSTALL_DIR}/${BINARY_NAME}"
        log_success "Binary removed"
    else
        log_warn "Binary not found"
    fi
    
    # Ask about configuration
    if [ -f "${CONFIG_DIR}/enterprise-shield.yaml" ]; then
        if confirm "Remove configuration file?"; then
            rm -f "${CONFIG_DIR}/enterprise-shield.yaml"
            log_success "Configuration removed"
        else
            log_info "Configuration kept"
        fi
    fi
    
    # Ask about logs
    if [ -d "$LOG_DIR" ]; then
        if confirm "Remove log files?"; then
            rm -rf "$LOG_DIR"
            log_success "Logs removed"
        else
            log_info "Logs kept"
        fi
    fi
    
    echo ""
    log_success "Enterprise Shield has been uninstalled"
    echo ""
    log_info "To reinstall: curl -sSL https://raw.githubusercontent.com/yourorg/enterprise-shield/main/install.sh | bash"
    echo ""
}

main "$@"

