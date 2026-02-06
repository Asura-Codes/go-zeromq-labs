# Lab 21 Orchestration (Quorum Mode)

# Build
.\build.ps1

# 1. Start 3 Independent Lock Servers
Start-Process ".\lock_server.exe" -ArgumentList "-port 5555" -NoNewWindow
Start-Process ".\lock_server.exe" -ArgumentList "-port 5556" -NoNewWindow
Start-Process ".\lock_server.exe" -ArgumentList "-port 5557" -NoNewWindow

Write-Host "Started 3 Lock Servers (Ports: 5555, 5556, 5557)."
Start-Sleep -Seconds 2

# 2. Start Competing Clients
$serverList = "tcp://127.0.0.1:5555,tcp://127.0.0.1:5556,tcp://127.0.0.1:5557"
Start-Process ".\lock_client.exe" -ArgumentList "-id ALICE -resource shared-file -servers $serverList" -NoNewWindow
Start-Process ".\lock_client.exe" -ArgumentList "-id BOB -resource shared-file -servers $serverList" -NoNewWindow

Write-Host "Alice and Bob are competing for a majority (Quorum) across the 3 servers."
Write-Host "Press Ctrl+C to stop."

trap {
    Write-Host "Stopping processes..."
    Stop-Process -Name lock_server, lock_client -ErrorAction SilentlyContinue
    exit
}

while($true) { Start-Sleep -Seconds 1 }