# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o server ./cmd/server

# Runtime stage
FROM gcr.io/distroless/base-debian12

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/server .

# Copy .env.example as fallback (can be overridden by volume mount)
COPY .env.example .env.example

# Expose port
EXPOSE 8080

# Use non-root user
USER nonroot:nonroot

# Run the application
ENTRYPOINT ["./server"]
