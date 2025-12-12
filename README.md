# Text Expander

A lightweight text expansion tool for Windows that automatically replaces shortcuts with longer text snippets. Works across all applications.

## Features

- 120+ built-in expansions for code, symbols, and common text
- Multi-language code snippet support (Python, JavaScript, Go, C, HTML, CSS, SQL, Java, TypeScript)
- System tray integration with GUI editor
- Hot-reload configuration
- Auto-start on Windows boot
- Works in any application

## Quick Start

### Running the Application

**Option 1: Double-click launcher**
```
Launch-TextExpander.vbs
```

**Option 2: Set up auto-start**
1. Right-click `setup-auto-start.ps1`
2. Select "Run with PowerShell"
3. Confirm with 'Y'

The application runs silently in the background. Look for the icon in your system tray (bottom-right corner).

## Usage

Type a trigger followed by Space, Tab, or Enter:

| Trigger | Result |
|---------|--------|
| `;email` | your@email.com |
| `;date` | 2024-12-13 |
| `;sig` | Email signature |

### Code Snippets

**Python**
- `;pydef` - Function definition
- `;pyclass` - Class template
- `;pyboiler` - Full Python file template

**JavaScript**
- `;jsfunc` - Function declaration
- `;jsarrow` - Arrow function
- `;jsasync` - Async/await template

**C**
- `;cfunc` - Function template
- `;cmain` - Main function with includes
- `;cmalloc` - Memory allocation with error checking
- `;cboiler` - Complete C file template

**Go**
- `;gofunc` - Go function
- `;goerr` - Error handling pattern
- `;goboiler` - Go file template

**Web Development**
- `;html5` - HTML5 boilerplate
- `;cssflex` - Flexbox layout
- `;cssgrid` - CSS grid

**Database**
- `;sqlselect` - SELECT query
- `;sqlinsert` - INSERT statement
- `;sqljoin` - JOIN query

**Git**
- `;gitcommit` - Conventional commit prefix
- `;gitfix` - Fix commit prefix

Complete list available in `config/expansions.json`

## Configuration

### GUI Editor
Right-click system tray icon > Configure

### Manual Editing
Edit `config/expansions.json`:

```json
{
  "trigger": ";hello",
  "replacement": "Hello, World!",
  "case_sensitive": false,
  "description": "Greeting"
}
```

Changes apply automatically.

### Template Variables

| Variable | Result |
|----------|--------|
| `{DATE}` | Current date (YYYY-MM-DD) |
| `{TIME}` | Current time (HH:MM:SS) |
| `{DATETIME}` | Date and time |
| `{CLIPBOARD}` | Clipboard contents |
| `{CURSOR}` | Cursor position marker |

Example:
```json
{
  "trigger": ";meeting",
  "replacement": "Meeting notes for {DATE}:\n\n{CURSOR}"
}
```

## System Tray Menu

- Enable/Disable - Toggle expansions
- Configure - Open GUI editor
- Statistics - View usage stats
- View Logs - Open log file
- Reload Configuration - Refresh config
- About - Version information
- Quit - Exit application

## Project Structure

```
Text-expander-main/
├── Launch-TextExpander.vbs    # Silent launcher
├── setup-auto-start.ps1        # Auto-start configuration
├── TextExpander.exe            # Main application
├── README.md                   # Documentation
├── LICENSE                     # MIT License
├── config/
│   └── expansions.json         # Expansion definitions
├── logs/
│   └── expander.log           # Activity log
├── main.go                     # Source code
├── expander/                   # Core engine
├── gui/                        # GUI editor
└── utils/                      # Utilities
```

## Sharing the Application

### Simple Method
1. Create a folder containing:
   - `TextExpander.exe`
   - `Launch-TextExpander.vbs`
   - `config/` folder
   - `setup-auto-start.ps1`
   - `README.md`

2. Compress to ZIP

3. Share via email, cloud storage, or USB drive

### Privacy
Before sharing, review `config/expansions.json` and remove:
- Personal email addresses
- Phone numbers
- Private information

## Development

### Requirements
- Go 1.21+
- 64-bit MinGW-w64 compiler (Windows)

### Building
```bash
# Install dependencies
go mod tidy

# Build
go build -o TextExpander.exe main.go

# Run tests
go test ./...
```

### Source Structure
```
├── main.go              # Entry point, system tray
├── expander/
│   ├── expander.go      # Expansion engine
│   ├── buffer.go        # Input buffer
│   ├── keyboard.go      # Keyboard hooks
│   └── template.go      # Variable processing
├── config/
│   └── config.go        # Configuration management
├── gui/
│   └── editor.go        # GUI editor
└── utils/
    ├── logger.go        # Logging
    └── security.go      # Security checks
```

## Troubleshooting

**Expansions not working**
- Check system tray icon shows "Enabled"
- Verify trigger spelling
- Ensure Space/Tab/Enter pressed after trigger
- Review `logs/expander.log`

**System tray icon not visible**
- Click arrow (^) in system tray to show hidden icons
- Check Task Manager for running process

**Application closes when terminal closes**
- Use `Launch-TextExpander.vbs` instead of running executable directly
- Avoid `go run main.go` for regular use

## License

MIT License - See LICENSE file for details

## Expansion Categories

- Personal Information (8 expansions)
- Python (9 expansions)
- JavaScript (12 expansions)
- Go (10 expansions)
- C (12 expansions)
- HTML (8 expansions)
- CSS (5 expansions)
- SQL (5 expansions)
- Git (7 expansions)
- Symbols (15+ expansions)
- Professional Templates (6 expansions)

Total: 120+ ready-to-use expansions