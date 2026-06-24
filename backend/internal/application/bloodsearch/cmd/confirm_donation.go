package cmd

import (
	"context"
	"log/slog"
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"
	bloodsearchevent "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/events"
	bloodmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bonus"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor"
	donorevent "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/events"
	donormodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/ports"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/presistance"
)

type ConfirmDonationHandler struct {
	bloodRepo bloodsearch.BloodRequestRepository
	donorRepo donor.Repository
	petRepo   pet.Repository
	userRepo  user.Repository
	txManager *presistance.TxManager
	publisher ports.EventPublisher
	bonusSvc  *bonus.BonusService
}

func NewConfirmDonationHandler(
	bloodRepo bloodsearch.BloodRequestRepository,
	donorRepo donor.Repository,
	petRepo pet.Repository,
	userRepo user.Repository,
	txManager *presistance.TxManager,
	publisher ports.EventPublisher,
	bonusSvc *bonus.BonusService,
) *ConfirmDonationHandler {
	return &ConfirmDonationHandler{
		bloodRepo: bloodRepo,
		donorRepo: donorRepo,
		petRepo:   petRepo,
		userRepo:  userRepo,
		txManager: txManager,
		publisher: publisher,
		bonusSvc:  bonusSvc,
	}
}

func (h *ConfirmDonationHandler) Handle(ctx context.Context, donorResponseID string, factAmount float64) error {
	var bloodReq *bloodmodel.BloodRequestWithApplications
	var application *donormodel.DonorResponse
	var rejectedDonorIDs []string

	application, err := h.donorRepo.GetDonorResponseByID(ctx, donorResponseID)
	if err != nil {
		return err
	}

	err = h.txManager.WithTx(ctx, func(txCtx context.Context) error {
		if err := application.Confirm(factAmount); err != nil {
			return err
		}
		if err := h.donorRepo.Confirm(txCtx, donorResponseID, factAmount); err != nil {
			return err
		}
		var err error
		bloodReq, err = h.bloodRepo.GetByApplicationID(txCtx, donorResponseID, false)
		if err != nil {
			return err
		}

		bloodReq.RecalculateBloodAmount()
		bloodReq.RecalculateStatus()

		if err := h.bloodRepo.UpdateStatus(txCtx, bloodReq.ID, bloodReq.Status); err != nil {
			return err
		}

		if bloodReq.IsClosed() {
			for _, app := range bloodReq.DonorApplications {
				if app.ID != donorResponseID && app.IsActiveForDonation() {
					if err := app.Reject("other"); err != nil {
						return err
					}
					if err := h.donorRepo.Reject(txCtx, &app); err != nil {
						return err
					}
					rejectedDonorIDs = append(rejectedDonorIDs, app.DonorID)
				}
			}
			// Установить флаг переливания для recipient'а
			if err := h.petRepo.SetTransfused(txCtx, bloodReq.PetID, true); err != nil {
				return err
			}
		}

		application, err = h.donorRepo.GetDonorResponseByID(txCtx, donorResponseID)
		if err != nil {
			return err
		}

		now := time.Now()
		if err := h.petRepo.SetLastDonation(txCtx, application.DonorID, &now); err != nil {
			return err
		}

		// Confirm reserved bonuses for the donor
		donorPet, err := h.petRepo.GetByID(txCtx, application.DonorID, pet.PetPreloadOptions{})
		if err != nil {
			return err
		}
		if err := h.bonusSvc.ConfirmBonuses(txCtx, donorPet.OwnerID, donorPet.Type); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	// Get donor data
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

	// Get recipient data
	bloodReq, err = h.bloodRepo.GetByApplicationID(ctx, donorResponseID, false)
	if err != nil {
		return err
	}
	recipientPet, err := h.petRepo.GetByID(ctx, bloodReq.PetID, pet.PetPreloadOptions{})
	if err != nil {
		return err
	}

	// Extract Provider IDs
	donorMaxID, donorTelegramID := extractProviderIDs(donorUser)

	event := bloodsearchevent.DonationConfirmed{
		DonorData: bloodsearchevent.DonorInfo{
			UserName:         donorUser.FullName,
			PetName:          donorPet.Name,
			ProviderMaxID:    donorMaxID,
			ProviderTelegram: donorTelegramID,
			Phone:            donorUser.Phone,
			BloodGroup:       donorPet.BloodGroupName,
		},
		RecipientData: bloodsearchevent.RecipientInfo{
			PetName:    recipientPet.Name,
			BloodGroup: recipientPet.BloodGroupName,
		},
		Volume:    factAmount,
		CreatedAt: time.Now(),
	}

	if err := h.publisher.PublishDonationConfirmed(ctx, event); err != nil {
		slog.Error("failed to publish donation confirmed notification", "err", err, "donorResponseID", donorResponseID)
	}

	// Publish rejection events for donors who were auto-rejected when the request closed
	if len(rejectedDonorIDs) > 0 {
		for _, rejectedDonorID := range rejectedDonorIDs {
			rejectedDonorPet, err := h.petRepo.GetByID(ctx, rejectedDonorID, pet.PetPreloadOptions{})
			if err != nil {
				return err
			}
			rejectedDonorUser, err := h.userRepo.GetByID(ctx, rejectedDonorPet.OwnerID, user.UserPreloadOptions{
				WithIdentities: true,
			})
			if err != nil {
				return err
			}

			donorMaxID, donorTelegramID := extractProviderIDs(rejectedDonorUser)

			rejectEvent := donorevent.DonorReject{
				RecipientPetName:        recipientPet.Name,
				RecipientBloodGroup:     recipientPet.BloodGroupName,
				DonorProviderMaxID:      donorMaxID,
				DonorProviderTelegramID: donorTelegramID,
				DonorPetName:            rejectedDonorPet.Name,
				RejectedReason:          "other",
				CreatedAt:               time.Now(),
			}

			if err := h.publisher.PublishDonorReject(ctx, rejectEvent); err != nil {
				slog.Error("failed to publish donor reject notification", "err", err, "donorID", rejectedDonorID)
			}
		}
	}

	return nil
}
