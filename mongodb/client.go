package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*
this parameter should come from config
*/
const (
	DefaultMongodbConnectTimeout = time.Second
)

func GetClient(mongodbURL string, timeout time.Duration) (*mongo.Client, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(mongodbURL))
	if err != nil {
		return nil, fmt.Errorf(`failed to create mongodb client: [url: %s, err: %w]`, mongodbURL, err)
	}
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	err = client.Connect(ctx)
	if err != nil {
		return nil, fmt.Errorf(`failed to connect mongodb: [url: %s, err: %w]`, mongodbURL, err)
	}
	return client, nil
}
