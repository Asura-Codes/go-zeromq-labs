# Lab 08: Resilient Edge Uplink (Paranoid Pirate)

## Description
This lab demonstrates the **Paranoid Pirate** pattern (also known as Lazy Pirate in simpler forms), which provides reliability for the Request-Reply pattern over unstable networks.
- **Central Receiver (Server):** Simulates a flaky service that randomly drops requests ("crashes") or processes them very slowly ("overload").
- **Field Device (Client):** Implements a robust retry mechanism. If it doesn't receive a reply within a timeout, it assumes the request failed, closes the socket (to clear any stuck state), opens a new socket, and retries the request.

## Architecture
- **Protocol:** TCP
- **Socket Types:**
  - Server: `ROUTER` (Simulating an async backend, though processing is synchronous here).
  - Client: `REQ`.
- **Pattern:** Paranoid Pirate (Reliable Request-Reply).

## Key Concepts
1.  **Socket Cycling:** In ZeroMQ, a `REQ` socket strictly expects a Reply after a Request. If the request is lost, the socket is "stuck". The only way to recover is to destroy the socket and create a new one.
2.  **Linger Period:** When closing a socket, we set `Linger` to 0. This ensures the library doesn't block trying to send pending messages to a dead server.
3.  **Idempotency:** Because the client retries, the server might receive the same message twice (if the *Reply* was lost, not the Request). The application protocol must handle this (e.g., using unique Sequence IDs), though this lab focuses on the transport mechanics.

## Execution
Run `run.ps1`.
- You will see the Client sending requests.
- The Server will log "SIMULATING CRASH" or "SIMULATING OVERLOAD" occasionally.
- When this happens, the Client will log "Timeout" and "Retrying connection...".
- Eventually, the Client should succeed despite the server's flakiness.
