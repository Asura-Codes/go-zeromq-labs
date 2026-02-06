# build.ps1
$ErrorActionPreference = "Stop"

Write-Host "Building Lab 22 Hub & Swarm..." -ForegroundColor Cyan

go mod tidy

go build -o monitored_broker.exe ./cmd/monitored_broker
if ($LASTEXITCODE -ne 0) { Write-Error "Build monitored_broker failed"; exit 1 }

go build -o client_app.exe ./cmd/client_app
if ($LASTEXITCODE -ne 0) { Write-Error "Build client_app failed"; exit 1 }

go build -o worker_app.exe ./cmd/worker_app
if ($LASTEXITCODE -ne 0) { Write-Error "Build worker_app failed"; exit 1 }

Write-Host "Build complete." -ForegroundColor Green