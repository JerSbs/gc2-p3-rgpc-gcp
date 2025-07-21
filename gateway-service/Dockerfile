# STAGE 1: BUILD
FROM golang:1.24 AS builder

WORKDIR /app

# Download Go module dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy entire project
COPY . .

# Build static binary from cmd/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o gateway-server ./cmd

# STAGE 2: RUNTIME
FROM debian:bullseye-slim

WORKDIR /app

# Copy binary only
COPY --from=builder /app/gateway-server .

# Expose REST gateway port
EXPOSE 8082

# Run gateway server
CMD ["./gateway-server"]
