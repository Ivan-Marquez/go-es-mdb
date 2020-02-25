package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/ivan-marquez/es-mdb/pkg/storage"
	"go.mongodb.org/mongo-driver/bson"
)

func getUsers(store *storage.Storage) ([]*User, error) {
	client := store.DBClient
	db := client.Database(os.Getenv("DB_NAME"))

	var users []*User
	cur, err := db.Collection("user").Find(context.TODO(), bson.D{{}})
	if err != nil {
		return nil, err
	}

	for cur.Next(context.TODO()) {
		var doc User
		err := cur.Decode(&doc)
		if err != nil {
			log.Fatal(err)
		}

		users = append(users, &doc)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	cur.Close(context.TODO())

	return users, nil
}

func importUsersDataES(store *storage.Storage, users []*User) error {
	var (
		r  map[string]interface{}
		wg sync.WaitGroup
	)

	es := store.ESClient

	res, err := es.Info()
	if err != nil {
		return fmt.Errorf("Error getting ElasticSearch info: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("Error: %s", res.String())
	}
	// Deserialize the response into a map.
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return fmt.Errorf("Error parsing the response body: %s", err)
	}
	// Print client and server version numbers.
	log.Printf("Importing data to ElasticSearch")

	md, err := getUsers(store)
	if err != nil {
		return fmt.Errorf("Error getting data from MongoDB: %v", err)
	}

	for i, user := range md {
		wg.Add(1)

		go func(i int, u *User) error {
			defer wg.Done()

			// Build the request body.
			body, _ := json.Marshal(u)

			// Set up the request object.
			req := esapi.IndexRequest{
				Index:      "users",
				DocumentID: u.ID.Hex(),
				Body:       bytes.NewReader(body),
				Refresh:    "true",
			}

			// Perform the request with the client.
			res, err := req.Do(context.Background(), es)
			if err != nil {
				return fmt.Errorf("Error getting response: %s", err)
			}
			defer res.Body.Close()

			if res.IsError() {
				log.Printf("[%s] Error indexing document ID=%d", res.Status(), i+1)
			} else {
				// Deserialize the response into a map.
				var r map[string]interface{}
				if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
					log.Printf("Error parsing the response body: %s", err)
				}
			}

			return nil
		}(i, user)
	}
	wg.Wait()

	log.Printf("Data successfully imported to ElasticSearch")

	return nil
}
