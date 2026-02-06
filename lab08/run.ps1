# Lab 08 Orchestration Script
$ErrorActionPreference = "Stop"

Write-Host "Initializing Lab 08..."
go mod tidy

# Build first
./build.ps1

Write-Host "Starting Lab 08 processes (output streamed to this console)..."

$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Definition

# --- Managed helpers (C#) -------------------------------------------------
$cs = @"
using System;
using System.Diagnostics;
using System.Collections.Generic;

public class StreamWriterHandler {
	private string label;
	public StreamWriterHandler(string label) { this.label = label; }
	public void OnOutput(object sender, DataReceivedEventArgs e) { if (e.Data != null) Console.WriteLine("[{0}] {1}", label, e.Data); }
	public void OnError(object sender, DataReceivedEventArgs e) { if (e.Data != null) Console.WriteLine("[{0}][ERR] {1}", label, e.Data); }
}

public static class CtrlCHandler {
	private static List<int> pids = new List<int>();
	public static void RegisterPids(int[] ids) { foreach (var id in ids) if (!pids.Contains(id)) pids.Add(id); }
	public static void OnCancel(object sender, ConsoleCancelEventArgs e) {
		e.Cancel = true;
		foreach (var id in pids.ToArray()) {
			try { var proc = Process.GetProcessById(id); if (!proc.HasExited) proc.Kill(); } catch { }
		}
	}
}
"@

if (-not ([System.Management.Automation.PSTypeName]'StreamWriterHandler').Type) {
	Add-Type -TypeDefinition $cs -Language CSharp
}

# --- PowerShell helper functions -------------------------------------------------
function Start-ChildProcess([string]$exeName, [string]$label) {
	$fullPath = Join-Path $scriptDir $exeName
	if (-not (Test-Path $fullPath)) { throw "Executable not found: $fullPath" }

	$psi = New-Object System.Diagnostics.ProcessStartInfo $fullPath
	$psi.WorkingDirectory = $scriptDir
	$psi.RedirectStandardOutput = $true
	$psi.RedirectStandardError = $true
	$psi.UseShellExecute = $false
	$psi.CreateNoWindow = $true

	$p = New-Object System.Diagnostics.Process
	$p.StartInfo = $psi
	$p.EnableRaisingEvents = $true

	$handlerObj = New-Object StreamWriterHandler $label
	$delOut = [System.Delegate]::CreateDelegate([System.Diagnostics.DataReceivedEventHandler], $handlerObj, 'OnOutput')
	$delErr = [System.Delegate]::CreateDelegate([System.Diagnostics.DataReceivedEventHandler], $handlerObj, 'OnError')
	$p.add_OutputDataReceived($delOut)
	$p.add_ErrorDataReceived($delErr)

	if (-not $p.Start()) { throw "Failed to start $fullPath" }

	$p.BeginOutputReadLine()
	$p.BeginErrorReadLine()
	return $p
}

function Register-CtrlCHandler([System.Diagnostics.Process[]]$procs) {
	$ids = $procs | ForEach-Object { $_.Id }
	[CtrlCHandler]::RegisterPids(@($ids))
	return [System.Delegate]::CreateDelegate([System.ConsoleCancelEventHandler], [CtrlCHandler], 'OnCancel')
}

function Stop-Processes([System.Diagnostics.Process[]]$procs) {
	foreach ($p in $procs) {
		try { if ($p -and -not $p.HasExited) { $p.Kill() } } catch { }
	}
}

function Run-Lab {
	$pServer = $null; $pClient = $null; $cancelDel = $null
	try {
		$pServer = Start-ChildProcess "central_receiver.exe" "Server"
		Start-Sleep -Seconds 1
		$pClient = Start-ChildProcess "field_device.exe" "Client"

		$cancelDel = Register-CtrlCHandler -procs @($pServer, $pClient)
		[console]::add_CancelKeyPress($cancelDel)

		Write-Host "Lab 08 running. Press Ctrl+C to stop."

		while (-not $pServer.HasExited) { 
			Start-Sleep -Milliseconds 200 
		}
	}
	catch {
		Write-Host "Error: $($_.Exception.Message)"
	}
	finally {
		Stop-Processes @($pServer, $pClient)
		if ($cancelDel) { [console]::remove_CancelKeyPress($cancelDel) }
		Write-Host "All processes stopped."
		Write-Host "Exiting Lab 08."
	}
}

Run-Lab
