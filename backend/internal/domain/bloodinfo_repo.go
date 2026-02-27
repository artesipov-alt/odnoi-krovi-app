package domain

import "context"

// BloodInfoRepository определяет интерфейс для работы с типами крови
type BloodInfoRepository interface {
	AllComponents(ctx context.Context) ([]*ent.BloodComponent, error)
	ComponentByID(ctx context.Context, id string) (*ent.BloodComponent, error)
	BloodGroupsByPetType(ctx context.Context, petType bloodgroup.PetType) ([]*ent.BloodGroup, error)
	FindByTypeAndBloodGroup(ctx context.Context, petType bloodgroup.PetType, bloodGroup string) (*ent.BloodGroup, error)
	FindByBloodGroup(ctx context.Context, bloodGroup string) (*ent.BloodGroup, error)
	FindByName(ctx context.Context, name string) (*ent.BloodGroup, error)
}
