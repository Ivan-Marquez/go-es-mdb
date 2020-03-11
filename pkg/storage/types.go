package storage

import (
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/ivan-marquez/es-mdb/pkg/domain"
	"go.mongodb.org/mongo-driver/mongo"
)

// MDBUser extends domain.User to add
// decoding functionality
type MDBUser domain.User

// StreamHandler handles MongoDB change stream
type StreamHandler func(cs *mongo.ChangeStream)

// DecodeResult represents the result of decoding
// a MongoDB change stream
type DecodeResult struct {
	ID  string
	Doc interface{}
}

// Storage type with MongoDB and ElasticSearch config
type Storage struct {
	mdb *mongo.Client
	es  *elasticsearch.Client
}

type doc struct {
	Doc    *domain.User `json:"doc"`
	Upsert *domain.User `json:"upsert"`
}

// DecodedStream for Json decoding
type DecodedStream map[string]interface{}
