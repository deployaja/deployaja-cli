FROM golang:1.20-alpine AS builder

WORKDIR /app

# Install git (required for go mod)
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the CLI binary
RUN go build -o aja main.go

# Final image
FROM alpine:latest

WORKDIR /app

# Copy the built binary from builder
COPY --from=builder /app/aja /usr/local/bin/aja

# Set entrypoint
ENTRYPOINT ["aja"]
