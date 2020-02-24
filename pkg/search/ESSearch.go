package search

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
)

// User schema
type User struct {
	FirstName string
	LastName  string
	Email     string
	Gender    string
	IPAddress string
}

// ESSearch performs a query string search on ElasticSearch
func ESSearch(client *elasticsearch.Client, term string) ([]*User, error) {
	var (
		buf bytes.Buffer
		r   map[string]interface{}
	)

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"query_string": map[string]interface{}{
				"fields": []string{
					"Email", "FirstName", "LastName", "IPAddress",
				},
				"query": term,
			},
		},
	}

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, fmt.Errorf("Error encoding query: %s", err)
	}

	res, err := client.Search(
		client.Search.WithContext(context.Background()),
		client.Search.WithIndex("users"),
		client.Search.WithBody(&buf),
		client.Search.WithTrackTotalHits(true),
		client.Search.WithPretty(),
	)
	if err != nil {
		return nil, fmt.Errorf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err == nil {
			log.Fatalf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}

		return nil, fmt.Errorf("Error parsing the response body: %s", err)
	}

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}

	// TODO: convert hits to User struct
	// for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
	// 	log.Printf(" * ID=%s, %s", hit.(map[string]interface{})["_id"], hit.(map[string]interface{})["_source"])
	// }

	log.Println(strings.Repeat("=", 37))

	return []*User{}, nil
}
