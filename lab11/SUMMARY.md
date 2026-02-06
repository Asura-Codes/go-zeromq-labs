# Lab 11: High-Availability Broker (Binary Star Pattern)

## Description
This lab demonstrates the **Binary Star** pattern for High Availability (HA). It consists of two broker nodes (Primary and Backup) that monitor each other using a Finite State Machine (FSM). 

- **Normal Operation:** Primary is `Active`, Backup is `Passive`.
- **Failover:** If Primary stops sending heartbeats, Backup becomes `Active`.
- **Recovery:** When Primary returns, it detects the active Backup and stays `Passive` (or takes over depending on logic, here simplified to non-preemptive).

## Architecture
- **Primary Broker:** Starts as Primary. Publishes state on port `6001`, Subscribes to Backup on `6002`. Binds Frontend on `5001`.
- **Backup Broker:** Starts as Backup. Publishes state on port `6002`, Subscribes to Primary on `6001`. Binds Frontend on `5002`.
- **Client App:** Tries to connect to Primary. If it times out, it reconnects to Backup.

## Key Concepts
- **Split-Brain Avoidance:** The FSM handles logic to determine who should be active.
- **Heartbeating:** Continuous stream of `HEARTBEAT` or `I_AM_ACTIVE` messages.
- **Client-Side Reliability:** The client must be smart enough to detect timeouts and switch endpoints.

## Usage
1. **Build:**
    ```powershell
    ./build.ps1
    ```
2. **Run:**
    ```powershell
    ./run.ps1
    ```
    This starts both brokers and the client. You should see "ACK from PRIMARY".

3. **Test Failover:**
    - Kill the `ha_broker_primary.exe` process (e.g., via Task Manager or Ctrl+C if running manually).
    - Watch the Client output. It should detect a timeout and switch to "ACK from BACKUP".
