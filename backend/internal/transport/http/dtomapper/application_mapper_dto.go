// Package mapper provides conversion functions between domain models and DTOs.
package mapper

import (
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/transport/http/dto"
)

// ApplicationToDTO converts a domain Donorresponse model to a Application DTO.
func ApplicationToDTO(d model.DonorResponse) dto.DonorApplication {

	return dto.DonorApplication{
		ID:               d.ID,
		RequestID:        d.RequestID,
		DonorID:          d.DonorID,
		DonorName:        d.DonorName,
		DonorPhotos:      d.DonorPhotos,
		DonorBloodGroup:  d.DonorBloodGroup,
		Amount:           d.Amount,
		CompensationType: d.DonorPrefs.CompensationType,
		TaxiCompensation: d.DonorPrefs.TaxiCompensation,
		Status:           string(d.Status),
		IsConfirmed:      d.IsConfirmed,
		CreatedAt:        d.CreatedAt,
		UpdatedAt:        d.UpdatedAt,
	}
}
