package main

import (
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// DecodedStream represents a MongoDB change stream
type DecodedStream map[string]interface{}

// TODO: description
type DecodeResult struct {
	ID  string
	Doc interface{}
}

// TODO: description
type Decoder interface {
	DecodeStream(s *mongo.ChangeStream, d DecodedStream) *DecodeResult // TODO: description
}

// User schema
type User struct {
	FirstName string
	LastName  string
	Email     string
	Gender    string
	IPAddress string
}

// TODO: description
func (u *User) DecodeStream(
	s *mongo.ChangeStream,
	d DecodedStream,
) *DecodeResult {
	if e := s.Decode(&d); e != nil {
		log.Printf("error decoding: %s", e)
	}

	fd, ok := d["fullDocument"]

	if !ok {
		// TODO: what if fullDocument is not available?
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
