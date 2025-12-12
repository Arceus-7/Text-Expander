@echo off
REM Setup script to configure Go to use TDM-GCC-64 from the current directory
REM Run this script before building/running the project

setlocal

set "PROJECT_ROOT=%~dp0.."
set "GCC_PATH=%PROJECT_ROOT%\bin\gcc.exe"
set "GPP_PATH=%PROJECT_ROOT%\bin\g++.exe"

if not exist "%GCC_PATH%" (
    echo ERROR: TDM-GCC-64 not found at %GCC_PATH%
    echo Make sure TDM-GCC-64 is installed in the project directory.
    exit /b 1
)

echo Configuring Go to use TDM-GCC-64...
echo   GCC: %GCC_PATH%
echo   G++: %GPP_PATH%

REM Set Go environment variables
set "CC=%GCC_PATH%"
set "CXX=%GPP_PATH%"

REM Add to PATH for this session
set "BIN_PATH=%PROJECT_ROOT%\bin"
set "PATH=%BIN_PATH%;%PATH%"

echo.
echo Go compiler configured for this session
echo.
echo You can now run:
echo   go run main.go
echo   go build main.go
echo.
echo Note: This configuration is only for the current CMD session.
echo To make it permanent, set these environment variables in System Properties.

endlocal

