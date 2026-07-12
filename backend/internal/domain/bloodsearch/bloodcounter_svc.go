package bloodsearch

import (
	"math"

	bloodreqmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/model"
	donormodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/model"
)

type BloodCounterService struct{}

func NewBloodCounterService() *BloodCounterService {
	return &BloodCounterService{}
}

func (b *BloodCounterService) RecalculateBloodAmount(bloodreq bloodreqmodel.BloodRequest, applications []donormodel.DonorResponse) (donated float64, reserved float64) {
	for _, app := range applications {
		if app.IsFullyCompleted() {
			donated += app.Amount
		}
	}
	donated = math.Round(donated*10) / 10

	reserved = donated
	for _, app := range applications {
		if app.IsAwaitingConfirmation() {
			reserved += app.Amount
			// Обрезаем до максимального
			if reserved >= bloodreq.BloodVolumeNeeded {
				reserved = bloodreq.BloodVolumeNeeded
			}
		}
	}

	reserved = math.Round(reserved*10) / 10

	return donated, reserved
}

func (b *BloodCounterService) HasActiveDonorApplications(applications []donormodel.DonorResponse) bool {
	for _, app := range applications {
		if app.IsActiveForDonation() {
			return true
		}
	}
	return false
}
