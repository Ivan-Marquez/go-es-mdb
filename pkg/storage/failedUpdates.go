package storage

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// HandleFailedUpdates sends ES failed updates to a separate MongoDB collection
func (s *Storage) HandleFailedUpdates(decoded *DecodeResult) (*mongo.InsertOneResult, error) {
	dbName := os.Getenv("DB_NAME")
	collection := s.mdb.Database(dbName).Collection("es_failed_updates")
	ctx := context.Background()

	doc := bson.D{
		primitive.E{Key: "es_id", Value: decoded.ID},
		primitive.E{Key: "doc", Value: decoded.Doc},
	}

	res, err := collection.InsertOne(ctx, doc)
	if err != nil {
		return nil, fmt.Errorf("Error while inserting ES failed update: %v", err)
	}

	// log.Printf("ES failed update inserted ID: %s", res.InsertedID.(primitive.ObjectID).Hex())
	return res, nil
}
