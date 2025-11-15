BINARY_NAME = pr-service

.PHONY: build run test tidy up down logs

build:
	go build -o bin/$(BINARY_NAME) ./cmd/server

run:
	go run ./cmd/server

test:
	go test ./...

tidy:
	go mod tidy

lint:
	golangci-lint run ./...

up:
	docker-compose up --build

down:
	docker-compose down -v

logs:
	docker-compose logs -f
