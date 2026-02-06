# Lab 05: Asynchronous Audit Gateway (Router-Dealer)

## Description
This lab demonstrates the **ROUTER-DEALER** pattern, which allows for a fully asynchronous, non-blocking service gateway.
- **Audit Gateway (Broker):** Sits between Clients and Workers. It accepts requests from multiple clients and dispatches them to available workers.
- **Archival Worker (Dealer):** Connects to the backend of the Gateway. It processes audit logs (simulating slow disk I/O) and sends acknowledgments back.
- **Audit Client (Req):** Sends audit logs to the Gateway and awaits confirmation.

## Architecture
- **Protocol:** TCP
- **Socket Types:** 
  - Gateway Frontend: `ROUTER` (Accepts concurrent connections, tracks identity).
  - Gateway Backend: `DEALER` (Load balances to workers).
  - Worker: `DEALER` (Connects to backend, handles asynchronous reply).
  - Client: `REQ` (Standard synchronous request).
- **Pattern:** ROUTER-DEALER Broker (also known as a Shared Queue).

## Key Concepts
1.  **Identity Frames:** The `ROUTER` socket automatically prepends the sender's identity to the message. The `DEALER` socket (at the broker backend) preserves this identity when forwarding to workers. The Worker must manually handle this "Envelope" to reply to the correct client.
2.  **Asynchronous Handling:** Unlike `REQ-REP` (which is strictly lock-step), `ROUTER-DEALER` allows the Gateway to process thousands of requests in parallel without blocking.
3.  **Load Balancing:** The `DEALER` backend automatically distributes messages to connected workers in a round-robin fashion.

## Execution
Run `run.ps1` to start the Gateway, two Workers, and a group of Clients.
You will observe that:
- Clients send requests simultaneously.
- Workers pick them up in a round-robin fashion.
- Responses find their way back to the correct Client ID.
