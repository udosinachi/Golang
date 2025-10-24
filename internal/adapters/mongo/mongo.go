package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Client struct {
	DB *mongo.Database
}

func Connect(ctx context.Context, uri, dbName string) (*Client, error) {
	opts := options.Client().ApplyURI(uri)
	c, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}
	if err := c.Ping(ctx, nil); err != nil {
		return nil, err
	}
	return &Client{DB: c.Database(dbName)}, nil
}
