# Lab 19: Native Ironhouse (CurveZMQ)

## Scenario
A high-security environment requiring native transport-layer encryption to protect command and control channels.

## Pattern: Ironhouse (CurveZMQ)
This lab utilizes ZeroMQ's built-in Curve25519 security mechanism for authentication and encryption, leveraging the native `libzmq` library via CGO.

### Architecture
- **Key Generation:** A `keygen` utility generates separate keypairs for the server and client, saving them to a shared JSON format.
- **Secure Server:** Loads its keys from `server_keys.json`, proves its identity, and encrypts traffic. Auto-shutdown is configurable.
- **Secure Client:** Loads its keys from `client_keys.json`, verifies the server's public key, and performs a secure handshake before exchanging data.

## Key Concepts
- **CurveZMQ:** Elliptic curve cryptography (Curve25519) integrated directly into the transport layer.
- **Z85 Encoding:** A binary-to-text encoding optimized for ZeroMQ key representation.
- **Docker Orchestration:** Simplifies the complex CGO/ZeroMQ dependency management and securely shares keys between containers.

## Running the Lab

### Option 1: Docker (Recommended)
This method avoids installing local C compilers and ZeroMQ libraries. It automatically generates keys and orchestrates the exchange.

1. Ensure Docker Desktop is running.
2. Run `./run_docker.ps1` (or `docker-compose up --build`).

### Option 2: Local Execution
> **Note:** Requires `libzmq` and a C compiler (MinGW/GCC) in your PATH.

1. Build binaries: `./build.ps1`
2. Generate keys: `go run ./cmd/keygen`
3. Start server: `./secure_server.exe -timeout 35s`
4. Start client: `./secure_client.exe -duration 30s`

## Configuration
Both binaries support command-line flags for automation:
- **Server:** `-timeout <duration>` (e.g., `35s`) - Shuts down automatically after the specified time.
- **Client:** `-duration <duration>` (e.g., `30s`) - Keeps sending messages for the specified time.