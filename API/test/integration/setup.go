// Package integration contains integration tests for the Vida Plus API
package integration

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/vida-plus/api/internal/auth"
	"github.com/vida-plus/api/internal/handlers"
	"github.com/vida-plus/api/internal/middleware"
	"github.com/vida-plus/api/internal/repository"
	"github.com/vida-plus/api/internal/user"
	"github.com/vida-plus/api/models"
	"github.com/vida-plus/api/pkg"
)

// TestContainer holds the MongoDB test container and related resources
type TestContainer struct {
	Container    testcontainers.Container
	MongoClient  *mongo.Client
	Database     *mongo.Database
	DatabaseName string
	URI          string
}

// TestApp holds the application dependencies for testing
type TestApp struct {
	Echo             *echo.Echo
	UserRepo         models.UserRepository
	AuthService      models.AuthService
	JWTManager       models.JWTManager
	AuthHandler      *handlers.AuthHandler
	ProtectedHandler *handlers.ProtectedHandler
	HealthHandler    *handlers.HealthHandler
}

// SetupMongoDB creates a MongoDB test container
func SetupMongoDB(ctx context.Context, t *testing.T) *TestContainer {
	t.Helper()

	req := testcontainers.ContainerRequest{
		Image:        "mongo:latest",
		ExposedPorts: []string{"27017/tcp"},
		WaitingFor:   wait.ForLog("Waiting for connections").WithStartupTimeout(60 * time.Second),
	}

	mongoContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("Failed to start MongoDB container: %v", err)
	}

	host, err := mongoContainer.Host(ctx)
	if err != nil {
		t.Fatalf("Failed to get container host: %v", err)
	}

	port, err := mongoContainer.MappedPort(ctx, "27017")
	if err != nil {
		t.Fatalf("Failed to get container port: %v", err)
	}

	uri := fmt.Sprintf("mongodb://%s:%s", host, port.Port())

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Test connection
	if err := client.Ping(ctx, nil); err != nil {
		t.Fatalf("Failed to ping MongoDB: %v", err)
	}

	databaseName := "vida_plus_test"
	db := client.Database(databaseName)

	return &TestContainer{
		Container:    mongoContainer,
		MongoClient:  client,
		Database:     db,
		DatabaseName: databaseName,
		URI:          uri,
	}
}

// TeardownMongoDB cleans up the MongoDB test container
func (tc *TestContainer) TeardownMongoDB(ctx context.Context, t *testing.T) {
	t.Helper()

	if tc.MongoClient != nil {
		if err := tc.MongoClient.Disconnect(ctx); err != nil {
			log.Printf("Error disconnecting from MongoDB: %v", err)
		}
	}

	if tc.Container != nil {
		if err := tc.Container.Terminate(ctx); err != nil {
			log.Printf("Error terminating MongoDB container: %v", err)
		}
	}
}

// SetupTestApp creates a test application with all dependencies
func SetupTestApp(tc *TestContainer) *TestApp {
	// Initialize repositories
	userRepo := repository.NewUserRepository(tc.Database)

	// Initialize services
	jwtManager := pkg.NewJWTManager()
	userService := user.NewUserService(userRepo)
	authService := auth.NewAuthService(userService, jwtManager)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	protectedHandler := handlers.NewProtectedHandler()
	healthHandler := handlers.NewHealthHandler(tc.MongoClient)

	// Setup Echo app
	e := echo.New()

	// Configure routes
	setupTestRoutes(e, jwtManager, authHandler, protectedHandler, healthHandler)

	return &TestApp{
		Echo:             e,
		UserRepo:         userRepo,
		AuthService:      authService,
		JWTManager:       jwtManager,
		AuthHandler:      authHandler,
		ProtectedHandler: protectedHandler,
		HealthHandler:    healthHandler,
	}
}

// setupTestRoutes configures all routes for testing
func setupTestRoutes(e *echo.Echo, jwtManager models.JWTManager, authHandler *handlers.AuthHandler,
	protectedHandler *handlers.ProtectedHandler, healthHandler *handlers.HealthHandler) {

	// Health check
	e.GET("/health", healthHandler.Check)

	// Auth routes
	v1 := e.Group("/v1")
	v1.POST("/auth/register", authHandler.Register)
	v1.POST("/auth/login", authHandler.Login)

	// Protected routes
	protected := v1.Group("", middleware.JWTMiddleware(jwtManager))
	protected.GET("/protected", protectedHandler.GetProtectedInfo)

	// Simple profile endpoint to demonstrate user differentiation
	protected.GET("/profile", func(c echo.Context) error {
		claims, err := models.GetAuthClaims(c.Get("claims"))
		if err != nil {
			return c.JSON(401, models.NewAPIError(401, err.Error()))
		}

		return c.JSON(200, map[string]interface{}{
			"user_id": claims.UserID,
			"email":   claims.Email,
			"type":    claims.UserType,
			"message": "Perfil do usuário - acesso baseado no tipo: " + string(claims.UserType),
		})
	})
}

// CleanDatabase removes all data from test database
func (tc *TestContainer) CleanDatabase(ctx context.Context, t *testing.T) {
	t.Helper()

	collections, err := tc.Database.ListCollectionNames(ctx, map[string]interface{}{})
	if err != nil {
		t.Fatalf("Failed to list collections: %v", err)
	}

	for _, collection := range collections {
		if err := tc.Database.Collection(collection).Drop(ctx); err != nil {
			t.Fatalf("Failed to drop collection %s: %v", collection, err)
		}
	}
}

// Cleanup cleans up the test container and database
func (tc *TestContainer) Cleanup(ctx context.Context) error {
	if tc.MongoClient != nil {
		tc.MongoClient.Disconnect(ctx)
	}
	if tc.Container != nil {
		return tc.Container.Terminate(ctx)
	}
	return nil
}
