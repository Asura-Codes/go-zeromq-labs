Get-ChildItem -Path . -Filter '*.exe' | Remove-Item -Force -ErrorAction SilentlyContinue

# Clean up old data
if (Test-Path "titanic_data") {
    Remove-Item -Recurse -Force "titanic_data"
}

Write-Host 'Cleaned executables in lab13'
