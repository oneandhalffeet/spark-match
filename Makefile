.PHONY: run test build docker-build docker-run seed clean

# Run the application
run:
	go run main.go

# Run tests
test:
	go test ./... -v

# Build the application
build:
	go build -o matchmaking-engine .

# Build Docker image
docker-build:
	docker build -t matchmaking-engine .

# Run Docker container
docker-run:
	docker run -p 8080:8080 matchmaking-engine

# Seed data (requires running server)
seed:
	curl http://localhost:8080/seed?count=1000

# Clean build artifacts
clean:
	rm -f matchmaking-engine
	go clean

# Install dependencies
deps:
	go mod download

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run