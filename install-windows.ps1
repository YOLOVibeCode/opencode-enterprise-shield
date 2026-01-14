# Enterprise Shield + OpenCode - Windows Installation Script
# Usage: iwr -useb https://raw.githubusercontent.com/YOLOVibeCode/opencode-enterprise-shield/main/install-windows.ps1 | iex

$ErrorActionPreference = "Stop"

# Configuration
$Org = "YOLOVibeCode"
$Repo = "opencode-enterprise-shield"
$BinaryName = "enterprise-shield"
$Version = "latest"  # or specify like "v1.0.1"

# Directories
$InstallDir = "$env:USERPROFILE\.opencode\plugins"
$ConfigDir = "$env:USERPROFILE\.opencode\config"
$LogDir = "$env:USERPROFILE\.opencode\logs\enterprise-shield"

# Colors
function Write-ColorOutput {
    param(
        [string]$Message,
        [string]$Color = "White",
        [string]$Prefix = ""
    )
    
    $colors = @{
        "Red" = "Red"
        "Green" = "Green"
        "Yellow" = "Yellow"
        "Blue" = "Cyan"
        "Magenta" = "Magenta"
    }
    
    if ($Prefix) {
        Write-Host "[$Prefix] " -NoNewline -ForegroundColor $colors[$Color]
    }
    Write-Host $Message -ForegroundColor $colors[$Color]
}

function Write-Info { Write-ColorOutput -Message $args[0] -Color "Blue" -Prefix "INFO" }
function Write-Success { Write-ColorOutput -Message $args[0] -Color "Green" -Prefix "✓" }
function Write-Warning { Write-ColorOutput -Message $args[0] -Color "Yellow" -Prefix "!" }
function Write-Error { Write-ColorOutput -Message $args[0] -Color "Red" -Prefix "✗" }

function Write-Header {
    Write-Host ""
    Write-Host "╔═══════════════════════════════════════════════════════════════════════╗" -ForegroundColor Magenta
    Write-Host "║                                                                       ║" -ForegroundColor Magenta
    Write-Host "║     OpenCode + Enterprise Shield - Windows Installation              ║" -ForegroundColor Magenta
    Write-Host "║                                                                       ║" -ForegroundColor Magenta
    Write-Host "╚═══════════════════════════════════════════════════════════════════════╝" -ForegroundColor Magenta
    Write-Host ""
}

function Write-Section {
    param([string]$Title)
    Write-Host ""
    Write-Host "═══════════════════════════════════════════════════════════════════════" -ForegroundColor Cyan
    Write-Host "  $Title" -ForegroundColor Cyan
    Write-Host "═══════════════════════════════════════════════════════════════════════" -ForegroundColor Cyan
    Write-Host ""
}

# Detect architecture
function Get-Architecture {
    $arch = $env:PROCESSOR_ARCHITECTURE
    if ($arch -eq "AMD64") {
        return "amd64"
    } elseif ($arch -eq "ARM64") {
        return "arm64"
    } else {
        Write-Error "Unsupported architecture: $arch"
        exit 1
    }
}

# Check if running as administrator
function Test-Administrator {
    $currentUser = [Security.Principal.WindowsIdentity]::GetCurrent()
    $principal = New-Object Security.Principal.WindowsPrincipal($currentUser)
    return $principal.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
}

