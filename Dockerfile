# syntax=docker/dockerfile:1

# --- Build Stage ---
FROM golang:1.27-alpine AS builder

# Install git and ca-certificates (needed if fetching private modules or verifying TLS)
RUN apk add --no-cache git ca-certificates

WORKDIR /src

# Download dependencies first to leverage Docker layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build a statically linked binary
# CGO_ENABLED=0 ensures no dynamic C libraries are required at runtime
# -ldflags="-s -w" strips debug symbols to minimize binary size
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/server .

# --- Runtime Stage ---
FROM gcr.io/distroless/static-debian12:nonroot

# Set Gin to release mode to disable debug logging and improve performance
ENV GIN_MODE=release
WORKDIR /app

# Copy the compiled binary and assets from the builder stage
COPY --from=builder /app/server /app/server
COPY --from=builder /src/public /app/public
COPY --from=builder /src/templates /app/templates

# Document the default port Gin listens on
EXPOSE 8080

# Run as the non-root user provided by distroless
USER nonroot:nonroot

ENTRYPOINT ["/app/server"]

