# GitHub Secrets Setup Guide

To enable full automation of Homebrew and NPM publishing, you need to configure GitHub secrets.

---

## 📋 Required Secrets

### 1. **HOMEBREW_TAP_TOKEN** (Required for Homebrew auto-publish)

**Purpose:** Allows the main repository to push formula updates to your Homebrew tap repository.

**Setup Steps:**

1. **Create Personal Access Token (Classic):**
   - Go to GitHub → Settings → Developer settings → Personal access tokens → Tokens (classic)
   - Click "Generate new token (classic)"
   - Name: `HOMEBREW_TAP_TOKEN`
   - Expiration: No expiration (or 1 year with renewal reminder)
   - Scopes: Select **only** `public_repo` (or `repo` if tap is private)
   - Click "Generate token"
   - **Copy the token** (you won't see it again!)

2. **Add Secret to Main Repository:**
   - Go to your `opencode-enterprise-shield` repository
   - Settings → Secrets and variables → Actions
   - Click "New repository secret"
   - Name: `HOMEBREW_TAP_TOKEN`
   - Value: Paste the token you copied
   - Click "Add secret"

**Alternative (If tap is in same org):**
- You can use `${{ secrets.GITHUB_TOKEN }}` instead
- No extra secret needed
- Workflow will use the default GITHUB_TOKEN

---

### 2. **NPM_TOKEN** (Required for NPM auto-publish)

**Purpose:** Allows GitHub Actions to publish packages to the NPM registry.

**Setup Steps:**

1. **Create NPM Access Token:**
   - Go to [npmjs.com](https://www.npmjs.com/)
   - Log in to your account
   - Click your profile → Access Tokens
   - Click "Generate New Token" → "Classic Token"
   - Type: **Automation** (not "Publish")
   - Click "Generate Token"
   - **Copy the token**

2. **Add Secret to GitHub Repository:**
   - Go to your `opencode-enterprise-shield` repository
   - Settings → Secrets and variables → Actions
   - Click "New repository secret"
   - Name: `NPM_TOKEN`
   - Value: Paste the NPM token
   - Click "Add secret"

3. **Verify NPM Package Name is Available:**
   ```bash
   npm search @YOLOVibeCode/opencode-enterprise-shield
   ```
   - If not found, name is available
   - If found, choose different name or claim it

---

## 🔐 Security Best Practices

### Token Permissions

| Secret | Minimum Scope | Why |
|--------|--------------|-----|
| `HOMEBREW_TAP_TOKEN` | `public_repo` | Only needs to push to tap repo |
| `NPM_TOKEN` | Automation | Allows publishing packages |

### Token Rotation

**Recommended schedule:**
- **HOMEBREW_TAP_TOKEN**: Rotate every 90 days (or set 1-year expiration)
- **NPM_TOKEN**: Rotate every 90 days

**How to rotate:**
1. Generate new token with same permissions
2. Update secret in GitHub repository settings
3. Delete old token

### Monitoring

**Check secret usage:**
- GitHub → Repository → Settings → Secrets → Actions
- Each secret shows "Last used" timestamp
- Review regularly for unauthorized access

---

## 🧪 Testing Secrets

### Test Homebrew Publishing

```bash
# Manually trigger workflow to test
gh workflow run publish-homebrew.yml \
  --ref main \
  -f version=v1.0.0
```

**Check:**
1. Go to Actions tab
2. Find "Update Homebrew Formula" workflow run
3. Check for errors
4. Verify tap repository was updated

### Test NPM Publishing

```bash
# Manually trigger workflow to test
gh workflow run publish-npm.yml \
  --ref main \
  -f version=v1.0.0
```

**Check:**
1. Go to Actions tab
2. Find "Publish to NPM" workflow run
3. Check for errors
4. Verify package appears on npmjs.com

---

## 🚀 One-Time Setup Checklist

### Homebrew Automation

- [ ] Create Homebrew tap repository: `homebrew-opencode-enterprise-shield`
- [ ] Copy `Formula/enterprise-shield.rb` to tap
- [ ] Push formula to tap repository
- [ ] Create `HOMEBREW_TAP_TOKEN` (or use default `GITHUB_TOKEN`)
- [ ] Add secret to main repository
- [ ] Test with manual workflow trigger

### NPM Automation

- [ ] Create NPM account (if needed)
- [ ] Verify package name is available
- [ ] Generate NPM automation token
- [ ] Add `NPM_TOKEN` secret to GitHub
- [ ] Test with manual workflow trigger
- [ ] Verify package appears on npmjs.com

---

## 📊 Verification

### After Setup, Verify:

```bash
# 1. Check secrets exist
# GitHub → Repository → Settings → Secrets → Actions
# Should see:
#   - HOMEBREW_TAP_TOKEN ✅ (or using GITHUB_TOKEN)
#   - NPM_TOKEN ✅

# 2. Create a test release
git tag -a v1.0.1-test -m "Test release"
git push origin v1.0.1-test

# 3. Watch workflows
# https://github.com/YOLOVibeCode/opencode-enterprise-shield/actions
# Should see 3 workflows run:
#   - Release (builds binaries)
#   - Update Homebrew Formula
#   - Publish to NPM

# 4. Verify Homebrew formula updated
# https://github.com/YOLOVibeCode/homebrew-opencode-enterprise-shield/commits/main
# Should see commit: "Update enterprise-shield to v1.0.1-test"

# 5. Verify NPM package published
npm view @YOLOVibeCode/opencode-enterprise-shield version
# Should show: 1.0.1-test

# 6. Clean up test release
gh release delete v1.0.1-test --yes
git tag -d v1.0.1-test
git push origin :v1.0.1-test
```

---

## 🔧 Troubleshooting

### Homebrew Formula Update Fails

**Error:** "Authentication failed"

**Solution:**
- Verify `HOMEBREW_TAP_TOKEN` has `public_repo` scope
- Check token hasn't expired
- Ensure tap repository exists and is accessible

**Error:** "remote: Permission denied"

**Solution:**
- Token needs write access to tap repository
- If tap is in different org, token must be from that org
- Consider using GitHub App instead of PAT for cross-org

### NPM Publish Fails

**Error:** "401 Unauthorized"

**Solution:**
- Verify `NPM_TOKEN` is correct
- Check token hasn't expired
- Token must be "Automation" type

**Error:** "403 Forbidden: package name taken"

**Solution:**
- Package name already exists
- Either claim it or choose different name
- Update `package.json` name field

**Error:** "402 Payment Required"

**Solution:**
- Trying to publish scoped package (@YOLOVibeCode/...) to free account
- Either publish unscoped or upgrade to paid NPM account
- For organizations, need NPM Orgs subscription

---

## 🎯 Alternative: No Secrets Needed

If you don't want to set up secrets:

### Homebrew

**Manual update script** (already provided):
```bash
# After each release
./scripts/update-homebrew-formula.sh v1.0.0
cd ../homebrew-opencode-enterprise-shield
git push
```

### NPM

**Manual publish:**
```bash
# After each release
npm version $(cat VERSION) --no-git-tag-version
npm publish
```

**Still easy, just not automatic.**

---

## ✅ Recommended Setup

### Minimal (Free, No Extra Secrets)

```
GitHub Releases   ✅ Auto (GITHUB_TOKEN - automatic)
Docker           ✅ Auto (GITHUB_TOKEN - automatic)  
Go Install       ✅ Auto (public repository)
Install Script   ✅ Auto (public repository)
Homebrew         ⚠️ Manual (run script after release)
NPM              ⚠️ Manual (npm publish after release)
```

### Full Automation (Requires Secrets)

```
GitHub Releases   ✅ Auto (GITHUB_TOKEN)
Docker           ✅ Auto (GITHUB_TOKEN)
Go Install       ✅ Auto (public repo)
Install Script   ✅ Auto (public repo)
Homebrew         ✅ Auto (HOMEBREW_TAP_TOKEN)
NPM              ✅ Auto (NPM_TOKEN)
```

**Cost:** Free (just need to create tokens)

---

## 📞 Support

If you encounter issues:
- Check workflow logs in Actions tab
- Verify secrets are correctly named
- Test with manual workflow triggers
- See troubleshooting section above

---

*Set up once, automated forever* 🔐

