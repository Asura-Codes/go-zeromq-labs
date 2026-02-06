Get-ChildItem -Path . -Filter '*.exe' | Remove-Item -Force -ErrorAction SilentlyContinue
Write-Host 'Cleaned executables in lab18'
