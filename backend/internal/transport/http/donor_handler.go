package http

import (
	"context"
	"net/http"
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/application/donor/cmd"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/application/donor/query"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/common"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/filestorage"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/transport/http/dto"
	commondto "github.com/artesipov-alt/odnoi-krovi-app/internal/transport/http/dto/common"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/transport/http/middleware"

	"github.com/danielgtaylor/huma/v2"
)

type DonorHandler struct {
	recipientsListHandler     *query.ListRequestsHandler
	recipientDetailsHandler   *query.RecipientDetailHandler
	applyHandler              *cmd.ApplyForRequestHandler
	plannedDonationsList      *query.PlannedDonationsHandler
	completedDonationsHandler *query.CompletedDonationsHandler
	completeDonationHandler   *cmd.CompleteDonationHandler
	cancelDonationHandler     *cmd.CancelDonationHandler
	storage                   filestorage.Repository
}

// NewDonorHandler creates a new handler for donor-related operations.
func NewDonorHandler(
	recipientsListHandler *query.ListRequestsHandler,
	recipientDetailsHandler *query.RecipientDetailHandler,
	applyHandler *cmd.ApplyForRequestHandler,
	plannedDonationsList *query.PlannedDonationsHandler,
	completedDonationsHandler *query.CompletedDonationsHandler,
	completeDonationHandler *cmd.CompleteDonationHandler,
	cancelDonationHandler *cmd.CancelDonationHandler,
	storage filestorage.Repository,
) *DonorHandler {
	return &DonorHandler{
		recipientsListHandler:     recipientsListHandler,
		recipientDetailsHandler:   recipientDetailsHandler,
		applyHandler:              applyHandler,
		plannedDonationsList:      plannedDonationsList,
		completedDonationsHandler: completedDonationsHandler,
		completeDonationHandler:   completeDonationHandler,
		cancelDonationHandler:     cancelDonationHandler,
		storage:                   storage,
	}
}

// Register регистрирует маршруты заявок на поиск крови в Huma API
func (h *DonorHandler) Register(api huma.API) {
	// Получить список заявок на поиск крови
	huma.Register(api, huma.Operation{
		OperationID: "get-recipients",
		Method:      http.MethodGet,
		Path:        "/v1/donor/recipient-list/{user_id}", // Изменено
		Summary:     "Получить список заявок на поиск крови",
		Description: "Возвращает список реципиентов по фильтрам",
		Tags:        []string{"donor-v1"},
	}, h.GetRecipientsList)

	// Получить детальную заявку на поиск крови
	huma.Register(api, huma.Operation{
		OperationID: "get-recipient-details",
		Method:      http.MethodGet,
		Path:        "/v1/donor/recipient-details/{req_id}", // Изменено
		Summary:     "Получить детальные данные по заявке на поиск крови",
		Description: "Возвращает детальную информацию по заявке на поиск крови",
		Tags:        []string{"donor-v1"},
	}, h.GetRecipientDetails)

	// Откликнуться на заявку
	huma.Register(api, huma.Operation{
		OperationID:   "apply-for-blood-request",
		Method:        http.MethodPost,
		Path:          "/v1/donor/recipient/{req_id}/apply", // Изменено
		Summary:       "Откликнуться на заявку на поиск крови",
		Description:   "Позволяет донору откликнуться на существующую заявку на поиск крови.",
		Tags:          []string{"donor-v1"},
		DefaultStatus: http.StatusCreated,
	}, h.ApplyForBloodRequest)

	// Получить список планируемых донаций
	huma.Register(api, huma.Operation{
		OperationID: "get-planned-donations",
		Method:      http.MethodGet,
		Path:        "/v1/donor/planned-donations/{user_id}",
		Summary:     "Получить список планируемых донаций",
		Description: "Возвращает список планируемых донаций по user ID",
		Tags:        []string{"donor-v1"},
	}, h.GetPlannedDonations)

	// Получить список завершенных донаций
	huma.Register(api, huma.Operation{
		OperationID: "get-completed-donations",
		Method:      http.MethodGet,
		Path:        "/v1/donor/completed-donations/{user_id}",
		Summary:     "Получить список завершенных донаций",
		Description: "Возвращает список завершенных донаций по user ID",
		Tags:        []string{"donor-v1"},
	}, h.GetCompletedDonations)

	// Подтвердить донацию
	huma.Register(api, huma.Operation{
		OperationID: "complete-donation",
		Method:      http.MethodPost,
		Path:        "/v1/donor/donation/{res_id}/complete",
		Summary:     "Подтвердить донацию",
		Description: "Помечает донацию как состоявшуюся",
		Tags:        []string{"donor-v1"},
	}, h.CompleteDonation)

	// Отменить донацию
	huma.Register(api, huma.Operation{
		OperationID: "cancel-donation",
		Method:      http.MethodPost,
		Path:        "/v1/donor/donation/{res_id}/cancel",
		Summary:     "Отменить донацию",
		Description: "Отменяет запланированную донацию",
		Tags:        []string{"donor-v1"},
	}, h.CancelDonation)

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
	now := time.Now()
	items := make([]dto.RecipientDetail, len(result))
	for i, r := range result {
		matching := make([]dto.MatchingDonor, len(r.MatchingDonors))
		for j, md := range r.MatchingDonors {
			matching[j] = dto.MatchingDonor{
				PetID:           md.PetID,
				PetName:         md.PetName,
				DonorBloodGroup: md.DonorBloodGroup,
				Amount:          md.Amount,
				PhotoURLs:       h.storage.BuildPhotoURLs(md.PhotoURLs, now),
			}
		}
		items[i] = dto.RecipientDetail{
			ID:                       r.ID,
			PetID:                    r.PetID,
			PetName:                  r.RecipientData.PetName,
			SmallPetsNotifyAllowed:   r.SmallPetsNotifyAllowed,
			IncludeUnknownBloodGroup: r.IncludeUnknownBloodGroup,
			PetType:                  string(r.RecipientData.PetType),
			BloodVolumeRemaining:     r.BloodVolumeNeeded - r.BloodVolumeReserved,
			PhotoURLs:                h.storage.BuildPhotoURLs(r.RecipientData.PhotoURLs, now),
			BloodGroupName:           r.RecipientData.BloodGroupName,
			PrioritySearch:           r.PrioritySearch,
			Status:                   string(r.Status),
			MatchingDonors:           matching,
		}
	}

	return &dto.ListRecipientsOutput{Body: dto.RecipientsList{Items: items, Total: len(items)}}, nil
}

