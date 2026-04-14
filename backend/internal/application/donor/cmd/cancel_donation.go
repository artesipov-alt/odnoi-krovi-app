package cmd

import (
	"context"
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	authmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/auth/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor"
	donorevent "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/events"
	donormodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/ports"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user"
	usermodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/presistance"
)

type CancelDonationHandler struct {
	donorRepo      donor.Repository
	bloodRepo      bloodsearch.BloodRequestRepository
	txManager      *presistance.TxManager
	eventPublisher ports.EventPublisher
	petRepo        pet.PetReadRepository
	userRepo       user.Repository
}

func NewCancelDonationHandler(
	donorRepo donor.Repository,
	bloodRepo bloodsearch.BloodRequestRepository,
	txManager *presistance.TxManager,
	eventPublisher ports.EventPublisher,
	petRepo pet.PetReadRepository,
	userRepo user.Repository,
) *CancelDonationHandler {
	return &CancelDonationHandler{
		donorRepo:      donorRepo,
		bloodRepo:      bloodRepo,
		txManager:      txManager,
		eventPublisher: eventPublisher,
		petRepo:        petRepo,
		userRepo:       userRepo,
	}
}

func extractProviderIDs(user *usermodel.User) (maxID, telegramID string) {
	for _, identity := range user.Identities {
		if identity.ProviderName == authmodel.ProviderMax {
			maxID = identity.ProviderUserID
		}
		if identity.ProviderName == authmodel.ProviderTelegram {
			telegramID = identity.ProviderUserID
		}
	}
	return
}

func (h *CancelDonationHandler) Handle(ctx context.Context, resID string) error {
	// Получаем DonorResponse
	donorResponse, err := h.donorRepo.GetDonorResponseByID(ctx, resID)
	if err != nil {
		return apperrors.Internal(err, "failed to get donor response")
	}
	if donorResponse.Status == donormodel.DonorResponseStatusCompleted {
		return apperrors.BadRequest("donor response status is invalid").WithMessage("cannot cancel a completed donation")
	}

	// Donor name and blood group from response

	err = h.txManager.WithTx(ctx, func(txCtx context.Context) error {
		if err := h.donorRepo.Cancel(txCtx, resID); err != nil {
			return apperrors.Internal(err, "failed to cancel donation")
		}

		// Пересчитываем статус заявки после отмены отклика
		bloodReq, err := h.bloodRepo.GetByApplicationID(txCtx, resID)
		if err != nil {
			return apperrors.Internal(err, "failed to get blood request after cancel")
		}
		bloodReq.RecalculateBloodAmount()
		bloodReq.RecalculateStatus()
		if err := h.bloodRepo.UpdateStatus(txCtx, bloodReq.ID, bloodReq.Status); err != nil {
			return apperrors.Internal(err, "failed to update blood request status after cancel")
		}

		return nil
	})

	if err != nil {
		return err
	}

	// Publish DonorCancel event
	bloodReq, err := h.bloodRepo.GetByApplicationID(ctx, resID)
	if err != nil {
		return apperrors.Internal(err, "failed to get blood request for event")
	}

	recipientPet, err := h.petRepo.GetByID(ctx, bloodReq.PetID, pet.PetPreloadOptions{})
	if err != nil {
		return apperrors.Internal(err, "failed to get recipient pet")
	}

	recipientUser, err := h.userRepo.GetByID(ctx, recipientPet.OwnerID, user.UserPreloadOptions{
		WithIdentities: true,
	})
	if err != nil {
		return apperrors.Internal(err, "failed to get recipient user")
	}

	recipientProviderMaxID, _ := extractProviderIDs(recipientUser)

	event := donorevent.DonorCancel{
		DonorName:              donorResponse.DonorName,
		DonorBloodGroup:        donorResponse.DonorBloodGroup,
		RecipientProviderMaxID: recipientProviderMaxID,
		RecipientPetName:       recipientPet.Name,
		CreatedAt:              time.Now(),
	}

	if err := h.eventPublisher.PublishDonorCancel(ctx, event); err != nil {
		return apperrors.Internal(err, "failed to publish donor cancel event")
	}

	return nil
}
