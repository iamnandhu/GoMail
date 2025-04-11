package user

import (
	"GoMail/app/repository/models"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Save saves a user to the database
func (r *repository) Save(ctx context.Context, user *models.User) error {
	collection := r.db.Collection(userCollection)
	
	if user.ID.IsZero() {
		user.CreatedAt = time.Now()
		user.ID = primitive.NewObjectID()
	}
	
	user.UpdatedAt = time.Now()
	
	filter := bson.M{"_id": user.ID}
	update := bson.M{"$set": user}
	opts := options.Update().SetUpsert(true)
	
	_, err := collection.UpdateOne(ctx, filter, update, opts)
	return err
} 