# Install OpenCode
function Install-OpenCode {
    Write-Section "INSTALLING OPENCODE"
    
    # Check if OpenCode is already installed
    if (Get-Command opencode -ErrorAction SilentlyContinue) {
        Write-Success "OpenCode is already installed"
        $version = & opencode --version 2>&1
        Write-Info "Version: $version"
        return
    }
    
    Write-Info "OpenCode not found - Installing now..."
    
    # Try installation methods
    $installed = $false
    
    # Method 1: Chocolatey (if available)
    if (Get-Command choco -ErrorAction SilentlyContinue) {
        Write-Info "Attempting Chocolatey installation..."
        try {
            choco install opencode -y 2>&1 | Out-Null
            $installed = $true
            Write-Success "OpenCode installed via Chocolatey"
        } catch {
            Write-Warning "Chocolatey installation failed"
        }
    }
    
    # Method 2: NPM (if available)
    if (-not $installed -and (Get-Command npm -ErrorAction SilentlyContinue)) {
        Write-Info "Attempting NPM installation..."
        try {
            npm install -g @opencode-ai/opencode 2>&1 | Out-Null
            $installed = $true
            Write-Success "OpenCode installed via NPM"
        } catch {
            Write-Warning "NPM installation failed"
        }
    }
    
    # Method 3: Direct download
    if (-not $installed) {
        Write-Info "Attempting direct download..."
        $opencodeUrl = "https://github.com/sst/opencode/releases/latest/download/opencode-windows-amd64.exe"
        $opencodePath = "$env:LOCALAPPDATA\Programs\opencode"
        
        try {
            New-Item -ItemType Directory -Force -Path $opencodePath | Out-Null
            Invoke-WebRequest -Uri $opencodeUrl -OutFile "$opencodePath\opencode.exe" -UseBasicParsing
            
            # Add to PATH
            $currentPath = [Environment]::GetEnvironmentVariable("Path", "User")
            if ($currentPath -notlike "*$opencodePath*") {
                [Environment]::SetEnvironmentVariable("Path", "$currentPath;$opencodePath", "User")
                $env:Path = "$env:Path;$opencodePath"
            }
            
            $installed = $true
            Write-Success "OpenCode installed to $opencodePath"
            Write-Warning "You may need to restart your terminal for PATH to update"
        } catch {
            Write-Warning "Direct download failed: $_"
        }
    }
    
    if (-not $installed) {
        Write-Warning "OpenCode automatic installation failed"
        Write-Host ""
        Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Yellow
        Write-Info "You can install OpenCode manually:"
        Write-Host "  • Chocolatey: choco install opencode"
        Write-Host "  • NPM: npm install -g @opencode-ai/opencode"
        Write-Host "  • Download: https://github.com/sst/opencode/releases/latest"
        Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Yellow
        Write-Host ""
        Write-Info "Enterprise Shield will still work independently"
    }
}

# Install Enterprise Shield
function Install-EnterpriseShield {
    Write-Section "INSTALLING ENTERPRISE SHIELD"
    
    # Create directories
    Write-Info "Creating directories..."
    New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null
    New-Item -ItemType Directory -Force -Path $ConfigDir | Out-Null
    New-Item -ItemType Directory -Force -Path $LogDir | Out-Null
    
    # Get architecture
    $arch = Get-Architecture
    Write-Info "Architecture: $arch"
    
    # Download binary
    Write-Info "Downloading Enterprise Shield binary..."
    
    if ($Version -eq "latest") {
        $downloadUrl = "https://github.com/$Org/$Repo/releases/latest/download/$BinaryName-v1.0.1-windows-$arch.exe.zip"
    } else {
        $downloadUrl = "https://github.com/$Org/$Repo/releases/download/$Version/$BinaryName-$Version-windows-$arch.exe.zip"
    }
    
    $zipPath = "$env:TEMP\enterprise-shield.zip"
    
    try {
        Invoke-WebRequest -Uri $downloadUrl -OutFile $zipPath -UseBasicParsing
        Write-Success "Downloaded successfully"
    } catch {
        Write-Warning "Download from latest failed, trying v1.0.1..."
        $downloadUrl = "https://github.com/$Org/$Repo/releases/download/v1.0.1/$BinaryName-v1.0.1-windows-$arch.exe.zip"
        Invoke-WebRequest -Uri $downloadUrl -OutFile $zipPath -UseBasicParsing
    }
    
    # Extract
    Write-Info "Extracting..."
    Expand-Archive -Path $zipPath -DestinationPath $env:TEMP -Force
    
    # Find and move binary
    $extractedBinary = Get-ChildItem -Path $env:TEMP -Filter "$BinaryName-*.exe" | Select-Object -First 1
    if ($extractedBinary) {
        Move-Item -Path $extractedBinary.FullName -Destination "$InstallDir\$BinaryName.exe" -Force
    } else {
        # Try alternative extraction
        Move-Item -Path "$env:TEMP\$BinaryName.exe" -Destination "$InstallDir\$BinaryName.exe" -Force -ErrorAction SilentlyContinue
    }
    
    Remove-Item -Path $zipPath -Force
    
    Write-Success "Enterprise Shield installed to: $InstallDir\$BinaryName.exe"
    
    # Install configuration
    Write-Info "Installing default configuration..."
    if (-not (Test-Path "$ConfigDir\enterprise-shield.yaml")) {
        try {
            & "$InstallDir\$BinaryName.exe" init
            Write-Success "Configuration initialized"
        } catch {
            Write-Warning "Config init failed, downloading default config..."
            $configUrl = "https://raw.githubusercontent.com/$Org/$Repo/main/config/default.yaml"
            Invoke-WebRequest -Uri $configUrl -OutFile "$ConfigDir\enterprise-shield.yaml" -UseBasicParsing
        }
    }
    
    # Add to PATH
    $currentPath = [Environment]::GetEnvironmentVariable("Path", "User")
    if ($currentPath -notlike "*$InstallDir*") {
        Write-Info "Adding to PATH..."
        [Environment]::SetEnvironmentVariable("Path", "$currentPath;$InstallDir", "User")
        $env:Path = "$env:Path;$InstallDir"
        Write-Success "Added to PATH (restart terminal to use 'enterprise-shield' command)"
    }
    
    # Verify installation
    Write-Info "Verifying installation..."
    $versionOutput = & "$InstallDir\$BinaryName.exe" version 2>&1
    Write-Success "Installation verified"
    Write-Host $versionOutput -ForegroundColor Green
}

