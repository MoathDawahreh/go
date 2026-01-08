package users

import (
	"context"
	"sync"

	appErr "example.com/myapp/internal/errors"
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
		repo.Create(context.Background(), user)
	}

	return repo
}

// Create adds a new user to the repository
func (r *InMemoryRepository) Create(ctx context.Context, user *User) error {
	// Check context cancellation
	if err := ctx.Err(); err != nil {
		return appErr.Internal("context cancelled", err)
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	user.ID = r.nextID
	r.users[r.nextID] = user
	r.nextID++
	return nil
}

// GetByID retrieves a user by ID
func (r *InMemoryRepository) GetByID(ctx context.Context, id int) (*User, error) {
	if err := ctx.Err(); err != nil {
		return nil, appErr.Internal("context cancelled", err)
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, appErr.NotFound("user not found")
	}
	return user, nil
}

// GetAll retrieves all users
func (r *InMemoryRepository) GetAll(ctx context.Context) ([]*User, error) {
	if err := ctx.Err(); err != nil {
		return nil, appErr.Internal("context cancelled", err)
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]*User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}
	return users, nil
}

// Update modifies an existing user
func (r *InMemoryRepository) Update(ctx context.Context, user *User) error {
	if err := ctx.Err(); err != nil {
		return appErr.Internal("context cancelled", err)
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[user.ID]; !exists {
		return appErr.NotFound("user not found")
	}
	r.users[user.ID] = user
	return nil
}

// Delete removes a user from the repository
func (r *InMemoryRepository) Delete(ctx context.Context, id int) error {
	if err := ctx.Err(); err != nil {
		return appErr.Internal("context cancelled", err)
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[id]; !exists {
		return appErr.NotFound("user not found")
	}
	delete(r.users, id)
	return nil
}
