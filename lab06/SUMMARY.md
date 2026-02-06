# Lab 06: Scalable Malware Scanning Cluster (Load Balancing)

## Description
This lab implements a **Load Balancing Broker** using the **ROUTER-ROUTER** pattern. This pattern allows for precise control over message routing, enabling features like dynamic worker registration and Least-Recently-Used (LRU) distribution.
- **Scanner Broker:** Accepts files from clients and distributes them to the next available antivirus engine. It maintains a queue of ready workers.
- **AV Engine:** Connects to the broker, signals readiness, scans files (simulated), and returns the verdict.
- **Upload Service:** Simulates multiple users uploading files for scanning.

## Architecture
- **Protocol:** TCP
- **Socket Types:**
  - Broker Frontend: `ROUTER` (Clients).
  - Broker Backend: `ROUTER` (Workers).
  - Worker: `DEALER` (Connects to Backend).
  - Client: `REQ` (Connects to Frontend).
- **Pattern:** Load Balancing Pattern (ROUTER-ROUTER).

## Key Concepts
1.  **Dynamic Registration:** Workers are not statically known. They connect to the broker and announce their presence (`READY` signal).
2.  **LRU Routing:** The broker maintains a queue of available workers. When a worker finishes a task, it is added to the back of the queue. New tasks are assigned to the worker at the front.
3.  **Manual Routing:** Unlike `DEALER` (which round-robins automatically), the `ROUTER` socket requires the application to specify the destination identity for every message. This gives the application full control over the load-balancing algorithm.

## Execution
Run `run.ps1`.
- The Broker starts.
- 3 AV Engines connect and send `READY`.
- The Uploader sends 5 concurrent file requests.
- The Broker dispatches them to available workers.
- You will see different "AV-x" engines handling the files.
