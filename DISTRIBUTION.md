# Enterprise Shield - Distribution Guide

Complete guide for distributing and installing Enterprise Shield across all platforms.

---

## 📦 Distribution Methods

### 1. **One-Line Install Script** (Recommended for End Users)

**For users:**
```bash
curl -sSL https://raw.githubusercontent.com/yourorg/opencode-enterprise-shield/main/install.sh | bash
```

**Features:**
- ✅ Auto-detects OS and architecture
- ✅ Downloads correct binary
- ✅ Verifies checksums
- ✅ Installs to `~/.opencode/plugins/`
- ✅ Creates default configuration
- ✅ Adds to PATH
- ✅ Runs health check

**Uninstall:**
```bash
curl -sSL https://raw.githubusercontent.com/yourorg/opencode-enterprise-shield/main/uninstall.sh | bash
```

---

### 2. **Homebrew** (macOS/Linux)

**Setup Homebrew Tap (One-time):**
```bash
# 1. Create tap repository
mkdir -p homebrew-opencode-enterprise-shield
cd homebrew-opencode-enterprise-shield

# 2. Copy formula
cp ../opencode-enterprise-shield/Formula/enterprise-shield.rb Formula/

# 3. Create README
cat > README.md << 'EOF'
# Homebrew Tap for Enterprise Shield

## Installation

```bash
brew tap yourorg/opencode-enterprise-shield
brew install enterprise-shield
```

## Updating

```bash
brew update
brew upgrade enterprise-shield
```
EOF

# 4. Push to GitHub
git init
git add .
git commit -m "Initial formula"
git remote add origin https://github.com/yourorg/homebrew-opencode-enterprise-shield
git push -u origin main
```

**For users:**
```bash
# Install
brew tap yourorg/opencode-enterprise-shield
brew install enterprise-shield

# Update
brew upgrade enterprise-shield

# Uninstall
brew uninstall enterprise-shield
brew untap yourorg/opencode-enterprise-shield
```

**Updating Formula After Release:**
```bash
# After creating a new release, update the formula:
# 1. Update version in Formula/enterprise-shield.rb
# 2. Update SHA256 checksums (from release artifacts)
# 3. Commit and push

# Formula auto-update script
./scripts/update-homebrew-formula.sh v1.0.1
```

---

### 3. **Go Install** (Universal)

**For users:**
```bash
# Install latest
go install github.com/yourorg/opencode-enterprise-shield/cmd/plugin@latest

# Install specific version
go install github.com/yourorg/opencode-enterprise-shield/cmd/plugin@v1.0.0

# The binary will be in $GOPATH/bin or $HOME/go/bin
# Symlink to OpenCode plugins directory
ln -s $(go env GOPATH)/bin/plugin ~/.opencode/plugins/enterprise-shield
```

**Setup (Maintainer):**
- No special setup needed
- Just tag releases with `vX.Y.Z`
- Users can install directly from GitHub

---

### 4. **Pre-built Binaries** (GitHub Releases)

**Download and Install Manually:**
```bash
# 1. Download for your platform
VERSION=v1.0.0
OS=darwin  # or linux, windows
ARCH=arm64 # or amd64

wget https://github.com/yourorg/opencode-enterprise-shield/releases/download/${VERSION}/enterprise-shield-${VERSION}-${OS}-${ARCH}.tar.gz

# 2. Verify checksum
wget https://github.com/yourorg/opencode-enterprise-shield/releases/download/${VERSION}/enterprise-shield-${VERSION}-${OS}-${ARCH}.tar.gz.sha256
sha256sum -c enterprise-shield-${VERSION}-${OS}-${ARCH}.tar.gz.sha256

# 3. Extract
tar -xzf enterprise-shield-${VERSION}-${OS}-${ARCH}.tar.gz

# 4. Install
mkdir -p ~/.opencode/plugins
mv enterprise-shield-${VERSION}-${OS}-${ARCH} ~/.opencode/plugins/enterprise-shield
chmod +x ~/.opencode/plugins/enterprise-shield
```

---

### 5. **Docker Image** (Optional)

**Build Image:**
```dockerfile
# Dockerfile
FROM golang:1.22-alpine AS builder
WORKDIR /build
COPY . .
RUN go build -ldflags="-w -s" -o enterprise-shield ./cmd/plugin

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /build/enterprise-shield /usr/local/bin/
COPY config/default.yaml /etc/enterprise-shield/config.yaml
ENTRYPOINT ["enterprise-shield"]
CMD ["serve"]
```

**For users:**
```bash
docker pull yourorg/enterprise-shield:latest
docker run -v ~/.opencode:/root/.opencode yourorg/enterprise-shield:latest
```

---

### 6. **NPM Package** (OpenCode Native Integration)

