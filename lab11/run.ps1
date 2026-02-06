$ErrorActionPreference = "Stop"

Write-Host "Starting Lab 11 (High Availability)..."

Start-Process -NoNewWindow -FilePath ".\ha_broker_primary.exe"
Start-Process -NoNewWindow -FilePath ".\ha_broker_backup.exe"

Start-Sleep -Seconds 2

Write-Host "Starting Client..."
./client_app.exe
