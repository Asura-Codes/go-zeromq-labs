$ErrorActionPreference = "Stop"

go mod tidy

Write-Host "Building Lab 05 binaries..."
go build -o audit_gateway.exe ./cmd/audit_gateway
go build -o archival_worker.exe ./cmd/archival_worker
go build -o audit_client.exe ./cmd/audit_client
Write-Host "Build complete."
