Write-Host "Building Lab 18 binaries..."

go mod tidy

go build -o camera_feed.exe ./cmd/camera_feed
go build -o analytics_engine.exe ./cmd/analytics_engine
Write-Host "Build complete."
