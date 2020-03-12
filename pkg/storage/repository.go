package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/ivan-marquez/es-mdb/pkg/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetAllUsers retrieves all users
func (s *Storage) GetAllUsers() ([]*domain.User, error) {
	db := s.mdb.Database(os.Getenv("DB_NAME"))

	var users []*domain.User
	cur, err := db.Collection("user").Find(context.TODO(), bson.D{{}})
	if err != nil {
		return nil, err
	}

	for cur.Next(context.TODO()) {
		var doc domain.User
		err := cur.Decode(&doc)
		if err != nil {
			return nil, err
		}

		users = append(users, &doc)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	cur.Close(context.TODO())

	return users, nil
}

// GetUsersByTerm retrieves users that match specified term
func (s *Storage) GetUsersByTerm(term string) ([]*domain.User, error) {
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

	res, err := s.es.Search(
		s.es.Search.WithContext(context.Background()),
		s.es.Search.WithIndex("users"),
		s.es.Search.WithBody(&buf),
		s.es.Search.WithTrackTotalHits(true),
		s.es.Search.WithPretty(),
	)
	if err != nil {
		return nil, fmt.Errorf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err == nil {
			return nil, fmt.Errorf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}

		return nil, fmt.Errorf("Error parsing the response body: %v", err)
	}

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf("Error parsing the response body: %v", err)
	}

	var users []*domain.User
	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		src := hit.(map[string]interface{})["_source"]
		ID := hit.(map[string]interface{})["_id"].(string)

		docID, err := primitive.ObjectIDFromHex(ID)
		if err != nil {
			return nil, err
		}

		users = append(users, &domain.User{
			ID:        &docID,
			Email:     src.(map[string]interface{})["Email"].(string),
			FirstName: src.(map[string]interface{})["FirstName"].(string),
			LastName:  src.(map[string]interface{})["LastName"].(string),
			Gender:    src.(map[string]interface{})["Gender"].(string),
			IPAddress: src.(map[string]interface{})["IPAddress"].(string),
		})
	}

	return users, nil
}

// UpdateUser updates specified document on ElasticSearch
func (s *Storage) UpdateUser(ID string, u *domain.User) (string, error) {
	body, _ := json.Marshal(doc{
		Doc:    u,
		Upsert: u,
	})

	var i = 1
	req := esapi.UpdateRequest{
		Index:           "users",
		DocumentID:      ID,
		RetryOnConflict: &i,
		Body:            bytes.NewReader(body),
	}

	res, err := req.Do(context.Background(), s.es)
	if err != nil {
		return "", fmt.Errorf("Error sending ES update: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return "", fmt.Errorf("[%s] Error updating document ID=%s", res.Status(), ID)
	}

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return "", fmt.Errorf("Error parsing the response body: %v", err)
	}

	return fmt.Sprintf("Document %s update result: %s", ID, r["result"].(string)), nil
}

// NewStorage creates a Storage instance with
// ElasticSearch and MongoDB configuration
func NewStorage() (*Storage, error) {
	// MongoDB client config
	clientOptions := options.Client().ApplyURI(os.Getenv("MDB_URL"))
	mdb, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, fmt.Errorf("Error creating MongoDB client: %v", err)
	}

	// ElasticSearch client config
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		return nil, fmt.Errorf("Error creating ES client: %v", err)
	}

	return &Storage{mdb, es}, nil
}
