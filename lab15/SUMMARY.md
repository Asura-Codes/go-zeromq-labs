# Lab 15: MapReduce Cluster (Scatter-Gather)

## Scenario
Parallel processing of a large dataset (text chunks) to perform a distributed word count.

## Pattern
Scatter-Gather (Parallel Pipeline).

## Key Concepts
- **Ventilator (Master):** Splits the input data into chunks and pushes them to Map workers.
- **Map Phase:** Workers process text chunks into word-count pairs.
- **Reduce Phase:** Workers aggregate counts for specific words.
- **Sink (Master):** Collects the final aggregated results from the Reduce workers.
- **Control Flow:** Uses a PUB-SUB channel to broadcast a `STOP` signal to finalize reduction.

## Deliverables
- `mr_master.exe`: Orchestrates the cluster and collects final results.
- `map_worker.exe`: Performs word frequency mapping.
- `reduce_worker.exe`: Performs count aggregation.
