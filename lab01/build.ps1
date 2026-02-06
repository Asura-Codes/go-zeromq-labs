# build.ps1
$ErrorActionPreference = "Stop"

Write-Host "Building Lab 01 binaries..." -ForegroundColor Cyan

go mod tidy

# Create bin directory if it doesn't exist (optional, but good practice, here we build to root for simplicity or specific bin folder)
# keeping to root to match run.ps1 expectations for now.

go build -o monitor_agent.exe ./cmd/monitor_agent
if ($LASTEXITCODE -ne 0) { Write-Error "Build monitor_agent failed"; exit 1 }

go build -o dashboard.exe ./cmd/dashboard
if ($LASTEXITCODE -ne 0) { Write-Error "Build dashboard failed"; exit 1 }

Write-Host "Build complete." -ForegroundColor Green
