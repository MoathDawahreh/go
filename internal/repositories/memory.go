package repositories

import (
	"errors"
	"sync"

	"example.com/myapp/internal/models"
)

// InMemoryUserRepository is an in-memory implementation of UserRepository
type InMemoryUserRepository struct {
	mu     sync.RWMutex
	users  map[int]*models.User
	nextID int
}

// NewInMemoryUserRepository creates a new in-memory user repository
func NewInMemoryUserRepository() *InMemoryUserRepository {
	repo := &InMemoryUserRepository{
		users:  make(map[int]*models.User),
		nextID: 1,
	}

	// Add 3 dummy users for testing
	dummyUsers := []*models.User{
		{Name: "John Doe", Email: "john@example.com", Age: 30},
		{Name: "Jane Smith", Email: "jane@example.com", Age: 28},
		{Name: "Bob Johnson", Email: "bob@example.com", Age: 35},
	}

	for _, user := range dummyUsers {
		repo.Create(user)
	}

	return repo
}

// Create adds a new user to the repository
func (r *InMemoryUserRepository) Create(user *models.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	user.ID = r.nextID
	r.users[r.nextID] = user
	r.nextID++
	return nil
}

// GetByID retrieves a user by ID
func (r *InMemoryUserRepository) GetByID(id int) (*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

// GetAll retrieves all users
func (r *InMemoryUserRepository) GetAll() ([]*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]*models.User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}
	return users, nil
}

// Update modifies an existing user
func (r *InMemoryUserRepository) Update(user *models.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[user.ID]; !exists {
		return errors.New("user not found")
	}
	r.users[user.ID] = user
	return nil
}

// Delete removes a user from the repository
func (r *InMemoryUserRepository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[id]; !exists {
		return errors.New("user not found")
	}
	delete(r.users, id)
	return nil
}

// InMemoryMediaRepository is an in-memory implementation of MediaRepository
type InMemoryMediaRepository struct {
	mu    sync.RWMutex
	media map[string]*models.Media
}

// NewInMemoryMediaRepository creates a new in-memory media repository
func NewInMemoryMediaRepository() *InMemoryMediaRepository {
	return &InMemoryMediaRepository{
		media: make(map[string]*models.Media),
	}
}

// Save stores a media file
func (r *InMemoryMediaRepository) Save(media *models.Media) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if media.ID == "" {
		return errors.New("media ID is required")
	}

	r.media[media.ID] = media
	return nil
}

// GetByID retrieves a media file by ID
func (r *InMemoryMediaRepository) GetByID(id string) (*models.Media, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	media, exists := r.media[id]
	if !exists {
		return nil, errors.New("media not found")
	}
	return media, nil
}

// GetAll retrieves all media files
func (r *InMemoryMediaRepository) GetAll() ([]*models.Media, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	mediaList := make([]*models.Media, 0, len(r.media))
	for _, m := range r.media {
		mediaList = append(mediaList, m)
	}
	return mediaList, nil
}

// Delete removes a media file
func (r *InMemoryMediaRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.media[id]; !exists {
		return errors.New("media not found")
	}
	delete(r.media, id)
	return nil
}
