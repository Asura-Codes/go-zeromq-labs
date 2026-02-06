$ErrorActionPreference = "Stop"

$Root = $PSScriptRoot
if ($Root -eq $null -or $Root -eq "") { $Root = "." }

Write-Host "Starting Lab 15 (MapReduce Cluster) - Robustness & Timeout Test..."

# Build
& "$Root\build.ps1"

# 1. Start Master (Binds MapVent, Control, Sink)
Write-Host "Starting Master (Listening)..."
$master = Start-Process "$Root\mr_master.exe" -ArgumentList "-ready" -PassThru -NoNewWindow

Start-Sleep -Seconds 1

# 2. Start Map Workers
Write-Host "Starting Map Workers..."
$m1 = Start-Process "$Root\map_worker.exe" -PassThru -NoNewWindow
$m2 = Start-Process "$Root\map_worker.exe" -PassThru -NoNewWindow

Start-Sleep -Seconds 1

# 3. Start Reduce Worker
Write-Host "Starting Reduce Worker..."
$r1 = Start-Process "$Root\reduce_worker.exe" -PassThru -NoNewWindow

# Wait for Master to finish, but time out after 15s idleness (handled by master)
# We also add a safety timeout here in PowerShell
Write-Host "Waiting for Master to complete (max 30s total)..."
if (-not $master.WaitForExit(30000)) {
    Write-Host "Master timed out in PowerShell! Forcing termination."
    Stop-Process -Id $master.Id -Force
}

# Cleanup
Write-Host "Cleaning up..."
try { Stop-Process -Id $m1.Id -Force -ErrorAction SilentlyContinue } catch {}
try { Stop-Process -Id $m2.Id -Force -ErrorAction SilentlyContinue } catch {}
try { Stop-Process -Id $r1.Id -Force -ErrorAction SilentlyContinue } catch {}

Write-Host "Done."