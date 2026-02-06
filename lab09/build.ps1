$ErrorActionPreference = "Stop"

go mod tidy

Write-Host "Building Lab 09 binaries..."
go build -o c2_server.exe ./cmd/c2_server
go build -o secure_agent.exe ./cmd/secure_agent
go build -o keygen.exe ./cmd/keygen
Write-Host "Build complete."