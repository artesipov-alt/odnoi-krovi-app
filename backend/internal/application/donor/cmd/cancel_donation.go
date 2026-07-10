package cmd

import (
	"context"
	"log/slog"
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"
	bloodreqmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bonus"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor"
	donorevent "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/events"
	donormodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/ports"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/presistance"
)

type CancelDonationHandler struct {
	donorRepo      donor.Repository
	bloodRepo      bloodsearch.Repository
	txManager      *presistance.TxManager
	eventPublisher ports.EventPublisher
	petRepo        pet.PetReadRepository
	userRepo       user.Repository
	bonusSvc       *bonus.BonusService
}

func NewCancelDonationHandler(
	donorRepo donor.Repository,
	bloodRepo bloodsearch.Repository,
	txManager *presistance.TxManager,
	eventPublisher ports.EventPublisher,
	petRepo pet.PetReadRepository,
	userRepo user.Repository,
	bonusSvc *bonus.BonusService,
) *CancelDonationHandler {
	return &CancelDonationHandler{
		donorRepo:      donorRepo,
		bloodRepo:      bloodRepo,
		txManager:      txManager,
		eventPublisher: eventPublisher,
		petRepo:        petRepo,
		userRepo:       userRepo,
		bonusSvc:       bonusSvc,
	}
}

func (h *CancelDonationHandler) Handle(ctx context.Context, resID string, reason string) error {
	// Получаем DonorResponse
	donorResponse, err := h.donorRepo.GetDonorResponseByID(ctx, resID)
	if err != nil {
		return apperrors.Internal(err, "failed to get donor response")
	}
	if donorResponse.IsCompleted() {
		return apperrors.BadRequest("donor response status is invalid").WithMessage("cannot cancel a completed donation")
	}

	// Получаем данные донора
	donorPet, err := h.petRepo.GetByID(ctx, donorResponse.DonorID, pet.PetPreloadOptions{})
	if err != nil {
		return apperrors.Internal(err, "failed to get donor pet")
	}

	// Собираем read-only данные для уведомления до транзакции
	recipientPet, recipientProviderMaxID, recipientProviderTelegramID, err := h.collectRecipientData(ctx, donorResponse)
	if err != nil {
		return err
	}

	// Транзакция: отмена отклика с пересчётом заявки (свежие applications) и отмена бонусов
	var bloodReq *bloodreqmodel.BloodRequestWithApplications
	err = h.txManager.WithTx(ctx, func(txCtx context.Context) error {
		if err := donorResponse.Cancel(reason); err != nil {
			return apperrors.Internal(err, "failed to cancel donation")
		}
		if err := h.donorRepo.Update(txCtx, donorResponse); err != nil {
			return apperrors.Internal(err, "failed to cancel donation")
		}

		// Запрос внутри транзакции — подтягивает актуальный список DonorApplications
		bloodReq, err = h.bloodRepo.GetByApplicationID(txCtx, resID, false)
		if err != nil {
			return apperrors.Internal(err, "failed to get blood request after cancel")
		}
		bloodReq.RecalculateBloodAmount()
		bloodReq.BloodRequest.RecalculateStatus()
		if err := h.bloodRepo.UpdateStatus(txCtx, bloodReq.BloodRequest.ID, bloodReq.BloodRequest.Status); err != nil {
			return apperrors.Internal(err, "failed to update blood request status after cancel")
		}

		if err := h.bonusSvc.UnassignReservedBonuses(txCtx, donorPet.OwnerID, donorPet.Type); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	// Уведомление реципиента после успешной транзакции.
	// Ошибка публикации не фатальна — логируем и продолжаем.
	if err := h.eventPublisher.PublishEvent(ctx, ports.EventDonorCancel, donorevent.DonorCancel{
		DonorName:                   donorResponse.DonorName,
		DonorBloodGroup:             donorResponse.DonorBloodGroup,
		RecipientProviderMaxID:      recipientProviderMaxID,
		RecipientProviderTelegramID: recipientProviderTelegramID,
		RecipientPetName:            recipientPet.Name,
		CreatedAt:                   time.Now(),
	}); err != nil {
		slog.Error("failed to publish donor cancel notification", "err", err, "responseID", resID)
	}

	return nil
}

func (h *CancelDonationHandler) collectRecipientData(ctx context.Context, response *donormodel.DonorResponse) (*model.Pet, string, string, error) {
	bloodReq, err := h.bloodRepo.GetByApplicationID(ctx, response.ID, false)
	if err != nil {
		return nil, "", "", apperrors.Internal(err, "failed to get blood request")
	}

	recipientPet, err := h.petRepo.GetByID(ctx, bloodReq.BloodRequest.PetID, pet.PetPreloadOptions{})
	if err != nil {
		return nil, "", "", apperrors.Internal(err, "failed to get recipient pet")
	}

	recipientUser, err := h.userRepo.GetByID(ctx, recipientPet.OwnerID, user.UserPreloadOptions{
		WithIdentities: true,
	})
	if err != nil {
		return nil, "", "", apperrors.Internal(err, "failed to get recipient user")
	}

	maxID, telegramID := recipientUser.MessengerContacts()
	return recipientPet, maxID, telegramID, nil
}
