package middleware

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"

	"teammate/server/modules/user/domain/entities"
	"teammate/server/modules/user/domain/repositories"
	"teammate/server/modules/user/infrastructure/services"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware provides authentication functionality
type AuthMiddleware struct {
	userRepo            repositories.UserRepository
	firebaseAuthService *services.FirebaseAuthService
}

// NewAuthMiddleware creates a new auth middleware instance
func NewAuthMiddleware(userRepo repositories.UserRepository, firebaseAuthService *services.FirebaseAuthService) *AuthMiddleware {
	return &AuthMiddleware{
		userRepo:            userRepo,
		firebaseAuthService: firebaseAuthService,
	}
}

// FirebaseAuth is a middleware that verifies Firebase ID tokens
func (m *AuthMiddleware) FirebaseAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		// Check if it has the Bearer prefix
		idToken := strings.TrimSpace(strings.Replace(authHeader, "Bearer", "", 1))
		if idToken == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			return
		}

		// Debug logging (remove in production)
		if os.Getenv("APP_ENV") == "development" {
			tokenPreview := idToken
			if len(idToken) > 20 {
				tokenPreview = idToken[:20] + "..."
			}
			log.Printf("DEBUG: Received token (first 20 chars): %s", tokenPreview)
		}

		// Try to verify as ID token first
		token, err := m.firebaseAuthService.VerifyIDToken(context.Background(), idToken)
		if err != nil {
			// If ID token verification fails, check if it's a custom token error
			if strings.Contains(err.Error(), "custom token") {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error":   "Custom token detected",
					"message": "Please exchange your custom token for an ID token using Firebase Client SDK before making API requests",
					"details": err.Error(),
				})
				return
			}

			// Add more detailed error logging for debugging
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Invalid token",
				"details": err.Error(), // This will help debug the specific issue
			})
			return
		}

		// Get the user from the repository using the Firebase UID as the primary key
		user, err := m.userRepo.FindByID(token.UID)
		if err != nil {
			// User doesn't exist in our database, we can create them automatically
			// based on the token claims or require explicit registration

			// Option 1: Auto-create user
			var email string
			if emailClaim, ok := token.Claims["email"]; ok {
				email = emailClaim.(string)
			}

			var name string
			if nameClaim, ok := token.Claims["name"]; ok {
				name = nameClaim.(string)
			} else if emailClaim, ok := token.Claims["email"]; ok {
				// Use email as fallback for name
				name = strings.Split(emailClaim.(string), "@")[0]
			}

			emailVO, err := entities.NewEmail(email)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid email"})
				return
			}

			newUser := entities.NewUser(token.UID, name, emailVO)

			if err := m.userRepo.Create(&newUser); err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError,
					gin.H{"error": "Failed to create user"})
				return
			}

			user = &newUser

			// Option 2: Require explicit registration (uncomment to use this instead)
			// c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "User not registered"})
			// return
		}

		// Store both user and raw token in the context for handlers
		c.Set("user", user)
		c.Set("firebase_uid", token.UID)
		c.Set("firebase_token", token)

		c.Next()
	}
}

// RequireRole checks if the authenticated user has a specific role
func (m *AuthMiddleware) RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user from context (set by FirebaseAuth middleware)
		_, exists := c.Get("user")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// TODO: Implement role checking when roles are added to the user model
		// This is just a placeholder for future role-based auth

		c.Next()
	}
}
