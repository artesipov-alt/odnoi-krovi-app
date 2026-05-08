package cmd

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"
	bloodreqmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bonus"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor"
	donormodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/ports"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/presistance"
)

type CloseRequestHandler struct {
	bloodRepo bloodsearch.BloodRequestRepository
	donorRepo donor.Repository
	txManager *presistance.TxManager
	publisher ports.EventPublisher
	petRepo   pet.PetReadRepository
	bonusSvc  *bonus.BonusService
}

func NewCloseRequestHandler(
	bloodRepo bloodsearch.BloodRequestRepository,
	donorRepo donor.Repository,
	txManager *presistance.TxManager,
	publisher ports.EventPublisher,
	petRepo pet.PetReadRepository,
	bonusSvc *bonus.BonusService,
) *CloseRequestHandler {
	return &CloseRequestHandler{
		bloodRepo: bloodRepo,
		donorRepo: donorRepo,
		txManager: txManager,
		publisher: publisher,
		petRepo:   petRepo,
		bonusSvc:  bonusSvc,
	}
}

func (h *CloseRequestHandler) Handle(ctx context.Context, bloodReqID string) error {
	bloodReq, err := h.bloodRepo.GetByID(ctx, bloodReqID)
	if err != nil {
		return err
	}

	err = h.txManager.WithTx(ctx, func(txCtx context.Context) error {
		if err := h.bloodRepo.UpdateStatus(txCtx, bloodReq.ID, bloodreqmodel.BloodRequestStatusClosed); err != nil {
			return err
		}
		for _, application := range bloodReq.DonorApplications {
			if application.Status == donormodel.DonorResponseStatusPending || application.Status == donormodel.DonorResponseStatusAccepted || (application.Status == donormodel.DonorResponseStatusCompleted && application.IsConfirmed != true) {
				if err := application.Reject("other"); err != nil {
					return err
				}
				if err := h.donorRepo.Reject(txCtx, &application); err != nil {
					return err
				}
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

	return nil
}
