# Installation Guide - All Platforms

## 🚀 One-Command Installation

Enterprise Shield can be installed on **macOS, Linux, and Windows** with a single command.

---

## 🍎 macOS / 🐧 Linux

### Complete Setup (OpenCode + Enterprise Shield)

```bash
curl -sSL https://raw.githubusercontent.com/YOLOVibeCode/opencode-enterprise-shield/main/install-complete.sh | bash
```

### Enterprise Shield Only

```bash
curl -sSL https://raw.githubusercontent.com/YOLOVibeCode/opencode-enterprise-shield/main/install.sh | bash
```

**Time:** 2-3 minutes  
**Includes:** Automatic protection tests

---

## 🪟 Windows

### Complete Setup (OpenCode + Enterprise Shield)

Open **PowerShell** and run:

```powershell
iwr -useb https://raw.githubusercontent.com/YOLOVibeCode/opencode-enterprise-shield/main/install-windows.ps1 | iex
```

### Enterprise Shield Only (Quick)

```powershell
iwr -useb https://raw.githubusercontent.com/YOLOVibeCode/opencode-enterprise-shield/main/install-windows-simple.ps1 | iex
```

**Time:** 2-3 minutes  
**Note:** Restart terminal after installation for PATH update

**Detailed Windows Guide:** See `WINDOWS_INSTALLATION.md`

---

## 📦 Package Managers

### Homebrew (macOS/Linux)

```bash
brew tap YOLOVibeCode/opencode-enterprise-shield
brew install enterprise-shield
```

### NPM (All Platforms)

```bash
npm install -g @yolovibeCode/opencode-enterprise-shield
```

### Go Install (All Platforms)

```bash
go install github.com/YOLOVibeCode/opencode-enterprise-shield/cmd/plugin@latest
```

### Chocolatey (Windows - Future)

```powershell
choco install enterprise-shield
```

### Scoop (Windows - Future)

```powershell
scoop bucket add yolovibeCode https://github.com/YOLOVibeCode/scoop-bucket
scoop install enterprise-shield
```

---

## 🐳 Docker (All Platforms)

```bash
docker pull ghcr.io/yolovibeCode/opencode-enterprise-shield:latest

# Run
docker run ghcr.io/yolovibeCode/opencode-enterprise-shield:latest version

# With local config
docker run -v ~/.opencode:/root/.opencode ghcr.io/yolovibeCode/opencode-enterprise-shield:latest
```

---

## 📥 Manual Download

### All Platforms

1. **Download for your platform:**
   - https://github.com/YOLOVibeCode/opencode-enterprise-shield/releases/latest

2. **Choose your binary:**
   - **macOS (Intel):** `enterprise-shield-v1.0.1-darwin-amd64.tar.gz`
   - **macOS (Apple Silicon):** `enterprise-shield-v1.0.1-darwin-arm64.tar.gz`
   - **Linux (x64):** `enterprise-shield-v1.0.1-linux-amd64.tar.gz`
   - **Linux (ARM):** `enterprise-shield-v1.0.1-linux-arm64.tar.gz`
   - **Windows (x64):** `enterprise-shield-v1.0.1-windows-amd64.exe.zip`

3. **Extract and install:**

**macOS/Linux:**
```bash
tar -xzf enterprise-shield-*.tar.gz
mkdir -p ~/.opencode/plugins
mv enterprise-shield-* ~/.opencode/plugins/enterprise-shield
chmod +x ~/.opencode/plugins/enterprise-shield
enterprise-shield init
```

**Windows (PowerShell):**
```powershell
Expand-Archive enterprise-shield-*.zip
New-Item -ItemType Directory -Force -Path "$env:USERPROFILE\.opencode\plugins"
Move-Item enterprise-shield-*.exe "$env:USERPROFILE\.opencode\plugins\enterprise-shield.exe"
& "$env:USERPROFILE\.opencode\plugins\enterprise-shield.exe" init
```

---

## ✅ Verification

### All Platforms

After installation, verify it works:

**macOS/Linux:**
```bash
enterprise-shield version
enterprise-shield scan "Test with ServerDB01"
```

**Windows (PowerShell):**
```powershell
enterprise-shield version
enterprise-shield scan "Test with ServerDB01"
```

