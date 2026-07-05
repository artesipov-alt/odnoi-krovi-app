package cmd

import (
	"context"
	"log/slog"
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor"
	donorevent "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/events"
	donormodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet"
	petmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/ports"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user"
)

type CompleteDonationHandler struct {
	donorRepo donor.Repository
	bloodRepo bloodsearch.Repository
	petRepo   pet.Repository
	userRepo  user.Repository
	publisher ports.EventPublisher
}

func NewCompleteDonationHandler(
	donorRepo donor.Repository,
	bloodRepo bloodsearch.Repository,
	petRepo pet.Repository,
	userRepo user.Repository,
	publisher ports.EventPublisher,
) *CompleteDonationHandler {
	return &CompleteDonationHandler{
		donorRepo: donorRepo,
		bloodRepo: bloodRepo,
		petRepo:   petRepo,
		userRepo:  userRepo,
		publisher: publisher,
	}
}

func (h *CompleteDonationHandler) Handle(ctx context.Context, resID string, amount float64) error {
	// Получаем DonorResponse
	donorResponse, err := h.donorRepo.GetDonorResponseByID(ctx, resID)
	if err != nil {
		return apperrors.Internal(err, "failed to get donor response")
	}
	if donorResponse.Status != donormodel.DonorResponseStatusAccepted {
		return apperrors.BadRequest("donor response status is invalid").WithMessage("donor response must be accepted to complete donation")
	}

	// Собираем read-only данные для уведомления до записи
	recipientPet, recipientProviderMaxID, recipientProviderTelegramID, err := h.collectRecipientData(ctx, resID)
	if err != nil {
		return err
	}

	// Выполняем complete
	if err := donorResponse.Complete(amount); err != nil {
		return err
	}
	if err := h.donorRepo.Update(ctx, donorResponse); err != nil {
		return err
	}

	// Уведомление реципиента после успешного завершения.
	// Ошибка публикации не фатальна — логируем и продолжаем.
	if err := h.publisher.PublishEvent(ctx, ports.EventDonorCompleted, donorevent.DonorCompleted{
		DonorPetName:                donorResponse.DonorName,
		DonorBloodGroup:             donorResponse.DonorBloodGroup,
		RecipientProviderMaxID:      recipientProviderMaxID,
		RecipientProviderTelegramID: recipientProviderTelegramID,
		RecipientPetName:            recipientPet.Name,
		Amount:                      amount,
		CreatedAt:                   time.Now(),
	}); err != nil {
		slog.Error("failed to publish donor completed notification", "err", err, "responseID", resID)
	}

	return nil
}

func (h *CompleteDonationHandler) collectRecipientData(ctx context.Context, resID string) (*petmodel.Pet, string, string, error) {
	bloodReq, err := h.bloodRepo.GetByApplicationID(ctx, resID, false)
	if err != nil {
		return nil, "", "", apperrors.Internal(err, "failed to get blood request")
	}

	recipientPet, err := h.petRepo.GetByID(ctx, bloodReq.PetID, pet.PetPreloadOptions{})
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
