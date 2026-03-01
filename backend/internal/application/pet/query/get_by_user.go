package query

import (
	"context"
	"errors"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/filestorage"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent"
)

type GetByUserHandler struct {
	petReadRepo  pet.PetReadRepository
	userRepo     user.Repository
	bloodReqRepo bloodsearch.BloodRequestRepository
	storage      filestorage.Repository
}

func NewGetByUserHandler(
	petReadRepo pet.PetReadRepository,
	userRepo user.Repository,
	bloodReqRepo bloodsearch.BloodRequestRepository,
	storage filestorage.Repository,
) *GetByUserHandler {
	return &GetByUserHandler{
		petReadRepo:  petReadRepo,
		userRepo:     userRepo,
		bloodReqRepo: bloodReqRepo,
		storage:      storage,
	}
}

func (h *GetByUserHandler) Handle(ctx context.Context, userID string, opts pet.PetPreloadOptions) ([]*model.Pet, error) {
	_, _, err := h.userRepo.GetByID(ctx, userID, user.UserPreloadOptions{})
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, apperrors.ErrUserNotFound
		}
		return nil, apperrors.Internal(err, "failed to get user")
	}

	pets, err := h.petReadRepo.GetByUserID(ctx, userID, opts)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to get pets")
	}

	for i := range pets {
		pets[i].PhotoURLs = h.storage.BuildPhotoURLs(pets[i].PhotoURLs, *pets[i].UpdatedAt)

		bloodReq, err := h.bloodReqRepo.GetByPetID(ctx, pets[i].ID)
		if err != nil && !errors.Is(err, apperrors.ErrBloodRequestNotFound) {
			return nil, err
		}

		// Set status based on blood request and stored DonorRestrictions
		if bloodReq != nil {
			if len(bloodReq.ResponseIDs) > 0 {
				pets[i].PetStatus = model.PetStatusBloodFound
			} else {
				pets[i].PetStatus = model.PetStatusRecipient
			}
		} else {
			if len(pets[i].StopFactors) > 0 {
				pets[i].PetStatus = model.PetStatusNone
			} else {
				pets[i].PetStatus = model.PetStatusDonor
			}
		}
	}

	return pets, nil
}
