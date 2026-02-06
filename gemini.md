# ZeroMQ with Go: Advanced Distributed Systems Labs

## Course Overview
This series of labs focuses on building high-performance, resilient distributed systems using ZeroMQ and Go. The scenarios are drawn from **cybersecurity**, **high-throughput data processing**, and **advanced microservices architecture**.

### Core Mandates
- **Independent Modules:** Each lab is a standalone Go project with its own `go.mod`.
- **Separation of Concerns:** Publishers/Servers and Subscribers/Clients are distinct executables.
- **Production Readiness:** Implement structured logging, signal handling (SIGINT/TERM), and configuration via flags/env vars.
- **SOLID Architecture:** Adhere to clean architecture principles.
- **Automation & IDE Integration:** Use `run.ps1` for local orchestration and VS Code's `tasks.json` / `launch.json` for streamlined development and debugging.

## Lab Curriculum

### Lab 01: Cluster Heartbeat Monitor (Basic Pub-Sub)
**Scenario:** A large server farm needs to broadcast health status ("heartbeats") to a monitoring dashboard and an alert system.
**Pattern:** Publish-Subscribe (PUB-SUB).
**Key Concepts:** Topic filtering (subscribing to specific racks or alert levels), decoupling producers from consumers.
**Deliverable:** `monitor_agent.exe` (Publisher - runs on nodes), `dashboard.exe` (Subscriber - visualizes state).

### Lab 02: High-Volume Log Ingestion Pipeline (Parallel Pipeline)
**Scenario:** A Security Information and Event Management (SIEM) system ingests terabytes of raw logs. Logs must be anonymized and formatted before storage.
**Pattern:** Pipeline (PUSH-PULL).
**Key Concepts:** Parallel processing, backpressure handling, ventilator-worker-sink architecture.
**Deliverable:** `log_collector.exe` (Ventilator), `log_parser.exe` (Worker), `storage_writer.exe` (Sink).

### Lab 03: Remote Node Inspector (Request-Reply)
**Scenario:** An administrator needs to query the active process list or memory usage of a remote server securely and synchronously.
**Pattern:** Request-Reply (REQ-REP).
**Key Concepts:** Blocking vs. non-blocking I/O, synchronous remote procedure calls (RPC).
**Deliverable:** `admin_cli.exe` (Client), `node_agent.exe` (Server).

### Lab 04: Real-Time Telemetry Feed (Pub-Sub + Multipart)
**Scenario:** A high-frequency trading or industrial IoT system broadcasts multipart telemetry data (Header + Payload). Subscribers need the "last known value" immediately upon connection.
**Pattern:** XPUB-XSUB (Broker) with Last Value Caching (LVC).
**Key Concepts:** Multipart messages, message envelopes, handling late joiners.
**Deliverable:** `telemetry_source.exe` (Publisher), `lvc_broker.exe` (Proxy), `analyst_terminal.exe` (Subscriber).

### Lab 05: Asynchronous Audit Gateway (Router-Dealer)
**Scenario:** A secure gateway receives audit logs from thousands of concurrent clients. It must not block the clients while the backend writes to slow archival storage.
**Pattern:** Router-Dealer (ROUTER-DEALER).
**Key Concepts:** Asynchronous Request-Reply, identity frames, non-blocking sockets, request correlation.
**Deliverable:** `audit_gateway.exe` (Router), `archival_worker.exe` (Dealer).

### Lab 06: Scalable Malware Scanning Cluster (Load Balancing)
**Scenario:** A file upload service requires incoming files to be scanned by a dynamic cluster of antivirus engines. Distribution must be load-balanced based on worker availability.
**Pattern:** Load Balancing Broker (ROUTER-ROUTER).
**Key Concepts:** Least-recently-used (LRU) routing, high availability, dynamic worker registration.
**Deliverable:** `upload_service.exe` (Client), `scanner_broker.exe` (Broker), `av_engine.exe` (Worker).

### Lab 07: Global Policy Synchronization (Clone Pattern)
**Scenario:** A distributed firewall system needs to sync security policies (Access Control Lists) across all edge nodes in real-time. New nodes must fetch the full policy state on startup.
**Pattern:** Clone (Server-to-Client State Replication).
**Key Concepts:** State snapshots, delta updates, eventual consistency, key-value storage.
**Deliverable:** `policy_master.exe` (Publisher), `firewall_node.exe` (Subscriber).

