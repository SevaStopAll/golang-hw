BIN := "./bin/calendar"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	@echo "Building calendar API..."
	go build -o bin/calendar.exe ./cmd/calendar
	@echo "Building calendar scheduler..."
	go build -o bin/calendar_scheduler.exe ./cmd/scheduler
	@echo "Building calendar storer..."
	go build -o bin/calendar_storer.exe ./cmd/storer
	@echo "Build complete! Executables are in bin/ directory"

run-calendar:
	./bin/calendar --config=configs/calendar.yaml

run-scheduler:
	./bin/calendar_scheduler --config=configs/scheduler.yaml

run-storer:
	./bin/calendar_storer --config=configs/storer.yaml

clean:
	rm -rf bin/

version: build
	$(BIN) version

test:
	go test -race ./internal/...

lint:
	go vet ./...
	go fmt ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest ./...

.PHONY: build run build-img run-img version test lint