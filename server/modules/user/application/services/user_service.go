package services

import (
	"teammate/server/modules/user/application/commands"
	"teammate/server/modules/user/application/queries"
	"teammate/server/modules/user/domain/entities"
	"teammate/server/modules/user/domain/repositories"
)

// UserService orchestrates user-related operations
type UserService struct {
	// Command handlers
	createUserHandler *commands.CreateUserHandler
	updateUserHandler *commands.UpdateUserHandler
	deleteUserHandler *commands.DeleteUserHandler

	// Query handlers
	getUserHandler   *queries.GetUserHandler
	listUsersHandler *queries.ListUsersHandler
}

// NewUserService creates a new user service with all dependencies
func NewUserService(userRepo repositories.UserRepository) *UserService {
	return &UserService{
		// Initialize command handlers
		createUserHandler: commands.NewCreateUserHandler(userRepo),
		updateUserHandler: commands.NewUpdateUserHandler(userRepo),
		deleteUserHandler: commands.NewDeleteUserHandler(userRepo),

		// Initialize query handlers
		getUserHandler:   queries.NewGetUserHandler(userRepo),
		listUsersHandler: queries.NewListUsersHandler(userRepo),
	}
}

// Command operations

// CreateUser creates a new user
func (s *UserService) CreateUser(cmd commands.CreateUserCommand) (*entities.User, error) {
	return s.createUserHandler.Handle(cmd)
}

// UpdateUser updates an existing user
func (s *UserService) UpdateUser(cmd commands.UpdateUserCommand) (*entities.User, error) {
	return s.updateUserHandler.Handle(cmd)
}

// DeleteUser deletes a user
func (s *UserService) DeleteUser(cmd commands.DeleteUserCommand) error {
	return s.deleteUserHandler.Handle(cmd)
}

// Query operations

// GetUserByID retrieves a user by ID
func (s *UserService) GetUserByID(id string) (*entities.User, error) {
	return s.getUserHandler.HandleGetUser(queries.GetUserQuery{ID: id})
}

// GetUserByEmail retrieves a user by email
func (s *UserService) GetUserByEmail(email string) (*entities.User, error) {
	return s.getUserHandler.HandleGetUserByEmail(queries.GetUserByEmailQuery{Email: email})
}

// ListUsers retrieves all users
func (s *UserService) ListUsers() (queries.ListUsersResult, error) {
	return s.listUsersHandler.Handle(queries.ListUsersQuery{})
}
