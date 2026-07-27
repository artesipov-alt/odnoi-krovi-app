package pet

import (
	bloodreqmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/model"
	donormodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet/model"
)

func BuildDonationContext(application *donormodel.DonorResponse, bloodReq *bloodreqmodel.BloodRequestWithApplications) model.DonationContext {
	isRecipient := bloodReq != nil && !bloodReq.IsClosed()
	return model.DonationContext{
		IsRecipient:                isRecipient,
		HasActiveDonorApplications: isRecipient && bloodReq.HasActiveDonorApplications(),
		IsPlanningDonation:         application != nil && application.IsActiveForDonation(),
		IsPrioritySearch:           isRecipient && bloodReq.BloodRequest.PrioritySearch,
	}
}
