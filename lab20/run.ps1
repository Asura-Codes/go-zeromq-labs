# Lab 20 Orchestration

# 1. Build
.\build.ps1

# 2. Wait for collector to bind
Start-Sleep -Seconds 1

# 3. Start Monitored Services
Start-Process ".\monitored_service.exe" -ArgumentList "-name gateway-service -collector tcp://127.0.0.1:5555" -NoNewWindow
Start-Process ".\monitored_service.exe" -ArgumentList "-name auth-service -collector tcp://127.0.0.1:5555" -NoNewWindow
Start-Process ".\monitored_service.exe" -ArgumentList "-name db-service -collector tcp://127.0.0.1:5555" -NoNewWindow

Write-Host "Distributed Tracing Simulation running. Press Ctrl+C to stop."

.\trace_collector.exe -port 5555

trap {
    Write-Host "Stopping processes..."
    Stop-Process -Name monitored_service, trace_collector -ErrorAction SilentlyContinue
    exit
}

while($true) { Start-Sleep -Seconds 1 }


