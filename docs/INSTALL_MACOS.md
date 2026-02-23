# Installing GitFlow TUI on macOS

## Method 1: Homebrew (Recommended)

```bash
# Add tap
brew tap vahaj/gitflow

# Install
brew install gitflow-tui

# Update
brew upgrade gitflow-tui
```

## Method 2: Install Script

```bash
curl -sSL https://raw.githubusercontent.com/VahajUrRehman/git-flow/main/install.sh | bash
```

This installs to `/usr/local/bin/`.

## Method 3: Manual Download

1. Download for your Mac:
   - [Intel Mac (AMD64)](https://github.com/VahajUrRehman/git-flow/releases/latest/download/gitflow-tui-darwin-amd64.tar.gz)
   - [Apple Silicon (ARM64)](https://github.com/VahajUrRehman/git-flow/releases/latest/download/gitflow-tui-darwin-arm64.tar.gz)

2. Install:
   ```bash
   # For Intel Mac
   tar -xzf gitflow-tui-darwin-amd64.tar.gz
   chmod +x gitflow-tui-darwin-amd64
   sudo mv gitflow-tui-darwin-amd64 /usr/local/bin/gitflow-tui
   
   # For Apple Silicon
   tar -xzf gitflow-tui-darwin-arm64.tar.gz
   chmod +x gitflow-tui-darwin-arm64
   sudo mv gitflow-tui-darwin-arm64 /usr/local/bin/gitflow-tui
   ```

## Method 4: Go Install

```bash
go install github.com/VahajUrRehman/git-flow/cmd/gitflow-tui@latest
```

Make sure `~/go/bin` is in your PATH.

## Method 5: Build from Source

Requirements:
- Go 1.21+
- Xcode Command Line Tools: `xcode-select --install`

```bash
# Clone
git clone https://github.com/VahajUrRehman/git-flow.git
cd git-flow

# Build
go build -o gitflow-tui ./cmd/gitflow-tui

# Install
sudo cp gitflow-tui /usr/local/bin/
```

## Verification

```bash
gitflow-tui --version
```

## Usage

```bash
# Navigate to git repo
cd ~/Projects/my-project

# Launch
gitflow-tui
```

## Recommended Terminals

1. **iTerm2** - Best terminal for macOS
2. **Terminal.app** - Built-in, works well
3. **VSCode Terminal** - Good integration
4. **Hyper** - Modern terminal

All support 256 colors and Unicode perfectly.

## Troubleshooting

### "command not found"
```bash
# Check if /usr/local/bin is in PATH
echo $PATH | grep /usr/local/bin

# If not, add to ~/.zshrc or ~/.bash_profile
export PATH="/usr/local/bin:$PATH"
```

### Permission denied
```bash
chmod +x /usr/local/bin/gitflow-tui
```

### Colors not showing properly
Check terminal supports 256 colors:
```bash
echo $TERM
# Should show: xterm-256color or similar
```

## Uninstall

```bash
# Homebrew
brew uninstall gitflow-tui
brew untap vahaj/gitflow

# Manual
sudo rm /usr/local/bin/gitflow-tui
rm -rf ~/.config/gitflow-tui
```
