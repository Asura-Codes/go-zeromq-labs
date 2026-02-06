FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build argument for which binary to build
ARG SERVICE_NAME
RUN go build -o /app/service ./cmd/${SERVICE_NAME}

# Final stage
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/service .

CMD ["./service"]
