build-server:
	@echo "Building Server..."
	@go build -o main cmd/server/main.go

build-client: 
	@echo "Building Client..."
	@go build -o main cmd/client/main.go

# Run just the API
run-server:
	@go run cmd/server/main.go

# Run just the frontend
run-client:
	@go run cmd/client/main.go

