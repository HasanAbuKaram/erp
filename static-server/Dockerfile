# Use the official Golang image as a base
FROM golang:1.23-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o static-server .

# Use a minimal base image for the final build
FROM alpine:3.18

# Set the working directory inside the container
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/static-server .

# Copy static files
COPY static /app/static

# Expose the port on which the service will run
EXPOSE 8081

# Command to run the binary
CMD ["./static-server"]
