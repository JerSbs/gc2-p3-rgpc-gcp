# STAGE 1: BUILD
FROM golang:1.24 as builder

WORKDIR /app

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build static gRPC binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o grpc-server

# STAGE 2: RUNTIME
FROM debian:bullseye-slim

WORKDIR /app

# Copy gRPC binary only
COPY --from=builder /app/grpc-server .

# Expose gRPC port
EXPOSE 50051

# Run gRPC server
CMD ["./grpc-server"]
