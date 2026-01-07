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
