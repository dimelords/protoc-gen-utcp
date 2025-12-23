.PHONY: build install test clean examples help

# Build configuration
BINARY_NAME=protoc-gen-utcp
INSTALL_PATH=$(shell go env GOPATH)/bin
BUILD_DIR=bin

help: ## Show this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@awk 'BEGIN {FS = ":.*##"; printf "\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  %-15s %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

build: ## Build the protoc plugin
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/protoc-gen-utcp
	@echo "✓ Built to $(BUILD_DIR)/$(BINARY_NAME)"

install: build ## Install the plugin to GOPATH/bin
	@echo "Installing $(BINARY_NAME) to $(INSTALL_PATH)..."
	@cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "✓ Installed to $(INSTALL_PATH)/$(BINARY_NAME)"

test: ## Run tests
	@echo "Running tests..."
	@go test -v -race -cover ./...
	@echo "✓ Tests passed"

test-coverage: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	@go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...
	@go tool cover -html=coverage.txt -o coverage.html
	@echo "✓ Coverage report generated: coverage.html"

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.txt coverage.html
	@find examples -name "*.utcp.json" -delete
	@find examples -name "*.pb.go" -delete
	@echo "✓ Cleaned"

examples: install ## Generate UTCP for example protos
	@echo "Generating UTCP for examples..."
	@echo ""
	@echo "→ Simple example (generic HTTP):"
	@protoc \
		--utcp_out=. \
		--utcp_opt=base_url=https://api.example.com \
		--utcp_opt=provider_type=http \
		--utcp_opt=auth_type=bearer \
		examples/simple/service.proto
	@echo "  Generated: examples/simple/service.utcp.json"
	@echo ""
	@echo "→ Twirp example:"
	@protoc \
		--utcp_out=. \
		--utcp_opt=base_url=https://teamleader.se \
		--utcp_opt=provider_type=http \
		--utcp_opt=auth_type=bearer \
		examples/twirp/documents.proto
	@echo "  Generated: examples/twirp/documents.utcp.json"
	@echo ""
	@echo "✓ Examples generated"

lint: ## Run linters
	@echo "Running linters..."
	@go vet ./...
	@echo "✓ Linting passed"

fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...
	@echo "✓ Code formatted"

mod-tidy: ## Tidy go.mod
	@echo "Tidying go.mod..."
	@go mod tidy
	@echo "✓ go.mod tidied"

version: ## Show version
	@$(BUILD_DIR)/$(BINARY_NAME) -version || echo "Build first: make build"

.DEFAULT_GOAL := help
