package cmd

import (
	"context"
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor"
	donorevent "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/events"
	donormodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/ports"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user"
)

type CompleteDonationHandler struct {
	donorRepo donor.Repository
	bloodRepo bloodsearch.BloodRequestRepository
	petRepo   pet.Repository
	userRepo  user.Repository
	publisher ports.EventPublisher
}

func NewCompleteDonationHandler(
	donorRepo donor.Repository,
	bloodRepo bloodsearch.BloodRequestRepository,
	petRepo pet.Repository,
	userRepo user.Repository,
	publisher ports.EventPublisher,
) *CompleteDonationHandler {
	return &CompleteDonationHandler{
		donorRepo: donorRepo,
		bloodRepo: bloodRepo,
		petRepo:   petRepo,
		userRepo:  userRepo,
		publisher: publisher,
	}
}

func (h *CompleteDonationHandler) Handle(ctx context.Context, resID string, amount float64) error {
	// Получаем DonorResponse
	donorResponse, err := h.donorRepo.GetDonorResponseByID(ctx, resID)
	if err != nil {
		return apperrors.Internal(err, "failed to get donor response")
	}
	if donorResponse.Status != donormodel.DonorResponseStatusAccepted {
		return apperrors.BadRequest("donor response status is invalid").WithMessage("donor response must be accepted to complete donation")
	}
	if err := h.donorRepo.Complete(ctx, donorResponse.ID, amount); err != nil {
		return err
	}

	// Publish DonorCompleted event
	bloodReq, err := h.bloodRepo.GetByApplicationID(ctx, resID)
	if err != nil {
		return apperrors.Internal(err, "failed to get blood request")
	}

	donorPet, err := h.petRepo.GetByID(ctx, donorResponse.DonorID, pet.PetPreloadOptions{})
	if err != nil {
		return apperrors.Internal(err, "failed to get donor pet")
	}

	recipientPet, err := h.petRepo.GetByID(ctx, bloodReq.PetID, pet.PetPreloadOptions{})
	if err != nil {
		return apperrors.Internal(err, "failed to get recipient pet")
	}

	recipientUser, err := h.userRepo.GetByID(ctx, recipientPet.OwnerID, user.UserPreloadOptions{
		WithIdentities: true,
	})
	if err != nil {
		return apperrors.Internal(err, "failed to get recipient user")
	}

	recipientProviderMaxID, _ := extractProviderIDs(recipientUser)

	event := donorevent.DonorCompleted{
		DonorPetName:           donorPet.Name,
		DonorBloodGroup:        donorPet.BloodGroupName,
		RecipientProviderMaxID: recipientProviderMaxID,
		RecipientPetName:       recipientPet.Name,
		Amount:                 amount,
		CreatedAt:              time.Now(),
	}

	if err := h.publisher.PublishDonorCompleted(ctx, event); err != nil {
		return apperrors.Internal(err, "failed to publish donor completed event")
	}

	return nil
}
