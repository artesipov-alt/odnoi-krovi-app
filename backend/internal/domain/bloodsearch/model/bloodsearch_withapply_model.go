package model

import donormodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/model"

type BloodRequestWithApplications struct {
	BloodRequest      BloodRequest
	DonorApplications []donormodel.DonorResponse
}

func (b *BloodRequestWithApplications) HasActiveDonorApplications() bool {
	if b == nil {
		return false
	}
	for _, app := range b.DonorApplications {
		if app.IsActiveForDonation() {
			return true
		}
	}
	return false
}

// IsClosed checks if the blood request is closed. Nil-safe.
func (b *BloodRequestWithApplications) IsClosed() bool {
	if b == nil {
		return false
	}
	return b.BloodRequest.IsClosed()
}

// IsActive checks if the blood request is active. Nil-safe.
func (b *BloodRequestWithApplications) IsActive() bool {
	if b == nil {
		return false
	}
	return b.BloodRequest.IsActive()
}

// RecalculateStatus recalculates the blood request status. Nil-safe.
func (b *BloodRequestWithApplications) RecalculateStatus() {
	if b == nil {
		return
	}
	b.BloodRequest.RecalculateStatus()
}

func (r *BloodRequestWithApplications) SearchingBloodGroupNames() []string {
	var searchingBloodGroupNames []string
	searchingBloodGroupNames = append(searchingBloodGroupNames, r.BloodRequest.BloodGroupNames...)
	if r.BloodRequest.IncludeUnknownBloodGroup {
		searchingBloodGroupNames = append(searchingBloodGroupNames, "UNKNOWN")
	}
	return searchingBloodGroupNames
}

func (r *BloodRequestWithApplications) SearchingRegions() []string {
	var searchingRegions []string
	searchingRegions = append(searchingRegions, r.BloodRequest.Regions...)
	return searchingRegions
}
