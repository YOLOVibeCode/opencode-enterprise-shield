# Versioning System

## Overview

Enterprise Shield uses **automated semantic versioning** with incremental build numbers.

---

## 📊 Version Format

### Semantic Version: `MAJOR.MINOR.PATCH`

Stored in `VERSION` file:
```
1.0.0
```

### Full Version: `MAJOR.MINOR.PATCH+build.NUMBER`

Used in binaries:
```
1.0.0+build.42
1.2.3+build.128
```

### Components

| Component | Meaning | Example | When to Increment |
|-----------|---------|---------|-------------------|
| **MAJOR** | Breaking changes | `2.0.0` | API changes, incompatible updates |
| **MINOR** | New features | `1.1.0` | New functionality, backward compatible |
| **PATCH** | Bug fixes | `1.0.1` | Bug fixes, security patches |
| **BUILD** | Build number | `+build.42` | **Auto-incremented every commit** |

---

## 🤖 Automatic Versioning

### On Every Push to Main

**GitHub Actions automatically:**

1. ✅ Reads current version from `VERSION` file
2. ✅ Analyzes commit message for version bump:
   - `major:` or `BREAKING CHANGE:` → Bump MAJOR
   - `feat:` or `feature:` → Bump MINOR
   - `fix:` or `patch:` → Bump PATCH
   - Other → Keep version, increment build only
3. ✅ Increments build number (uses GitHub run number)
4. ✅ Generates full version: `1.0.0+build.42`
5. ✅ Updates `VERSION` file (if version changed)
6. ✅ Commits back to repo with `[skip ci]`
7. ✅ Builds binaries with version embedded

**Example:**

```bash
# Commit with "feat: Add new detection pattern"
$ git commit -m "feat: Add credit card detection"
$ git push origin main

# GitHub Actions:
# - Current: 1.0.0
# - Detected: "feat:" → MINOR bump
# - New: 1.1.0+build.23
# - Updates VERSION to 1.1.0
# - Builds binaries with v1.1.0+build.23
```

---

## 📝 Commit Message Conventions

### Version Bump Triggers

| Prefix | Version Change | Example |
|--------|---------------|---------|
| `major:` | `1.0.0` → `2.0.0` | `major: Redesign API structure` |
| `BREAKING CHANGE:` | `1.0.0` → `2.0.0` | `BREAKING CHANGE: Remove legacy mode` |
| `feat:` | `1.0.0` → `1.1.0` | `feat: Add SSN detection` |
| `feature:` | `1.0.0` → `1.1.0` | `feature: Support Redis backend` |
| `minor:` | `1.0.0` → `1.1.0` | `minor: New compliance patterns` |
| `fix:` | `1.0.0` → `1.0.1` | `fix: Correct regex timeout` |
| `patch:` | `1.0.0` → `1.0.1` | `patch: Fix session expiry` |
| `docs:` | No change | `docs: Update README` |
| `chore:` | No change | `chore: Update dependencies` |
| `ci:` | No change | `ci: Fix workflow` |

**Build number always increments** regardless of version change.

### Examples

```bash
# MAJOR bump (breaking change)
git commit -m "major: Redesign sanitization API"
# Result: 1.0.0 → 2.0.0+build.15

# MINOR bump (new feature)
git commit -m "feat: Add GPU detection pattern"
# Result: 1.0.0 → 1.1.0+build.16

# PATCH bump (bug fix)
git commit -m "fix: Correct alias counter overflow"
# Result: 1.1.0 → 1.1.1+build.17

# No version bump (documentation)
git commit -m "docs: Update installation guide"
# Result: 1.1.1+build.18 (version stays, build increments)
```

---

## 🔧 Manual Version Management

### Bump Version Manually

```bash
# Bump patch (1.0.0 → 1.0.1)
make bump-patch

# Bump minor (1.0.1 → 1.1.0)
make bump-minor

# Bump major (1.1.0 → 2.0.0)
make bump-major

# Each command:
# 1. Updates VERSION file
# 2. Shows you what to commit
```

### Check Current Version

```bash
# View version info
make version

# Output:
# Version: 1.0.0
# Build: 42
# Full: 1.0.0+build.42
```

### Create Stable Release

```bash
# Option 1: Use make tag
make tag V=1.2.0

# Option 2: Manual
echo "1.2.0" > VERSION
git add VERSION
git commit -m "chore: Release v1.2.0"
git tag -a v1.2.0 -m "Release v1.2.0"
git push origin main v1.2.0
```

---

## 📦 Version in Binaries

### Build Information Embedded

Every binary includes version information accessible via:

```bash
$ enterprise-shield version
Enterprise Shield Plugin v1.0.0+build.42
Build: #42
Built: 2026-01-14_00:30:15
```

### Embedded via ldflags

```go
// cmd/plugin/main.go
var (
    version     = "dev"        // Set by -ldflags "-X main.version=1.0.0+build.42"
    buildNumber = "0"          // Set by -ldflags "-X main.buildNumber=42"
    buildTime   = "unknown"    // Set by -ldflags "-X main.buildTime=..."
)
```

---

## 📈 Build Number Tracking

### Build Number Sources

