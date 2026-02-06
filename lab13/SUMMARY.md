# Lab 13: Persistent Message Queue (Titanic Pattern)

## Scenario
Implement a "store-and-forward" mechanism where client requests are persisted to disk and processed even if the target service is temporarily offline.

## Pattern
Titanic (Persistent Queuing) + Majordomo (MDP).

## Key Concepts
- **Persistence:** Requests and replies are stored in `titanic_data/` using unique UUIDs.
- **Asynchronous Processing:** A background loop in the Titanic Broker polls the storage and dispatches pending requests to Majordomo workers.
- **Polling:** Clients submit a request and poll for the result using the assigned request ID.

## Deliverables
- `titanic_broker.exe`: The primary gateway that stores requests and manages dispatch.
- `storage_service.exe`: A dedicated service for file-based CRUD operations on messages.
- `patient_client.exe`: A client that submits tasks and waits patiently for results.
- `mock_majordomo.exe`: A simplified REP-based worker for testing.
