Write-Host "Starting Lab 17 (Distributed Consensus - PUB/SUB Mesh)..."

# 1. Build
.\build.ps1

if ($LASTEXITCODE -ne 0) {
    Write-Error "Build failed."
    exit 1
}

# 2. Start Cluster (Single Process Mode)
Write-Host "Starting Consensus Cluster..."
$cluster = Start-Process .\consensus_node.exe -ArgumentList "-id 0" -PassThru -NoNewWindow

# Wait for election
Start-Sleep -Seconds 5

# 3. Run Client
Write-Host "Sending Client Command..."
.\client_proposer.exe -cmd "SET value=100"

# Wait
Start-Sleep -Seconds 5

# 4. Cleanup
Stop-Process -Id $cluster.Id -Force
Write-Host "Cluster stopped."
