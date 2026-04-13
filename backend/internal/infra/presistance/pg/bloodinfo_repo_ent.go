package pg

import (
	"context"
	"fmt"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent"
)

// EntBloodInfoRepository implements BloodInfoRepository using ENT
type EntBloodInfoRepository struct {
	client *ent.Client
}

// NewEntBloodInfoRepository creates a new ENT blood info repository
func NewEntBloodInfoRepository(client *ent.Client) *EntBloodInfoRepository {
	return &EntBloodInfoRepository{
		client: client,
	}
}

// AllComponents returns all blood components
func (r *EntBloodInfoRepository) AllComponents(ctx context.Context) ([]*ent.BloodComponent, error) {
	components, err := r.client.BloodComponent.Query().All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all blood components: %w", err)
	}
	return components, nil
}

// ComponentByID returns a blood component by ID
func (r *EntBloodInfoRepository) ComponentByID(ctx context.Context, id string) (*ent.BloodComponent, error) {
	component, err := r.client.BloodComponent.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("blood component with id %s not found: %w", id, err)
		}
		return nil, fmt.Errorf("failed to get blood component by id %s: %w", id, err)
	}
	return component, nil
}
