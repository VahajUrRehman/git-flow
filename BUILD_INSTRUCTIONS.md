# GitFlow TUI - Build Instructions (PowerShell)

This guide provides step-by-step instructions for building GitFlow TUI on Windows using PowerShell.

## Prerequisites

- **Go 1.21+** installed (via winget or official installer)
- **PowerShell** (Windows PowerShell or PowerShell Core)
- **Git** (for cloning the repository)

## Quick Build Steps

### Step 1: Add Go to PATH (Temporary)

If `go` command is not recognized, add it to your session PATH:

```powershell
$env:PATH = "$env:PATH;C:\Program Files\Go\bin"
```

To verify Go is accessible:
```powershell
go version
```

### Step 2: Navigate to Project Directory

```powershell
cd D:\MyOpenSource\gitflow-tui
```

### Step 3: Download Dependencies

```powershell
go mod download
```

### Step 4: Tidy Module Dependencies

This updates `go.sum` with all required checksums:

```powershell
go mod tidy
```

### Step 5: Build the Binary

```powershell
go build -o build/gitflow-tui.exe ./cmd/gitflow-tui
```

The executable will be created at `build/gitflow-tui.exe`.

---

## Multi-Platform Build

To build for different operating systems and architectures:

### Windows (AMD64)
```powershell
$env:GOOS="windows"; $env:GOARCH="amd64"; go build -o build/gitflow-tui-windows-amd64.exe ./cmd/gitflow-tui
```

### Linux (AMD64)
```powershell
$env:GOOS="linux"; $env:GOARCH="amd64"; go build -o build/gitflow-tui-linux-amd64 ./cmd/gitflow-tui
```

### macOS (AMD64 - Intel)
```powershell
$env:GOOS="darwin"; $env:GOARCH="amd64"; go build -o build/gitflow-tui-darwin-amd64 ./cmd/gitflow-tui
```

### macOS (ARM64 - Apple Silicon)
```powershell
$env:GOOS="darwin"; $env:GOARCH="arm64"; go build -o build/gitflow-tui-darwin-arm64 ./cmd/gitflow-tui
```

---

## Running the Application

After building, run the executable:

```powershell
.\build\gitflow-tui.exe
```

Or with debug mode:

```powershell
.\build\gitflow-tui.exe --debug
```

---

## Permanent PATH Setup (Optional)

To permanently add Go to your system PATH (run as Administrator):

```powershell
[Environment]::SetEnvironmentVariable("PATH", $env:PATH + ";C:\Program Files\Go\bin", "Machine")
```

Then restart PowerShell for changes to take effect.

---

## Troubleshooting

### Issue: "go: command not found"
**Solution:** Ensure Go is installed and PATH is set correctly:
```powershell
Test-Path "C:\Program Files\Go\bin\go.exe"
```

### Issue: "missing go.sum entry for module"
**Solution:** Run `go mod tidy` to update dependencies.

### Issue: Build fails with module errors
**Solution:** Clean and rebuild:
```powershell
go clean -cache
go mod tidy
go build -o build/gitflow-tui.exe ./cmd/gitflow-tui
```

### Issue: Permission denied when running executable
**Solution:** Check Windows Defender or antivirus settings, or run PowerShell as Administrator.

---

## Makefile Alternative (if Make is installed)

If you have `make` available (via Chocolatey, scoop, or Git Bash):

```powershell
# Build for current platform
make build

# Build for all platforms
make build-all

# Run tests
make test

# Clean build artifacts
make clean

# Show all available targets
make help
```

---

## Verification

After successful build, verify the binary:

```powershell
Get-Item .\build\gitflow-tui.exe | Select-Object Name, Length, LastWriteTime
.\build\gitflow-tui.exe --version
```

---

## Next Steps

- Read [README.md](README.md) for usage instructions
- Check [QUICKSTART.md](QUICKSTART.md) for getting started guide
- See [CONTRIBUTING.md](CONTRIBUTING.md) for contribution guidelines
