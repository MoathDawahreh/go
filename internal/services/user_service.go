package services

import (
	"example.com/myapp/internal/models"
	"example.com/myapp/internal/repositories"
)

type UserService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

// Create a new user
func (s *UserService) CreateUser(req *models.CreateUserRequest) (*models.User, error) {
	user := &models.User{
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
func (s *UserService) GetUser(id int) (*models.User, error) {
	return s.repo.GetByID(id)
}

// Get all users
func (s *UserService) GetAllUsers() ([]*models.User, error) {
	return s.repo.GetAll()
}

// Update user
func (s *UserService) UpdateUser(id int, req *models.UpdateUserRequest) (*models.User, error) {
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
func (s *UserService) DeleteUser(id int) error {
	return s.repo.Delete(id)
}
