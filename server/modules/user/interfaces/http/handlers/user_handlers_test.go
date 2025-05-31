package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"teammate/server/modules/user/application/services"
	"teammate/server/modules/user/domain/entities"
	"teammate/server/modules/user/domain/repositories"
	infraRepos "teammate/server/modules/user/infrastructure/repositories"
	infraServices "teammate/server/modules/user/infrastructure/services"
	"teammate/server/modules/user/interfaces/http/dtos"
	"teammate/server/seedwork/infrastructure/config"
	"teammate/server/seedwork/infrastructure/database"
	"teammate/server/seedwork/infrastructure/firebase"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

// UserHandlersTestSuite defines the test suite for user handlers
type UserHandlersTestSuite struct {
	suite.Suite
	router              *gin.Engine
	handlers            *UserHandlers
	userService         *services.UserService
	userRepo            repositories.UserRepository
	firebaseAuthService *infraServices.FirebaseAuthService
}

// SetupSuite runs once before all tests
func (suite *UserHandlersTestSuite) SetupSuite() {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Initialize test database
	err := database.Initialize()
	suite.Require().NoError(err, "Failed to initialize test database")

	// Run migrations
	err = database.RunMigrations("../../../../../seedwork/infrastructure/database/migrations")
	suite.Require().NoError(err, "Failed to run migrations")

	// Load test configuration
	cfg := &config.Config{
		User: config.UserConfig{
			RepositoryType: "gorm", // Use GORM for tests
		},
	}

	// Create Firebase client
	firebaseClient, err := firebase.NewClient(cfg)
	suite.Require().NoError(err, "Failed to create Firebase client")

	// Create Firebase auth service
	suite.firebaseAuthService = infraServices.NewFirebaseAuthService(firebaseClient)

	// Create user repository (GORM for tests)
	suite.userRepo = infraRepos.NewGormUserRepository()

	// Create user service
	suite.userService = services.NewUserService(suite.userRepo)

	// Create handlers
	suite.handlers = NewUserHandlers(suite.userService, suite.firebaseAuthService)

	// Setup router
	suite.setupRouter()
}

// TearDownSuite runs once after all tests
func (suite *UserHandlersTestSuite) TearDownSuite() {
	// Clean up test data
	suite.cleanupTestData()
}

// SetupTest runs before each test
func (suite *UserHandlersTestSuite) SetupTest() {
	// Clean up any existing test data
	suite.cleanupTestData()
}

// TearDownTest runs after each test
func (suite *UserHandlersTestSuite) TearDownTest() {
	// Clean up test data after each test
	suite.cleanupTestData()
}

// setupRouter configures the test router with all user routes
func (suite *UserHandlersTestSuite) setupRouter() {
	suite.router = gin.New()

	// Add routes
	suite.router.GET("/users", suite.handlers.GetUsers)
	suite.router.GET("/users/:id", suite.handlers.GetUserByID)
	suite.router.GET("/me", suite.handlers.GetCurrentUser)
	suite.router.POST("/register", suite.handlers.CreateUser)
	suite.router.POST("/signup", suite.handlers.SignUpUser)
	suite.router.PUT("/users/:id", suite.handlers.UpdateUser)
	suite.router.DELETE("/users/:id", suite.handlers.DeleteUser)
}

// cleanupTestData removes all test users from the database
func (suite *UserHandlersTestSuite) cleanupTestData() {
	// Delete all users with test email domains
	db := database.DB
	db.Exec("DELETE FROM users WHERE email LIKE '%@test.com' OR email LIKE '%@example.com'")
}

// createTestUser creates a test user for testing purposes
func (suite *UserHandlersTestSuite) createTestUser(id, name, email string) *entities.User {
	emailVO, err := entities.NewEmail(email)
	suite.Require().NoError(err)

	user := entities.NewUser(id, name, emailVO)

	err = suite.userRepo.Create(&user)
	suite.Require().NoError(err)

	// Return the created user by fetching it from repository
	savedUser, err := suite.userRepo.FindByID(id)
	suite.Require().NoError(err)

	return savedUser
}

