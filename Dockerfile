# Stage 1: Build the application
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./
# Download dependencies
RUN go mod download

# Copy the entire application
COPY . .

# Build the Go application
RUN go build -o flatchecker-scheduler.exe .

# Stage 2: Create a minimal image
FROM alpine:latest

# Copy the binary from the builder stage
COPY --from=builder /app/flatchecker-scheduler.exe .
COPY prod-config.txt .

# Command to run the application
CMD ["./flatchecker-scheduler.exe", "prod-config.txt"]