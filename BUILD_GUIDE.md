# GitFlow TUI - Build Guide

This guide explains how to build GitFlow TUI for all platforms.

## Quick Start

### Windows (PowerShell)

```powershell
# Build for all platforms
.\scripts\build-all.ps1 -Version 0.1.0

# Build for WSL only
.\scripts\build-all.ps1 -Version 0.1.0 -WSLOnly
```

### Cross-Platform (Make)

```bash
# Build everything (requires Make)
make all

# Build for specific platforms
make build-windows
make build-macos
make build-linux
make build-wsl
```

### Manual Build

```bash
# Windows
GOOS=windows GOARCH=amd64 go build -o dist/windows/gitflow-tui-windows-amd64.exe ./cmd/gitflow-tui

# macOS Intel
GOOS=darwin GOARCH=amd64 go build -o dist/macos/gitflow-tui-darwin-amd64 ./cmd/gitflow-tui

# macOS ARM (M1/M2)
GOOS=darwin GOARCH=arm64 go build -o dist/macos/gitflow-tui-darwin-arm64 ./cmd/gitflow-tui

# Linux
GOOS=linux GOARCH=amd64 go build -o dist/linux/gitflow-tui-linux-amd64 ./cmd/gitflow-tui

# Linux ARM
GOOS=linux GOARCH=arm64 go build -o dist/linux/gitflow-tui-linux-arm64 ./cmd/gitflow-tui
```

## Build Output

After building, you'll have:

```
dist/
├── windows/
│   ├── gitflow-tui-windows-amd64.exe
│   └── gitflow-tui-windows-arm64.exe
├── macos/
│   ├── gitflow-tui-darwin-amd64
│   └── gitflow-tui-darwin-arm64
└── linux/
    ├── gitflow-tui-linux-amd64
    └── gitflow-tui-linux-arm64

release/
├── gitflow-tui-0.1.0-windows-amd64.zip
├── gitflow-tui-0.1.0-darwin-amd64.tar.gz
├── gitflow-tui-0.1.0-linux-amd64.tar.gz
├── homebrew-formula.rb
├── scoop-manifest.json
└── install.sh
```

## Platform-Specific Instructions

### Windows

1. **Build:**
   ```powershell
   go build -o gitflow-tui.exe ./cmd/gitflow-tui
   ```

2. **Run:**
   ```powershell
   .\gitflow-tui.exe
   ```

3. **Install via Scoop (after release):**
   ```powershell
   scoop bucket add gitflow https://github.com/VahajUrRehman/gitflow-tui
   scoop install gitflow-tui
   ```

### macOS

1. **Build:**
   ```bash
   GOOS=darwin GOARCH=arm64 go build -o gitflow-tui ./cmd/gitflow-tui
   ```

2. **Install:**
   ```bash
   chmod +x gitflow-tui
   sudo mv gitflow-tui /usr/local/bin/
   ```

3. **Install via Homebrew (after release):**
   ```bash
   brew tap VahajUrRehman/gitflow
   brew install gitflow-tui
   ```

### Linux / WSL

1. **Build:**
   ```bash
   GOOS=linux GOARCH=amd64 go build -o gitflow-tui ./cmd/gitflow-tui
   ```

2. **Install:**
   ```bash
   chmod +x gitflow-tui
   sudo mv gitflow-tui /usr/local/bin/
   ```

3. **One-liner install:**
   ```bash
   curl -sSL https://raw.githubusercontent.com/VahajUrRehman/gitflow-tui/main/install.sh | bash
   ```

## Build Tags

| Platform | Architecture | Output Filename                 |
| -------- | ------------ | ------------------------------- |
| Windows  | amd64        | `gitflow-tui-windows-amd64.exe` |
| Windows  | arm64        | `gitflow-tui-windows-arm64.exe` |
| macOS    | amd64        | `gitflow-tui-darwin-amd64`      |
| macOS    | arm64        | `gitflow-tui-darwin-arm64`      |
| Linux    | amd64        | `gitflow-tui-linux-amd64`       |
| Linux    | arm64        | `gitflow-tui-linux-arm64`       |

## Version Information

Build with version info:

```bash
LDFLAGS="-s -w -X main.Version=0.1.0 -X main.BuildTime=$(date -u '+%Y-%m-%d_%H:%M:%S') -X main.Commit=$(git rev-parse --short HEAD)"
go build -ldflags "$LDFLAGS" -o gitflow-tui ./cmd/gitflow-tui
```

Check version:

```bash
gitflow-tui --version
# Output: GitFlow TUI 0.1.0 (commit: abc123, built: 2026-02-23_22:00:00)
```

## Creating a Release

1. **Update version** in code if needed
2. **Run full build:**
   ```powershell
   .\scripts\build-all.ps1 -Version 0.1.0
   ```
3. **Test binaries**
4. **Upload packages** from `release/` folder
5. **Update Homebrew formula** with SHA256 hashes
6. **Update Scoop manifest** with new URLs and hashes

## Troubleshooting

### Build fails on Windows

- Install Go from https://golang.org/dl/
- Ensure `go` is in your PATH

### Build fails for macOS on Windows

- Go supports cross-compilation, no macOS needed
- Some features may not work when cross-compiled

### WSL build issues

- Build Linux binary on Windows, then copy to WSL
- Or run the build inside WSL directly

## Makefile Targets

| Target               | Description                             |
| -------------------- | --------------------------------------- |
| `make all`           | Full build + packages for all platforms |
| `make build`         | Build for current platform              |
| `make build-all`     | Build all platforms                     |
| `make build-windows` | Build Windows binaries                  |
| `make build-macos`   | Build macOS binaries                    |
| `make build-linux`   | Build Linux binaries                    |
| `make build-wsl`     | Build for WSL                           |
| `make package-all`   | Create all packages                     |
| `make release`       | Create complete release                 |
| `make clean`         | Clean build artifacts                   |
| `make test`          | Run tests                               |
| `make install`       | Install locally                         |
