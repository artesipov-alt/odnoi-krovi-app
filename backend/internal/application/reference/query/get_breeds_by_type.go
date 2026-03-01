package query

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/reference"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent/breed"
)

type GetBreedsByPetTypeHandler struct {
	breedRepo reference.BreedInfoRepository
}

func NewGetBreedsByPetTypeHandler(breedRepo reference.BreedInfoRepository) *GetBreedsByPetTypeHandler {
	return &GetBreedsByPetTypeHandler{
		breedRepo: breedRepo,
	}
}

func (h *GetBreedsByPetTypeHandler) Handle(ctx context.Context, petType breed.Type) ([]*ent.Breed, error) {
	return h.breedRepo.GetByPetType(ctx, petType)
}
