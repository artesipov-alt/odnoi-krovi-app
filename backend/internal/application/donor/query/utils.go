package query

import donormodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/model"

// findActiveApplication finds the most recent active donor application from the list.
// Active means status is Accepted, Pending, or Completed but not confirmed.
func findActiveApplication(applications []*donormodel.DonorResponse) *donormodel.DonorResponse {
	for _, app := range applications {
		if app.Status == donormodel.DonorResponseStatusAccepted ||
			app.Status == donormodel.DonorResponseStatusPending ||
			(app.Status == donormodel.DonorResponseStatusCompleted && !app.IsConfirmed) {
			return app
		}
	}
	return nil
}
