package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/events"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bonus"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor"
	donormodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet"
	petmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/ports"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/presistance"
)

type CreateRequestHandler struct {
	bloodRepo  bloodsearch.BloodRequestRepository
	petRepo    pet.Repository
	donorRepo  donor.Repository
	userRepo   user.Repository
	bonusRepo  bonus.Repository
	publisher  ports.EventPublisher
	petService *pet.PetService
	txManager  *presistance.TxManager
}

func NewCreateRequestHandler(
	bloodRepo bloodsearch.BloodRequestRepository,
	petRepo pet.Repository,
	donorRepo donor.Repository,
	userRepo user.Repository,
	bonusRepo bonus.Repository,
	publisher ports.EventPublisher,
	petService *pet.PetService,
	txManager *presistance.TxManager,
) *CreateRequestHandler {
	return &CreateRequestHandler{
		bloodRepo:  bloodRepo,
		petRepo:    petRepo,
		donorRepo:  donorRepo,
		userRepo:   userRepo,
		bonusRepo:  bonusRepo,
		publisher:  publisher,
		petService: petService,
		txManager:  txManager,
	}
}

func (h *CreateRequestHandler) Handle(ctx context.Context, req *model.BloodRequest) (*model.BloodRequestWithApplications, error) {
	// Проверяем существование питомца
	petRecipient, err := h.petRepo.GetByID(ctx, req.PetID, pet.PetPreloadOptions{})
	if err != nil {
		return nil, apperrors.Internal(err, "Ошибка поиска питомца")
	}

	// Проверяем, нет ли уже активной заявки для этого питомца
	activeExists, err := h.bloodRepo.ExistsByPetID(ctx, req.PetID)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to check request existence")
	}
	if activeExists {
		return nil, apperrors.ErrBloodRequestAlreadyExists
	}

	var newReq *model.BloodRequestWithApplications
	err = h.txManager.WithTx(ctx, func(txCtx context.Context) error {
		var err error
		if req.PrioritySearch && petRecipient.Privilege == "" {
			err = h.bonusRepo.SubtractPrioritySearch(txCtx, req.OwnerID)
			if err != nil {
				return apperrors.Internal(err, "failed to subtract priority search")
			}
		}
		newReq, err = h.bloodRepo.Create(txCtx, req)
		if err != nil {
			return apperrors.Internal(err, "failed to create blood request")
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	pets, err := h.petRepo.GetPetsByBloodGroupAndRegion(ctx, petRecipient.Type, newReq.SearchingBloodGroupNames(), req.Regions)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to get pets")
	}
	slog.Info("found pets by blood group and region", "count", len(pets), "bloodGroups", req.BloodGroupNames, "regions", req.Regions)

	// Collect pet IDs for batch queries
	petIDs := make([]string, len(pets))
	for i, pet := range pets {
		petIDs[i] = pet.ID
	}

	// Batch fetch applications and blood requests
	applicationsMap, err := h.donorRepo.GetByPetIDs(ctx, petIDs, false)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to get donor applications")
	}

	bloodReqsMap, err := h.bloodRepo.GetByPetIDs(ctx, petIDs)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to get blood requests")
	}

	timeNow := time.Now()
	for _, pet := range pets {
		applications := applicationsMap[pet.ID]
		var donorApplication *donormodel.DonorResponse
		for _, app := range applications {
			if app.IsActiveForDonation() {
				donorApplication = app
				break
			}
		}
		donorBloodReq := bloodReqsMap[pet.ID]
		h.petService.RecalculateFactorsAndStatus(pet, timeNow, donorApplication, donorBloodReq)
	}
	// Логика события
	// TODO: Вынести отдельно.

	var avilableDonors []petmodel.Pet
	for _, pet := range pets {
		if pet.PetStatus == petmodel.PetStatusDonor {
			avilableDonors = append(avilableDonors, *pet)
		}
	}

	// Get peers for available donors
	peersMap := make(map[string]events.Peers)
	for _, donorPet := range avilableDonors {
		donorUser, err := h.userRepo.GetByID(ctx, donorPet.OwnerID, user.UserPreloadOptions{
			WithIdentities: true,
		})
		if err != nil {
			slog.Error("failed to get donor user", "err", err, "petID", donorPet.ID)
			continue
		}
		maxID, telegramID := extractProviderIDs(donorUser)
		if maxID != "" || telegramID != "" {
			key := fmt.Sprintf("%s|%s", maxID, telegramID)
			if _, exists := peersMap[key]; !exists {
				peersMap[key] = events.Peers{
					MaxID:      maxID,
					TelegramID: telegramID,
				}
			}
		}
	}
	peers := make([]events.Peers, 0, len(peersMap))
	for _, p := range peersMap {
		peers = append(peers, p)
	}

	if err := h.publisher.PublishBloodRequestCreated(ctx, events.BloodRequestCreated{
		RequestID:      newReq.ID,
		BloodTypes:     req.BloodGroupNames,
		Regions:        req.Regions,
		AvilableDonors: peers,
		CreatedAt:      *newReq.CreatedAt,
	}); err != nil {
		slog.Error("failed to publish blood request created event", "err", err)
		return newReq, nil
	}

	return newReq, nil
}
