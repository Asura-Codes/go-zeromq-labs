Write-Host "Building Lab 16 binaries..." -ForegroundColor Cyan

go mod tidy

go build -o telemetry_producer.exe ./cmd/telemetry_producer
go build -o zmq_bridge.exe ./cmd/zmq_bridge

Write-Host "Build complete." -ForegroundColor Green
