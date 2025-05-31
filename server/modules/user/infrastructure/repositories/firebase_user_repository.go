package repositories

import (
	"context"
	"errors"
	"time"

	"teammate/server/modules/user/domain/entities"
	"teammate/server/modules/user/domain/repositories"
	"teammate/server/seedwork/infrastructure/firebase"

	"firebase.google.com/go/v4/auth"
)

// FirebaseUserRepository handles user operations using Firebase Auth
type FirebaseUserRepository struct {
	client *firebase.Client
}

// Ensure FirebaseUserRepository implements UserRepository
var _ repositories.UserRepository = (*FirebaseUserRepository)(nil)

// NewFirebaseUserRepository creates a new Firebase-based user repository
func NewFirebaseUserRepository(client *firebase.Client) *FirebaseUserRepository {
	return &FirebaseUserRepository{
		client: client,
	}
}

// FindAll retrieves all users from Firebase Auth
// Note: Firebase Auth doesn't support pagination in the same way as SQL databases
func (r *FirebaseUserRepository) FindAll() ([]*entities.User, error) {
	ctx := context.Background()

	// List users with a reasonable limit
	iter := r.client.ListUsers(ctx, "")
	var users []*entities.User

	for {
		exportedUser, err := iter.Next()
		if err != nil {
			break // End of iteration
		}

		user, err := r.exportedUserToEntity(exportedUser)
		if err != nil {
			continue // Skip invalid users
		}

		users = append(users, &user)
	}

	return users, nil
}

// FindByID retrieves a user by ID (Firebase UID)
func (r *FirebaseUserRepository) FindByID(id string) (*entities.User, error) {
	ctx := context.Background()

	userRecord, err := r.client.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	user, err := r.firebaseUserToEntity(userRecord)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// FindByEmail retrieves a user by email
func (r *FirebaseUserRepository) FindByEmail(email string) (*entities.User, error) {
	ctx := context.Background()

	userRecord, err := r.client.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	user, err := r.firebaseUserToEntity(userRecord)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Create creates a new user in Firebase Auth
// Note: This method cannot set a password since it's not in the User model
// For password-based signup, use the firebase.CreateUser function directly
func (r *FirebaseUserRepository) Create(user *entities.User) error {
	ctx := context.Background()

	params := (&auth.UserToCreate{}).
		Email(user.Email.String()).
		DisplayName(user.Name)

	// If ID is already set, use it; otherwise let Firebase generate one
	if user.GetID() != "" {
		params = params.UID(user.GetID())
	}

	userRecord, err := r.client.CreateUser(ctx, user.Email.String(), "", user.Name)
	if err != nil {
		return err
	}

	// Update the user with the Firebase UID and timestamps
	user.SetID(userRecord.UID)
	user.CreatedAt = time.Unix(userRecord.UserMetadata.CreationTimestamp, 0)
	user.UpdatedAt = time.Unix(userRecord.UserMetadata.LastLogInTimestamp, 0)

	return nil
}

// Update updates an existing user in Firebase Auth
func (r *FirebaseUserRepository) Update(user *entities.User) error {
	if user.GetID() == "" {
		return errors.New("user ID is required for update")
	}

	ctx := context.Background()

	params := (&auth.UserToUpdate{}).
		Email(user.Email.String()).
		DisplayName(user.Name)

	_, err := r.client.UpdateUser(ctx, user.GetID(), params)
	if err != nil {
		return err
	}

	user.UpdatedAt = time.Now()
	return nil
}

// Delete disables a user in Firebase Auth (Firebase doesn't have soft delete)
func (r *FirebaseUserRepository) Delete(id string) error {
	ctx := context.Background()

	params := (&auth.UserToUpdate{}).Disabled(true)
	_, err := r.client.UpdateUser(ctx, id, params)
	return err
}

// HardDelete permanently deletes a user from Firebase Auth
func (r *FirebaseUserRepository) HardDelete(id string) error {
	ctx := context.Background()
	return r.client.DeleteUser(ctx, id)
}

// Count returns the total number of users in Firebase Auth
func (r *FirebaseUserRepository) Count() (int64, error) {
	ctx := context.Background()

	iter := r.client.ListUsers(ctx, "")
	var count int64

	for {
		_, err := iter.Next()
		if err != nil {
			break // End of iteration
		}
		count++
	}

	return count, nil
}

// firebaseUserToEntity converts a Firebase UserRecord to our User entity
func (r *FirebaseUserRepository) firebaseUserToEntity(userRecord *auth.UserRecord) (entities.User, error) {
	email, err := entities.NewEmail(userRecord.Email)
	if err != nil {
		return entities.User{}, err
	}

	user := entities.NewUser(userRecord.UID, userRecord.DisplayName, email)
	user.CreatedAt = time.Unix(userRecord.UserMetadata.CreationTimestamp, 0)
	user.UpdatedAt = time.Unix(userRecord.UserMetadata.LastLogInTimestamp, 0)

	return user, nil
}

// exportedUserToEntity converts a Firebase ExportedUserRecord to our User entity
func (r *FirebaseUserRepository) exportedUserToEntity(exportedUser *auth.ExportedUserRecord) (entities.User, error) {
	email, err := entities.NewEmail(exportedUser.Email)
	if err != nil {
		return entities.User{}, err
	}

	user := entities.NewUser(exportedUser.UID, exportedUser.DisplayName, email)
	user.CreatedAt = time.Unix(exportedUser.UserMetadata.CreationTimestamp, 0)
	user.UpdatedAt = time.Unix(exportedUser.UserMetadata.LastLogInTimestamp, 0)

	return user, nil
}
