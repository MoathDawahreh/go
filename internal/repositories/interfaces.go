package repositories

import "example.com/myapp/internal/models"

// UserRepository defines the interface for user persistence
type UserRepository interface {
	Create(user *models.User) error
	GetByID(id int) (*models.User, error)
	GetAll() ([]*models.User, error)
	Update(user *models.User) error
	Delete(id int) error
}

// MediaRepository defines the interface for media persistence
type MediaRepository interface {
	Save(media *models.Media) error
	GetByID(id string) (*models.Media, error)
	GetAll() ([]*models.Media, error)
	Delete(id string) error
}
