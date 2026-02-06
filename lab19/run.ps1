Write-Host "Running Lab 19 via Docker..."

if (-not (Get-Command docker-compose -ErrorAction SilentlyContinue)) {
    Write-Error "docker-compose not found. Please install Docker Desktop."
    exit 1
}

docker-compose up --build --abort-on-container-exit

Write-Host "Cleaning up..."
# Stop and remove containers, networks, volumes, and images created by docker-compose
docker-compose down --volumes --remove-orphans 2>$null

Write-Host "Lab 19 Docker Run Complete."