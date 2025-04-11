package user

import (
	"GoMail/app/repository/models"
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

const (
	userCollection = "users"
)

// Repository defines the user repository interface
type Repository interface {
	Save(ctx context.Context, user *models.User) error
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByID(ctx context.Context, id string) (*models.User, error)
	FindAll(ctx context.Context, filter interface{}, page, limit int) ([]*models.User, int64, error)
}

type repository struct {
	db *mongo.Database
}

// New creates a new user repository
func New(db *mongo.Database) Repository {
	return &repository{
		db: db,
	}
} 