SHELL := /bin/zsh

APP_NAME := gateway
BIN_DIR := bin

.PHONY: init fmt lint test build run run-license-svc clean docker-up docker-down

init:
	go mod tidy

fmt:
	gofmt -w $(shell find . -name '*.go' -not -path './vendor/*')

lint:
	go vet ./...

test:
	go test ./...

build:
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/gateway ./cmd/gateway
	go build -o $(BIN_DIR)/license-svc ./cmd/license-svc

run:
	go run ./cmd/gateway

run-license-svc:
	go run ./cmd/license-svc

docker-up:
	docker compose -f deploy/docker/docker-compose.yml up -d

docker-down:
	docker compose -f deploy/docker/docker-compose.yml down

clean:
	rm -rf $(BIN_DIR)
