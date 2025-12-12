' Text Expander Launcher - Runs without terminal window
' Double-click this file to start Text Expander silently

Set objShell = CreateObject("WScript.Shell")
Set fso = CreateObject("Scripting.FileSystemObject")

' Get the directory where this script is located
scriptDir = fso.GetParentFolderName(WScript.ScriptFullName)

' Path to the executable
exePath = scriptDir & "\TextExpander.exe"

' Check if exe exists, if not try main.exe
If Not fso.FileExists(exePath) Then
    exePath = scriptDir & "\main.exe"
End If

' Launch the application completely hidden (0 = hidden window)
' vbHide = 0, doesn't wait for the program to finish
objShell.Run """" & exePath & """", 0, False

' Clean up
Set objShell = Nothing
Set fso = Nothing
