package domain

// UserService interface
type UserService interface {
	GetAll() ([]*User, error)                  // GetAll returns all users
	GetByTerm(term string) ([]*User, error)    // GetByTerm returns users that match specified term
	Update(ID string, u *User) (string, error) // Update updates user by specified ID and doc
}

type userService struct {
	r UserRepository
}

func (us *userService) GetAll() ([]*User, error) {
	users, err := us.r.GetAllUsers()
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (us *userService) GetByTerm(term string) ([]*User, error) {
	users, err := us.r.GetUsersByTerm(term)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (us *userService) Update(ID string, u *User) (string, error) {
	res, err := us.r.UpdateUser(ID, u)
	if err != nil {
		return "", err
	}

	return res, nil
}

// NewUserService creates a user service with repository implementation
func NewUserService(ur UserRepository) UserService {
	return &userService{ur}
}
