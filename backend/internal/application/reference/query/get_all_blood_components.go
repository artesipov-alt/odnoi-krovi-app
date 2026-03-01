package query

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/reference"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent"
)

type GetAllBloodComponentsHandler struct {
	bloodInfoRepo reference.BloodInfoRepository
}

func NewGetAllBloodComponentsHandler(bloodInfoRepo reference.BloodInfoRepository) *GetAllBloodComponentsHandler {
	return &GetAllBloodComponentsHandler{
		bloodInfoRepo: bloodInfoRepo,
	}
}

func (h *GetAllBloodComponentsHandler) Handle(ctx context.Context) ([]*ent.BloodComponent, error) {
	return h.bloodInfoRepo.AllComponents(ctx)
}
