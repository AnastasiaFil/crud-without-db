# Stage 1: Build the Go application
FROM golang:1.23 AS builder

WORKDIR /app

# Copy go.mod and go.sum to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -o crud-without-db ./cmd/main.go

# Stage 2: Create the runtime image
FROM amazonlinux:2

# Install necessary dependencies
RUN yum update -y && yum install -y shadow-utils curl && yum clean all

# Create directories for the application and logs
RUN mkdir -p /usr/local/bin /var/log/crud-without-db
RUN chmod 755 /usr/local/bin /var/log/crud-without-db

# Copy the compiled binary from the builder stage
COPY --from=builder /app/crud-without-db /usr/local/bin/crud-without-db

# Expose the application port
EXPOSE 3000

# Command to run the application
CMD ["/usr/local/bin/crud-without-db"]