**Package.json:**
```json
{
  "name": "@yourorg/opencode-enterprise-shield",
  "version": "1.0.0",
  "description": "Enterprise Shield plugin for OpenCode",
  "main": ".opencode/plugin.js",
  "bin": {
    "enterprise-shield": "bin/enterprise-shield"
  },
  "scripts": {
    "postinstall": "node scripts/download-binary.js"
  }
}
```

**For users:**
```bash
npm install -g @yourorg/opencode-enterprise-shield
```

---

## 🚀 Release Process

### Creating a New Release

**1. Update Version:**
```bash
# Update version in:
# - cmd/plugin/main.go (version const)
# - Formula/enterprise-shield.rb
# - .opencode/plugin.yaml
# - package.json (if using NPM)
# - CHANGELOG.md

NEW_VERSION="v1.0.1"
```

**2. Update CHANGELOG:**
```bash
# Add new section to CHANGELOG.md
cat >> CHANGELOG.md << EOF

## [1.0.1] - $(date +%Y-%m-%d)

### Added
- New feature X

### Fixed
- Bug Y

### Changed
- Improvement Z
EOF
```

**3. Commit and Tag:**
```bash
git add .
git commit -m "Release ${NEW_VERSION}"
git tag -a "${NEW_VERSION}" -m "Release ${NEW_VERSION}"
git push origin main --tags
```

**4. GitHub Actions Auto-Build:**
- Workflow triggers on tag push
- Builds binaries for all platforms
- Generates checksums
- Creates GitHub release
- Uploads artifacts

**5. Update Homebrew Formula:**
```bash
# After release is created, update checksums in formula
./scripts/update-homebrew-formula.sh ${NEW_VERSION}
```

---

## 📋 Platform Support

| Platform | Architecture | Status | Install Method |
|----------|-------------|--------|----------------|
| macOS    | Intel (amd64) | ✅ Supported | Homebrew, Install Script, Go Install |
| macOS    | Apple Silicon (arm64) | ✅ Supported | Homebrew, Install Script, Go Install |
| Linux    | amd64 | ✅ Supported | Install Script, Go Install, Binary |
| Linux    | arm64 | ✅ Supported | Install Script, Go Install, Binary |
| Windows  | amd64 | ✅ Supported | Binary, Go Install |
| Windows  | arm64 | 🔄 Planned | - |

---

## 🔧 Troubleshooting

### Install Script Fails

**Problem:** `curl: command not found`

**Solution:**
```bash
# macOS
brew install curl

# Linux (Ubuntu/Debian)
sudo apt-get install curl

# Linux (RHEL/CentOS)
sudo yum install curl
```

**Problem:** `Permission denied`

**Solution:**
```bash
chmod +x ~/.opencode/plugins/enterprise-shield
```

---

### Binary Won't Run

**Problem:** `cannot execute binary file`

**Solution:** Downloaded wrong architecture
```bash
# Check your architecture
uname -m

# Re-download correct binary
```

**Problem (macOS):** `"enterprise-shield" cannot be opened because the developer cannot be verified`

**Solution:**
```bash
xattr -d com.apple.quarantine ~/.opencode/plugins/enterprise-shield
```

---

### OpenCode Won't Load Plugin

**Problem:** Plugin not detected

**Solution:**
```bash
# Ensure plugin.yaml exists
ls ~/.opencode/plugins/enterprise-shield
ls .opencode/plugin.yaml

# Check permissions
chmod +x ~/.opencode/plugins/enterprise-shield

# Check configuration
cat ~/.opencode/config/enterprise-shield.yaml
```

---

## 📊 Distribution Checklist

Before releasing a new version:

- [ ] Tests pass (`make test`)
- [ ] Version updated in all files
- [ ] CHANGELOG.md updated
- [ ] Git tagged with version
- [ ] GitHub release created
- [ ] All platform binaries built
- [ ] Checksums verified
- [ ] Homebrew formula updated
- [ ] Documentation updated
- [ ] Release notes written
- [ ] Announced in changelog

---

## 🌐 Publishing Locations

| Location | URL | Purpose |
|----------|-----|---------|
| GitHub Releases | `github.com/yourorg/opencode-enterprise-shield/releases` | Primary distribution |
| Homebrew Tap | `github.com/yourorg/homebrew-opencode-enterprise-shield` | macOS/Linux package manager |
| NPM Registry | `npmjs.com/package/@yourorg/opencode-enterprise-shield` | Node/OpenCode native |
| Docker Hub | `hub.docker.com/r/yourorg/enterprise-shield` | Container distribution |
| Go Packages | `pkg.go.dev/github.com/yourorg/opencode-enterprise-shield` | Go module repository |

---

## 📞 Support

For distribution issues:
- GitHub Issues: https://github.com/yourorg/opencode-enterprise-shield/issues
- Documentation: https://github.com/yourorg/opencode-enterprise-shield
- Email: support@yourorg.com

