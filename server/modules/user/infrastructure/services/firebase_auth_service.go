package services

import (
	"context"

	"teammate/server/seedwork/infrastructure/firebase"

	"firebase.google.com/go/v4/auth"
)

// FirebaseAuthService provides Firebase authentication operations
type FirebaseAuthService struct {
	client *firebase.Client
}

// NewFirebaseAuthService creates a new Firebase auth service
func NewFirebaseAuthService(client *firebase.Client) *FirebaseAuthService {
	return &FirebaseAuthService{
		client: client,
	}
}

// VerifyIDToken verifies a Firebase ID token
func (s *FirebaseAuthService) VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error) {
	return s.client.VerifyIDToken(ctx, idToken)
}

// CreateUser creates a new user in Firebase Auth
func (s *FirebaseAuthService) CreateUser(ctx context.Context, email, password, name string) (*auth.UserRecord, error) {
	return s.client.CreateUser(ctx, email, password, name)
}

// GetUser retrieves a user by ID
func (s *FirebaseAuthService) GetUser(ctx context.Context, uid string) (*auth.UserRecord, error) {
	return s.client.GetUser(ctx, uid)
}

// GetUserByEmail retrieves a user by email
func (s *FirebaseAuthService) GetUserByEmail(ctx context.Context, email string) (*auth.UserRecord, error) {
	return s.client.GetUserByEmail(ctx, email)
}

// UpdateUser updates an existing user
func (s *FirebaseAuthService) UpdateUser(ctx context.Context, uid string, params *auth.UserToUpdate) (*auth.UserRecord, error) {
	return s.client.UpdateUser(ctx, uid, params)
}

// DeleteUser deletes a user
func (s *FirebaseAuthService) DeleteUser(ctx context.Context, uid string) error {
	return s.client.DeleteUser(ctx, uid)
}

// ListUsers returns an iterator for listing users
func (s *FirebaseAuthService) ListUsers(ctx context.Context, nextPageToken string) *auth.UserIterator {
	return s.client.ListUsers(ctx, nextPageToken)
}

// CreateCustomToken creates a custom token for a user
func (s *FirebaseAuthService) CreateCustomToken(ctx context.Context, uid string) (string, error) {
	return s.client.CreateCustomToken(ctx, uid)
}

// RevokeRefreshTokens revokes all refresh tokens for a user
func (s *FirebaseAuthService) RevokeRefreshTokens(ctx context.Context, uid string) error {
	return s.client.RevokeRefreshTokens(ctx, uid)
}

// GeneratePasswordResetLink generates a password reset link
func (s *FirebaseAuthService) GeneratePasswordResetLink(ctx context.Context, email string) (string, error) {
	return s.client.GeneratePasswordResetLink(ctx, email)
}

// ConfirmPasswordReset confirms a password reset with code
func (s *FirebaseAuthService) ConfirmPasswordReset(ctx context.Context, oobCode, newPassword string) error {
	return s.client.ConfirmPasswordReset(ctx, oobCode, newPassword)
}

// VerifyEmail verifies an email with verification code
func (s *FirebaseAuthService) VerifyEmail(ctx context.Context, oobCode string) error {
	return s.client.VerifyEmail(ctx, oobCode)
}
