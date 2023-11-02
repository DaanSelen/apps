#Requires -RunAsAdministrator
cmd.exe /c "WinSW.exe install ..\CSMTA.xml"
Write-Output "Starting Service"
cmd.exe /c "sc.exe start CSMTA"
Write-Output "Done"
Pause