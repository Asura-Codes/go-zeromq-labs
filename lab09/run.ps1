# Lab 09 Orchestration Script
$ErrorActionPreference = "Stop"

Write-Host "Initializing Lab 09..."
go mod tidy

# Build
./build.ps1

Write-Host "Starting Encrypted C2 Server and Agent..."

$pServer = Start-Process ./c2_server.exe -PassThru -NoNewWindow
Start-Sleep -Seconds 2
$pAgent = Start-Process ./secure_agent.exe -PassThru -NoNewWindow

Write-Host "Running for 10 seconds..."
Start-Sleep -Seconds 10

Stop-Process $pAgent.Id -Force
Stop-Process $pServer.Id -Force
Write-Host "Done."
