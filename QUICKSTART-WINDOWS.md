# Quick Start Guide for Windows

## If TDM-GCC-64 is Installed in Project Directory

Since TDM-GCC-64 is installed locally in this project, follow these steps:

### PowerShell (Recommended)

1. **Open PowerShell in the project directory**

2. **Run the setup script:**
   ```powershell
   .\scripts\setup-tdm-gcc.ps1
   ```

3. **Run the project:**
   ```powershell
   go run main.go
   ```

### Command Prompt (CMD)

1. **Open CMD in the project directory**

2. **Run the setup script:**
   ```cmd
   scripts\setup-tdm-gcc.bat
   ```

3. **Run the project:**
   ```cmd
   go run main.go
   ```

### Manual Setup (Alternative)

If you prefer to set environment variables manually:

**PowerShell:**
```powershell
$env:CC = "$PWD\bin\gcc.exe"
$env:CXX = "$PWD\bin\g++.exe"
go run main.go
```

**CMD:**
```cmd
set CC=%CD%\bin\gcc.exe
set CXX=%CD%\bin\g++.exe
go run main.go
```

### Important Notes

- **You must run the setup script in each new terminal session** before running `go run` or `go build`
- The setup script only affects the current terminal session
- To make it permanent, you can:
  - Add the environment variables to your PowerShell profile
  - Set them in System Environment Variables
  - Or create an alias/function in your shell profile

### Verify It's Working

After running the setup script, verify Go is using the correct compiler:

```powershell
go env CC CXX
```

Should show:
```
D:\Text-expander-main\bin\gcc.exe
D:\Text-expander-main\bin\g++.exe
```

### Building for Production

```powershell
# After running setup-tdm-gcc.ps1
go build -ldflags="-s -w -H=windowsgui" -o text-expander.exe main.go
```

