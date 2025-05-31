package repositories

import (
	"teammate/server/modules/user/domain/entities"
)

// UserRepository defines the contract for user repository implementations
type UserRepository interface {
	FindAll() ([]*entities.User, error)
	FindByID(id string) (*entities.User, error)
	FindByEmail(email string) (*entities.User, error)
	Create(user *entities.User) error
	Update(user *entities.User) error
	Delete(id string) error
	HardDelete(id string) error
	Count() (int64, error)
}
