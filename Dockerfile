FROM golang:1.20-alpine AS builder

# Install build dependencies
RUN apk add --no-cache build-base libwebp-dev

# Set the working directory
WORKDIR /app

# Copy the Go modules manifests
COPY go.mod go.sum ./

# Download dependencies
RUN GO111MODULE=on go mod download

# Copy the entire source code
COPY . .

# Set the command directory
WORKDIR /app

# Build the Go application
RUN GO111MODULE=on go build -o /main ./cmd/main.go

# Use a minimal image for the runtime
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache libwebp

# Set the working directory
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /main .

# Expose the application port
EXPOSE 8080

# Command to run the executable
CMD ["./main"]