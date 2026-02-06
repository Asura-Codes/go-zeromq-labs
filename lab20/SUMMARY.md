# Lab 20: Distributed Tracing System (The "Spy" Pattern)

## Scenario
Monitoring request flow through a microservices architecture without adding latency to the critical path.

## Pattern: PUSH-PULL (Telemetry Plane)
Services use a unidirectional **PUSH** socket to send telemetry data (spans) to a centralized **PULL** collector. 

### Key Features
- **Side-Channel Logging:** Telemetry is sent via a dedicated transport plane.
- **Asynchronous Delivery:** ZeroMQ's internal queuing allows the application logic to continue without waiting for the collector to acknowledge receipt.
- **Horizontal Scaling:** Multiple collectors can be placed behind a load balancer (or multiple PULL instances can bind to the same address if using a proxy) to handle high telemetry volume.

## Implementation Details
- **Collector:** Listens on a PULL socket, unmarshals JSON spans, and logs them with Trace ID correlation.
- **Services:** Every "request" generates a Span object containing metadata, timing, and status, which is then PUSHed to the collector.

## How to Run
1. Build: `.\build.ps1`
2. Run: `.\run.ps1`
3. Observe the collector output as it receives spans from `gateway-service`, `auth-service`, and `db-service`.
