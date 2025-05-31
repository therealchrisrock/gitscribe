package routes

import (
	"teammate/server/modules/user/interfaces/http/handlers"
	"teammate/server/modules/user/interfaces/http/middleware"

	"github.com/gin-gonic/gin"
)

// UserRoutes sets up all user-related routes
type UserRoutes struct {
	userHandlers   *handlers.UserHandlers
	authMiddleware *middleware.AuthMiddleware
}

// NewUserRoutes creates a new user routes instance
func NewUserRoutes(userHandlers *handlers.UserHandlers, authMiddleware *middleware.AuthMiddleware) *UserRoutes {
	return &UserRoutes{
		userHandlers:   userHandlers,
		authMiddleware: authMiddleware,
	}
}

// SetupPublicRoutes sets up public user routes (no authentication required)
func (ur *UserRoutes) SetupPublicRoutes(public *gin.RouterGroup) {
	// Auth endpoints
	public.POST("/signup", ur.userHandlers.SignUpUser)                               // Create new account
	public.POST("/login", ur.userHandlers.LoginWithToken)                            // Primary login (Firebase ID token)
	public.POST("/login/email", ur.userHandlers.LoginWithEmailPassword)              // Email/password login
	public.POST("/login/email/direct", ur.userHandlers.LoginWithEmailPasswordDirect) // Email/password login (development)
	public.POST("/register", ur.userHandlers.CreateUser)                             // Register existing Firebase user
	public.POST("/refresh", ur.userHandlers.RefreshToken)                            // Refresh Firebase token
	public.POST("/logout", ur.userHandlers.Logout)                                   // Logout user
	public.POST("/forgot-password", ur.userHandlers.ForgotPassword)                  // Password reset
	public.POST("/reset-password", ur.userHandlers.ResetPassword)                    // Confirm password reset
	public.POST("/verify-email", ur.userHandlers.VerifyEmail)                        // Verify email address
	public.POST("/auth/exchange-token", ur.userHandlers.ExchangeCustomToken)         // Custom token exchange guidance
	public.POST("/auth/get-id-token", ur.userHandlers.GetIDTokenFromCustomToken)     // Exchange custom token for ID token

	// User endpoints (typically would be protected in production)
	public.GET("/users", ur.userHandlers.GetUsers)        // List all users
	public.GET("/users/:id", ur.userHandlers.GetUserByID) // Get user by ID

	// Debug endpoints (development only)
	public.GET("/debug/firebase", ur.userHandlers.DebugFirebaseConfig) // Firebase debug info
}

// SetupProtectedRoutes sets up protected user routes (authentication required)
func (ur *UserRoutes) SetupProtectedRoutes(protected *gin.RouterGroup) {
	// Apply auth middleware to protected routes
	protected.Use(ur.authMiddleware.FirebaseAuth())

	// Current user endpoints
	protected.GET("/me", ur.userHandlers.GetCurrentUser)

	// User management (typically admin only, but simplified for demo)
	protected.PUT("/users/:id", ur.userHandlers.UpdateUser)
	protected.DELETE("/users/:id", ur.userHandlers.DeleteUser)
}
