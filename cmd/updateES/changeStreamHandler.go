package main

import (
	"context"
	"log"
	"os"

	"github.com/ivan-marquez/es-mdb/pkg/storage"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func changeStreamHandler() {
	ctx := context.TODO()
	client, err := storage.New()
	if err != nil {
		log.Fatal(err)
	}

	dbName := os.Getenv("DB_NAME")
	collection := client.DBClient.Database(dbName).Collection("user")
	streamOptions := options.ChangeStream().SetFullDocument(options.UpdateLookup)

	stream, err := collection.Watch(ctx, mongo.Pipeline{}, streamOptions)
	if err != nil {
		log.Fatal(err)
	}

	log.Print("waiting for changes")
	var changeDoc map[string]interface{}

	for stream.Next(ctx) {
		if e := stream.Decode(&changeDoc); e != nil {
			log.Printf("error decoding: %s", e)
		}
		log.Printf("change: %+v\n", changeDoc)
	}
}
