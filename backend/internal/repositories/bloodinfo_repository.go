package repositories

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/ent"
	"github.com/artesipov-alt/odnoi-krovi-app/ent/bloodgroup"
)

// BloodInfoRepository определяет интерфейс для работы с типами крови
type BloodInfoRepository interface {
	AllComponents(ctx context.Context) ([]*ent.BloodComponent, error)
	ComponentByID(ctx context.Context, id string) (*ent.BloodComponent, error)
	BloodGroupsByPetType(ctx context.Context, petType bloodgroup.PetType) ([]*ent.BloodGroup, error)
	FindByTypeAndBloodGroup(ctx context.Context, petType bloodgroup.PetType, bloodGroup string) (*ent.BloodGroup, error)
}
