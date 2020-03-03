package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ChangeStreamListener listens for changes in the user collection
// and executes the StreamHandler callback
func ChangeStreamListener(dbClient *mongo.Client, colName string, pl mongo.Pipeline, sh StreamHandler) error {
	ctx := context.TODO()

	dbName := os.Getenv("DB_NAME")
	collection := dbClient.Database(dbName).Collection(colName)
	streamOptions := options.ChangeStream().SetFullDocument(options.UpdateLookup)

	cs, err := collection.Watch(ctx, pl, streamOptions)
	if err != nil {
		return fmt.Errorf("Error while watching %s collection: %v", colName, err)
	}

	log.Println("waiting for changes")
	for cs.Next(ctx) {
		sh(cs)
	}

	return nil
}
