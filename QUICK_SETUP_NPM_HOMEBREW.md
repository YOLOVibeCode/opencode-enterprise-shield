# Quick Setup: NPM & Homebrew Auto-Publishing

## 🎯 Overview

Two optional distribution channels need one-time setup:
1. **Homebrew** - Popular for macOS/Linux users (5 minutes)
2. **NPM** - For Node.js developers (3 minutes)

**Once configured, both auto-publish on every release tag!**

---

## 🍺 Homebrew Setup (5 Minutes)

### Step 1: Create Homebrew Tap Repository

**Go to GitHub and create a new repository:**
- URL: https://github.com/new
- Repository name: `homebrew-opencode-enterprise-shield`
- Description: "Homebrew tap for Enterprise Shield"
- Visibility: **Public** (Homebrew requires public taps)
- **Don't** initialize with README

### Step 2: Push Formula to Tap

```bash
# Clone your new tap repository
cd /tmp
git clone https://github.com/YOLOVibeCode/homebrew-opencode-enterprise-shield.git
cd homebrew-opencode-enterprise-shield

# Create Formula directory
mkdir Formula

# Copy the formula
cp /Users/admin/Dev/YOLOProjects/opencode-enterprise-shield/Formula/enterprise-shield.rb Formula/

# Create README
cat > README.md << 'EOF'
# Homebrew Tap for Enterprise Shield

Enterprise-grade security plugin for OpenCode AI assistants.

## Installation

```bash
brew tap YOLOVibeCode/opencode-enterprise-shield
brew install enterprise-shield
```

## Updating

```bash
brew update
brew upgrade enterprise-shield
```

## Uninstalling

```bash
brew uninstall enterprise-shield
brew untap YOLOVibeCode/opencode-enterprise-shield
```

## About

See main repository: https://github.com/YOLOVibeCode/opencode-enterprise-shield
EOF

# Commit and push
git add .
git commit -m "Initial Homebrew formula for enterprise-shield"
git push origin main
```

### Step 3: Configure GitHub Secret (Optional)

**Option A: Use default GITHUB_TOKEN (easiest, no setup needed)**
- Already works if tap is in same organization (YOLOVibeCode)
- No additional configuration needed!
- ✅ **Recommended for same-org setups**

**Option B: Create dedicated token (more control)**

1. Go to GitHub Settings → Developer settings → Personal access tokens → Tokens (classic)
2. Click "Generate new token (classic)"
3. Settings:
   - Name: `HOMEBREW_TAP_TOKEN`
   - Expiration: No expiration (or 1 year)
   - Scope: Check **only** `public_repo`
