package pet

import (
	"slices"
	"time"

	bloodreqmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/common"
	donormodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet/model"
)

// PetService provides business logic for pets
type PetService struct{}

// NewPetService creates a new PetService
func NewPetService() *PetService {
	return &PetService{}
}

// CalculateAndSetStatus calculates and sets the pet's status based on related aggregates
func (s *PetService) CalculateAndSetStatus(pet *model.Pet, application *donormodel.DonorResponse, bloodReq *bloodreqmodel.BloodRequestWithApplications) {
	if s.hasActiveBloodRequest(bloodReq) {
		if s.hasActiveDonorApplications(bloodReq) {
			pet.PetStatus = model.PetStatusBloodFound
		} else {
			pet.PetStatus = model.PetStatusRecipient
		}
	} else if s.canBeDonor(pet) {
		pet.PetStatus = model.PetStatusDonor
	}

	if application.IsActiveForDonation() {
		pet.PetStatus = model.PetStatusPlannedDonation
	}

	if s.shouldBeRecovering(pet) && pet.PetStatus != model.PetStatusRecipient && pet.PetStatus != model.PetStatusBloodFound {
		pet.PetStatus = model.PetStatusRecovering
	}
}

// hasActiveBloodRequest checks if there is an active blood request
func (s *PetService) hasActiveBloodRequest(bloodReq *bloodreqmodel.BloodRequestWithApplications) bool {
	return bloodReq != nil && bloodReq.Status != bloodreqmodel.BloodRequestStatusClosed
}

// hasActiveDonorApplications checks if there are active donor applications
func (s *PetService) hasActiveDonorApplications(bloodReq *bloodreqmodel.BloodRequestWithApplications) bool {
	if bloodReq == nil {
		return false
	}
	for _, app := range bloodReq.DonorApplications {
		if app.IsActiveForDonation() {
			return true
		}
	}
	return false
}

// canBeDonor checks if the pet can be a donor (no stop factors)
func (s *PetService) canBeDonor(pet *model.Pet) bool {
	return len(pet.StopFactors) == 0
}

// shouldBeRecovering checks if the pet should be in recovering status
func (s *PetService) shouldBeRecovering(pet *model.Pet) bool {
	if pet.Health != nil && pet.Health.Transfused != nil && *pet.Health.Transfused {
		return false
	}
	return slices.Contains(pet.StopFactors, string(model.StopFactorDonationTooRecent))
}

// RecalculateFactorsAndStatus recalculates pet's factors and sets status based on related aggregates
func (s *PetService) RecalculateFactorsAndStatus(pet *model.Pet, now time.Time, application *donormodel.DonorResponse, bloodReq *bloodreqmodel.BloodRequestWithApplications) {
	isRecipient := s.hasActiveBloodRequest(bloodReq)
	pet.RecalculateFactors(now, isRecipient, application.IsActiveForDonation())
	s.CalculateAndSetStatus(pet, application, bloodReq)

	if bloodReq != nil && !bloodReq.IsClosed() && bloodReq.PrioritySearch && pet.Privilege == "" {
		pet.Privilege = common.PrivilegePrioritySearch
	}
}

// CalculateRecoveryDays calculates remaining recovery days after donation (date-only comparison)
func (s *PetService) CalculateRecoveryDays(pet *model.Pet, recoveryPeriodMonths int, now time.Time) *int {
	if pet.Health == nil || pet.Health.LastDonation == nil || recoveryPeriodMonths <= 0 {
		return nil
	}
	// Truncate to date only (ignore time)
	nowDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	lastDonationDate := time.Date(pet.Health.LastDonation.Year(), pet.Health.LastDonation.Month(), pet.Health.LastDonation.Day(), 0, 0, 0, 0, pet.Health.LastDonation.Location())
	endDate := lastDonationDate.AddDate(0, recoveryPeriodMonths, 0)
	if endDate.After(nowDate) {
		days := int(endDate.Sub(nowDate).Hours() / 24)
		return &days
	}
	return nil
}
