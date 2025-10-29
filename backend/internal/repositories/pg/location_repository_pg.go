package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/models"
	"gorm.io/gorm"
)

// PostgresLocationRepository implements LocationRepository for PostgreSQL
type PostgresLocationRepository struct {
	db *gorm.DB
}

// NewPostgresLocationRepository creates a new PostgreSQL location repository
func NewPostgresLocationRepository(db *gorm.DB) *PostgresLocationRepository {
	return &PostgresLocationRepository{
		db: db,
	}
}

// Create creates a new location in the database
func (r *PostgresLocationRepository) Create(ctx context.Context, location *models.Location) error {
	if location == nil {
		return errors.New("location cannot be nil")
	}

	result := r.db.WithContext(ctx).Create(location)
	if result.Error != nil {
		return fmt.Errorf("failed to create location: %w", result.Error)
	}

	return nil
}

// GetByID retrieves a location by its ID
func (r *PostgresLocationRepository) GetByID(ctx context.Context, id int) (*models.Location, error) {
	if id <= 0 {
		return nil, errors.New("invalid location ID")
	}

	var location models.Location
	result := r.db.WithContext(ctx).First(&location, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("location with id %d not found: %w", id, gorm.ErrRecordNotFound)
		}
		return nil, fmt.Errorf("failed to get location by id %d: %w", id, result.Error)
	}

	return &location, nil
}

// GetAll retrieves all locations from the database
func (r *PostgresLocationRepository) GetAll(ctx context.Context) ([]*models.Location, error) {
	var locations []*models.Location
	result := r.db.WithContext(ctx).Find(&locations)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get all locations: %w", result.Error)
	}

	return locations, nil
}

// Update updates an existing location in the database
func (r *PostgresLocationRepository) Update(ctx context.Context, location *models.Location) error {
	if location == nil {
		return errors.New("location cannot be nil")
	}

	if location.ID <= 0 {
		return errors.New("invalid location ID")
	}

	result := r.db.WithContext(ctx).Save(location)
	if result.Error != nil {
		return fmt.Errorf("failed to update location: %w", result.Error)
	}

	return nil
}

// Delete deletes a location by its ID
func (r *PostgresLocationRepository) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("invalid location ID")
	}

	// First get the location to ensure it exists
	var location models.Location
	result := r.db.WithContext(ctx).First(&location, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return fmt.Errorf("location with id %d not found", id)
		}
		return fmt.Errorf("failed to get location by id %d: %w", id, result.Error)
	}

	// Perform delete
	result = r.db.WithContext(ctx).Delete(&location)
	if result.Error != nil {
		return fmt.Errorf("failed to delete location: %w", result.Error)
	}

	return nil
}

// ExistsByName checks if a location with the given name exists
func (r *PostgresLocationRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	if name == "" {
		return false, errors.New("location name cannot be empty")
	}

	var count int64
	result := r.db.WithContext(ctx).Model(&models.Location{}).Where("name = ?", name).Count(&count)
	if result.Error != nil {
		return false, fmt.Errorf("failed to check location existence by name %s: %w", name, result.Error)
	}

	return count > 0, nil
}
