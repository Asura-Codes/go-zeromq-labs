# Build all binaries

go mod tidy

go build -o lock_server.exe ./cmd/lock_server
go build -o lock_client.exe ./cmd/lock_client
