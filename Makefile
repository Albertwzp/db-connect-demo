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
ui-build:
	@if [ -d frontend ]; then \
		echo "building frontend..."; \
		cd frontend && \
		if command -v npm >/dev/null 2>&1; then \
			if [ -f package-lock.json ]; then \
				npm ci; \
			elif [ -f yarn.lock ]; then \
				yarn install; \
			else \
				npm install; \
			fi; \
			npm run build; \
		elif command -v yarn >/dev/null 2>&1; then \
			yarn install && yarn build; \
		else \
			echo "npm/yarn not found, skipping frontend build"; \
		fi; \
		cd - >/dev/null; \
	else \
		echo "frontend dir not found, skipping"; \
	fi

ui-clean:
	@if [ -d frontend/dist ]; then rm -rf frontend/dist; else if exist frontend\dist rmdir /S /Q frontend\dist; fi; fi

run-service: build ui-build
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
