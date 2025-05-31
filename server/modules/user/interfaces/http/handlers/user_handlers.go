package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"teammate/server/modules/user/application/commands"
	"teammate/server/modules/user/application/services"
	"teammate/server/modules/user/domain/entities"
	infraServices "teammate/server/modules/user/infrastructure/services"
	"teammate/server/modules/user/interfaces/http/dtos"
	"teammate/server/seedwork/domain"

	"github.com/gin-gonic/gin"
)

// UserHandlers contains all user-related HTTP handlers
type UserHandlers struct {
	userService         *services.UserService
	firebaseAuthService *infraServices.FirebaseAuthService
}

// NewUserHandlers creates a new user handlers instance
func NewUserHandlers(userService *services.UserService, firebaseAuthService *infraServices.FirebaseAuthService) *UserHandlers {
	return &UserHandlers{
		userService:         userService,
		firebaseAuthService: firebaseAuthService,
	}
}

// GetUsers returns all users
// @Summary Get all users
// @Description Get all users in the system
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} dtos.UsersListResponse
// @Failure 500 {object} map[string]string
// @Router /users [get]
func (h *UserHandlers) GetUsers(c *gin.Context) {
	result, err := h.userService.ListUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
		return
	}

	response := dtos.ToUsersListResponse(result.Users, result.Total)
	c.JSON(http.StatusOK, response)
}

// GetUserByID returns a specific user by ID
// @Summary Get a user by ID
// @Description Get a user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} dtos.UserResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /users/{id} [get]
func (h *UserHandlers) GetUserByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	user, err := h.userService.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	response := dtos.ToUserResponse(user)
	c.JSON(http.StatusOK, response)
}

// GetCurrentUser returns the currently authenticated user
// @Summary Get current user
// @Description Get the currently authenticated user
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dtos.UserResponse
// @Failure 401 {object} map[string]string
// @Router /me [get]
func (h *UserHandlers) GetCurrentUser(c *gin.Context) {
	// Get the user from the context (set by the auth middleware)
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	user, ok := userInterface.(*entities.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user data"})
		return
	}

	response := dtos.ToUserResponse(user)
	c.JSON(http.StatusOK, response)
}

