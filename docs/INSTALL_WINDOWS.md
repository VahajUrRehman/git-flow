# Installing GitFlow TUI on Windows

## Method 1: PowerShell Install Script (Recommended)

Open **PowerShell** and run:

```powershell
irm https://raw.githubusercontent.com/VahajUrRehman/git-flow/main/install.ps1 | iex
```

This will:
1. Download the latest Windows binary
2. Install to `C:\Program Files\GitFlow TUI\`
3. Add to your system PATH

## Method 2: Scoop Package Manager

```powershell
# Install Scoop if not already installed
irm get.scoop.sh | iex

# Add bucket and install
scoop bucket add gitflow https://github.com/VahajUrRehman/git-flow
scoop install gitflow-tui

# Update
scoop update gitflow-tui
```

## Method 3: Manual Download

1. Download the latest release:
   - [Windows AMD64](https://github.com/VahajUrRehman/git-flow/releases/latest/download/gitflow-tui-windows-amd64.zip)
   - [Windows ARM64](https://github.com/VahajUrRehman/git-flow/releases/latest/download/gitflow-tui-windows-arm64.zip)

2. Extract the ZIP file

3. Add to PATH:
   ```powershell
   # Option A: Move to a folder in PATH
   Move-Item gitflow-tui-windows-amd64.exe C:\Windows\System32\gitflow-tui.exe
   
   # Option B: Add to user PATH
   [Environment]::SetEnvironmentVariable(
       "Path",
       [Environment]::GetEnvironmentVariable("Path", "User") + ";C:\Tools",
       "User"
   )
   ```

## Method 4: Build from Source

Requirements:
- [Go 1.21+](https://golang.org/dl/)
- Git

```powershell
# Clone repository
git clone https://github.com/VahajUrRehman/git-flow.git
cd git-flow

# Build
go build -o gitflow-tui.exe ./cmd/gitflow-tui

# Install
Move-Item gitflow-tui.exe C:\Windows\System32\
```

## Post-Installation

### Verify Installation

```powershell
gitflow-tui --version
```

Output:
```
GitFlow TUI 0.1.0 (commit: abc123, built: 2026-02-23)
```

### First Run

```powershell
# Navigate to a git repository
cd C:\Projects\my-project

# Launch
gitflow-tui
```

## Recommended Terminal

For the best experience, use one of these terminals:

1. **Windows Terminal** (Microsoft Store) - Recommended
2. **PowerShell 7+** - Good colors and Unicode support
3. **VSCode Terminal** - Good integration

Avoid plain `cmd.exe` - limited color and Unicode support.

## Troubleshooting

### "gitflow-tui is not recognized"
```powershell
# Check if in PATH
$env:Path -split ";" | Select-String "gitflow"

# Add to PATH manually
$env:Path += ";C:\Program Files\GitFlow TUI"
```

### "Execution policy prevents running scripts"
```powershell
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
```

### Colors look wrong
- Use Windows Terminal or PowerShell 7+
- Set terminal color mode to 256 colors or true color
- Enable UTF-8: `chcp 65001`

## Uninstall

```powershell
# If installed via script
Remove-Item -Recurse "C:\Program Files\GitFlow TUI"

# If installed via Scoop
scoop uninstall gitflow-tui

# Remove from PATH (manual edit required)
```
