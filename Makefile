BINARY_NAME = pr-service

.PHONY: build run test tidy up down logs loadtest

build:
	go build -o bin/$(BINARY_NAME) ./cmd/server

run:
	go run ./cmd/server

test:
	go test ./...

test-reset:
	docker-compose --profile test up -d --build
	go test ./...
	docker-compose --profile test down -v

tidy:
	go mod tidy

lint:
	golangci-lint run ./...

up:
	docker-compose --profile dev up -d --build

up-test:
	docker-compose --profile test up -d --build

down:
	docker-compose --profile dev --profile test down -v

logs:
	docker-compose --profile dev logs -f

logs-test:
	docker-compose --profile test logs -f

loadtest:
	k6 run loadtest/pr_loadtest.js