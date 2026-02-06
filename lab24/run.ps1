# Lab 24 Enterprise Orchestration
$ErrorActionPreference = "Stop"

Write-Host "Initializing Lab 24 (Enterprise Discovery)..."
./build.ps1

Write-Host "Starting Device Manager (Scanner)..."
$Scanner = Start-Process ".\network_scanner.exe" -NoNewWindow -PassThru

Write-Host "Launching Multi-Device Environment..."
$Devices = @()
$BasePort = 5550

$DeviceTypes = "Camera", "Sensor", "Controller", "Gateway"

for ($i = 0; $i -lt $DeviceTypes.Length; $i++) {
    $type = $DeviceTypes[$i]
    $port = $BasePort + $i
    Write-Host "Starting $type on port $port..."
    $Devices += Start-Process ".\secure_device.exe" -ArgumentList "-type $type -service-port $port" -NoNewWindow -PassThru
}

Write-Host "`nEnvironment is active. Observe the Scanner inventory."
Write-Host "Press Ctrl+C to stop all processes."

trap {
    Write-Host "`nStopping environment..."
    Stop-Process -Id ($Scanner.Id) -Force -ErrorAction SilentlyContinue
    foreach ($d in $Devices) { Stop-Process -Id ($d.Id) -Force -ErrorAction SilentlyContinue }
    exit
}

while($true) { Start-Sleep -Seconds 1 }