package cmd

import (
	"context"
	"log/slog"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/application/bloodsearch/service"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bonus"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/presistance"
)

type CreateRequestHandler struct {
	bloodRepo bloodsearch.Repository
	petRepo   pet.Repository
	bonusRepo bonus.Repository
	notifySvc *service.DonorMatchNotifier
	txManager *presistance.TxManager
}

func NewCreateRequestHandler(
	bloodRepo bloodsearch.Repository,
	petRepo pet.Repository,
	bonusRepo bonus.Repository,
	notifySvc *service.DonorMatchNotifier,
	txManager *presistance.TxManager,
) *CreateRequestHandler {
	return &CreateRequestHandler{
		bloodRepo: bloodRepo,
		petRepo:   petRepo,
		bonusRepo: bonusRepo,
		notifySvc: notifySvc,
		txManager: txManager,
	}
}

func (h *CreateRequestHandler) Handle(ctx context.Context, req *model.BloodRequest) (*model.BloodRequestWithApplications, error) {
	// Проверяем существование питомца
	petRecipient, err := h.petRepo.GetByID(ctx, req.PetID, pet.PetPreloadOptions{})
	if err != nil {
		return nil, err
	}

	// Проверяем, нет ли уже активной заявки для этого питомца
	activeExists, err := h.bloodRepo.ExistsByPetID(ctx, req.PetID)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to check request existence")
	}
	if activeExists {
		return nil, apperrors.ErrBloodRequestAlreadyExists
	}

	var newReq *model.BloodRequestWithApplications
	err = h.txManager.WithTx(ctx, func(txCtx context.Context) error {
		var err error
		if req.PrioritySearch && petRecipient.Privilege == "" {
			err = h.bonusRepo.SubtractPrioritySearch(txCtx, petRecipient.OwnerID)
			if err != nil {
				return apperrors.Internal(err, "failed to subtract priority search")
			}
		}
		newReq, err = h.bloodRepo.Create(txCtx, req)
		if err != nil {
			return apperrors.Internal(err, "failed to create blood request")
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				slog.Error("recovered from panic", "error", r)
			}
		}()
		detachedCtx := context.WithoutCancel(ctx)
		if err := h.notifySvc.NotifyMatchDonors(detachedCtx, newReq); err != nil {
			slog.Error("failed to notify match donors", "error", err)
		}
	}()

	return newReq, nil
}
