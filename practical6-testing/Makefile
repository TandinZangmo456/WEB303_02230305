.PHONY: help test-unit test-integration test-e2e docker-up docker-down

help:
	@echo "Available commands:"
	@echo "  make test-unit          - Run all unit tests"
	@echo "  make test-unit-user     - Run user service unit tests"
	@echo "  make test-unit-menu     - Run menu service unit tests"
	@echo "  make test-unit-order    - Run order service unit tests"
	@echo "  make test-integration   - Run integration tests"
	@echo "  make test-e2e           - Run E2E tests (requires services running)"
	@echo "  make test-all           - Run all tests"
	@echo "  make docker-up          - Start all services with Docker"
	@echo "  make docker-down        - Stop all services"
	@echo "  make docker-logs        - Show Docker logs"
	@echo "  make proto-generate     - Generate proto files"

# Unit Tests
test-unit:
	@echo "=== Running Unit Tests ==="
	@cd user-service && go test ./grpc/... -v
	@cd menu-service && go test ./grpc/... -v
	@cd order-service && go test ./grpc/... -v

test-unit-user:
	@echo "=== User Service Unit Tests ==="
	@cd user-service && go test ./grpc/... -v

test-unit-menu:
	@echo "=== Menu Service Unit Tests ==="
	@cd menu-service && go test ./grpc/... -v

test-unit-order:
	@echo "=== Order Service Unit Tests ==="
	@cd order-service && go test ./grpc/... -v

# Integration Tests
test-integration:
	@echo "=== Running Integration Tests ==="
	@cd tests/integration && go test -v

# E2E Tests
test-e2e:
	@echo "=== Running E2E Tests ==="
	@cd tests/e2e && go test -v

# Run all tests
test-all: test-unit test-integration
	@echo "=== All Tests Completed ==="

# Docker commands
docker-up:
	docker compose up -d

docker-down:
	docker compose down

docker-logs:
	docker compose logs -f

# Proto generation
proto-generate:
	./generate_proto.sh

# Coverage
test-coverage:
	@echo "=== Generating Coverage Reports ==="
	@cd user-service && go test ./grpc/... -coverprofile=coverage.out
	@cd user-service && go tool cover -html=coverage.out -o coverage.html
	@cd menu-service && go test ./grpc/... -coverprofile=coverage.out
	@cd menu-service && go tool cover -html=coverage.out -o coverage.html
	@cd order-service && go test ./grpc/... -coverprofile=coverage.out
	@cd order-service && go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage reports generated in each service directory"
