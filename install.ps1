# GitFlow TUI - Windows Installation Script
# Usage: irm https://raw.githubusercontent.com/VahajUrRehman/git-flow/main/install.ps1 | iex

$ErrorActionPreference = "Stop"

# Configuration
$Repo = "VahajUrRehman/git-flow"
$BinaryName = "gitflow-tui"
$InstallDir = "$env:LOCALAPPDATA\GitFlow TUI"

# Colors
function Write-Color($Text, $Color) {
    Write-Host $Text -ForegroundColor $Color
}

function Show-Banner {
    Write-Color @"

╔═══════════════════════════════════════════════════════════════╗
║                                                               ║
║   ██████╗ ██╗████████╗███████╗██╗      ██████╗ ██╗    ██╗    ║
║  ██╔════╝ ██║╚══██╔══╝██╔════╝██║     ██╔═══██╗██║    ██║    ║
║  ██║  ███╗██║   ██║   █████╗  ██║     ██║   ██║██║ █╗ ██║    ║
║  ██║   ██║██║   ██║   ██╔══╝  ██║     ██║   ██║██║███╗██║    ║
║  ╚██████╔╝██║   ██║   ██║     ███████╗╚██████╔╝╚███╔███╔╝    ║
║   ╚═════╝ ╚═╝   ╚═╝   ╚═╝     ╚══════╝ ╚═════╝  ╚══╝╚══╝     ║
║                                                               ║
║              Installing GitFlow TUI for Windows               ║
║                                                               ║
╚═══════════════════════════════════════════════════════════════╝

"@ "Cyan"
}

function Get-LatestVersion {
    Write-Color "Checking for latest version..." "Yellow"
    
    try {
        $Release = Invoke-RestMethod -Uri "https://api.github.com/repos/$Repo/releases/latest" -TimeoutSec 10
        $Version = $Release.tag_name
        Write-Color "Latest version: $Version" "Green"
        return $Version
    } catch {
        Write-Color "Warning: Could not fetch latest version. Using 'latest'..." "Yellow"
        return "latest"
    }
}

function Get-Architecture {
    if ($env:PROCESSOR_ARCHITECTURE -eq "ARM64") {
        return "arm64"
    }
    return "amd64"
}

function Download-Binary($Version, $Arch) {
    $Url = "https://github.com/$Repo/releases/download/$Version/${BinaryName}-windows-${Arch}.exe"
    $TempFile = "$env:TEMP\${BinaryName}.exe"
    
    Write-Color "Downloading from GitHub..." "Yellow"
    Write-Color "URL: $Url" "DarkGray"
    
    try {
        Invoke-WebRequest -Uri $Url -OutFile $TempFile -UseBasicParsing -TimeoutSec 120
        Write-Color "Download complete!" "Green"
        return $TempFile
    } catch {
        Write-Color "Download failed! $_" "Red"
        exit 1
    }
}

function Install-Binary($Source, $Destination) {
    Write-Color "Installing to $Destination..." "Yellow"
    
    # Create directory if needed
    if (!(Test-Path $Destination)) {
        New-Item -ItemType Directory -Path $Destination -Force | Out-Null
    }
    
    $Target = Join-Path $Destination "${BinaryName}.exe"
    
    # Remove old version if exists
    if (Test-Path $Target) {
        Remove-Item $Target -Force
    }
    
    # Move binary
    Move-Item $Source $Target -Force
    
    Write-Color "Binary installed!" "Green"
    return $Target
}

function Add-ToPath($Directory) {
    Write-Color "Adding to PATH..." "Yellow"
    
    $UserPath = [Environment]::GetEnvironmentVariable("Path", "User")
    
    if ($UserPath -notlike "*$Directory*") {
        $NewPath = $UserPath + ";" + $Directory
        [Environment]::SetEnvironmentVariable("Path", $NewPath, "User")
        Write-Color "Added to PATH!" "Green"
        Write-Color "Note: Restart your terminal for PATH changes to take effect." "Yellow"
    } else {
        Write-Color "Already in PATH!" "Green"
    }
}

function Verify-Installation {
    Write-Color "Verifying installation..." "Yellow"
    
    try {
        $Version = & "$InstallDir\${BinaryName}.exe" --version
        Write-Color "✓ Installation successful!" "Green"
        Write-Color $Version "Cyan"
    } catch {
        Write-Color "✗ Verification failed" "Red"
        exit 1
    }
}

function Show-Completion {
    Write-Color @"

╔═══════════════════════════════════════════════════════════════╗
║                    Installation Complete!                     ║
╚═══════════════════════════════════════════════════════════════╝

To get started:

  1. Open a new terminal window (for PATH to update)
  2. Navigate to a git repository
  3. Run: gitflow-tui

Keyboard Shortcuts:
  • Arrow keys / hjkl  Navigate
  • Enter              Select
  • Tab                Next tab
  • Space              Stage/Unstage
  • q                  Quit
  • ?                  Help

Documentation:
  • Quick Start: https://github.com/$Repo/blob/main/QUICKSTART.md
  • Full Guide:  https://github.com/$Repo/blob/main/README.md

Support:
  • Issues: https://github.com/$Repo/issues

"@ "Green"
}

# Main installation flow
Show-Banner

$Version = Get-LatestVersion
$Arch = Get-Architecture
Write-Color "Architecture: $Arch" "DarkGray"

$TempFile = Download-Binary $Version $Arch
$InstalledPath = Install-Binary $TempFile $InstallDir
Add-ToPath $InstallDir
Verify-Installation
Show-Completion
