package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ Repository = (*repository)(nil)

type repository struct {
	col *mongo.Collection
}

func NewUserRepository(db *mongo.Database) Repository {
	return &repository{col: db.Collection("user")}
}

func (r *repository) GetAllUsers(ctx context.Context, page int, pageSize int, filter bson.M) ([]User, error) {

	skip := (page - 1) * pageSize

	opts := options.Find().
		SetSort(bson.M{"createdAt": -1}).
		SetSkip(int64(skip)).
		SetLimit(int64(pageSize))

	cursor, err := r.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users: %v", err)
	}
	defer cursor.Close(ctx)

	var users []User
	if err = cursor.All(ctx, &users); err != nil {
		return nil, fmt.Errorf("failed to decode users: %v", err)
	}

	return users, nil
}

func (r *repository) GetUserByID(ctx context.Context, id string) (*User, error) {
	var foundUser User
	if err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&foundUser); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, err
		}
		return nil, err
	}
	user := foundUser
	return &user, nil
}

func (r *repository) Create(ctx context.Context, u User) (*User, error) {
	doc := u
	res, err := r.col.InsertOne(ctx, doc)
	if err != nil {
		return nil, err
	}
	if uid, ok := res.InsertedID.(primitive.ObjectID); ok {
		doc.ID = uid
	}
	out := doc
	return &out, nil
}

func (r *repository) GetByEmail(ctx context.Context, email string) (*User, error) {
	var foundUser User
	fmt.Println(email)
	if err := r.col.FindOne(ctx, bson.M{"email": email}).Decode(&foundUser); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, err
		}
		return nil, fmt.Errorf("failed to query user: %w", err)
	}
	return &foundUser, nil
}

func (r *repository) UpdateUser(ctx context.Context, userId string, update bson.M) (*User, error) {

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var d User
	if err := r.col.FindOneAndUpdate(ctx, bson.M{"_id": userId}, update, opts).Decode(&d); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, err
		}
		return nil, err
	}
	out := d
	return &out, nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	res, err := r.col.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return err
	}
	return nil
}

func (r *repository) GetUserCount(ctx context.Context, filter bson.M) (int, error) {

	count, err := r.col.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}
	return int(count), nil
}

func (r *repository) GetGoogleUserInfo(ctx context.Context, accessToken string) (map[string]interface{}, error) {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + accessToken)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userInfo map[string]interface{}
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, err
	}

	return userInfo, nil
}
