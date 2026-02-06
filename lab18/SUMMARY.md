# Lab 18: Zero-Copy Video Streaming

## Scenario
A live video surveillance system needs to distribute 4K video frames to multiple analytics engines with minimal latency.

## Pattern: Pub-Sub with Multipart Messages
This lab demonstrates high-throughput data transmission using ZeroMQ's efficient message handling.

### Architecture
- **Camera Feed (Publisher):** Generates simulated 1080p raw frames (~6MB per frame) and broadcasts them at 30 FPS.
- **Analytics Engine (Subscriber):** Consumes the frame stream and calculates throughput.

## Key Concepts
- **Multipart Messages:** Each message consists of a header frame (ID) and a payload frame (Raw Data).
- **Buffer Management:** Reusing byte slices to minimize allocation overhead.
- **Throughput:** Capable of saturating local network bandwidth (multiple hundreds of MB/s).

## Running the Lab
1. Run `./run.ps1`.
2. Observe the analytics engine logs to see the MB/s throughput.
