# Lab 03: Remote Node Inspector (Request-Reply)

## Description
This lab implements the **Request-Reply (REQ-REP)** pattern to simulate a secure remote administration tool.
- **Node Agent (Server/Replier):** Runs on a target machine, listening for commands. It uses `gopsutil` to fetch real system metrics (CPU usage, Memory stats, Host info).
- **Admin CLI (Client/Requester):** A command-line interface that sends specific commands (`CPU`, `MEM`, `HOST`) to the agent and displays the returned JSON data.

## Architecture
- **Protocol:** TCP
- **Socket Types:** `REP` (Server), `REQ` (Client).
- **Flow:** Synchronous blocking call. Client sends -> Client blocks -> Server receives -> Server processes -> Server replies -> Client unblocks.

## Advantages
1.  **Reliability:** Strict send-receive-send-receive cycle ensures the client knows the server processed the specific request.
2.  **Simplicity:** Easy to reason about control flow; similar to HTTP.

## Disadvantages
1.  **Blocking:** If the Server crashes or hangs, the Client blocks indefinitely (unless a timeout is implemented, as done in this lab).
2.  **Low Throughput:** The lock-step nature prevents high-volume message processing compared to PUSH-PULL or PUB-SUB.

## Code / Implementation Notes
- Uses `github.com/shirou/gopsutil` for **real** system data.
- Implements a **Client-side Timeout** using `context.WithTimeout` to handle server unavailability gracefully.
- **Potential Issue:** If the server restarts while the client is waiting, the REQ socket might get stuck in a state expecting a reply. ZMQ REQ sockets are sensitive to the send/recv cycle.
