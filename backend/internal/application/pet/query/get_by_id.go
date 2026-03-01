package query

import (
	"context"
	"errors"
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet/model"
)

type GetByIDHandler struct {
	petReadRepo  pet.PetReadRepository
	bloodReqRepo bloodsearch.BloodRequestRepository
}

func NewGetByIDHandler(petReadRepo pet.PetReadRepository, bloodReqRepo bloodsearch.BloodRequestRepository) *GetByIDHandler {
	return &GetByIDHandler{
		petReadRepo:  petReadRepo,
		bloodReqRepo: bloodReqRepo,
	}
}

func (h *GetByIDHandler) Handle(ctx context.Context, petID string, opts pet.PetPreloadOptions) (*model.Pet, error) {
	p, err := h.petReadRepo.GetByID(ctx, petID, opts)
	if err != nil {
		return nil, err
	}

	bloodReq, err := h.bloodReqRepo.GetByPetID(ctx, petID)
	if err != nil && !errors.Is(err, apperrors.ErrBloodRequestNotFound) {
		return nil, err
	}

	// Calculate status using domain method
	hasActiveRequest := bloodReq != nil
	hasResponses := hasActiveRequest && len(bloodReq.ResponseIDs) > 0
	p.PetStatus = p.CalculateStatus(time.Now(), hasActiveRequest, hasResponses)

	return p, nil
}
