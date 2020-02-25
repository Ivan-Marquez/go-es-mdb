package main

import (
	"context"
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
	store, err := storage.New()
	defer store.DBClient.Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	}

	users, err := getMockData()
	if err != nil {
		log.Fatal(err)
	}

	err = importUsersDataMDB(store, users)
	if err != nil {
		log.Fatal(err)
	}

	err = importUsersDataES(store)
	if err != nil {
		log.Fatal(err)
	}
}
