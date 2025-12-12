# How to Use Text Expander

## After Running the Application

When you run `go run main.go`, the application starts in the background. **There is no visible window** - it runs as a system tray application.

## What to Look For

1. **System Tray Icon**
   - Look in the **system tray** (notification area) in the bottom-right corner of your screen
   - On Windows, click the **^** (up arrow) to show hidden icons
   - You should see a "Text Expander" icon
   - Hover over it to see the tooltip: "Text Expander (Enabled)" or "Text Expander (Disabled)"

2. **Right-Click the Icon**
   - Right-click the system tray icon to see the menu:
     - **Enable / Disable** - Toggle expansions on/off
     - **Configure...** - Open the configuration editor
     - **Statistics** - View usage statistics
     - **View Logs** - Open the log file
     - **Reload Configuration** - Reload config from disk
     - **About** - Version information
     - **Quit** - Exit the application

## Testing the Text Expander

1. **Make sure it's enabled** (check the system tray tooltip or menu)

2. **Open any text editor** (Notepad, VS Code, browser, etc.)

3. **Type a trigger** followed by Space, Tab, or Enter:
   - `;email` + Space → expands to `your.email@example.com`
   - `;date` + Space → expands to current date (e.g., `2025-12-07`)
   - `;sig` + Space → expands to `Best regards,\nYour Name`
   - `;shrug` + Space → expands to `¯\_(ツ)_/¯`

4. **The trigger text will be deleted** and replaced with the expansion

## If You Don't See the System Tray Icon

1. **Check if the process is running:**
   ```powershell
   Get-Process | Where-Object {$_.ProcessName -eq "main"}
   ```

2. **Check for errors in the log:**
   ```powershell
   Get-Content logs\expander.log
   ```

3. **Try running it again** - sometimes the icon takes a moment to appear

4. **On Windows**, the system tray icon might be hidden:
   - Click the **^** (up arrow) in the system tray
   - Look for "Text Expander" or a generic application icon
   - You can drag it to make it always visible

## Configuration

- **Default triggers** are in `config/expansions.json`
- **Edit via GUI**: Right-click system tray icon → "Configure..."
- **Edit manually**: Edit `config/expansions.json` (hot-reloads automatically)

## Troubleshooting

- **Nothing happens when typing triggers:**
  - Check if the app is enabled (system tray tooltip)
  - Make sure you're typing the exact trigger (e.g., `;email` not `;Email`)
  - Make sure you press Space, Tab, or Enter after the trigger
  - Check `logs\expander.log` for errors

- **Can't find the system tray icon:**
  - Windows may hide it - click the **^** arrow to show hidden icons
  - The icon might be a generic application icon
  - Try restarting the application

- **App won't start:**
  - Check for errors in the terminal where you ran `go run main.go`
  - Make sure you ran `.\scripts\setup-tdm-gcc.ps1` first
  - Check `logs\expander.log` for error messages

