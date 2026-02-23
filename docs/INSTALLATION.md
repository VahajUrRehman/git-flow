# Installation Guide

## Table of Contents
- [Prerequisites](#prerequisites)
- [Binary Installation](#binary-installation)
- [Building from Source](#building-from-source)
- [Editor Integration](#editor-integration)
- [Verification](#verification)
- [Troubleshooting](#troubleshooting)

## Prerequisites

### Required
- **Go 1.21+** (for building from source)
- **Git 2.20+**

### Optional
- **Neovim 0.8+** (for Neovim plugin)
- **VSCode 1.74+** (for VSCode extension)

## Binary Installation

### macOS

#### Homebrew
```bash
brew tap gitflow/tui
brew install gitflow-tui
```

#### Manual
```bash
# Download latest release
curl -L -o gitflow-tui.tar.gz https://github.com/gitflow/tui/releases/latest/download/gitflow-tui-darwin-amd64.tar.gz

# Extract
tar -xzf gitflow-tui.tar.gz

# Move to PATH
sudo mv gitflow-tui /usr/local/bin/

# Make executable
chmod +x /usr/local/bin/gitflow-tui
```

### Linux

#### Debian/Ubuntu
```bash
# Download .deb package
wget https://github.com/gitflow/tui/releases/latest/download/gitflow-tui_amd64.deb

# Install
sudo dpkg -i gitflow-tui_amd64.deb
```

#### Arch Linux (AUR)
```bash
# Using yay
yay -S gitflow-tui

# Using paru
paru -S gitflow-tui
```

#### Fedora/RHEL
```bash
# Download .rpm package
wget https://github.com/gitflow/tui/releases/latest/download/gitflow-tui_x86_64.rpm

# Install
sudo rpm -i gitflow-tui_x86_64.rpm
```

#### Generic
```bash
# Download
curl -L -o gitflow-tui.tar.gz https://github.com/gitflow/tui/releases/latest/download/gitflow-tui-linux-amd64.tar.gz

# Extract
tar -xzf gitflow-tui.tar.gz

# Install
sudo install -Dm755 gitflow-tui /usr/local/bin/gitflow-tui
```

### Windows

#### Scoop
```powershell
scoop bucket add gitflow https://github.com/gitflow/scoop-bucket
scoop install gitflow-tui
```

#### Chocolatey
```powershell
choco install gitflow-tui
```

#### Manual
1. Download `gitflow-tui-windows-amd64.zip` from [releases](https://github.com/gitflow/tui/releases)
2. Extract to a folder (e.g., `C:\Program Files\gitflow-tui`)
3. Add to PATH:
   ```powershell
   [Environment]::SetEnvironmentVariable("Path", $env:Path + ";C:\Program Files\gitflow-tui", "User")
   ```

## Building from Source

### Clone Repository
```bash
git clone https://github.com/gitflow/tui.git
cd tui
```

### Install Dependencies
```bash
make deps
```

### Build
```bash
# Build for current platform
make build

# Build for all platforms
make build-all
```

### Install
```bash
# Install globally
sudo make install

# Install to custom location
make install PREFIX=/opt/gitflow
```

### Using Go Install
```bash
go install github.com/gitflow/tui/cmd/gitflow-tui@latest
```

## Editor Integration

### Neovim

#### Using lazy.nvim
```lua
{
  'gitflow/tui',
  config = function()
    require('gitflow').setup({
      keymaps = {
        open = '<leader>gg',
      },
    })
  end,
}
```

#### Using packer.nvim
```lua
use {
  'gitflow/tui',
  config = function()
    require('gitflow').setup()
  end
}
```

#### Manual Installation
```bash
# Clone plugin
git clone https://github.com/gitflow/tui.git ~/.config/nvim/pack/plugins/start/gitflow-tui

# Or using the built-in plugin
cp -r editors/nvim/lua/gitflow ~/.config/nvim/lua/
```

Add to `init.lua`:
```lua
require('gitflow').setup()
```

### VSCode

#### From Marketplace
1. Open VSCode
2. Go to Extensions (Ctrl+Shift+X)
3. Search for "GitFlow TUI"
4. Click Install

#### From VSIX
```bash
cd editors/vscode
npm install
npm run package
code --install-extension gitflow-tui-*.vsix
```

#### From Source
```bash
cd editors/vscode
npm install
npm run compile
# Press F5 to launch extension host
```

## Verification

### Check Installation
```bash
# Check binary
gitflow-tui --version

# Expected output:
# GitFlow TUI v1.0.0 - Open Source Git Management
# Supports: Neovim | VSCode | Terminal
```

### Test in Repository
```bash
# Navigate to a git repository
cd /path/to/your/repo

# Launch TUI
gitflow-tui
```

### Check Editor Integration

#### Neovim
```vim
:GitFlow
```

#### VSCode
```
Ctrl+Shift+P → "GitFlow: Open GitFlow TUI"
```

## Troubleshooting

### Binary Not Found

**Problem**: `gitflow-tui: command not found`

**Solutions**:
1. Check if binary is in PATH:
   ```bash
   which gitflow-tui
   ```

2. Add to PATH manually:
   ```bash
   # Add to ~/.bashrc or ~/.zshrc
   export PATH="$PATH:/usr/local/bin"
   ```

3. Use full path:
   ```bash
   /usr/local/bin/gitflow-tui
   ```

### Permission Denied

**Problem**: `permission denied: gitflow-tui`

**Solution**:
```bash
chmod +x /path/to/gitflow-tui
```

### Git Repository Not Found

**Problem**: `not a git repository`

**Solutions**:
1. Initialize git repository:
   ```bash
   git init
   ```

2. Navigate to correct directory:
   ```bash
   cd /path/to/git/repo
   ```

3. Specify directory:
   ```bash
   gitflow-tui --cwd /path/to/repo
   ```

### Neovim Plugin Not Loading

**Problem**: `module 'gitflow' not found`

**Solutions**:
1. Check plugin installation path
2. Ensure binary is in PATH
3. Configure binary path:
   ```lua
   require('gitflow').setup({
     bin_path = '/usr/local/bin/gitflow-tui'
   })
   ```

### VSCode Extension Not Working

**Problem**: Commands not appearing

**Solutions**:
1. Reload VSCode window
2. Check extension is enabled
3. Configure binary path in settings:
   ```json
   {
     "gitflow.binaryPath": "/usr/local/bin/gitflow-tui"
   }
   ```

### Terminal Display Issues

**Problem**: Colors not showing correctly

**Solutions**:
1. Check terminal supports 256 colors:
   ```bash
   echo $TERM
   ```

2. Set TERM variable:
   ```bash
   export TERM=xterm-256color
   ```

3. Disable colors in config:
   ```json
   {
     "theme": {
       "colors": {
         "primary": "white"
       }
     }
   }
   ```

### Mouse Not Working

**Problem**: Mouse clicks not registering

**Solutions**:
1. Enable mouse in config:
   ```json
   {
     "mouse_enabled": true
   }
   ```

2. Check terminal supports mouse:
   - iTerm2: Preferences → Profiles → Terminal → Enable mouse reporting
   - Terminal.app: Limited support

### Performance Issues

**Problem**: Slow rendering with large repositories

**Solutions**:
1. Limit commit history:
   ```json
   {
     "max_commits": 100
   }
   ```

2. Disable animations:
   ```json
   {
     "animations": false
   }
   ```

3. Use compact graph style:
   ```json
   {
     "graph_style": "compact"
   }
   ```

## Getting Help

If you continue to experience issues:

1. Check [GitHub Issues](https://github.com/gitflow/tui/issues)
2. Join our [Discord](https://discord.gg/gitflow)
3. Open a new issue with:
   - Operating system
   - Terminal emulator
   - GitFlow TUI version
   - Error messages
   - Steps to reproduce