**Expected output:**
```
Enterprise Shield Plugin v1.0.1
Build: #XX
Built: 2026-01-14...
```

---

## 🧪 Run Protection Tests

### macOS/Linux

```bash
curl -sSL https://raw.githubusercontent.com/YOLOVibeCode/opencode-enterprise-shield/main/test-protection.sh | bash
```

### Windows

```powershell
# Download test script
iwr -useb https://raw.githubusercontent.com/YOLOVibeCode/opencode-enterprise-shield/main/test-protection.ps1 -OutFile test.ps1

# Run tests
.\test.ps1
```

**Expected:** 8/8 tests pass ✅

---

## 🗺️ Installation Comparison

| Platform | One-Liner | Package Manager | Docker | Manual |
|----------|-----------|----------------|--------|--------|
| **macOS** | ✅ Bash | ✅ Homebrew | ✅ | ✅ |
| **Linux** | ✅ Bash | ✅ Homebrew | ✅ | ✅ |
| **Windows** | ✅ PowerShell | ⏳ Soon | ✅ | ✅ |

---

## 🎯 Which Method to Choose?

### For End Users

| You have | Use this |
|----------|----------|
| macOS/Linux | `curl -sSL ... \| bash` |
| Windows | `iwr -useb ... \| iex` |
| Docker | `docker pull ...` |

### For Developers

| You prefer | Use this |
|------------|----------|
| Homebrew | `brew install` |
| NPM | `npm install -g` |
| Go | `go install` |

### For Enterprises

| Scenario | Use this |
|----------|----------|
| Mass deployment | One-liner in deployment script |
| Controlled environments | Manual download + verification |
| Container-based | Docker images |

---

## 📊 Platform-Specific Features

### macOS
- ✅ Homebrew support (native)
- ✅ Universal binary (Intel + Apple Silicon)
- ✅ Keychain integration (future)

### Linux
- ✅ Multi-arch (amd64, arm64)
- ✅ Systemd service (future)
- ✅ Package managers (Homebrew, apt, yum - future)

### Windows
- ✅ PowerShell scripts
- ✅ PATH auto-configuration
- ✅ Chocolatey support (future)
- ✅ Scoop support (future)
- ✅ WinGet support (future)

---

## 🔐 Security Notes

### Checksum Verification

**All platforms:**

1. Download checksum file:
   ```
   https://github.com/YOLOVibeCode/opencode-enterprise-shield/releases/download/v1.0.1/enterprise-shield-v1.0.1-PLATFORM-ARCH.tar.gz.sha256
   ```

2. Verify:

**macOS/Linux:**
```bash
sha256sum -c enterprise-shield-v1.0.1-darwin-arm64.tar.gz.sha256
```

**Windows (PowerShell):**
```powershell
$hash = Get-FileHash enterprise-shield-v1.0.1-windows-amd64.exe.zip -Algorithm SHA256
$expected = Get-Content enterprise-shield-v1.0.1-windows-amd64.exe.zip.sha256 -First 1
$hash.Hash -eq $expected.Split()[0]
```

---

## 📞 Platform-Specific Support

### macOS Issues
- Check: Gatekeeper may block unsigned binaries
- Fix: `xattr -d com.apple.quarantine ~/.opencode/plugins/enterprise-shield`

### Linux Issues
- Check: File permissions
- Fix: `chmod +x ~/.opencode/plugins/enterprise-shield`

### Windows Issues
- Check: Execution policy
- Fix: `Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser`

---

## ✅ Success Checklist

After installation on any platform:

- [ ] Binary installed in correct location
- [ ] `enterprise-shield version` works
- [ ] Configuration file exists
- [ ] Protection tests pass (8/8)
- [ ] Ready to use with OpenCode

---

## 🎉 Quick Start After Installation

**All platforms:**

```bash
# Check version
enterprise-shield version

# Test protection
enterprise-shield scan "Test with ServerDB01 and 192.168.1.100"

# Process request
enterprise-shield process your@email.com "Your query" openai
```

---

**Enterprise Shield works on all major platforms!** 🛡️

- **macOS:** ✅ Intel & Apple Silicon
- **Linux:** ✅ x64 & ARM64
- **Windows:** ✅ x64 & ARM64

