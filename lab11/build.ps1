$ErrorActionPreference = "Stop"

go mod tidy

Write-Host "Building Lab 11 binaries..."
go build -o ha_broker_primary.exe ./cmd/ha_broker_primary
go build -o ha_broker_backup.exe ./cmd/ha_broker_backup
go build -o client_app.exe ./cmd/client_app
Write-Host "Build complete."
