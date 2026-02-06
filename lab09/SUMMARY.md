# Lab 09: Encrypted Command & Control (Application-Layer)

## Description
This lab demonstrates **Ironhouse** security (Authentication + Confidentiality) using **Application-Layer Encryption** with Curve25519 (NaCl Box).

**Crucially, the ZeroMQ implementation used in this project (`github.com/go-zeromq/zmq4`) does not offer built-in encryption for sessions.** While the official ZeroMQ protocol (ZMTP 3.0) defines the *CurveZMQ* security mechanism, it is not implemented in the version of the pure-Go driver we are restricted to.

Therefore, to achieve a secure "Ironhouse" architecture, we must manually implement encryption at the application layer. This ensures that:
1.  **Confidentiality:** All payloads are sealed before entering the socket.
2.  **Authentication:** Only clients with the correct private key (matching the server's known public key) can send valid commands.
3.  **Integrity:** Messages are signed and cannot be tampered with.

## Architecture
- **Protocol:** TCP carrying NaCl Boxed Payloads.
- **Crypto:** `golang.org/x/crypto/nacl/box` (Curve25519, XSalsa20, Poly1305).
- **Components:**
  - **KeyGen:** A utility to generate valid Curve25519 keypairs.
  - **Server:** Holds `ServerSecret` + `ClientPublic`. Decrypts requests, encrypts responses.
  - **Client:** Holds `ClientSecret` + `ServerPublic`. Encrypts requests, decrypts responses.

## Setup & Usage

### 1. Generate Keys
In a secure system, keys must be provisioned. Run the included utility to generate a fresh set of keys:
```powershell
go run ./cmd/keygen/main.go
```
Copy the output (Public/Secret keys for Server and Client).

### 2. Configure
Open `internal/config/config.go` and replace the constants with the values you generated.
```go
const (
    ServerPublicKey = "..."
    ServerSecretKey = "..."
    ClientPublicKey = "..."
    ClientSecretKey = "..."
)
```

### 3. Run
Start the simulation:
```powershell
./run.ps1
```
You will see:
- The **Agent** sending encrypted payloads.
- The **Server** successfully decrypting and processing them.
- The **Agent** receiving and decrypting the reply.

## Key Concepts
- **Application vs. Transport Security:** Since the transport (ZeroMQ session) is cleartext in this implementation, we treat the network as hostile and encrypt the data *before* handing it to ZeroMQ.
- **Public Key Cryptography:** The server and client never share their secret keys. They only need to know each other's *Public* keys to establish a secure channel.
- **Nonce Management:** Unique nonces are essential for security to prevent replay attacks and allow the same message to be encrypted differently each time. Our `security` package handles this by prepending the nonce to the ciphertext.
