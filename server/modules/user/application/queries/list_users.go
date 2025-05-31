package queries

import (
	"teammate/server/modules/user/domain/entities"
	"teammate/server/modules/user/domain/repositories"
	"teammate/server/seedwork/domain"
)

// ListUsersQuery represents a query to list all users
type ListUsersQuery struct {
	// Add pagination parameters in the future
}

// ListUsersResult represents the result of listing users
type ListUsersResult struct {
	Users []*entities.User
	Total int64
}

// ListUsersHandler handles the list users query
type ListUsersHandler struct {
	userRepo repositories.UserRepository
}

// NewListUsersHandler creates a new list users handler
func NewListUsersHandler(userRepo repositories.UserRepository) *ListUsersHandler {
	return &ListUsersHandler{
		userRepo: userRepo,
	}
}

// Handle executes the list users query
func (h *ListUsersHandler) Handle(query ListUsersQuery) (ListUsersResult, error) {
	users, err := h.userRepo.FindAll()
	if err != nil {
		return ListUsersResult{}, domain.NewDomainError("LIST_USERS_FAILED", "Failed to list users", err)
	}

	total, err := h.userRepo.Count()
	if err != nil {
		return ListUsersResult{}, domain.NewDomainError("COUNT_USERS_FAILED", "Failed to count users", err)
	}

	return ListUsersResult{
		Users: users,
		Total: total,
	}, nil
}