// TestGetUsers tests the GET /users endpoint
func (suite *UserHandlersTestSuite) TestGetUsers() {
	// Create test users
	user1 := suite.createTestUser("test-id-1", "Test User 1", "user1@test.com")
	user2 := suite.createTestUser("test-id-2", "Test User 2", "user2@test.com")

	// Make request
	req, _ := http.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Assert response
	suite.Equal(http.StatusOK, w.Code)

	var response dtos.UsersListResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.NoError(err)

	// Should have at least our 2 test users
	suite.GreaterOrEqual(len(response.Users), 2)
	suite.GreaterOrEqual(response.Total, int64(2))

	// Check if our test users are in the response
	userEmails := make(map[string]bool)
	for _, user := range response.Users {
		userEmails[user.Email] = true
	}
	suite.True(userEmails[user1.GetEmail().String()])
	suite.True(userEmails[user2.GetEmail().String()])
}

// TestGetUserByID tests the GET /users/:id endpoint
func (suite *UserHandlersTestSuite) TestGetUserByID() {
	// Create test user
	user := suite.createTestUser("test-get-by-id", "Test User", "getbyid@test.com")

	// Test successful case
	req, _ := http.NewRequest("GET", "/users/"+user.GetID(), nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)

	var response dtos.UserResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.NoError(err)

	suite.Equal(user.GetID(), response.ID)
	suite.Equal(user.GetName(), response.Name)
	suite.Equal(user.GetEmail().String(), response.Email)

	// Test not found case
	req, _ = http.NewRequest("GET", "/users/non-existent-id", nil)
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusNotFound, w.Code)
}

// TestGetCurrentUser tests the GET /me endpoint
func (suite *UserHandlersTestSuite) TestGetCurrentUser() {
	// Create test user
	user := suite.createTestUser("test-current-user", "Current User", "current@test.com")

	// Test with user in context (simulating auth middleware)
	req, _ := http.NewRequest("GET", "/me", nil)
	w := httptest.NewRecorder()

	// Create context with user
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("user", user)

	suite.handlers.GetCurrentUser(c)

	suite.Equal(http.StatusOK, w.Code)

	var response dtos.UserResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.NoError(err)

	suite.Equal(user.GetID(), response.ID)
	suite.Equal(user.GetName(), response.Name)
	suite.Equal(user.GetEmail().String(), response.Email)

	// Test without user in context
	req, _ = http.NewRequest("GET", "/me", nil)
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusUnauthorized, w.Code)
}

// TestCreateUser tests the POST /register endpoint
func (suite *UserHandlersTestSuite) TestCreateUser() {
	// Test successful creation
	createReq := dtos.CreateUserRequest{
		ID:    "test-create-user",
		Name:  "New User",
		Email: "newuser@test.com",
	}

	reqBody, _ := json.Marshal(createReq)
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusCreated, w.Code)

	var response dtos.UserResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.NoError(err)

	suite.Equal(createReq.ID, response.ID)
	suite.Equal(createReq.Name, response.Name)
	suite.Equal(createReq.Email, response.Email)
}

// TestMain sets up the test environment
func TestMain(m *testing.M) {
	// Set test environment variables
	os.Setenv("GIN_MODE", "test")
	os.Setenv("USER_REPOSITORY_TYPE", "gorm")

	// Set test database URL if not already set
	if os.Getenv("DATABASE_URL") == "" {
		os.Setenv("DATABASE_URL", "postgres://localhost/teammate_test?sslmode=disable")
	}

	// Run tests
	code := m.Run()

	// Exit with the test result code
	os.Exit(code)
}

// TestUserHandlersTestSuite runs the test suite
func TestUserHandlersTestSuite(t *testing.T) {
	suite.Run(t, new(UserHandlersTestSuite))
}
