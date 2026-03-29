package query

import (
	"context"
	"errors"
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user"
)

type GetByUserHandler struct {
	petReadRepo   pet.PetReadRepository
	userRepo      user.Repository
	donorRespRepo donor.Repository
	bloodReqRepo  bloodsearch.BloodRequestRepository
}

func NewGetByUserHandler(
	petReadRepo pet.PetReadRepository,
	userRepo user.Repository,
	donorRespRepo donor.Repository,
	bloodReqRepo bloodsearch.BloodRequestRepository,

) *GetByUserHandler {
	return &GetByUserHandler{
		petReadRepo:   petReadRepo,
		userRepo:      userRepo,
		bloodReqRepo:  bloodReqRepo,
		donorRespRepo: donorRespRepo,
	}
}

func (h *GetByUserHandler) Handle(ctx context.Context, userID string, opts pet.PetPreloadOptions) ([]*model.Pet, error) {
	exists, err := h.userRepo.ExistsByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, apperrors.ErrUserNotFound
	}

	pets, err := h.petReadRepo.GetByUserID(ctx, userID, opts)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to get pets")
	}

	for _, pet := range pets {
		application, err := h.donorRespRepo.GetByPetID(ctx, pet.ID)
		if err != nil && !errors.Is(err, apperrors.ErrDonorResponseNotFound) {
			return nil, apperrors.Internal(err, "failed to get donor application")
		}
		bloodReq, err := h.bloodReqRepo.GetByPetID(ctx, pet.ID)
		if err != nil && !errors.Is(err, apperrors.ErrBloodRequestNotFound) {
			return nil, apperrors.Internal(err, "failed to get blood request")
		}
		pet.RecalculateFactors(time.Now(), application, bloodReq)
		pet.CalculateStatus(application, bloodReq)
	}

	return pets, nil
}
