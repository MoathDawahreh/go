package users

// Repository defines the interface for user persistence
type Repository interface {
	Create(user *User) error
	GetByID(id int) (*User, error)
	GetAll() ([]*User, error)
	Update(user *User) error
	Delete(id int) error
}
