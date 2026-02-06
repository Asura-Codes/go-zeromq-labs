$ErrorActionPreference = "Stop"

go mod tidy

Write-Host "Building Lab 06 binaries..."
go build -o scanner_broker.exe ./cmd/scanner_broker
go build -o av_engine.exe ./cmd/av_engine
go build -o upload_service.exe ./cmd/upload_service
Write-Host "Build complete."
