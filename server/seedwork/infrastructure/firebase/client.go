package firebase

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"teammate/server/seedwork/infrastructure/config"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

// Client represents a Firebase client wrapper
type Client struct {
	Auth *auth.Client
	app  *firebase.App
}

// NewClient creates a new Firebase client based on configuration
func NewClient(cfg *config.Config) (*Client, error) {
	app, err := initializeFirebaseApp(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Firebase app: %w", err)
	}

	authClient, err := app.Auth(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Firebase Auth client: %w", err)
	}

	log.Println("Firebase Auth initialized successfully")

	return &Client{
		Auth: authClient,
		app:  app,
	}, nil
}

// initializeFirebaseApp initializes the Firebase app based on configuration
func initializeFirebaseApp(cfg *config.Config) (*firebase.App, error) {
	var app *firebase.App
	var err error

	if cfg.Firebase.CredentialsPath != "" {
		// Use credentials file if specified
		opt := option.WithCredentialsFile(cfg.Firebase.CredentialsPath)
		app, err = firebase.NewApp(context.Background(), nil, opt)
	} else if credJSON := os.Getenv("FIREBASE_CREDENTIALS_JSON"); credJSON != "" {
		// Use credentials JSON from env var if specified
		opt := option.WithCredentialsJSON([]byte(credJSON))
		app, err = firebase.NewApp(context.Background(), nil, opt)
	} else if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") != "" {
		// Use GOOGLE_APPLICATION_CREDENTIALS env var (standard method)
		app, err = firebase.NewApp(context.Background(), nil)
	} else if cfg.Server.Env == "development" {
		// In development, create an empty credentials file if it doesn't exist
		if err = createEmptyCredentialsIfNeeded(); err != nil {
			return nil, fmt.Errorf("failed to create empty credentials: %w", err)
		}

		// Use the development credentials
		opt := option.WithCredentialsFile("firebase-credentials-dev.json")
		app, err = firebase.NewApp(context.Background(), nil, opt)
	} else {
		return nil, fmt.Errorf("no Firebase credentials provided")
	}

	return app, err
}

// createEmptyCredentialsIfNeeded creates a development Firebase credentials file
func createEmptyCredentialsIfNeeded() error {
	filename := "firebase-credentials-dev.json"

	// Check if file already exists
	if _, err := os.Stat(filename); err == nil {
		return nil // File exists, nothing to do
	}

	// Create minimal valid credentials file for development
	// Note: This is just for development; it won't connect to a real Firebase project
	credentials := map[string]interface{}{
		"type":                        "service_account",
		"project_id":                  "development-project",
		"private_key_id":              "development-key-id",
		"private_key":                 "-----BEGIN PRIVATE KEY-----\nMIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDev...\n-----END PRIVATE KEY-----\n",
		"client_email":                "firebase-adminsdk-dev@development-project.iam.gserviceaccount.com",
		"client_id":                   "123456789",
		"auth_uri":                    "https://accounts.google.com/o/oauth2/auth",
		"token_uri":                   "https://oauth2.googleapis.com/token",
		"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
		"client_x509_cert_url":        "https://www.googleapis.com/robot/v1/metadata/x509/firebase-adminsdk-dev%40development-project.iam.gserviceaccount.com",
	}

	jsonBytes, err := json.MarshalIndent(credentials, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, jsonBytes, 0600)
}

// VerifyIDToken verifies the Firebase ID token
func (c *Client) VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error) {
	return c.Auth.VerifyIDToken(ctx, idToken)
}

// CreateUser creates a new user in Firebase Auth
func (c *Client) CreateUser(ctx context.Context, email, password, name string) (*auth.UserRecord, error) {
	params := (&auth.UserToCreate{}).
		Email(email).
		Password(password).
		DisplayName(name)
	return c.Auth.CreateUser(ctx, params)
}

// GetUser retrieves a user by ID
func (c *Client) GetUser(ctx context.Context, uid string) (*auth.UserRecord, error) {
	return c.Auth.GetUser(ctx, uid)
}

// GetUserByEmail retrieves a user by email
func (c *Client) GetUserByEmail(ctx context.Context, email string) (*auth.UserRecord, error) {
	return c.Auth.GetUserByEmail(ctx, email)
}

// UpdateUser updates an existing user
func (c *Client) UpdateUser(ctx context.Context, uid string, params *auth.UserToUpdate) (*auth.UserRecord, error) {
	return c.Auth.UpdateUser(ctx, uid, params)
}

// DeleteUser deletes a user
func (c *Client) DeleteUser(ctx context.Context, uid string) error {
	return c.Auth.DeleteUser(ctx, uid)
}

// ListUsers returns an iterator for listing users
func (c *Client) ListUsers(ctx context.Context, nextPageToken string) *auth.UserIterator {
	return c.Auth.Users(ctx, nextPageToken)
}

// CreateCustomToken creates a custom token for a user
func (c *Client) CreateCustomToken(ctx context.Context, uid string) (string, error) {
	return c.Auth.CustomToken(ctx, uid)
}

// RevokeRefreshTokens revokes all refresh tokens for a user
func (c *Client) RevokeRefreshTokens(ctx context.Context, uid string) error {
	return c.Auth.RevokeRefreshTokens(ctx, uid)
}

// GeneratePasswordResetLink generates a password reset link
func (c *Client) GeneratePasswordResetLink(ctx context.Context, email string) (string, error) {
	link, err := c.Auth.PasswordResetLink(ctx, email)
	return link, err
}

// ConfirmPasswordReset confirms a password reset with code
func (c *Client) ConfirmPasswordReset(ctx context.Context, oobCode, newPassword string) error {
	// Note: Firebase Admin SDK doesn't directly support password reset confirmation
	// This would typically be handled by Firebase Auth REST API or Client SDK
	return fmt.Errorf("password reset confirmation should be handled by Firebase Client SDK or REST API")
}

// VerifyEmail verifies an email with verification code
func (c *Client) VerifyEmail(ctx context.Context, oobCode string) error {
	// Note: Firebase Admin SDK doesn't directly support email verification
	// This would typically be handled by Firebase Auth REST API or Client SDK
	return fmt.Errorf("email verification should be handled by Firebase Client SDK or REST API")
}
