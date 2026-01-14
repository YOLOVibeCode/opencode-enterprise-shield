# Continuous Deployment Guide

## Overview

Enterprise Shield has automated CI/CD pipelines that ensure every change is tested and deployed across all distribution channels.

---

## 🔄 Automated Workflows

### 1. **Main Branch CI/CD** (`.github/workflows/main.yml`)

**Triggers:** Push to `main` or Pull Request

**What it does:**

#### On Pull Request:
1. ✅ **Tests** on Ubuntu and macOS
2. ✅ **Linting** with golangci-lint
3. ✅ **Coverage** reporting to Codecov
4. ⏸️ Does NOT deploy (PR only validates)

#### On Push to Main:
1. ✅ **Tests** on Ubuntu and macOS
2. ✅ **Build** binaries for all 5 platforms
3. ✅ **Create/Update** `dev` release tag
4. ✅ **Upload** development binaries
5. ✅ **Make available** via install script

**Development Build Access:**
```bash
# Install latest development build
curl -sSL https://raw.githubusercontent.com/YOLOVibeCode/enterprise-shield/main/install.sh | VERSION=dev bash

# Or download manually from
https://github.com/YOLOVibeCode/enterprise-shield/releases/tag/dev
```

---

### 2. **Release Workflow** (`.github/workflows/release.yml`)

**Triggers:** Git tag matching `v*` (e.g., `v1.0.0`, `v1.0.1`)

**What it does:**
1. ✅ **Build** binaries for all 5 platforms
2. ✅ **Generate** SHA256 checksums
3. ✅ **Create** GitHub release with:
   - Release notes from CHANGELOG.md
   - Installation instructions
   - All platform binaries
   - Checksums file
4. ✅ **Run** integration tests
5. ✅ **Update** Homebrew formula (if configured)

**Stable Release Access:**
```bash
# Install latest stable release
curl -sSL https://raw.githubusercontent.com/YOLOVibeCode/enterprise-shield/main/install.sh | bash

# Or specific version
curl -sSL https://raw.githubusercontent.com/YOLOVibeCode/enterprise-shield/main/install.sh | VERSION=v1.0.0 bash
```

---

### 3. **Docker Publishing** (`.github/workflows/docker-publish.yml`)

**Triggers:** Push to `main` OR tag `v*`

**What it does:**
1. ✅ **Build** multi-arch Docker images (amd64, arm64)
2. ✅ **Push** to GitHub Container Registry
3. ✅ **Tag** appropriately:
   - `latest` - Latest stable release
   - `dev` - Latest main branch
   - `v1.0.0` - Specific version tags
   - `1.0` - Major.minor tags
   - `1` - Major version tags

**Docker Image Access:**
```bash
# Latest stable
docker pull ghcr.io/YOLOVibeCode/opencode-enterprise-shield:latest

# Development
docker pull ghcr.io/YOLOVibeCode/opencode-enterprise-shield:dev

# Specific version
docker pull ghcr.io/YOLOVibeCode/opencode-enterprise-shield:v1.0.0
```

---

### 4. **Test Workflow** (`.github/workflows/test.yml`)

**Triggers:** Push to any branch, Pull Request

**What it does:**
1. ✅ **Tests** on Ubuntu and macOS
2. ✅ **Linting** with golangci-lint
3. ✅ **Coverage** reporting
4. ✅ **Build** verification (ensures it compiles)

---

## 📦 Distribution Channels

### Automatic Updates on Main Push

| Channel | Update Trigger | Access |
|---------|---------------|--------|
| **GitHub Releases** | Every main push | `dev` tag (prerelease) |
| **Install Script** | Every main push | `VERSION=dev` flag |
| **Docker (dev)** | Every main push | `:dev` tag |
| **Docker (latest)** | Every main push | `:latest` tag (if on main) |

### Automatic Updates on Version Tag

| Channel | Update Trigger | Access |
|---------|---------------|--------|
| **GitHub Releases** | Git tag `v*` | Stable release |
| **Install Script** | Git tag `v*` | Default (latest) |
| **Docker (versioned)** | Git tag `v*` | `:v1.0.0`, `:1.0`, `:1` tags |
| **Docker (latest)** | Git tag `v*` | `:latest` tag |
| **Go Modules** | Git tag `v*` | `go install ...@v1.0.0` |
| **Homebrew** | Manual update | Formula update needed |
| **NPM** | Manual publish | `npm publish` |

---

## 🚀 Deployment Workflow

### For Development (Continuous)

```
Developer pushes to main
    ↓
GitHub Actions triggered
    ↓
[1] Run tests (Ubuntu + macOS)
    ├─ Unit tests
    ├─ Integration tests  
    └─ Linting
    ↓
[2] Build binaries (if tests pass)
    ├─ Linux (amd64, arm64)
    ├─ macOS (amd64, arm64)
    └─ Windows (amd64)
    ↓
[3] Create checksums
    └─ SHA256 for each binary
    ↓
[4] Update 'dev' release
    ├─ Delete old dev release
    ├─ Create new dev release
    ├─ Upload all binaries
    └─ Generate release notes
    ↓
[5] Build & push Docker images
    ├─ Multi-arch (amd64, arm64)
    ├─ Tag as 'dev'
    └─ Tag as 'latest' (if main)
    ↓
✅ Available via all channels within ~5 minutes
```

### For Production (On Tag)

