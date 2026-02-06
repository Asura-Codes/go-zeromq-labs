$ErrorActionPreference = "Stop"

$Root = $PSScriptRoot
if ($Root -eq $null -or $Root -eq "") { $Root = "." }

Write-Host "Starting Lab 16 (WebSocket Bridge)..."

# Build
& "$Root\build.ps1"

# 1. Start Producer
Write-Host "Starting Telemetry Producer..."
$producer = Start-Process "$Root\telemetry_producer.exe" -PassThru -NoNewWindow

# 2. Start Bridge
Write-Host "Starting ZMQ Bridge..."
Write-Host "`nDashboard available at: http://localhost:8080" -ForegroundColor Cyan

& "$Root\zmq_bridge.exe"

# Cleanup
Write-Host "Cleaning up..."
Stop-Process -Id $producer.Id -Force

Write-Host "Done."
