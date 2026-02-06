param (
    [int]$Duration = 15
)

$ErrorActionPreference = "Stop"

Write-Host "Starting Lab 12 (Majordomo Pattern) with Multiple Workers and Clients..."
go mod tidy

# Build
./build.ps1

Write-Host "Test Duration: $Duration seconds"

# Start Broker
$broker = Start-Process .\majordomo_broker.exe -PassThru -NoNewWindow
Start-Sleep -Seconds 1

# Start Multiple Workers
$workers = @()
$workers += Start-Process .\echo_worker.exe -ArgumentList "-name Alpha" -PassThru -NoNewWindow
$workers += Start-Process .\echo_worker.exe -ArgumentList "-name Beta" -PassThru -NoNewWindow
$workers += Start-Process .\echo_worker.exe -ArgumentList "-name Gamma" -PassThru -NoNewWindow

Start-Sleep -Seconds 2

# Start Multiple Clients
$clients = @()
$clients += Start-Process .\client_requester.exe -ArgumentList "-name Alice" -PassThru -NoNewWindow
$clients += Start-Process .\client_requester.exe -ArgumentList "-name Bob" -PassThru -NoNewWindow

# Wait for completion
Start-Sleep -Seconds $Duration

# Cleanup
Write-Host "Stopping all Majordomo components..."
$clients | ForEach-Object { Stop-Process -Id $_.Id -Force -ErrorAction SilentlyContinue }
$workers | ForEach-Object { Stop-Process -Id $_.Id -Force -ErrorAction SilentlyContinue }
Stop-Process -Id $broker.Id -Force -ErrorAction SilentlyContinue
Write-Host "Done."
