package users

import (
	"context"
	"log/slog"

	appErr "example.com/myapp/internal/errors"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// Create a new user
func (s *Service) CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error) {
	if err := ctx.Err(); err != nil {
		return nil, appErr.Internal("context cancelled", err)
	}

	if req.Name == "" || req.Email == "" {
		return nil, appErr.BadRequest("name and email are required")
	}

	user := &User{
		Name:  req.Name,
		Email: req.Email,
		Age:   req.Age,
	}
	err := s.repo.Create(ctx, user)
	if err != nil {
		slog.Error("Failed to create user", "error", err)
		return nil, appErr.Internal("failed to create user", err)
	}
	return user, nil
}

// Get user by ID
func (s *Service) GetUser(ctx context.Context, id int) (*User, error) {
	if err := ctx.Err(); err != nil {
		return nil, appErr.Internal("context cancelled", err)
	}

	if id <= 0 {
		return nil, appErr.InvalidID("user id must be positive")
	}

	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		slog.Error("Failed to get user", "id", id, "error", err)
		return nil, err
	}
	return user, nil
}

// Get all users
func (s *Service) GetAllUsers(ctx context.Context) ([]*User, error) {
	if err := ctx.Err(); err != nil {
		return nil, appErr.Internal("context cancelled", err)
	}

	users, err := s.repo.GetAll(ctx)
	if err != nil {
		slog.Error("Failed to get all users", "error", err)
		return nil, appErr.Internal("failed to retrieve users", err)
	}
	return users, nil
}

// Update user
func (s *Service) UpdateUser(ctx context.Context, id int, req *UpdateUserRequest) (*User, error) {
	if err := ctx.Err(); err != nil {
		return nil, appErr.Internal("context cancelled", err)
	}

	if id <= 0 {
		return nil, appErr.InvalidID("user id must be positive")
	}

	if req.Name == "" || req.Email == "" {
		return nil, appErr.BadRequest("name and email are required")
	}

	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		slog.Error("Failed to get user for update", "id", id, "error", err)
		return nil, err
	}

	user.Name = req.Name
	user.Email = req.Email
	user.Age = req.Age

	err = s.repo.Update(ctx, user)
	if err != nil {
		slog.Error("Failed to update user", "id", id, "error", err)
		return nil, appErr.Internal("failed to update user", err)
	}
	return user, nil
}

// Delete user
func (s *Service) DeleteUser(ctx context.Context, id int) error {
	if err := ctx.Err(); err != nil {
		return appErr.Internal("context cancelled", err)
	}

	if id <= 0 {
		return appErr.InvalidID("user id must be positive")
	}

	err := s.repo.Delete(ctx, id)
	if err != nil {
		slog.Error("Failed to delete user", "id", id, "error", err)
		return err
	}
	return nil
}
