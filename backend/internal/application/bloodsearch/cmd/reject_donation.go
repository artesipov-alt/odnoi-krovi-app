package cmd

import (
	"context"
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bonus"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor"
	donorevent "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/events"
	donormodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet"
	petmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/ports"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/presistance"
)

type RejectDonationHandler struct {
	bloodRepo bloodsearch.BloodRequestRepository
	donorRepo donor.Repository
	petRepo   pet.Repository
	userRepo  user.Repository
	txManager *presistance.TxManager
	publisher ports.EventPublisher
	bonusSvc  *bonus.BonusService
}

func NewRejectDonationHandler(
	bloodRepo bloodsearch.BloodRequestRepository,
	donorRepo donor.Repository,
	petRepo pet.Repository,
	userRepo user.Repository,
	txManager *presistance.TxManager,
	publisher ports.EventPublisher,
	bonusSvc *bonus.BonusService,
) *RejectDonationHandler {
	return &RejectDonationHandler{
		bloodRepo: bloodRepo,
		donorRepo: donorRepo,
		petRepo:   petRepo,
		userRepo:  userRepo,
		txManager: txManager,
		publisher: publisher,
		bonusSvc:  bonusSvc,
	}
}

func (h *RejectDonationHandler) Handle(ctx context.Context, donorResponseID string, rejectedReason string) error {
	application, err := h.donorRepo.GetDonorResponseByID(ctx, donorResponseID)
	if err != nil {
		return err
	}

	if err := application.Reject(rejectedReason); err != nil {
		return err
	}

	var donorPet *petmodel.Pet
	err = h.txManager.WithTx(ctx, func(txCtx context.Context) error {
		if err := h.donorRepo.Reject(txCtx, application); err != nil {
			return err
		}
		bloodReq, err := h.bloodRepo.GetByApplicationID(txCtx, donorResponseID, false)
		if err != nil {
			return err
		}
		bloodReq.RecalculateBloodAmount()
		bloodReq.RecalculateStatus()
		if err := h.bloodRepo.UpdateStatus(txCtx, bloodReq.ID, bloodReq.Status); err != nil {
			return err
		}
		donorPet, err = h.petRepo.GetByID(txCtx, application.DonorID, pet.PetPreloadOptions{})
		if err != nil {
			return err
		}
		if application.Status == donormodel.DonorResponseStatusRejected {
			if err := h.bonusSvc.UnassignReservedBonuses(txCtx, donorPet.OwnerID, donorPet.Type); err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	// Publish DonorReject event
	bloodReq, err := h.bloodRepo.GetByApplicationID(ctx, donorResponseID, false)
	if err != nil {
		return err
	}

	recipientPet, err := h.petRepo.GetByID(ctx, bloodReq.PetID, pet.PetPreloadOptions{})
	if err != nil {
		return err
	}

	donorUser, err := h.userRepo.GetByID(ctx, donorPet.OwnerID, user.UserPreloadOptions{
		WithIdentities: true,
	})
	if err != nil {
		return err
	}

	recipientUser, err := h.userRepo.GetByID(ctx, recipientPet.OwnerID, user.UserPreloadOptions{
		WithIdentities: true,
	})
	if err != nil {
		return err
	}

	donorProviderMaxID, _ := extractProviderIDs(donorUser)

	switch application.Status {
	case donormodel.DonorResponseStatusRejected:
		event := donorevent.DonorReject{
			RecipientPetName:    recipientPet.Name,
			RecipientBloodGroup: recipientPet.BloodGroupName,
			DonorProviderMaxID:  donorProviderMaxID,
			DonorPetName:        donorPet.Name,
			RejectedReason:      rejectedReason,
			CreatedAt:           time.Now(),
		}

		if err := h.publisher.PublishDonorReject(ctx, event); err != nil {
			return err
		}
	case donormodel.DonorResponseStatusAccepted:
		recipientProviderMaxID, _ := extractProviderIDs(recipientUser)

		notConfirmedEvent := donorevent.DonorNotConfirmed{
			DonorPetName:        donorPet.Name,
			DonorBloodGroup:     donorPet.BloodGroupName,
			DonorProviderMaxID:  donorProviderMaxID,
			RecipientPetName:    recipientPet.Name,
			RecipientBloodGroup: recipientPet.BloodGroupName,
			RecipientUserData: donorevent.ContactData{
				Name:             recipientUser.FullName,
				ProviderMaxID:    recipientProviderMaxID,
				ProviderTelegram: "", // if needed
				Phone:            recipientUser.Phone,
			},
			DonorUserData: donorevent.ContactData{
				Name:             donorUser.FullName,
				ProviderMaxID:    donorProviderMaxID,
				ProviderTelegram: "", // if needed
				Phone:            donorUser.Phone,
			},
			CreatedAt: time.Now(),
		}

		if err := h.publisher.PublishDonorNotConfirmed(ctx, notConfirmedEvent); err != nil {
			return err
		}
	default:
		// Unexpected status, do nothing or log
	}

	return nil
}
