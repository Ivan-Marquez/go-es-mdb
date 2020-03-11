package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/ivan-marquez/es-mdb/pkg/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// ImportUsersDataES imports data from MongoDB
// to "users" ElasticSearch index
func (s *Storage) ImportUsersDataES(users []*domain.User) error {
	var (
		r  DecodedStream
		wg sync.WaitGroup
	)

	res, err := s.es.Info()
	if err != nil {
		return fmt.Errorf("Error getting ElasticSearch info: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("Error: %s", res.String())
	}
	// Deserialize the response into a map
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return fmt.Errorf("Error parsing the response body: %s", err)
	}

	for i, user := range users {
		wg.Add(1)

		go func(i int, u *domain.User) {
			defer wg.Done()
			// TODO: improve error handling within goroutines
			// https://www.atatus.com/blog/goroutines-error-handling/
			s.indexDocument(u.ID.Hex(), u)
		}(i, user)
	}

	wg.Wait()

	return nil
}

// indexDocument indexes the passed document in ElasticSearch
func (s *Storage) indexDocument(ID string, user *domain.User) (DecodedStream, error) {
	// Build the request body.
	body, _ := json.Marshal(user)

	// Set up the request object.
	req := esapi.IndexRequest{
		Index:      "users",
		DocumentID: ID,
		Body:       bytes.NewReader(body),
		Refresh:    "true",
	}

	// Perform the request with the client
	res, err := req.Do(context.Background(), s.es)
	if err != nil {
		return nil, fmt.Errorf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("[%s] Error indexing document ID=%s", res.Status(), ID)
	}

	// Deserialize the response into a DecodedStream
	var r DecodedStream
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf("Error parsing the response body: %s", err)
	}

	return r, nil
}

// ImportUsersDataMDB imports data from CSV file to MongoDB
func (s *Storage) ImportUsersDataMDB(data []*domain.User) (*mongo.BulkWriteResult, error) {
	dbName := os.Getenv("DB_NAME")

	db := s.mdb.Database(dbName)

	var inserts []mongo.WriteModel

	for _, u := range data {
		// TODO: pending review
		doc := mongo.NewInsertOneModel().SetDocument(bson.M{
			"firstName": u.FirstName,
			"lastName":  u.LastName,
			"email":     u.Email,
			"gender":    u.Gender,
			"ipAddress": u.IPAddress,
		})

		inserts = append(inserts, doc)
	}

	_, err := db.Collection("user").DeleteMany(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}

	bwr, err := db.Collection("user").BulkWrite(context.Background(), inserts)
	if err != nil {
		return nil, err
	}

	return bwr, nil
}
