# Dashboard build stage
FROM node:20-alpine AS dashboard-builder

WORKDIR /dashboard

# Copy dashboard files
COPY dashboard/package*.json ./
RUN npm install --legacy-peer-deps

COPY dashboard/ ./
RUN npm run build

# Go build stage
FROM golang:1.24-alpine AS builder

WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Copy built dashboard from previous stage to the location expected by embed directive
COPY --from=dashboard-builder /dashboard/out ./internal/dashboard/out

# Build the binary with embedded dashboard
RUN CGO_ENABLED=0 GOOS=linux go build -tags embed_dashboard -o netkit ./cmd/netkit

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/netkit .

# Expose ports (proxy, admin, dashboard)
EXPOSE 8080 8081 3000

# Run the proxy server with embedded dashboard
CMD ["./netkit", "serve", "--port", "8080", "--admin-port", "8081", "--dashboard", "--dashboard-port", "3000", "--log-level", "info"]
