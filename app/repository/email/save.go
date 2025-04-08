package email

import (
	"context"
	"time"

	"GoMail/app/repository/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Save creates or updates an email in the database
func (m *mongoDB) Save(ctx context.Context, email *models.Email) error {
	// Set update timestamp
	now := time.Now()
	email.UpdatedAt = now

	// Handle new vs existing document
	if email.ID.IsZero() {
		// This is a create operation
		email.CreatedAt = now

		// Insert document
		result, err := m.collection.InsertOne(ctx, email)
		if err != nil {
			return err
		}

		// Set the ID from the inserted document
		if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
			email.ID = oid
		}

		return nil
	} else {
		// This is an update operation
		filter := bson.M{"_id": email.ID}
		result, err := m.collection.ReplaceOne(ctx, filter, email)
		if err != nil {
			return err
		}

		// Check if document was found and updated
		if result.MatchedCount == 0 {
			return ErrEmailNotFound
		}

		return nil
	}
} 