4. Click "Generate token"
5. **Copy the token** (you won't see it again!)

6. Add to your main repository:
   - Go to: https://github.com/YOLOVibeCode/opencode-enterprise-shield/settings/secrets/actions
   - Click "New repository secret"
   - Name: `HOMEBREW_TAP_TOKEN`
   - Value: Paste the token
   - Click "Add secret"

### Step 4: Test

```bash
# Go back to main repo
cd /Users/admin/Dev/YOLOProjects/opencode-enterprise-shield

# Manually trigger the Homebrew workflow to test
gh workflow run publish-homebrew.yml --repo YOLOVibeCode/opencode-enterprise-shield -f version=v1.0.0

# Check if it worked
# Visit: https://github.com/YOLOVibeCode/homebrew-opencode-enterprise-shield/commits
# Should see: "Update enterprise-shield to v1.0.0"
```

### Step 5: Users Can Now Install

```bash
brew tap YOLOVibeCode/opencode-enterprise-shield
brew install enterprise-shield
```

**Future releases auto-update the formula!**

---

## 📦 NPM Setup (3 Minutes)

### Step 1: Create NPM Account (if needed)

**If you don't have an NPM account:**
1. Go to https://www.npmjs.com/signup
2. Create free account
3. Verify email

**If you have an account:**
```bash
npm login
# Enter your credentials
```

### Step 2: Verify Package Name is Available

```bash
npm search @yolovibeCode/opencode-enterprise-shield

# If "No matches found" → Name is available ✅
# If found → Choose different name
```

**Alternative: Use unscoped name** (easier for free accounts)

Update `package.json`:
```json
{
  "name": "opencode-enterprise-shield",  // Remove @yolovibeCode/
  ...
}
```

### Step 3: Generate NPM Token

1. Go to https://www.npmjs.com/settings/YOUR_USERNAME/tokens
2. Click "Generate New Token" → "Classic Token"
3. Token type: **Automation** (important!)
4. Click "Generate Token"
5. **Copy the token** (starts with `npm_...`)

### Step 4: Add Token to GitHub

1. Go to: https://github.com/YOLOVibeCode/opencode-enterprise-shield/settings/secrets/actions
2. Click "New repository secret"
3. Name: `NPM_TOKEN` (exact name)
4. Value: Paste the token
5. Click "Add secret"

### Step 5: Test

```bash
# Manually trigger NPM publish workflow to test
cd /Users/admin/Dev/YOLOProjects/opencode-enterprise-shield

gh workflow run publish-npm.yml --repo YOLOVibeCode/opencode-enterprise-shield -f version=v1.0.0

# Watch workflow
gh run watch

# Check NPM registry (after workflow completes)
npm view @yolovibeCode/opencode-enterprise-shield version
# Should show: 1.0.0
```

### Step 6: Users Can Now Install

```bash
npm install -g @yolovibeCode/opencode-enterprise-shield

# Or if using unscoped:
npm install -g opencode-enterprise-shield
```

**Future releases auto-publish to NPM!**

---

## 🚀 Even Quicker: Test with Current Release

### Homebrew

```bash
# 1. Create tap repo (1 command)
gh repo create YOLOVibeCode/homebrew-opencode-enterprise-shield --public -y

# 2. Clone and set up
cd /tmp
git clone https://github.com/YOLOVibeCode/homebrew-opencode-enterprise-shield.git
cd homebrew-opencode-enterprise-shield
mkdir Formula
cp /Users/admin/Dev/YOLOProjects/opencode-enterprise-shield/Formula/enterprise-shield.rb Formula/
cat > README.md << 'EOF'
# Homebrew Tap for Enterprise Shield
brew tap YOLOVibeCode/opencode-enterprise-shield && brew install enterprise-shield
EOF
git add . && git commit -m "Initial formula" && git push

# 3. Trigger update for v1.0.0
cd /Users/admin/Dev/YOLOProjects/opencode-enterprise-shield
gh workflow run publish-homebrew.yml -f version=v1.0.0

# Done! Formula will auto-update on future releases
```

### NPM

```bash
# 1. Login to NPM
npm login

# 2. Create token at: https://www.npmjs.com/settings/YOUR_USERNAME/tokens
#    - Type: Automation
#    - Copy token

# 3. Add to GitHub
gh secret set NPM_TOKEN --repo YOLOVibeCode/opencode-enterprise-shield
# Paste token when prompted

# 4. Trigger publish for v1.0.0
gh workflow run publish-npm.yml -f version=v1.0.0

# Done! NPM will auto-publish on future releases
```

---

## 📋 Decision Matrix

### Option A: Set Up Both (Recommended)

**Time:** 8 minutes total
**Benefit:** Users can install via **all 6 methods**
**Future:** Fully automated forever

```bash
# Quick script to set up both
cd /Users/admin/Dev/YOLOProjects/opencode-enterprise-shield
./scripts/setup-all-channels.sh  # (I can create this script if you want)
```

### Option B: Set Up Homebrew Only

**Time:** 5 minutes
**Benefit:** macOS/Linux users happy (large audience)
**Future:** Auto-updates on releases

```bash
# Just run the Homebrew setup steps above
```

### Option C: Set Up NPM Only

**Time:** 3 minutes
**Benefit:** Node.js developers happy
**Future:** Auto-publishes on releases

```bash
# Just run the NPM setup steps above
```

### Option D: Skip Both (Keep Current State)

**Time:** 0 minutes
**Benefit:** 4 working methods already (GitHub, Docker, Go, Install Script)
**Future:** Can set up anytime later

**Users can still install via:**
- ✅ One-line install script
- ✅ Go install
- ✅ Docker
- ✅ Manual download

---

## 🎯 My Recommendation

### For Maximum Reach: Set Up Both

**Why Homebrew:**
- Very popular in developer community
- macOS users expect it
- "brew install" is the easiest command
- Takes 5 minutes

**Why NPM:**
- Huge Node.js/JavaScript community
- Many OpenCode users likely use Node
- "npm install -g" is familiar
- Takes 3 minutes

**Combined:** 8 minutes of setup → Automated forever

---

## 🛠️ I Can Help Set These Up

### Want me to create automated setup scripts?

I can create:
1. `scripts/setup-homebrew-tap.sh` - Automates tap creation
2. `scripts/setup-npm-token.sh` - Guides NPM token setup
3. `scripts/setup-all-channels.sh` - Does both

### Want me to guide you through it?

Just say:
- "Set up Homebrew" - I'll guide you step-by-step
- "Set up NPM" - I'll guide you step-by-step
- "Set up both" - I'll create automation scripts
- "Skip for now" - Current release works fine as-is

---

## ✅ Current Status Summary

**Your release v1.0.0 is LIVE and functional!**

**Working NOW (no setup needed):**
- ✅ GitHub releases
- ✅ Install script (one-line)
- ✅ Go install
- ✅ Docker images

**Optional to enable:**
- ⏸️ Homebrew (5 min setup)
- ⏸️ NPM (3 min setup)

**Users can already install and use Enterprise Shield via 4 methods!**

---

**What would you like to do?**
1. Set up both Homebrew & NPM now (8 minutes)
2. Set up just Homebrew (5 minutes)
3. Set up just NPM (3 minutes)
4. Skip for now (release works fine as-is)

Let me know and I can help with whichever you choose!

