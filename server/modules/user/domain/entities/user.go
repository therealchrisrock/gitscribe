package entities

import (
	"teammate/server/seedwork/domain"
)

// User represents a user entity in the domain
type User struct {
	domain.BaseEntity
	Name  string `json:"name" binding:"required" gorm:"column:name"`
	Email Email  `json:"email" binding:"required" gorm:"column:email"`
}

// NewUser creates a new User entity
func NewUser(id, name string, email Email) User {
	user := User{
		Name:  name,
		Email: email,
	}
	user.SetID(id)
	return user
}

// GetEmail returns the user's email
func (u *User) GetEmail() Email {
	return u.Email
}

// SetEmail sets the user's email
func (u *User) SetEmail(email Email) {
	u.Email = email
}

// GetName returns the user's name
func (u *User) GetName() string {
	return u.Name
}

// SetName sets the user's name
func (u *User) SetName(name string) {
	u.Name = name
}

// TableName sets the table name for GORM
func (User) TableName() string {
	return "users"
}
