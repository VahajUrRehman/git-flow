# GitFlow TUI - Universal Build Script for Windows
# Builds for all platforms: Windows, macOS, Linux

param(
    [string]$Version = "dev",
    [switch]$SkipPackages,
    [switch]$WSLOnly
)

$ErrorActionPreference = "Stop"

# Colors
function Write-Color($Text, $Color) {
    Write-Host $Text -ForegroundColor $Color
}

Write-Color "========================================" "Cyan"
Write-Color "  GitFlow TUI - Universal Build System" "Cyan"
Write-Color "========================================" "Cyan"
Write-Host ""

# Build info
$BuildTime = Get-Date -Format "yyyy-MM-dd_HH:mm:ss"
try { $Commit = git rev-parse --short HEAD 2>$null } catch { $Commit = "unknown" }
if (-not $Commit) { $Commit = "unknown" }
$Ldflags = "-s -w -X main.Version=$Version -X main.BuildTime=$BuildTime -X main.Commit=$Commit"

Write-Color "Build Configuration:" "Yellow"
Write-Host "  Version:    $Version"
Write-Host "  Commit:     $Commit"
Write-Host "  Build Time: $BuildTime"
Write-Host ""

# Directories
$BuildDir = "build"
$DistDir = "dist"
$ReleaseDir = "release"

# Create directories
New-Item -ItemType Directory -Force -Path $BuildDir | Out-Null
New-Item -ItemType Directory -Force -Path $DistDir | Out-Null
New-Item -ItemType Directory -Force -Path $ReleaseDir | Out-Null
New-Item -ItemType Directory -Force -Path "$DistDir\windows" | Out-Null
New-Item -ItemType Directory -Force -Path "$DistDir\macos" | Out-Null
New-Item -ItemType Directory -Force -Path "$DistDir\linux" | Out-Null

function Build-Platform {
    param($GOOS, $GOARCH, $OutputDir, $BinaryName)
    
    $env:GOOS = $GOOS
    $env:GOARCH = $GOARCH
    
    $Output = "$OutputDir\gitflow-tui-${GOOS}-${GOARCH}"
    if ($GOOS -eq "windows") {
        $Output += ".exe"
    }
    
    Write-Color "Building for $GOOS/$GOARCH..." "Blue"
    
    $ErrorActionPreference = "Continue"
    go build -ldflags "$Ldflags" -o "$Output" .\cmd\gitflow-tui 2>&1 | Out-Null
    $ErrorActionPreference = "Stop"
    
    if ($LASTEXITCODE -eq 0) {
        Write-Color "  Built: $Output" "Green"
        return $true
    } else {
        Write-Color "  Failed: $GOOS/$GOARCH" "Red"
        return $false
    }
}

# Build all platforms
$Platforms = @(
    @{OS="windows"; Arch="amd64"; Dir="$DistDir\windows"},
    @{OS="windows"; Arch="arm64"; Dir="$DistDir\windows"},
    @{OS="darwin"; Arch="amd64"; Dir="$DistDir\macos"},
    @{OS="darwin"; Arch="arm64"; Dir="$DistDir\macos"},
    @{OS="linux"; Arch="amd64"; Dir="$DistDir\linux"},
    @{OS="linux"; Arch="arm64"; Dir="$DistDir\linux"}
)

if ($WSLOnly) {
    $Platforms = @(@{OS="linux"; Arch="amd64"; Dir="$DistDir\linux"})
    Write-Color "WSL-only build mode" "Yellow"
}

$SuccessCount = 0
$FailCount = 0

foreach ($Platform in $Platforms) {
    $Result = Build-Platform -GOOS $Platform.OS -GOARCH $Platform.Arch -OutputDir $Platform.Dir
    if ($Result) {
        $SuccessCount++
    } else {
        $FailCount++
    }
}

Write-Host ""
Write-Color "Build Summary:" "Yellow"
Write-Color "  Success: $SuccessCount" "Green"
if ($FailCount -gt 0) {
    Write-Color "  Failed:  $FailCount" "Red"
}

# Create packages if not skipped
if (-not $SkipPackages) {
    Write-Host ""
    Write-Color "Creating packages..." "Blue"
    
    # Windows packages (ZIP)
    if (Test-Path "$DistDir\windows") {
        Write-Color "  Packaging Windows..." "Blue"
        Compress-Archive -Path "$DistDir\windows\gitflow-tui-windows-amd64.exe" -DestinationPath "$ReleaseDir\gitflow-tui-$Version-windows-amd64.zip" -Force
        Compress-Archive -Path "$DistDir\windows\gitflow-tui-windows-arm64.exe" -DestinationPath "$ReleaseDir\gitflow-tui-$Version-windows-arm64.zip" -Force
        Write-Color "  ✓ Windows packages created" "Green"
    }
    
    # macOS packages (tar.gz)
    if (Test-Path "$DistDir\macos") {
        Write-Color "  Packaging macOS..." "Blue"
        tar -czf "$ReleaseDir\gitflow-tui-$Version-darwin-amd64.tar.gz" -C "$DistDir\macos" "gitflow-tui-darwin-amd64"
        tar -czf "$ReleaseDir\gitflow-tui-$Version-darwin-arm64.tar.gz" -C "$DistDir\macos" "gitflow-tui-darwin-arm64"
        Write-Color "  ✓ macOS packages created" "Green"
    }
    
    # Linux packages (tar.gz)
    if (Test-Path "$DistDir\linux") {
        Write-Color "  Packaging Linux..." "Blue"
        tar -czf "$ReleaseDir\gitflow-tui-$Version-linux-amd64.tar.gz" -C "$DistDir\linux" "gitflow-tui-linux-amd64"
        tar -czf "$ReleaseDir\gitflow-tui-$Version-linux-arm64.tar.gz" -C "$DistDir\linux" "gitflow-tui-linux-arm64"
        Write-Color "  ✓ Linux packages created" "Green"
    }
    
    # Generate checksums
    Write-Host ""
    Write-Color "Generating checksums..." "Blue"
    $Checksums = @()
    Get-ChildItem $ReleaseDir | ForEach-Object {
        $Hash = (Get-FileHash $_.FullName -Algorithm SHA256).Hash
        $Checksums += "$Hash  $($_.Name)"
    }
    $Checksums | Out-File "$ReleaseDir\checksums.txt"
    Write-Color "  ✓ Checksums saved" "Green"
    
    # Copy install scripts
    Write-Host ""
    Write-Color "Copying install scripts..." "Blue"
    Copy-Item "scripts\homebrew-formula.rb" "$ReleaseDir\" -ErrorAction SilentlyContinue
    Copy-Item "scripts\scoop-manifest.json" "$ReleaseDir\" -ErrorAction SilentlyContinue
    Copy-Item "install.sh" "$ReleaseDir\" -ErrorAction SilentlyContinue
    Write-Color "  ✓ Scripts copied" "Green"
}

Write-Host ""
Write-Color "========================================" "Green"
Write-Color "  Build Complete!" "Green"
Write-Color "========================================" "Green"
Write-Host ""
Write-Host "Output locations:"
Write-Host "  Binaries: $DistDir\"
Write-Host "  Packages: $ReleaseDir\"
Write-Host ""

if ($WSLOnly) {
    Write-Color "To install in WSL:" "Yellow"
    Write-Host "  cp $DistDir\linux\gitflow-tui-linux-amd64 /mnt/c/temp/gitflow-tui"
    Write-Host "  wsl -e sudo cp /mnt/c/temp/gitflow-tui /usr/local/bin/"
    Write-Host "  wsl -e sudo chmod +x /usr/local/bin/gitflow-tui"
}
