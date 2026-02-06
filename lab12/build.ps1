$ErrorActionPreference = "Stop"

go mod tidy

Write-Host "Building Lab 12 binaries..."
go build -o majordomo_broker.exe ./cmd/majordomo_broker
go build -o echo_worker.exe ./cmd/echo_worker
go build -o client_requester.exe ./cmd/client_requester
Write-Host "Build complete."