func (h *DonorHandler) GetRecipientDetails(ctx context.Context, input *commondto.BloodRequestIDPath) (*dto.RecipientDetailsOutput, error) {
	userID := middleware.GetUserID(ctx)
	if userID == "" {
		return nil, apperrors.Unauthorized("user ID is missing in context")
	}

	recipientData, err := h.recipientDetailsHandler.Handle(ctx, input.ID, userID)
	if err != nil {
		return nil, err
	}
	// Собираем подходящих доноров в ДТО.
	now := time.Now()
	matchingDonors := make([]dto.MatchingDonor, len(recipientData.Recipient.MatchingDonors))
	for i, md := range recipientData.Recipient.MatchingDonors {
		matchingDonors[i] = dto.MatchingDonor{
			PetID:           md.PetID,
			PetName:         md.PetName,
			Amount:          md.Amount,
			DonorBloodGroup: md.DonorBloodGroup,
			PhotoURLs:       h.storage.BuildPhotoURLs(md.PhotoURLs, now),
		}
	}
	// Собираем подходящие бонусы в ДТО.
	avilableBonuses := make([]common.Bonus, len(recipientData.AvilableBonuses))
	for i, bns := range recipientData.AvilableBonuses {
		avilableBonuses[i] = common.Bonus{
			Partner:     bns.PartnerName,
			Description: bns.Description,
			Type:        bns.Category,
		}
	}

	// Собираем дефолтные настройки донороа в ДТО.
	var defaultPrefs *dto.DefaultDonorPrefs
	if recipientData.Recipient.DefaultDonorPrefs != nil {
		defaultPrefs = &dto.DefaultDonorPrefs{
			CompensationType: string(recipientData.Recipient.DefaultDonorPrefs.CompensationType),
			Bonuses:          recipientData.Recipient.DefaultDonorPrefs.Bonuses,
			TaxiCompensation: recipientData.Recipient.DefaultDonorPrefs.TaxiCompensation,
		}
	}

	recipientDetail := dto.RecipientDetail{
		ID:                       recipientData.Recipient.ID,
		PetID:                    recipientData.Recipient.PetID,
		PetName:                  recipientData.Recipient.RecipientData.PetName,
		PetType:                  string(recipientData.Recipient.RecipientData.PetType),
		OwnerName:                recipientData.Recipient.RecipientData.OwnerName,
		SearchRegions:            recipientData.Recipient.Regions,
		BloodVolumeNeeded:        recipientData.Recipient.BloodVolumeNeeded,
		BloodVolumeReserved:      recipientData.Recipient.BloodVolumeReserved,
		SearchingBloodNames:      recipientData.Recipient.BloodGroupNames,
		SmallPetsNotifyAllowed:   recipientData.Recipient.SmallPetsNotifyAllowed,
		IncludeUnknownBloodGroup: recipientData.Recipient.IncludeUnknownBloodGroup,
		PhotoURLs:                h.storage.BuildPhotoURLs(recipientData.Recipient.RecipientData.PhotoURLs, now),
		BloodGroupName:           recipientData.Recipient.RecipientData.BloodGroupName,
		PrioritySearch:           recipientData.Recipient.PrioritySearch,
		Status:                   string(recipientData.Recipient.Status),
		MatchingDonors:           matchingDonors,
		DefaultDonorPrefs:        defaultPrefs,
		AvailableBonuses:         avilableBonuses,
		AdvancedInfo: &dto.AdvancedInfo{
			Description: recipientData.Recipient.AdvancedInfo.Description,
			PhotoURLs:   h.storage.BuildPhotoURLs(recipientData.Recipient.AdvancedInfo.PhotoURLs, now),
		},
	}

	return &dto.RecipientDetailsOutput{Body: recipientDetail}, nil
}

