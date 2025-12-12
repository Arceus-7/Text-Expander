# Setup script to configure Go to use TDM-GCC-64 from the current directory
# Run this script before building/running the project

$projectRoot = Split-Path -Parent $PSScriptRoot
$gccPath = Join-Path $projectRoot "bin\gcc.exe"
$gppPath = Join-Path $projectRoot "bin\g++.exe"

if (-not (Test-Path $gccPath)) {
    Write-Host "ERROR: TDM-GCC-64 not found at $gccPath" -ForegroundColor Red
    Write-Host "Make sure TDM-GCC-64 is installed in the project directory." -ForegroundColor Yellow
    exit 1
}

Write-Host "Configuring Go to use TDM-GCC-64..." -ForegroundColor Cyan
Write-Host "  GCC: $gccPath" -ForegroundColor Gray
Write-Host "  G++: $gppPath" -ForegroundColor Gray

# Set Go environment variables for this session
$env:CC = $gccPath
$env:CXX = $gppPath

# Also add to PATH for this session (prepend to prioritize)
$binPath = Join-Path $projectRoot "bin"
if ($env:PATH -notlike "*$binPath*") {
    $env:PATH = "$binPath;$env:PATH"
}

Write-Host ""
Write-Host "Go compiler configured for this session" -ForegroundColor Green
Write-Host ""
Write-Host "You can now run:" -ForegroundColor Yellow
Write-Host "  go run main.go" -ForegroundColor White
Write-Host "  go build main.go" -ForegroundColor White
Write-Host ""
Write-Host "Note: This configuration is only for the current PowerShell session." -ForegroundColor Gray
Write-Host "To make it permanent, add these to your PowerShell profile or set them globally." -ForegroundColor Gray
