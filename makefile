.PHONY: build run test swag clean dev dev-down local

build:
	go mod download && go build -o crud-without-db ./cmd/main.go

# Run with Docker Compose (includes PostgreSQL)
run: build
	docker-compose up --remove-orphans

# Run development environment with Docker Compose
dev:
	docker-compose -f docker-compose.yml up --build --remove-orphans

# Stop development environment
dev-down:
	docker-compose -f docker-compose.yml down

# Run locally (requires local PostgreSQL)
local: build
	./crud-without-db

# Run PostgreSQL only for local development
db-only:
	docker run --name crud-postgres -d \
		-e POSTGRES_DB=crud_db \
		-e POSTGRES_USER=postgres \
		-e POSTGRES_PASSWORD=qwerty123 \
		-p 5432:5432 \
		postgres:17-alpine

# Stop and remove PostgreSQL container
db-stop:
	docker stop crud-postgres || true
	docker rm crud-postgres || true

test:
	go test -v ./...

swag:
	swag init -g ./cmd/main.go

clean:
	rm -f crud-without-db
	docker-compose down --volumes --remove-orphans 2>/dev/null || true
	docker-compose -f docker-compose.yml down --volumes --remove-orphans 2>/dev/null || true

# Install dependencies
deps:
	go mod download
	go mod tidy

# Run with live reload (requires air: go install github.com/cosmtrek/air@latest)
live:
	air