# WebLauncher Windows Build Script
param(
    [string]$OutputName = "weblauncher.exe",
    [string]$OutputDir = "dist",
    [switch]$Hidden
)

$ErrorActionPreference = "Stop"

# Get script directory
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$ProjectRoot = Split-Path -Parent $ScriptDir

# Set paths
$SrcDir = Join-Path $ProjectRoot "src"
$OutputPath = Join-Path $ProjectRoot $OutputDir

# Ensure output directory exists
if (!(Test-Path $OutputPath)) {
    New-Item -ItemType Directory -Path $OutputPath | Out-Null
}

# Set ldflags
$LdfFlags = "-s -w"
if ($Hidden) {
    $LdfFlags += " -H=windowsgui"
}

# Generate resource file with rsrc (icon only)
$RsrcPath = Join-Path (go env GOPATH) "bin\rsrc.exe"
if (!(Test-Path $RsrcPath)) {
    Write-Host "Installing rsrc..."
    go install github.com/akavel/rsrc@latest
}

Write-Host "Generating resource file..."
Push-Location $SrcDir
& $RsrcPath -ico="assets/icon.ico" -o="rsrc.syso"
Pop-Location

# Build executable
Write-Host "Building $OutputName ..."
$FullOutputPath = Join-Path $OutputPath $OutputName

Set-Location $SrcDir
go build -ldflags="$LdfFlags" -o "$FullOutputPath"
Set-Location $ProjectRoot

# Cleanup syso file
$SysoPath = Join-Path $SrcDir "rsrc.syso"
if (Test-Path $SysoPath) {
    Remove-Item -Path $SysoPath -ErrorAction SilentlyContinue
}

Write-Host "Build complete: $FullOutputPath" -ForegroundColor Green
