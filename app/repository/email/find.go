package email

import (
	"context"

	"GoMail/app/repository/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// FindAll retrieves emails with optional filtering and pagination
func (m *mongoDB) FindAll(ctx context.Context, filter interface{}, page, limit int) ([]*models.Email, int64, error) {
	// Default filter
	queryFilter := bson.M{}

	// Apply filter if provided
	if filter != nil {
		// Check if filter is a primitive.ObjectID (for finding by ID)
		if oid, ok := filter.(primitive.ObjectID); ok {
			queryFilter = bson.M{"_id": oid}
		} else if status, ok := filter.(models.EmailStatus); ok {
			// Check if filter is an EmailStatus
			queryFilter = bson.M{"status": status}
		} else if mapFilter, ok := filter.(map[string]interface{}); ok {
			// Check if filter is a map
			queryFilter = bson.M(mapFilter)
		} else if bsonFilter, ok := filter.(bson.M); ok {
			// Check if filter is already a bson.M
			queryFilter = bsonFilter
		}
	}

	// Configure pagination
	skip := int64((page - 1) * limit)
	opts := options.Find().
		SetSkip(skip).
		SetLimit(int64(limit)).
		SetSort(bson.M{"created_at": -1}) // Sort by creation date, newest first

	// Execute query
	cursor, err := m.collection.Find(ctx, queryFilter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	// Decode results
	var emails []*models.Email
	if err := cursor.All(ctx, &emails); err != nil {
		return nil, 0, err
	}

	// Get total count
	total, err := m.collection.CountDocuments(ctx, queryFilter)
	if err != nil {
		return nil, 0, err
	}

	return emails, total, nil
} 