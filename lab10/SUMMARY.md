# Lab 10: SOC Simulation (Orchestration)

## Description
This lab simulates a **Security Operations Center (SOC)** environment where multiple microservices interact using various ZeroMQ patterns. It demonstrates how to orchestrate a complex distributed system using **Docker Compose**.

## Architecture
The simulation consists of four main components:
1.  **Intel Provider (PUB):** Simulates a threat intelligence feed broadcasting suspect IP addresses.
2.  **Anomaly Detector (REP):** A synchronous service that performs "deep analysis" on an IP and returns a safety verdict.
3.  **Alert Logger (PULL):** A centralized logging service that receives critical security alerts.
4.  **SOC Processor (SUB / REQ / PUSH):** The "brain" of the system.
    - **Subscribes** to the Intel Provider.
    - **Requests** analysis from the Anomaly Detector for every received IP.
    - **Pushes** an alert to the Logger only if the verdict is `MALICIOUS`.

## Key Concepts
- **Mixed Patterns:** Combining PUB-SUB, REQ-REP, and PUSH-PULL in a single workflow.
- **Service Discovery:** Using Docker's internal DNS (e.g., `tcp://intel-provider:5555`) to connect services without hardcoding IPs.
- **Orchestration:** Using `docker-compose` to manage the lifecycle of multiple containers, environment variables, and networks.
- **Protocol Decoupling:** Services interact via shared constants and well-defined messaging patterns.

## Usage
### Docker Compose (Recommended)
```powershell
./run.ps1
```
This will build and start all four containers. You will see the logs interleaved, showing the end-to-end flow from IP detection to alert logging.

### Local Execution (Manual)
If you don't have Docker, you can run the services manually:
1. `go run ./cmd/intel_provider`
2. `go run ./cmd/anomaly_detector`
3. `go run ./cmd/alert_logger`
4. `go run ./cmd/soc_processor`

## Verification
Watch the `alert-logger` container logs. Approximately 30% of the IPs broadcast by the `intel-provider` should result in a `[CRITICAL ALERT]` being logged by the `alert-logger`.
