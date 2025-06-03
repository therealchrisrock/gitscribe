package domain

import (
	"time"

	"gorm.io/gorm"
)

// RepositoryModel defines the contract that all repository models must implement
type RepositoryModel interface {
	// GetID returns the repository model's primary key
	GetID() string

	// SetID sets the repository model's primary key
	SetID(id string)

	// TableName returns the database table name for this model
	TableName() string
}

// BaseRepositoryModel provides common fields and methods for repository models
// This handles persistence concerns with appropriate GORM tags
type BaseRepositoryModel struct {
	ID        string         `json:"id" gorm:"primaryKey;type:varchar(128)"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// GetID returns the repository model's ID
func (r *BaseRepositoryModel) GetID() string {
	return r.ID
}

// SetID sets the repository model's ID
func (r *BaseRepositoryModel) SetID(id string) {
	r.ID = id
}

// TableName must be implemented by concrete repository models
// This is intentionally not implemented here to force concrete types to define their table names

// DomainMapper provides bidirectional mapping between domain entities and repository models
// Now constrained to only work with types that implement RepositoryModel
type DomainMapper[D Entity, R RepositoryModel] interface {
	// ToRepository converts domain entity to repository model
	ToRepository(domain D) R

	// ToDomain converts repository model to domain entity
	ToDomain(repo R) D

	// ToRepositoryList converts slice of domain entities to repository models
	ToRepositoryList(domains []D) []R

	// ToDomainList converts slice of repository models to domain entities
	ToDomainList(repos []R) []D
}

// BaseDomainMapper provides common mapping utilities
type BaseDomainMapper struct{}

// StringToPointer converts string to *string, returning nil for empty strings
func (BaseDomainMapper) StringToPointer(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// PointerToString converts *string to string, returning empty string for nil
func (BaseDomainMapper) PointerToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// TimeToPointer converts time.Time to *time.Time, returning nil for zero time
func (BaseDomainMapper) TimeToPointer(t time.Time) *time.Time {
	if t.IsZero() {
		return nil
	}
	return &t
}

// PointerToTime converts *time.Time to time.Time, returning zero time for nil
func (BaseDomainMapper) PointerToTime(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}
	return *t
}

// ValidatedMapper wraps a mapper with validation
type ValidatedMapper[D Entity, R RepositoryModel] struct {
	mapper    DomainMapper[D, R]
	validator EntityValidator[D]
}

// NewValidatedMapper creates a mapper that validates domain entities
func NewValidatedMapper[D Entity, R RepositoryModel](mapper DomainMapper[D, R], validator EntityValidator[D]) *ValidatedMapper[D, R] {
	return &ValidatedMapper[D, R]{
		mapper:    mapper,
		validator: validator,
	}
}

// ToRepository converts domain entity to repository model with validation
func (vm *ValidatedMapper[D, R]) ToRepository(domain D) (R, error) {
	if err := vm.validator.Validate(domain); err != nil {
		var zero R
		return zero, err
	}
	return vm.mapper.ToRepository(domain), nil
}

// ToDomain converts repository model to domain entity
func (vm *ValidatedMapper[D, R]) ToDomain(repo R) D {
	return vm.mapper.ToDomain(repo)
}

// ToRepositoryList converts slice of domain entities to repository models with validation
func (vm *ValidatedMapper[D, R]) ToRepositoryList(domains []D) ([]R, error) {
	for _, domain := range domains {
		if err := vm.validator.Validate(domain); err != nil {
			return nil, err
		}
	}
	return vm.mapper.ToRepositoryList(domains), nil
}

// ToDomainList converts slice of repository models to domain entities
func (vm *ValidatedMapper[D, R]) ToDomainList(repos []R) []D {
	return vm.mapper.ToDomainList(repos)
}

// EntityValidator validates domain entities
type EntityValidator[D Entity] interface {
	Validate(entity D) error
}

// NoOpValidator provides a no-operation validator
type NoOpValidator[D Entity] struct{}

// Validate performs no validation (always returns nil)
func (NoOpValidator[D]) Validate(entity D) error {
	return nil
}
