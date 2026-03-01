package query

import (
	"context"
	"errors"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/filestorage"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet/model"
)

type GetByIDHandler struct {
	petReadRepo  pet.PetReadRepository
	bloodReqRepo bloodsearch.BloodRequestRepository
	storage      filestorage.Repository
}

func NewGetByIDHandler(petReadRepo pet.PetReadRepository, bloodReqRepo bloodsearch.BloodRequestRepository, storage filestorage.Repository) *GetByIDHandler {
	return &GetByIDHandler{
		petReadRepo:  petReadRepo,
		bloodReqRepo: bloodReqRepo,
		storage:      storage,
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

	p.PhotoURLs = h.storage.BuildPhotoURLs(p.PhotoURLs, *p.UpdatedAt)

	// Set status based on blood request and stored DonorRestrictions
	if bloodReq != nil {
		if len(bloodReq.ResponseIDs) > 0 {
			p.PetStatus = model.PetStatusBloodFound
		} else {
			p.PetStatus = model.PetStatusRecipient
		}
	} else {
		if len(p.StopFactors) > 0 {
			p.PetStatus = model.PetStatusNone
		} else {
			p.PetStatus = model.PetStatusDonor
		}
	}

	return p, nil
}
