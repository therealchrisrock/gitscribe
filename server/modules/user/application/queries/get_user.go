package queries

import (
	"teammate/server/modules/user/domain/entities"
	"teammate/server/modules/user/domain/repositories"
	"teammate/server/seedwork/domain"
)

// GetUserQuery represents a query to get a user by ID
type GetUserQuery struct {
	ID string
}

// GetUserByEmailQuery represents a query to get a user by email
type GetUserByEmailQuery struct {
	Email string
}

// GetUserHandler handles user queries
type GetUserHandler struct {
	userRepo repositories.UserRepository
}

// NewGetUserHandler creates a new get user handler
func NewGetUserHandler(userRepo repositories.UserRepository) *GetUserHandler {
	return &GetUserHandler{
		userRepo: userRepo,
	}
}

// HandleGetUser executes the get user by ID query
func (h *GetUserHandler) HandleGetUser(query GetUserQuery) (*entities.User, error) {
	user, err := h.userRepo.FindByID(query.ID)
	if err != nil {
		return nil, domain.NewDomainError("USER_NOT_FOUND", "User not found", domain.ErrNotFound)
	}

	return user, nil
}

// HandleGetUserByEmail executes the get user by email query
func (h *GetUserHandler) HandleGetUserByEmail(query GetUserByEmailQuery) (*entities.User, error) {
	user, err := h.userRepo.FindByEmail(query.Email)
	if err != nil {
		return nil, domain.NewDomainError("USER_NOT_FOUND", "User not found", domain.ErrNotFound)
	}

	return user, nil
}
