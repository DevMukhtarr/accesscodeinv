package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client

func ConnectDB() (*mongo.Database, error) {
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_LOCAL_URI"))
	var err error
	// Connect to MongoDB
	mongoClient, err = mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
		fmt.Println(err)
	}

	err = mongoClient.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
		fmt.Println(err)
	}
	fmt.Println("Connected to MongoDB!")

	return mongoClient.Database("taskscheduler"), nil
}

func GetCollection(collectionName string) *mongo.Collection {
	if mongoClient == nil {
		log.Fatal("MongoDB client is not initialized")
	}
	return mongoClient.Database("admin").Collection(collectionName)
}