func (h *DonorHandler) ApplyForBloodRequest(ctx context.Context, input *dto.ApplyForBloodRequestInput) (*dto.ApplyForBloodRequestOutput, error) {
	resp, err := h.applyHandler.Handle(ctx, input.ID, input.Body.DonorID, input.Body.CompensationType, input.Body.TaxiCompensation)
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

// GetPlannedDonations returns a list of planned donations.
func (h *DonorHandler) GetPlannedDonations(ctx context.Context, input *commondto.UserIDPath) (*dto.ListPlannedDonationsOutput, error) {
	results, err := h.plannedDonationsList.Handle(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	donationCards := make([]dto.DonationCardForDonor, 0, len(results))
	for _, res := range results {
		// В планируемой донации не должно возвращаться отмененный статус донором и отмененный репертипиентом и прошедшие
		if res.ApplicationData.Status != model.DonorResponseStatusCancelled && res.ApplicationData.Status != model.DonorResponseStatusRejected && res.ApplicationData.IsConfirmed != true {
			bonusDTOs := make([]common.Bonus, len(res.Bonuses))
			for i, b := range res.Bonuses {
				bonusDTOs[i] = common.Bonus{
					Partner:     b.PartnerName,
					Description: b.Description,
					Type:        string(b.Category),
				}
			}
			application := dto.ApplicationShort{
				ID:               res.ApplicationData.ID,
				PetName:          res.DonorPetData.Name,
				Amount:           res.ApplicationData.Amount,
				PhotoURLs:        h.storage.BuildPhotoURLs(res.DonorPetData.PhotoURLs, *res.ApplicationData.UpdatedAt),
				CompensationType: res.ApplicationData.CompensationType,
				TaxiCompensation: res.ApplicationData.TaxiCompensation,
				IsConfirmed:      res.ApplicationData.IsConfirmed,
				Bonuses:          bonusDTOs,
				RejectedReason:   res.ApplicationData.RejectedReason,
				Status:           string(res.ApplicationData.Status),
			}

			recipient := dto.RecipientForDonor{
				ID:                       res.BloodSearchData.ID,
				PetName:                  res.RecipientPetData.Name,
				PetType:                  string(res.RecipientPetData.Type),
				OwnerName:                res.RecipientPetData.OwnerName,
				OwnerID:                  res.RecipientPetData.OwnerID,
				BloodGroup:               res.RecipientPetData.BloodGroupName,
				Regions:                  res.BloodSearchData.Regions,
				BloodVolumeNeeded:        res.BloodSearchData.BloodVolumeNeeded,
				BloodVolumeReserved:      res.BloodSearchData.BloodVolumeReserved,
				BloodVolumeDonated:       res.BloodSearchData.BloodVolumeDonated,
				IncludeUnknownBloodGroup: res.BloodSearchData.IncludeUnknownBloodGroup,
				PhotoURLs:                h.storage.BuildPhotoURLs(res.RecipientPetData.PhotoURLs, *res.RecipientPetData.UpdatedAt),
				SearchingBloodNames:      res.BloodSearchData.BloodGroupNames,
				AdvancedInfo: &dto.AdvancedInfoDTO{
					PhotoURLs:   h.storage.BuildPhotoURLs(res.BloodSearchData.AdvancedInfo.PhotoURLs, *res.BloodSearchData.UpdatedAt),
					Description: res.BloodSearchData.AdvancedInfo.Description,
				},
				Status:    string(res.BloodSearchData.Status),
				CreatedAt: res.BloodSearchData.CreatedAt,
				UpdatedAt: res.BloodSearchData.UpdatedAt,
			}

			donationCards = append(donationCards, dto.DonationCardForDonor{
				ApplicationData: application,
				RecipientData:   recipient,
			})
		}
	}

	return &dto.ListPlannedDonationsOutput{
		Body: dto.PlannedDonationsList{
			Items: donationCards,
			Total: len(donationCards),
		},
	}, nil
}

func (h *DonorHandler) CompleteDonation(ctx context.Context, input *dto.CompleteDonationInput) (*commondto.DefaultMessageOutput, error) {
	err := h.completeDonationHandler.Handle(ctx, input.ID, input.Body.Amount)
	if err != nil {
		return nil, err
	}
	return &commondto.DefaultMessageOutput{Body: commondto.ResultMessage{Message: "Донация успешно завершена"}}, nil
}

// CancelDonation отменяет запланированную донацию.
func (h *DonorHandler) CancelDonation(ctx context.Context, input *commondto.DonorApplicationIDPath) (*commondto.DefaultMessageOutput, error) {
	err := h.cancelDonationHandler.Handle(ctx, input.ID)
	if err != nil {
		return nil, err
	}
	return &commondto.DefaultMessageOutput{Body: commondto.ResultMessage{Message: "Донация успешно отменена"}}, nil
}

// GetCompletedDonations returns a list of completed donations.
func (h *DonorHandler) GetCompletedDonations(ctx context.Context, input *commondto.UserIDPath) (*dto.ListCompletedDonationsOutput, error) {
	results, err := h.completedDonationsHandler.Handle(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	donationCards := make([]dto.DonationCardForDonor, 0, len(results))
	var totalDonatedVolume float64
	var totalCompletedDonations int
	for _, res := range results {
		bonusDTOs := make([]common.Bonus, len(res.Bonuses))
		for i, b := range res.Bonuses {
			bonusDTOs[i] = common.Bonus{
				Partner:     b.PartnerName,
				Description: b.Description,
				Type:        string(b.Category),
			}
		}
		application := dto.ApplicationShort{
			ID:               res.ApplicationData.ID,
			PetName:          res.DonorPetData.Name,
			Amount:           res.ApplicationData.Amount,
			PhotoURLs:        h.storage.BuildPhotoURLs(res.DonorPetData.PhotoURLs, *res.ApplicationData.UpdatedAt),
			CompensationType: res.ApplicationData.CompensationType,
			TaxiCompensation: res.ApplicationData.TaxiCompensation,
			IsConfirmed:      res.ApplicationData.IsConfirmed,
			Bonuses:          bonusDTOs,
			RejectedReason:   res.ApplicationData.RejectedReason,
			Status:           string(res.ApplicationData.Status),
			CreatedAt:        res.ApplicationData.CreatedAt,
			UpdatedAt:        res.ApplicationData.UpdatedAt,
		}

		recipient := dto.RecipientForDonor{
			ID:                  res.BloodSearchData.ID,
			PetName:             res.RecipientPetData.Name,
			PetType:             string(res.RecipientPetData.Type),
			OwnerName:           res.RecipientPetData.OwnerName,
			OwnerID:             res.RecipientPetData.OwnerID,
			BloodGroup:          res.RecipientPetData.BloodGroupName,
			Regions:             res.BloodSearchData.Regions,
			BloodVolumeNeeded:   res.BloodSearchData.BloodVolumeNeeded,
			BloodVolumeReserved: res.BloodSearchData.BloodVolumeReserved,
			BloodVolumeDonated:  res.BloodSearchData.BloodVolumeDonated,
			PhotoURLs:           h.storage.BuildPhotoURLs(res.RecipientPetData.PhotoURLs, *res.RecipientPetData.UpdatedAt),
			SearchingBloodNames: res.BloodSearchData.BloodGroupNames,
			AdvancedInfo: &dto.AdvancedInfoDTO{
				PhotoURLs:   h.storage.BuildPhotoURLs(res.BloodSearchData.AdvancedInfo.PhotoURLs, *res.BloodSearchData.UpdatedAt),
				Description: res.BloodSearchData.AdvancedInfo.Description,
			},
			Status:    string(res.BloodSearchData.Status),
			CreatedAt: res.BloodSearchData.CreatedAt,
			UpdatedAt: res.BloodSearchData.UpdatedAt,
		}

		donationCards = append(donationCards, dto.DonationCardForDonor{
			ApplicationData: application,
			RecipientData:   recipient,
		})

		if res.ApplicationData.Status != model.DonorResponseStatusRejected && res.ApplicationData.Status != model.DonorResponseStatusCancelled {
			totalDonatedVolume += res.ApplicationData.Amount
			totalCompletedDonations++
		}
	}

	return &dto.ListCompletedDonationsOutput{
		Body: dto.CompletedDonationsList{
			Items:                   donationCards,
			Total:                   len(donationCards),
			TotalCompletedDonations: totalCompletedDonations,
			TotalDonatedVolume:      totalDonatedVolume,
		},
	}, nil
}
