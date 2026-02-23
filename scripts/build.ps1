# GitFlow TUI Build Script for Windows
# Usage: .\scripts\build.ps1 [version]

param(
    [string]$Version = "dev"
)

$ErrorActionPreference = "Stop"

# Build info
$BuildTime = Get-Date -Format "yyyy-MM-dd_HH:mm:ss"
$Commit = (git rev-parse --short HEAD 2>$null) || "unknown"
$Ldflags = "-s -w -X main.Version=$Version -X main.BuildTime=$BuildTime -X main.Commit=$Commit"

Write-Host "Building GitFlow TUI..." -ForegroundColor Cyan
Write-Host "Version: $Version" -ForegroundColor Gray
Write-Host "Commit: $Commit" -ForegroundColor Gray
Write-Host "Build Time: $BuildTime" -ForegroundColor Gray

# Create build directory
New-Item -ItemType Directory -Force -Path "build" | Out-Null

# Build
$Output = "build\gitflow-tui.exe"
go build -ldflags "$Ldflags" -o "$Output" .\cmd\gitflow-tui

if ($LASTEXITCODE -eq 0) {
    Write-Host "✓ Built: $Output" -ForegroundColor Green
} else {
    Write-Host "✗ Build failed" -ForegroundColor Red
    exit 1
}
