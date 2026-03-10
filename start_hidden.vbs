Set WshShell = CreateObject("WScript.Shell")
' Get the directory of this script
strPath = WScript.ScriptFullName
Set objFSO = CreateObject("Scripting.FileSystemObject")
Set objFile = objFSO.GetFile(strPath)
strFolder = objFile.ParentFolder.Path

' Run docker-compose up -d silently
WshShell.CurrentDirectory = strFolder
WshShell.Run "cmd /c docker-compose up -d", 0, False
