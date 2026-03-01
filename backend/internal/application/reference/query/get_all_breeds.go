package query

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/reference"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent"
)

type GetAllBreedsHandler struct {
	breedRepo reference.BreedInfoRepository
}

func NewGetAllBreedsHandler(breedRepo reference.BreedInfoRepository) *GetAllBreedsHandler {
	return &GetAllBreedsHandler{
		breedRepo: breedRepo,
	}
}

func (h *GetAllBreedsHandler) Handle(ctx context.Context) ([]*ent.Breed, error) {
	return h.breedRepo.GetAll(ctx)
}
