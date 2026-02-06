# Lab 04: Real-Time Telemetry Feed (Pub-Sub + Multipart)

## Description
This lab simulates a high-frequency telemetry system using the **XPUB-XSUB** pattern to create a Last Value Caching (LVC) Broker.
- **Telemetry Source (Publisher):** Broadcasts random updates for multiple sensors (`sensors/temp`, `sensors/pressure`).
- **LVC Broker (Proxy):** sits between Publishers and Subscribers. It caches the *last* message seen for every topic. When a new subscriber joins, the broker immediately re-sends the cached value for that topic.
- **Analyst Terminal (Subscriber):** connects to the broker and subscribes to `sensors/temp`. It receives the "Last Known Value" immediately upon connection, even if the source publishes slowly.

## Architecture
- **Protocol:** TCP
- **Socket Types:** `XPUB` (Broker Frontend), `XSUB` (Broker Backend), `PUB` (Source), `SUB` (Terminal).
- **Pattern:** Publish-Subscribe with Intermediary (Broker).

## Advantages
1.  **Late Joiner Support:** Subscribers don't have to wait for the next update to know the current state.
2.  **Decoupling:** Publishers and Subscribers are completely isolated by the Broker.
3.  **Network Efficiency:** The "Re-publish on Subscription" mechanism is handled locally by the broker, not burdening the original publisher.

## Disadvantages
1.  **Complexity:** Requires a custom Proxy loop instead of the standard `zmq_proxy`.
2.  **Stale Data:** If the publisher dies, the broker continues to serve the "Last Value" which might be old. (Can be mitigated with heartbeats or TTLs).

## Code / Implementation Notes
- The Broker manually handles `XPUB` subscription frames (first byte `0x01`).
- Uses a `map[string][][]byte` for caching multipart messages.
- **Key Concept:** `XPUB` sockets allow the application to receive subscription events as messages.