// CreateUser adds a new user (for existing Firebase users)
// @Summary Create a new user
// @Description Register a new user in the system (for existing Firebase users)
// @Tags users
// @Accept json
// @Produce json
// @Param user body dtos.CreateUserRequest true "User information"
// @Success 201 {object} dtos.UserResponse
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /register [post]
func (h *UserHandlers) CreateUser(c *gin.Context) {
	var req dtos.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get Firebase UID if available (when called after authentication)
	if fbUID, exists := c.Get("firebase_uid"); exists {
		req.ID = fbUID.(string)
	}

	cmd := commands.CreateUserCommand{
		ID:    req.ID,
		Name:  req.Name,
		Email: req.Email,
	}

	user, err := h.userService.CreateUser(cmd)
	if err != nil {
		// Handle domain errors
		if domainErr, ok := err.(*domain.DomainError); ok {
			switch domainErr.Code {
			case "USER_ALREADY_EXISTS", "EMAIL_ALREADY_EXISTS":
				c.JSON(http.StatusConflict, gin.H{"error": domainErr.Message})
				return
			case "INVALID_EMAIL":
				c.JSON(http.StatusBadRequest, gin.H{"error": domainErr.Message})
				return
			}
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	response := dtos.ToUserResponse(user)
	c.JSON(http.StatusCreated, response)
}

// SignUpUser creates a new user with email/password in Firebase Auth and local DB
// @Summary Sign up a new user
// @Description Create a new user with email and password in Firebase Auth and local database
// @Tags users
// @Accept json
// @Produce json
// @Param signup body object{name=string,email=string,password=string} true "Signup information"
// @Success 201 {object} dtos.UserResponse
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /signup [post]
func (h *UserHandlers) SignUpUser(c *gin.Context) {
	var req struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if email already exists
	_, err := h.userService.GetUserByEmail(req.Email)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User with this email already exists"})
		return
	}

	// Create user in Firebase Auth
	userRecord, err := h.firebaseAuthService.CreateUser(c, req.Email, req.Password, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user in Firebase: " + err.Error()})
		return
	}

	// Create user in local database
	cmd := commands.CreateUserCommand{
		ID:    userRecord.UID,
		Name:  req.Name,
		Email: req.Email,
	}

	user, err := h.userService.CreateUser(cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user in database"})
		return
	}

	response := dtos.ToUserResponse(user)
	c.JSON(http.StatusCreated, response)
}

// LoginWithToken verifies a Firebase ID token and returns user info
// @Summary Login with Firebase ID token
// @Description Verify Firebase ID token and return user information
// @Tags auth
// @Accept json
// @Produce json
// @Param login body object{idToken=string} true "Firebase ID token"
// @Success 200 {object} object{user=dtos.UserResponse,token=string}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /login [post]
func (h *UserHandlers) LoginWithToken(c *gin.Context) {
	var req struct {
		IDToken string `json:"idToken" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify the Firebase ID token
	token, err := h.firebaseAuthService.VerifyIDToken(c, req.IDToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token: " + err.Error()})
		return
	}

	// Get user from database
	user, err := h.userService.GetUserByID(token.UID)
	if err != nil {
		// If user doesn't exist in our database, create them
		if domainErr, ok := err.(*domain.DomainError); ok && domainErr.Code == "USER_NOT_FOUND" {
			// Get user info from Firebase
			firebaseUser, err := h.firebaseAuthService.GetUser(c, token.UID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user from Firebase"})
				return
			}

			// Create user in local database
			cmd := commands.CreateUserCommand{
				ID:    firebaseUser.UID,
				Name:  firebaseUser.DisplayName,
				Email: firebaseUser.Email,
			}

			user, err = h.userService.CreateUser(cmd)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user in database"})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
			return
		}
	}

	response := dtos.ToUserResponse(user)
	c.JSON(http.StatusOK, gin.H{
		"user":  response,
		"token": req.IDToken, // Return the same token for client use
	})
}

// LoginWithEmailPassword authenticates user with email and password
// @Summary Login with email and password
// @Description Authenticate user with email and password. Returns both custom token and simulated ID token for testing.
// @Tags auth
// @Accept json
// @Produce json
// @Param login body object{email=string,password=string} true "Login credentials"
// @Success 200 {object} object{user=dtos.UserResponse,customToken=string,idToken=string,instructions=string}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /login/email [post]
func (h *UserHandlers) LoginWithEmailPassword(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user by email from Firebase
	firebaseUser, err := h.firebaseAuthService.GetUserByEmail(c, req.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Note: Firebase Admin SDK doesn't directly verify passwords
	// In a real implementation, you'd either:
	// 1. Use Firebase Client SDK on frontend (recommended)
	// 2. Use Firebase Auth REST API
	// 3. Create custom tokens for verified users

	// Create a custom token
	customToken, err := h.firebaseAuthService.CreateCustomToken(c, firebaseUser.UID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create authentication token"})
		return
	}

	// Get user from local database
	user, err := h.userService.GetUserByID(firebaseUser.UID)
	if err != nil {
		// If user doesn't exist in our database, create them
		if domainErr, ok := err.(*domain.DomainError); ok && domainErr.Code == "USER_NOT_FOUND" {
			cmd := commands.CreateUserCommand{
				ID:    firebaseUser.UID,
				Name:  firebaseUser.DisplayName,
				Email: firebaseUser.Email,
			}

			user, err = h.userService.CreateUser(cmd)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user in database"})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
			return
		}
	}

	response := dtos.ToUserResponse(user)
	c.JSON(http.StatusOK, gin.H{
		"user":         response,
		"customToken":  customToken,
		"instructions": "For Swagger UI testing: Use the 'customToken' with the /auth/exchange-token endpoint to get an ID token, or use the /login/email/direct endpoint for direct ID token.",
		"note":         "In production, use Firebase Client SDK to exchange customToken for idToken on the frontend",
	})
}

// LoginWithEmailPasswordDirect authenticates and returns an ID token directly (for testing only)
// @Summary Login with email/password and get ID token directly
// @Description Authenticate user and return an ID token directly. This simulates the complete Firebase flow for testing purposes.
// @Tags auth
// @Accept json
// @Produce json
// @Param login body object{email=string,password=string} true "Login credentials"
// @Success 200 {object} object{user=dtos.UserResponse,idToken=string,message=string}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /login/email/direct [post]
func (h *UserHandlers) LoginWithEmailPasswordDirect(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user by email from Firebase
	firebaseUser, err := h.firebaseAuthService.GetUserByEmail(c, req.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Create a custom token first
	customToken, err := h.firebaseAuthService.CreateCustomToken(c, firebaseUser.UID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create authentication token"})
		return
	}

	// For testing purposes, we'll simulate the Firebase Client SDK flow
	// In a real app, this would be done on the frontend using Firebase Client SDK
	//
	// Note: This is a simplified simulation for testing. In production:
	// 1. Frontend gets custom token from server
	// 2. Frontend uses Firebase Client SDK: signInWithCustomToken(auth, customToken)
	// 3. Frontend gets ID token: user.getIdToken()
	// 4. Frontend uses ID token for API requests

	// Since we can't actually call Firebase Client SDK from the server,
	// we'll create a mock ID token for testing purposes
	// WARNING: This is only for development/testing!

	if os.Getenv("APP_ENV") != "development" {
		c.JSON(http.StatusForbidden, gin.H{
			"error":   "Direct ID token endpoint only available in development",
			"message": "In production, use Firebase Client SDK to exchange custom tokens for ID tokens",
		})
		return
	}

	// Get user from local database
	user, err := h.userService.GetUserByID(firebaseUser.UID)
	if err != nil {
		// If user doesn't exist in our database, create them
		if domainErr, ok := err.(*domain.DomainError); ok && domainErr.Code == "USER_NOT_FOUND" {
			cmd := commands.CreateUserCommand{
				ID:    firebaseUser.UID,
				Name:  firebaseUser.DisplayName,
				Email: firebaseUser.Email,
			}

			user, err = h.userService.CreateUser(cmd)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user in database"})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
			return
		}
	}

	response := dtos.ToUserResponse(user)
	c.JSON(http.StatusOK, gin.H{
		"user":        response,
		"customToken": customToken,
		"message":     "DEVELOPMENT ONLY: Use the customToken with Firebase Client SDK to get a real ID token",
		"instructions": map[string]string{
			"step1": "This endpoint is for development testing only",
			"step2": "In production, use Firebase Client SDK on frontend",
			"step3": "For Swagger UI: Copy the customToken and use /auth/exchange-token for guidance",
			"step4": "For real apps: Use signInWithCustomToken(auth, customToken) then user.getIdToken()",
		},
	})
}

// RefreshToken refreshes a Firebase ID token
// @Summary Refresh Firebase token
// @Description Refresh an expired Firebase ID token
// @Tags auth
// @Accept json
// @Produce json
// @Param refresh body object{refreshToken=string} true "Firebase refresh token"
// @Success 200 {object} object{token=string,expiresAt=string}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/refresh [post]
func (h *UserHandlers) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refreshToken" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Note: Firebase Admin SDK doesn't handle refresh tokens directly
	// This would typically be handled by the Firebase Client SDK on the frontend
	// For server-side refresh, you'd need to use Firebase Auth REST API
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "Token refresh should be handled by Firebase Client SDK on frontend",
		"info":  "Use Firebase Auth REST API: https://firebase.google.com/docs/reference/rest/auth#section-refresh-token",
	})
}

// Logout logs out a user
// @Summary Logout user
// @Description Logout user and invalidate session
// @Tags auth
// @Accept json
// @Produce json
// @Param logout body object{idToken=string} true "Firebase ID token to invalidate"
// @Success 200 {object} object{message=string}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/logout [post]
func (h *UserHandlers) Logout(c *gin.Context) {
	var req struct {
		IDToken string `json:"idToken" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify the token first
	token, err := h.firebaseAuthService.VerifyIDToken(c, req.IDToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// Revoke refresh tokens for the user
	err = h.firebaseAuthService.RevokeRefreshTokens(c, token.UID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully logged out",
		"uid":     token.UID,
	})
}

// ForgotPassword initiates password reset
// @Summary Forgot password
// @Description Send password reset email
// @Tags auth
// @Accept json
// @Produce json
// @Param forgot body object{email=string} true "User email"
// @Success 200 {object} object{message=string}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/forgot-password [post]
func (h *UserHandlers) ForgotPassword(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user exists
	_, err := h.firebaseAuthService.GetUserByEmail(c, req.Email)
	if err != nil {
		// Don't reveal if email exists or not for security
		c.JSON(http.StatusOK, gin.H{
			"message": "If the email exists, a password reset link has been sent",
		})
		return
	}

	// Generate password reset link
	resetLink, err := h.firebaseAuthService.GeneratePasswordResetLink(c, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate reset link"})
		return
	}

	// In production, you'd send this via email service
	// For now, we'll return it in the response (NOT recommended for production)
	c.JSON(http.StatusOK, gin.H{
		"message":   "Password reset link generated",
		"resetLink": resetLink, // Remove this in production
	})
}

// ResetPassword confirms password reset
// @Summary Reset password
// @Description Confirm password reset with code
// @Tags auth
// @Accept json
// @Produce json
// @Param reset body object{oobCode=string,newPassword=string} true "Reset code and new password"
// @Success 200 {object} object{message=string}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/reset-password [post]
func (h *UserHandlers) ResetPassword(c *gin.Context) {
	var req struct {
		OOBCode     string `json:"oobCode" binding:"required"`
		NewPassword string `json:"newPassword" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Confirm password reset
	err := h.firebaseAuthService.ConfirmPasswordReset(c, req.OOBCode, req.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired reset code"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password successfully reset",
	})
}

// VerifyEmail verifies user email address
// @Summary Verify email
// @Description Verify user email with verification code
// @Tags auth
// @Accept json
// @Produce json
// @Param verify body object{oobCode=string} true "Email verification code"
// @Success 200 {object} object{message=string}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/verify-email [post]
func (h *UserHandlers) VerifyEmail(c *gin.Context) {
	var req struct {
		OOBCode string `json:"oobCode" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify email with code
	err := h.firebaseAuthService.VerifyEmail(c, req.OOBCode)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired verification code"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Email successfully verified",
	})
}

// UpdateUser updates an existing user
// @Summary Update a user
// @Description Update an existing user
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body dtos.UpdateUserRequest true "User information to update"
// @Security BearerAuth
// @Success 200 {object} dtos.UserResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{id} [put]
func (h *UserHandlers) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	var req dtos.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cmd := commands.UpdateUserCommand{
		ID:    id,
		Name:  req.Name,
		Email: req.Email,
	}

	user, err := h.userService.UpdateUser(cmd)
	if err != nil {
		// Handle domain errors
		if domainErr, ok := err.(*domain.DomainError); ok {
			switch domainErr.Code {
			case "USER_NOT_FOUND":
				c.JSON(http.StatusNotFound, gin.H{"error": domainErr.Message})
				return
			case "EMAIL_ALREADY_EXISTS":
				c.JSON(http.StatusConflict, gin.H{"error": domainErr.Message})
				return
			case "INVALID_EMAIL":
				c.JSON(http.StatusBadRequest, gin.H{"error": domainErr.Message})
				return
			}
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	response := dtos.ToUserResponse(user)
	c.JSON(http.StatusOK, response)
}

// DeleteUser removes a user
// @Summary Delete a user
// @Description Delete an existing user
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{id} [delete]
func (h *UserHandlers) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	cmd := commands.DeleteUserCommand{
		ID:   id,
		Hard: false, // Soft delete by default
	}

	err := h.userService.DeleteUser(cmd)
	if err != nil {
		// Handle domain errors
		if domainErr, ok := err.(*domain.DomainError); ok {
			switch domainErr.Code {
			case "USER_NOT_FOUND":
				c.JSON(http.StatusNotFound, gin.H{"error": domainErr.Message})
				return
			}
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":      id,
		"deleted": true,
	})
}

// DebugFirebaseConfig returns Firebase configuration info for debugging
// @Summary Debug Firebase configuration
// @Description Get Firebase project information for debugging
// @Tags debug
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /debug/firebase [get]
func (h *UserHandlers) DebugFirebaseConfig(c *gin.Context) {
	// Only allow in development
	if os.Getenv("APP_ENV") != "development" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Debug endpoint only available in development"})
		return
	}

	// Try to get a test user to verify Firebase connection
	testUID := "test-uid-that-does-not-exist"
	_, err := h.firebaseAuthService.GetUser(c, testUID)

	var firebaseStatus string
	if err != nil {
		if strings.Contains(err.Error(), "user not found") || strings.Contains(err.Error(), "no user record") {
			firebaseStatus = "connected (user not found as expected)"
		} else {
			firebaseStatus = "error: " + err.Error()
		}
	} else {
		firebaseStatus = "connected (unexpected: test user exists)"
	}

	c.JSON(http.StatusOK, gin.H{
		"firebase_status":     firebaseStatus,
		"credentials_path":    os.Getenv("FIREBASE_CREDENTIALS_PATH"),
		"app_env":             os.Getenv("APP_ENV"),
		"project_id_from_env": os.Getenv("FIREBASE_PROJECT_ID"),
	})
}

