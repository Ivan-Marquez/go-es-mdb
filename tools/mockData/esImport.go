package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/elastic/go-elasticsearch/v8"
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

func importUsersDataES(store *storage.Storage) error {
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

	log.Printf("Importing data to ElasticSearch")

	users, err := getUsers(store)
	if err != nil {
		return fmt.Errorf("Error getting data from MongoDB: %v", err)
	}

	for i, user := range users {
		wg.Add(1)

		go func(i int, u *User) {
			defer wg.Done()
			// TODO: improve error handling within goroutines
			// https://www.atatus.com/blog/goroutines-error-handling/
			indexDocument(es, u.ID.Hex(), u)
		}(i, user)
	}

	wg.Wait()
	log.Printf("Data successfully imported to ElasticSearch")

	return nil
}

func indexDocument(es *elasticsearch.Client, ID string, user *User) error {
	// Build the request body.
	body, _ := json.Marshal(user)

	// Set up the request object.
	req := esapi.IndexRequest{
		Index:      "users",
		DocumentID: user.ID.Hex(),
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
		return fmt.Errorf("[%s] Error indexing document ID=%s", res.Status(), ID)
	}

	// Deserialize the response into a map.
	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return fmt.Errorf("Error parsing the response body: %s", err)
	}

	return nil
}
