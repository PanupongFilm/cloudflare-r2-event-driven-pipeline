.PHONY: help run seed seed-reset build clean

help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

run: ## Run the API server
	go run cmd/api/main.go

seed: ## Run database seeds
	go run cmd/seed/main.go

seed-reset: ## Reset database and run seeds
	go run cmd/seed/main.go -reset

build: ## Build all binaries
	@echo "Building API server..."
	@go build -o bin/api cmd/api/main.go
	@echo "Building seed command..."
	@go build -o bin/seed cmd/seed/main.go
	@echo "Build completed!"

clean: ## Remove built binaries
	@rm -rf bin/
	@echo "Cleaned!"

dev: ## Run server with auto-reload (requires air)
	air

install-tools: ## Install development tools
	go install github.com/cosmtrek/air@latest
