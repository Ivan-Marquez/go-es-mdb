package main

import (
	"context"
	"log"
	"os"

	"github.com/ivan-marquez/es-mdb/pkg/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func importUsersDataMDB(store *storage.Storage, data []*User) error {
	dbName := os.Getenv("DB_NAME")

	client := store.DBClient

	db := client.Database(dbName)
	defer client.Disconnect(context.TODO())

	var inserts []mongo.WriteModel

	for _, u := range data {
		doc := mongo.NewInsertOneModel().SetDocument(bson.M{
			"firstName": u.FirstName,
			"lastName":  u.LastName,
			"email":     u.Email,
			"gender":    u.Gender,
			"ipAddress": u.IPAddress,
		})

		inserts = append(inserts, doc)
	}

	log.Printf("Importing mock data to MongoDB")
	_, err := db.Collection("user").DeleteMany(context.Background(), bson.M{})
	if err != nil {
		return err
	}

	bwr, err := db.Collection("user").BulkWrite(context.Background(), inserts)
	if err != nil {
		return err
	}

	log.Printf("Data successfully imported: %d", bwr.InsertedCount)

	return nil
}
