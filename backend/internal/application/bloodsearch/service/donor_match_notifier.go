package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/events"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor"
	donormodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet"
	petmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/ports"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user"
)

type DonorMatchNotifier struct {
	petRepo         pet.Repository
	userRepo        user.Repository
	bloodSearchRepo bloodsearch.Repository
	donorRepo       donor.Repository
	publisher       ports.EventPublisher
}

func NewDonorMatchNotifier(
	petRepo pet.Repository,
	userRepo user.Repository,
	bloodSearchRepo bloodsearch.Repository,
	donorRepo donor.Repository,
	publisher ports.EventPublisher,
) *DonorMatchNotifier {
	return &DonorMatchNotifier{
		petRepo:         petRepo,
		userRepo:        userRepo,
		bloodSearchRepo: bloodSearchRepo,
		donorRepo:       donorRepo,
		publisher:       publisher,
	}
}

func (n *DonorMatchNotifier) NotifyMatchDonors(ctx context.Context, req *model.BloodRequestWithApplications) error {
	bloodGroupNames := req.SearchingBloodGroupNames()
	regions := req.BloodRequest.Regions

	initiatorPet, err := n.petRepo.GetByID(ctx, req.BloodRequest.PetID, pet.PetPreloadOptions{})
	if err != nil {
		return apperrors.Internal(err, "failed to get initiator pet")
	}
	initiatorUserID := initiatorPet.OwnerID

	petsIDs, err := n.petRepo.GetPetIDsByBloodGroupAndRegion(ctx, bloodGroupNames, regions)
	if err != nil {
		return apperrors.Internal(err, "failed to get pet IDs")
	}

	// Batch fetch applications and blood requests
	applicationsMap, err := n.donorRepo.GetByPetIDs(ctx, petsIDs, false)
	if err != nil {
		return apperrors.Internal(err, "failed to get donor applications")
	}

	bloodReqsMap, err := n.bloodSearchRepo.GetByPetIDs(ctx, petsIDs)
	if err != nil {
		return apperrors.Internal(err, "failed to get blood requests")
	}

	pets, err := n.petRepo.GetByIDs(ctx, petsIDs, pet.PetPreloadOptions{
		WithHealth:     true,
		WithAnalyses:   true,
		WithTreatments: true,
	})
	if err != nil {
		return apperrors.Internal(err, "failed to get pets")
	}

	for _, matchingPet := range pets {
		applications := applicationsMap[matchingPet.ID]
		var donorApplication *donormodel.DonorResponse
		for _, app := range applications {
			if app.IsActiveForDonation() {
				donorApplication = app
				break
			}
		}
		donorBloodReq := bloodReqsMap[matchingPet.ID]

		matchingPet.RecalculateStatus(time.Now(), pet.BuildDonationContext(donorApplication, donorBloodReq))
	}

	var avilableDonors []petmodel.Pet
	for _, pet := range pets {
		if !pet.IsDonor() {
			continue
		}
		if !req.BloodRequest.IsCoversNededAmount(pet.CalculateDonationAmount()) {
			continue
		}
		avilableDonors = append(avilableDonors, *pet)
	}

	// Дедуп по OwnerID ДО похода в userRepo — чтобы не запрашивать
	// одного и того же владельца несколько раз, если у него несколько
	// подходящих питомцев.
	ownerSeen := make(map[string]struct{})
	for _, donorPet := range avilableDonors {
		if donorPet.OwnerID == initiatorUserID {
			continue
		}
		if _, exists := ownerSeen[donorPet.OwnerID]; exists {
			continue
		}
		ownerSeen[donorPet.OwnerID] = struct{}{}

		donorUser, err := n.userRepo.GetByID(ctx, donorPet.OwnerID, user.UserPreloadOptions{
			WithIdentities: true,
		})
		if err != nil {
			slog.Error("failed to get donor user", "err", err, "petID", donorPet.ID)
			continue
		}

		maxID, telegramID := donorUser.MessengerContacts()
		if maxID == "" && telegramID == "" {
			continue
		}

		if err := n.publisher.PublishEvent(ctx, ports.EventBloodRequestCreated, events.BloodRequestCreated{
			BloodTypes: bloodGroupNames,
			Regions:    regions,
			TelegramID: telegramID,
			MaxID:      maxID,
			CreatedAt:  time.Now(),
		}); err != nil {
			slog.Error("failed to publish blood request created", "err", err, "petID", donorPet.ID)
		}
	}

	return nil
}
