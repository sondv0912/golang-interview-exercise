package mongodb

import (
	"fmt"
	"golang-interview-exercise/utils/context_utils"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	clientInstance *mongo.Client
	once           sync.Once
	errInit        error
)

func GetMongoClient() (*mongo.Client, error) {
	once.Do(func() {
		clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

		ctx, cancel := context_utils.CreateTimeoutContext(5 * time.Second)
		defer cancel()

		clientInstance, errInit = mongo.Connect(ctx, clientOptions)
		if errInit != nil {
			errInit = fmt.Errorf("failed to connect to MongoDB: %w", errInit)
			return
		}

		errInit = clientInstance.Ping(ctx, nil)
		if errInit != nil {
			errInit = fmt.Errorf("failed to ping MongoDB: %w", errInit)
		}
	})
	return clientInstance, errInit
}

func GetCollection(databaseName, collectionName string) (*mongo.Collection, error) {
	client, err := GetMongoClient()
	if err != nil {
		return nil, fmt.Errorf("error getting MongoDB client: %w", err)
	}

	return client.Database(databaseName).Collection(collectionName), nil
}
