# Use the official Golang image as a base
FROM golang:1.23-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod ./

# Copy the rest of the application code
COPY . .

# Download dependencies
RUN go mod tidy

# Build the Go application
RUN go build -o main .

# Use a minimal base image for the final build
FROM alpine:3.18

# Set the working directory inside the container
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/main .
COPY --from=builder /app/templates /app/templates
COPY --from=builder /app/sql /app/sql

# Expose the port on which the service will run
EXPOSE 8082

# Command to run the binary
CMD ["./main"]
