# Exemplos de comandos cURL para a API Vida Plus

Este arquivo contém exemplos práticos de como testar todos os endpoints da API usando cURL.

## Variáveis de ambiente úteis

```bash
export API_BASE_URL="http://localhost:8080"
export JWT_TOKEN=""  # Será preenchido após o login
```

## 1. Health Check

```bash
curl -X GET "$API_BASE_URL/health" \
     -H "Content-Type: application/json" \
     | jq
```

**Resposta esperada:**
```json
{
  "status": "healthy"
}
```

## 2. Registrar novo usuário

```bash
curl -X POST "$API_BASE_URL/v1/auth/register" \
     -H "Content-Type: application/json" \
     -d '{
       "email": "user@example.com",
       "password": "mypassword123"
     }' \
     | jq
```

**Resposta esperada:**
```json
{
  "id": "user123",
  "email": "user@example.com"
}
```

## 3. Login

```bash
# Fazer login e capturar o token
RESPONSE=$(curl -X POST "$API_BASE_URL/v1/auth/login" \
     -H "Content-Type: application/json" \
     -d '{
       "email": "user@example.com",
       "password": "mypassword123"
     }')

echo $RESPONSE | jq

# Extrair o token para uso posterior
export JWT_TOKEN=$(echo $RESPONSE | jq -r '.token')
echo "Token JWT: $JWT_TOKEN"
```

**Resposta esperada:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

## 4. Acessar rota protegida

```bash
curl -X GET "$API_BASE_URL/v1/protected" \
     -H "Authorization: Bearer $JWT_TOKEN" \
     -H "Content-Type: application/json" \
     | jq
```

**Resposta esperada:**
```json
{
  "message": "This is a protected endpoint",
  "user": "authenticated"
}
```

## Scripts de teste automatizado

### Teste completo do fluxo

```bash
#!/bin/bash

API_BASE_URL="http://localhost:8080"
EMAIL="test$(date +%s)@example.com"  # Email único para cada teste
PASSWORD="testpassword123"

echo "🏥 Testando API Vida Plus..."
echo "📧 Email de teste: $EMAIL"

# 1. Health Check
echo "1️⃣ Verificando saúde da API..."
HEALTH=$(curl -s "$API_BASE_URL/health")
echo "Health: $HEALTH"

# 2. Registrar usuário
echo "2️⃣ Registrando usuário..."
REGISTER_RESPONSE=$(curl -s -X POST "$API_BASE_URL/v1/auth/register" \
     -H "Content-Type: application/json" \
     -d "{
       \"email\": \"$EMAIL\",
       \"password\": \"$PASSWORD\"
     }")
echo "Registro: $REGISTER_RESPONSE"

# 3. Login
echo "3️⃣ Fazendo login..."
LOGIN_RESPONSE=$(curl -s -X POST "$API_BASE_URL/v1/auth/login" \
     -H "Content-Type: application/json" \
     -d "{
       \"email\": \"$EMAIL\",
       \"password\": \"$PASSWORD\"
     }")
echo "Login: $LOGIN_RESPONSE"

# Extrair token
JWT_TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.token')

if [ "$JWT_TOKEN" != "null" ] && [ "$JWT_TOKEN" != "" ]; then
    echo "✅ Token JWT obtido com sucesso"
    
    # 4. Acessar rota protegida
    echo "4️⃣ Acessando rota protegida..."
    PROTECTED_RESPONSE=$(curl -s -X GET "$API_BASE_URL/v1/protected" \
         -H "Authorization: Bearer $JWT_TOKEN")
    echo "Protegida: $PROTECTED_RESPONSE"
    
    echo "✅ Teste completo finalizado!"
else
    echo "❌ Falha ao obter token JWT"
fi
```

## Testes de erro

### Login com credenciais inválidas

```bash
curl -X POST "$API_BASE_URL/v1/auth/login" \
     -H "Content-Type: application/json" \
     -d '{
       "email": "invalid@example.com",
       "password": "wrongpassword"
     }' \
     | jq
```

### Acessar rota protegida sem token

```bash
curl -X GET "$API_BASE_URL/v1/protected" \
     -H "Content-Type: application/json" \
     | jq
```

### Registrar usuário com email inválido

```bash
curl -X POST "$API_BASE_URL/v1/auth/register" \
     -H "Content-Type: application/json" \
     -d '{
       "email": "invalid-email",
       "password": "mypassword123"
     }' \
     | jq
```

## Notas

- Todos os comandos assumem que a API está rodando em `http://localhost:8080`
- O `jq` é usado para formatar as respostas JSON (instale com `sudo apt install jq`)
- Para testes automatizados, salve o script em um arquivo e execute com `bash script.sh`
- Os tokens JWT têm validade de 24 horas por padrão
