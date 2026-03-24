package cmd

import (
	"context"

	authmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/auth/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/events"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/ports"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor"
	donormodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user"
)

type ApplyResponseHandler struct {
	bloodRepo bloodsearch.BloodRequestRepository
	donorRepo donor.Repository
	petRepo   pet.Repository
	userRepo  user.Repository
	publisher ports.EventPublisher
}

func NewApplyResponseHandler(
	bloodRepo bloodsearch.BloodRequestRepository,
	donorRepo donor.Repository,
	petRepo pet.Repository,
	userRepo user.Repository,
	publisher ports.EventPublisher,
) *ApplyResponseHandler {
	return &ApplyResponseHandler{
		bloodRepo: bloodRepo,
		donorRepo: donorRepo,
		petRepo:   petRepo,
		userRepo:  userRepo,
		publisher: publisher,
	}
}

func (h *ApplyResponseHandler) Handle(ctx context.Context, donorResponseID string) error {
	req, err := h.bloodRepo.GetByApplicationID(ctx, donorResponseID)
	if err != nil {
		return err
	}
	if req.BloodVolumeReserved >= req.BloodVolumeNeeded {
		req.Close()
	}
	if err := h.donorRepo.UpdateDonorResponseStatus(ctx, donorResponseID, donormodel.DonorResponseStatusAccepted); err != nil {
		return err
	}

	donorPet, err := h.petRepo.GetByID(ctx, req.PetID, pet.PetPreloadOptions{})
	if err != nil {
		return err
	}
	donorUser, err := h.userRepo.GetByID(ctx, donorPet.OwnerID, user.UserPreloadOptions{
		WithIdentities: true,
	})
	if err != nil {
		return err
	}
	recipientPet, err := h.petRepo.GetByID(ctx, req.PetID, pet.PetPreloadOptions{})
	if err != nil {
		return err
	}
	recipientUser, err := h.userRepo.GetByID(ctx, recipientPet.OwnerID, user.UserPreloadOptions{
		WithIdentities: true,
	})
	if err != nil {
		return err
	}

	event := events.ApplyDonor{}
	for _, identity := range donorUser.Identities {
		if identity.ProviderName == authmodel.ProviderMax {
			event.DonorData.ProviderMaxID = identity.ProviderUserID
		}
		if identity.ProviderName == authmodel.ProviderTelegram {
			event.DonorData.ProviderTelegram = identity.ProviderUserID
		}
	}
	event.DonorData.Phone = donorUser.Phone
	for _, identity := range recipientUser.Identities {
		if identity.ProviderName == authmodel.ProviderMax {
			event.RecipientData.ProviderMaxID = identity.ProviderUserID
		}
		if identity.ProviderName == authmodel.ProviderTelegram {
			event.RecipientData.ProviderTelegram = identity.ProviderUserID
		}

	}
	event.RecipientData.Phone = recipientUser.Phone

	if err := h.publisher.PublishDonorApply(ctx, event); err != nil {
		return err
	}

	return nil
}
