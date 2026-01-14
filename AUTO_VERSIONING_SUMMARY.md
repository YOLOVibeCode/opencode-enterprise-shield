# ✅ Auto-Versioning System - Complete

## 🎉 Intelligent Build Numbering Implemented

Every push to main now **automatically**:
- ✅ Analyzes commit message
- ✅ Bumps version based on semantic keywords
- ✅ Increments build number
- ✅ Updates VERSION file
- ✅ Builds binaries with full version
- ✅ Deploys to all channels

---

## 🔢 How It Works

### Version Format

**Semantic Version:** `MAJOR.MINOR.PATCH` (stored in `VERSION` file)

**Full Version:** `MAJOR.MINOR.PATCH+build.NUMBER`

**Example:** `1.0.0+build.42`

### Build Number

- ✅ **Auto-increments** on every push (uses GitHub run number)
- ✅ **Never decreases** (monotonically increasing)
- ✅ **Unique per build** (no collisions)
- ✅ **Embedded in binary** (visible via `--version`)

---

## 🤖 Automatic Version Bumping

### Commit Message Triggers

| Commit Prefix | Version Change | Example |
|--------------|----------------|---------|
| `major:` or `BREAKING CHANGE:` | `1.0.0` → `2.0.0` | Incompatible changes |
| `feat:` or `feature:` or `minor:` | `1.0.0` → `1.1.0` | New features |
| `fix:` or `patch:` | `1.0.0` → `1.0.1` | Bug fixes |
| `docs:`, `chore:`, `ci:` | No change | Build number increments |

**Build number ALWAYS increments** regardless of version change.

### Example Flow

```bash
# Current state: VERSION = 1.0.0, Build = 5

# Developer commits with "feat:" prefix
$ git commit -m "feat: Add GPU detection pattern"
$ git push origin main

# GitHub Actions automatically:
# 1. Reads VERSION file (1.0.0)
# 2. Detects "feat:" in commit message
# 3. Bumps MINOR: 1.0.0 → 1.1.0
# 4. Increments build: 5 → 6
# 5. Full version: 1.1.0+build.6
# 6. Updates VERSION file to 1.1.0
# 7. Commits back: "chore: Bump version to 1.1.0 [skip ci]"
# 8. Builds all binaries with v1.1.0+build.6
# 9. Creates dev release: "Development Build 1.1.0+build.6 (Build #6)"
# 10. All binaries available for download
```

---

## 📊 Real Examples

### Scenario 1: Documentation Update

```bash
$ git commit -m "docs: Update installation guide"
$ git push origin main

GitHub Actions:
├─ Current: 1.0.0 (from VERSION file)
├─ Commit: "docs:" (no version bump)
├─ Build: 6 → 7
├─ Result: 1.0.0+build.7
├─ VERSION file: UNCHANGED (stays 1.0.0)
└─ Binary: enterprise-shield v1.0.0+build.7
```

### Scenario 2: Bug Fix

```bash
$ git commit -m "fix: Correct regex timeout handling"
$ git push origin main

GitHub Actions:
├─ Current: 1.0.0
├─ Commit: "fix:" (PATCH bump)
├─ Version: 1.0.0 → 1.0.1
├─ Build: 7 → 8
├─ Result: 1.0.1+build.8
├─ VERSION file: UPDATED to 1.0.1
├─ Auto-commit: "chore: Bump version to 1.0.1 [skip ci]"
└─ Binary: enterprise-shield v1.0.1+build.8
```

### Scenario 3: New Feature

```bash
$ git commit -m "feat: Add Azure Key Vault support"
$ git push origin main

GitHub Actions:
├─ Current: 1.0.1
├─ Commit: "feat:" (MINOR bump)
├─ Version: 1.0.1 → 1.1.0
├─ Build: 8 → 9
├─ Result: 1.1.0+build.9
├─ VERSION file: UPDATED to 1.1.0
├─ Auto-commit: "chore: Bump version to 1.1.0 [skip ci]"
└─ Binary: enterprise-shield v1.1.0+build.9
```

### Scenario 4: Breaking Change

```bash
$ git commit -m "major: Redesign configuration format"
$ git push origin main

GitHub Actions:
├─ Current: 1.1.0
├─ Commit: "major:" (MAJOR bump)
├─ Version: 1.1.0 → 2.0.0
├─ Build: 9 → 10
├─ Result: 2.0.0+build.10
├─ VERSION file: UPDATED to 2.0.0
├─ Auto-commit: "chore: Bump version to 2.0.0 [skip ci]"
└─ Binary: enterprise-shield v2.0.0+build.10
```

