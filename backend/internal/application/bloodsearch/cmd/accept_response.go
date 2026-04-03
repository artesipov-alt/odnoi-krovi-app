package cmd

import (
	"context"
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/events"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor"
	donormodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/ports"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user"
	userevent "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user/events"
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
	application, err := h.donorRepo.GetDonorResponseByID(ctx, donorResponseID)
	if err != nil {
		return err
	}

	var petID string
	err = h.txManager.WithTx(ctx, func(txCtx context.Context) error {
		if err := h.donorRepo.UpdateDonorResponseStatus(txCtx, donorResponseID, donormodel.DonorResponseStatusAccepted); err != nil {
			return err
		}

		bloodreq, err := h.bloodRepo.GetByApplicationID(txCtx, donorResponseID)
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

	eventRecipient := userevent.GenerateContact(recipientUser)
	eventDonor := userevent.GenerateContact(donorUser)

	event := events.ApplyDonor{
		DonorData:     eventDonor.UserData,
		RecipientData: eventRecipient.UserData,
		CreatedAt:     time.Now(),
	}

	if err := h.publisher.PublishDonorApply(ctx, event); err != nil {
		return err
	}

	return nil
}
