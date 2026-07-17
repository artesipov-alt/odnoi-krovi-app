package cmd

import (
	"context"
	"log/slog"
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/application/pet/enrich"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/events"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bonus"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor"
	donormodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet"
	petmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/ports"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/presistance"
)

// derefDonorResponses конвертирует []*DonorResponse в []DonorResponse для совместимости с BloodCounterService.
func derefDonorResponses(apps []*donormodel.DonorResponse) []donormodel.DonorResponse {
	result := make([]donormodel.DonorResponse, len(apps))
	for i, app := range apps {
		if app != nil {
			result[i] = *app
		}
	}
	return result
}

// SelectDonorHandler обрабатывает выбор донора реципиентом из списка потенциальных.
type SelectDonorHandler struct {
	bloodRepo    bloodsearch.Repository
	donorRepo    donor.Repository
	petRepo      pet.Repository
	petEnricher  enrich.PetEnricher
	userRepo     user.Repository
	bonusSvc     *bonus.BonusService
	publisher    ports.EventPublisher
	txManager    *presistance.TxManager
	bloodCounter *bloodsearch.BloodCounterService
}

func NewSelectDonorHandler(
	bloodRepo bloodsearch.Repository,
	donorRepo donor.Repository,
	petRepo pet.Repository,
	petEnricher enrich.PetEnricher,
	userRepo user.Repository,
	bonusSvc *bonus.BonusService,
	publisher ports.EventPublisher,
	txManager *presistance.TxManager,
) *SelectDonorHandler {
	return &SelectDonorHandler{
		bloodRepo:    bloodRepo,
		donorRepo:    donorRepo,
		petRepo:      petRepo,
		petEnricher:  petEnricher,
		userRepo:     userRepo,
		bonusSvc:     bonusSvc,
		publisher:    publisher,
		txManager:    txManager,
		bloodCounter: bloodsearch.NewBloodCounterService(),
	}
}

// Handle — реципиент выбирает конкретного донора из списка потенциальных.
// Создаёт DonorResponse сразу со статусом accepted (минуя pending),
// пересчитывает объёмы крови и статус заявки.
func (h *SelectDonorHandler) Handle(
	ctx context.Context,
	callerUserID string,
	requestID string,
	donorPetID string,
	compensationType string,
	taxiCompensation bool,
) (*donormodel.DonorResponse, error) {
	// 1. Получаем заявку
	req, err := h.bloodRepo.GetByID(ctx, requestID)
	if err != nil {
		return nil, err
	}
	if !req.IsActive() {
		return nil, apperrors.ErrInvalidBloodRequestStatus.WithMessage("blood request is not active")
	}

	// 2. Проверяем, что caller — владелец реципиента
	recipientPet, err := h.petRepo.GetByID(ctx, req.BloodRequest.PetID, pet.PetPreloadOptions{})
	if err != nil {
		return nil, apperrors.Internal(err, "failed to get recipient pet")
	}
	if recipientPet.OwnerID != callerUserID {
		return nil, apperrors.Forbidden("only the recipient owner can select donors")
	}

	// 3. Получаем питомца-донора и проверяем актуальный статус
	donorPet, err := h.petRepo.GetByID(ctx, donorPetID, pet.PetPreloadOptions{})
	if err != nil {
		return nil, apperrors.NotFound("donor pet not found")
	}
	// Пересчёт актуального статуса (БД может быть устаревшей)
	h.petEnricher.RecalculateOne(donorPet, nil, req, enrich.Options{RecoveryPeriodMonths: 0})
	if donorPet.PetStatus != petmodel.PetStatusDonor {
		return nil, apperrors.ErrInvalidBloodRequestStatus.WithMessage("donor pet is not currently available for donation")
	}

	// 4. Создаём DonorResponse напрямую со статусом accepted
	donorResponse, err := donormodel.NewDonorResponse(
		req.BloodRequest.ID,
		donorPet.ID,
		compensationType,
		donorPet.CalculateDonationAmount(),
		taxiCompensation,
	)
	if err != nil {
		return nil, apperrors.Validation(err.Error(), map[string]any{"field": "donor_response"})
	}
	donorResponse.Status = donormodel.DonorResponseStatusAccepted

	// 5. В транзакции: создаём response + пересчитываем объёмы + статус заявки
	err = h.txManager.WithTx(ctx, func(txCtx context.Context) error {
		saved, err := h.donorRepo.CreateDonorResponse(txCtx, donorResponse)
		if err != nil {
			return err
		}
		donorResponse = saved

		if err := h.bonusSvc.AssignBonuses(txCtx, donorPet.OwnerID, donorPet.Type, donorResponse.ID); err != nil {
			return err
		}

		// Перечитываем заявку со всеми приложениями, чтобы пересчитать объёмы
		fresh, err := h.bloodRepo.GetByID(txCtx, req.BloodRequest.ID)
		if err != nil {
			return err
		}

		applications, err := h.donorRepo.GetDonorResponsesByRequestID(txCtx, req.BloodRequest.ID)
		if err != nil {
			return err
		}

		donated, reserved := h.bloodCounter.RecalculateBloodAmount(fresh.BloodRequest, derefDonorResponses(applications))
		fresh.BloodRequest.SetBloodVolume(donated, reserved)
		fresh.RecalculateStatus()

		if err := h.bloodRepo.UpdateStatus(txCtx, fresh.BloodRequest.ID, fresh.BloodRequest.Status); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	// 6. Уведомление донору (вне транзакции — ошибка некритична)
	recipientUser, err := h.userRepo.GetByID(ctx, recipientPet.OwnerID, user.UserPreloadOptions{WithIdentities: true})
	if err != nil {
		slog.Error("failed to get recipient user for notification", "err", err)
		return donorResponse, nil
	}
	donorUser, err := h.userRepo.GetByID(ctx, donorPet.OwnerID, user.UserPreloadOptions{WithIdentities: true})
	if err != nil {
		slog.Error("failed to get donor user for notification", "err", err)
		return donorResponse, nil
	}

	donorMaxID, donorTgID := donorUser.MessengerContacts()
	recipientMaxID, recipientTgID := recipientUser.MessengerContacts()

	donorData := events.DonorData{
		UserName:         donorUser.FullName,
		PetName:          donorPet.Name,
		Phone:            donorUser.Phone,
		BloodGroup:       donorPet.BloodGroupName,
		ProviderMaxID:    donorMaxID,
		ProviderTelegram: donorTgID,
	}
	recipientData := events.RecipientData{
		UserName:         recipientUser.FullName,
		PetName:          recipientPet.Name,
		Phone:            recipientUser.Phone,
		BloodGroup:       recipientPet.BloodGroupName,
		Volume:           req.BloodRequest.BloodVolumeNeeded,
		ProviderMaxID:    recipientMaxID,
		ProviderTelegram: recipientTgID,
	}
	if err := h.publisher.PublishEvent(ctx, ports.EventDonorApply, events.ApplyDonor{
		DonorData:     donorData,
		RecipientData: recipientData,
		CreatedAt:     time.Now(),
	}); err != nil {
		slog.Error("failed to publish donor selected notification", "err", err, "donorResponseID", donorResponse.ID)
	}

	return donorResponse, nil
}