---

## 📦 Binaries Get Proper Versions

### Binary Naming

Development builds:
```
enterprise-shield-1.0.0+build.42-darwin-arm64.tar.gz
enterprise-shield-1.1.0+build.43-linux-amd64.tar.gz
enterprise-shield-2.0.0+build.50-windows-amd64.zip
```

Stable releases (git tags):
```
enterprise-shield-v1.0.0-darwin-arm64.tar.gz
enterprise-shield-v1.0.1-linux-amd64.tar.gz
```

### Version in Binary

```bash
$ ./enterprise-shield version
Enterprise Shield Plugin v1.0.0+build.42
Build: #42
Built: 2026-01-14_03:15:30
```

### Version in GitHub Release

**Title:** "Development Build 1.0.0+build.42 (Build #42)"

**Description includes:**
- Full version with build number
- Commit SHA
- Build date/time
- Link to CI run
- Recent changes (last 10 commits)
- Checksums for all binaries

---

## 🔧 Manual Version Control

### Bump Version Manually

```bash
# Patch bump: 1.0.0 → 1.0.1
make bump-patch

# Minor bump: 1.0.0 → 1.1.0
make bump-minor

# Major bump: 1.0.0 → 2.0.0
make bump-major

# Check current version
make version
# Output:
# Version: 1.0.0
# Build: 6
# Full: 1.0.0+build.6
```

### Create Stable Release

```bash
# Set version and create tag
make tag V=1.2.0
git push origin main v1.2.0

# Or manually
echo "1.2.0" > VERSION
git add VERSION
git commit -m "chore: Release v1.2.0"
git tag -a v1.2.0 -m "Release v1.2.0"
git push origin main v1.2.0
```

---

## 📈 Build Number Properties

### Always Increments

```
Commit 1 → Build 1 (1.0.0+build.1)
Commit 2 → Build 2 (1.0.0+build.2)
Commit 3 → Build 3 (1.0.1+build.3)  ← Version bumped, build increments
Commit 4 → Build 4 (1.0.1+build.4)
...
Commit 42 → Build 42 (1.2.0+build.42)
```

### Source

| Environment | Build Number From |
|-------------|-------------------|
| **GitHub Actions** | `${{ github.run_number }}` (workflow run count) |
| **Local Build** | Git commit count: `git rev-list --count HEAD` |
| **Manual Override** | `BUILD_NUMBER=123 make build` |

### Visibility

Build number appears in:
- ✅ Binary version output
- ✅ GitHub release title
- ✅ Docker image labels
- ✅ Artifact filenames
- ✅ Release notes

---

## 🎯 Complete Workflow Example

### Developer Workflow

```bash
# 1. Make changes
$ vim pkg/sanitizer/engine.go

# 2. Commit with semantic prefix
$ git commit -m "feat: Add new detection pattern for Azure resources"

# 3. Push to main
$ git push origin main
```

### GitHub Actions (Automatic)

```
[1] Version Calculation (30 seconds)
    ├─ Read VERSION: 1.0.0
    ├─ Analyze commit: "feat:" found
    ├─ Bump MINOR: 1.0.0 → 1.1.0
    ├─ Build number: 42
    ├─ Full version: 1.1.0+build.42
    └─ Update VERSION file → 1.1.0

[2] Run Tests (1 minute)
    ├─ Ubuntu: 24 tests ✅
    ├─ macOS: 24 tests ✅
    └─ Coverage: 95%+ ✅

[3] Build Binaries (2 minutes)
    ├─ linux-amd64: enterprise-shield-1.1.0+build.42-linux-amd64.tar.gz
    ├─ linux-arm64: enterprise-shield-1.1.0+build.42-linux-arm64.tar.gz
    ├─ darwin-amd64: enterprise-shield-1.1.0+build.42-darwin-amd64.tar.gz
    ├─ darwin-arm64: enterprise-shield-1.1.0+build.42-darwin-arm64.tar.gz
    └─ windows-amd64: enterprise-shield-1.1.0+build.42-windows-amd64.zip

[4] Generate Checksums (10 seconds)
    └─ SHA256 for each binary

[5] Create Dev Release (30 seconds)
    ├─ Delete old 'dev' release
    ├─ Create new 'dev' release
    ├─ Title: "Development Build 1.1.0+build.42 (Build #42)"
    ├─ Upload 5 binaries
    ├─ Upload checksums
    └─ Generate notes with recent commits

[6] Build Docker (1 minute)
    ├─ Multi-arch: linux/amd64, linux/arm64
    ├─ Push tags: dev, latest, build-42, 1.1.0
    └─ Labels: version, build number, commit SHA

Total time: ~5 minutes
```

