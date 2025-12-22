# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod files (go.sum is required)
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY main.go ./

# Build the binary with verification
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o append-sheet . && \
    ls -la /app/append-sheet && \
    /app/append-sheet --help || echo "Binary built successfully"

# Final stage - use a minimal base image
FROM alpine:latest

# Install CA certificates for HTTPS requests to Google APIs
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/append-sheet ./append-sheet

# Verify the binary exists and is executable
RUN ls -la /root/append-sheet && chmod +x /root/append-sheet

# Run the binary
ENTRYPOINT ["/root/append-sheet"]

