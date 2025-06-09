# Vida Plus API

Sistema de gestão de saúde e bem-estar desenvolvido em Go, com arquitetura limpa, autenticação JWT e diferenciação de tipos de usuário para ambientes hospitalares e clínicos.

> **⚠️ Status do Projeto**: Este projeto está em desenvolvimento ativo. Algumas configurações estão hardcoded para facilitar o desenvolvimento local.

## 🚀 Características Principais

- **🔐 Autenticação JWT**: Sistema completo de autenticação com tokens seguros
- **👥 Tipos de Usuário**: Suporte para múltiplos perfis (paciente, médico, enfermeiro, admin, recepcionista)
- **🛡️ Autorização por Papel**: Middleware para controle de acesso baseado em função
- **📦 MongoDB**: Integração robusta com padrão repository
- **📚 Documentação Swagger**: API documentada automaticamente com OpenAPI 3.0
- **🧪 Testes de Integração**: Cobertura completa usando testcontainers-go
- **💊 Health Check**: Monitoramento de conectividade do banco de dados
## 📁 Estrutura do Projeto

```
API/
├── cmd/api/                    # Ponto de entrada da aplicação
│   └── main.go                 # Aplicação principal com rotas simplificadas
├── internal/                   # Código interno (não exportável)
│   ├── domain/                 # Modelos de domínio e regras de negócio
│   │   ├── auth.go             # Estruturas de autenticação
│   │   ├── errors.go           # Definições de erros customizados
│   │   ├── repository.go       # Interfaces de repositório
│   │   ├── requests.go         # Modelos de requisição/resposta
│   │   └── user.go             # Modelo de usuário
│   ├── handler/                # Handlers HTTP
│   │   ├── admin_handler.go    # Endpoints administrativos
│   │   ├── auth_handler.go     # Endpoints de autenticação
│   │   ├── health_handler.go   # Endpoints de health check
│   │   ├── protected_handler.go # Rotas protegidas de exemplo
│   │   └── validator.go        # Validação de requisições
│   ├── healthcheck/            # Serviço de health check
│   │   └── healthcheck.go      # Implementação do health check
│   ├── middleware/             # Middlewares
│   │   ├── authorization.go    # Autorização baseada em papel
│   │   └── jwt.go              # Autenticação JWT
│   ├── repository/             # Camada de acesso a dados
│   │   └── user_repository.go  # Repositório de usuários
│   └── service/                # Camada de serviços
│       ├── auth_service.go     # Lógica de autenticação
│       └── user_service.go     # Lógica de usuários
├── mocks/                      # Mocks para testes
│   ├── auth_service_mocks.go   # Mocks do serviço de auth
│   ├── jwt_manager_mocks.go    # Mocks do gerenciador JWT
│   ├── repository_mocks.go     # Mocks de repositório
│   ├── user_repository_mocks.go # Mocks do repositório de usuários
│   └── user_store_mocks.go     # Mocks do store de usuários
├── pkg/                        # Pacotes utilitários (exportáveis)
│   ├── id.go                   # Geração de IDs
│   ├── jwt.go                  # Utilitários JWT
│   └── database/               # Utilitários de banco
│       └── mongodb.go          # Cliente MongoDB
├── test/integration/           # Testes de integração
│   ├── auth_test.go            # Testes de autenticação
│   ├── authorization_test.go   # Testes de autorização
│   ├── core_test.go            # Testes de funcionalidade core
│   ├── handlers_test.go        # Testes de handlers
│   ├── health_test.go          # Testes de health check
│   └── setup.go                # Infraestrutura de testes
├── doc/                        # Documentação Swagger
│   ├── docs.go                 # Documentação gerada
│   ├── postman-collection.json # Coleção do Postman
│   ├── swagger-config.json     # Configuração do Swagger
│   ├── swagger.json            # Especificação OpenAPI JSON
│   └── swagger.yaml            # Especificação OpenAPI YAML
├── docker-compose.yml          # Ambiente de desenvolvimento
├── Dockerfile                  # Configuração do container
├── Makefile                    # Comandos de automação
├── go.mod                      # Definição do módulo Go
└── go.sum                      # Checksums das dependências
```

## 👥 Tipos de Usuário

O sistema suporta os seguintes tipos de usuário com diferentes níveis de permissão:

| Tipo | Descrição | Permissões |
|------|-----------|------------|
| **👤 Patient** | Paciente do sistema | Acesso básico, visualização do próprio perfil |
| **👨‍⚕️ Doctor** | Médico | Acesso a pacientes, prescrições, consultas |
| **👩‍⚕️ Nurse** | Enfermeiro(a) | Cuidados com pacientes, registros médicos |
| **👨‍💼 Admin** | Administrador | Acesso total ao sistema, gestão de usuários |
| **🏥 Receptionist** | Recepcionista | Agendamentos, cadastros, atendimento |

### Campos Específicos por Tipo

- **Médicos**: CRM, especialidade
- **Enfermeiros**: COREN, setor
- **Pacientes**: Data de nascimento, histórico médico
- **Funcionários**: Departamento, cargo

## 🛠️ API Endpoints

### 🔐 Autenticação
- `POST /v1/auth/register` - Cadastro de usuário com tipo específico
- `POST /v1/auth/login` - Login de usuário

### 🔒 Rotas Protegidas
- `GET /v1/protected` - Exemplo de endpoint protegido

### 👨‍💼 Administração (Admin apenas)
- `GET /v1/admin/users` - Listar todos os usuários
- `GET /v1/admin/stats` - Estatísticas do sistema

