package emaillog

import (
	"context"
	"time"

	"GoMail/app/repository/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Save inserts or updates an email log in the database
func (db *mongoDB) Save(ctx context.Context, emailLog *models.EmailLog) error {
	// Set creation time if this is a new document
	if emailLog.ID.IsZero() {
		emailLog.ID = primitive.NewObjectID()
		emailLog.CreatedAt = time.Now()
	}

	// Insert document
	_, err := db.collection.InsertOne(ctx, emailLog)
	if err != nil {
		return err
	}

	return nil
} 