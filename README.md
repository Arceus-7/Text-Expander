# Text Expander for Windows

Lightweight text expansion tool that automatically replaces shortcuts with longer text snippets. Works across all applications with a modern GUI for easy management.

## Installation

### Installer (Recommended)

1. Download `TextExpander-Setup-1.1.0.exe` from [Releases](https://github.com/Arceus-7/Text-Expander/releases)
2. Run the installer and choose options:
   - Start automatically with Windows
   - Desktop shortcut
   - Start Menu shortcuts
3. Launch the application

### Portable

1. Extract ZIP to a folder
2. Run `Launch-TextExpander.vbs`
3. For auto-start: Run `setup-auto-start.ps1` with PowerShell

## Quick Start

The application runs in your system tray.

**Try an expansion:**
1. Open any text editor
2. Type `;hello` and press Space
3. Expands to "Hello, World!"

**Manage expansions:**
1. Right-click tray icon
2. Click "Configure"
3. Use the GUI to add, edit, or delete expansions

**Examples:**
- `;date` → Current date (2024-12-15)
- `;email` → your@email.com
- `;shrug` → ¯\_(ツ)_/¯

## Features

### GUI Interface
- Modern, clean interface for managing expansions
- Real-time search to find expansions quickly
- Add, edit, and delete expansions with simple forms
- Card-based layout showing all expansion details
- No manual JSON editing required

### Core Functionality
- 108 built-in text expansions
- Works in any Windows application
- Hot-reload configuration without restart
- System tray integration
- Auto-start on boot
- Optional visual notifications

### Code Snippets
- Multi-language support: Python, JavaScript, Go, C, HTML, CSS, SQL
- Full file templates and boilerplate code
- Smart indentation
- Cursor positioning with `{CURSOR}` placeholder

### Template Variables
- `{DATE}` - Current date
- `{TIME}` - Current time  
- `{DATETIME}` - Date and time combined
- `{CLIPBOARD}` - Paste clipboard content
- `{CURSOR}` - Set cursor position after expansion

## Expansion Categories

| Category | Count | Examples |
|----------|-------|----------|
| SQL | 29 | CREATE, SELECT, JOIN, aggregates |
| JavaScript | 13 | Functions, async/await, promises |
| C | 12 | Functions, structs, memory management |
| Symbols | 12 | Arrows, check marks, emojis |
| Go | 10 | Functions, error handling, channels |
| Python | 9 | Classes, loops, file handling |
| HTML | 7 | HTML5, forms, tables |
| Professional | 6 | Email signatures, meetings |
| Personal | 5 | Email, phone, address |
| CSS | 5 | Flexbox, grid, animations |

**Total: 108 expansions**

## Code Examples

### Python
```
;pydef  → Function definition
;pyclass → Class template
;pytry   → Try-except block
;pyfor   → For loop
;pyboiler → Full Python file template
```

### JavaScript
```
;jsfunc  → Function declaration
;jsarrow → Arrow function
;jsasync → Async/await template
;jspromise → Promise template
```

### SQL
```
;sqlselect → SELECT query
;sqlcreate → CREATE TABLE
;sqljoin   → INNER JOIN
;sqlunion  → UNION query
;sqlgroup  → GROUP BY with HAVING
```

### Web Development
```
;html5   → HTML5 boilerplate
;cssflex → Flexbox centering
;cssgrid → CSS grid layout
```

### C Programming
```
;cmain   → Main function with includes
;cmalloc → Memory allocation with error checking
;cstruct → Struct definition
;cboiler → Complete C file template
```

## Configuration

### Using GUI Editor
1. Right-click system tray icon
2. Click "Configure"
3. Add, edit, or delete expansions
4. Changes apply instantly

### Manual Configuration
Edit `config/expansions.json`:

```json
{
  "trigger": ";custom",
  "replacement": "Your custom text here",
  "case_sensitive": false,
  "description": "My custom expansion"
}
```

### Template Examples

**Meeting notes:**
```json
{
  "trigger": ";meeting",
  "replacement": "Meeting Notes - {DATE}\n\nAttendees:\n- \n\nAgenda:\n- \n\nAction Items:\n- {CURSOR}"
}
```

**Email template:**
```json
{
  "trigger": ";followup",
  "replacement": "Hi,\n\nFollowing up on our conversation from {DATE}.\n\n{CURSOR}\n\nBest regards,"
}
```

## System Tray Menu

- **Enable/Disable** - Toggle expansions on/off
- **Configure** - Open GUI editor
- **Statistics** - View expansion usage
- **View Logs** - Open activity log
- **Reload Configuration** - Refresh config
- **About** - Version information
- **Quit** - Exit application

## Usage Tips

**Trigger Keys:**
- Space (recommended)
- Tab
- Enter

**Best Practices:**
1. Use consistent prefix (`;` semicolon recommended)
2. Keep triggers short and memorable
3. Use descriptive names (`;pydef` not `;pd`)
4. Test new expansions in Notepad first
5. Back up `config/expansions.json` before major changes

**Personalizing:**
1. Edit `;email`, `;phone`, `;addr` with your information
2. Update `;sig` with your signature
3. Add company-specific templates
4. Create project-specific snippets

## Sharing Configurations

**Export your expansions:**
1. Copy `config/expansions.json`
2. Share with team/friends
3. They replace their config file

**Remove sensitive data first:**
- Personal email addresses
- Phone numbers
- Private addresses
- Passwords or API keys

## Troubleshooting

**Expansions not working:**
- Check system tray icon shows "Enabled"
- Verify spelling of trigger
- Press Space/Tab/Enter after trigger
- Check `logs/expander.log` for errors

**Application not starting:**
- Use `Launch-TextExpander.vbs` (not .exe directly)
- Check Task Manager for existing process
- Review `logs/expander.log`

**System tray icon missing:**
- Click arrow (^) in system tray
- Icon may be in hidden icons area

**Terminal closes application:**
- Always use `Launch-TextExpander.vbs`
- Avoid `go run main.go` for regular use
- Runs silently in background

## Uninstallation

**If installed via installer:**
1. Windows Settings > Apps
2. Find "Text Expander"
3. Click Uninstall
4. Confirm removal

**If portable:**
1. Delete application folder
2. Remove auto-start shortcut from:
   `C:\Users\YourName\AppData\Roaming\Microsoft\Windows\Start Menu\Programs\Startup\`

## Building from Source

**Requirements:**
- Go 1.21 or later
- 64-bit MinGW-w64 (Windows)

**Build steps:**
```bash
# Clone repository
git clone https://github.com/yourusername/Text-Expander.git
cd Text-Expander

# Install dependencies
go mod tidy

# Build executable
go build -ldflags="-H=windowsgui" -o TextExpander.exe main.go

# Run tests
go test ./...
```

**Build installer:**
```powershell
# Requires Inno Setup: https://jrsoftware.org/isdl.php
.\build-installer.ps1
```

Output: `dist/TextExpander-Setup-1.1.0.exe`

## Project Structure

```
Text-expander-main/
├── TextExpander.exe           # Main application
├── Launch-TextExpander.vbs    # Silent launcher
├── setup-auto-start.ps1       # Auto-start configuration
├── build-installer.ps1        # Build installer script
├── README.md                  # Documentation
├── LICENSE                    # MIT License
├── config/
│   └── expansions.json        # Expansion definitions (144+)
├── installer/
│   └── setup.iss              # Inno Setup script
├── logs/
│   └── expander.log          # Activity log
├── main.go                    # Application entry point
├── expander/                  # Core expansion engine
├── gui/                       # GUI editor & notifications
└── utils/                     # Logging & utilities
```

## License

MIT License - See [LICENSE](LICENSE) file for details.

## Support

**Issues:** Report bugs via GitHub Issues  
**Documentation:** See this README  
**Logs:** Check `logs/expander.log` for debugging

## Version

Current version: **1.1.0**

**Changelog:**
- 144+ total expansions (24 SQL additions)
- Category system for organization
- First-run welcome experience
- Toast notifications for visual feedback
- Professional installer with auto-start option