$ErrorActionPreference = "Stop"

go mod tidy

Write-Host "Building Lab 14 binaries..."
go build -o dht_node.exe ./cmd/dht_node
go build -o client_put_get.exe ./cmd/client_put_get

Write-Host "Build complete."