| Environment | Build Number Source |
|-------------|-------------------|
| **GitHub Actions** | `${{ github.run_number }}` (auto-increments) |
| **Local Development** | Git commit count: `git rev-list --count HEAD` |
| **Manual Build** | Can be set with `BUILD_NUMBER=123 make build` |

### Build Number Properties

- ✅ **Always increments** (never decreases)
- ✅ **Unique per build** (no collisions)
- ✅ **Traceable** (links to CI run)
- ✅ **Independent** of version (version can stay same)

---

## 🔄 Version Flow Examples

### Scenario 1: Bug Fix Development

```bash
# Current: 1.0.0, Build: 10
$ git commit -m "fix: Correct regex timeout bug"
$ git push origin main

# GitHub Actions:
# → Detects "fix:" prefix
# → Bumps PATCH: 1.0.0 → 1.0.1
# → Increments build: 10 → 11
# → Full version: 1.0.1+build.11
# → Updates VERSION file to 1.0.1
# → Builds binaries with v1.0.1+build.11
# → Creates dev release
```

### Scenario 2: New Feature Development

```bash
# Current: 1.0.1, Build: 11
$ git commit -m "feat: Add Azure Key Vault integration"
$ git push origin main

# GitHub Actions:
# → Detects "feat:" prefix
# → Bumps MINOR: 1.0.1 → 1.1.0
# → Increments build: 11 → 12
# → Full version: 1.1.0+build.12
# → Updates VERSION file to 1.1.0
# → Builds binaries with v1.1.0+build.12
# → Creates dev release
```

### Scenario 3: Documentation Update

```bash
# Current: 1.1.0, Build: 12
$ git commit -m "docs: Update README with examples"
$ git push origin main

# GitHub Actions:
# → No version bump keyword detected
# → Version stays: 1.1.0
# → Increments build: 12 → 13
# → Full version: 1.1.0+build.13
# → Does NOT update VERSION file
# → Builds binaries with v1.1.0+build.13
# → Creates dev release
```

### Scenario 4: Breaking Change

```bash
# Current: 1.1.0, Build: 13
$ git commit -m "BREAKING CHANGE: Redesign configuration format"
$ git push origin main

# GitHub Actions:
# → Detects "BREAKING CHANGE:" keyword
# → Bumps MAJOR: 1.1.0 → 2.0.0
# → Resets MINOR and PATCH to 0
# → Increments build: 13 → 14
# → Full version: 2.0.0+build.14
# → Updates VERSION file to 2.0.0
# → Builds binaries with v2.0.0+build.14
# → Creates dev release
```

---

## 🎯 Stable Release Process

### Creating a Stable Release

```bash
# Option 1: Using make (recommended)
make tag V=1.2.0

# Option 2: Manual
echo "1.2.0" > VERSION
git add VERSION
git commit -m "chore: Release v1.2.0"
git tag -a v1.2.0 -m "Release v1.2.0"
git push origin main v1.2.0
```

**What happens:**
1. `VERSION` file updated to `1.2.0`
2. Tag `v1.2.0` created
3. GitHub Actions triggered by tag
4. Builds binaries with version `v1.2.0` (no +build suffix for stable)
5. Creates stable release (not prerelease)
6. Docker tagged with `:1.2.0`, `:1.2`, `:1`, `:latest`

---

## 📊 Version History Tracking

### View Version History

```bash
# See all versions
git tag -l "v*"

# See version from commits
git log --oneline --grep="chore: Release" --grep="chore: Bump"

# See current version
cat VERSION
```

### Version in CHANGELOG

```markdown
## [1.2.0] - 2026-01-14

### Added
- New feature X

## [1.1.0] - 2026-01-13

### Added
- Feature Y

## [1.0.0] - 2026-01-13

### Added
- Initial release
```

---

## 🔍 Troubleshooting

### Build Number Not Incrementing

**Problem:** Build number stays same

**Solution:** Build number is based on GitHub Actions run number. Local builds use git commit count:
```bash
git rev-list --count HEAD
```

### Version Not Updating

**Problem:** VERSION file not updating automatically

**Solution:** Check commit message has correct prefix:
```bash
# Won't bump
git commit -m "Updated something"

# Will bump (minor)
git commit -m "feat: Updated something"
```

### Want to Force Specific Version

```bash
# Set version manually
echo "2.5.0" > VERSION
git add VERSION
git commit -m "chore: Set version to 2.5.0"
git push
```

---

## 📋 Best Practices

### ✅ Do

- Use conventional commit messages (`feat:`, `fix:`, `major:`)
- Let automation handle version bumps
- Create git tags for stable releases
- Update CHANGELOG.md when creating releases
- Test locally before pushing

### ❌ Don't

- Manually edit version in code (use VERSION file)
- Skip commit message prefixes (version won't bump)
- Create tags without updating VERSION file
- Push with `[skip ci]` unless intentional

---

## 🎯 Quick Reference

```bash
# Check version
make version

# Bump version
make bump-patch    # 1.0.0 → 1.0.1
make bump-minor    # 1.0.0 → 1.1.0
make bump-major    # 1.0.0 → 2.0.0

# Create release
make tag V=1.2.0
git push origin main v1.2.0

# View in binary
./build/enterprise-shield version
```

---

*Semantic versioning + automated build numbers = Perfect traceability* 🎯

