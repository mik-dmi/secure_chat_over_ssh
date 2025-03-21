# Step 1: Build the Go application using the golang image
FROM golang:1.22.6 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum first to take advantage of Docker cache
COPY go.mod go.sum ./

# Download dependencies
RUN go mod tidy

# Copy all the necessary Go source files into the container
COPY . .

# Build the Go application
RUN go build -o ssh-chat-server ./cmd

# Step 2: Use the base Ubuntu image to run the Go application
FROM ubuntu:latest

# Install OpenSSH server (optional, for debugging purposes)
RUN apt-get update && apt-get install -y openssh-server

# Copy over private/public keys for SSH (if needed)
WORKDIR /app
RUN mkdir /keys
COPY ./keys/dockerkey /keys

# Setup authorized keys for SSH (for Go SSH server)
COPY ./keys/dockerkey.pub /root/.ssh/authorized_keys
RUN chmod 700 /root/.ssh && chmod 600 /root/.ssh/authorized_keys

# Copy the compiled Go application from the builder stage
COPY --from=builder /app/ssh-chat-server /usr/local/bin/ssh-chat-server

# Create necessary directories for SSH server (optional)
RUN mkdir /var/run/sshd

# Expose port 2222 for SSH
EXPOSE 2222

# Start the Go SSH chat server
CMD ["/usr/local/bin/ssh-chat-server"]
