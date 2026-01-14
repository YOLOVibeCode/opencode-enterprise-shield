# 🤖 Automated Deployment - How It Works

## ✅ What Happens When You Push to Main

Every time you push code to the `main` branch, this happens **automatically**:

---

## 🔄 The Automated Pipeline

```
┌─────────────────────────────────────────────────────────────────┐
│  Developer pushes to main branch                                 │
│  $ git push origin main                                          │
└─────────────────────┬───────────────────────────────────────────┘
                      │
        ┌─────────────┴─────────────┐
        │                           │
        ▼                           ▼
┌───────────────────┐       ┌───────────────────┐
│   Run Tests       │       │   Build Binaries  │
│   (Parallel)      │       │   (After tests)   │
├───────────────────┤       ├───────────────────┤
│ • Ubuntu (Go 1.22)│       │ • linux-amd64     │
│ • macOS (Go 1.22) │       │ • linux-arm64     │
│ • 24 unit tests   │       │ • darwin-amd64    │
│ • Linting         │       │ • darwin-arm64    │
│ • Coverage report │       │ • windows-amd64   │
└─────────┬─────────┘       └─────────┬─────────┘
          │                           │
          │ ✅ Tests pass             │ ✅ Binaries built
          │                           │
          └─────────────┬─────────────┘
                        │
                        ▼
        ┌───────────────────────────────────┐
        │   Generate Checksums              │
        │   • SHA256 for each binary        │
        └───────────────┬───────────────────┘
                        │
                        ▼
        ┌───────────────────────────────────┐
        │   Update 'dev' Release            │
        ├───────────────────────────────────┤
        │ 1. Delete old 'dev' release       │
        │ 2. Create new 'dev' release       │
        │ 3. Upload all 5 binaries          │
        │ 4. Upload checksums               │
        │ 5. Generate release notes         │
        └───────────────┬───────────────────┘
                        │
                        ▼
        ┌───────────────────────────────────┐
        │   Build Docker Images             │
        ├───────────────────────────────────┤
        │ • Multi-arch (amd64, arm64)       │
        │ • Push to GitHub Container Registry│
        │ • Tags: dev, latest, main-<sha>   │
        └───────────────┬───────────────────┘
                        │
                        ▼
        ┌───────────────────────────────────┐
        │   ✅ Deployment Complete!          │
        ├───────────────────────────────────┤
        │ All channels updated:             │
        │ • GitHub releases/tag/dev         │
        │ • Docker ghcr.io/...:dev          │
        │ • Docker ghcr.io/...:latest       │
        │ • Install script (VERSION=dev)    │
        └───────────────────────────────────┘
                        │
                        ▼
        ⏱️  Total time: ~5 minutes
```

---

## 📦 What Gets Updated Automatically

### 1. GitHub Development Release

**Location:** `https://github.com/yourorg/opencode-enterprise-shield/releases/tag/dev`

**Contains:**
- ✅ Binaries for all 5 platforms
- ✅ SHA256 checksums
- ✅ Auto-generated release notes
- ✅ Recent commit history
- ✅ Installation instructions

**Users access with:**
```bash
curl -sSL https://raw.githubusercontent.com/yourorg/enterprise-shield/main/install.sh | VERSION=dev bash
```

### 2. Docker Images (GitHub Container Registry)

**Images created/updated:**
```bash
ghcr.io/yourorg/opencode-enterprise-shield:dev        # Latest main branch
ghcr.io/yourorg/opencode-enterprise-shield:latest    # Latest main branch
ghcr.io/yourorg/opencode-enterprise-shield:main-abc123 # Specific commit
```

**Multi-architecture support:**
- ✅ linux/amd64
- ✅ linux/arm64

**Users access with:**
```bash
docker pull ghcr.io/yourorg/opencode-enterprise-shield:dev
```

### 3. Install Script

**Updated behavior:**
```bash
# Latest development (auto-updated)
VERSION=dev bash install.sh

# Latest stable (from releases)
VERSION=latest bash install.sh

# Specific version
VERSION=v1.0.0 bash install.sh
```

### 4. Go Module (pkg.go.dev)

**Auto-indexed** when you push to main:
```bash
go install github.com/yourorg/opencode-enterprise-shield/cmd/plugin@main
```

---

## 🎯 User Experience

### For End Users

**Latest Stable (Recommended):**
```bash
curl -sSL https://raw.githubusercontent.com/yourorg/enterprise-shield/main/install.sh | bash
# Installs latest tagged version (e.g., v1.0.0)
```

**Latest Development (Bleeding Edge):**
```bash
curl -sSL https://raw.githubusercontent.com/yourorg/enterprise-shield/main/install.sh | VERSION=dev bash
# Installs latest main branch build
```

### For Beta Testers

```bash
# Get development builds automatically
export ENTERPRISE_SHIELD_CHANNEL=dev
curl -sSL ... | bash

# Or via Docker
docker pull ghcr.io/yourorg/opencode-enterprise-shield:dev
```

---

