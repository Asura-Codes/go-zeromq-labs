# Lab 17: Distributed Consensus (Raft Prototype)

## Scenario
A cluster of critical database nodes must agree on a leader to handle operations. This implementation uses a robust decentralized mesh to elect a leader and replicate simple commands.

## Pattern: PUB-SUB Mesh (Gossip Consensus)
This lab implements a high-reliability version of the Raft **Leader Election** phase using a full-mesh PUB-SUB architecture.

### Architecture
- **Topology:** Full Mesh (Every node is both a Publisher and a Subscriber).
- **Transport:** 
    - **PUB:** Broadcasts state changes (Vote Requests, Votes, Heartbeats).
    - **SUB:** Listens to all peers for cluster-wide events.
    - **Gateway (PULL):** Entry point for external client commands.
- **States:**
    - **Follower:** Listens for leader heartbeats; grants votes to valid candidates.
    - **Candidate:** Broadcasts election requests and collects public votes.
    - **Leader:** Heartbeats the cluster to maintain authority and broadcasts client commands.

## Key Concepts
- **Decentralized Messaging:** No central broker or point-to-point locks. The state is gossiped across the mesh.
- **Randomized Election Timers:** Prevents split votes by staggering candidate starts.
- **Quorum-based Authority:** Leaders are only elected when a majority of the cluster publicly broadcasts their support.

## Running the Lab
1. Run `./run.ps1`.
2. Observe the "Election" logs as nodes negotiate the leader for the current term.
3. Once a leader is elected, the `client_proposer` sends a command to Node 1's gateway.
4. Observe the leader broadcasting the command and other nodes receiving it.