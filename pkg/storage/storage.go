package storage

import (
	"context"
	"fmt"
	"os"

	"github.com/elastic/go-elasticsearch/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Storage struct {
	DBClient *mongo.Client
	ESClient *elasticsearch.Client
}

func New() (*Storage, error) {
	// MongoDB client config
	clientOptions := options.Client().ApplyURI(os.Getenv("MDB_URL"))
	DBClient, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, fmt.Errorf("Error creating MongoDB client: %v", err)
	}

	// ElasticSearch client config
	ESClient, err := elasticsearch.NewDefaultClient()
	if err != nil {
		return nil, fmt.Errorf("Error creating ES client: %v", err)
	}

	return &Storage{
		DBClient,
		ESClient,
	}, nil
}
