# build.ps1
$ErrorActionPreference = "Stop"

Write-Host "Building Lab 23 binaries..." -ForegroundColor Cyan

go mod tidy

go build -o gossip_node.exe ./cmd/gossip_node
if ($LASTEXITCODE -ne 0) { Write-Error "Build gossip_node failed"; exit 1 }

Write-Host "Build complete." -ForegroundColor Green
