package queries

import (
	"context"
	"errors"
	"fmt"
	"time"
	"udo-golang/database"
	models "udo-golang/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "users")

func newCtx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}

func toObjectID(id string) (primitive.ObjectID, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("invalid ObjectID: %w", err)
	}
	return objID, nil
}

func GetAllUsers(page int, pageSize int, filter bson.M) ([]models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	skip := (page - 1) * pageSize

	opts := options.Find().
		SetSort(bson.M{"createdAt": -1}).
		SetSkip(int64(skip)).
		SetLimit(int64(pageSize))

	cursor, err := userCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users: %v", err)
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err = cursor.All(ctx, &users); err != nil {
		return nil, fmt.Errorf("failed to decode users: %v", err)
	}

	return users, nil
}

func GetUserByID(id string) (*models.User, error) {
	ctx, cancel := newCtx()
	defer cancel()

	objID, err := toObjectID(id)
	if err != nil {
		return nil, err
	}

	var foundUser models.User
	err = userCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&foundUser)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to query user: %w", err)
	}

	return &foundUser, nil
}

func GetUserByEmail(email string) (*models.User, error) {
	ctx, cancel := newCtx()
	defer cancel()

	var foundUser models.User
	err := userCollection.FindOne(ctx, bson.M{"email": email}).Decode(&foundUser)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to query user: %w", err)
	}

	return &foundUser, nil
}

func UpdateUser(userId string, update bson.M) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	filter := bson.M{"_id": objID}
	updateDoc := bson.M{"$set": update}

	result, err := userCollection.UpdateOne(ctx, filter, updateDoc)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("no user found with the given ID")
	}

	return nil
}

func DeleteUserById(userId string) error {
	ctx, cancel := newCtx()
	defer cancel()

	objID, err := toObjectID(userId)
	if err != nil {
		return err
	}

	result := userCollection.FindOneAndDelete(ctx, bson.M{"_id": objID})
	if err := result.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Errorf("user not found")
		}
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

func GetUserCount(filter bson.M) (int, error) {
	ctx, cancel := newCtx()
	defer cancel()

	count, err := userCollection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}
	return int(count), nil
}
