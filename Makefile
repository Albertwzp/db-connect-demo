# Makefile for db-bench

BIN ?= ./db-bench.exe
PORT ?= 8080
BACKENDS ?= backends.json

.PHONY: all tidy build clean test-build run-service run-postgres run-mysql run-sqlite run-kafka run-solace help

all: build

tidy:
	@echo "tidying modules..."
	go mod tidy

build:
	@echo "building db-bench..."
	go build -o $(BIN)

clean:
	@echo "cleaning..."
	@if [ -f $(BIN) ]; then rm -f $(BIN); else if exist $(BIN) del /Q $(BIN); fi; fi

# cross-build example
test-build:
	@echo "building linux binary..."
	GOOS=linux GOARCH=amd64 go build -o db-bench-linux

# Run service with backends file (set BACKENDS and PORT as needed)
run-service: build
	@echo "starting service with backends=$(BACKENDS) port=$(PORT) - UI at http://localhost:$(PORT)/ui"
	$(BIN) -backends-file=$(BACKENDS) -port=$(PORT)

# Example run targets (override DSN by providing DSN variable)
run-postgres: build
	@echo "running against Postgres"
	$(BIN) -driver=postgres -dsn="$(DSN)"

run-mysql: build
	@echo "running against MySQL"
	$(BIN) -driver=mysql -dsn="$(DSN)"

run-sqlite: build
	@echo "running against SQLite (in-memory)"
	$(BIN) -driver=sqlite -dsn="$(DSN)"

run-kafka: build
	@echo "running Kafka producer"
	$(BIN) -driver=kafka -dsn="$(DSN)"

run-solace: build
	@echo "running Solace via MQTT"
	$(BIN) -driver=solace -dsn="$(DSN)"

help:
	@echo "Available targets: all tidy build clean test-build run-service run-postgres run-mysql run-sqlite run-kafka run-solace"
	@echo "Examples: make run-postgres DSN='postgres://user:pass@localhost:5432/db?sslmode=disable'"
