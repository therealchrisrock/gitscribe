package dtos

import (
	"time"

	"teammate/server/modules/user/domain/entities"
)

// CreateUserRequest represents the request to create a user
type CreateUserRequest struct {
	ID    string `json:"id" binding:"required"`
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

// UpdateUserRequest represents the request to update a user
type UpdateUserRequest struct {
	Name  *string `json:"name,omitempty"`
	Email *string `json:"email,omitempty" binding:"omitempty,email"`
}

// UserResponse represents the response containing user data
type UserResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UsersListResponse represents the response containing a list of users
type UsersListResponse struct {
	Users []UserResponse `json:"users"`
	Total int64          `json:"total"`
}

// ToUserResponse converts a User entity to UserResponse DTO
func ToUserResponse(user *entities.User) UserResponse {
	return UserResponse{
		ID:        user.GetID(),
		Name:      user.GetName(),
		Email:     user.GetEmail().String(),
		CreatedAt: user.GetCreatedAt(),
		UpdatedAt: user.GetUpdatedAt(),
	}
}

// ToUsersListResponse converts a slice of User entities to UsersListResponse DTO
func ToUsersListResponse(users []*entities.User, total int64) UsersListResponse {
	userResponses := make([]UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = ToUserResponse(user)
	}

	return UsersListResponse{
		Users: userResponses,
		Total: total,
	}
}
