SHELL := /bin/zsh

APP_NAME := gateway
BIN_DIR := bin
CMD_DIR := ./cmd/$(APP_NAME)

.PHONY: init fmt lint test build run clean

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
	go build -o $(BIN_DIR)/$(APP_NAME) $(CMD_DIR)

run:
	go run $(CMD_DIR)

clean:
	rm -rf $(BIN_DIR)
