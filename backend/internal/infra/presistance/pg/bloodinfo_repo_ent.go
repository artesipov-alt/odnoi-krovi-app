package pg

import (
	"context"
	"fmt"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent/bloodgroup"
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

// BloodGroupsByPetType returns blood groups by pet type
func (r *EntBloodInfoRepository) BloodGroupsByPetType(ctx context.Context, petType bloodgroup.PetType) ([]*ent.BloodGroup, error) {
	groups, err := r.client.BloodGroup.Query().
		Where(bloodgroup.PetTypeEQ(petType)).
		All(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get blood groups for pet type %s: %w", petType, err)
	}
	return groups, nil
}

// FindByTypeAndBloodGroup returns a blood group by pet type and blood group value
func (r *EntBloodInfoRepository) FindByTypeAndBloodGroup(ctx context.Context, petType bloodgroup.PetType, bloodGroup string) (*ent.BloodGroup, error) {
	group, err := r.client.BloodGroup.Query().
		Where(bloodgroup.PetTypeEQ(petType)).
		Where(bloodgroup.BloodGroupEQ(bloodGroup)).
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("blood group %s for pet type %s not found: %w", bloodGroup, petType, err)
		}
		return nil, fmt.Errorf("failed to get blood group %s for pet type %s: %w", bloodGroup, petType, err)
	}
	return group, nil
}

// FindByBloodGroup returns a blood group by blood group value only
func (r *EntBloodInfoRepository) FindByBloodGroup(ctx context.Context, bloodGroup string) (*ent.BloodGroup, error) {
	group, err := r.client.BloodGroup.Query().
		Where(bloodgroup.BloodGroupEQ(bloodGroup)).
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("blood group %s not found: %w", bloodGroup, err)
		}
		return nil, fmt.Errorf("failed to get blood group %s: %w", bloodGroup, err)
	}
	return group, nil
}

// FindByName returns a blood group by name
func (r *EntBloodInfoRepository) FindByName(ctx context.Context, bloodGroup string) (*ent.BloodGroup, error) {
	group, err := r.client.BloodGroup.Query().
		Where(bloodgroup.BloodGroupEQ(bloodGroup)).
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("blood group %s not found: %w", bloodGroup, err)
		}
		return nil, fmt.Errorf("failed to get blood group %s: %w", bloodGroup, err)
	}
	return group, nil
}
