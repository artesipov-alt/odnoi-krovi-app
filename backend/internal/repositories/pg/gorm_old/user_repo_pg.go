package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/models"
	"gorm.io/gorm"
)

// PostgresUserRepository implements UserRepository for PostgreSQL
type PostgresUserRepository struct {
	db *gorm.DB
}

// NewPostgresUserRepository creates a new PostgreSQL user repository
func NewPostgresUserRepository(db *gorm.DB) *PostgresUserRepository {
	return &PostgresUserRepository{
		db: db,
	}
}

// Create creates a new user in the database
func (r *PostgresUserRepository) Create(ctx context.Context, user *models.User) error {
	if user == nil {
		return errors.New("user cannot be nil")
	}

	result := r.db.WithContext(ctx).Create(user)
	if result.Error != nil {
		return fmt.Errorf("failed to create user: %w", result.Error)
	}

	return nil
}

// GetByID retrieves a user by their ID
func (r *PostgresUserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	if id == "" {
		return nil, errors.New("invalid user ID")
	}

	var user models.User
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user with id %s not found: %w", id, gorm.ErrRecordNotFound)
		}
		return nil, fmt.Errorf("failed to get user by id %s: %w", id, result.Error)
	}

	return &user, nil
}

// GetByTelegramID retrieves a user by their Telegram ID
func (r *PostgresUserRepository) GetByTelegramID(ctx context.Context, telegramID int64) (*models.User, error) {
	if telegramID <= 0 {
		return nil, errors.New("invalid telegram ID")
	}

	var user models.User
	result := r.db.WithContext(ctx).Where("telegram_id = ?", telegramID).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user with telegram id %d not found: %w", telegramID, gorm.ErrRecordNotFound)
		}
		return nil, fmt.Errorf("failed to get user by telegram id %d: %w", telegramID, result.Error)
	}

	return &user, nil
}

// Update updates an existing user in the database
func (r *PostgresUserRepository) Update(ctx context.Context, user *models.User) error {
	if user == nil {
		return errors.New("user cannot be nil")
	}

	if user.ID == "" {
		return errors.New("invalid user ID")
	}

	result := r.db.WithContext(ctx).Save(user)
	if result.Error != nil {
		return fmt.Errorf("failed to update user: %w", result.Error)
	}

	return nil
}

// Delete deletes a user by their ID (soft delete)
func (r *PostgresUserRepository) Delete(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("invalid user ID")
	}

	// First get the user to ensure it exists
	var user models.User
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return fmt.Errorf("user with id %s not found", id)
		}
		return fmt.Errorf("failed to get user by id %s: %w", id, result.Error)
	}

	// Perform soft delete
	result = r.db.WithContext(ctx).Delete(&user)
	if result.Error != nil {
		return fmt.Errorf("failed to delete user: %w", result.Error)
	}

	return nil
}

// ExistsByTelegramID checks if a user with the given Telegram ID exists
func (r *PostgresUserRepository) ExistsByTelegramID(ctx context.Context, telegramID int64) (bool, error) {
	if telegramID <= 0 {
		return false, errors.New("invalid telegram ID")
	}

	var count int64
	result := r.db.WithContext(ctx).Model(&models.User{}).Where("telegram_id = ?", telegramID).Count(&count)
	if result.Error != nil {
		return false, fmt.Errorf("failed to check user existence by telegram id %d: %w", telegramID, result.Error)
	}

	return count > 0, nil
}

// Reset resets user's email and phone number by ID
func (r *PostgresUserRepository) ResetUser(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("invalid user ID")
	}

	result := r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", id).Updates(map[string]any{
		"email": "",
		"phone": "",
	})
	if result.Error != nil {
		return fmt.Errorf("failed to reset user data: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("user with id %s not found", id)
	}

	return nil
}

// Restore restores a soft-deleted user by setting deleted_at to NULL
func (r *PostgresUserRepository) RestoreUser(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("invalid user ID")
	}

	result := r.db.WithContext(ctx).Unscoped().Model(&models.User{}).Where("id = ? AND deleted_at IS NOT NULL", id).Update("deleted_at", nil)
	if result.Error != nil {
		return fmt.Errorf("failed to restore user: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("user with id %s not found", id)
	}

	return nil
}

// GetDeletedUsers retrieves all soft-deleted users
func (r *PostgresUserRepository) GetDeletedUsers(ctx context.Context) ([]*models.User, error) {
	var users []*models.User
	result := r.db.WithContext(ctx).Unscoped().Where("deleted_at IS NOT NULL").Find(&users)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get deleted users: %w", result.Error)
	}

	return users, nil
}
