package cmd

import (
	"context"
	"log/slog"
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/events"
	bloodmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/ports"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/presistance"
)

type ApplyResponseHandler struct {
	bloodRepo    bloodsearch.Repository
	donorRepo    donor.Repository
	petRepo      pet.Repository
	userRepo     user.Repository
	publisher    ports.EventPublisher
	txManager    *presistance.TxManager
	bloodCounter *bloodsearch.BloodCounterService
}

func NewApplyResponseHandler(
	bloodRepo bloodsearch.Repository,
	donorRepo donor.Repository,
	petRepo pet.Repository,
	userRepo user.Repository,
	publisher ports.EventPublisher,
	txManager *presistance.TxManager,
) *ApplyResponseHandler {
	return &ApplyResponseHandler{
		bloodRepo:    bloodRepo,
		donorRepo:    donorRepo,
		petRepo:      petRepo,
		userRepo:     userRepo,
		publisher:    publisher,
		txManager:    txManager,
		bloodCounter: bloodsearch.NewBloodCounterService(),
	}
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
		if err := h.donorRepo.Update(txCtx, application); err != nil {
			return err
		}

		bloodreq, err = h.bloodRepo.GetByApplicationID(txCtx, donorResponseID, false)
		if err != nil {
			return err
		}

		donated, reserved := h.bloodCounter.RecalculateBloodAmount(bloodreq.BloodRequest, bloodreq.DonorApplications)

		bloodreq.BloodRequest.SetBloodVolume(donated, reserved)
		bloodreq.RecalculateStatus()

		if err := h.bloodRepo.UpdateStatus(txCtx, bloodreq.BloodRequest.ID, bloodreq.BloodRequest.Status); err != nil {
			return err
		}
		petID = bloodreq.BloodRequest.PetID
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

	donorMaxID, donorTelegramID := donorUser.MessengerContacts()
	recipientMaxID, recipientTelegramID := recipientUser.MessengerContacts()

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
		Volume:           bloodreq.BloodRequest.BloodVolumeNeeded,
		ProviderMaxID:    recipientMaxID,
		ProviderTelegram: recipientTelegramID,
	}

	event := events.ApplyDonor{
		DonorData:     donorData,
		RecipientData: recipientData,
		CreatedAt:     time.Now(),
	}

	if err := h.publisher.PublishEvent(ctx, ports.EventDonorApply, event); err != nil {
		slog.Error("failed to publish donor apply notification", "err", err, "donorResponseID", donorResponseID)
	}

	return nil
}
