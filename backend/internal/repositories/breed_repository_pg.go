package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/models"
	"gorm.io/gorm"
)

// PostgresBreedRepository implements BreedRepository for PostgreSQL
type PostgresBreedRepository struct {
	db *gorm.DB
}

// NewPostgresBreedRepository creates a new PostgreSQL breed repository
func NewPostgresBreedRepository(db *gorm.DB) *PostgresBreedRepository {
	return &PostgresBreedRepository{
		db: db,
	}
}

// GetAll returns all breeds from the database
func (r *PostgresBreedRepository) GetAll(ctx context.Context) ([]*models.Breed, error) {
	var breeds []*models.Breed
	result := r.db.WithContext(ctx).Find(&breeds)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get all breeds: %w", result.Error)
	}

	return breeds, nil
}

// GetByID retrieves a breed by its ID
func (r *PostgresBreedRepository) GetByID(ctx context.Context, id int) (*models.Breed, error) {
	if id <= 0 {
		return nil, errors.New("invalid breed ID")
	}

	var breed models.Breed
	result := r.db.WithContext(ctx).First(&breed, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("breed with id %d not found: %w", id, gorm.ErrRecordNotFound)
		}
		return nil, fmt.Errorf("failed to get breed by id %d: %w", id, result.Error)
	}

	return &breed, nil
}

// GetByPetType retrieves breeds by pet type
func (r *PostgresBreedRepository) GetByPetType(ctx context.Context, petType models.PetType) ([]*models.Breed, error) {
	if petType == "" {
		return nil, errors.New("pet type cannot be empty")
	}

	var breeds []*models.Breed
	result := r.db.WithContext(ctx).Where("type = ?", petType).Find(&breeds)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get breeds for pet type %s: %w", petType, result.Error)
	}

	return breeds, nil
}

// Create creates a new breed in the database
func (r *PostgresBreedRepository) Create(ctx context.Context, breed *models.Breed) error {
	if breed == nil {
		return errors.New("breed cannot be nil")
	}

	result := r.db.WithContext(ctx).Create(breed)
	if result.Error != nil {
		return fmt.Errorf("failed to create breed: %w", result.Error)
	}

	return nil
}

// Update updates an existing breed in the database
func (r *PostgresBreedRepository) Update(ctx context.Context, breed *models.Breed) error {
	if breed == nil {
		return errors.New("breed cannot be nil")
	}

	if breed.ID <= 0 {
		return errors.New("invalid breed ID")
	}

	result := r.db.WithContext(ctx).Save(breed)
	if result.Error != nil {
		return fmt.Errorf("failed to update breed: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("breed with id %d not found", breed.ID)
	}

	return nil
}

// Delete deletes a breed by its ID
func (r *PostgresBreedRepository) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("invalid breed ID")
	}

	// First get the breed to ensure it exists
	var breed models.Breed
	result := r.db.WithContext(ctx).First(&breed, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return fmt.Errorf("breed with id %d not found", id)
		}
		return fmt.Errorf("failed to get breed by id %d: %w", id, result.Error)
	}

	// Perform delete
	result = r.db.WithContext(ctx).Delete(&breed)
	if result.Error != nil {
		return fmt.Errorf("failed to delete breed: %w", result.Error)
	}

	return nil
}

// ExistsByName checks if a breed with the given name exists
func (r *PostgresBreedRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	if name == "" {
		return false, errors.New("breed name cannot be empty")
	}

	var count int64
	result := r.db.WithContext(ctx).Model(&models.Breed{}).Where("name = ?", name).Count(&count)
	if result.Error != nil {
		return false, fmt.Errorf("failed to check breed existence by name %s: %w", name, result.Error)
	}

	return count > 0, nil
}
