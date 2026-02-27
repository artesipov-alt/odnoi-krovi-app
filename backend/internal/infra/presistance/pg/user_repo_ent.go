package pg

import (
	"context"
	"errors"
	"fmt"

	"github.com/artesipov-alt/odnoi-krovi-app/ent"
	"github.com/artesipov-alt/odnoi-krovi-app/ent/schema"
	"github.com/artesipov-alt/odnoi-krovi-app/ent/user"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/services"
)

// EntUserRepository implements UserRepository using ENT
type EntUserRepository struct {
	client *ent.Client
}

// NewEntUserRepository creates a new ENT user repository
func NewEntUserRepository(client *ent.Client) *EntUserRepository {
	return &EntUserRepository{
		client: client,
	}
}

// Create creates a new user in the database
func (r *EntUserRepository) Create(ctx context.Context, input *ent.CreateUserInput) (*ent.User, error) {
	newUser, err := r.client.User.
		Create().
		SetInput(*input). // МАГИЯ! Ent сам вызовет все нужные Set-методы
		Save(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return newUser, nil
}

// GetQuery returns a query for eager loading
func (r *EntUserRepository) GetByID(ctx context.Context, id string, opts services.UserPreloadOptions) (*ent.User, error) {
	quser := r.client.User.Query().Where(user.ID(id))

	if opts.WithPets {
		quser = quser.WithPets()
	}

	user, err := quser.Only(ctx)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetByTelegram returns a user by Telegram ID
func (r *EntUserRepository) GetByTelegram(ctx context.Context, telegramID int64, opts services.UserPreloadOptions) (*ent.User, error) {
	quser := r.client.User.Query().Where(user.TelegramID(telegramID))

	if opts.WithPets {
		quser = quser.WithPets()
	}

	user, err := quser.Only(ctx)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// ExistsByID checks if a user with the given ID exists
func (r *EntUserRepository) ExistsByID(ctx context.Context, id string) (bool, error) {
	if id == "" {
		return false, errors.New("invalid user ID")
	}

	exists, err := r.client.User.Query().
		Where(user.ID(id)).
		Exist(ctx)

	if err != nil {
		return false, fmt.Errorf("failed to check user existence by id %s: %w", id, err)
	}

	return exists, nil
}

// Update updates an existing user in the database
func (r *EntUserRepository) Update(ctx context.Context, id string, input *ent.UpdateUserInput) error {
	if input == nil {
		return errors.New("user cannot be nil")
	}

	if id == "" {
		return errors.New("invalid user ID")
	}

	_, err := r.client.User.UpdateOneID(id).
		SetInput(*input).
		Save(ctx)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// Delete deletes a user by their ID (soft delete via SoftDeleteMixin)
func (r *EntUserRepository) Delete(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("invalid user ID")
	}

	// Soft delete via SoftDeleteMixin hook
	err := r.client.User.DeleteOneID(id).Exec(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return fmt.Errorf("user with id %s not found", id)
		}
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// ExistsByTelegramID checks if a user with the given Telegram ID exists
func (r *EntUserRepository) ExistsByTelegramID(ctx context.Context, telegramID int64) (bool, error) {
	if telegramID <= 0 {
		return false, errors.New("invalid telegram ID")
	}

	exists, err := r.client.User.Query().
		Where(user.TelegramID(telegramID)).
		Exist(ctx)

	if err != nil {
		return false, fmt.Errorf("failed to check user existence by telegram id %d: %w", telegramID, err)
	}

	return exists, nil
}

// ResetUser resets user's email and phone number by ID
func (r *EntUserRepository) ResetUser(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("invalid user ID")
	}

	err := r.client.User.UpdateOneID(id).
		SetEmail("").
		SetPhone("").
		Exec(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return fmt.Errorf("user with id %s not found", id)
		}
		return fmt.Errorf("failed to reset user data: %w", err)
	}

	return nil
}

// RestoreUser restores a soft-deleted user by setting deleted_at to NULL
func (r *EntUserRepository) RestoreUser(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("invalid user ID")
	}

	// Use SkipSoftDelete context to find the deleted record
	ctxWithSkip := schema.SkipSoftDelete(ctx)

	err := r.client.User.UpdateOneID(id).
		ClearDeletedAt().
		Exec(ctxWithSkip)

	if err != nil {
		if ent.IsNotFound(err) {
			return fmt.Errorf("user with id %s not found", id)
		}
		return fmt.Errorf("failed to restore user: %w", err)
	}

	return nil
}

// GetDeletedUsers retrieves all soft-deleted users
func (r *EntUserRepository) GetDeletedUsers(ctx context.Context) ([]*ent.User, error) {
	// Use SkipSoftDelete context to see deleted records
	ctxWithSkip := schema.SkipSoftDelete(ctx)

	users, err := r.client.User.Query().
		Where(user.DeletedAtNotNil()).
		All(ctxWithSkip)

	if err != nil {
		return nil, fmt.Errorf("failed to get deleted users: %w", err)
	}

	return users, nil
}

// AddPhotoURLs adds new photo paths to the user's PhotoUrls array
func (r *EntUserRepository) AddPhotoURLs(ctx context.Context, id string, paths []string) error {
	if id == "" {
		return errors.New("invalid user ID")
	}

	// Replace photo URLs with new paths
	newPhotoUrls := paths

	// Update user
	err := r.client.User.UpdateOneID(id).
		SetPhotoUrls(newPhotoUrls).
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to update user photo URLs: %w", err)
	}

	return nil
}