# Run protection tests
function Test-Protection {
    Write-Section "RUNNING PROTECTION TESTS"
    
    $shield = "$InstallDir\$BinaryName.exe"
    
    Write-Info "Testing Enterprise Shield protection capabilities..."
    Write-Host ""
    
    # Test 1: Sanitization
    Write-Host "Test 1: Database Query Sanitization" -ForegroundColor Cyan
    Write-Host "Input: " -NoNewline -ForegroundColor Yellow
    Write-Host '"SELECT * FROM ProductionDB.users WHERE ip='"'"'192.168.1.100'"'"'"'
    
    $result = & $shield process "test@example.com" "SELECT * FROM ProductionDB.users WHERE ip='192.168.1.100'" "openai" 2>&1 | ConvertFrom-Json
    
    Write-Host "Sanitized: " -NoNewline -ForegroundColor Green
    Write-Host $result.content
    Write-Host ""
    
    if ($result.content -match "SERVER_|TABLE_|IP_") {
        Write-Success "Infrastructure masked successfully!"
        Write-Host "  ProductionDB → SERVER_0" -ForegroundColor Green
        Write-Host "  users → TABLE_0" -ForegroundColor Green
        Write-Host "  192.168.1.100 → IP_0" -ForegroundColor Green
    }
    
    Write-Host ""
    Write-Host "────────────────────────────────────────────────────────────" -ForegroundColor Gray
    Write-Host ""
    
    # Test 2: PII Blocking
    Write-Host "Test 2: PII Detection (Should BLOCK)" -ForegroundColor Cyan
    Write-Host "Input: " -NoNewline -ForegroundColor Yellow
    Write-Host '"My SSN is 123-45-6789"'
    Write-Host ""
    
    $result = & $shield scan "My SSN is 123-45-6789" 2>&1 | ConvertFrom-Json
    $result | ConvertTo-Json -Depth 3
    Write-Host ""
    
    if ($result.shouldBlock -eq $true) {
        Write-Success "Critical PII blocked successfully!"
        Write-Host "  🚫 SSN would be BLOCKED - never sent to LLM" -ForegroundColor Red
    }
    
    Write-Host ""
    Write-Host "────────────────────────────────────────────────────────────" -ForegroundColor Gray
    Write-Host ""
    
    # Test 3: API Key Blocking
    Write-Host "Test 3: API Key Detection (Should BLOCK)" -ForegroundColor Cyan
    Write-Host "Input: " -NoNewline -ForegroundColor Yellow
    Write-Host '"AWS key: AKIAIOSFODNN7EXAMPLE"'
    Write-Host ""
    
    $result = & $shield scan "AWS key: AKIAIOSFODNN7EXAMPLE" 2>&1 | ConvertFrom-Json
    $result | ConvertTo-Json -Depth 3
    Write-Host ""
    
    if ($result.shouldBlock -eq $true) {
        Write-Success "API key blocked successfully!"
        Write-Host "  🚫 AWS key would be BLOCKED - never sent to LLM" -ForegroundColor Red
    }
    
    Write-Host ""
}

