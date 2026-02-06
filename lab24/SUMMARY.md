# Lab 24: Autonomous Secure Beaconing (UDP + REP/REQ)

## Scenario
In industrial environments, devices must discover each other automatically on a LAN without a central registry. This lab demonstrates how to use UDP broadcasts for discovery and ZeroMQ for the subsequent communication.

## Pattern: UDP Discovery + ZeroMQ
- **UDP Beaconing:** The device periodically sends a UDP broadcast containing its service port.
- **UDP Scanning:** The scanner listens for UDP broadcasts to learn the device's IP and port.
- **Secure Communication:** Once discovered, a `REQ/REP` connection is established. While this demo uses a standard connection, `zmq4` supports CURVE for full encryption.

## Key Features
- **Pure-Go Implementation:** Uses standard `net` package for UDP and `zmq4` for messaging.
- **Zero-Config Discovery:** No hardcoded IP addresses are needed for the scanner to find the device.
- **Cross-Platform Compatibility:** Works on any system with a standard network stack without needing C libraries.

## How to Run
1. Start the secure device:
   ```powershell
   go run ./cmd/secure_device
   ```
2. Start the network scanner:
   ```powershell
   go run ./cmd/network_scanner
   ```
3. Observe how the scanner detects the device's IP and port via UDP and successfully sends a request.