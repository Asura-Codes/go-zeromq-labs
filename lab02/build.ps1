# build.ps1
$ErrorActionPreference = "Stop"

Write-Host "Building Lab 02 binaries..." -ForegroundColor Cyan

go mod tidy

go build -o log_collector.exe ./cmd/log_collector
if ($LASTEXITCODE -ne 0) { Write-Error "Build log_collector failed"; exit 1 }

go build -o log_parser.exe ./cmd/log_parser
if ($LASTEXITCODE -ne 0) { Write-Error "Build log_parser failed"; exit 1 }

go build -o storage_writer.exe ./cmd/storage_writer
if ($LASTEXITCODE -ne 0) { Write-Error "Build storage_writer failed"; exit 1 }

Write-Host "Build complete." -ForegroundColor Green
