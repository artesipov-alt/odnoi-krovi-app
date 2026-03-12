package query

import (
	"context"
	"errors"
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user"
)

type GetByUserHandler struct {
	petReadRepo  pet.PetReadRepository
	userRepo     user.Repository
	bloodReqRepo bloodsearch.BloodRequestRepository
}

func NewGetByUserHandler(
	petReadRepo pet.PetReadRepository,
	userRepo user.Repository,
	bloodReqRepo bloodsearch.BloodRequestRepository,

) *GetByUserHandler {
	return &GetByUserHandler{
		petReadRepo:  petReadRepo,
		userRepo:     userRepo,
		bloodReqRepo: bloodReqRepo,
	}
}

func (h *GetByUserHandler) Handle(ctx context.Context, userID string, opts pet.PetPreloadOptions) ([]*model.Pet, error) {
	_, err := h.userRepo.GetByID(ctx, userID, user.UserPreloadOptions{})
	if err != nil {
		// Репозиторий уже возвращает доменные ошибки
		return nil, err
	}

	pets, err := h.petReadRepo.GetByUserID(ctx, userID, opts)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to get pets")
	}

	for i := range pets {
		bloodReq, err := h.bloodReqRepo.GetByPetID(ctx, pets[i].ID)
		if err != nil && !errors.Is(err, apperrors.ErrBloodRequestNotFound) {
			return nil, err
		}
		pets[i].RecalculateFactors(time.Now())
		// Calculate status using domain method
		hasActiveRequest := bloodReq != nil
		hasResponses := hasActiveRequest && len(bloodReq.ResponseIDs) > 0
		pets[i].PetStatus = pets[i].CalculateStatus(time.Now(), hasActiveRequest, hasResponses)
	}

	return pets, nil
}
