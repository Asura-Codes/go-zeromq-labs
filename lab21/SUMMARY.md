# Lab 21: Distributed Lock Manager (The "Chubby" Clone)

## Scenario
Coordinating access to a shared resource among multiple writers. Only one writer can hold the lock at a time.

## Pattern: ROUTER-DEALER (Lease Management)
A central **ROUTER** server manages lock states, while **DEALER** clients request, maintain (heartbeat), and release locks.

### Key Features
- **Lease-Based Locking:** Locks are not permanent; they have a Time-To-Live (TTL).
- **Heartbeating:** Clients must send periodic heartbeats to renew their lease.
- **Auto-Expiration:** If a client crashes or loses connection, the server automatically releases the lock after the TTL expires, preventing deadlocks.
- **Asynchronous Communication:** The ROUTER socket allows the server to handle lock requests from many clients without blocking.

## Implementation Details
- **Server:** Tracks active locks in a map and runs a background goroutine to purge expired ones.
- **Client:** Uses a state machine (Request -> Work+Heartbeat -> Release).

## How to Run
1. Build: `.\build.ps1`
2. Run: `.\run.ps1`
3. Watch the logs to see how Alice and Bob negotiate access to the `data-file`. When one holds the lock, the other is denied until the lock is released or expires.
