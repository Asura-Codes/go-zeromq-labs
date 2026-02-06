# Lab 04 Orchestration
$ErrorActionPreference = "Stop"

Write-Host "Initializing Lab 04..."
go mod tidy
./build.ps1

Write-Host "Starting LVC Broker..."
Start-Process ".\lvc_broker.exe" -NoNewWindow

Start-Sleep -Seconds 1
Write-Host "Starting Telemetry Source..."
Start-Process ".\telemetry_source.exe" -NoNewWindow

Write-Host "Waiting 3 seconds to populate LVC cache..."
Start-Sleep -Seconds 3

Write-Host "Starting Analyst Terminal..."
Start-Process ".\analyst_terminal.exe" -ArgumentList "sensors/temp" -NoNewWindow

Write-Host "Lab 04 running. Press Ctrl+C to stop."

trap {
    Write-Host "Stopping processes..."
    Stop-Process -Name lvc_broker, telemetry_source, analyst_terminal -ErrorAction SilentlyContinue
    exit
}

while($true) { Start-Sleep -Seconds 1 }