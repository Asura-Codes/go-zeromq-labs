# Lab 02: High-Volume Log Ingestion Pipeline (Parallel Pipeline)

## Description
This lab demonstrates the **Ventilator-Worker-Sink (Push-Pull)** pattern for parallel data processing. It simulates a log ingestion pipeline:
- **Log Collector (Ventilator):** Generates a batch of raw log entries and pushes them to workers.
- **Log Parser (Worker):** Receives raw logs, performs CPU-intensive tasks (regex redaction of IPs), and pushes the result to the sink.
- **Storage Writer (Sink):** Collects processed logs and aggregates statistics.

## Architecture
- **Protocol:** TCP
- **Socket Types:** `PUSH` (Ventilator/Worker Output), `PULL` (Worker Input/Sink).
- **Topology:** One-to-Many-to-One (Fan-out / Fan-in).

## Advantages
1.  **Parallelism:** The "Worker" stage can scale horizontally. Running 10 workers processes the batch roughly 10x faster than 1 worker.
2.  **Load Balancing:** ZeroMQ's `PUSH` socket automatically load-balances messages among connected workers using a Round-Robin strategy.
3.  **Pipeline Construction:** Easy to chain stages together.

## Disadvantages
1.  **Unidirectional:** There is no feedback loop. If the Sink is slow, the Workers might eventually block, back-propagating pressure to the Ventilator.
2.  **No Reliability:** If a Worker crashes while holding a message, that message is lost. ZeroMQ does not re-queue unacknowledged messages automatically in this pattern.

## Code / Implementation Notes
- Uses `config` package for centralized port management.
- **Potential Issue:** The `Sink` logic relies on detecting the "first" message to start the timer. In a real system, explicit signaling (Start/End of Batch messages) is often robust.
