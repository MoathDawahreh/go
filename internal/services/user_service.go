package services

import (
	"errors"

	"example.com/myapp/internal/models"
)

type UserService struct {
	users map[int]*models.User
	nextID int
}

func NewUserService() *UserService {
	return &UserService{
		users: make(map[int]*models.User),
		nextID: 1,
	}
}

// Create a new user
func (s *UserService) CreateUser(req *models.CreateUserRequest) *models.User {
	user := &models.User{
		ID:    s.nextID,
		Name:  req.Name,
		Email: req.Email,
		Age:   req.Age,
	}
	s.users[s.nextID] = user
	s.nextID++
	return user
}

// Get user by ID
func (s *UserService) GetUser(id int) (*models.User, error) {
	user, exists := s.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

// Get all users
func (s *UserService) GetAllUsers() []*models.User {
	users := make([]*models.User, 0, len(s.users))
	for _, user := range s.users {
		users = append(users, user)
	}
	// add a mock user for demonstration
	users = append(users, &models.User{
		ID:    0,
		Name:  "Mock User",
		Email: "Maz@gmail.com",
		Age:   30,
	})


	return users
}

// Update user
func (s *UserService) UpdateUser(id int, req *models.UpdateUserRequest) (*models.User, error) {
	user, exists := s.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	
	user.Name = req.Name
	user.Email = req.Email
	user.Age = req.Age
	
	return user, nil
}

// Delete user
func (s *UserService) DeleteUser(id int) error {
	if _, exists := s.users[id]; !exists {
		return errors.New("user not found")
	}
	delete(s.users, id)
	return nil
}
