# Multi-stage Dockerfile for Enterprise Shield

# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /build

# Install build dependencies
RUN apk add --no-cache git make

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s -X main.version=$(git describe --tags --always)" \
    -o enterprise-shield \
    ./cmd/plugin

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1000 shield && \
    adduser -D -u 1000 -G shield shield

# Create directories
RUN mkdir -p /home/shield/.opencode/config \
             /home/shield/.opencode/logs \
             /home/shield/.opencode/plugins && \
    chown -R shield:shield /home/shield/.opencode

# Copy binary from builder
COPY --from=builder /build/enterprise-shield /usr/local/bin/enterprise-shield

# Copy default config
COPY --chown=shield:shield config/default.yaml /home/shield/.opencode/config/enterprise-shield.yaml

# Switch to non-root user
USER shield
WORKDIR /home/shield

# Expose metrics port (if needed)
# EXPOSE 9090

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ["/usr/local/bin/enterprise-shield", "version"]

# Default command
ENTRYPOINT ["/usr/local/bin/enterprise-shield"]
CMD ["serve"]

