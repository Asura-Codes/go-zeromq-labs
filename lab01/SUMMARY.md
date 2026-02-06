# Lab 01: Cluster Heartbeat Monitor (Basic Pub-Sub)

## Description
This lab implements a basic **Publish-Subscribe (PUB-SUB)** pattern using ZeroMQ. It simulates a distributed monitoring system where:
- **Monitor Agent (Publisher):** Runs on a node, generates mock hardware metrics (CPU, Memory), and publishes them to a specific topic.
- **Dashboard (Subscriber):** Connects to the agent and filters messages based on the topic to visualize the system state.

## Architecture
- **Protocol:** TCP
- **Socket Types:** `PUB` (Publisher), `SUB` (Subscriber)
- **Data Format:** JSON payload inside a Multipart ZMQ message (`[Topic, Payload]`).

## Advantages
1.  **Decoupling:** The Agent does not know who is listening. New Dashboards can be added without modifying the Agent.
2.  **Scalability:** Efficient multicasting of state data.
3.  **Simplicity:** Minimal boilerplate for 1-to-N communication.

## Disadvantages
1.  **"Slow Joiner" Syndrome:** If the Subscriber starts after the Publisher, it misses all messages sent before it connected. ZeroMQ does not buffer history by default.
2.  **Unreliable Transport:** In PUB-SUB, if a subscriber disconnects or the network drops, messages are lost. There are no ACKs.

## Code / Implementation Notes
- Uses `time.Ticker` for periodic updates.
- Uses `slog` for structured logging.
- Includes a basic `run.ps1` for orchestration.
- **Potential Issue:** If `monitor_agent` is started long before `dashboard`, the dashboard will show nothing initially. This is expected behavior in raw PUB-SUB.