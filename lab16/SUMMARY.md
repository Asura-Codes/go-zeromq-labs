# Lab 16: WebSocket Bridge

## Scenario
A security operations center (SOC) requires a web-based real-time dashboard to monitor system alerts. Since web browsers cannot connect directly to ZeroMQ sockets, a bridge is required to translate ZeroMQ messages into a format browsers understand: WebSockets.

## Pattern: Protocol Bridging
This lab demonstrates the **Protocol Bridging** pattern (often referred to in ZMQ as a "Streamer" when bridging different transports).

### Architecture
1.  **Telemetry Producer:** Simulates security events (Login, File Access, etc.) and publishes them over a ZeroMQ `PUB` socket using JSON serialization.
2.  **ZMQ Bridge (The "Streamer"):**
    -   Acts as a ZeroMQ `SUB` client, receiving messages from the producer.
    -   Runs a WebSocket server (using `gorilla/websocket`).
    -   Maintains a `Hub` of connected browser clients and broadcasts ZMQ payloads to all of them.
    -   Serves the static HTML/JS dashboard.
3.  **Web Dashboard:** A Bootstrap-based frontend that connects to the WebSocket and dynamically updates a table with incoming events.

## Key Concepts
- **Interoperability:** Extending ZeroMQ's reach to non-ZMQ environments (Web).
- **Serialization:** Using JSON as the common denominator between Go backend and JS frontend.
- **Concurrency:** Handling multiple concurrent WebSocket connections while processing a high-frequency ZMQ stream.

## Running the Lab
1. Run `./run.ps1`.
2. Open `http://localhost:8080` in your browser.
3. Observe the real-time security events flowing into the table.
