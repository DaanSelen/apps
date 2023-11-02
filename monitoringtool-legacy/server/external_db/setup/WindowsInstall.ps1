#Requires -RunAsAdministrator
cmd.exe /c "WinSW.exe install ..\CSMTS.xml"
Write-Output "Starting Service"
cmd.exe /c "sc.exe start CSMTS"
Write-Output "Done"
Pause