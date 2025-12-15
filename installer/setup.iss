; Text Expander - Inno Setup Installer Script
; Requires Inno Setup 6.0 or later: https://jrsoftware.org/isdl.php

#define MyAppName "Text Expander"
#define MyAppVersion "1.1.0"
#define MyAppPublisher "Text Expander"
#define MyAppURL "https://github.com/yourusername/text-expander"
#define MyAppExeName "TextExpander.exe"

[Setup]
AppId={{8F9C4B2A-1D3E-4F5C-9A2B-7E6D8C4F1A3B}
AppName={#MyAppName}
AppVersion={#MyAppVersion}
AppPublisher={#MyAppPublisher}
AppPublisherURL={#MyAppURL}
AppSupportURL={#MyAppURL}
AppUpdatesURL={#MyAppURL}
DefaultDirName={autopf}\{#MyAppName}
DefaultGroupName={#MyAppName}
AllowNoIcons=yes
LicenseFile=..\LICENSE
OutputDir=..\dist
OutputBaseFilename=TextExpander-Setup-{#MyAppVersion}
Compression=lzma
SolidCompression=yes
WizardStyle=modern
PrivilegesRequired=lowest
UninstallDisplayIcon={app}\{#MyAppExeName}

[Languages]
Name: "english"; MessagesFile: "compiler:Default.isl"

[Types]
Name: "full"; Description: "Full installation"
Name: "compact"; Description: "Compact installation"
Name: "custom"; Description: "Custom installation"; Flags: iscustom

[Components]
Name: "main"; Description: "Text Expander Application"; Types: full compact custom; Flags: fixed
Name: "autostart"; Description: "Start automatically with Windows"; Types: full
Name: "shortcuts"; Description: "Desktop and Start Menu shortcuts"; Types: full

[Files]
Source: "..\TextExpander.exe"; DestDir: "{app}"; Flags: ignoreversion; Components: main
Source: "..\Launch-TextExpander.vbs"; DestDir: "{app}"; Flags: ignoreversion; Components: main
Source: "..\config\expansions.json"; DestDir: "{app}\config"; Flags: ignoreversion; Components: main
Source: "..\README.md"; DestDir: "{app}"; Flags: ignoreversion isreadme; Components: main
Source: "..\LICENSE"; DestDir: "{app}"; Flags: ignoreversion; Components: main
Source: "..\setup-auto-start.ps1"; DestDir: "{app}"; Flags: ignoreversion; Components: main
Source: "..\app-icon.ico"; DestDir: "{app}"; Flags: ignoreversion; Components: main
Source: "..\app-icon.png"; DestDir: "{app}"; Flags: ignoreversion; Components: main

[Icons]
Name: "{group}\{#MyAppName}"; Filename: "{app}\Launch-TextExpander.vbs"; IconFilename: "{app}\app-icon.ico"; Components: shortcuts
Name: "{group}\Configure {#MyAppName}"; Filename: "{app}\config\expansions.json"
Name: "{group}\{cm:UninstallProgram,{#MyAppName}}"; Filename: "{uninstallexe}"
Name: "{autodesktop}\{#MyAppName}"; Filename: "{app}\Launch-TextExpander.vbs"; IconFilename: "{app}\app-icon.ico"; Components: shortcuts

[Run]
Filename: "{app}\Launch-TextExpander.vbs"; Description: "Launch {#MyAppName} now"; Flags: postinstall shellexec skipifsilent nowait

[Registry]
Root: HKCU; Subkey: "Software\Microsoft\Windows\CurrentVersion\Run"; ValueType: string; ValueName: "TextExpander"; ValueData: """{app}\Launch-TextExpander.vbs"""; Flags: uninsdeletevalue; Components: autostart

[Code]
procedure CurStepChanged(CurStep: TSetupStep);
begin
  if CurStep = ssPostInstall then
  begin
    // Create logs directory
    if not DirExists(ExpandConstant('{app}\logs')) then
      CreateDir(ExpandConstant('{app}\logs'));
      
    // Create config\app_settings.json directory
    if not DirExists(ExpandConstant('{app}\config')) then
      CreateDir(ExpandConstant('{app}\config'));
  end;
end;

[UninstallDelete]
Type: filesandordirs; Name: "{app}\logs"
Type: filesandordirs; Name: "{app}\config\app_settings.json"
