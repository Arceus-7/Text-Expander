# Text Expander Auto-Start Script for Windows
# Run this as Administrator to create startup task

$scriptPath = Split-Path -Parent $MyInvocation.MyCommand.Path

# Use the VBS launcher instead of direct exe
$launcherPath = Join-Path $scriptPath "Launch-TextExpander.vbs"

# Fallback to exe if VBS doesn't exist
if (-not (Test-Path $launcherPath)) {
    $launcherPath = Join-Path $scriptPath "main.exe"
    if (-not (Test-Path $launcherPath)) {
        $launcherPath = Join-Path $scriptPath "TextExpander.exe"
    }
}

$startupFolder = [Environment]::GetFolderPath('Startup')
$shortcutPath = Join-Path $startupFolder "Text Expander.lnk"

Write-Host "Creating startup shortcut..." -ForegroundColor Cyan

# Create a WScript Shell object
$WScriptShell = New-Object -ComObject WScript.Shell

# Create the shortcut
$Shortcut = $WScriptShell.CreateShortcut($shortcutPath)
$Shortcut.TargetPath = $launcherPath
$Shortcut.WorkingDirectory = $scriptPath
$Shortcut.Description = "Text Expander - Auto Text Replacement Tool"
$Shortcut.WindowStyle = 7  # Minimized window
$Shortcut.Save()

Write-Host ""
Write-Host "✓ Success! Text Expander will now start automatically when you log in." -ForegroundColor Green
Write-Host ""
Write-Host "Shortcut created at:" -ForegroundColor Yellow
Write-Host "  $shortcutPath"
Write-Host ""
Write-Host "To disable auto-start, simply delete this shortcut from:" -ForegroundColor Cyan
Write-Host "  $startupFolder"
Write-Host ""

# Ask if user wants to start it now
$response = Read-Host "Would you like to start Text Expander now? (Y/N)"
if ($response -eq 'Y' -or $response -eq 'y') {
    if ($launcherPath -match "\.vbs$") {
        Start-Process -FilePath "wscript.exe" -ArgumentList """$launcherPath""" -WindowStyle Hidden
    }
    else {
        Start-Process -FilePath $launcherPath -WorkingDirectory $scriptPath -WindowStyle Hidden
    }
    Write-Host "✓ Text Expander is now running!" -ForegroundColor Green
}
