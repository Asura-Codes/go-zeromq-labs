Write-Host "Building Lab 17 binaries..." -ForegroundColor Cyan

go mod tidy

Write-Host "Building binaries..."
go build -o consensus_node.exe ./cmd/consensus_node
go build -o client_proposer.exe ./cmd/client_proposer

Write-Host "Build complete." -ForegroundColor Green
