# Lab 02 Orchestration
$ErrorActionPreference = "Stop"

Write-Host "Initializing Lab 02..."
go mod tidy
./build.ps1

Write-Host "Starting Storage Writer..."
Start-Process ".\storage_writer.exe" -NoNewWindow

Write-Host "Starting Log Parsers (2 workers)..."
Start-Process ".\log_parser.exe" -NoNewWindow
Start-Process ".\log_parser.exe" -NoNewWindow

Start-Sleep -Seconds 1
Write-Host "Starting Log Collector..."
Start-Process ".\log_collector.exe" -NoNewWindow

Write-Host "Lab 02 running. Press Ctrl+C to stop."

trap {
    Write-Host "Stopping processes..."
    Stop-Process -Name storage_writer, log_parser, log_collector -ErrorAction SilentlyContinue
    exit
}

while($true) { Start-Sleep -Seconds 1 }
