package cmd

import (
	"context"
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"
	bloodsearchevent "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/events"
	bloodmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor"
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
}

func NewConfirmDonationHandler(
	bloodRepo bloodsearch.BloodRequestRepository,
	donorRepo donor.Repository,
	petRepo pet.Repository,
	userRepo user.Repository,
	txManager *presistance.TxManager,
	publisher ports.EventPublisher,
) *ConfirmDonationHandler {
	return &ConfirmDonationHandler{
		bloodRepo: bloodRepo,
		donorRepo: donorRepo,
		petRepo:   petRepo,
		userRepo:  userRepo,
		txManager: txManager,
		publisher: publisher,
	}
}

func (h *ConfirmDonationHandler) Handle(ctx context.Context, donorResponseID string, factAmount float64) error {
	var bloodReq *bloodmodel.BloodRequestWithApplications
	var application *donormodel.DonorResponse

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
		bloodReq, err = h.bloodRepo.GetByApplicationID(txCtx, donorResponseID)
		if err != nil {
			return err
		}

		bloodReq.RecalculateBloodAmount()
		bloodReq.RecalculateStatus()

		if err := h.bloodRepo.UpdateStatus(txCtx, bloodReq.ID, bloodReq.Status); err != nil {
			return err
		}

		if bloodReq.Status == bloodmodel.BloodRequestStatusClosed {
			for _, app := range bloodReq.DonorApplications {
				if app.ID != donorResponseID && (app.Status == donormodel.DonorResponseStatusPending || app.Status == donormodel.DonorResponseStatusAccepted || (app.Status == donormodel.DonorResponseStatusCompleted && app.IsConfirmed != true)) {
					if err := app.Reject("other"); err != nil {
						return err
					}
					if err := h.donorRepo.Reject(txCtx, &app); err != nil {
						return err
					}
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

	// Extract Provider IDs
	donorMaxID, donorTelegramID := extractProviderIDs(donorUser)

	donorBloodGroup := ""
	if donorPet.BloodGroupName != nil {
		donorBloodGroup = *donorPet.BloodGroupName
	}

	event := bloodsearchevent.DonationConfirmed{
		DonorData: bloodsearchevent.DonorInfo{
			UserName:         donorUser.FullName,
			PetName:          donorPet.Name,
			ProviderMaxID:    donorMaxID,
			ProviderTelegram: donorTelegramID,
			Phone:            donorUser.Phone,
			BloodGroup:       donorBloodGroup,
		},
		Volume:    factAmount,
		CreatedAt: time.Now(),
	}

	if err := h.publisher.PublishDonationConfirmed(ctx, event); err != nil {
		return err
	}

	return nil
}
