package search

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/elastic/go-elasticsearch/v8"
)

// User schema
type User struct {
	ID        string
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
				"query": fmt.Sprintf("*%s*", term),
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

	var users []*User
	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		src := hit.(map[string]interface{})["_source"]
		ID := hit.(map[string]interface{})["_id"].(string)

		users = append(users, &User{
			ID:        ID,
			Email:     src.(map[string]interface{})["Email"].(string),
			FirstName: src.(map[string]interface{})["FirstName"].(string),
			LastName:  src.(map[string]interface{})["LastName"].(string),
			Gender:    src.(map[string]interface{})["Gender"].(string),
			IPAddress: src.(map[string]interface{})["IPAddress"].(string),
		})
	}

	return users, nil
}
