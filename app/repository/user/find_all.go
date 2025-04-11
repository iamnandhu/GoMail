package user

import (
	"GoMail/app/repository/models"
	"GoMail/app/utils"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// FindAll retrieves users from the database
func (r *repository) FindAll(ctx context.Context, filter interface{}, page, limit int) ([]*models.User, int64, error) {
	collection := r.db.Collection(userCollection)
	
	// Default values
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	
	skip := int64((page - 1) * limit)
	
	// Sanitize the filter
	sanitizedFilter, err := utils.SanitizeMongoFilter(filter)
	if err != nil {
		log.Printf("MongoDB filter sanitization error: %v", err)
		return nil, 0, err
	}
	
	// Get total count
	total, err := collection.CountDocuments(ctx, sanitizedFilter)
	if err != nil {
		return nil, 0, err
	}
	
	// Set options
	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(skip)
	findOptions.SetSort(bson.M{"createdAt": -1})
	
	// Execute query
	cursor, err := collection.Find(ctx, sanitizedFilter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)
	
	var users []*models.User
	if err = cursor.All(ctx, &users); err != nil {
		return nil, 0, err
	}
	
	return users, total, nil
} 