# Build Text Expander Installer
# Requires Inno Setup installed: https://jrsoftware.org/isdl.php

param(
    [switch]$SkipBuild
)

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "Text Expander Installer Build Script" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# Step 1: Build the executable
if (-not $SkipBuild) {
    Write-Host "[1/4] Building TextExpander.exe..." -ForegroundColor Yellow
    go build -ldflags="-H=windowsgui" -o TextExpander.exe main.go
    
    if ($LASTEXITCODE -ne 0) {
        Write-Host "Build failed!" -ForegroundColor Red
        exit 1
    }
    Write-Host "  Build successful!" -ForegroundColor Green
}
else {
    Write-Host "[1/4] Skipping build (using existing exe)..." -ForegroundColor Yellow
}

Write-Host ""

# Step 2: Create dist directory
Write-Host "[2/4] Creating dist directory..." -ForegroundColor Yellow
New-Item -ItemType Directory -Path "dist" -Force | Out-Null
Write-Host "  Directory created!" -ForegroundColor Green
Write-Host ""

# Step 3: Find Inno Setup compiler
Write-Host "[3/4] Looking for Inno Setup..." -ForegroundColor Yellow
$isccPaths = @(
    "C:\Program Files (x86)\Inno Setup 6\ISCC.exe",
    "C:\Program Files\Inno Setup 6\ISCC.exe",
    "$env:ProgramFiles\Inno Setup 6\ISCC.exe",
    "${env:ProgramFiles(x86)}\Inno Setup 6\ISCC.exe"
)

$iscc = $null
foreach ($path in $isccPaths) {
    if (Test-Path $path) {
        $iscc = $path
        break
    }
}

if (-not $iscc) {
    Write-Host "  Inno Setup not found!" -ForegroundColor Red
    Write-Host ""
    Write-Host "Please install Inno Setup from: https://jrsoftware.org/isdl.php" -ForegroundColor Yellow
    Write-Host "After installation, run this script again." -ForegroundColor Yellow
    exit 1
}

Write-Host "  Found: $iscc" -ForegroundColor Green
Write-Host ""

# Step 4: Compile installer
Write-Host "[4/4] Compiling installer..." -ForegroundColor Yellow
& $iscc "installer\setup.iss"

if ($LASTEXITCODE -ne 0) {
    Write-Host "  Installer compilation failed!" -ForegroundColor Red
    exit 1
}

Write-Host ""
Write-Host "========================================" -ForegroundColor Green
Write-Host "SUCCESS!" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Green
Write-Host ""
Write-Host "Installer created: dist\TextExpander-Setup-1.1.0.exe" -ForegroundColor Cyan
Write-Host ""
Write-Host "You can now distribute this installer!" -ForegroundColor White
