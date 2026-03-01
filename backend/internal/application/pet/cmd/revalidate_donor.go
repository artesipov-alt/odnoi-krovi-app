package cmd

import (
	"context"
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet/model"
)

type RevalidateDonorHandler struct {
	petReadRepo  pet.PetReadRepository
	petWriteRepo pet.PetWriteRepository
}

func NewRevalidateDonorHandler(petReadRepo pet.PetReadRepository, petWriteRepo pet.PetWriteRepository) *RevalidateDonorHandler {
	return &RevalidateDonorHandler{
		petReadRepo:  petReadRepo,
		petWriteRepo: petWriteRepo,
	}
}

func (h *RevalidateDonorHandler) Handle(ctx context.Context, petID string) (*model.Pet, error) {
	p, err := h.petReadRepo.GetByID(ctx, petID, pet.PetPreloadOptions{
		WithAll: true,
	})
	if err != nil {
		return nil, err
	}

	stopFactors := p.GetStaticStopFactors()
	warnFactors := p.GetWarnFactors(time.Now())

	// Создаём структуру только с полями для обновления
	updatePet := &model.Pet{}
	updatePet.StopFactors = make([]string, len(stopFactors))
	for i, f := range stopFactors {
		updatePet.StopFactors[i] = string(f)
	}
	updatePet.WarnFactors = make([]string, len(warnFactors))
	for i, f := range warnFactors {
		updatePet.WarnFactors[i] = string(f)
	}

	return h.petWriteRepo.Update(ctx, petID, updatePet)
}
