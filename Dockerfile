# Stage 1: Build the Go application
FROM golang:1.23 AS builder

WORKDIR /app

# Copy go.mod and go.sum to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the Go binary with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o crud-without-db ./cmd/main.go

# Stage 2: Create the runtime image
FROM amazonlinux:2

# Install necessary dependencies including CA certificates for SSL connections
RUN yum update -y && \
    yum install -y shadow-utils curl ca-certificates && \
    yum clean all && \
    rm -rf /var/cache/yum

# Create directories for the application and logs
RUN mkdir -p /usr/local/bin /var/log/crud-without-db
RUN chmod 755 /usr/local/bin /var/log/crud-without-db

# Copy the compiled binary from the builder stage
COPY --from=builder /app/crud-without-db /usr/local/bin/crud-without-db

# Create a non-root user for security
RUN useradd -r -s /bin/false appuser
RUN chown appuser:appuser /usr/local/bin/crud-without-db /var/log/crud-without-db

# Switch to non-root user
USER appuser

# Expose the application port
EXPOSE 3000

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:3000/health || exit 1

# Command to run the application
CMD ["/usr/local/bin/crud-without-db"]