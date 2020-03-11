package domain

// UserRepository interface with required methods
type UserRepository interface {
	GetAllUsers() ([]*User, error)                 // GetAllUsers returns all users
	GetUsersByTerm(term string) ([]*User, error)   // GetUsersByTerm return users that match term
	UpdateUser(ID string, u *User) (string, error) // UpdateUser updates user with specified ID and doc
}
