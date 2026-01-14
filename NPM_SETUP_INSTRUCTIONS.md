# NPM Setup - Step-by-Step

## 🎯 Goal
Enable automated NPM publishing so future releases auto-publish to the NPM registry.

---

## 📋 Prerequisites

You need an NPM account. If you don't have one:
- Go to: https://www.npmjs.com/signup
- Create free account
- Verify email

---

## ⚡ Quick Steps (3 Minutes)

### Step 1: Create NPM Automation Token (2 minutes)

1. **Login to NPM:**
   - Go to: https://www.npmjs.com/login
   - Or run locally: `npm login`

2. **Generate token:**
   - Go to: https://www.npmjs.com/settings/YOUR_USERNAME/tokens
   - Click "Generate New Token"
   - Select "Classic Token"
   - Type: **Automation** (important! Not "Publish")
   - Click "Generate Token"
   - **Copy the token** (starts with `npm_...`)
   - ⚠️ Save it somewhere - you won't see it again!

### Step 2: Add Token to GitHub (1 minute)

1. **Open GitHub secrets page:**
   - Go to: https://github.com/YOLOVibeCode/opencode-enterprise-shield/settings/secrets/actions
   - Or click: Settings → Secrets and variables → Actions

2. **Add new secret:**
   - Click "New repository secret"
   - Name: `NPM_TOKEN` (exact name, case-sensitive)
   - Value: Paste your NPM token
   - Click "Add secret"

### Step 3: Test (Optional)

Trigger a test publish of v1.0.0:

```bash
gh workflow run publish-npm.yml \
  --repo YOLOVibeCode/opencode-enterprise-shield \
  -f version=v1.0.0
```

Watch the workflow:
```bash
gh run watch --repo YOLOVibeCode/opencode-enterprise-shield
```

Verify after completion:
```bash
npm view @yolovibeCode/opencode-enterprise-shield version
# Should show: 1.0.0
```

---

## ✅ After Setup

**Future releases automatically publish to NPM!**

When you create v1.0.1:
```bash
git tag -a v1.0.1 -m "Release v1.0.1"
git push origin v1.0.1
```

GitHub Actions will automatically:
1. Build all binaries
2. Create GitHub release
3. Publish to NPM registry ✨
4. Update Homebrew formula ✨
5. Push Docker images
6. Update all 6 channels

Users can install with:
```bash
npm install -g @yolovibeCode/opencode-enterprise-shield
```

---

## 🔍 Troubleshooting

**Problem:** "Package name already taken"

**Solution:** Either:
- Use unscoped name: Change package.json to `"name": "opencode-enterprise-shield"`
- Or pick different name: `@yolovibeCode/enterprise-shield-plugin`

**Problem:** "Authentication failed"

**Solution:**
- Verify token is type "Automation" (not "Publish")
- Check token wasn't revoked
- Try generating new token

**Problem:** "403 Forbidden"

**Solution:**
- Scoped packages (@org/name) require paid account
- Either upgrade or use unscoped name

---

## 📝 Quick Checklist

- [ ] Have NPM account
- [ ] Generate automation token (not publish token)
- [ ] Copy token securely
- [ ] Add NPM_TOKEN to GitHub secrets
- [ ] Test with manual workflow trigger
- [ ] Verify package published

---

**See also:** `docs/SECRETS_SETUP.md` for detailed troubleshooting
