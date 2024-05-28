BIN := "./bin/calendar"
BIN_DATABASE := "./bin/calendar-database"
BIN_SCHEDULER := "./bin/calendar-scheduler"
BIN_SENDER := "./bin/calendar-sender"
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
	go build -v -o $(BIN_SCHEDULER) -ldflags "$(LDFLAGS)" ./cmd/scheduler
	go build -v -o $(BIN_SENDER) -ldflags "$(LDFLAGS)" ./cmd/sender

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
		--publish="8070:8070" \
		--publish="8080:8080" \
		$(DOCKER_IMG)

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
	go test -race -v -count 100 -timeout=2m ./internal/...

install-go:
	wget -c https://go.dev/dl/go1.22.3.linux-amd64.tar.gz \
	&& sudo tar -xvf go1.22.3.linux-amd64.tar.gz -C /usr/local \
	&& rm go1.22.3.linux-amd64.tar.gz

go-version:
	go version

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(LINTER_PATH) v1.55.2

lint-version: install-lint-deps
	$(LINTER_BIN) version

lint-fix: install-lint-deps
	$(LINTER_BIN) run --fix ./...

lint: install-lint-deps
	$(LINTER_BIN) run ./...

generate:
	protoc --go_out=./ --go_opt=paths=source_relative --go-grpc_out=./ --go-grpc_opt=paths=source_relative api/proto/calendar.proto

.PHONY: build build-database run run-database build-img run-img up down restart version migrate-up migrate-down test install-go go-version lint-version lint-fix lint generate
