package repository

// Repository represent database quries.
type Repository interface {
	Add(args ...string) error
	Update(args ...string) error
	Delete(name string) (interface{}, error)
	GetByID(args ...string) (interface{}, error)
	GetAll() (interface{}, error)
}
