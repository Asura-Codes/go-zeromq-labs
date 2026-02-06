Write-Host "Starting Lab 18: Zero-Copy Video Streaming..."

.\build.ps1

$cam = Start-Process .\camera_feed.exe -PassThru -NoNewWindow
$ana = Start-Process .\analytics_engine.exe -PassThru -NoNewWindow

Start-Sleep -Seconds 10

Stop-Process -Id $cam.Id -Force
Stop-Process -Id $ana.Id -Force

Write-Host "Lab 18 Complete."
