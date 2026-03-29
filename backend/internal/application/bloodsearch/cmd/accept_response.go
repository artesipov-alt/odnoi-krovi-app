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

func (h *ApplyResponseHandler) Handle(ctx context.Context, donorResponseID string) error {
	bloodreq, err := h.bloodRepo.GetByApplicationID(ctx, donorResponseID)
	if err != nil {
		return err
	}

	application, err := h.donorRepo.GetDonorResponseByID(ctx, donorResponseID)
	if err != nil {
		return err
	}

	bloodreq.ReserveVolume(application.Amount)

	err = h.txManager.WithTx(ctx, func(txCtx context.Context) error {
		if err := h.donorRepo.UpdateDonorResponseStatus(txCtx, donorResponseID, donormodel.DonorResponseStatusAccepted); err != nil {
			return err
		}
		if err := h.bloodRepo.UpdateReservedVolume(txCtx, bloodreq.ID, bloodreq.BloodVolumeReserved); err != nil {
			return err
		}
		if err := h.bloodRepo.UpdateStatus(txCtx, bloodreq.ID, bloodreq.Status); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	recipientPet, err := h.petRepo.GetByID(ctx, bloodreq.PetID, pet.PetPreloadOptions{})
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

	event := events.ApplyDonor{}
	for _, identity := range donorUser.Identities {
		if identity.ProviderName == authmodel.ProviderMax {
			event.DonorData.ProviderMaxID = identity.ProviderUserID
		}
		if identity.ProviderName == authmodel.ProviderTelegram {
			event.DonorData.ProviderTelegram = identity.ProviderUserID
		}
	}
	event.DonorData.Name = donorUser.FullName
	event.DonorData.Phone = donorUser.Phone
	for _, identity := range recipientUser.Identities {
		if identity.ProviderName == authmodel.ProviderMax {
			event.RecipientData.ProviderMaxID = identity.ProviderUserID
		}
		if identity.ProviderName == authmodel.ProviderTelegram {
			event.RecipientData.ProviderTelegram = identity.ProviderUserID
		}

	}
	event.RecipientData.Name = recipientUser.FullName
	event.RecipientData.Phone = recipientUser.Phone

	if err := h.publisher.PublishDonorApply(ctx, event); err != nil {
		return err
	}

	return nil
}
