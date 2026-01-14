# Windows Installation Guide

## 🚀 Quick Install (PowerShell One-Liner)

### Method 1: Complete Installation (Recommended)

Open **PowerShell** and run:

```powershell
iwr -useb https://raw.githubusercontent.com/YOLOVibeCode/opencode-enterprise-shield/main/install-windows.ps1 | iex
```

**This installs:**
- ✅ OpenCode (tries Chocolatey, NPM, direct download)
- ✅ Enterprise Shield
- ✅ Configuration
- ✅ Runs 3 protection tests
- ✅ Adds to PATH

**Time:** 2-3 minutes

---

### Method 2: Enterprise Shield Only (Quick)

If you already have OpenCode:

```powershell
iwr -useb https://raw.githubusercontent.com/YOLOVibeCode/opencode-enterprise-shield/main/install-windows-simple.ps1 | iex
```

**Time:** 30 seconds

---

## 📋 Step-by-Step Manual Installation

### Prerequisites

- Windows 10/11
- PowerShell 5.1 or later
- Internet connection

### Step 1: Download Binary

1. Go to: https://github.com/YOLOVibeCode/opencode-enterprise-shield/releases/latest

2. Download for your architecture:
   - **Most Windows PCs:** `enterprise-shield-v1.0.1-windows-amd64.exe.zip`
   - **ARM Windows:** `enterprise-shield-v1.0.1-windows-arm64.exe.zip`

### Step 2: Extract and Install

```powershell
# Extract the zip
Expand-Archive -Path "Downloads\enterprise-shield-v1.0.1-windows-amd64.exe.zip" -DestinationPath "Downloads\"

# Create plugins directory
New-Item -ItemType Directory -Force -Path "$env:USERPROFILE\.opencode\plugins"

# Move binary
Move-Item -Path "Downloads\enterprise-shield-v1.0.1-windows-amd64.exe" `
          -Destination "$env:USERPROFILE\.opencode\plugins\enterprise-shield.exe"
```

### Step 3: Add to PATH

```powershell
# Add to user PATH
$installDir = "$env:USERPROFILE\.opencode\plugins"
$currentPath = [Environment]::GetEnvironmentVariable("Path", "User")
[Environment]::SetEnvironmentVariable("Path", "$currentPath;$installDir", "User")
```

### Step 4: Initialize Configuration

```powershell
& "$env:USERPROFILE\.opencode\plugins\enterprise-shield.exe" init
```

### Step 5: Verify

```powershell
# Restart PowerShell, then:
enterprise-shield version
```

---

## 🎯 Alternative Installation Methods

### Via Chocolatey (If Available)

```powershell
# Install Chocolatey first (if not installed)
Set-ExecutionPolicy Bypass -Scope Process -Force
iex ((New-Object System.Net.WebClient).DownloadString('https://chocolatey.org/install.ps1'))

# Then install (future - not yet published)
choco install enterprise-shield
```

### Via Scoop (If Available)

```powershell
# Install Scoop first (if not installed)
iwr -useb get.scoop.sh | iex

# Add bucket (future)
scoop bucket add yolovibeCode https://github.com/YOLOVibeCode/scoop-bucket
scoop install enterprise-shield
```

### Via WinGet (Windows Package Manager)

```powershell
# Future - once published to winget
winget install YOLOVibeCode.EnterpriseShield
```

---

## 🧪 Testing on Windows

### Run Protection Tests

```powershell
# Method 1: Download and run test script
iwr -useb https://raw.githubusercontent.com/YOLOVibeCode/opencode-enterprise-shield/main/test-protection.ps1 -OutFile test.ps1
.\test.ps1

# Method 2: Manual tests
enterprise-shield scan "SSN: 123-45-6789"
enterprise-shield process test@example.com "Query ServerDB01" openai
```

---

## 🔧 Troubleshooting

### "Cannot be loaded because running scripts is disabled"

**Solution:** Enable script execution

```powershell
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
```

### "enterprise-shield is not recognized"

**Solution 1:** Restart PowerShell (PATH not updated yet)

**Solution 2:** Use full path temporarily
```powershell
& "$env:USERPROFILE\.opencode\plugins\enterprise-shield.exe" version
```

**Solution 3:** Manually add to PATH
```powershell
$env:Path += ";$env:USERPROFILE\.opencode\plugins"
```

### "Download failed"

**Solution:** Download manually
1. Visit: https://github.com/YOLOVibeCode/opencode-enterprise-shield/releases/latest
2. Download Windows .zip file
3. Extract to `%USERPROFILE%\.opencode\plugins\`
4. Rename to `enterprise-shield.exe`

---

## 🎯 Quick Reference

### Installation Commands

```powershell
# Complete install (OpenCode + Enterprise Shield)
iwr -useb https://raw.githubusercontent.com/YOLOVibeCode/opencode-enterprise-shield/main/install-windows.ps1 | iex

# Enterprise Shield only
iwr -useb https://raw.githubusercontent.com/YOLOVibeCode/opencode-enterprise-shield/main/install-windows-simple.ps1 | iex
```

### Usage Commands

```powershell
# Check version
enterprise-shield version

# Scan for violations
enterprise-shield scan "Your content here"

# Process request
enterprise-shield process user@company.com "Your query" openai

# Initialize config
enterprise-shield init
```

### Locations

| Item | Windows Path |
|------|--------------|
| **Binary** | `%USERPROFILE%\.opencode\plugins\enterprise-shield.exe` |
| **Config** | `%USERPROFILE%\.opencode\config\enterprise-shield.yaml` |
| **Logs** | `%USERPROFILE%\.opencode\logs\enterprise-shield\` |

---

## 📊 Platform Compatibility

| Windows Version | Status | Tested |
|----------------|--------|--------|
| Windows 11 | ✅ Supported | Yes |
| Windows 10 (1809+) | ✅ Supported | Yes |
| Windows Server 2019+ | ✅ Supported | Yes |
| Windows 10 (older) | ⚠️ May work | Not tested |

| Architecture | Status |
|--------------|--------|
| x64 (amd64) | ✅ Supported |
| ARM64 | ✅ Supported |

---

## 🎉 Success Indicators

After installation, you should see:

✅ Binary at: `%USERPROFILE%\.opencode\plugins\enterprise-shield.exe`  
✅ Config at: `%USERPROFILE%\.opencode\config\enterprise-shield.yaml`  
✅ `enterprise-shield version` works (after terminal restart)  
✅ Protection tests pass  

---

## 📞 Need Help?

- **Documentation:** https://github.com/YOLOVibeCode/opencode-enterprise-shield
- **Issues:** https://github.com/YOLOVibeCode/opencode-enterprise-shield/issues
- **Windows-specific:** Tag issue with `windows` label

---

**Enterprise Shield works great on Windows!** 🛡️

