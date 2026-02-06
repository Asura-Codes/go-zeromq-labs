# build.ps1
$ErrorActionPreference = "Stop"

Write-Host "Building Lab 24 binaries..." -ForegroundColor Cyan

go mod tidy

go build -o secure_device.exe ./cmd/secure_device
if ($LASTEXITCODE -ne 0) { Write-Error "Build secure_device failed"; exit 1 }

go build -o network_scanner.exe ./cmd/network_scanner
if ($LASTEXITCODE -ne 0) { Write-Error "Build network_scanner failed"; exit 1 }

Write-Host "Build complete." -ForegroundColor Green
