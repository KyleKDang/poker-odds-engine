.PHONY: build run test clean docker-build docker-run docker-stop

# Build the application
build:
	go build -o bin/poker-odds-engine cmd/server/main.go

# Run the application
run:
	go run cmd/server/main.go

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Format code
fmt:
	go fmt ./...

# Docker commands
docker-build:
	docker build -t poker-odds-engine .

docker-run:
	docker-compose up -d

docker-stop:
	docker-compose down

docker-logs:
	docker-compose logs -f
	