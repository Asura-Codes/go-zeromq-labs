# Lab 23: Decentralized Cluster Gossip (PUB-SUB Pattern)

## Scenario
In a peer-to-peer network without a central broker, nodes need a way to share their status and discover each other. This lab implements a simple gossip mechanism where each node is both a publisher of its own state and a subscriber to its peers' states.

## Pattern: Decentralized PUB-SUB
- Each node creates a `PUB` socket to broadcast its status.
- Each node creates a `SUB` socket to listen for status updates from other nodes.
- By connecting to known peers, a mesh network is formed where state information propagates across the cluster.

## Key Features
- **Pure-Go Implementation:** No CGO or external C libraries required.
- **Node Self-Reporting:** Nodes periodically broadcast an "ALIVE" status message.
- **Local State Tracking:** Each node maintains an in-memory map of the cluster's current state.

## How to Run
1. Start three nodes in separate terminals:
   ```powershell
   # Node A
   go run ./cmd/gossip_node --pub tcp://*:6666 --name Node-A
   
   # Node B (connects to A)
   go run ./cmd/gossip_node --pub tcp://*:6667 --name Node-B tcp://localhost:6666
   
   # Node C (connects to B)
   go run ./cmd/gossip_node --pub tcp://*:6668 --name Node-C tcp://localhost:6667
   ```
2. Observe how `Node-A` sees updates from `Node-B`, and `Node-C` sees updates from `Node-B`. To see all nodes, ensure they are cross-connected or use a discovery mechanism.