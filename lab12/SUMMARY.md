# Lab 12: Service-Oriented Broker (Majordomo Pattern)

## Description
This lab demonstrates the **Majordomo Pattern** (MDP), a robust Service-Oriented Architecture (SOA) pattern using ZeroMQ. It allows clients to request services by name rather than connecting to specific workers, and handles worker registration, load balancing, and reliability.

## Architecture
- **Majordomo Broker:** The central hub.
    - **Frontend (ROUTER):** Accepts requests from Clients (`MDPC01` protocol).
    - **Backend (ROUTER):** Manages connections to Workers (`MDPW01` protocol).
    - Maintains a registry of services and queues requests if no workers are available.
- **Echo Worker:** Implements the `echo` service.
    - Connects to the Broker Backend.
    - Sends `READY` to register.
    - Replies to `REQUEST` messages.
    - Sends periodic `HEARTBEAT`s.
- **Client Requester:**
    - Connects to the Broker Frontend.
    - Sends requests for the `echo` service.

## Key Concepts
- **Service Discovery:** Clients don't know where workers are; they only know the Service Name (e.g., "echo").
- **Load Balancing:** The broker distributes requests among available workers for a service.
- **Protocol Envelopes:** Using `MDPC01` and `MDPW01` headers to structure complex interactions.

## Usage
1. **Build:**
    ```powershell
    ./build.ps1
    ```
2. **Run:**
    ```powershell
    ./run.ps1
    ```
    You should see the Client sending "Hello Majordomo X" and receiving the echo reply.
