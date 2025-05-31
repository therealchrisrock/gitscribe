package commands

import (
	"teammate/server/modules/user/domain/entities"
	"teammate/server/modules/user/domain/repositories"
	"teammate/server/seedwork/domain"
)

// UpdateUserCommand represents a command to update a user
type UpdateUserCommand struct {
	ID    string
	Name  *string
	Email *string
}

// UpdateUserHandler handles the update user command
type UpdateUserHandler struct {
	userRepo repositories.UserRepository
}

// NewUpdateUserHandler creates a new update user handler
func NewUpdateUserHandler(userRepo repositories.UserRepository) *UpdateUserHandler {
	return &UpdateUserHandler{
		userRepo: userRepo,
	}
}

// Handle executes the update user command
func (h *UpdateUserHandler) Handle(cmd UpdateUserCommand) (*entities.User, error) {
	// Find existing user
	user, err := h.userRepo.FindByID(cmd.ID)
	if err != nil {
		return nil, domain.NewDomainError("USER_NOT_FOUND", "User not found", domain.ErrNotFound)
	}

	// Update name if provided
	if cmd.Name != nil {
		user.SetName(*cmd.Name)
	}

	// Update email if provided
	if cmd.Email != nil {
		// Check if email is already in use by another user
		existingUser, err := h.userRepo.FindByEmail(*cmd.Email)
		if err == nil && existingUser.GetID() != cmd.ID {
			return nil, domain.NewDomainError("EMAIL_ALREADY_EXISTS", "Email already in use", domain.ErrAlreadyExists)
		}

		// Create email value object
		email, err := entities.NewEmail(*cmd.Email)
		if err != nil {
			return nil, domain.NewDomainError("INVALID_EMAIL", "Invalid email format", err)
		}

		user.SetEmail(email)
	}

	// Save user
	if err := h.userRepo.Update(user); err != nil {
		return nil, domain.NewDomainError("UPDATE_USER_FAILED", "Failed to update user", err)
	}

	return user, nil
}
