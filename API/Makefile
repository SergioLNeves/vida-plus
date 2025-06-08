# Makefile para automação de tarefas Go

.PHONY: build run test lint fmt clean mocks

build:
	go build -o bin/api ./cmd/api

run:
	go run ./cmd/api

test:
	go test ./...

lint:
	golangci-lint run ./...

fmt:
	gofmt -s -w .

clean:
	rm -rf bin/

mocks:
	@echo "Generating mocks..."
	@rm -rf mocks/*.go
	@mockery
	@echo "Mocks generated successfully!"
