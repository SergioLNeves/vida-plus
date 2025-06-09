# Documentação da API Vida Plus

Esta pasta contém a documentação da API gerada automaticamente usando Swagger/OpenAPI.

## Arquivos

- `docs.go` - Código Go gerado automaticamente pelo Swag
- `swagger.json` - Especificação da API em formato JSON
- `swagger.yaml` - Especificação da API em formato YAML
- `swagger-config.json` - Configuração adicional do Swagger
- `postman-collection.json` - Coleção do Postman para testes
- `curl-examples.md` - Exemplos de comandos cURL para teste
- `README.md` - Esta documentação

## Como visualizar a documentação

### 1. Interface Web (Swagger UI)

Após iniciar a aplicação, acesse:
```
http://localhost:8080/swagger/index.html
```

### 2. Arquivo JSON/YAML

Você pode usar os arquivos `swagger.json` ou `swagger.yaml` em qualquer ferramenta que suporte OpenAPI/Swagger, como:
- [Swagger Editor](https://editor.swagger.io/)
- [Postman](https://www.postman.com/)
- [Insomnia](https://insomnia.rest/)

### 3. Coleção do Postman

Importe o arquivo `postman-collection.json` no Postman para ter uma coleção completa com todos os endpoints configurados. A coleção inclui:
- Variáveis de ambiente pré-configuradas
- Scripts para capturar automaticamente o token JWT após login
- Exemplos de todas as requisições

### 4. Comandos cURL

Para desenvolvedores que preferem linha de comando, consulte o arquivo `curl-examples.md` que contém:
- Exemplos completos de todas as rotas
- Scripts automatizados para teste
- Casos de teste para cenários de erro

## Regenerar a documentação

Para regenerar a documentação após alterações nos comentários do código:

```bash
swag init -g cmd/api/main.go -o doc
```

## Endpoints disponíveis

### Autenticação
- `POST /v1/auth/register` - Registrar novo usuário
- `POST /v1/auth/login` - Login de usuário

### Rotas protegidas
- `GET /v1/protected` - Informações protegidas (requer autenticação)

### Utilitários
- `GET /health` - Verificação de saúde da API

## Autenticação

A API usa JWT (JSON Web Tokens) para autenticação. Para acessar rotas protegidas:

1. Faça login para obter um token
2. Inclua o token no header Authorization: `Bearer <token>`

## Exemplos de uso

### Registrar usuário
```bash
curl -X POST "http://localhost:8080/v1/auth/register" \
     -H "Content-Type: application/json" \
     -d '{
       "email": "user@example.com",
       "password": "mypassword123"
     }'
```

### Login
```bash
curl -X POST "http://localhost:8080/v1/auth/login" \
     -H "Content-Type: application/json" \
     -d '{
       "email": "user@example.com",
       "password": "mypassword123"
     }'
```

### Acessar rota protegida
```bash
curl -X GET "http://localhost:8080/v1/protected" \
     -H "Authorization: Bearer <seu-jwt-token>"
```
