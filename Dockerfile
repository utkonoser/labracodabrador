# Build stage
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make gcc musl-dev linux-headers

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o labracodabrador ./cmd/labracodabrador

# Runtime stage
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache ca-certificates curl

# Create app user
RUN addgroup -g 1000 app && \
    adduser -D -u 1000 -G app app

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/labracodabrador .
COPY --from=builder /build/config.yaml .
COPY --from=builder /build/genesis.json .

# Create data directory
RUN mkdir -p /app/data && \
    chown -R app:app /app

# Switch to app user
USER app

# Expose ports
EXPOSE 8545 8546 30303

# Run the application
ENTRYPOINT ["./labracodabrador"]
CMD ["-config", "config.yaml", "-genesis", "genesis.json"]

