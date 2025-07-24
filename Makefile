.PHONY: build clean run docker-build docker-run tidy test

# Binary name
BINARY_NAME=shipwright-build-mcp-server

# Build the binary
build:
	go build -o $(BINARY_NAME) .

# Clean build artifacts
clean:
	go clean
	rm -f $(BINARY_NAME)

# Run the server
run: build
	./$(BINARY_NAME)

# Run from source
run-dev:
	go run main.go

# Tidy dependencies
tidy:
	go mod tidy

# Run tests
test:
	go test -v ./...

# Build Docker image
docker-build:
	docker build -t $(BINARY_NAME) .

# Run Docker container
docker-run: docker-build
	docker run -i $(BINARY_NAME)

# Development helpers
fmt:
	go fmt ./...

vet:
	go vet ./...

# Install dependencies
deps:
	go mod download

# Build for multiple platforms
build-all: build-linux build-darwin build-windows

build-linux:
	GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME)-linux-amd64 .

build-darwin:
	GOOS=darwin GOARCH=amd64 go build -o $(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build -o $(BINARY_NAME)-darwin-arm64 .

build-windows:
	GOOS=windows GOARCH=amd64 go build -o $(BINARY_NAME)-windows-amd64.exe .

# Default target
all: tidy fmt vet test build 