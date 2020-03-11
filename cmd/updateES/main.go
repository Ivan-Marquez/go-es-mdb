package main

import (
	"log"

	"github.com/ivan-marquez/es-mdb/pkg/domain"
	"github.com/ivan-marquez/es-mdb/pkg/storage"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	store, err := storage.NewStorage()
	if err != nil {
		log.Fatal(err)
	}

	col := "user"
	sh := func(cs *mongo.ChangeStream) {
		user := new(domain.User)

		var ds storage.DecodedStream
		mdbu := storage.MDBUser(*user)
		decoded := mdbu.DecodeStream(cs, ds)

		res, err := store.UpdateUser(decoded.ID, decoded.Doc.(*domain.User))
		if err != nil {
			log.Printf("Error updating ES: %v", err)
			log.Println("Storing failed update on database")
			r, err := store.HandleFailedUpdates(decoded)
			if err != nil {
				log.Fatal(err)
			}

			log.Printf("%s", r.InsertedID.(string))
		}

		log.Println(res)
	}

	store.ChangeStreamListener(col, sh)
}
