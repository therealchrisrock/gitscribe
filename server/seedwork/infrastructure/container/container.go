package container

import (
	"teammate/server/modules/user/application/services"
	"teammate/server/modules/user/domain/repositories"
	userInfraRepos "teammate/server/modules/user/infrastructure/repositories"
	userInfraServices "teammate/server/modules/user/infrastructure/services"
	userMiddleware "teammate/server/modules/user/interfaces/http/middleware"
	"teammate/server/seedwork/infrastructure/config"
	"teammate/server/seedwork/infrastructure/firebase"
)

// Container holds all application dependencies
type Container struct {
	Config *config.Config

	// Infrastructure
	FirebaseClient *firebase.Client

	// Repositories
	UserRepository repositories.UserRepository

	// Services
	UserService         *services.UserService
	FirebaseAuthService *userInfraServices.FirebaseAuthService

	// Middleware
	AuthMiddleware *userMiddleware.AuthMiddleware
}

// NewContainer creates and wires up all dependencies
func NewContainer() (*Container, error) {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	// Create Firebase client
	firebaseClient, err := firebase.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	// Create Firebase auth service
	firebaseAuthService := userInfraServices.NewFirebaseAuthService(firebaseClient)

	// Create user repository based on configuration
	var userRepo repositories.UserRepository
	switch cfg.User.RepositoryType {
	case "gorm":
		userRepo = userInfraRepos.NewGormUserRepository()
	case "firebase":
		userRepo = userInfraRepos.NewFirebaseUserRepository(firebaseClient)
	default:
		userRepo = userInfraRepos.NewGormUserRepository() // Default to GORM
	}

	// Create user service
	userService := services.NewUserService(userRepo)

	// Create auth middleware
	authMiddleware := userMiddleware.NewAuthMiddleware(userRepo, firebaseAuthService)

	return &Container{
		Config:              cfg,
		FirebaseClient:      firebaseClient,
		UserRepository:      userRepo,
		UserService:         userService,
		FirebaseAuthService: firebaseAuthService,
		AuthMiddleware:      authMiddleware,
	}, nil
}

// GetUserService returns the user service
func (c *Container) GetUserService() *services.UserService {
	return c.UserService
}

// GetAuthMiddleware returns the auth middleware
func (c *Container) GetAuthMiddleware() *userMiddleware.AuthMiddleware {
	return c.AuthMiddleware
}

// GetFirebaseAuthService returns the Firebase auth service
func (c *Container) GetFirebaseAuthService() *userInfraServices.FirebaseAuthService {
	return c.FirebaseAuthService
}

// GetFirebaseClient returns the Firebase client
func (c *Container) GetFirebaseClient() *firebase.Client {
	return c.FirebaseClient
}

// GetConfig returns the configuration
func (c *Container) GetConfig() *config.Config {
	return c.Config
}
