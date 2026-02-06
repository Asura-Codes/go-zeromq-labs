# Build all binaries

go mod tidy

go build -o trace_collector.exe ./cmd/trace_collector
go build -o monitored_service.exe ./cmd/monitored_service
