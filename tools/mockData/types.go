package main

import "go.mongodb.org/mongo-driver/bson/primitive"

// User schema
type User struct {
	ID        *primitive.ObjectID `bson:"_id,omitempty" json:"-"`
	FirstName string
	LastName  string
	Email     string
	Gender    string
	IPAddress string
}
