package users

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// Create a new user
func (s *Service) CreateUser(req *CreateUserRequest) (*User, error) {
	user := &User{
		Name:  req.Name,
		Email: req.Email,
		Age:   req.Age,
	}
	err := s.repo.Create(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// Get user by ID
func (s *Service) GetUser(id int) (*User, error) {
	return s.repo.GetByID(id)
}

// Get all users
func (s *Service) GetAllUsers() ([]*User, error) {
	return s.repo.GetAll()
}

// Update user
func (s *Service) UpdateUser(id int, req *UpdateUserRequest) (*User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	user.Name = req.Name
	user.Email = req.Email
	user.Age = req.Age

	err = s.repo.Update(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// Delete user
func (s *Service) DeleteUser(id int) error {
	return s.repo.Delete(id)
}
