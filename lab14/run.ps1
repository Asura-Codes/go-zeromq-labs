$ErrorActionPreference = "Stop"

$Root = $PSScriptRoot
if ($Root -eq $null -or $Root -eq "") { $Root = "." }

Write-Host "Starting Lab 14 (DHT Ring with VNodes)..."
go mod tidy

# Build
./build.ps1

$Nodes = "tcp://127.0.0.1:5001,tcp://127.0.0.1:5002,tcp://127.0.0.1:5003"

# Start Nodes
$n1 = Start-Process "$Root\dht_node.exe" -ArgumentList "-port 5001 -peers $Nodes -vnodes 100" -PassThru -NoNewWindow
$n2 = Start-Process "$Root\dht_node.exe" -ArgumentList "-port 5002 -peers $Nodes -vnodes 100" -PassThru -NoNewWindow
$n3 = Start-Process "$Root\dht_node.exe" -ArgumentList "-port 5003 -peers $Nodes -vnodes 100" -PassThru -NoNewWindow

Start-Sleep -Seconds 2

# Run Client
Write-Host "Running DHT Client..."
& "$Root\client_put_get.exe" -port 5001

# Cleanup
Write-Host "Cleaning up..."
Stop-Process -Id $n1.Id -Force
Stop-Process -Id $n2.Id -Force
Stop-Process -Id $n3.Id -Force

Write-Host "Done."
