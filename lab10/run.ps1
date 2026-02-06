$ErrorActionPreference = "Stop"

if (Get-Command docker-compose -ErrorAction SilentlyContinue) {
    Write-Host "Starting Lab 10 SOC Simulation with Docker Compose..."
    docker-compose up --build
} else {
    Write-Host "Docker Compose not found. Starting services manually in background..."
    Start-Process go -ArgumentList "run ./cmd/intel_provider" -NoNewWindow
    Start-Process go -ArgumentList "run ./cmd/anomaly_detector" -NoNewWindow
    Start-Process go -ArgumentList "run ./cmd/alert_logger" -NoNewWindow
    Write-Host "Waiting for services to initialize..."
    Start-Sleep -Seconds 2
    go run ./cmd/soc_processor
}
