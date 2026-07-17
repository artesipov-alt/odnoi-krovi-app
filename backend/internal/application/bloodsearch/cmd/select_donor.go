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

// SelectDonorHandler обрабатывает выбор донора реципиентом из списка потенциальных.
type SelectDonorHandler struct {
	bloodRepo   bloodsearch.Repository
	donorRepo   donor.Repository
	petRepo     pet.PetReadRepository
	petEnricher enrich.PetEnricher
	userRepo    user.Repository
	bonusSvc    *bonus.BonusService
	publisher   ports.EventPublisher
	txManager   *presistance.TxManager
}

func NewSelectDonorHandler(
	bloodRepo bloodsearch.Repository,
	donorRepo donor.Repository,
	petRepo pet.PetReadRepository,
	petEnricher enrich.PetEnricher,
	userRepo user.Repository,
	bonusSvc *bonus.BonusService,
	publisher ports.EventPublisher,
	txManager *presistance.TxManager,
) *SelectDonorHandler {
	return &SelectDonorHandler{
		bloodRepo:   bloodRepo,
		donorRepo:   donorRepo,
		petRepo:     petRepo,
		petEnricher: petEnricher,
		userRepo:    userRepo,
		bonusSvc:    bonusSvc,
		publisher:   publisher,
		txManager:   txManager,
	}
}

// Handle — реципиент выбирает конкретного донора из списка потенциальных.
// Создаёт DonorResponse сразу со статусом accepted (минуя pending),
// пересчитывает объёмы крови и статус заявки.
// Условия донации (compensationType, taxiCompensation) берутся из
// DonorPreference владельца донора.
func (h *SelectDonorHandler) Handle(
	ctx context.Context,
	callerUserID string,
	requestID string,
	donorPetID string,
) (*donormodel.DonorResponse, error) {
	// 1. Получаем заявку
	req, err := h.bloodRepo.GetByID(ctx, requestID)
	if err != nil {
		return nil, err
	}
	if !req.IsActive() {
		return nil, apperrors.ErrInvalidBloodRequestStatus.WithMessage("blood request is not active")
	}

	// 2. Получаем реципиента и проверяем, что caller — владелец
	recipientPet, err := h.petRepo.GetByID(ctx, req.BloodRequest.PetID, pet.PetPreloadOptions{})
	if err != nil {
		return nil, apperrors.NotFound("recipient pet not found")
	}
	if recipientPet.OwnerID != callerUserID {
		return nil, apperrors.Forbidden("only the recipient owner can select donors")
	}

	// 3. Получаем питомца-донора и пересчитываем его актуальный статус.
	// ВАЖНО: статус донора считаем по ЕГО СОБСТВЕННЫМ данным (его DonorResponse
	// на чужие заявки + его собственные активные BloodRequest), а не по заявке
	// реципиента — иначе peekStatus в Pet.peekStatus переводит донора в Recipient/BloodFound.
	donorPet, err := h.petRepo.GetByID(ctx, donorPetID, pet.PetPreloadOptions{WithAll: true})
	if err != nil {
		return nil, apperrors.NotFound("donor pet not found")
	}
	// Тянем владельца донора один раз: и DonorPreference (для enrich + условий донации),
	// и Identities (для уведомления после транзакции) — два обхода БД не нужны.
	donorOwner, err := h.userRepo.GetByID(ctx, donorPet.OwnerID, user.UserPreloadOptions{
		WithDonorPreference: true,
		WithIdentities:      true,
	})
	if err != nil {
		return nil, apperrors.NotFound("donor owner not found")
	}
	fc, err := h.petEnricher.Fetch(ctx, []string{donorPet.ID})
	if err != nil {
		return nil, apperrors.Internal(err, "failed to load donor context")
	}
	// Если у владельца нет DonorPreference (новый аккаунт без настроек), берём дефолт
	// из usermodel.DefaultDonorPreference(), чтобы RecalculateRecoveryDays пересчитал дни.
	recoveryMonths := 2
	if donorOwner.DonorPreference != nil && donorOwner.DonorPreference.RecoveryPeriodMonths > 0 {
		recoveryMonths = donorOwner.DonorPreference.RecoveryPeriodMonths
	}
	h.petEnricher.Recalculate(donorPet, fc, enrich.Options{RecoveryPeriodMonths: recoveryMonths})
	if donorPet.PetStatus != petmodel.PetStatusDonor {
		return nil, apperrors.Validation(
			"donor pet is not currently available for donation",
			map[string]any{"field": "donor_pet", "donor_id": donorPetID},
		)
	}

	// Условия донации берём из DonorPreference владельца донора (а не из тела запроса).
	compensationType := ""
	taxiCompensation := false
	if donorOwner.DonorPreference != nil {
		compensationType = string(donorOwner.DonorPreference.CompensationType)
		taxiCompensation = donorOwner.DonorPreference.TaxiCompensation
	}

	// 4. Создаём DonorResponse со статусом accepted через доменный метод
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
	if err := donorResponse.Accept(); err != nil {
		return nil, apperrors.Validation(err.Error(), map[string]any{"field": "donor_response"})
	}

	// 5. Транзакция: создать response, начислить бонусы, пересчитать объёмы и статус
	err = h.txManager.WithTx(ctx, func(txCtx context.Context) error {
		saved, err := h.donorRepo.CreateDonorResponse(txCtx, donorResponse)
		if err != nil {
			return err
		}
		donorResponse = saved

		if err := h.bonusSvc.AssignBonuses(txCtx, donorPet.OwnerID, donorPet.Type, donorResponse.ID); err != nil {
			return err
		}

		// Перечитываем заявку вместе с приложениями (донор-респонсами),
		// чтобы пересчитать объёмы и сохранить актуальный статус заявки.
		fresh, err := h.bloodRepo.GetByID(txCtx, req.BloodRequest.ID)
		if err != nil {
			return err
		}

		bloodCounter := bloodsearch.NewBloodCounterService()
		donated, reserved := bloodCounter.RecalculateBloodAmount(fresh.BloodRequest, fresh.DonorApplications)
		fresh.BloodRequest.SetBloodVolume(donated, reserved)
		fresh.RecalculateStatus()

		return h.bloodRepo.UpdateStatus(txCtx, fresh.BloodRequest.ID, fresh.BloodRequest.Status)
	})
	if err != nil {
		return nil, err
	}

	// 6. Уведомление донору (вне транзакции — ошибка некритична).
	// donorOwner уже загружен на шаге 3 с WithIdentities=true — переиспользуем.
	recipientUser, err := h.userRepo.GetByID(ctx, recipientPet.OwnerID, user.UserPreloadOptions{WithIdentities: true})
	if err != nil {
		slog.Error("failed to get recipient user for notification", "err", err)
		return donorResponse, nil
	}

	donorMaxID, donorTgID := donorOwner.MessengerContacts()
	recipientMaxID, recipientTgID := recipientUser.MessengerContacts()

	donorData := events.DonorData{
		UserName:         donorOwner.FullName,
		PetName:          donorPet.Name,
		Phone:            donorOwner.Phone,
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
