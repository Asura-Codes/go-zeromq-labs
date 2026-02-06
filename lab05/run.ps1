# Lab 05 Orchestration
$ErrorActionPreference = "Stop"

Write-Host "Initializing Lab 05..."
go mod tidy
./build.ps1

Write-Host "Starting Audit Gateway..."
Start-Process ".\audit_gateway.exe" -NoNewWindow

Write-Host "Starting Archival Workers (2)..."
Start-Process ".\archival_worker.exe" -NoNewWindow
Start-Process ".\archival_worker.exe" -NoNewWindow

Start-Sleep -Seconds 2
Write-Host "Cluster active. Launching Client..."
Start-Process ".\audit_client.exe" -NoNewWindow

Write-Host "Lab 05 running. Press Ctrl+C to stop."

trap {
    Write-Host "Stopping processes..."
    Stop-Process -Name audit_gateway, archival_worker, audit_client -ErrorAction SilentlyContinue
    exit
}

while($true) { Start-Sleep -Seconds 1 }
