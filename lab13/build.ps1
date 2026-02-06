$ErrorActionPreference = "Stop"

go mod tidy

Write-Host "Building Lab 13 binaries..."
go build -o titanic_broker.exe ./cmd/titanic_broker
go build -o storage_service.exe ./cmd/storage_service
go build -o patient_client.exe ./cmd/patient_client
go build -o mock_majordomo.exe ./cmd/mock_majordomo

# We also need an echo worker that speaks MDP for the Titanic Broker to call
# But we can just use a simple REP socket to mock the Majordomo Broker + Service if we want
# Or implement a minimal version.

Write-Host "Build complete."
