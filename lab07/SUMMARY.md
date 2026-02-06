# Lab 07: Global Policy Synchronization (Clone Pattern)

## Description
This lab implements the **Clone Pattern**, a distributed state synchronization mechanism. It allows nodes to fetch a full point-in-time snapshot of a dataset and then maintain consistency via a real-time stream of updates (deltas).
- **Policy Master:** Maintains the authoritative state of firewall rules. It broadcasts updates on a PUB socket and serves full snapshots on a ROUTER socket.
- **Firewall Node:** Upon startup, it fetches the current state from the Master and then applies incremental updates, ensuring it never misses a change or applies an out-of-order update.

## Architecture
- **Protocol:** TCP / JSON
- **Socket Types:**
  - Master Snapshot: `ROUTER`
  - Master Updates: `PUB`
  - Node Snapshot Client: `DEALER`
  - Node Update Client: `SUB`
- **Pattern:** Clone Pattern (v1).

## Key Concepts
1.  **State vs. Stream:** The problem with pure PUB-SUB is that new subscribers miss all previous messages. The Clone pattern solves this by adding a "State" (Snapshot) side-channel.
2.  **Sequence Numbers:** Every update has a monotonic sequence number. The Node uses this to ensure that it only applies updates that are newer than its current local state.
3.  **Idempotency:** In this lab, updates are KV-overwrites. Applying the same update twice doesn't hurt, and sequence numbers prevent applying old updates over new ones.

## Execution
Run `run.ps1`.
- The Master starts and begins generating random firewall rules.
- After 5 seconds, a Node joins.
- You will see the Node fetch several "Snapshot items" (the rules generated while it was away).
- Then, you will see the Node transition to "real-time updates" as the Master continues to broadcast new rules.
