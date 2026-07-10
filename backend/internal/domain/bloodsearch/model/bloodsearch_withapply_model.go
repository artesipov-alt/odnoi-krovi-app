package model

import (
	"math"

	donormodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/model"
)

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

func (b *BloodRequestWithApplications) RecalculateBloodAmount() {
	var donated float64
	for _, app := range b.DonorApplications {
		if app.IsConfirmed && app.Status == donormodel.DonorResponseStatusCompleted {
			donated += app.Amount
		}
	}
	b.BloodRequest.BloodVolumeDonated = math.Round(donated*10) / 10

	reserved := b.BloodRequest.BloodVolumeDonated
	for _, app := range b.DonorApplications {
		// Ищем только откликнувшихся доноров
		if (app.IsConfirmed == false && app.Status == donormodel.DonorResponseStatusAccepted) || (app.IsConfirmed == false && app.Status == donormodel.DonorResponseStatusCompleted) {
			reserved += app.Amount
			// Обрезаем до максимального
			if reserved >= b.BloodRequest.BloodVolumeNeeded {
				reserved = b.BloodRequest.BloodVolumeNeeded
			}
		}
	}
	b.BloodRequest.BloodVolumeReserved = math.Round(reserved*10) / 10
}