## 🔄 Complete CI/CD Matrix

### Trigger: Pull Request

| Workflow | Tests | Build | Deploy | Status |
|----------|-------|-------|--------|--------|
| test.yml | ✅ Yes | ✅ Yes | ❌ No | Validate only |
| main.yml | ✅ Yes | ❌ No | ❌ No | Validate only |
| docker-publish.yml | ❌ No | ✅ Yes | ❌ No | Build test only |

**Result:** Changes validated but NOT deployed

### Trigger: Push to Main

| Workflow | Tests | Build | Deploy | Artifacts |
|----------|-------|-------|--------|-----------|
| main.yml | ✅ Yes | ✅ Yes | ✅ Yes | 5 binaries → `dev` release |
| docker-publish.yml | ❌ No | ✅ Yes | ✅ Yes | Docker → `:dev`, `:latest` |

**Result:** All channels updated with development builds (~5 min)

### Trigger: Git Tag (v*)

| Workflow | Tests | Build | Deploy | Artifacts |
|----------|-------|-------|--------|-----------|
| release.yml | ✅ Yes | ✅ Yes | ✅ Yes | 5 binaries → version release |
| docker-publish.yml | ❌ No | ✅ Yes | ✅ Yes | Docker → `:v1.0.0`, `:1.0`, `:1`, `:latest` |

**Result:** Stable release created on all channels (~5 min)

---

## 📊 Distribution Channel Status

| Channel | Development (main) | Stable (tags) | Auto-Update |
|---------|-------------------|---------------|-------------|
| **GitHub Releases** | `/releases/tag/dev` | `/releases/tag/v1.0.0` | ✅ Auto |
| **Install Script** | `VERSION=dev` | `VERSION=latest` | ✅ Auto |
| **Docker (GHCR)** | `:dev` tag | `:v1.0.0`, `:latest` | ✅ Auto |
| **Go Install** | `@main` | `@v1.0.0` | ✅ Auto |
| **Homebrew** | Manual | Manual | ⚠️ Script |
| **NPM** | Manual | Manual | ⚠️ Manual |

**Legend:**
- ✅ Auto = Fully automated via GitHub Actions
- ⚠️ Script = Semi-automated (run script after release)
- ⚠️ Manual = Manual publish required

---

## 🎬 Example: What Happens When You Push

```bash
# You make a change
$ vim pkg/sanitizer/engine.go
$ git add .
$ git commit -m "feat: Add new detection pattern"
$ git push origin main

# Within seconds...
```

**GitHub Actions logs:**
```
✓ Checkout code
✓ Set up Go 1.22
✓ Download dependencies
✓ Run tests (24 tests)
  ✓ TestSanitize_ServerNames
  ✓ TestSanitize_IPAddresses
  ✓ TestDesanitize_Basic
  ... (all 24 tests pass)
✓ Run linter
✓ Build linux-amd64 binary
✓ Build linux-arm64 binary
✓ Build darwin-amd64 binary
✓ Build darwin-arm64 binary
✓ Build windows-amd64 binary
✓ Generate checksums
✓ Delete old 'dev' release
✓ Create new 'dev' release
✓ Upload 5 binaries + checksums
✓ Build Docker image (multi-arch)
✓ Push to ghcr.io (tags: dev, latest, main-abc123)

🎉 Deployment complete! (4m 32s)
```

**Users can now install:**
```bash
# Immediately available
curl -sSL https://raw.githubusercontent.com/yourorg/enterprise-shield/main/install.sh | VERSION=dev bash

# Or via Docker
docker pull ghcr.io/yourorg/opencode-enterprise-shield:dev
```

---

## 🚀 Quick Reference

### To Deploy Development Build

```bash
# Just push to main - automatic!
git push origin main
```

### To Deploy Stable Release

```bash
# Create and push tag
git tag -a v1.0.1 -m "Release v1.0.1"
git push origin v1.0.1
```

### To Check Deployment Status

```bash
# View workflows
open https://github.com/yourorg/opencode-enterprise-shield/actions

# View releases
open https://github.com/yourorg/opencode-enterprise-shield/releases

# View Docker images
open https://github.com/yourorg/opencode-enterprise-shield/pkgs/container/opencode-enterprise-shield
```

---

## ✅ Summary

**Every push to `main` automatically:**
1. ✅ Runs 24 tests on 2 platforms
2. ✅ Builds binaries for 5 platforms
3. ✅ Creates `dev` GitHub release
4. ✅ Uploads binaries with checksums
5. ✅ Builds and pushes Docker images
6. ✅ Makes available via install script

**Every Git tag automatically:**
1. ✅ All of the above, plus:
2. ✅ Creates stable versioned release
3. ✅ Tags Docker with version numbers
4. ✅ Adds to Go package index
5. ✅ Generates release notes from CHANGELOG

**Zero manual intervention needed for deployment!** 🎉

---

*See [docs/CONTINUOUS_DEPLOYMENT.md](docs/CONTINUOUS_DEPLOYMENT.md) for complete details*

