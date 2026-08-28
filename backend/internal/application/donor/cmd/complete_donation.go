package cmd

import (
	"context"
	"log/slog"
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"
	bloodsearchmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/model"
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

	// Проверяем, что заявка ещё активна — донацию нельзя завершить на закрытой заявке.
	// Без этой проверки донор может повторно complete'нуть донацию, если отклик
	// вернулся в accepted после закрытия заявки (баг в CloseRequestHandler до фикса).
	bloodReq, err := h.bloodRepo.GetByApplicationID(ctx, resID, false)
	if err != nil {
		return apperrors.Internal(err, "failed to get blood request")
	}
	if !bloodReq.IsActive() {
		return apperrors.BadRequest("blood request is not active").WithMessage("нельзя завершить донацию на закрытой заявке")
	}

	// Собираем read-only данные для уведомления до записи.
	// Заявка уже загружена выше — переиспользуем её, чтобы не делать лишний запрос.
	recipientPet, recipientProviderMaxID, recipientProviderTelegramID, err := h.collectRecipientDataFromReq(ctx, bloodReq)
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

// collectRecipientDataFromReq собирает данные реципиента из уже загруженной заявки —
// экономит один запрос к БД по сравнению с collectRecipientData.
func (h *CompleteDonationHandler) collectRecipientDataFromReq(ctx context.Context, bloodReq *bloodsearchmodel.BloodRequestWithApplications) (*petmodel.Pet, string, string, error) {
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
