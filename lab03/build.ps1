# build.ps1
$ErrorActionPreference = "Stop"

Write-Host "Building Lab 03 binaries..." -ForegroundColor Cyan

go mod tidy

go build -o node_agent.exe ./cmd/node_agent
if ($LASTEXITCODE -ne 0) { Write-Error "Build node_agent failed"; exit 1 }

go build -o admin_cli.exe ./cmd/admin_cli
if ($LASTEXITCODE -ne 0) { Write-Error "Build admin_cli failed"; exit 1 }

Write-Host "Build complete." -ForegroundColor Green
