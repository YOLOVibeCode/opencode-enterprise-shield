# Enterprise Shield - Simple Windows Installer
# Usage: iwr -useb https://raw.githubusercontent.com/YOLOVibeCode/opencode-enterprise-shield/main/install-windows-simple.ps1 | iex

# Quick one-liner version
$InstallDir = "$env:USERPROFILE\.opencode\plugins"
New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null

# Download latest release
Write-Host "Downloading Enterprise Shield..." -ForegroundColor Cyan
$url = "https://github.com/YOLOVibeCode/opencode-enterprise-shield/releases/latest/download/enterprise-shield-v1.0.1-windows-amd64.exe.zip"
$zip = "$env:TEMP\enterprise-shield.zip"

Invoke-WebRequest -Uri $url -OutFile $zip -UseBasicParsing
Expand-Archive -Path $zip -DestinationPath $env:TEMP -Force

# Move binary
Move-Item -Path "$env:TEMP\enterprise-shield-*.exe" -Destination "$InstallDir\enterprise-shield.exe" -Force -ErrorAction SilentlyContinue
Move-Item -Path "$env:TEMP\enterprise-shield.exe" -Destination "$InstallDir\enterprise-shield.exe" -Force -ErrorAction SilentlyContinue

# Add to PATH
$currentPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($currentPath -notlike "*$InstallDir*") {
    [Environment]::SetEnvironmentVariable("Path", "$currentPath;$InstallDir", "User")
}

# Initialize config
& "$InstallDir\enterprise-shield.exe" init

Write-Host ""
Write-Host "✅ Enterprise Shield installed!" -ForegroundColor Green
Write-Host "   Location: $InstallDir\enterprise-shield.exe" -ForegroundColor Cyan
Write-Host ""
Write-Host "Test it:" -ForegroundColor Cyan
Write-Host "   enterprise-shield version"
Write-Host ""
Write-Host "⚠️  Restart your terminal to use the 'enterprise-shield' command" -ForegroundColor Yellow
Write-Host ""

