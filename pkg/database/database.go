package database

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type EnvConfig func(key string) string

func New(config EnvConfig) (*mongo.Database, error) {

	host := config("DATABASE_HOST")
	port := config("DATABASE_PORT")

	psn := fmt.Sprintf("mongodb://%s:%s", host, port)

	clientOptions := options.Client().ApplyURI(psn)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	return client.Database("notes-app"), nil
}

func Close(db *mongo.Database) {
	db.Client().Disconnect(context.TODO())
}
