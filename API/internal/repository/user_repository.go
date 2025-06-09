package repository

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/vida-plus/api/models"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) models.UserRepository {
	return &UserRepository{
		collection: db.Collection("users"),
	}
}

func (r *UserRepository) HealthCheck(ctx context.Context) error {
	return r.collection.Database().Client().Ping(ctx, nil)
}

func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	logger := slog.With(
		slog.String("repository", "UserRepository"),
		slog.String("method", "CreateUser"),
		slog.String("userID", user.ID),
	)

	_, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		logger.Error("failed to create user", slog.Any("error", err))
		return models.NewInternalError("failed to create user")
	}

	logger.Info("user created successfully")
	return nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	logger := slog.With(
		slog.String("repository", "UserRepository"),
		slog.String("method", "GetUserByEmail"),
		slog.String("email", email),
	)

	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			logger.Info("user not found")
			return nil, nil
		}
		logger.Error("failed to get user", slog.Any("error", err))
		return nil, models.NewInternalError("failed to get user")
	}

	logger.Info("user found successfully")
	return &user, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	logger := slog.With(
		slog.String("repository", "UserRepository"),
		slog.String("method", "GetByID"),
		slog.String("userID", id),
	)

	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			logger.Info("user not found")
			return nil, nil
		}
		logger.Error("failed to get user by ID", slog.Any("error", err))
		return nil, models.NewInternalError("failed to get user by ID")
	}

	logger.Info("user found successfully")
	return &user, nil
}

func (r *UserRepository) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	logger := slog.With(
		slog.String("repository", "UserRepository"),
		slog.String("method", "GetAllUsers"),
	)

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		logger.Error("failed to find users", slog.Any("error", err))
		return nil, models.NewInternalError("failed to find users")
	}
	defer cursor.Close(ctx)

	var users []*models.User
	if err = cursor.All(ctx, &users); err != nil {
		logger.Error("failed to decode users", slog.Any("error", err))
		return nil, models.NewInternalError("failed to decode users")
	}

	logger.Info("users retrieved successfully", slog.Int("count", len(users)))
	return users, nil
}
