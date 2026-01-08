package users

import (
	"errors"
	"sync"
)

// InMemoryRepository is an in-memory implementation of Repository
type InMemoryRepository struct {
	mu     sync.RWMutex
	users  map[int]*User
	nextID int
}

// NewInMemoryRepository creates a new in-memory user repository
func NewInMemoryRepository() *InMemoryRepository {
	repo := &InMemoryRepository{
		users:  make(map[int]*User),
		nextID: 1,
	}

	// Add 3 dummy users for testing
	dummyUsers := []*User{
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
func (r *InMemoryRepository) Create(user *User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	user.ID = r.nextID
	r.users[r.nextID] = user
	r.nextID++
	return nil
}

// GetByID retrieves a user by ID
func (r *InMemoryRepository) GetByID(id int) (*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

// GetAll retrieves all users
func (r *InMemoryRepository) GetAll() ([]*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]*User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}
	return users, nil
}

// Update modifies an existing user
func (r *InMemoryRepository) Update(user *User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[user.ID]; !exists {
		return errors.New("user not found")
	}
	r.users[user.ID] = user
	return nil
}

// Delete removes a user from the repository
func (r *InMemoryRepository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[id]; !exists {
		return errors.New("user not found")
	}
	delete(r.users, id)
	return nil
}
