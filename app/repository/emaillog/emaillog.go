package emaillog

import (
	"context"
	"errors"

	"GoMail/app/repository/models"

	"go.mongodb.org/mongo-driver/mongo"
)

const CollectionName = "email_logs"

var (
	ErrInvalidEmailLogType = errors.New("invalid email log type")
	ErrInvalidID           = errors.New("invalid ID type")
	ErrEmailLogNotFound    = errors.New("email log not found")
)

type EmailLogRepository interface {
	Save(ctx context.Context, emailLog *models.EmailLog) error
	FindAll(ctx context.Context, filter interface{}, page, limit int) ([]*models.EmailLog, int64, error)
	FindByID(ctx context.Context, id string) (*models.EmailLog, error)
}

type mongoDB struct {
	collection *mongo.Collection
}

func New(database *mongo.Database) EmailLogRepository {
	return &mongoDB{
		collection: database.Collection(CollectionName),
	}
} 