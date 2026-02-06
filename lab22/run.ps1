# Lab 22 Swarm Orchestration
$ErrorActionPreference = "Stop"

Write-Host "Initializing Lab 22 (The Swarm)..."
./build.ps1

Write-Host "Starting Central Message Hub..."
$Hub = Start-Process ".\monitored_broker.exe" -NoNewWindow -PassThru

Write-Host "Launching Workers..."
$Workers = @()
for ($i = 1; $i -le 3; $i++) {
    $Workers += Start-Process ".\worker_app.exe" -ArgumentList "-id worker-$i" -NoNewWindow -PassThru
}

Write-Host "Waiting for infrastructure to stabilize..."
Start-Sleep -Seconds 2

Write-Host "Launching Client Swarm..."
$Clients = @()
for ($i = 1; $i -le 5; $i++) {
    # Each client sends 5000 requests
    $Clients += Start-Process ".\client_app.exe" -ArgumentList "-id client-$i -n 5000" -NoNewWindow -PassThru
}

Write-Host "Swarm is active. Observe the console for throughput stats."
Write-Host "Press Ctrl+C to stop all processes."

trap {
    Write-Host "`nStopping all processes..."
    Stop-Process -Id ($Hub.Id) -Force -ErrorAction SilentlyContinue
    foreach ($w in $Workers) { Stop-Process -Id ($w.Id) -Force -ErrorAction SilentlyContinue }
    foreach ($c in $Clients) { Stop-Process -Id ($c.Id) -Force -ErrorAction SilentlyContinue }
    exit
}

while($true) { Start-Sleep -Seconds 1 }