package cmd

import (
	"context"
	"log/slog"
	"time"

	authmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/auth/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/events"
	bloodmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/ports"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user"
	usermodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/presistance"
)

type ApplyResponseHandler struct {
	bloodRepo bloodsearch.BloodRequestRepository
	donorRepo donor.Repository
	petRepo   pet.Repository
	userRepo  user.Repository
	publisher ports.EventPublisher
	txManager *presistance.TxManager
}

func NewApplyResponseHandler(
	bloodRepo bloodsearch.BloodRequestRepository,
	donorRepo donor.Repository,
	petRepo pet.Repository,
	userRepo user.Repository,
	publisher ports.EventPublisher,
	txManager *presistance.TxManager,
) *ApplyResponseHandler {
	return &ApplyResponseHandler{
		bloodRepo: bloodRepo,
		donorRepo: donorRepo,
		petRepo:   petRepo,
		userRepo:  userRepo,
		publisher: publisher,
		txManager: txManager,
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

func (h *ApplyResponseHandler) Handle(ctx context.Context, donorResponseID string) error {
	application, err := h.donorRepo.GetDonorResponseByID(ctx, donorResponseID)
	if err != nil {
		return err
	}

	var petID string
	var bloodreq *bloodmodel.BloodRequestWithApplications

	err = h.txManager.WithTx(ctx, func(txCtx context.Context) error {
		if err := application.Accept(); err != nil {
			return err
		}
		if err := h.donorRepo.Accept(txCtx, donorResponseID); err != nil {
			return err
		}

		bloodreq, err = h.bloodRepo.GetByApplicationID(txCtx, donorResponseID, false)
		if err != nil {
			return err
		}

		bloodreq.RecalculateBloodAmount()
		bloodreq.RecalculateStatus()

		if err := h.bloodRepo.UpdateStatus(txCtx, bloodreq.ID, bloodreq.Status); err != nil {
			return err
		}
		petID = bloodreq.PetID
		return nil
	})
	if err != nil {
		return err
	}

	recipientPet, err := h.petRepo.GetByID(ctx, petID, pet.PetPreloadOptions{})
	if err != nil {
		return err
	}
	recipientUser, err := h.userRepo.GetByID(ctx, recipientPet.OwnerID, user.UserPreloadOptions{
		WithIdentities: true,
	})
	if err != nil {
		return err
	}

	donorPet, err := h.petRepo.GetByID(ctx, application.DonorID, pet.PetPreloadOptions{})
	if err != nil {
		return err
	}
	donorUser, err := h.userRepo.GetByID(ctx, donorPet.OwnerID, user.UserPreloadOptions{
		WithIdentities: true,
	})
	if err != nil {
		return err
	}

	donorMaxID, donorTelegramID := extractProviderIDs(donorUser)
	recipientMaxID, recipientTelegramID := extractProviderIDs(recipientUser)

	donorData := events.DonorData{
		UserName:         donorUser.FullName,
		PetName:          donorPet.Name,
		Phone:            donorUser.Phone,
		BloodGroup:       donorPet.BloodGroupName,
		ProviderMaxID:    donorMaxID,
		ProviderTelegram: donorTelegramID,
	}

	recipientData := events.RecipientData{
		UserName:         recipientUser.FullName,
		PetName:          recipientPet.Name,
		Phone:            recipientUser.Phone,
		BloodGroup:       recipientPet.BloodGroupName,
		Volume:           bloodreq.BloodVolumeNeeded,
		ProviderMaxID:    recipientMaxID,
		ProviderTelegram: recipientTelegramID,
	}

	event := events.ApplyDonor{
		DonorData:     donorData,
		RecipientData: recipientData,
		CreatedAt:     time.Now(),
	}

	if err := h.publisher.PublishDonorApply(ctx, event); err != nil {
		slog.Error("failed to publish donor apply notification", "err", err, "donorResponseID", donorResponseID)
	}

	return nil
}
