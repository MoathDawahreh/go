package media

import "context"

// Repository defines the interface for media persistence
type Repository interface {
	Save(ctx context.Context, media *Media) error
	GetByID(ctx context.Context, id string) (*Media, error)
	GetAll(ctx context.Context) ([]*Media, error)
	Delete(ctx context.Context, id string) error
}