### Lab 08: Resilient Edge Uplink (Paranoid Pirate)
**Scenario:** Remote field devices (drones or IoT gateways) transmit critical data over unstable connections. If the central receiver fails, the device must retry and failover to a backup.
**Pattern:** Paranoid Pirate (Reliable Request-Reply).
**Key Concepts:** Heartbeating, client-side reliability, binary exponential backoff, circuit breaking.
**Deliverable:** `field_device.exe` (Client), `central_receiver.exe` (Server).

### Lab 09: Encrypted Command & Control (Ironhouse Pattern)
**Scenario:** A secure enclave requires a command channel that is authenticated and encrypted. Unloading the encryption to the application layer is not sufficient; the transport itself must be secure.
**Pattern:** Ironhouse (CurveZMQ Security).
**Key Concepts:** Curve25519 encryption, ZeroMQ authentication protocol (ZAP), client/server certificate management.
**Deliverable:** `c2_server.exe`, `secure_agent.exe`.

### Lab 10: SOC Simulation (Orchestration)
**Scenario:** A full Security Operations Center (SOC) simulation. Threat Intel feeds (Pub/Sub), Anomaly Detection (Req/Rep), and Alerting Services (Pipeline) interacting in a containerized environment.
**Pattern:** Mixed Architectures with Docker Compose.
**Key Concepts:** Service discovery, network isolation, integration testing, container orchestration.
**Deliverable:** `docker-compose.yml`, multiple interacting Go services.

### Lab 11: High-Availability Broker (Binary Star Pattern)
**Scenario:** The central message broker is a single point of failure. We need an active-passive pair of brokers that monitor each other and failover automatically.
**Pattern:** Binary Star (Primary-Backup).
**Key Concepts:** Finite state machines, primary/backup coordination, split-brain resolution.
**Deliverable:** `ha_broker_primary.exe`, `ha_broker_backup.exe`, `client_app.exe`.
                                                          
### Lab 12: Service-Oriented Broker (Majordomo Pattern)
**Scenario:** A robust service-oriented architecture where clients request services by name, and a broker dispatches these requests to available workers, handling retries and worker ilures transparently.
**Pattern:** Majordomo (Service Broker).
**Key Concepts:** Service discovery protocol (MDP), worker heartbeating, reliable request dispatching.
**Deliverable:** `majordomo_broker.exe`, `echo_worker.exe`, `client_requester.exe`.
                                                          
### Lab 13: Persistent Message Queue (Titanic Pattern)
**Scenario:** Clients need to send requests that must eventually be processed, even if the service is currently offline. The system requires a "store-and-forward" mechanism.
**Pattern:** Titanic (Persistent Queuing).
**Key Concepts:** Persistent storage of messages, asynchronous acknowledgement, eventual consistency.
**Deliverable:** `titanic_broker.exe`, `storage_service.exe`, `patient_client.exe`.
                                                          
### Lab 14: Distributed Hash Table Ring (P2P Architecture)
**Scenario:** A decentralized key-value store where data is distributed across a ring of nodes. Nodes join and leave dynamically.
**Pattern:** Distributed Hash Table (P2P Ring).
**Key Concepts:** Consistent hashing, peer discovery, gossip protocols (using ZMQ beacon), ring topology.
**Deliverable:** `dht_node.exe`, `client_put_get.exe`.

### Lab 15: MapReduce Cluster (Scatter-Gather)
**Scenario:** A massive dataset needs to be processed. A master node splits the data, distributes it to workers (Map), and then aggregates the results (Reduce).
**Pattern:** Scatter-Gather (Parallel Pipeline).
**Key Concepts:** Data partitioning, partial result aggregation, synchronization of parallel tasks.
**Deliverable:** `mr_master.exe` (Ventilator/Sink), `map_worker.exe` (Worker), `reduce_worker.exe` (Worker).

### Lab 16: WebSocket Bridge (Protocol Bridging)
**Scenario:** A web-based dashboard needs to visualize ZMQ streams in real-time. Since browsers don't speak native ZMQ, we need a bridge.
**Pattern:** Streamer (Bridging TCP to WebSocket).
**Key Concepts:** Protocol translation, integration with a web frontend, `gorilla/websocket` integration.
**Deliverable:** `zmq_bridge.exe`, `web_client.html`.

### Lab 17: Distributed Consensus (Raft Prototype)
**Scenario:** A cluster of critical database nodes must agree on a sequence of state changes. They need a consensus mechanism to elect a leader and replicate logs.
**Pattern:** PUB-SUB Mesh (Gossip Consensus).
**Key Concepts:** Leader election timeouts, quorum-based voting, heartbeat synchronization, PUB-SUB event bus.
**Deliverable:** `consensus_node.exe`, `client_proposer.exe`.

