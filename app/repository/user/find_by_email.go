package user

import (
	"GoMail/app/repository/models"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// FindByEmail retrieves a user by email
func (r *repository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	collection := r.db.Collection(userCollection)
	
	var user models.User
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	
	return &user, nil
} 