### Users Get Updated Binary

```bash
# Immediately available (within 5 minutes of push)
$ curl -sSL https://raw.githubusercontent.com/YOLOVibeCode/enterprise-shield/main/install.sh | VERSION=dev bash

Downloading: enterprise-shield-1.1.0+build.42-darwin-arm64.tar.gz
✓ Downloaded
✓ Verified checksum
✓ Installed to ~/.opencode/plugins/enterprise-shield

$ enterprise-shield version
Enterprise Shield Plugin v1.1.0+build.42
Build: #42
Built: 2026-01-14_03:15:30
```

---

## 📋 What Gets Versioned

### GitHub Releases

**Development (dev tag):**
- Title: `Development Build 1.0.0+build.42 (Build #42)`
- Binaries: `enterprise-shield-1.0.0+build.42-platform-arch.tar.gz`
- Auto-created on every main push

**Stable (vX.Y.Z tags):**
- Title: `Enterprise Shield v1.0.0`
- Binaries: `enterprise-shield-v1.0.0-platform-arch.tar.gz`
- Created when you push git tag

### Docker Images

**Tags created for build 42 with version 1.0.0:**
```
ghcr.io/YOLOVibeCode/opencode-enterprise-shield:dev          # Latest dev
ghcr.io/YOLOVibeCode/opencode-enterprise-shield:latest      # Latest main
ghcr.io/YOLOVibeCode/opencode-enterprise-shield:build-42    # Specific build
ghcr.io/YOLOVibeCode/opencode-enterprise-shield:1.0.0       # Semantic version
```

### Binaries

**Embedded version info:**
```bash
$ enterprise-shield version
Enterprise Shield Plugin v1.0.0+build.42
Build: #42
Built: 2026-01-14_03:15:30
```

---

## 🚀 The Magic

### Before (Manual)

```
Developer pushes to main
    ↓
Manually decide version
    ↓
Manually update version in code
    ↓
Manually build binaries
    ↓
Manually upload to releases
    ↓
Takes hours, error-prone
```

### After (Automated) ✅

```
Developer pushes to main with "feat: New feature"
    ↓
GitHub Actions automatically:
├─ Detects "feat:" → Bump MINOR
├─ Increments build number
├─ Updates VERSION file
├─ Builds 5 binaries with full version
├─ Creates dev release
├─ Uploads all binaries
├─ Builds Docker with version tags
└─ Everything ready in 5 minutes

Users can immediately download versioned binaries!
```

---

## 📊 Version Tracking

### View Version History

```bash
# Current version
$ cat VERSION
1.0.0

# Version with build number (local)
$ make version
Version: 1.0.0
Build: 6
Full: 1.0.0+build.6

# Version in binary
$ ./build/enterprise-shield version
Enterprise Shield Plugin v1.0.0+build.6
Build: #6
Built: 2026-01-14_03:21:53

# All git tags (stable releases)
$ git tag -l "v*"
v1.0.0
v1.0.1
v1.1.0

# Development builds
# Check: https://github.com/YOLOVibeCode/enterprise-shield/releases/tag/dev
```

### Version in Release Assets

Every release includes:
- Binary name: `enterprise-shield-1.0.0+build.42-darwin-arm64.tar.gz`
- Checksum file: `.sha256`
- Release notes with build info
- Link to exact commit

---

## 🎯 User Experience

### Installing Latest Dev Build

```bash
$ curl -sSL https://raw.githubusercontent.com/YOLOVibeCode/enterprise-shield/main/install.sh | VERSION=dev bash

Downloading: enterprise-shield-1.1.0+build.42-darwin-arm64.tar.gz
✓ Downloaded
✓ Verified checksum

$ enterprise-shield version
Enterprise Shield Plugin v1.1.0+build.42
Build: #42
Built: 2026-01-14_03:15:30

# User knows EXACT build they have
# Can reference "Build #42" in bug reports
```

