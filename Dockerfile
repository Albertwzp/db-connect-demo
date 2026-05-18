# Multi-stage Dockerfile for db-connect-demo

# Build stage
FROM golang:1.26-alpine AS builder

WORKDIR /workspace

# Install dependencies
RUN apk add --no-cache git gcc musl-dev sqlite-dev

# Copy source code
COPY . .

# Download dependencies
RUN go mod download

# Build single binary
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -a -ldflags="-w -s" -o /workspace/db-connect-demo .

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates

# Copy built binary
COPY --from=builder /workspace/db-connect-demo /app/db-connect-demo

EXPOSE 8080

ENTRYPOINT ["/app/db-connect-demo"]