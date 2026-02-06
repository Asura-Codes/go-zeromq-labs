# build.ps1
$ErrorActionPreference = "Stop"

Write-Host "Building Lab 04 binaries..." -ForegroundColor Cyan

go mod tidy

go build -o lvc_broker.exe ./cmd/lvc_broker
if ($LASTEXITCODE -ne 0) { Write-Error "Build lvc_broker failed"; exit 1 }

go build -o telemetry_source.exe ./cmd/telemetry_source
if ($LASTEXITCODE -ne 0) { Write-Error "Build telemetry_source failed"; exit 1 }

go build -o analyst_terminal.exe ./cmd/analyst_terminal
if ($LASTEXITCODE -ne 0) { Write-Error "Build analyst_terminal failed"; exit 1 }

Write-Host "Build complete." -ForegroundColor Green
