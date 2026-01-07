package pg

import (
	"context"
	"errors"
	"fmt"

	"github.com/artesipov-alt/odnoi-krovi-app/ent"
	"github.com/artesipov-alt/odnoi-krovi-app/ent/location"
)

// EntLocationRepository implements LocationRepository using ENT
type EntLocationRepository struct {
	client *ent.Client
}

// NewEntLocationRepository creates a new ENT location repository
func NewEntLocationRepository(client *ent.Client) *EntLocationRepository {
	return &EntLocationRepository{
		client: client,
	}
}

// GetByID retrieves a location by its ID
func (r *EntLocationRepository) GetByID(ctx context.Context, id int) (*ent.Location, error) {
	if id <= 0 {
		return nil, errors.New("invalid location ID")
	}

	l, err := r.client.Location.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("location with id %d not found: %w", id, err)
		}
		return nil, fmt.Errorf("failed to get location by id %d: %w", id, err)
	}

	return l, nil
}

// GetAll retrieves all locations from the database
func (r *EntLocationRepository) GetAll(ctx context.Context) ([]*ent.Location, error) {
	locations, err := r.client.Location.Query().
		Order(ent.Asc(location.FieldName)).
		All(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get all locations: %w", err)
	}

	return locations, nil
}
