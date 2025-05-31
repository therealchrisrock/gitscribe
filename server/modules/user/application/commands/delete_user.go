package commands

import (
	"teammate/server/modules/user/domain/repositories"
	"teammate/server/seedwork/domain"
)

// DeleteUserCommand represents a command to delete a user
type DeleteUserCommand struct {
	ID   string
	Hard bool // If true, performs hard delete
}

// DeleteUserHandler handles the delete user command
type DeleteUserHandler struct {
	userRepo repositories.UserRepository
}

// NewDeleteUserHandler creates a new delete user handler
func NewDeleteUserHandler(userRepo repositories.UserRepository) *DeleteUserHandler {
	return &DeleteUserHandler{
		userRepo: userRepo,
	}
}

// Handle executes the delete user command
func (h *DeleteUserHandler) Handle(cmd DeleteUserCommand) error {
	// Check if user exists
	_, err := h.userRepo.FindByID(cmd.ID)
	if err != nil {
		return domain.NewDomainError("USER_NOT_FOUND", "User not found", domain.ErrNotFound)
	}

	// Perform delete
	if cmd.Hard {
		if err := h.userRepo.HardDelete(cmd.ID); err != nil {
			return domain.NewDomainError("HARD_DELETE_USER_FAILED", "Failed to hard delete user", err)
		}
	} else {
		if err := h.userRepo.Delete(cmd.ID); err != nil {
			return domain.NewDomainError("DELETE_USER_FAILED", "Failed to delete user", err)
		}
	}

	return nil
}
