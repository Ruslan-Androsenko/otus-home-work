BIN := "./bin/calendar"
BIN_DATABASE := "./bin/calendar-database"
CONTAINER_NAME="calendar"
DOCKER_IMG="calendar:develop"
LINTER_PATH=/tmp/bin
LINTER_BIN=/tmp/bin/golangci-lint

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

build-database:
	go build -v -o $(BIN_DATABASE) -ldflags "$(LDFLAGS)" ./cmd/calendar

run: build
	$(BIN) -config ./configs/config-default.toml

run-database: build-database
	$(BIN_DATABASE) -config ./configs/config-database.toml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run --rm \
		--name=$(CONTAINER_NAME) \
		--network="host" $(DOCKER_IMG)

up:
	docker compose up -d --build

down:
	docker compose down

restart: down up

version: build
	$(BIN) version

migrate-up: build-database
	$(BIN_DATABASE) -config ./configs/config-database.toml migration-up

migrate-down: build-database
	$(BIN_DATABASE) -config ./configs/config-database.toml migration-down

test:
	go test -race -v -count 100 -timeout=20m ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(LINTER_PATH) v1.50.1

lint: install-lint-deps
	golangci-lint run ./...

.PHONY: build build-database run run-database build-img run-img up down restart version migrate-up migrate-down test lint
