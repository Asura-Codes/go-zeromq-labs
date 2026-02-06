$ErrorActionPreference = "Stop"

go mod tidy

$env:CGO_ENABLED = "0"
$env:GOOS = "windows"
$env:GOARCH = "amd64"
$ldflags = "-s -w"
$buildTags = "netgo"

Write-Host "Building Lab 10 binaries..."
if ( (Get-Command garble -ErrorAction SilentlyContinue) ) { # -and $false) {
    Write-Host "Using garble to build (may reduce size)..."
    & garble build -trimpath -ldflags $ldflags -o intel_provider.exe ./cmd/intel_provider
    & garble build -trimpath -ldflags $ldflags -o anomaly_detector.exe ./cmd/anomaly_detector
    & garble build -trimpath -ldflags $ldflags -o alert_logger.exe ./cmd/alert_logger
    & garble build -trimpath -ldflags $ldflags -o soc_processor.exe ./cmd/soc_processor
} else {
    go build -trimpath -tags $buildTags -ldflags $ldflags -o intel_provider.exe ./cmd/intel_provider
    go build -trimpath -tags $buildTags -ldflags $ldflags -o anomaly_detector.exe ./cmd/anomaly_detector
    go build -trimpath -tags $buildTags -ldflags $ldflags -o alert_logger.exe ./cmd/alert_logger
    go build -trimpath -tags $buildTags -ldflags $ldflags -o soc_processor.exe ./cmd/soc_processor
}

Write-Host "Build complete."
