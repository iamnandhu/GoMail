package email

import (
	"context"
	"errors"

	"GoMail/app/repository/models"

	"go.mongodb.org/mongo-driver/mongo"
)

const CollectionName = "emails"

var (
	ErrInvalidEmailType = errors.New("invalid email type")
	ErrInvalidID        = errors.New("invalid ID type")
	ErrEmailNotFound    = errors.New("email not found")
)

type EmailRepository interface {
	Save(ctx context.Context, email *models.Email) error
	FindAll(ctx context.Context, filter interface{}, page, limit int) ([]*models.Email, int64, error)
	Delete(ctx context.Context, id interface{}) error
}

type mongoDB struct {
	collection *mongo.Collection
}

func New(database *mongo.Database) EmailRepository {
	return &mongoDB{
		collection: database.Collection(CollectionName),
	}
}
