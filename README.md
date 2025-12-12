# Text Expander (Go)

A cross-platform text expander written in Go. It monitors keyboard input globally and replaces trigger phrases (for example `;email`) with expanded snippets (for example `john.doe@example.com`). The application runs in the system tray and provides a Fyne‑based configuration editor.

> Target Go version: **1.21+**

## Features

- Global keyboard hook using [`gohook`](https://github.com/robotn/gohook)
- Reliable text simulation using [`robotgo`](https://github.com/go-vgo/robotgo)
- System tray integration via [`systray`](https://github.com/getlantern/systray)
- Fyne-based configuration editor (add/edit/delete snippets, custom variables, settings)
- JSON configuration file with hot-reloading (`fsnotify`)
- Template variables:
  - `{DATE}` → `YYYY-MM-DD`
  - `{TIME}` → `HH:MM:SS`
  - `{DATETIME}` → `YYYY-MM-DD HH:MM:SS`
  - `{CLIPBOARD}` → current clipboard text
  - `{CURSOR}` → where the cursor should end up after expansion
  - Custom variables such as `{NAME}`, `{COMPANY}`, etc.
- Security and safety:
  - Simple password-window detection
  - Application blacklist (common password managers)
  - Rate limiting to avoid rapid-fire expansions
- Logging & statistics:
  - Per-expansion trigger logging (no expanded contents)
  - Stats: total expansions, today, most-used trigger, last expansion time

## Project Structure

```text
text-expander/
├── main.go                 # Entry point, system tray setup
├── expander/
│   ├── expander.go         # Core expansion engine
│   ├── buffer.go           # Input buffer management
│   ├── buffer_test.go
│   ├── keyboard.go         # Keyboard hook and simulation
│   ├── template.go         # Template processing
│   └── template_test.go
├── config/
│   ├── config.go           # Configuration loading/saving/hot-reload
│   ├── config_test.go
│   └── expansions.json     # Default expansions file
├── gui/
│   ├── editor.go           # Fyne-based config editor
│   └── styles.go           # UI theme helper
├── utils/
│   ├── logger.go           # Logging & statistics
│   └── security.go         # Security heuristics (password, blacklist, rate limit)
├── go.mod
└── README.md
```

## Development Setup

```bash
# Clone and enter the project
git clone https://github.com/yourusername/text-expander.git
cd text-expander

# Initialize (if not already done) and pull dependencies
go mod tidy

# Run in development
go run ./...
```

> Note: `robotgo`, `gohook`, and `systray` all require CGO and additional native libraries on Linux/macOS/Windows. See each project's README for OS-specific prerequisites.

### Windows Setup (Important!)

Before running the project on Windows, ensure you have a **64-bit MinGW-w64** compiler installed.

#### Option 1: Using TDM-GCC-64 in Project Directory (Recommended for Development)

If you've installed TDM-GCC-64 in the project directory:

1. **Configure Go to use the local compiler:**
   ```powershell
   # PowerShell
   .\scripts\setup-tdm-gcc.ps1
   ```
   Or for CMD:
   ```cmd
   scripts\setup-tdm-gcc.bat
   ```

2. **Then run the project:**
   ```powershell
   go run main.go
   ```

   **Note:** You need to run the setup script in each new terminal session, or set the environment variables manually:
   ```powershell
   $env:CC = "D:\Text-expander-main\bin\gcc.exe"
   $env:CXX = "D:\Text-expander-main\bin\g++.exe"
   ```

#### Option 2: System-Wide Installation

1. **Check your compiler setup:**
   ```powershell
   .\scripts\check-windows-compiler.ps1
   ```

2. **If you get a "64-bit mode not compiled in" error:**
   - Install a 64-bit MinGW-w64 toolchain (see [OS Requirements](#os-requirements-high-level) below)
   - Ensure it's in your system `PATH` before any 32-bit MinGW versions
   - Restart your terminal and try again

3. **Verify installation:**
   ```powershell
   gcc --version
   ```
   Should show "x86_64" or "mingw-w64" (NOT "i686" or "mingw32").

### OS Requirements (high level)

- **Linux**
  - `gcc`
  - `libgtk-3-dev`, `libayatana-appindicator3-dev` (for systray)
  - `x11` / `libXtst` / `xcb` / `libxkbcommon` (for robotgo + gohook)
  - `xsel` or `xclip` for clipboard support

- **macOS**
  - Xcode command‑line tools
  - Accessibility and Screen Recording permissions for keyboard/mouse control

- **Windows**
  - Go toolchain (64-bit)
  - **64-bit MinGW-w64 toolchain** (required for CGO)
    - **Recommended:** [TDM-GCC-64](https://jmeubank.github.io/tdm-gcc/) or [WinLibs](https://winlibs.com/)
    - **Alternative:** [MSYS2](https://www.msys2.org/) with MinGW-w64 (`pacman -S mingw-w64-x86_64-gcc`)
    - **Important:** The old 32-bit-only MinGW.org GCC will NOT work. You need MinGW-w64.
    - Verify installation: `gcc --version` should show "x86_64" or "mingw-w64" in the output

## Running

```bash
go run main.go
```

The app will:

1. Ensure `config/expansions.json` exists (creating a default version if needed).
2. Start a global keyboard listener.
3. Place an icon/menu in the system tray.
4. Begin monitoring for triggers and performing expansions.

## System Tray

The tray menu offers:

- **Enable / Disable** – toggles expansion engine
- **Configure…** – opens the Fyne configuration editor
- **Statistics** – shows usage stats in a simple dialog (robotgo alert)
- **View Logs** – opens the log file in the default editor/viewer
- **Reload Configuration** – forces configuration reload from disk
- **About** – version and project URL
- **Quit** – exits the application

The tooltip shows the current state:

- `Text Expander (Enabled)`
- `Text Expander (Disabled)`

## Configuration

### Config File

Default path:

```text
config/expansions.json
```

Structure:

```json
{
  "expansions": [
    {
      "trigger": ";email",
      "replacement": "your.email@example.com",
      "case_sensitive": false,
      "description": "Personal email"
    },
    {
      "trigger": ";date",
      "replacement": "{DATE}",
      "case_sensitive": false,
      "description": "Current date"
    },
    {
      "trigger": ";sig",
      "replacement": "Best regards,\nYour Name",
      "case_sensitive": false,
      "description": "Email signature"
    },
    {
      "trigger": ";shrug",
      "replacement": "¯\\_(ツ)_/¯",
      "case_sensitive": false,
      "description": "Shrug emoji"
    }
  ],
  "custom_variables": {
    "NAME": "John Doe",
    "COMPANY": "Acme Corp"
  },
  "settings": {
    "enabled": true,
    "trigger_on_space": true,
    "trigger_on_tab": true,
    "trigger_on_enter": true,
    "show_notifications": false,
    "log_expansions": true
  }
}
```

Configuration changes are hot‑reloaded using `fsnotify`. Edits via the GUI are saved atomically.

### Template Variables

Supported variables in `replacement` strings:

- `{DATE}` → `YYYY-MM-DD`
- `{TIME}` → `HH:MM:SS`
- `{DATETIME}` → `YYYY-MM-DD HH:MM:SS`
- `{CLIPBOARD}` → current clipboard text at expansion time
- `{CURSOR}` → marks where the cursor should be after expansion
- Any `{KEY}` where `KEY` is present in `custom_variables` is replaced with its value.

Cursor behaviour:

- `"{CURSOR}"` is removed from the resulting text.
- The expander moves the cursor left by `cursorOffset` positions after typing to place it at the requested location.

Example:

```json
{
  "trigger": ";sig",
  "replacement": "Best regards,\n{NAME}{CURSOR}",
  "case_sensitive": false,
  "description": "Signature with cursor before name"
}
```

After expansion, the cursor will end up before `{NAME}`.

## Security & Safety

Implemented heuristics in `utils/security.go`:

- **Password field detection**
  - Checks active window title for `password`, `passcode`, etc.
- **Blacklisted apps**
  - Disables expansions when the active window title appears to belong to common password managers (e.g. 1Password, LastPass, KeePass, Bitwarden).
- **Rate limiting**
  - Simple sliding‑window limit on expansions to avoid runaway or misconfigured triggers.

These checks are used in `expander.PerformExpansion` via `utils.ShouldAllowExpansion()`.

## Logging & Statistics

`utils/logger.go`:

- Logs each expansion trigger with a timestamp (never the expanded text) to:

  ```text
  logs/expander.log
  ```

- Maintains in-memory stats:

  - `TotalExpansions`
  - `TodayExpansions`
  - `MostUsedTrigger`
  - `LastExpansion`

The tray “Statistics” menu item displays these values in a simple alert dialog.

## Testing

Unit tests cover:

- **Buffer operations** (`expander/buffer_test.go`)
- **Template processing** (`expander/template_test.go`)
- **Configuration loading/saving** (`config/config_test.go`)
- **Expansion matching logic** (`expander/expander_test.go`)

Run all tests:

```bash
go test ./...
```

## Building for Production

### Linux

```bash
go build -ldflags="-s -w" -o text-expander main.go
```

### macOS

```bash
go build -ldflags="-s -w" -o text-expander main.go
# For a full app bundle, wrap the binary as described in the systray README.
```

### Windows

```bash
go build -ldflags="-s -w -H=windowsgui" -o text-expander.exe
```

You can then package the executable with NSIS, WiX, or `go-msi`.

### Cross-compilation

```bash
# Windows from Linux/Mac
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w -H=windowsgui" -o text-expander.exe main.go

# macOS from Linux/Windows
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o text-expander main.go

# Linux from Windows/Mac
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o text-expander main.go
```

> Note: Native dependencies of `robotgo` / `gohook` / `systray` may limit cross‑compilation in practice; consult their documentation if cross‑compiling.

## Troubleshooting

### Windows Compilation Issues

- **Error: `cc1.exe: sorry, unimplemented: 64-bit mode not compiled in`**
  - **Problem:** You have a 32-bit-only MinGW compiler installed, but Go with CGO requires a 64-bit compiler.
  - **Solution:**
    1. **If TDM-GCC-64 is installed in the project directory:**
       - Run `.\scripts\setup-tdm-gcc.ps1` (PowerShell) or `scripts\setup-tdm-gcc.bat` (CMD)
       - This configures Go to use the local 64-bit compiler
       - You need to run this in each new terminal session
    2. **For system-wide installation:**
       - Uninstall the old MinGW.org GCC (32-bit only).
       - Install a 64-bit MinGW-w64 toolchain:
         - **TDM-GCC-64:** Download from https://jmeubank.github.io/tdm-gcc/ (choose the 64-bit version)
         - **WinLibs:** Download from https://winlibs.com/ (choose the MinGW-w64 standalone build)
         - **MSYS2:** Install MSYS2, then run: `pacman -S mingw-w64-x86_64-gcc`
       - Ensure the 64-bit GCC is in your system `PATH` before any 32-bit versions.
       - Verify: `gcc --version` should show "x86_64" or "mingw-w64" (not "i686" or "mingw32").
       - Restart your terminal/PowerShell and try again.
  - **Quick check:** Run `scripts\check-windows-compiler.ps1` to diagnose your compiler setup.

- **Error: `gcc: command not found`**
  - Install a 64-bit MinGW-w64 toolchain (see above) and ensure it's in your system `PATH`.

- **CGO compilation fails with other errors**
  - Ensure you have the latest version of your MinGW-w64 toolchain.
  - Try setting `CGO_ENABLED=1` explicitly: `$env:CGO_ENABLED=1; go build`
  - Check that `go env CC` points to the correct 64-bit GCC.

### Application Issues

- **No expansions occur**
  - Ensure the app has required permissions (Accessibility, Screen Recording on macOS).
  - Confirm `settings.enabled` is true in `expansions.json` or via the tray.
  - Check that you are typing the exact trigger, followed by a configured trigger key (Space/Tab/Enter).

- **High CPU or lag**
  - Keep trigger buffer small (default is 50 characters).
  - Limit the number of expansions or keep triggers distinctive.
  - Ensure no conflicting global hotkey utilities are installed.

- **Nothing happens on Linux**
  - Verify all required native libraries for `robotgo`, `gohook`, and `systray` are installed.
  - Run from terminal to inspect any logged errors.

- **Clipboard variables not working**
  - On Linux, ensure `xsel` or `xclip` is installed and on `PATH`.

## Contributing

1. Fork the repository and create a feature branch.
2. Keep code idiomatic, gofmt’d, and well‑structured.
3. Add or update tests where appropriate.
4. Submit a pull request describing the change and rationale.

## Example Usage

After building and running the app:

1. Type `;email` followed by Space.
2. The app deletes `;email` and types your configured email address.
3. `{DATE}`, `{TIME}`, and other variables in the replacement are resolved at expansion time.
4. The tray icon/menu gives quick access to enabling/disabling, configuration, stats, and quitting.