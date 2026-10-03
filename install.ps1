# ==============================================================================
# PiltiSmart 'pilti' Enterprise CLI Installer for Windows
# Usage:
#   irm https://raw.githubusercontent.com/PiltiSmart/tb-stack-ps-cli/main/install.ps1 | iex
# ==============================================================================

[CmdletBinding()]
param()

$ErrorActionPreference = "Stop"

$repo = "PiltiSmart/tb-stack-ps-cli"
$installDir = "$env:LOCALAPPDATA\Programs\pilti"
$binPath = "$installDir\pilti.exe"
$mcPath = "$installDir\mc.exe"

Write-Host "==========================================================" -ForegroundColor Cyan
Write-Host "    PiltiSmart 'pilti' Enterprise CLI Installer" -ForegroundColor Cyan
Write-Host "    Target: Windows (x64 / ARM64)" -ForegroundColor Cyan
Write-Host "==========================================================" -ForegroundColor Cyan

# 1. Architecture Check
$arch = [System.Runtime.InteropServices.RuntimeInformation]::OSArchitecture.ToString().ToLower()
$binaryName = "pilti-windows-amd64.exe"

Write-Host "[✔] Detected Windows Architecture: $arch" -ForegroundColor Green

# 2. Ensure install directory exists
if (!(Test-Path $installDir)) {
    New-Item -ItemType Directory -Path $installDir -Force | Out-Null
}

# 3. Download pilti.exe
$releaseUrl = "https://github.com/$repo/releases/latest/download/$binaryName"
Write-Host "[1/3] Downloading 'pilti.exe' from GitHub Releases..." -ForegroundColor Cyan

$downloadSuccess = $false
try {
    Invoke-WebRequest -Uri $releaseUrl -OutFile $binPath -UseBasicParsing
    $downloadSuccess = $true
} catch {
    Write-Host "[!] Pre-compiled release binary not yet reachable at $releaseUrl" -ForegroundColor Yellow
    if (Get-Command "go" -ErrorAction SilentlyContinue) {
        Write-Host "[i] Detected Go toolchain. Building from source..." -ForegroundColor Cyan
        go build -ldflags="-s -w" -o $binPath main.go
        $downloadSuccess = $true
    }
}

if (-not $downloadSuccess -or !(Test-Path $binPath)) {
    Write-Host "[ERROR] Failed to download or build 'pilti.exe'." -ForegroundColor Red
    Write-Host "Please verify your network connection or install Go: https://go.dev/dl/"
    exit 1
}

# 4. Download MinIO Client ('mc.exe')
Write-Host "[2/3] Downloading MinIO Client ('mc.exe') for S3 operations..." -ForegroundColor Cyan
$mcUrl = "https://dl.min.io/client/mc/release/windows-amd64/mc.exe"
try {
    Invoke-WebRequest -Uri $mcUrl -OutFile $mcPath -UseBasicParsing
    Write-Host "[✔] Installed 'mc.exe' to $installDir" -ForegroundColor Green
} catch {
    Write-Host "[!] Note: Could not download 'mc.exe'. 'pilti s3' will attempt download when first run." -ForegroundColor Yellow
}

# 5. Add to User PATH if not present
Write-Host "[3/3] Configuring Environment PATH..." -ForegroundColor Cyan
$userPath = [Environment]::GetEnvironmentVariable("Path", "User")
$pathParts = $userPath -split ";"
if ($pathParts -notcontains $installDir) {
    $newUserPath = if ($userPath) { "$userPath;$installDir" } else { $installDir }
    [Environment]::SetEnvironmentVariable("Path", $newUserPath, "User")
    Write-Host "[✔] Added $installDir to User PATH." -ForegroundColor Green
} else {
    Write-Host "[✔] $installDir is already in User PATH." -ForegroundColor Green
}

# Update current process PATH
if (($env:Path -split ";") -notcontains $installDir) {
    $env:Path = "$env:Path;$installDir"
}

Write-Host "`n✔ Successfully installed 'pilti' to $binPath!`n" -ForegroundColor Green
& "$binPath" version

Write-Host "`nQuick Start:" -ForegroundColor Cyan
Write-Host "  • pilti list               (View software catalog and versions)"
Write-Host "  • pilti install            (Interactive software installer)"
Write-Host "  • pilti <software-id> status"
Write-Host "  • pilti s3 ls              (MinIO / AWS S3 storage explorer)"
Write-Host "  • pilti doctor             (System health & container readiness)`n"
