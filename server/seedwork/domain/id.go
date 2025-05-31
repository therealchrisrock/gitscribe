package domain

import (
	"github.com/google/uuid"
)

// GenerateID generates a new UUID string for entity IDs
func GenerateID() string {
	return uuid.New().String()
}

// IsValidID checks if a string is a valid UUID
func IsValidID(id string) bool {
	_, err := uuid.Parse(id)
	return err == nil
}
