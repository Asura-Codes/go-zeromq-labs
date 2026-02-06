$ErrorActionPreference = "Stop"

go mod tidy

Write-Host "Building Lab 15 binaries..."
go build -o mr_master.exe ./cmd/mr_master
go build -o map_worker.exe ./cmd/map_worker
go build -o reduce_worker.exe ./cmd/reduce_worker

Write-Host "Build complete."
