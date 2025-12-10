# Build stage
FROM golang:1.23.4-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

# Copy source code
COPY . .

# Build the application (use BuildKit caches, let platform decide)
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-w -s" -o server ./cmd/server

# Runtime stage
FROM gcr.io/distroless/base-debian12

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/server .

# Copy .env.example as fallback (can be overridden by volume mount)
COPY .env.example .env.example

# Expose port
EXPOSE 8000

# Use non-root user
USER nonroot:nonroot

# Run the application
ENTRYPOINT ["./server"]
