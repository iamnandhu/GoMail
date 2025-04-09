package emaillog

import (
	"context"
	"errors"

	"GoMail/app/repository/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// FindAll retrieves email logs based on filter and pagination
func (db *mongoDB) FindAll(ctx context.Context, filter interface{}, page, limit int) ([]*models.EmailLog, int64, error) {
	// Set default values if not provided
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	// Calculate skip value for pagination
	skip := int64((page - 1) * limit)

	// Set find options
	findOptions := options.Find()
	findOptions.SetSkip(skip)
	findOptions.SetLimit(int64(limit))
	findOptions.SetSort(bson.M{"sent_at": -1}) // Sort by most recent first

	// If filter is nil, use empty filter
	if filter == nil {
		filter = bson.M{}
	}

	// Execute query
	cursor, err := db.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	// Decode results
	logs := make([]*models.EmailLog, 0)
	if err := cursor.All(ctx, &logs); err != nil {
		return nil, 0, err
	}

	// Count total documents
	total, err := db.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// FindByID retrieves an email log by ID
func (db *mongoDB) FindByID(ctx context.Context, id string) (*models.EmailLog, error) {
	// Convert string ID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, ErrInvalidID
	}

	// Create filter
	filter := bson.M{"_id": objectID}

	// Execute query
	result := db.collection.FindOne(ctx, filter)
	if err := result.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrEmailLogNotFound
		}
		return nil, err
	}

	// Decode result
	emailLog := &models.EmailLog{}
	if err := result.Decode(emailLog); err != nil {
		return nil, err
	}

	return emailLog, nil
} 