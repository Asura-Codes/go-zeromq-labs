# Build Stage
FROM golang:1.25 AS builder

WORKDIR /app

# Install ZeroMQ development libraries (required for pebbe/zmq4)
RUN apt-get update && apt-get install -y libzmq3-dev pkg-config

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build all binaries
RUN go build -o keygen ./cmd/keygen
RUN go build -o secure_server ./cmd/secure_server
RUN go build -o secure_client ./cmd/secure_client

# Runtime Stage
FROM debian:bookworm-slim

WORKDIR /app

# Install ZeroMQ runtime library
RUN apt-get update && apt-get install -y libzmq5 && rm -rf /var/lib/apt/lists/*

# Copy binaries
COPY --from=builder /app/keygen /usr/local/bin/
COPY --from=builder /app/secure_server /usr/local/bin/
COPY --from=builder /app/secure_client /usr/local/bin/

# Default entrypoint (can be overridden)
CMD ["secure_server"]
