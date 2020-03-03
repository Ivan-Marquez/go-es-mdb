package main

import (
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// DecodedStream represents a MongoDB change stream
type DecodedStream map[string]interface{}

// DecodeResult represents the result of decoding
// a MongoDB change stream
type DecodeResult struct {
	ID  string
	Doc interface{}
}

// Decoder interface
type Decoder interface {
	DecodeStream(s *mongo.ChangeStream, d DecodedStream) *DecodeResult // method to decode a MongoDB change stream
}

// StreamHandler handles MongoDB change stream
type StreamHandler func(cs *mongo.ChangeStream)

// User schema
type User struct {
	FirstName string `bson:"firstName"`
	LastName  string `bson:"lastName"`
	Email     string `bson:"email"`
	Gender    string `bson:"gender"`
	IPAddress string `bson:"ipAddress"`
}

// DecodeStream decodes a MongoDB change stream
func (u *User) DecodeStream(
	s *mongo.ChangeStream,
	d DecodedStream,
) *DecodeResult {
	if e := s.Decode(&d); e != nil {
		log.Printf("error decoding: %s", e)
	}

	fd, ok := d["fullDocument"]

	if !ok {
		panic("fullDocument field on change stream required for this operation")
	}

	docID := fd.(DecodedStream)["_id"].(primitive.ObjectID).Hex()

	u = &User{
		FirstName: fd.(DecodedStream)["firstName"].(string),
		LastName:  fd.(DecodedStream)["lastName"].(string),
		Email:     fd.(DecodedStream)["email"].(string),
		Gender:    fd.(DecodedStream)["gender"].(string),
		IPAddress: fd.(DecodedStream)["ipAddress"].(string),
	}

	return &DecodeResult{
		ID:  docID,
		Doc: u,
	}
}
