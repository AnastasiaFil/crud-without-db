# Stage 1: Build the Go application
FROM public.ecr.aws/docker/library/golang:1.23 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Build from the correct path based on your Makefile
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

# Stage 2: Create the runtime image
FROM public.ecr.aws/amazonlinux/amazonlinux:2

# Install necessary dependencies including CA certificates for SSL connections
RUN yum update -y && \
    yum install -y ca-certificates && \
    yum clean all

WORKDIR /app
COPY --from=builder /app/main .

EXPOSE 3000
CMD ["./main"]