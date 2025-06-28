# Publishing DeployAja CLI to Homebrew

This guide explains how to publish the DeployAja CLI to Homebrew, either through a custom tap or to the official homebrew-core.

## Overview

The setup includes:
- ✅ **Cross-platform builds** for macOS (Intel/Apple Silicon) and Linux (x64/ARM64)
- ✅ **Automated releases** with GitHub Actions
- ✅ **Homebrew formula** with platform-specific binaries
- ✅ **Automated tap updates** when new releases are published
- ✅ **Checksum verification** for security

## Quick Start

### Option 1: Custom Homebrew Tap (Recommended for initial releases)

1. **Create a homebrew-tap repository:**
   ```bash
   # Create a new repository named 'homebrew-tap' under your organization
   # Example: https://github.com/deployaja/homebrew-tap
   ```

2. **Set up the TAP_TOKEN secret:**
   - Go to GitHub Settings → Developer settings → Personal access tokens
   - Create a token with `repo` permissions
   - Add it as `TAP_TOKEN` secret in your main repository

3. **Create your first release:**
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

4. **Users can then install with:**
   ```bash
   brew install deployaja/tap/aja
   ```

### Option 2: Submit to homebrew-core (For mature, popular tools)

After your tool is stable and has some adoption, you can submit to homebrew-core for wider distribution.

## Detailed Setup

### 1. Repository Structure

Your main repository should have:
```
deployaja-cli/
├── Formula/aja.rb              # Homebrew formula
├── scripts/update-homebrew.sh  # Manual update script
├── .github/workflows/
│   ├── release.yml             # Creates releases with binaries
│   └── update-homebrew.yml     # Updates formula automatically
└── ... (your source code)
```

### 2. Release Process

The release process is fully automated:

1. **Tag a release:**
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

2. **GitHub Actions will:**
   - Build binaries for all platforms
   - Create a GitHub release with assets
   - Generate checksums
   - Update the Homebrew formula
   - Push updates to your homebrew-tap repository (if TAP_TOKEN is set)

### 3. Manual Formula Updates

If you need to manually update the formula:

```bash
# Update formula for a specific version
./scripts/update-homebrew.sh v1.0.0

# Review changes
git diff Formula/aja.rb

# Commit and push
git add Formula/aja.rb
git commit -m "Update Homebrew formula to v1.0.0"
git push
```

### 4. Testing the Formula

Before publishing, test your formula:

```bash
# Test local formula
brew install --build-from-source Formula/aja.rb

# Test from your tap
brew install deployaja/tap/aja

# Run tests
aja version
aja --help
```

### 5. Homebrew Tap Repository Setup

If using a custom tap, create a repository named `homebrew-tap` under your organization:

```bash
# Clone and set up your tap repository
git clone https://github.com/deployaja/homebrew-tap.git
cd homebrew-tap

# Add the formula
mkdir -p Formula
cp ../deployaja-cli/Formula/aja.rb Formula/aja.rb

# Commit and push
git add Formula/aja.rb
git commit -m "Add aja formula"
git push
```

### 6. Users Installation Instructions

With your custom tap:
```bash
# Add the tap
brew tap deployaja/tap

# Install the CLI
brew install aja

# Or in one command
brew install deployaja/tap/aja
```

## Submitting to homebrew-core

Once your tool is stable and has good adoption, you can submit to homebrew-core:

### Prerequisites
- ✅ Tool is stable and well-tested
- ✅ Good documentation and user base
- ✅ Follows Homebrew guidelines
- ✅ Open source license
- ✅ No GUI components (CLI only)
- ✅ Not a duplicate of existing tools

### Submission Process

1. **Fork homebrew-core:**
   ```bash
   git clone https://github.com/Homebrew/homebrew-core.git
   cd homebrew-core
   git checkout -b add-aja
   ```

2. **Copy your formula:**
   ```bash
   cp ../deployaja-cli/Formula/aja.rb Formula/aja.rb
   ```

3. **Test thoroughly:**
   ```bash
   brew install --build-from-source Formula/aja.rb
   brew test aja
   brew audit --strict aja
   ```

4. **Create pull request:**
   - Follow the [Homebrew contributing guidelines](https://docs.brew.sh/How-To-Open-a-Homebrew-Pull-Request)
   - Include a clear description of what your tool does
   - Ensure all CI checks pass

## Maintenance

### Automated Updates
- New releases automatically trigger formula updates
- Checksums are automatically calculated and updated
- Both main repo and tap repository are updated

### Manual Maintenance
- Monitor Homebrew guidelines for changes
- Update formula syntax if needed
- Respond to user issues promptly

### Version Management
- Use semantic versioning (v1.0.0, v1.1.0, v2.0.0)
- Update version in `internal/version/version.go` for dev builds
- Actual version is injected during build time

## Troubleshooting

### Common Issues

1. **Checksum mismatch:**
   ```bash
   # Regenerate checksums
   ./scripts/update-homebrew.sh v1.0.0
   ```

2. **Formula syntax errors:**
   ```bash
   # Test formula syntax
   brew audit --strict Formula/aja.rb
   ```

3. **Binary not found:**
   - Check that release assets are properly uploaded
   - Verify URL structure matches formula

4. **TAP_TOKEN issues:**
   - Ensure token has `repo` permissions
   - Check token is not expired
   - Verify repository exists and is accessible

### Debug Commands

```bash
# Check what Homebrew is downloading
brew --cache aja

# Verbose install for debugging
brew install --verbose deployaja/tap/aja

# Check formula syntax
brew audit aja

# Test formula
brew test aja
```

## Resources

- [Homebrew Formula Cookbook](https://docs.brew.sh/Formula-Cookbook)
- [Homebrew Acceptable Formulae](https://docs.brew.sh/Acceptable-Formulae)
- [Creating Homebrew Taps](https://docs.brew.sh/How-to-Create-and-Maintain-a-Tap)

## Support

If you encounter issues with the Homebrew publishing process:
1. Check the GitHub Actions logs
2. Test the formula locally
3. Review Homebrew documentation
4. Open an issue in the repository 