# Homebrew Tap Setup

This document explains how to set up automated Homebrew formula publishing for releases.

## Prerequisites

1. A GitHub repository for your Homebrew tap (e.g., `dimelords/homebrew-tap`)
2. A GitHub Personal Access Token with appropriate permissions

## Setup Steps

### 1. Create Homebrew Tap Repository

If you don't already have a Homebrew tap:

```bash
# Create a new repository named "homebrew-tap"
gh repo create dimelords/homebrew-tap --public --description "Homebrew formulae for dimelords projects"

# Clone it
git clone https://github.com/dimelords/homebrew-tap.git
cd homebrew-tap

# Create Formula directory
mkdir -p Formula

# Create initial README
cat > README.md <<EOF
# Homebrew Tap for dimelords

## Installation

\`\`\`bash
brew tap dimelords/tap
brew install protoc-gen-utcp
\`\`\`
EOF

git add .
git commit -m "Initial commit"
git push
```

### 2. Create GitHub Personal Access Token

1. Go to https://github.com/settings/tokens/new
2. Token name: `HOMEBREW_TAP_TOKEN`
3. Expiration: Choose appropriate duration (e.g., 1 year)
4. Select scopes:
   - ✅ `repo` (Full control of private repositories)
     - Needed to push formula updates to homebrew-tap
5. Click "Generate token"
6. **Copy the token immediately** (you won't see it again)

### 3. Add Token to Repository Secrets

Add the token as a secret in your main project repository:

```bash
# Using GitHub CLI
gh secret set HOMEBREW_TAP_GITHUB_TOKEN -b "your_token_here"

# Or via web UI:
# 1. Go to https://github.com/dimelords/protoc-gen-utcp/settings/secrets/actions
# 2. Click "New repository secret"
# 3. Name: HOMEBREW_TAP_GITHUB_TOKEN
# 4. Value: Paste your token
# 5. Click "Add secret"
```

### 4. Update Release Workflow (Already Done)

The `.goreleaser.yml` is already configured to:
- Skip Homebrew publishing if token is not set
- Publish to `dimelords/homebrew-tap` when token is available

### 5. Verify Setup

After adding the token, the next release will automatically:
1. Build binaries for all platforms
2. Create GitHub release
3. **Push Homebrew formula to homebrew-tap**

Users can then install with:
```bash
brew tap dimelords/tap
brew install protoc-gen-utcp
```

## Updating Formula Manually

If needed, you can manually update the formula:

```bash
cd homebrew-tap
brew bump-formula-pr --url=https://github.com/dimelords/protoc-gen-utcp/archive/v0.3.0.tar.gz protoc-gen-utcp
```

## Troubleshooting

### Token Permissions Error

If you see `Resource not accessible by integration`:
- Token needs `repo` scope
- Token must be valid and not expired
- Repository name must match `.goreleaser.yml` configuration

### Formula Not Updated

Check:
1. Token is set as `HOMEBREW_TAP_GITHUB_TOKEN` secret
2. Repository `dimelords/homebrew-tap` exists
3. Release workflow completed successfully

### Testing Formula Locally

```bash
# Tap your repository
brew tap dimelords/tap

# Install from tap
brew install protoc-gen-utcp

# Test
protoc-gen-utcp -version
```

## Security Notes

- **Never commit tokens to git**
- Use GitHub Secrets for CI/CD
- Rotate tokens periodically
- Use fine-grained tokens when possible
- Limit token scope to minimum required permissions

## Alternative: Manual Homebrew Distribution

If you prefer not to use automated publishing, users can install via:

```bash
# Direct from release tarball
brew install https://github.com/dimelords/protoc-gen-utcp/releases/download/v0.3.0/protoc-gen-utcp_0.3.0_Darwin_arm64.tar.gz
```

Or via `go install`:
```bash
go install github.com/dimelords/protoc-gen-utcp/cmd/protoc-gen-utcp@latest
```
