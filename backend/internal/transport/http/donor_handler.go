package http

import (
	"context"
	"net/http"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/application/donor/cmd"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/application/donor/query"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/transport/http/dto"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/transport/http/middleware"

	"github.com/danielgtaylor/huma/v2"
)

type DonorHandler struct {
	recipientsListHandler   *query.ListRequestsHandler
	recipientDetailsHandler *query.RecipientDetailHandler
	applyHandler            *cmd.ApplyForRequestHandler
}

// NewDonorHandler creates a new handler for donor-related operations.
func NewDonorHandler(
	recipientsListHandler *query.ListRequestsHandler,
	recipientDetailsHandler *query.RecipientDetailHandler,
	applyHandler *cmd.ApplyForRequestHandler,
) *DonorHandler {
	return &DonorHandler{
		recipientsListHandler:   recipientsListHandler,
		recipientDetailsHandler: recipientDetailsHandler,
		applyHandler:            applyHandler,
	}
}

// Register регистрирует маршруты заявок на поиск крови в Huma API
func (h *DonorHandler) Register(api huma.API) {
	// Получить список заявок на поиск крови
	huma.Register(api, huma.Operation{
		OperationID: "get-recipients",
		Method:      http.MethodGet,
		Path:        "/v1/donor/recipient-list/{id}", // Изменено
		Summary:     "Получить список заявок на поиск крови",
		Description: "Возвращает список реципиентов по фильтрам",
		Tags:        []string{"donor-v1"},
	}, h.GetRecipientsList)

	// Получить детальную заявку на поиск крови
	huma.Register(api, huma.Operation{
		OperationID: "get-recipient-details",
		Method:      http.MethodGet,
		Path:        "/v1/donor/recipient-details/{id}", // Изменено
		Summary:     "Получить детальные данные по заявке на поиск крови",
		Description: "Возвращает детальную информацию по заявке на поиск крови",
		Tags:        []string{"donor-v1"},
	}, h.GetRecipientDetails)

	// Откликнуться на заявку
	huma.Register(api, huma.Operation{
		OperationID:   "apply-for-blood-request",
		Method:        http.MethodPost,
		Path:          "/v1/donor/recipient/{id}/apply", // Изменено
		Summary:       "Откликнуться на заявку на поиск крови",
		Description:   "Позволяет донору откликнуться на существующую заявку на поиск крови.",
		Tags:          []string{"donor-v1"},
		DefaultStatus: http.StatusCreated,
	}, h.ApplyForBloodRequest)

}

func (h *DonorHandler) GetRecipientsList(ctx context.Context, input *dto.GetRecipientsListInput) (*dto.ListRecipientsOutput, error) {
	result, err := h.recipientsListHandler.Handle(ctx, input.UserIDPath.ID, model.DonorPreloadFilter{
		Status: string(input.Status),
		Limit:  input.Limit,
		Offset: input.Offset,
	})
	if err != nil {
		return nil, err
	}

	items := make([]dto.RecipientDetail, len(result))
	for i, r := range result {
		matching := make([]dto.MatchingDonor, len(r.MatchingDonors))
		for j, md := range r.MatchingDonors {
			matching[j] = dto.MatchingDonor{
				PetID:           md.PetID,
				PetName:         md.PetName,
				DonorBloodGroup: md.DonorBloodGroup,
				PhotoURLs:       md.PhotoURLs,
			}
		}
		items[i] = dto.RecipientDetail{
			ID:                   r.ID,
			PetID:                r.PetID,
			PetName:              r.PetName,
			PetType:              string(r.PetType),
			BloodVolumeRemaining: r.BloodVolumeRemaining,
			PhotoURLs:            r.PhotoURLs,
			BloodGroupName:       r.BloodGroupName,
			PrioritySearch:       r.PrioritySearch,
			Status:               dto.BloodRequestStatus(r.Status),
			MatchingDonors:       matching,
		}
	}

	return &dto.ListRecipientsOutput{Body: dto.RecipientsList{Items: items, Total: len(items)}}, nil
}

func (h *DonorHandler) GetRecipientDetails(ctx context.Context, input *struct{ dto.BloodRequestIDPath }) (*dto.RecipientDetailsOutput, error) {
	userID := middleware.GetUserID(ctx)
	if userID == "" {
		return nil, apperrors.Unauthorized("user ID is missing in context")
	}

	recipient, err := h.recipientDetailsHandler.Handle(ctx, input.BloodRequestIDPath.ID, userID)
	if err != nil {
		return nil, err
	}

	matchingDonors := make([]dto.MatchingDonor, len(recipient.MatchingDonors))
	for i, md := range recipient.MatchingDonors {
		matchingDonors[i] = dto.MatchingDonor{
			PetID:           md.PetID,
			PetName:         md.PetName,
			PetType:         string(md.PetType),
			DonorBloodGroup: md.DonorBloodGroup,
			PhotoURLs:       md.PhotoURLs,
		}
	}

	var defaultPrefs *dto.DefaultDonorPrefs
	if recipient.DefaultDonorPrefs != nil {
		defaultPrefs = &dto.DefaultDonorPrefs{
			CompensationType: string(recipient.DefaultDonorPrefs.CompensationType),
			Bonuses:          recipient.DefaultDonorPrefs.Bonuses,
			TaxiCompensation: recipient.DefaultDonorPrefs.TaxiCompensation,
		}
	}

	recipientDetail := dto.RecipientDetail{
		ID:                   recipient.ID,
		PetID:                recipient.PetID,
		PetName:              recipient.PetName,
		PetType:              string(recipient.PetType),
		OwnerName:            recipient.OwnerName,
		SearchRegions:        recipient.SearchRegions,
		BloodVolumeRemaining: recipient.BloodVolumeRemaining,
		SearchingBloodNames:  recipient.SearchingBloodNames,
		PhotoURLs:            recipient.PhotoURLs,
		BloodGroupName:       recipient.BloodGroupName,
		PrioritySearch:       recipient.PrioritySearch,
		Status:               dto.BloodRequestStatus(recipient.Status),
		MatchingDonors:       matchingDonors,
		DefaultDonorPrefs:    defaultPrefs,
	}

	return &dto.RecipientDetailsOutput{Body: recipientDetail}, nil
}

func (h *DonorHandler) ApplyForBloodRequest(ctx context.Context, input *dto.ApplyForBloodRequestInput) (*dto.ApplyForBloodRequestOutput, error) {
	resp, err := h.applyHandler.Handle(ctx, input.ID, input.Body.DonorID, input.Body.Conditions)
	if err != nil {
		return nil, err
	}

	return &dto.ApplyForBloodRequestOutput{
		Body: dto.DonorApplicationResult{
			ID:        resp.ID,
			RequestID: resp.RequestID,
			DonorID:   resp.DonorID,
			Status:    string(resp.Status),
			CreatedAt: resp.CreatedAt,
		},
	}, nil
}
