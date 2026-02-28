package query

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent/bloodgroup"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/reference"
)

type GetBloodGroupsByPetTypeHandler struct {
	bloodInfoRepo reference.BloodInfoRepository
}

func NewGetBloodGroupsByPetTypeHandler(bloodInfoRepo reference.BloodInfoRepository) *GetBloodGroupsByPetTypeHandler {
	return &GetBloodGroupsByPetTypeHandler{
		bloodInfoRepo: bloodInfoRepo,
	}
}

func (h *GetBloodGroupsByPetTypeHandler) Handle(ctx context.Context, petType bloodgroup.PetType) ([]*ent.BloodGroup, error) {
	return h.bloodInfoRepo.BloodGroupsByPetType(ctx, petType)
}
