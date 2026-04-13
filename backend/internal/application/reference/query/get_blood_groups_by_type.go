package query

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/common"
)

type GetBloodGroupsByPetTypeHandler struct{}

func NewGetBloodGroupsByPetTypeHandler() *GetBloodGroupsByPetTypeHandler {
	return &GetBloodGroupsByPetTypeHandler{}
}

func (h *GetBloodGroupsByPetTypeHandler) Handle(ctx context.Context, petType common.PetType) ([]string, error) {
	return common.GetBloodGroupsByPetType(petType), nil
}
