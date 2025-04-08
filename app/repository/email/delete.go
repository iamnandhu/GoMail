package email

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Delete removes an email from the database
func (m *mongoDB) Delete(ctx context.Context, id interface{}) error {
	// Convert ID to ObjectID if needed
	var objectID primitive.ObjectID
	switch v := id.(type) {
	case primitive.ObjectID:
		objectID = v
	case string:
		var err error
		objectID, err = primitive.ObjectIDFromHex(v)
		if err != nil {
			return ErrInvalidID
		}
	default:
		return ErrInvalidID
	}

	// Execute delete
	result, err := m.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}

	// Check if document was found and deleted
	if result.DeletedCount == 0 {
		return ErrEmailNotFound
	}

	return nil
} 