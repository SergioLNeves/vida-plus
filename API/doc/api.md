# API Documentation

## Authentication Endpoints

### POST /v1/auth/register
Register a new user.

**Request:**
```json
{
  "email": "string",
  "password": "string"
}
```

**Response:**
```json
{
  "token": "string"
}
```

**Errors:**
- 400: Invalid request body
- 409: User already exists
- 500: Internal server error

### POST /v1/auth/login
Login with existing user.

**Request:**
```json
{
  "email": "string",
  "password": "string"
}
```

**Response:**
```json
{
  "token": "string"
}
```

**Errors:**
- 400: Invalid request body
- 401: Invalid credentials
- 500: Internal server error

## Protected Endpoints

### GET /v1/protected
Get protected information. Requires JWT token.

**Headers:**
```
Authorization: Bearer <token>
```

**Response:**
```json
{
  "message": "string"
}
```

**Errors:**
- 401: Unauthorized
- 500: Internal server error

## Health Check

### GET /health
Check API and MongoDB health status.

**Response:**
```json
{
  "status": "healthy"
}
```

**Errors:**
- 503: Service Unavailable (when database is down)
