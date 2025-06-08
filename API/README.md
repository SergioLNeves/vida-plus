# API de Autenticação e Autorização em Go

Este projeto implementa uma base para um sistema de autenticação e autorização utilizando Go, JWT e Echo, seguindo as melhores práticas de estrutura, modularidade e testes.

## Funcionalidades

- Registro de usuário com hash de senha (bcrypt)
- Login de usuário com geração de token JWT
- Middleware de autenticação JWT para rotas protegidas
- Estrutura modular: `cmd/`, `internal/`, `models/`, `pkg/`
- Handlers desacoplados e organizados
- Makefile para automação de build, test, lint e run

## Estrutura do Projeto

```
go.mod
Makefile
cmd/
  api/
    main.go
internal/
  auth/
    handler.go
    interface.go
    service.go
  middleware/
    jwt.go
  user/
    interface.go
    service.go
models/
  auth.go
  user.go
pkg/
  jwt.go
```

## Endpoints

- `POST /register` — Cria um novo usuário
  - Body: `{ "email": "user@example.com", "password": "senha" }`
- `POST /login` — Autentica usuário e retorna JWT
  - Body: `{ "email": "user@example.com", "password": "senha" }`
  - Response: `{ "token": "<jwt>" }`
- `GET /protected` — Exemplo de rota protegida (necessita JWT no header Authorization)

## Como rodar

```sh
make build      # Compila o binário
make run        # Executa a API
make test       # Roda os testes
make fmt        # Formata o código
make lint       # (Requer golangci-lint)
```

## Exemplo de uso

### Registro
```sh
curl -X POST http://localhost:8080/register \
  -H 'Content-Type: application/json' \
  -d '{"email":"user@example.com","password":"senha"}'
```

### Login
```sh
curl -X POST http://localhost:8080/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"user@example.com","password":"senha"}'
```

### Rota protegida
```sh
curl -H "Authorization: Bearer <token>" http://localhost:8080/protected
```

## Tecnologias e Bibliotecas
- [Echo](https://echo.labstack.com/) — Framework HTTP
- [bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt) — Hash de senha
- [JWT](https://github.com/golang-jwt/jwt) — Geração e validação de tokens

## Estrutura recomendada
- Código de domínio em inglês
- Comentários e documentação em português
- Modularização por responsabilidade
- Testes e automação via Makefile

---

> Projeto inicial para sistemas modernos de autenticação e autorização em Go.
