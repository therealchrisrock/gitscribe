package domain

import (
	"time"
)

// Entity represents the base interface for all domain entities
type Entity interface {
	GetID() string
	SetID(id string)
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time
}

// BaseEntity provides common fields and methods for all entities
// This is a pure domain object with no persistence concerns
type BaseEntity struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GetID returns the entity ID
func (e *BaseEntity) GetID() string {
	return e.ID
}

// SetID sets the entity ID
func (e *BaseEntity) SetID(id string) {
	e.ID = id
}

// GetCreatedAt returns the creation timestamp
func (e *BaseEntity) GetCreatedAt() time.Time {
	return e.CreatedAt
}

// GetUpdatedAt returns the last update timestamp
func (e *BaseEntity) GetUpdatedAt() time.Time {
	return e.UpdatedAt
}
