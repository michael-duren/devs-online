build:
	@echo "Building CLI..."
	@go build -o dol cmd/cli/main.go

# Run the websocket server
serve:
	@go run cmd/cli/main.go serve

# Run just the API
serve-ext:
	@go run cmd/cli/main.go serve --external

# Run just the frontend
chat:
	@go run cmd/cli/main.go chat

