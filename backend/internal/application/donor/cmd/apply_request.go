package cmd

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	authmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/auth/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"
	bloodmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/model"
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
	bloodRepo bloodsearch.BloodRequestRepository
	petRepo   pet.Repository
	donorRepo donor.Repository
	userRepo  user.Repository
	bonusSvc  *bonus.BonusService
	publisher ports.EventPublisher
	txManager *presistance.TxManager
}

func NewApplyForRequestHandler(
	bloodRepo bloodsearch.BloodRequestRepository,
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
	if req.Status != bloodmodel.BloodRequestStatusActive {
		return nil, apperrors.ErrInvalidBloodRequestStatus.WithMessage("blood request is not active")
	}

	// Проверяем существование донора
	donorPet, err := h.petRepo.GetByID(ctx, donorID, pet.PetPreloadOptions{})
	if err != nil {
		return nil, apperrors.Internal(err, "failed to check donor existence")
	}

	// Получаем данные пользователя донора
	donorUser, err := h.userRepo.GetByID(ctx, donorPet.OwnerID, user.UserPreloadOptions{
		WithDonorPreference: true,
	})
	if err != nil {
		return nil, apperrors.Internal(err, "failed to get donor user")
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
	var recipientProviderMaxID string
	for _, identity := range recipientUser.Identities {
		if identity.ProviderName == authmodel.ProviderMax {
			recipientProviderMaxID = identity.ProviderUserID
			break
		}
	}

	err = h.txManager.WithTx(ctx, func(ctx context.Context) error {
		donorResponse, err = h.donorRepo.CreateDonorResponse(ctx, donorResponse)
		if err != nil {
			return err
		}

		// Закрепляем бонусы за пользователем
		if err := h.bonusSvc.AssignBonuses(ctx, donorPet.OwnerID, donorPet.Type, donorUser.LastDonation); err != nil {
			return err
		}

		donorBloodGroup := donorPet.BloodGroupName

		if err := h.publisher.PublishRecipientApply(ctx, donorevent.RecipientApply{
			DonorName:                       donorPet.Name,
			DonorBloodGroup:                 donorBloodGroup,
			RecipientProviderMaxID:          recipientProviderMaxID,
			RecipientPetName:                recipientPet.Name,
			RecipientPetSearchingBloodGroup: req.BloodGroupNames,
			RecipientPetNeededVolume:        req.BloodVolumeNeeded,
			CreatedAt:                       *donorResponse.CreatedAt,
		}); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return donorResponse, nil
}
