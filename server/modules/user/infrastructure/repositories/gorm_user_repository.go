package repositories

import (
	"errors"
	"time"

	"teammate/server/modules/user/domain/entities"
	"teammate/server/modules/user/domain/repositories"
	"teammate/server/seedwork/infrastructure/database"

	"gorm.io/gorm"
)

// GormUserRepository handles database operations for users using GORM
type GormUserRepository struct {
	db *gorm.DB
}

// Ensure GormUserRepository implements UserRepository
var _ repositories.UserRepository = (*GormUserRepository)(nil)

// NewGormUserRepository creates a new GORM-based user repository
func NewGormUserRepository() *GormUserRepository {
	return &GormUserRepository{db: database.GetDB()}
}

// FindAll retrieves all users
func (r *GormUserRepository) FindAll() ([]*entities.User, error) {
	var users []*entities.User
	result := r.db.Find(&users)
	return users, result.Error
}

// FindByID retrieves a user by ID
func (r *GormUserRepository) FindByID(id string) (*entities.User, error) {
	var user entities.User
	result := r.db.First(&user, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// FindByEmail retrieves a user by email
func (r *GormUserRepository) FindByEmail(email string) (*entities.User, error) {
	var user entities.User
	result := r.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// Create creates a new user
func (r *GormUserRepository) Create(user *entities.User) error {
	return r.db.Create(user).Error
}

// Update updates an existing user
func (r *GormUserRepository) Update(user *entities.User) error {
	if user.GetID() == "" {
		return errors.New("user ID is required for update")
	}
	user.UpdatedAt = time.Now()
	return r.db.Save(user).Error
}

// Delete soft-deletes a user by ID
func (r *GormUserRepository) Delete(id string) error {
	return r.db.Delete(&entities.User{}, "id = ?", id).Error
}

// HardDelete permanently deletes a user by ID
func (r *GormUserRepository) HardDelete(id string) error {
	return r.db.Unscoped().Delete(&entities.User{}, "id = ?", id).Error
}

// Count returns the total number of users
func (r *GormUserRepository) Count() (int64, error) {
	var count int64
	result := r.db.Model(&entities.User{}).Count(&count)
	return count, result.Error
}
