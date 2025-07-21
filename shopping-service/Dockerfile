# STAGE 1: BUILD
FROM golang:1.24 AS builder

WORKDIR /app

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build static REST binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o rest-server ./cmd

# STAGE 2: RUNTIME
FROM debian:bullseye-slim

WORKDIR /app

# Copy REST binary only
COPY --from=builder /app/rest-server .

# Expose REST port
EXPOSE 8080

# Run REST server
CMD ["./rest-server"]