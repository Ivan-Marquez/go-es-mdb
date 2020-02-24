package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"sync"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/ivan-marquez/es-mdb/pkg/storage"
)

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
	log.Printf("Client: %s", elasticsearch.Version)
	log.Printf("Server: %s", r["version"].(map[string]interface{})["number"])

	md, err := getMockData()
	if err != nil {
		return fmt.Errorf("Error getting mock data: %v", err)
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
				DocumentID: strconv.Itoa(i + 1),
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
