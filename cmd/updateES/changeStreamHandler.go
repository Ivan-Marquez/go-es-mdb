package main

import (
	"context"
	"log"
	"os"

	"github.com/ivan-marquez/es-mdb/pkg/storage"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// changeStreamHandler listens for changes
// in the user collection and sends updates
// to ElasticSearch index
// TODO: refactor to receive updateES() as a callback
func changeStreamHandler(colName string, pl mongo.Pipeline, d Decoder) {
	ctx := context.TODO()
	store, err := storage.New()
	if err != nil {
		log.Fatal(err)
	}

	dbName := os.Getenv("DB_NAME")
	collection := store.DBClient.Database(dbName).Collection(colName)
	streamOptions := options.ChangeStream().SetFullDocument(options.UpdateLookup)

	cs, err := collection.Watch(ctx, pl, streamOptions)
	if err != nil {
		log.Fatal(err)
	}

	ch := make(chan *updateResult, 1)
	var ds DecodedStream

	log.Println("waiting for changes")
	for cs.Next(ctx) {
		decoded := d.DecodeStream(cs, ds)
		go updateES(store.ESClient, decoded.ID, decoded.Doc.(*User), ch)

		res := <-ch
		if res.err != nil {
			// TODO: handle failed updates
			log.Println(res.err)
		}

		log.Println(res.status)
	}
}
