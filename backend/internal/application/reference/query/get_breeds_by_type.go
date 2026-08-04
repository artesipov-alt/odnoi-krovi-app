package query

import (
	"cmp"
	"context"
	"slices"

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

func (h *GetBreedsByPetTypeHandler) Handle(ctx context.Context, petType string) ([]*ent.Breed, error) {
	breeds, err := h.breedRepo.GetByPetType(ctx, breed.Type(petType))
	if err != nil {
		return nil, err
	}

	// Сортировка: "МЕТИС" первым, остальные по алфавиту
	slices.SortStableFunc(breeds, func(a, b *ent.Breed) int {
		// "МЕТИС" всегда первый
		if a.Name == "МЕТИС" {
			return -1
		}
		if b.Name == "МЕТИС" {
			return 1
		}

		// Остальные — по алфавиту
		return cmp.Compare(a.Name, b.Name)
	})

	return breeds, nil
}
