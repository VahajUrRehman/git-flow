# Installing GitFlow TUI on Linux

## Method 1: Install Script (Recommended)

```bash
curl -sSL https://raw.githubusercontent.com/VahajUrRehman/git-flow/main/install.sh | bash
```

Or with `wget`:
```bash
wget -qO- https://raw.githubusercontent.com/VahajUrRehman/git-flow/main/install.sh | bash
```

## Method 2: Package Managers

### Debian/Ubuntu (APT)

Coming soon:
```bash
# Add repository (coming soon)
# sudo add-apt-repository ppa:vahaj/gitflow
# sudo apt update
# sudo apt install gitflow-tui
```

For now, use the install script or manual method.

### Arch Linux (AUR)

Coming soon:
```bash
# yay -S gitflow-tui
# or
# paru -S gitflow-tui
```

### Homebrew on Linux

```bash
brew tap vahaj/gitflow
brew install gitflow-tui
```

## Method 3: Manual Download

1. Download:
   ```bash
   # AMD64
   wget https://github.com/VahajUrRehman/git-flow/releases/latest/download/gitflow-tui-linux-amd64.tar.gz
   
   # ARM64
   wget https://github.com/VahajUrRehman/git-flow/releases/latest/download/gitflow-tui-linux-arm64.tar.gz
   ```

2. Install:
   ```bash
   # AMD64
   tar -xzf gitflow-tui-linux-amd64.tar.gz
   chmod +x gitflow-tui-linux-amd64
   sudo mv gitflow-tui-linux-amd64 /usr/local/bin/gitflow-tui
   
   # ARM64
   tar -xzf gitflow-tui-linux-arm64.tar.gz
   chmod +x gitflow-tui-linux-arm64
   sudo mv gitflow-tui-linux-arm64 /usr/local/bin/gitflow-tui
   ```

## Method 4: Go Install

```bash
go install github.com/VahajUrRehman/git-flow/cmd/gitflow-tui@latest
```

Ensure `~/go/bin` is in your PATH:
```bash
export PATH="$PATH:$HOME/go/bin"
```

## Method 5: Build from Source

Requirements:
- Go 1.21+
- Git
- Build tools: `sudo apt install build-essential` (Debian/Ubuntu)

```bash
# Clone
git clone https://github.com/VahajUrRehman/git-flow.git
cd git-flow

# Build
make build
# or
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
cd ~/projects/my-project

# Launch
gitflow-tui
```

## WSL (Windows Subsystem for Linux)

### Option 1: Install in WSL directly

```bash
# Inside WSL
curl -sSL https://raw.githubusercontent.com/VahajUrRehman/git-flow/main/install.sh | bash
```

### Option 2: Copy from Windows

```powershell
# In Windows PowerShell, build for Linux
$env:GOOS="linux"; $env:GOARCH="amd64"
go build -o gitflow-tui-linux ./cmd/gitflow-tui

# Copy to WSL
Copy-Item gitflow-tui-linux \\wsl$\Ubuntu\home\$env:USERNAME\
```

```bash
# Inside WSL
sudo mv ~/gitflow-tui-linux /usr/local/bin/gitflow-tui
chmod +x /usr/local/bin/gitflow-tui
```

### WSL Recommended Terminals

1. **Windows Terminal** - Best overall
2. **VSCode Terminal** - Good integration
3. **Hyper** - Modern, customizable

## Troubleshooting

### "command not found"
```bash
# Check PATH
echo $PATH

# Add to PATH (add to ~/.bashrc or ~/.zshrc)
export PATH="/usr/local/bin:$PATH"
```

### Missing dependencies
```bash
# Debian/Ubuntu
sudo apt install git

# RHEL/CentOS/Fedora
sudo dnf install git

# Arch
sudo pacman -S git
```

### Terminal colors not working
```bash
# Set TERM variable
export TERM=xterm-256color

# Add to ~/.bashrc or ~/.zshrc
echo 'export TERM=xterm-256color' >> ~/.bashrc
```

### Permission denied
```bash
sudo chmod +x /usr/local/bin/gitflow-tui
```

## Desktop Entry (Optional)

Create `~/.local/share/applications/gitflow-tui.desktop`:

```ini
[Desktop Entry]
Name=GitFlow TUI
Comment=Terminal Git UI
Exec=gnome-terminal -- gitflow-tui
Type=Application
Terminal=false
Icon=utilities-terminal
Categories=Development;
```

## Uninstall

```bash
# Remove binary
sudo rm /usr/local/bin/gitflow-tui

# Remove config
rm -rf ~/.config/gitflow-tui

# If installed via Homebrew
brew uninstall gitflow-tui
```
