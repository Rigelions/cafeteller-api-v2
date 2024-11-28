# Stage 1: Build
FROM golang:1.22 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the application source code
COPY . .

# Download dependencies
RUN go mod download

# Build the Go application
RUN go build -o main .

# Stage 2: Run
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /root/

# Copy the built binary from the builder stage
COPY --from=builder /app/main .

# Expose the application port (update if your app uses a different port)
EXPOSE 8080

# Command to run the application
CMD ["./main"]