### Installing Stable Release

```bash
$ curl -sSL https://raw.githubusercontent.com/YOLOVibeCode/enterprise-shield/main/install.sh | bash

Downloading: enterprise-shield-v1.0.0-darwin-arm64.tar.gz

$ enterprise-shield version
Enterprise Shield Plugin v1.0.0
# No build number for stable releases
```

---

## 📋 Complete Versioning Features

### ✅ Implemented

- [x] **VERSION file** - Single source of truth
- [x] **Auto-increment** - Build numbers always increase
- [x] **Semantic bumping** - Based on commit messages
- [x] **Auto-commit** - VERSION file updated automatically
- [x] **Binary embedding** - Version in compiled binary
- [x] **Manual bumps** - `make bump-[major|minor|patch]`
- [x] **Docker tags** - Version, build number, latest
- [x] **Release naming** - Descriptive with build number
- [x] **Checksum files** - SHA256 for every binary
- [x] **Full traceability** - Build → Commit → CI run

### 🎁 Bonus Features

- ✅ **Commit validation** - Enforces semantic prefixes (optional)
- ✅ **Changelog integration** - Auto-generates from commits
- ✅ **Build summary** - GitHub Actions summary shows version
- ✅ **Docker labels** - Metadata in images
- ✅ **Skip CI** - Auto-commits have `[skip ci]` to prevent loops

---

## 🔍 Debugging & Troubleshooting

### Check What Version Will Be

```bash
# View current VERSION file
cat VERSION

# See what next build will be
make version

# Simulate version bump
./scripts/bump-version.sh minor  # Shows what would happen
git checkout VERSION              # Undo changes
```

### Force Specific Version

```bash
# Set version manually
echo "2.5.0" > VERSION
git add VERSION
git commit -m "chore: Set version to 2.5.0"
git push origin main

# Next build will be 2.5.0+build.XX
```

### View All Builds

```bash
# GitHub releases
open https://github.com/YOLOVibeCode/enterprise-shield/releases

# Dev builds
open https://github.com/YOLOVibeCode/enterprise-shield/releases/tag/dev

# Docker builds
docker pull ghcr.io/YOLOVibeCode/opencode-enterprise-shield:build-42
```

---

## 📊 Comparison

### Before Auto-Versioning

```
Build process:
1. Manually decide if version should bump
2. Edit version in multiple files
3. Forget to update CHANGELOG
4. Build binaries with vague version "dev"
5. Upload to releases
6. Users get "dev" build (can't identify which one)

Problems:
❌ No traceability
❌ Manual errors
❌ Can't identify specific builds
❌ Time-consuming
```

### After Auto-Versioning ✅

```
Build process:
1. Push to main with semantic commit message
2. Everything else is automatic

Benefits:
✅ Every build has unique version
✅ Full traceability (build → commit → CI run)
✅ Users can report "Build #42" in bug reports
✅ Zero manual intervention
✅ Consistent versioning across all artifacts
✅ Takes 5 minutes, fully automated
```

---

## 🎉 Summary

**Every push to main branch now:**

1. ✅ **Intelligently analyzes** commit message
2. ✅ **Auto-bumps version** if `feat:`, `fix:`, or `major:`
3. ✅ **Always increments** build number
4. ✅ **Updates VERSION** file in repository
5. ✅ **Builds all binaries** with full version string
6. ✅ **Creates dev release** with descriptive title
7. ✅ **Tags Docker images** with version and build number
8. ✅ **Makes available** via all distribution channels

**Users get:**
- Unique version for every build
- Full traceability
- Easy bug reporting ("I'm on Build #42")
- Automatic updates via install script

**You get:**
- Zero manual version management
- Semantic versioning enforced
- Complete build history
- Professional release process

---

## 📝 Quick Reference

```bash
# Check version locally
make version

# Bump version manually
make bump-patch   # 1.0.0 → 1.0.1
make bump-minor   # 1.0.0 → 1.1.0
make bump-major   # 1.0.0 → 2.0.0

# Auto-bump based on commit
git commit -m "fix: Bug fix"    # Auto-bumps patch
git commit -m "feat: New thing" # Auto-bumps minor
git commit -m "major: Breaking" # Auto-bumps major
git push origin main            # Automation kicks in

# Create stable release
make tag V=1.2.0
git push origin main v1.2.0
```

---

**Perfect versioning: Automatic, semantic, traceable** 🎯

