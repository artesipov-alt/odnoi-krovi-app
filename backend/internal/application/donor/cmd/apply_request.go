package cmd

import (
	"context"
	"log/slog"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bonus"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor"
	donorevent "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/events"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/model"
	donormodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/ports"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/presistance"
)

type ApplyForRequestHandler struct {
	bloodRepo bloodsearch.Repository
	petRepo   pet.Repository
	donorRepo donor.Repository
	userRepo  user.Repository
	bonusSvc  *bonus.BonusService
	publisher ports.EventPublisher
	txManager *presistance.TxManager
}

func NewApplyForRequestHandler(
	bloodRepo bloodsearch.Repository,
	petRepo pet.Repository,
	donorRepo donor.Repository,
	userRepo user.Repository,
	bonusSvc *bonus.BonusService,
	publisher ports.EventPublisher,
	txManager *presistance.TxManager,
) *ApplyForRequestHandler {
	return &ApplyForRequestHandler{
		bloodRepo: bloodRepo,
		petRepo:   petRepo,
		donorRepo: donorRepo,
		userRepo:  userRepo,
		bonusSvc:  bonusSvc,
		publisher: publisher,
		txManager: txManager,
	}
}

func (h *ApplyForRequestHandler) Handle(ctx context.Context, reqID, donorID, compensationType string, taxiCompensation bool) (*model.DonorResponse, error) {
	// Проверяем существование и статус заявки
	req, err := h.bloodRepo.GetByID(ctx, reqID)
	if err != nil {
		return nil, err
	}
	if !req.IsActive() {
		return nil, apperrors.ErrInvalidBloodRequestStatus.WithMessage("blood request is not active")
	}

	// Проверяем существование донора
	donorPet, err := h.petRepo.GetByID(ctx, donorID, pet.PetPreloadOptions{})
	if err != nil {
		return nil, apperrors.Internal(err, "failed to check donor existence")
	}

	// Создаём новый отклик донора
	donorResponse, err := donormodel.NewDonorResponse(req.ID, donorPet.ID, compensationType, donorPet.CalculateDonationAmount(), taxiCompensation)
	if err != nil {
		return nil, apperrors.Validation(err.Error(), map[string]any{"field": "donor_response"})
	}

	// Получаем данные реципиента
	recipientPet, err := h.petRepo.GetByID(ctx, req.PetID, pet.PetPreloadOptions{})
	if err != nil {
		return nil, apperrors.Internal(err, "failed to get recipient pet")
	}

	recipientUser, err := h.userRepo.GetByID(ctx, recipientPet.OwnerID, user.UserPreloadOptions{
		WithIdentities: true,
	})
	if err != nil {
		return nil, apperrors.Internal(err, "failed to get recipient user")
	}

	err = h.txManager.WithTx(ctx, func(txCtx context.Context) error {
		donorResponse, err = h.donorRepo.CreateDonorResponse(txCtx, donorResponse)
		if err != nil {
			return err
		}

		// Закрепляем бонусы за пользователем
		if err := h.bonusSvc.AssignBonuses(txCtx, donorPet.OwnerID, donorPet.Type, donorResponse.ID); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	recipientProviderMaxID, recipientProviderTelegramID := recipientUser.MessengerContacts()

	// Отправляем уведомление реципиенту после успешной транзакции.
	// Ошибка публикации не фатальна — логируем и продолжаем.
	if err := h.publisher.PublishEvent(ctx, ports.EventRecipientApply, donorevent.RecipientApply{
		DonorName:                       donorPet.Name,
		DonorBloodGroup:                 donorPet.BloodGroupName,
		RecipientProviderMaxID:          recipientProviderMaxID,
		RecipientProviderTelegramID:     recipientProviderTelegramID,
		RecipientPetName:                recipientPet.Name,
		RecipientPetSearchingBloodGroup: req.BloodGroupNames,
		RecipientPetNeededVolume:        req.BloodVolumeNeeded,
		CreatedAt:                       *donorResponse.CreatedAt,
	}); err != nil {
		slog.Error("failed to publish recipient apply notification", "err", err, "donorResponseID", donorResponse.ID)
	}

	return donorResponse, nil
}
