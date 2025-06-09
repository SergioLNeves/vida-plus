# Makefile para automação de tarefas Go

.PHONY: build run test lint fmt clean mocks docs

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

docs:
	swag init -g cmd/api/main.go -o doc

docs-serve:
	@echo "Documentação Swagger estará disponível em: http://localhost:8080/swagger/index.html"
	@echo "Execute 'make run' em outro terminal para iniciar o servidor"

mocks:
	@echo "Generating mocks..."
	@rm -rf mocks/*.go
	@mockery
	@echo "Mocks generated successfully!"
