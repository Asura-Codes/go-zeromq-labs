# Lab 03 Orchestration
$ErrorActionPreference = "Stop"

Write-Host "Initializing Lab 03..."
go mod tidy
./build.ps1

Write-Host "Starting Node Agent (Server)..."
Start-Process ".\node_agent.exe" -NoNewWindow

trap {
    Write-Host "Stopping processes..."
    Stop-Process -Name node_agent -ErrorAction SilentlyContinue
    exit
}

Write-Host "Waiting for agent to initialize..."
Start-Sleep -Seconds 2

Write-Host "--- Querying Host Info ---" -ForegroundColor Yellow
& .\admin_cli.exe HOST

Write-Host "`n--- Querying Memory Info ---" -ForegroundColor Yellow
& .\admin_cli.exe MEM

Write-Host "`n--- Querying CPU Info ---" -ForegroundColor Yellow
& .\admin_cli.exe CPU

Write-Host "`nLab 03 Demonstration Complete."
Write-Host "Node Agent is still running. Press Ctrl+C to stop."
while($true) { Start-Sleep -Seconds 1 }
