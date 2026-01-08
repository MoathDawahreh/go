package media

import (
	"context"
	"sync"

	appErr "example.com/myapp/internal/errors"
)

// InMemoryRepository is an in-memory implementation of Repository
type InMemoryRepository struct {
	mu    sync.RWMutex
	media map[string]*Media
}

// NewInMemoryRepository creates a new in-memory media repository
func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		media: make(map[string]*Media),
	}
}

// Save stores a media file
func (r *InMemoryRepository) Save(ctx context.Context, media *Media) error {
	if err := ctx.Err(); err != nil {
		return appErr.Internal("context cancelled", err)
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if media.ID == "" {
		return appErr.BadRequest("media ID is required")
	}

	r.media[media.ID] = media
	return nil
}

// GetByID retrieves a media file by ID
func (r *InMemoryRepository) GetByID(ctx context.Context, id string) (*Media, error) {
	if err := ctx.Err(); err != nil {
		return nil, appErr.Internal("context cancelled", err)
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	media, exists := r.media[id]
	if !exists {
		return nil, appErr.NotFound("media not found")
	}
	return media, nil
}

// GetAll retrieves all media files
func (r *InMemoryRepository) GetAll(ctx context.Context) ([]*Media, error) {
	if err := ctx.Err(); err != nil {
		return nil, appErr.Internal("context cancelled", err)
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	mediaList := make([]*Media, 0, len(r.media))
	for _, m := range r.media {
		mediaList = append(mediaList, m)
	}
	return mediaList, nil
}

// Delete removes a media file
func (r *InMemoryRepository) Delete(ctx context.Context, id string) error {
	if err := ctx.Err(); err != nil {
		return appErr.Internal("context cancelled", err)
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.media[id]; !exists {
		return appErr.NotFound("media not found")
	}
	delete(r.media, id)
	return nil
}
