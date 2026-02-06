# Lab 23 Gossip Cluster Orchestration
$ErrorActionPreference = "Stop"

Write-Host "Initializing Lab 23 (Decentralized Mesh)..."
./build.ps1

$Nodes = @()
$BasePort = 6660

Write-Host "Launching 5-node Gossip Cluster..."

for ($i = 1; $i -le 5; $i++) {
    $port = $BasePort + $i
    $name = "Node-$i"
    $args = "-pub tcp://*:$port -name $name"
    
    # Connect to previous node to form a chain (Gossip will propagate through it)
    if ($i -gt 1) {
        $peerPort = $port - 1
        $args += " tcp://localhost:$peerPort"
    }
    
    Write-Host "Starting $name on port $port..."
    $Nodes += Start-Process ".\gossip_node.exe" -ArgumentList $args -NoNewWindow -PassThru
    Start-Sleep -Milliseconds 500
}

Write-Host "`nCluster is active. Nodes are exchanging load metrics."
Write-Host "Observe the 'Global Cluster View' printed by each node."
Write-Host "Press Ctrl+C to stop the cluster."

trap {
    Write-Host "`nStopping cluster..."
    foreach ($n in $Nodes) { Stop-Process -Id ($n.Id) -Force -ErrorAction SilentlyContinue }
    exit
}

while($true) { Start-Sleep -Seconds 1 }