# GitFlow TUI Release Checklist

## Pre-Release

- [ ] All tests passing (`make test`)
- [ ] Code formatted (`make fmt`)
- [ ] Linter passing (`make lint`)
- [ ] Manual testing done
- [ ] CHANGELOG.md updated
- [ ] Version bumped in:
  - [ ] `cmd/gitflow-tui/main.go` (if hardcoded anywhere)
  - [ ] `editors/vscode/package.json`

## Creating Release

1. **Create and push tag:**
   ```bash
   git tag -a v1.0.0 -m "Release v1.0.0"
   git push origin v1.0.0
   ```

2. **GitHub Actions will:**
   - Run tests
   - Build for all platforms
   - Create GitHub Release with binaries
   - Build and push Docker image

3. **Verify release:**
   - [ ] All binaries attached to release
   - [ ] Docker image pushed
   - [ ] Release notes generated

## Post-Release

### Homebrew (macOS/Linux)

1. Update formula in `homebrew-tap` repo:
   ```bash
   # Download release and calculate SHA256
   curl -sL https://github.com/gitflow/tui/releases/download/v1.0.0/gitflow-tui-darwin-amd64.tar.gz | shasum -a 256
   
   # Update formula with new version and SHA256
   ```

2. Test installation:
   ```bash
   brew tap gitflow/tap
   brew install gitflow-tui
   gitflow-tui --version
   ```

### Scoop (Windows)

1. Update manifest in `scoop-bucket` repo
2. Test installation:
   ```powershell
   scoop bucket add gitflow https://github.com/gitflow/scoop-bucket
   scoop install gitflow-tui
   gitflow-tui --version
   ```

### Documentation

- [ ] README.md updated with new features
- [ ] Installation docs updated
- [ ] Quickstart guide updated

### Announce

- [ ] GitHub Discussions
- [ ] Twitter/X
- [ ] Reddit (r/golang, r/commandline)
- [ ] Hacker News (for major releases)
