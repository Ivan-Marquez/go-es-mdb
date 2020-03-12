package storage

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DecodeStream decodes a MongoDB change stream
func (u *MDBUser) DecodeStream(
	s *mongo.ChangeStream,
	d map[string]interface{},
) *DecodeResult {
	if e := s.Decode(&d); e != nil {
		log.Printf("error decoding: %s", e)
	}

	fd, ok := d["fullDocument"]

	if !ok {
		panic("fullDocument field on change stream required for this operation")
	}

	docID := fd.(map[string]interface{})["_id"].(primitive.ObjectID).Hex()

	u = &MDBUser{
		FirstName: fd.(map[string]interface{})["firstName"].(string),
		LastName:  fd.(map[string]interface{})["lastName"].(string),
		Email:     fd.(map[string]interface{})["email"].(string),
		Gender:    fd.(map[string]interface{})["gender"].(string),
		IPAddress: fd.(map[string]interface{})["ipAddress"].(string),
	}

	return &DecodeResult{
		ID:  docID,
		Doc: u,
	}
}

// ChangeStreamListener listens for changes in the user collection
// and executes the StreamHandler callback
func (s *Storage) ChangeStreamListener(colName string, sh StreamHandler) error {
	ctx := context.TODO()

	dbName := os.Getenv("DB_NAME")
	collection := s.mdb.Database(dbName).Collection(colName)
	streamOptions := options.ChangeStream().SetFullDocument(options.UpdateLookup)

	pipeline := mongo.Pipeline{
		bson.D{primitive.E{
			Key: "$match", Value: bson.D{primitive.E{
				Key: "operationType", Value: bson.D{primitive.E{
					Key: "$in", Value: []string{"insert", "update"},
				}},
			}},
		}},
	}

	cs, err := collection.Watch(ctx, pipeline, streamOptions)
	if err != nil {
		return fmt.Errorf("Error while watching %s collection: %v", colName, err)
	}

	log.Println("waiting for changes")
	for cs.Next(ctx) {
		// TODO: implement as goroutine
		sh(cs)
	}

	return nil
}