```
Maintainer creates tag: v1.0.0
    ↓
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
    ↓
GitHub Actions triggered
    ↓
[1] Run tests
    ↓
[2] Build binaries (all platforms)
    ↓
[3] Generate checksums
    ↓
[4] Create GitHub Release
    ├─ Extract changelog for v1.0.0
    ├─ Add installation instructions
    ├─ Upload binaries
    └─ Upload checksums
    ↓
[5] Run integration tests
    └─ Test binaries on Ubuntu + macOS
    ↓
[6] Build & push Docker images
    ├─ Tag as v1.0.0
    ├─ Tag as 1.0
    ├─ Tag as 1
    └─ Tag as latest
    ↓
[7] Update Homebrew formula (optional)
    └─ Calculate new checksums
    ↓
✅ Stable release available
```

---

## 🔍 Testing Before Deployment

### Pull Request Checks

Every PR must pass:
1. ✅ All unit tests (24 tests)
2. ✅ Linting (golangci-lint)
3. ✅ Build verification
4. ✅ Code coverage maintained

### Main Branch Checks

Every push to main triggers:
1. ✅ Full test suite on 2 platforms
2. ✅ Build for 5 platforms
3. ✅ Integration tests
4. ✅ Binary verification

### Release Checks

Every tagged release includes:
1. ✅ All main branch checks
2. ✅ Multi-platform integration tests
3. ✅ Checksum verification
4. ✅ Installation script testing

---

## 📊 Monitoring Deployments

### GitHub Actions Dashboard

View all workflows:
```
https://github.com/YOLOVibeCode/opencode-enterprise-shield/actions
```

### Release Page

View all releases:
```
https://github.com/YOLOVibeCode/opencode-enterprise-shield/releases
```

### Docker Registry

View Docker images:
```
https://github.com/YOLOVibeCode/opencode-enterprise-shield/pkgs/container/opencode-enterprise-shield
```

---

## 🛠️ Manual Operations

### Force Update Development Build

```bash
# Trigger a new dev build
git commit --allow-empty -m "trigger: Force dev build"
git push origin main

# Wait ~5 minutes for build to complete
# New binaries available at /releases/tag/dev
```

### Create a New Release

```bash
# Update version in files
# Update CHANGELOG.md
# Commit changes

# Tag the release
git tag -a v1.0.1 -m "Release v1.0.1: Bug fixes"
git push origin v1.0.1

# GitHub Actions handles the rest automatically
```

### Rollback a Release

```bash
# Delete the tag locally and remotely
git tag -d v1.0.1
git push origin :refs/tags/v1.0.1

# Delete the GitHub release
gh release delete v1.0.1 --yes

# Users on install script will get previous version
```

---

## 🔐 Secrets & Permissions

### Required GitHub Secrets

| Secret | Purpose | Workflow |
|--------|---------|----------|
| `GITHUB_TOKEN` | Automatic (no setup needed) | All workflows |
| `HOMEBREW_TAP_TOKEN` | Update Homebrew formula | Release (optional) |
| `NPM_TOKEN` | Publish to NPM | Release (optional) |

### Required Permissions

Repository settings → Actions → General → Workflow permissions:
- ✅ **Read and write permissions**
- ✅ **Allow GitHub Actions to create and approve pull requests**

---

## 📈 Metrics & Analytics

### Build Success Rate

View in GitHub Actions:
- Success rate per workflow
- Average build time
- Platform-specific failures

### Download Statistics

View in GitHub Insights:
- Release download counts
- Asset popularity
- Geographic distribution

### Docker Pull Statistics

View in Container Registry:
- Pull count by tag
- Bandwidth usage
- Version distribution

---

## 🚨 Troubleshooting

### Build Fails on Main

1. Check Actions tab for error logs
2. Tests must pass before deployment
3. Fix locally and push again
4. CI automatically retries

### Dev Release Not Updating

1. Verify push reached `main` branch
2. Check GitHub Actions for failures
3. Ensure `dev` tag is deletable (not protected)
4. Force refresh: empty commit + push

### Docker Build Fails

1. Check Dockerfile syntax
2. Verify multi-arch support
3. Check GHCR permissions
4. Review build logs in Actions

---

## 📋 Deployment Checklist

### For Each Main Push

- [ ] Tests pass locally (`make test`)
- [ ] Linter passes (`golangci-lint run`)
- [ ] Commit message is clear
- [ ] Push to main triggers CI
- [ ] Monitor Actions for success
- [ ] Verify `dev` release updates

### For Each Release

- [ ] Version updated in all files
- [ ] CHANGELOG.md updated
- [ ] Tests pass locally
- [ ] Commit version bump
- [ ] Create annotated tag
- [ ] Push tag to trigger release
- [ ] Monitor release workflow
- [ ] Verify release created
- [ ] Test install script
- [ ] Update Homebrew formula (if needed)
- [ ] Announce release

---

## 🎯 Best Practices

### Branching Strategy

```
main (protected)
  ├── feature/new-detector (PR required)
  ├── fix/bug-123 (PR required)
  └── docs/improve-readme (PR required)
```

**Rules:**
- All changes via Pull Request
- PRs require passing tests
- Main is always deployable
- Tags create stable releases

### Commit Messages

```bash
# Good commit messages
feat: Add credit card detection with Luhn validation
fix: Correct session expiry calculation
docs: Update installation instructions
chore: Upgrade Go to 1.23

# Triggers
feat: Triggers full CI + dev deployment
fix: Triggers full CI + dev deployment
docs: Triggers full CI only (no build)
chore: Triggers full CI + dev deployment
```

### Versioning

Follow Semantic Versioning (semver):
- `v1.0.0` - Major release (breaking changes)
- `v1.1.0` - Minor release (new features)
- `v1.0.1` - Patch release (bug fixes)

---

## 📞 Support

For CI/CD issues:
- Check [GitHub Actions](https://github.com/YOLOVibeCode/opencode-enterprise-shield/actions)
- Review workflow logs
- Contact: devops@YOLOVibeCode.com

---

*Automated deployment ensures every commit is tested and ready for distribution* 🚀

