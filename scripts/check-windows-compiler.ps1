# Windows Compiler Check Script for Text Expander
# This script checks if you have a compatible 64-bit compiler for CGO

Write-Host "=== Text Expander - Windows Compiler Check ===" -ForegroundColor Cyan
Write-Host ""

$issues = @()
$warnings = @()

# Check if gcc is available
Write-Host "Checking for GCC compiler..." -ForegroundColor Yellow
try {
    $gccVersion = & gcc --version 2>&1 | Select-Object -First 1
    Write-Host "  Found: $gccVersion" -ForegroundColor Green
    
    # Check if it's 64-bit capable
    $gccFullOutput = & gcc --version 2>&1 | Out-String
    $gccPath = (Get-Command gcc -ErrorAction SilentlyContinue).Source
    
    if ($gccPath) {
        Write-Host "  Path: $gccPath" -ForegroundColor Gray
        
        # Check for 32-bit only indicators
        if ($gccFullOutput -match "MinGW\.org" -or $gccFullOutput -match "i686" -or $gccFullOutput -match "mingw32") {
            $issues += "You have a 32-bit-only MinGW compiler. This will NOT work with Go CGO on 64-bit Windows."
            Write-Host "  ERROR: 32-bit-only compiler detected!" -ForegroundColor Red
        }
        elseif ($gccFullOutput -match "x86_64" -or $gccFullOutput -match "mingw-w64" -or $gccFullOutput -match "w64") {
            Write-Host "  OK: 64-bit compiler detected" -ForegroundColor Green
        }
        else {
            $warnings += "Could not determine if compiler is 64-bit. Please verify manually."
            Write-Host "  WARNING: Could not verify 64-bit support" -ForegroundColor Yellow
        }
    }
} catch {
    $issues += "GCC compiler not found in PATH. Install a 64-bit MinGW-w64 toolchain."
    Write-Host "  ERROR: GCC not found!" -ForegroundColor Red
}

Write-Host ""

# Check Go CGO settings
Write-Host "Checking Go CGO configuration..." -ForegroundColor Yellow
try {
    $cgoEnabled = & go env CGO_ENABLED 2>&1
    $goCC = & go env CC 2>&1
    $goCXX = & go env CXX 2>&1
    
    Write-Host "  CGO_ENABLED: $cgoEnabled" -ForegroundColor $(if ($cgoEnabled -eq "1") { "Green" } else { "Yellow" })
    Write-Host "  CC: $goCC" -ForegroundColor Gray
    Write-Host "  CXX: $goCXX" -ForegroundColor Gray
    
    if ($cgoEnabled -ne "1") {
        $warnings += "CGO is disabled. This project requires CGO for robotgo, gohook, and systray."
    }
} catch {
    $issues += "Could not check Go environment. Is Go installed?"
    Write-Host "  ERROR: Go not found or not in PATH!" -ForegroundColor Red
}

Write-Host ""

# Test compilation (dry run)
Write-Host "Testing compiler compatibility..." -ForegroundColor Yellow
try {
    # Try to compile a simple CGO test
    $testFile = Join-Path $env:TEMP "cgo_test.go"
    @"
package main
/*
#include <stdio.h>
void hello() { printf("Hello from CGO\n"); }
*/
import "C"
import "fmt"
func main() {
    C.hello()
    fmt.Println("CGO test successful")
}
"@ | Out-File -FilePath $testFile -Encoding UTF8
    
    $testDir = Join-Path $env:TEMP "cgo_test_$(Get-Random)"
    New-Item -ItemType Directory -Path $testDir -Force | Out-Null
    Move-Item -Path $testFile -Destination (Join-Path $testDir "main.go") -Force
    
    Push-Location $testDir
    try {
        $buildOutput = & go build -o test.exe main.go 2>&1
        if ($LASTEXITCODE -eq 0) {
            Write-Host "  OK: CGO compilation test passed" -ForegroundColor Green
            Remove-Item -Path "test.exe" -ErrorAction SilentlyContinue
        } else {
            $errorMsg = $buildOutput -join "`n"
            if ($errorMsg -match "64-bit mode not compiled in") {
                $issues += "Compiler test failed: 64-bit mode not available. Install MinGW-w64 (64-bit)."
                Write-Host "  ERROR: 64-bit mode not available!" -ForegroundColor Red
            } else {
                $issues += "Compiler test failed: $errorMsg"
                Write-Host "  ERROR: Compilation test failed!" -ForegroundColor Red
            }
        }
    } finally {
        Pop-Location
        Remove-Item -Path $testDir -Recurse -Force -ErrorAction SilentlyContinue
    }
} catch {
    $warnings += "Could not run compiler test: $_"
    Write-Host "  WARNING: Could not test compilation" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "=== Summary ===" -ForegroundColor Cyan

if ($issues.Count -eq 0 -and $warnings.Count -eq 0) {
    Write-Host "✓ Your compiler setup looks good!" -ForegroundColor Green
    Write-Host ""
    Write-Host "You should be able to build the project with:" -ForegroundColor Green
    Write-Host "  go run main.go" -ForegroundColor White
} else {
    if ($issues.Count -gt 0) {
        Write-Host "✗ Issues found:" -ForegroundColor Red
        foreach ($issue in $issues) {
            Write-Host "  - $issue" -ForegroundColor Red
        }
        Write-Host ""
        Write-Host "To fix these issues:" -ForegroundColor Yellow
        Write-Host "  1. Install a 64-bit MinGW-w64 toolchain:" -ForegroundColor White
        Write-Host "     - TDM-GCC-64: https://jmeubank.github.io/tdm-gcc/" -ForegroundColor Gray
        Write-Host "     - WinLibs: https://winlibs.com/" -ForegroundColor Gray
        Write-Host "     - MSYS2: https://www.msys2.org/ (then: pacman -S mingw-w64-x86_64-gcc)" -ForegroundColor Gray
        Write-Host "  2. Ensure the 64-bit GCC is in your PATH before any 32-bit versions" -ForegroundColor White
        Write-Host "  3. Restart your terminal and run this script again" -ForegroundColor White
    }
    
    if ($warnings.Count -gt 0) {
        Write-Host ""
        Write-Host "⚠ Warnings:" -ForegroundColor Yellow
        foreach ($warning in $warnings) {
            Write-Host "  - $warning" -ForegroundColor Yellow
        }
    }
}

Write-Host ""

