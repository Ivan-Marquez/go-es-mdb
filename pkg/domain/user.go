package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

// User entity
type User struct {
	ID        *primitive.ObjectID `bson:"_id,omitempty" json:"-"` // record ID
	FirstName string              `bson:"firstName"`              // first name
	LastName  string              `bson:"lastName"`               // last name
	Email     string              `bson:"email"`                  // email
	Gender    string              `bson:"gender"`                 // gender
	IPAddress string              `bson:"ipAddress"`              // IP Address
}
