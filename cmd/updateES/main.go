package main

import (
	"log"

	"github.com/ivan-marquez/es-mdb/pkg/storage"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	store, err := storage.New()
	if err != nil {
		log.Fatal(err)
	}

	col := "user"
	pipeline := mongo.Pipeline{
		bson.D{primitive.E{
			Key: "$match", Value: bson.D{primitive.E{
				Key: "operationType", Value: bson.D{primitive.E{
					Key: "$in", Value: []string{"insert", "update"},
				}},
			}},
		}},
	}

	sh := func(cs *mongo.ChangeStream) {
		user := new(User)
		var ds DecodedStream
		decoded := user.DecodeStream(cs, ds)

		ch := make(chan *opResult, 1)
		go updateESUser(store.ESClient, decoded.ID, decoded.Doc.(*User), ch)

		res := <-ch
		if res.err == nil {
			log.Printf("Error updating ES: %v", res.err)
			log.Println("Storing failed update on database")
			HandleFailedUpdates(store.DBClient, decoded)
		}

		log.Println(res.status)
	}

	ChangeStreamListener(store.DBClient, col, pipeline, sh)
}
