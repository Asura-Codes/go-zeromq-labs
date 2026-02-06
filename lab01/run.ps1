# Lab 01 Orchestration
$ErrorActionPreference = "Stop"

Write-Host "Initializing Lab 01..."
go mod tidy
./build.ps1

Write-Host "Starting Dashboard..."
Start-Process ".\dashboard.exe" -NoNewWindow

Write-Host "Starting Monitor Agent..."
Start-Process ".\monitor_agent.exe" -NoNewWindow

Write-Host "Lab 01 running. Press Ctrl+C to stop."

trap {
    Write-Host "Stopping processes..."
    Stop-Process -Name dashboard, monitor_agent -ErrorAction SilentlyContinue
    exit
}

while($true) { Start-Sleep -Seconds 1 }
