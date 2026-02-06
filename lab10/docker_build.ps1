$ErrorActionPreference = "Stop"

Write-Host "Preparing Docker images for Lab 10 SOC Simulation..."
docker-compose build
Write-Host "Docker images are ready."
