package model

import "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/common"

// Recipient представляет модель чтения реципиент
type BloodRequestWithMatchingDonors struct {
	BloodRequest      BloodRequest
	RecipientData     RecipientData
	MatchingDonors    []MatchingDonorReadModel
	DefaultDonorPrefs *DefaultDonorPrefs
}

type RecipientData struct {
	PetName        string
	PetType        common.PetType
	BloodGroupName string
	OwnerName      string
	OwnerID        string
	Privilege      common.Privilege
	PhotoURLs      []string
}

type DefaultDonorPrefs struct {
	CompensationType common.CompensationType
	Bonuses          []string
	TaxiCompensation bool
}

// MatchingDonorReadModel представляет модель чтения для подходящего донора
type MatchingDonorReadModel struct {
	PetID           string
	PetName         string
	Amount          float64
	DonorBloodGroup string
	PhotoURLs       []string
}

func (r *BloodRequestWithMatchingDonors) SetDefaultPrefs(compensationType common.CompensationType, taxiCompensation bool) {
	r.DefaultDonorPrefs = &DefaultDonorPrefs{
		CompensationType: compensationType,
		TaxiCompensation: taxiCompensation,
		Bonuses:          []string{},
	}
}

// SyncPrivilegeAndPriority synchronizes privilege and priority search based on business rules
func (r *BloodRequestWithMatchingDonors) SyncPrivilegeAndPriority() {
	if r.RecipientData.Privilege != "" {
		r.BloodRequest.PrioritySearch = true
	} else if r.BloodRequest.PrioritySearch {
		r.RecipientData.Privilege = common.PrivilegePrioritySearch
	}
}

func (r *BloodRequestWithApplications) SearchingBloodGroupNames() []string {
	var searchingBloodGroupNames []string
	searchingBloodGroupNames = append(searchingBloodGroupNames, r.BloodRequest.BloodGroupNames...)
	if r.BloodRequest.IncludeUnknownBloodGroup {
		searchingBloodGroupNames = append(searchingBloodGroupNames, "UNKNOWN")
	}
	return searchingBloodGroupNames
}
