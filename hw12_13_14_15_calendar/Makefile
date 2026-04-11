# Build variables
BIN_DIR := ./bin
CALENDAR_BIN := $(BIN_DIR)/calendar
SCHEDULER_BIN := $(BIN_DIR)/scheduler
STORER_BIN := $(BIN_DIR)/storer

# Docker variables
DOCKER_COMPOSE := docker-compose -f deployments/docker-compose.yaml

# Version variables
GIT_HASH := $(shell git log --format="%h" -n 1 2>/dev/null || echo "unknown")
BUILD_DATE := $(shell date -u +%Y-%m-%dT%H:%M:%S)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(BUILD_DATE) -X main.gitHash=$(GIT_HASH)

## build: Build all binaries
.PHONY: build
build: clean
	@echo "Building calendar API..."
	@mkdir -p $(BIN_DIR)
	@go build -ldflags "$(LDFLAGS)" -o $(CALENDAR_BIN) ./cmd/calendar
	@echo "Building scheduler..."
	@go build -ldflags "$(LDFLAGS)" -o $(SCHEDULER_BIN) ./cmd/scheduler
	@echo "Building storer..."
	@go build -ldflags "$(LDFLAGS)" -o $(STORER_BIN) ./cmd/storer
	@echo "Build complete! Binaries are in $(BIN_DIR)/"

## run-calendar: Run calendar API service
.PHONY: run-calendar
run-calendar: build
	$(CALENDAR_BIN) --config=configs/calendar.yaml

## run-scheduler: Run scheduler service
.PHONY: run-scheduler
run-scheduler: build
	$(SCHEDULER_BIN) --config=configs/scheduler.yaml

## run-storer: Run storer service
.PHONY: run-storer
run-storer: build
	$(STORER_BIN) --config=configs/storer.yaml

## clean: Remove build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BIN_DIR)
	@rm -f calendar.exe
	@echo "Clean complete"

## test: Run tests
.PHONY: test
test:
	@echo "Running tests..."
	@go test -race -v ./internal/...

## lint: Run linter
.PHONY: lint
lint:
	@echo "Running linter..."
	@go vet ./...
	@gofmt -l . | (! grep .)

## docker-build: Build Docker images
.PHONY: docker-build
docker-build:
	@echo "Building Docker images..."
	@$(DOCKER_COMPOSE) build
	@echo "Docker images built"

## docker-up: Start all services with Docker Compose
.PHONY: docker-up
docker-up:
	@echo "Starting services..."
	@$(DOCKER_COMPOSE) up -d
	@echo "Services started"
	@echo "Calendar API: http://localhost:8080"

## docker-down: Stop all services
.PHONY: docker-down
docker-down:
	@echo "Stopping services..."
	@$(DOCKER_COMPOSE) down
	@echo "Services stopped"

## docker-logs: Show logs from all services
.PHONY: docker-logs
docker-logs:
	@$(DOCKER_COMPOSE) logs -f

## docker-clean: Stop services and remove volumes
.PHONY: docker-clean
docker-clean:
	@echo "Cleaning Docker environment..."
	@$(DOCKER_COMPOSE) down -v
	@echo "Docker environment cleaned"

## up: Start all services (alias for docker-up)
.PHONY: up
up: docker-up

## down: Stop all services (alias for docker-down)
.PHONY: down
down: docker-down

## integration-tests: Run integration tests
.PHONY: integration-tests
integration-tests:
	@echo "Starting integration tests..."
	@echo "Building and starting services..."
	@$(DOCKER_COMPOSE) up -d --build
	@echo "Waiting for services to be ready..."
	@powershell -Command "Start-Sleep -Seconds 20"
	@echo "Building integration tests..."
	@docker build -f build/Dockerfile.integration -t calendar-integration-tests .
	@echo "Running integration tests..."
	@docker run --rm --network calendar-network calendar-integration-tests && echo "Tests passed!" || (echo "Tests failed!" && $(DOCKER_COMPOSE) down && exit 1)
	@echo "Cleaning up..."
	@$(DOCKER_COMPOSE) down
	@echo "Integration tests completed successfully!"
