package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/models"
	"gorm.io/gorm"
)

// PostgresPetRepository implements PetRepository for PostgreSQL
type PostgresPetRepository struct {
	db *gorm.DB
}

// NewPostgresPetRepository creates a new PostgreSQL pet repository
func NewPostgresPetRepository(db *gorm.DB) *PostgresPetRepository {
	return &PostgresPetRepository{
		db: db,
	}
}

// Create creates a new pet in the database
func (r *PostgresPetRepository) Create(ctx context.Context, pet *models.Pet) error {
	if pet == nil {
		return errors.New("pet cannot be nil")
	}

	result := r.db.WithContext(ctx).Create(pet)
	if result.Error != nil {
		return fmt.Errorf("failed to create pet: %w", result.Error)
	}

	return nil
}

// GetByID retrieves a pet by their ID
func (r *PostgresPetRepository) GetByID(ctx context.Context, id int) (*models.Pet, error) {
	if id <= 0 {
		return nil, errors.New("invalid pet ID")
	}

	var pet models.Pet
	result := r.db.WithContext(ctx).First(&pet, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("pet with id %d not found: %w", id, gorm.ErrRecordNotFound)
		}
		return nil, fmt.Errorf("failed to get pet by id %d: %w", id, result.Error)
	}

	return &pet, nil
}

// GetByUserID retrieves all pets for a specific user
func (r *PostgresPetRepository) GetByUserID(ctx context.Context, userID int) ([]*models.Pet, error) {
	if userID <= 0 {
		return nil, errors.New("invalid user ID")
	}

	var pets []*models.Pet
	result := r.db.WithContext(ctx).Where("owner_id = ?", userID).Find(&pets)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get pets for user %d: %w", userID, result.Error)
	}

	return pets, nil
}

// Update updates an existing pet in the database
func (r *PostgresPetRepository) Update(ctx context.Context, pet *models.Pet) error {
	if pet == nil {
		return errors.New("pet cannot be nil")
	}

	if pet.ID <= 0 {
		return errors.New("invalid pet ID")
	}

	result := r.db.WithContext(ctx).Save(pet)
	if result.Error != nil {
		return fmt.Errorf("failed to update pet: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("pet with id %d not found", pet.ID)
	}

	return nil
}

// Delete deletes a pet by their ID
func (r *PostgresPetRepository) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("invalid pet ID")
	}

	// First get the pet to ensure it exists
	var pet models.Pet
	result := r.db.WithContext(ctx).First(&pet, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return fmt.Errorf("pet with id %d not found", id)
		}
		return fmt.Errorf("failed to get pet by id %d: %w", id, result.Error)
	}

	// Perform delete
	result = r.db.WithContext(ctx).Delete(&pet)
	if result.Error != nil {
		return fmt.Errorf("failed to delete pet: %w", result.Error)
	}

	return nil
}

// ExistsByID checks if a pet with the given ID exists
func (r *PostgresPetRepository) ExistsByID(ctx context.Context, id int) (bool, error) {
	if id <= 0 {
		return false, errors.New("invalid pet ID")
	}

	var count int64
	result := r.db.WithContext(ctx).Model(&models.Pet{}).Where("id = ?", id).Count(&count)
	if result.Error != nil {
		return false, fmt.Errorf("failed to check pet existence by id %d: %w", id, result.Error)
	}

	return count > 0, nil
}
