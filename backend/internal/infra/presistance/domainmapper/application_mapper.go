package domainmapper

import (
	donormodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent"
)

// ApplicationToDomain converts ENT DonorResponse to domain DonorResponse
func ApplicationToDomain(entResp *ent.DonorResponse) *donormodel.DonorResponse {
	if entResp == nil {
		return nil
	}

	return &donormodel.DonorResponse{
		ID:               entResp.ID,
		RequestID:        entResp.Edges.Request.ID,
		DonorID:          entResp.Edges.Donor.ID,
		Amount:           entResp.Amount,
		CompensationType: string(entResp.CompensationType),
		TaxiCompensation: entResp.TaxiCompensation,
		Status:           donormodel.DonorResponseStatus(entResp.Status),
		IsConfirmed:      entResp.IsConfirmed,
		RejectedReason:   entResp.RejectedReason,
		CreatedAt:        &entResp.CreatedAt,
		UpdatedAt:        &entResp.UpdatedAt,
	}
}