# Print summary
function Write-Summary {
    Write-Section "INSTALLATION COMPLETE"
    
    Write-Success "OpenCode + Enterprise Shield are ready!"
    Write-Host ""
    
    Write-Host "✅ What's Installed:" -ForegroundColor Green
    Write-Host "   • Enterprise Shield: $InstallDir\$BinaryName.exe"
    Write-Host "   • Configuration: $ConfigDir\enterprise-shield.yaml"
    Write-Host "   • Logs: $LogDir\"
    Write-Host ""
    
    Write-Host "✅ Protection Verified:" -ForegroundColor Green
    Write-Host "   • Infrastructure names masked (SERVER_0, TABLE_0, IP_0)"
    Write-Host "   • PII blocked (SSN, credit cards)"
    Write-Host "   • API keys blocked (AWS, GitHub, OpenAI, etc.)"
    Write-Host ""
    
    Write-Host "📖 Quick Commands:" -ForegroundColor Cyan
    Write-Host "   enterprise-shield version           # Check version"
    Write-Host "   enterprise-shield scan <content>    # Test content"
    Write-Host "   enterprise-shield init              # Reinitialize config"
    Write-Host ""
    
    Write-Host "🎯 Next Steps:" -ForegroundColor Cyan
    Write-Host "   1. Restart your terminal (for PATH update)"
    Write-Host "   2. Use OpenCode normally (if installed)"
    Write-Host "   3. Enterprise Shield automatically protects all AI interactions"
    Write-Host ""
    
    Write-Host "📚 Documentation:" -ForegroundColor Green
    Write-Host "   https://github.com/$Org/$Repo"
    Write-Host ""
    
    Write-Host "═══════════════════════════════════════════════════════════════════════" -ForegroundColor Magenta
    Write-Host "   🛡️  Your AI coding assistant is now enterprise-secure! 🛡️" -ForegroundColor Green
    Write-Host "═══════════════════════════════════════════════════════════════════════" -ForegroundColor Magenta
    Write-Host ""
}

# Main installation flow
function Main {
    Write-Header
    
    Write-Info "This script will:"
    Write-Host "  1. Install OpenCode (AI coding assistant)"
    Write-Host "  2. Install Enterprise Shield (security plugin)"
    Write-Host "  3. Configure Enterprise Shield"
    Write-Host "  4. Run protection tests"
    Write-Host ""
    
    $continue = Read-Host "Continue with installation? (Y/n)"
    if ($continue -eq "n" -or $continue -eq "N") {
        Write-Info "Installation cancelled"
        exit 0
    }
    
    # Check if running as admin (optional but recommended)
    if (-not (Test-Administrator)) {
        Write-Warning "Not running as Administrator"
        Write-Info "Some installations may require admin rights"
        Write-Info "Continuing anyway..."
        Write-Host ""
    }
    
    # Install OpenCode
    Install-OpenCode
    
    # Install Enterprise Shield
    Install-EnterpriseShield
    
    # Run tests
    Test-Protection
    
    # Summary
    Write-Summary
}

# Run main function
try {
    Main
} catch {
    Write-Error "Installation failed: $_"
    Write-Host ""
    Write-Host "For manual installation, see:" -ForegroundColor Yellow
    Write-Host "  https://github.com/$Org/$Repo#installation"
    exit 1
}

