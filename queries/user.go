package queries

import (
	"context"
	"fmt"
	"time"
	"udo-golang/database"
	models "udo-golang/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "users")

func GetUser(c *gin.Context) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var foundUser models.User

	// Retrieve the ID from the context
	id, exists := c.Get("uid") // Updated from "uid" to "id"
	if !exists {
		return foundUser, fmt.Errorf("ID not found in context")
	}

	// Convert the ID to a string
	idStr, ok := id.(string)
	if !ok {
		return foundUser, fmt.Errorf("invalid ID format")
	}

	// Convert the string ID to a MongoDB ObjectID
	objID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return foundUser, fmt.Errorf("invalid ObjectID: %v", err)
	}

	// Query MongoDB using the ObjectID
	err = userCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&foundUser)
	if err != nil {
		return foundUser, fmt.Errorf("user not found: %v", err)
	}

	return foundUser, nil
}

func GetUserByID(id string) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var foundUser models.User

	// Convert the string ID to a MongoDB ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return foundUser, fmt.Errorf("invalid ObjectID: %v", err)
	}

	// Query MongoDB using the ObjectID
	err = userCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&foundUser)
	if err != nil {
		return foundUser, fmt.Errorf("user not found: %v", err)
	}

	return foundUser, nil
}

func GetUserByEmail(email string) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var foundUser models.User

	// Query MongoDB using the email
	err := userCollection.FindOne(ctx, bson.M{"email": email}).Decode(&foundUser)
	if err != nil {
		return foundUser, fmt.Errorf("user not found: %v", err)
	}

	return foundUser, nil
}

func UpdateUser(userId string, user models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	filter := bson.M{"_id": objID}

	update := bson.M{"$set": user}

	result, err := userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update goal description: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("no goal found with the given ID")
	}

	return nil
}

func DeleteUser(userId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	filter := bson.M{"_id": objID}

	result := userCollection.FindOneAndDelete(ctx, filter)
	if err := result.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("user not found: %w", err)
		}
		return err
	}

	return nil
}

func GetUserCount(filter bson.M) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	total, err := userCollection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return int(total), nil
}
