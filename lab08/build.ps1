$ErrorActionPreference = "Stop"

go mod tidy

Write-Host "Building Lab 08 binaries..."
go build -o central_receiver.exe ./cmd/central_receiver
go build -o field_device.exe ./cmd/field_device
Write-Host "Build complete."
