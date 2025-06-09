package repository

import (
	"context"
	"log/slog"

	"github.com/vida-plus/api/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) domain.UserRepository {
	return &UserRepository{
		collection: db.Collection("users"),
	}
}

func (r *UserRepository) HealthCheck(ctx context.Context) error {
	return r.collection.Database().Client().Ping(ctx, nil)
}

func (r *UserRepository) CreateUser(ctx context.Context, user *domain.User) error {
	logger := slog.With(
		slog.String("repository", "UserRepository"),
		slog.String("method", "CreateUser"),
		slog.String("userID", user.ID),
	)

	_, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		logger.Error("failed to create user", slog.Any("error", err))
		return domain.NewInternalError("failed to create user")
	}

	logger.Info("user created successfully")
	return nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	logger := slog.With(
		slog.String("repository", "UserRepository"),
		slog.String("method", "GetUserByEmail"),
		slog.String("email", email),
	)

	var user domain.User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			logger.Info("user not found")
			return nil, nil
		}
		logger.Error("failed to get user", slog.Any("error", err))
		return nil, domain.NewInternalError("failed to get user")
	}

	logger.Info("user found successfully")
	return &user, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	logger := slog.With(
		slog.String("repository", "UserRepository"),
		slog.String("method", "GetByID"),
		slog.String("userID", id),
	)

	var user domain.User
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			logger.Info("user not found")
			return nil, nil
		}
		logger.Error("failed to get user by ID", slog.Any("error", err))
		return nil, domain.NewInternalError("failed to get user by ID")
	}

	logger.Info("user found successfully")
	return &user, nil
}

func (r *UserRepository) GetAllUsers(ctx context.Context) ([]*domain.User, error) {
	logger := slog.With(
		slog.String("repository", "UserRepository"),
		slog.String("method", "GetAllUsers"),
	)

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		logger.Error("failed to find users", slog.Any("error", err))
		return nil, domain.NewInternalError("failed to find users")
	}
	defer cursor.Close(ctx)

	var users []*domain.User
	if err = cursor.All(ctx, &users); err != nil {
		logger.Error("failed to decode users", slog.Any("error", err))
		return nil, domain.NewInternalError("failed to decode users")
	}

	logger.Info("users retrieved successfully", slog.Int("count", len(users)))
	return users, nil
}
