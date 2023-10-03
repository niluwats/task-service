package db

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func ConnectDB() *mongo.Client {
	credentials := options.Credential{
		Username: os.Getenv("MONGO_USERNAME"),
		Password: os.Getenv("MONGO_PASSWORD"),
	}

	uri := "mongodb://mongodb:27017"
	// uri := fmt.Sprintf("mongodb://%s:%s@localhost:%s/?authSource=admin", os.Getenv("MONGO_USERNAME"), os.Getenv("MONGO_PASSWORD"), os.Getenv("MONGO_PORT"))

	clientOptions := options.Client().ApplyURI(uri).SetAuth(credentials)

	log.Println("before ctx")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Error connecting to mongo : ", err)
	}

	log.Println("after ctx")
	defer client.Disconnect(ctx)

	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatal("Ping error : ", err)
	}

	return client
}
