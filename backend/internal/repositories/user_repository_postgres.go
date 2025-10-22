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
func (r *PostgresUserRepository) GetByID(ctx context.Context, id int) (*models.User, error) {
	if id <= 0 {
		return nil, errors.New("invalid user ID")
	}

	var user models.User
	result := r.db.WithContext(ctx).First(&user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user with id %d not found: %w", id, gorm.ErrRecordNotFound)
		}
		return nil, fmt.Errorf("failed to get user by id %d: %w", id, result.Error)
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

	if user.ID <= 0 {
		return errors.New("invalid user ID")
	}

	result := r.db.WithContext(ctx).Save(user)
	if result.Error != nil {
		return fmt.Errorf("failed to update user: %w", result.Error)
	}

	return nil
}

// Delete deletes a user by their ID
func (r *PostgresUserRepository) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("invalid user ID")
	}

	result := r.db.WithContext(ctx).Delete(&models.User{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete user: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("user with id %d not found", id)
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
