package pg

import (
	"context"
	"fmt"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent/location"
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
func (r *EntLocationRepository) GetByID(ctx context.Context, id string) (*ent.Location, error) {
	l, err := r.client.Location.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("location with id %s not found: %w", id, err)
		}
		return nil, fmt.Errorf("failed to get location by id %s: %w", id, err)
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

// Exists checks if a location with the given ID exists
func (r *EntLocationRepository) Exists(ctx context.Context, id string) (bool, error) {
	exists, err := r.client.Location.Query().Where(location.ID(id)).Exist(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to check if location with id %s exists: %w", id, err)
	}
	return exists, nil
}
