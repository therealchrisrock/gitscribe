package domain

// Repository represents the base interface for all repositories
type Repository[T Entity] interface {
	FindAll() ([]T, error)
	FindByID(id string) (T, error)
	Create(entity *T) error
	Update(entity *T) error
	Delete(id string) error
	Count() (int64, error)
}

// ReadRepository represents read-only operations
type ReadRepository[T Entity] interface {
	FindAll() ([]T, error)
	FindByID(id string) (T, error)
	Count() (int64, error)
}

// WriteRepository represents write operations
type WriteRepository[T Entity] interface {
	Create(entity *T) error
	Update(entity *T) error
	Delete(id string) error
}
