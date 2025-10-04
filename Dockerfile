# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o netkit ./cmd/netkit

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/netkit .

# Expose ports
EXPOSE 8080 8081

# Run the proxy server
CMD ["./netkit", "serve", "--port", "8080", "--admin-port", "8081", "--log-level", "info"]
