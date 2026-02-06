$ErrorActionPreference = "Stop"

# Proactively kill any lingering processes from previous runs
$processes = "storage_service", "mock_majordomo", "titanic_broker", "patient_client"
foreach ($p in $processes) {
    Get-Process $p -ErrorAction SilentlyContinue | Stop-Process -Force -ErrorAction SilentlyContinue
}

# Clean up old data
if (Test-Path "titanic_data") {
    Remove-Item -Recurse -Force "titanic_data"
}

Write-Host "Starting Lab 13 (Titanic Pattern)..."
go mod tidy

# Build
./build.ps1

# Configuration
$titanicFrontend = "tcp://127.0.0.1:5555"
$storageAddr = "tcp://127.0.0.1:5557"
$worker1Addr = "tcp://127.0.0.1:5560"
$worker2Addr = "tcp://127.0.0.1:5561"
$workerList = "$worker1Addr,$worker2Addr"

# 1. Start Storage Service and Titanic Broker
$s = Start-Process .\storage_service.exe -ArgumentList "-addr $storageAddr" -PassThru -NoNewWindow
$t = Start-Process .\titanic_broker.exe -ArgumentList "-frontend $titanicFrontend -storage $storageAddr -workers $workerList" -PassThru -NoNewWindow
Start-Sleep -Seconds 2

# 2. Run Patient Client to submit requests while NO workers are running
Write-Host "--- Submitting requests (Workers are OFFLINE) ---"
$c = Start-Process .\patient_client.exe -ArgumentList "-titanic $titanicFrontend" -PassThru -NoNewWindow
Start-Sleep -Seconds 5
Write-Host "Requests are safely stored in 'titanic_data/queue' even though no one is processing them."

# 3. Start multiple workers to demonstrate redundancy and load balancing
Write-Host "--- Starting Workers ---"
$m1 = Start-Process .\mock_majordomo.exe -ArgumentList "-addr $worker1Addr -name Worker-Alpha" -PassThru -NoNewWindow
$m2 = Start-Process .\mock_majordomo.exe -ArgumentList "-addr $worker2Addr -name Worker-Beta" -PassThru -NoNewWindow

# Wait for client to finish
$c | Wait-Process
Start-Sleep -Seconds 2

# Cleanup
Write-Host "Cleaning up..."
Stop-Process -Id $t.Id -Force -ErrorAction SilentlyContinue
Stop-Process -Id $m1.Id -Force -ErrorAction SilentlyContinue
Stop-Process -Id $m2.Id -Force -ErrorAction SilentlyContinue
Stop-Process -Id $s.Id -Force -ErrorAction SilentlyContinue

Write-Host "Done."