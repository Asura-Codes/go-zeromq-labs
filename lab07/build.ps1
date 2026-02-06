$ErrorActionPreference = "Stop"

go mod tidy

Write-Host "Building Lab 07 binaries..."
go build -o policy_master.exe ./cmd/policy_master
go build -o firewall_node.exe ./cmd/firewall_node
Write-Host "Build complete."
