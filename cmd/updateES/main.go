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

		var ds map[string]interface{}
		mdbu := storage.MDBUser(*user)
		decoded := mdbu.DecodeStream(cs, ds)
		doc := domain.User(*decoded.Doc.(*storage.MDBUser))

		res, err := store.UpdateUser(decoded.ID, &doc)
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
