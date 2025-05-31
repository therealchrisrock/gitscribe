package commands

import (
	"teammate/server/modules/user/domain/entities"
	"teammate/server/modules/user/domain/repositories"
	"teammate/server/seedwork/domain"
)

// CreateUserCommand represents a command to create a new user
type CreateUserCommand struct {
	ID    string
	Name  string
	Email string
}

// CreateUserHandler handles the create user command
type CreateUserHandler struct {
	userRepo repositories.UserRepository
}

// NewCreateUserHandler creates a new create user handler
func NewCreateUserHandler(userRepo repositories.UserRepository) *CreateUserHandler {
	return &CreateUserHandler{
		userRepo: userRepo,
	}
}

// Handle executes the create user command
func (h *CreateUserHandler) Handle(cmd CreateUserCommand) (*entities.User, error) {
	// Check if user already exists
	existingUser, err := h.userRepo.FindByID(cmd.ID)
	if err == nil {
		return existingUser, domain.NewDomainError("USER_ALREADY_EXISTS", "User already exists", domain.ErrAlreadyExists)
	}

	// Check if email is already in use
	_, err = h.userRepo.FindByEmail(cmd.Email)
	if err == nil {
		return nil, domain.NewDomainError("EMAIL_ALREADY_EXISTS", "Email already in use", domain.ErrAlreadyExists)
	}

	// Create email value object
	email, err := entities.NewEmail(cmd.Email)
	if err != nil {
		return nil, domain.NewDomainError("INVALID_EMAIL", "Invalid email format", err)
	}

	// Create new user
	user := entities.NewUser(cmd.ID, cmd.Name, email)

	// Save user
	if err := h.userRepo.Create(&user); err != nil {
		return nil, domain.NewDomainError("CREATE_USER_FAILED", "Failed to create user", err)
	}

	return &user, nil
}
