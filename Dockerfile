# Step 1: Use the official Golang image to build the Go application
FROM golang:1.20-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Initialize go.mod (without specifying the Go version directly)
RUN go mod init my-go-app

# Get the PostgreSQL package
RUN go get github.com/lib/pq

# Copy the Go application source code
COPY . .

# Build the Go application
RUN go build -o my-go-app

# Step 2: Use a smaller base image for the final container
FROM alpine:latest

# Set the working directory in the container
WORKDIR /app

# Install PostgreSQL client to interact with the database
RUN apk --no-cache add postgresql-client

# Copy the compiled Go binary from the builder stage
COPY --from=builder /app/my-go-app .

# Expose the application port (if needed)
EXPOSE 8080

# Command to run the Go application
CMD ["./my-go-app"]

