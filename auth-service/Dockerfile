# STAGE 1: BUILD
FROM golang:1.24 as builder

WORKDIR /app

# Set GOPROXY langsung ke sumber
ENV GOPROXY=direct
ENV GOFLAGS=-mod=vendor


# Download dependencies
COPY go.mod go.sum ./
COPY vendor/ ./vendor/

# Copy source code
COPY . .

# Build static gRPC binary from cmd/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o grpc-server ./cmd

# STAGE 2: RUNTIME
FROM debian:bullseye-slim

WORKDIR /app

# Copy gRPC binary only
COPY --from=builder /app/grpc-server .

# Expose gRPC port
EXPOSE 50052

# Run gRPC server
CMD ["./grpc-server"]
