# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN make build

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache git ca-certificates

# Create non-root user
RUN adduser -D -s /bin/sh gitflow

# Copy binary from builder
COPY --from=builder /build/build/gitflow-tui /usr/local/bin/gitflow-tui

# Set ownership
RUN chown gitflow:gitflow /usr/local/bin/gitflow-tui

# Switch to non-root user
USER gitflow

# Set working directory
WORKDIR /workspace

# Entry point
ENTRYPOINT ["gitflow-tui"]
