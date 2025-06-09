# Makefile para automação de tarefas Go

.PHONY: build run test lint fmt clean mocks docs docker-up docker-down dev

build:
	go build -o bin/api ./cmd/api

docker-up:
	docker compose up -d
	@echo "Aguardando serviços Docker ficarem prontos..."
	@sleep 5

docker-down:
	docker compose down

run: docker-up
	@echo "Iniciando servidor da API..."
	go run ./cmd/api

dev:
	@echo "Iniciando ambiente de desenvolvimento..."
	@make docker-up
	@echo "Serviços Docker iniciados. Iniciando servidor da API..."
	@go run ./cmd/api

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

# Comandos auxiliares
check-deps:
	@which docker > /dev/null || (echo "Docker não encontrado. Instale Docker primeiro." && exit 1)
	@docker compose version > /dev/null || (echo "Docker Compose não encontrado. Instale Docker Compose primeiro." && exit 1)

# Limpar tudo (containers + build artifacts)
clean-all: docker-down clean
	docker system prune -f

# Ver logs dos serviços Docker
logs:
	docker compose logs -f

# Reiniciar serviços Docker
restart: docker-down docker-up

# Ajuda
help:
	@echo "Comandos disponíveis:"
	@echo "  build       - Compila a aplicação"
	@echo "  run         - Inicia Docker + servidor da API"
	@echo "  dev         - Alias para 'run' (ambiente de desenvolvimento)"
	@echo "  test        - Executa testes"
	@echo "  lint        - Executa linter"
	@echo "  fmt         - Formata código Go"
	@echo "  clean       - Remove arquivos de build"
	@echo "  clean-all   - Para Docker + remove arquivos de build"
	@echo "  docs        - Gera documentação Swagger"
	@echo "  docs-serve  - Mostra como acessar documentação"
	@echo "  mocks       - Gera mocks para testes"
	@echo "  docker-up   - Inicia apenas serviços Docker"
	@echo "  docker-down - Para serviços Docker"
	@echo "  logs        - Mostra logs dos serviços Docker"
	@echo "  restart     - Reinicia serviços Docker"
	@echo "  check-deps  - Verifica dependências necessárias"
