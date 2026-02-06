# Lab 14: Distributed Hash Table Ring (P2P Architecture)

## Scenario
A decentralized key-value store where data is distributed across a ring of nodes using consistent hashing.

## Pattern
P2P Ring / Distributed Hash Table (DHT).

## Key Concepts
- **Consistent Hashing:** Using `crc32` to map both node addresses and keys to a 32-bit circular space.
- **Range-based Responsibility:** Nodes handle keys in the range `(PredecessorID, NodeID]`.
- **Request Forwarding:** If a node is not responsible for a key, it proxies the request to its successor in the ring.
- **Self-Healing:** (Conceptual) Nodes form a ring and maintain connectivity.

## Deliverables
- `dht_node.exe`: A member of the DHT ring.
- `client_put_get.exe`: A test client that interacts with the ring via any node.