### 💊 Health Check
- `GET /health` - Status de conectividade do banco de dados

### 📚 Documentação
- `GET /swagger/index.html` - Interface Swagger UI
- `GET /swagger/doc.json` - Especificação OpenAPI JSON

## 🚀 Início Rápido

### 🐳 Usando Docker Compose (Recomendado)

```bash
# Iniciar o ambiente de desenvolvimento
docker-compose up -d

# A API estará disponível em http://localhost:8080
# Documentação Swagger em http://localhost:8080/swagger/index.html
```

### 🔧 Configuração Manual

1. **Instalar Dependências**
   ```bash
   go mod tidy
   ```

2. **Iniciar MongoDB**
   ```bash
   # Usando Docker
   docker run -d -p 27017:27017 --name mongodb mongo:latest
   ```

3. **Executar a Aplicação**
   ```bash
   go run cmd/api/main.go
   ```

4. **Verificar se está funcionando**
   ```bash
   curl http://localhost:8080/health
   ```

## 🧪 Testes

### Testes de Integração

```bash
# Executar todos os testes de integração
go test ./test/integration/... -v

# Executar com cobertura
go test ./test/integration/... -v -cover
```

### Testes Unitários

```bash
# Todos os testes unitários
go test ./internal/... -v

# Testes de um pacote específico
go test ./internal/service/... -v
```

### Cobertura de Testes

```bash
# Cobertura geral
go test -cover ./...

# Cobertura detalhada com HTML
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

## 💻 Desenvolvimento

### Gerar Documentação Swagger

```bash
# Instalar swag
go install github.com/swaggo/swag/cmd/swag@latest

# Gerar documentação
swag init -g cmd/api/main.go -o doc
```

### Build para Produção

```bash
# Build nativo
go build -o bin/api cmd/api/main.go

# Build com Docker
docker build -t vida-plus-api .
```

## 📝 Exemplos de Uso

### Cadastrar um Novo Paciente

```bash
curl -X POST http://localhost:8080/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "paciente@exemplo.com",
    "password": "senha123",
    "type": "patient",
    "profile": {
      "first_name": "João",
      "last_name": "Silva",
      "cpf": "12345678901",
      "phone": "+5511999999999",
      "date_of_birth": "1990-01-01"
    }
  }'
```

### Cadastrar um Médico

```bash
curl -X POST http://localhost:8080/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "medico@exemplo.com",
    "password": "senha123",
    "type": "doctor",
    "profile": {
      "first_name": "Dra. Maria",
      "last_name": "Santos",
      "cpf": "98765432101",
      "phone": "+5511888888888",
      "crm": "CRM/SP 123456",
      "speciality": "Cardiologia"
    }
  }'
```

### Login

```bash
curl -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "paciente@exemplo.com",
    "password": "senha123"
  }'
```

### Acessar Rota Protegida

```bash
curl -X GET http://localhost:8080/v1/protected \
  -H "Authorization: Bearer SEU_TOKEN_JWT"
```

### Acessar Estatísticas (Admin apenas)

```bash
curl -X GET http://localhost:8080/v1/admin/stats \
  -H "Authorization: Bearer TOKEN_DO_ADMIN"
```

## 🔧 Configuração

### Configurações Atuais (Hardcoded)

| Configuração | Valor | Arquivo |
|--------------|-------|---------|
| **MongoDB URI** | `mongodb://localhost:27017/vida_plus` | `pkg/database/mongodb.go` |
| **JWT Secret** | `local-development-secret-key` | `pkg/jwt.go` |
| **Porta do Servidor** | `8080` | `cmd/api/main.go` |
| **Nome do Banco** | `vida_plus` | `cmd/api/main.go` |

## 🔒 Recursos de Segurança

- **🔐 Autenticação JWT**: Tokens seguros com tempo de expiração de 24 horas
- **🛡️ Hash de Senhas**: bcrypt para armazenamento seguro de senhas
- **👮 Autorização por Papel**: Middleware para controle de acesso baseado em função
- **✅ Validação de Entrada**: Validação rigorosa usando go-playground/validator

## 🛠️ Tecnologias Utilizadas

| Tecnologia | Descrição |
|------------|-----------|
| **🐹 Go** | Linguagem principal do backend |
| **⚡ Echo** | Framework HTTP de alta performance |
| **🍃 MongoDB** | Banco de dados NoSQL |
| **🔑 JWT** | Autenticação baseada em tokens |
| **🔐 bcrypt** | Hash seguro de senhas |
| **🧪 Testify** | Framework de testes |
| **🐳 Testcontainers** | Testes de integração com containers |
| **📚 Swagger** | Documentação automática da API |

## 🔧 Solução de Problemas

**❌ Erro de Conexão com MongoDB**
```bash
# Verificar se o MongoDB está rodando
docker ps | grep mongo

# Reiniciar o container
docker-compose restart mongodb
```

**❌ Token JWT Inválido**
```bash
# Fazer login novamente para obter novo token
curl -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password"}'
```

**❌ Porta já em uso**
```bash
# Verificar o que está usando a porta 8080
lsof -i :8080
```

## 📄 Licença

Este projeto faz parte do sistema Vida Plus de gestão de saúde e bem-estar.

---

<div align="center">
  <h3>🏥 Vida Plus API</h3>
  <p><em>Sistema de gestão de saúde e bem-estar desenvolvido com ❤️ em Go</em></p>
</div>
