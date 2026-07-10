package pet

import (
	"time"

	bloodreqmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/common"
	donormodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet/model"
)

type PetService interface {
	CalculateStatus(pet *model.Pet, application *donormodel.DonorResponse, bloodReq *bloodreqmodel.BloodRequestWithApplications) model.PetStatus
	RecalculateFactorsAndStatus(pet *model.Pet, now time.Time, application *donormodel.DonorResponse, bloodReq *bloodreqmodel.BloodRequestWithApplications)
	CalculateRecoveryDays(pet *model.Pet, recoveryPeriodMonths int, now time.Time) *int
}

// PetService provides business logic for pets
type PetServiceV2 struct{}

// NewPetService creates a new PetService
func NewPetServiceV2() *PetServiceV2 {
	return &PetServiceV2{}
}

// CalculateStatus calculates and returns the pet's status based on related aggregates.
func (s *PetServiceV2) CalculateStatus(pet *model.Pet, application *donormodel.DonorResponse, bloodReq *bloodreqmodel.BloodRequestWithApplications) model.PetStatus {
	// Planned donation takes precedence over everything else.
	if application != nil && application.IsActiveForDonation() {
		return model.PetStatusPlannedDonation
	}

	// Active blood request overrides donor/recovering status.
	if !bloodReq.IsClosed() {
		if bloodReq.HasActiveDonorApplications() {
			return model.PetStatusBloodFound
		}
		return model.PetStatusRecipient
	}

	// Recovering pet that isn't in an active request.
	if pet.IsRecovering() {
		return model.PetStatusRecovering
	}

	// No stop factors → eligible donor.
	if !pet.HasStopFactors() {
		return model.PetStatusDonor
	}

	return model.PetStatusNone
}

// RecalculateFactorsAndStatus recalculates pet's factors and sets status based on related aggregates
func (s *PetServiceV2) RecalculateFactorsAndStatus(pet *model.Pet, now time.Time, application *donormodel.DonorResponse, bloodReq *bloodreqmodel.BloodRequestWithApplications) {
	isRecipient := !bloodReq.IsClosed()
	isPlaningDonation := application != nil && application.IsActiveForDonation()

	pet.RecalculateFactors(now, isRecipient, isPlaningDonation)
	status := s.CalculateStatus(pet, application, bloodReq)
	pet.PetStatus = status

	if !bloodReq.IsClosed() && bloodReq.BloodRequest.PrioritySearch && pet.Privilege == "" {
		pet.Privilege = common.PrivilegePrioritySearch
	}
}

// CalculateRecoveryDays calculates remaining recovery days after donation (date-only comparison)
func (s *PetServiceV2) CalculateRecoveryDays(pet *model.Pet, recoveryPeriodMonths int, now time.Time) *int {
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
