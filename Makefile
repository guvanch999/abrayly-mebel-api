BINARY := abrayly-mebel-backend
SWAG_VERSION := v1.16.4

.PHONY: run build test cover vet fmt swagger docker-up docker-down

run: ## Run the API locally (reads .env)
	go run .

build: ## Build the binary into ./bin
	go build -o bin/$(BINARY) .

test: ## Run unit tests (no database required)
	go test -race -count=1 ./...

cover: ## Run tests with a coverage summary
	go test -count=1 -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out | tail -n 1

vet: ## Static analysis
	go vet ./...

fmt: ## Format the code
	gofmt -s -w .

swagger: ## Regenerate docs/ from the handler annotations
	go run github.com/swaggo/swag/cmd/swag@$(SWAG_VERSION) init -g main.go -o docs

docker-up: ## Start PostgreSQL + API with Docker Compose
	docker compose up --build

docker-down: ## Stop the stack (keeps volumes)
	docker compose down
