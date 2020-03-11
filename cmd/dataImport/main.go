package main

import (
	"log"

	"github.com/ivan-marquez/es-mdb/pkg/storage"
	"github.com/joho/godotenv"
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

	users, err := getMockData()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Inserting data to MongoDB…")

	res, err := store.ImportUsersDataMDB(users)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Insert to MongoDB successful. Inserted: %d\n", res.InsertedCount)

	users, err = store.GetAllUsers()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Inserting data to ElasticSearch…")

	err = store.ImportUsersDataES(users)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Insert to ElasticSearch successful.")
}
