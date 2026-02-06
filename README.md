# Advanced Distributed Systems with ZeroMQ and Go

[![Go Reference](https://pkg.go.dev/badge/github.com/go-zeromq/zmq4.svg)](https://pkg.go.dev/github.com/go-zeromq/zmq4)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Status: Active](https://img.shields.io/badge/Status-Active-success.svg)]()

A comprehensive curriculum of **22 standalone laboratories** exploring advanced distributed system patterns. This project demonstrates how to build resilient, high-performance architectures using **Go** and **ZeroMQ**.

> **Note:** This curriculum and the implementation code were developed with the assistance of **Google Gemini**, serving as an AI pair programmer and architectural consultant.

## 📚 Overview

This repository moves beyond basic "Hello World" examples to tackle real-world scenarios found in cybersecurity, high-frequency trading, and microservices orchestration. Each lab is a self-contained Go module focusing on a specific architectural pattern or problem.

**Key Topics Covered:**
*   **Messaging Patterns:** XPUB/XSUB, ROUTER-DEALER, PUSH-PULL, PAIR.
*   **Reliability:** Binary Exponential Backoff, Heartbeating, Circuit Breaking.
*   **Security:** Ironhouse Pattern (Curve25519), Native Transport Encryption.
*   **Consensus & State:** Raft-like Leader Election, Distributed Locks, Distributed Hash Tables (DHT).
*   **Performance:** Zero-Copy Networking, Load Balancing, Binary Star Failover.

## 🛠️ Prerequisites

*   **Go 1.22+**
*   **ZeroMQ (libzmq)** (Required only for Lab 19 and Lab 16 CGO bindings; most labs use pure-Go `zmq4`)
*   **Docker & Docker Compose** (For orchestration labs like Lab 10)
*   **PowerShell** (For automated build/run scripts on Windows)

## 🚀 Getting Started

Each lab is an independent module. You can run them individually or use the provided orchestration scripts.

### Running a Lab (Automated)
Most labs include a `run.ps1` script that builds and launches all necessary components (Publisher, Subscriber, Broker, etc.) in the correct order.

```powershell
cd lab01
./run.ps1
```

### Running Manually
You can also run components separately to see their output clearly:

```powershell
# Terminal 1
cd lab01
go run ./cmd/monitor_agent

# Terminal 2
cd lab01
go run ./cmd/dashboard
```

---

## 🗺️ Roadmap & Curriculum

| Lab | Pattern / Scenario | Key Concepts | Documentation |
| :--- | :--- | :--- | :--- |
| **01** | **Cluster Heartbeat** (PUB-SUB) | Topic filtering, Health monitoring | [Summary](./lab01/SUMMARY.md) |
| **02** | **Log Ingestion Pipeline** (PUSH-PULL) | Parallel processing, Backpressure | [Summary](./lab02/SUMMARY.md) |
| **03** | **Remote Node Inspector** (REQ-REP) | RPC, Blocking vs Non-blocking I/O | [Summary](./lab03/SUMMARY.md) |
| **04** | **Real-Time Telemetry** (XPUB-XSUB) | Last Value Caching (LVC), Multipart messages | [Summary](./lab04/SUMMARY.md) |
| **05** | **Audit Gateway** (ROUTER-DEALER) | Async Request-Reply, Identity frames | [Summary](./lab05/SUMMARY.md) |
| **06** | **Malware Scanning Cluster** (Load Balancer) | LRU Routing, Dynamic workers | [Summary](./lab06/SUMMARY.md) |
| **07** | **Policy Sync** (Clone Pattern) | State replication, Eventual consistency | [Summary](./lab07/SUMMARY.md) |
| **08** | **Resilient Edge Uplink** (Paranoid Pirate) | Binary Exponential Backoff, Retry logic | [Summary](./lab08/SUMMARY.md) |
| **09** | **Encrypted C2** (Ironhouse) | Curve25519 (Pure Go), ZAP Authentication | [Summary](./lab09/SUMMARY.md) |
| **10** | **SOC Simulation** (Docker Orchestration) | Containerization, Service Discovery | [Summary](./lab10/SUMMARY.md) |
| **11** | **HA Broker** (Binary Star) | Active-Passive Failover, FSM | [Summary](./lab11/SUMMARY.md) |
| **12** | **Service Broker** (Majordomo) | Reliable Service Oriented Architecture | [Summary](./lab12/SUMMARY.md) |
| **13** | **Persistent Queue** (Titanic) | Store-and-forward, Asynchronous ACK | [Summary](./lab13/SUMMARY.md) |
| **14** | **DHT Ring** (P2P Architecture) | Consistent Hashing, Peer Discovery | [Summary](./lab14/SUMMARY.md) |
| **15** | **MapReduce** (Scatter-Gather) | Data Partitioning, Result Aggregation | [Summary](./lab15/SUMMARY.md) |
| **16** | **Web Dashboard** (Protocol Gateway) | WebSocket bridging, Protocol translation | [Summary](./lab16/SUMMARY.md) |
| **17** | **Distributed Consensus** (Raft Prototype) | Leader Election, Quorum Voting | [Summary](./lab17/SUMMARY.md) |
| **18** | **4K Video Stream** (Zero-Copy) | High-throughput binary data, Memory pooling | [Summary](./lab18/SUMMARY.md) |
| **19** | **Native Ironhouse** (CGO CurveZMQ) | Native C bindings, Transport Encryption | [Summary](./lab19/SUMMARY.md) |
| **20** | **Distributed Tracing** (Spy Pattern) | Asynchronous Telemetry, Side-channel logging | [Summary](./lab20/SUMMARY.md) |
| **21** | **Distributed Lock Manager** (Quorum) | Client-side Quorum, Distributed Mutex | [Summary](./lab21/SUMMARY.md) |
| **22** | **Federated Bridge** (Zone Bridging) | WAN Optimization, Subscription Forwarding | [Summary](./lab22/SUMMARY.md) |

## 📂 Project Structure

```text
/
├── lab01/             # Lab Directory
│   ├── cmd/           # Executables (Main packages)
│   ├── internal/      # Private library code
│   ├── go.mod         # Go Module definition
│   ├── run.ps1        # Orchestration script
│   └── SUMMARY.md     # Documentation
├── .vscode/           # Editor configuration
└── README.md          # This file
```

## 🤝 Contributing

Contributions are welcome! Please ensure any Pull Requests adhere to the existing architectural patterns and include a corresponding `SUMMARY.md` update.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
