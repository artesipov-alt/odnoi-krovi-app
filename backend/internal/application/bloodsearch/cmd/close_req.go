package cmd

import (
	"context"
	"log/slog"
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"
	bloodreqmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bonus"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor"
	donorevent "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/events"
	donormodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/ports"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/presistance"
)

type CloseRequestHandler struct {
	bloodRepo bloodsearch.Repository
	donorRepo donor.Repository
	petRepo   pet.PetReadRepository
	userRepo  user.Repository
	txManager *presistance.TxManager
	publisher ports.EventPublisher
	bonusSvc  *bonus.BonusService
}

func NewCloseRequestHandler(
	bloodRepo bloodsearch.Repository,
	donorRepo donor.Repository,
	petRepo pet.PetReadRepository,
	userRepo user.Repository,
	txManager *presistance.TxManager,
	publisher ports.EventPublisher,
	bonusSvc *bonus.BonusService,
) *CloseRequestHandler {
	return &CloseRequestHandler{
		bloodRepo: bloodRepo,
		donorRepo: donorRepo,
		petRepo:   petRepo,
		userRepo:  userRepo,
		txManager: txManager,
		publisher: publisher,
		bonusSvc:  bonusSvc,
	}
}

func (h *CloseRequestHandler) Handle(ctx context.Context, bloodReqID string) error {
	bloodReq, err := h.bloodRepo.GetByID(ctx, bloodReqID)
	if err != nil {
		return err
	}

	var rejectedDonorIDs []string

	err = h.txManager.WithTx(ctx, func(txCtx context.Context) error {
		if err := h.bloodRepo.UpdateStatus(txCtx, bloodReq.BloodRequest.ID, bloodreqmodel.BloodRequestStatusClosed); err != nil {
			return err
		}
		for i := range bloodReq.DonorApplications {
			application := &bloodReq.DonorApplications[i]
			if application.IsActiveForDonation() {
				if err := application.Reject("other"); err != nil {
					return err
				}
				if err := h.donorRepo.Update(txCtx, application); err != nil {
					return err
				}
				rejectedDonorIDs = append(rejectedDonorIDs, application.DonorID)
				// Unassign reserved bonuses if the application was not completed
				if application.Status != donormodel.DonorResponseStatusCompleted {
					donorPet, err := h.petRepo.GetByID(txCtx, application.DonorID, pet.PetPreloadOptions{})
					if err != nil {
						return err
					}
					if err := h.bonusSvc.UnassignReservedBonuses(txCtx, donorPet.OwnerID, donorPet.Type); err != nil {
						return err
					}
				}
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	// Publish rejection events for donors who were auto-rejected when the request closed
	if len(rejectedDonorIDs) > 0 {
		recipientPet, err := h.petRepo.GetByID(ctx, bloodReq.BloodRequest.PetID, pet.PetPreloadOptions{})
		if err != nil {
			slog.Error("failed to get recipient pet for close request notifications", "err", err, "bloodReqID", bloodReqID)
			return nil
		}

		for _, rejectedDonorID := range rejectedDonorIDs {
			rejectedDonorPet, err := h.petRepo.GetByID(ctx, rejectedDonorID, pet.PetPreloadOptions{})
			if err != nil {
				slog.Error("failed to get rejected donor pet", "err", err, "donorID", rejectedDonorID)
				continue
			}
			rejectedDonorUser, err := h.userRepo.GetByID(ctx, rejectedDonorPet.OwnerID, user.UserPreloadOptions{
				WithIdentities: true,
			})
			if err != nil {
				slog.Error("failed to get rejected donor user", "err", err, "donorID", rejectedDonorID)
				continue
			}

			donorMaxID, donorTelegramID := rejectedDonorUser.MessengerContacts()

			rejectEvent := donorevent.DonorReject{
				RecipientPetName:        recipientPet.Name,
				RecipientBloodGroup:     recipientPet.BloodGroupName,
				DonorProviderMaxID:      donorMaxID,
				DonorProviderTelegramID: donorTelegramID,
				DonorPetName:            rejectedDonorPet.Name,
				RejectedReason:          "other",
				CreatedAt:               time.Now(),
			}

			if err := h.publisher.PublishEvent(ctx, ports.EventDonorReject, rejectEvent); err != nil {
				slog.Error("failed to publish donor reject notification", "err", err, "donorID", rejectedDonorID)
			}
		}
	}

	return nil
}