### Lab 18: Zero-Copy Video Streaming (High Performance)
**Scenario:** A live video surveillance system needs to distribute 4K video frames to multiple analytics engines with minimal latency and CPU usage.
**Pattern:** Pub-Sub with Zero-Copy.
**Key Concepts:** Memory pools, avoiding GC pressure, zero-copy message handling, high-throughput binary data.
**Deliverable:** `camera_feed.exe`, `analytics_engine.exe`.

### Lab 19: Native Ironhouse (CurveZMQ with CGO)
**Scenario:** High-security environment requiring native transport-layer encryption. Unlike application-layer encryption, this utilizes ZeroMQ's built-in Curve25519 implementation for transparent, high-performance security.
**Pattern:** Ironhouse (Native CurveZMQ).
**Key Concepts:** CGO bindings (`pebbe/zmq4`), Z85 key encoding, native transport encryption, `ZMQ_CURVE_SERVER`.
**Deliverable:** `secure_server.exe`, `secure_client.exe`, `keygen.exe`.

### Lab 20: Distributed Tracing System (The "Spy" Pattern)
**Scenario:** A microservices request travels through multiple services. We need to visualize the latency and success of each hop without impacting the performance of the main request path.
**Pattern:** PUSH-PULL (Telemetry Plane).
**Key Concepts:** Side-channel logging, asynchronous telemetry, non-blocking I/O, span correlation.
**Deliverable:** `trace_collector.exe`, `monitored_service.exe`.

### Lab 21: Distributed Lock Manager (The "Chubby" Clone)
**Scenario:** Multiple writer nodes need to update a shared resource, but only one can hold the lock at a time. The system must handle client crashes by expiring stale locks.
**Pattern:** ROUTER-DEALER (Lease Management).
**Key Concepts:** Distributed mutual exclusion, lease TTL (Time-To-Live), heartbeating, asynchronous queueing.
**Deliverable:** `lock_server.exe`, `lock_client.exe`.

### Lab 22: High-Observability Message Broker (Manual Proxy)
**Scenario:** A central message hub that requires real-time monitoring of traffic flow for audit and debugging, implemented without high-level actors.
**Pattern:** Manual Proxy (Router-Dealer).
**Key Concepts:** Manual message pumping, asynchronous forwarding, traffic observability, custom event logging.
**Deliverable:** `monitored_broker.exe`, `client_app.exe`.

### Lab 23: Decentralized Cluster Gossip (PUB-SUB Pattern)
**Scenario:** A peer-to-peer network of compute nodes where nodes join and leave dynamically without a central broker, using decentralized state sharing.
**Pattern:** Decentralized PUB-SUB.
**Key Concepts:** Peer discovery, decentralized state sharing, PUB-SUB mesh, eventually consistent cluster state.
**Deliverable:** `gossip_node.exe`.

### Lab 24: Autonomous Secure Beaconing (UDP + REP/REQ)
**Scenario:** Industrial IoT devices that need to find each other on a local network automatically and establish secure channels using standard network protocols.
**Pattern:** UDP Discovery + REP/REQ.
**Key Concepts:** UDP discovery (beaconing), zero-config networking, native messaging, transport-layer security.
**Deliverable:** `secure_device.exe`, `network_scanner.exe`.

## Project Structure

Each lab is a **completely separate Go module**.

```text
/
├── .vscode/
│   ├── tasks.json
│   └── launch.json
├── gemini.md
├── lab01/
│   ├── go.mod
│   ├── go.sum
│   ├── build.ps1
│   ├── run.ps1
│   └── SUMMARY.md
...
```

## Setup Instructions

1.  **Initialize Lab:**
    Enter the lab directory and initialize the module.
    Most labs (01-18, 20-24) use the pure-Go driver:
    ```powershell
    cd lab01
    go mod init gemini-zeromq-labs/lab01
    go get github.com/go-zeromq/zmq4
    ```
    **Note for Lab 19 (Native Encryption):**
    Lab 19 requires CGO and the `pebbe/zmq4` library. It is designed to be run in a containerized environment to handle C dependencies:
    ```powershell
    cd lab19
    ./run.ps1  # Orchestrates Docker Compose
    ```

2.  **Develop:**
    Implement the solution in `cmd/` and `internal/`. Use `internal/` for shared logic between the binaries within that specific lab. Leverage VS Code's **Run and Debug** (F5) using the provided `launch.json`.

3.  **Run & Orchestrate:**
    Use the provided PowerShell script for automated execution:
    ```powershell
    ./run.ps1
    ```
    Alternatively, run the executables manually in separate terminals:
    ```powershell
    go run ./cmd/monitor_agent
    go run ./cmd/dashboard
    ```