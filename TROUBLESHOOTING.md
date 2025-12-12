# Troubleshooting Guide

## Triggers Not Working

If typing triggers like `;email` followed by Space doesn't expand, check the following:

### 1. Check if the App is Enabled

- Right-click the system tray icon
- Verify "Enable" is shown (not "Disable")
- If it says "Disable", the app is enabled
- If it says "Enable", click it to enable the app

### 2. Check Configuration

Verify your `config/expansions.json` has the correct settings:

```json
{
  "settings": {
    "enabled": true,
    "trigger_on_space": true,
    "trigger_on_tab": true,
    "trigger_on_enter": true
  }
}
```

### 3. Windows Permissions

On Windows, global keyboard hooks may require:
- **Administrator privileges** - Try running the app as Administrator
- **No conflicting software** - Some antivirus or security software may block global hooks

### 4. Test the Keyboard Hook

1. **Check the log file:**
   ```powershell
   Get-Content logs\expander.log
   ```
   Look for "keyboard hook started successfully"

2. **Try typing a trigger:**
   - Open Notepad or any text editor
   - Type `;email` (exactly as shown)
   - Press **Space** (not Enter)
   - It should expand to `your.email@example.com`

3. **Verify the trigger format:**
   - Triggers are case-insensitive by default (unless `case_sensitive: true`)
   - Must be typed exactly as defined (e.g., `;email` not `;Email` or `; email`)
   - Must be followed by Space, Tab, or Enter (depending on settings)

### 5. Common Issues

**Issue: Nothing happens when typing triggers**
- ✅ Check if app is enabled (system tray)
- ✅ Verify trigger is typed correctly (e.g., `;email`)
- ✅ Make sure you press Space/Tab/Enter after the trigger
- ✅ Check if you're in a password field (expansions are disabled for security)
- ✅ Try running as Administrator

**Issue: "failed to start keyboard hook" in logs**
- Run the app as Administrator
- Check if antivirus is blocking the hook
- Restart the application

**Issue: Expansions work sometimes but not always**
- Check rate limiting (max 20 expansions per 2 seconds)
- Verify you're not in a blacklisted app (password managers)
- Check if the active window is detected as a password field

### 6. Debug Steps

1. **Check if the process is running:**
   ```powershell
   Get-Process | Where-Object {$_.ProcessName -eq "main"}
   ```

2. **Check the log file:**
   ```powershell
   Get-Content logs\expander.log -Tail 50
   ```

3. **Verify configuration:**
   ```powershell
   Get-Content config\expansions.json
   ```

4. **Test with a simple trigger:**
   - Edit `config/expansions.json`
   - Add a simple test expansion:
     ```json
     {
       "trigger": "test",
       "replacement": "TEST WORKED!",
       "case_sensitive": false,
       "description": "Test expansion"
     }
     ```
   - Type `test` + Space in any text editor
   - Should expand to "TEST WORKED!"

## Fyne Configuration Editor Error

If you see: `*** Error in Fyne call thread, this should have been called in fyne.Do [AndWait] ***`

**This has been fixed!** The error was caused by Fyne operations running on the wrong thread. The fix ensures the editor runs on its own OS thread.

If you still see this error:
1. Make sure you've rebuilt the application:
   ```powershell
   .\scripts\setup-tdm-gcc.ps1
   go build main.go
   ```
2. Restart the application
3. Try opening the configuration editor again

## Still Not Working?

1. **Restart the application:**
   - Right-click system tray icon → Quit
   - Run `go run main.go` again

2. **Check for errors:**
   - Look at the terminal where you ran `go run main.go`
   - Check `logs\expander.log` for error messages

3. **Verify Windows permissions:**
   - Try running PowerShell as Administrator
   - Then run the application from there

4. **Test in a simple text editor:**
   - Use Notepad (not VS Code or other complex editors initially)
   - Type `;email` + Space
   - Should expand immediately

