package media

// Repository defines the interface for media persistence
type Repository interface {
	Save(media *Media) error
	GetByID(id string) (*Media, error)
	GetAll() ([]*Media, error)
	Delete(id string) error
}
