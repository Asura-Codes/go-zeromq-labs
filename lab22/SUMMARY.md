# Lab 22: High-Observability Message Broker (Manual Proxy + Events)

## Scenario
In a production environment, message brokers often become "black boxes". This lab implements a broker with manual message pumping, allowing for custom event logging and observability of the traffic flow.

## Pattern: Manual Proxy
Instead of using a high-level `ZProxy` actor, we use two separate sockets (`ROUTER` and `DEALER`) and manually forward messages between them using Go routines. This provides full control over the message lifecycle.

## Key Features
- **Pure-Go Implementation:** No CGO or external C libraries required.
- **Traffic Logging:** Every message forwarded through the broker is logged with its frame count.
- **Asynchronous Flow:** Frontend-to-Backend and Backend-to-Frontend flows run in parallel.

## How to Run
1. Start the broker:
   ```powershell
   go run ./cmd/monitored_broker
   ```
2. Start one or more clients:
   ```powershell
   go run ./cmd/client_app
   ```
3. Observe the `[EVENT]` logs in the broker terminal as messages are proxied.