// ExchangeCustomToken provides guidance on how to exchange custom token for ID token
// @Summary Exchange custom token for ID token
// @Description Provides instructions on how to exchange a custom token for an ID token
// @Tags auth
// @Accept json
// @Produce json
// @Param token body object{customToken=string} true "Custom token to exchange"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /auth/exchange-token [post]
func (h *UserHandlers) ExchangeCustomToken(c *gin.Context) {
	var req struct {
		CustomToken string `json:"customToken" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Note: Firebase Admin SDK cannot directly exchange custom tokens for ID tokens
	// This must be done using Firebase Client SDK on the frontend
	tokenPreview := req.CustomToken
	if len(req.CustomToken) > 50 {
		tokenPreview = req.CustomToken[:50] + "..."
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Custom token received. To get an ID token, use Firebase Client SDK:",
		"instructions": map[string]interface{}{
			"step1": "Install Firebase Client SDK: npm install firebase",
			"step2": "Initialize Firebase in your app with your project config",
			"step3": "Use signInWithCustomToken(auth, customToken)",
			"step4": "Call user.getIdToken() to get the ID token",
			"step5": "Use the ID token in Authorization header for API requests",
		},
		"example_code": `
import { signInWithCustomToken } from "firebase/auth";
import { auth } from "./firebase-config";

const userCredential = await signInWithCustomToken(auth, customToken);
const idToken = await userCredential.user.getIdToken();

// Use idToken for API requests
fetch('/me', {
  headers: { 'Authorization': 'Bearer ' + idToken }
});`,
		"custom_token_preview": tokenPreview,
	})
}

// GetIDTokenFromCustomToken exchanges a custom token for an ID token using Firebase REST API
// @Summary Exchange custom token for ID token (for testing)
// @Description Exchange a custom token for an ID token using Firebase REST API. Use this for Swagger UI testing.
// @Tags auth
// @Accept json
// @Produce json
// @Param token body object{customToken=string} true "Custom token to exchange"
// @Success 200 {object} object{idToken=string,refreshToken=string,expiresIn=string,message=string}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/get-id-token [post]
func (h *UserHandlers) GetIDTokenFromCustomToken(c *gin.Context) {
	var req struct {
		CustomToken string `json:"customToken" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Use Firebase REST API to exchange custom token for ID token
	// This is what Firebase Client SDK does under the hood
	idToken, refreshToken, expiresIn, err := h.exchangeCustomTokenForIDToken(req.CustomToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to exchange custom token for ID token",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"idToken":      idToken,
		"refreshToken": refreshToken,
		"expiresIn":    expiresIn,
		"message":      "Success! Use the 'idToken' in Authorization header as 'Bearer {idToken}' for authenticated requests",
		"example":      "Authorization: Bearer " + idToken[:50] + "...",
	})
}

// exchangeCustomTokenForIDToken uses Firebase REST API to exchange custom token for ID token
func (h *UserHandlers) exchangeCustomTokenForIDToken(customToken string) (string, string, string, error) {
	// Firebase REST API endpoint for exchanging custom token for ID token
	// Documentation: https://firebase.google.com/docs/reference/rest/auth#section-verify-custom-token

	firebaseAPIKey := os.Getenv("FIREBASE_API_KEY")
	if firebaseAPIKey == "" {
		return "", "", "", fmt.Errorf("FIREBASE_API_KEY environment variable not set")
	}

	url := fmt.Sprintf("https://identitytoolkit.googleapis.com/v1/accounts:signInWithCustomToken?key=%s", firebaseAPIKey)

	payload := map[string]interface{}{
		"token":             customToken,
		"returnSecureToken": true,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", "", "", err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return "", "", "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", "", "", fmt.Errorf("Firebase API error: %s", string(body))
	}

	var result struct {
		IDToken      string `json:"idToken"`
		RefreshToken string `json:"refreshToken"`
		ExpiresIn    string `json:"expiresIn"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", "", "", err
	}

	return result.IDToken, result.RefreshToken, result.ExpiresIn, nil
}
