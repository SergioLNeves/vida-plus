# Vida Plus API - Healthcare System with User Differentiation

A simplified Go API with user type differentiation for healthcare systems, built with Echo framework, MongoDB, and comprehensive integration testing.

## Features

- **User Authentication**: Complete JWT-based authentication system
- **User Type Differentiation**: Support for multiple user types (patient, doctor, nurse, admin, receptionist)
- **Role-Based Authorization**: Middleware for role-based access control
- **MongoDB Integration**: Repository pattern with error handling
- **Swagger Documentation**: Auto-generated API documentation
- **Integration Testing**: Comprehensive tests using testcontainers-go
- **Health Check**: Database connectivity monitoring
## Project Structure

```
API/
├── cmd/api/                    # Application entry point
│   └── main.go                 # Main application with simplified routes
├── internal/
│   ├── auth/                   # Authentication service
│   │   ├── service.go          # AuthService implementation
│   │   └── service_test.go     # Unit tests
│   ├── handlers/               # HTTP handlers
│   │   ├── auth_handler.go     # Authentication endpoints
│   │   ├── health_handler.go   # Health check endpoints
│   │   ├── protected_handler.go # Protected route examples
│   │   └── validator.go        # Request validation
│   ├── middleware/             # Middleware components
│   │   ├── authorization.go    # Role-based authorization
│   │   └── jwt.go              # JWT authentication
│   ├── repository/             # Data access layer
│   │   └── user_repository.go  # User repository
│   └── user/                   # User service
│       └── service.go          # User business logic
├── models/                     # Data models
│   ├── auth.go                 # Authentication models
│   ├── requests.go             # Request/response models
│   └── user.go                 # User model with types
├── pkg/                        # Utility packages
│   ├── id.go                   # ID generation
│   └── jwt.go                  # JWT utilities
├── test/integration/           # Integration tests
│   ├── setup.go                # Test infrastructure
│   ├── auth_test.go            # Authentication tests
│   ├── core_test.go            # Core functionality tests
│   └── health_test.go          # Health check tests
├── doc/                        # Swagger documentation
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── mocks/                      # Test mocks
├── docker-compose.yml          # Development environment
├── Dockerfile                  # Container configuration
├── go.mod                      # Go module definition
└── go.sum                      # Go module checksums
```

## User Types

The system supports the following user types with different permission levels:

- **Patient**: Basic user with limited access
- **Doctor**: Medical professional with patient access
- **Nurse**: Healthcare provider with specific permissions
- **Admin**: System administrator with full access
- **Receptionist**: Front desk staff with administrative tasks

## API Endpoints

### Authentication
- `POST /v1/auth/register` - User registration with type
- `POST /v1/auth/login` - User login

### Protected Routes
- `GET /v1/profile` - Get user profile (authenticated)
- `GET /v1/protected` - Example protected endpoint

### Health Check
- `GET /health` - Database connectivity status

### Documentation
- `GET /swagger/index.html` - Swagger UI

## Quick Start

### Using Docker Compose (Recommended)

```bash
# Start the development environment
docker-compose up -d

# The API will be available at http://localhost:8080
# Swagger documentation at http://localhost:8080/swagger/index.html
```

### Manual Setup

1. **Install Dependencies**
   ```bash
   go mod tidy
   ```

2. **Start MongoDB**
   ```bash
   # Using Docker
   docker run -d -p 27017:27017 --name mongodb mongo:latest
   ```

3. **Set Environment Variables**
   ```bash
   export MONGODB_URI="mongodb://localhost:27017"
   export JWT_SECRET="your-secret-key"
   export PORT="8080"
   ```

4. **Run the Application**
   ```bash
   go run cmd/api/main.go
   ```

## Testing

### Integration Tests

Run comprehensive integration tests with testcontainers:

```bash
go test ./test/integration/... -v
```

### Unit Tests

Run unit tests for individual components:

```bash
go test ./internal/auth/... -v
```

### Test Coverage

```bash
go test -cover ./...
```

## Development

### Generate Swagger Documentation

```bash
# Install swag
go install github.com/swaggo/swag/cmd/swag@latest

# Generate docs
swag init -g cmd/api/main.go -o doc
```

### Build for Production

```bash
go build -o bin/api cmd/api/main.go
```

## Usage Examples

### Register a New Patient

```bash
curl -X POST http://localhost:8080/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "patient@example.com",
    "password": "secure123",
    "name": "John Doe",
    "type": "patient",
    "profile": {
      "dateOfBirth": "1990-01-01",
      "phoneNumber": "+1234567890"
    }
  }'
```

### Login

```bash
curl -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "patient@example.com",
    "password": "secure123"
  }'
```

### Access Protected Route

```bash
curl -X GET http://localhost:8080/v1/profile \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `MONGODB_URI` | MongoDB connection string | `mongodb://localhost:27017` |
| `JWT_SECRET` | JWT signing secret | `your-secret-key` |
| `PORT` | Server port | `8080` |
| `DB_NAME` | Database name | `vida_plus` |

## Architecture

The project follows clean architecture principles:

- **Domain Models**: Core business entities in `models/`
- **Repository Pattern**: Data access abstraction in `internal/repository/`
- **Service Layer**: Business logic in `internal/auth/` and `internal/user/`
- **Handlers**: HTTP transport layer in `internal/handlers/`
- **Middleware**: Cross-cutting concerns in `internal/middleware/`

## Security Features

- **JWT Authentication**: Secure token-based authentication
- **Password Hashing**: bcrypt for secure password storage
- **Role-Based Authorization**: Middleware for access control
- **Input Validation**: Request validation using go-playground/validator

## Future Expansion

The system is designed for easy expansion. To add role-specific functionality:

1. Create new handlers in `internal/handlers/`
2. Add role-specific routes in `cmd/api/main.go`
3. Use existing authorization middleware for access control
4. Add integration tests in `test/integration/`

## Contributing

1. Follow the existing code structure and patterns
2. Add integration tests for new features
3. Update Swagger documentation
4. Follow Go best practices and the Uber Go Style Guide

## Technologies Used

- **Framework**: Echo v4 - High performance HTTP framework
- **Database**: MongoDB - Document database with Go driver
- **Authentication**: JWT tokens with bcrypt password hashing
- **Documentation**: Swagger/OpenAPI 3.0 with auto-generation
- **Testing**: Testcontainers-go for integration testing
- **Validation**: go-playground/validator for request validation
- **Containerization**: Docker and Docker Compose

## License

This project is part of the Vida Plus healthcare system.
