package http

import (
	"context"
	"log/slog"
	"net/http"

	bloodcmd "github.com/artesipov-alt/odnoi-krovi-app/internal/application/bloodsearch/cmd"
	bloodquery "github.com/artesipov-alt/odnoi-krovi-app/internal/application/bloodsearch/query"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/filestorage"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/transport/http/dto"
	mapper "github.com/artesipov-alt/odnoi-krovi-app/internal/transport/http/dtomapper"

	"github.com/danielgtaylor/huma/v2"
)

// BloodRequestHandler обрабатывает HTTP запросы для операций с заявками на поиск крови
type BloodRequestHandler struct {
	createHandler        *bloodcmd.CreateRequestHandler
	updateHandler        *bloodcmd.UpdateRequestHandler
	updateStatusHandler  *bloodcmd.UpdateStatusHandler
	deleteHandler        *bloodcmd.DeleteRequestHandler
	getByIDHandler       *bloodquery.GetByIDHandler
	getByPetIDHandler    *bloodquery.GetByPetIDHandler
	getDonorByIDHandler  *bloodquery.GetDonorByIDHandler
	getDonationHandler   *bloodquery.GetDonationHandler
	applyResponseHandler *bloodcmd.ApplyResponseHandler
	bloodRequestMapper   *mapper.BloodRequestMapper
	petMapper            *mapper.PetMapper
	storage              filestorage.Repository
}

func NewBloodRequestHandler(
	createHandler *bloodcmd.CreateRequestHandler,
	updateHandler *bloodcmd.UpdateRequestHandler,
	updateStatusHandler *bloodcmd.UpdateStatusHandler,
	deleteHandler *bloodcmd.DeleteRequestHandler,
	getByIDHandler *bloodquery.GetByIDHandler,
	getByPetIDHandler *bloodquery.GetByPetIDHandler,
	getDonorByIDHandler *bloodquery.GetDonorByIDHandler,
	getDonationHandler *bloodquery.GetDonationHandler,
	applyResponseHandler *bloodcmd.ApplyResponseHandler,
	storage filestorage.Repository,
) *BloodRequestHandler {
	return &BloodRequestHandler{
		createHandler:        createHandler,
		updateHandler:        updateHandler,
		updateStatusHandler:  updateStatusHandler,
		deleteHandler:        deleteHandler,
		getByIDHandler:       getByIDHandler,
		getByPetIDHandler:    getByPetIDHandler,
		getDonorByIDHandler:  getDonorByIDHandler,
		getDonationHandler:   getDonationHandler,
		applyResponseHandler: applyResponseHandler,
		bloodRequestMapper:   mapper.NewBloodRequestMapper(storage),
		petMapper:            mapper.NewPetMapper(storage),
		storage:              storage,
	}
}

// Register регистрирует маршруты заявок на поиск крови в Huma API
func (h *BloodRequestHandler) Register(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID:   "add-pet-to-blood-request-pool",
		Method:        http.MethodPost,
		Path:          "/v1/blood-request/pool",
		Summary:       "Добавить питомца в пул поиска крови",
		Description:   "Создает новую заявку на поиск крови для питомца",
		Tags:          []string{"blood-request-v1"},
		DefaultStatus: http.StatusCreated,
	}, h.AddPetToBloodRequestPool)

	huma.Register(api, huma.Operation{
		OperationID: "get-blood-request-by-id",
		Method:      http.MethodGet,
		Path:        "/v1/blood-request/{id}",
		Summary:     "Получить заявку по ID",
		Description: "Возвращает информацию о конкретной заявке",
		Tags:        []string{"blood-request-v1"},
	}, h.GetBloodRequestByID)

	huma.Register(api, huma.Operation{
		OperationID: "get-blood-request-by-pet-id",
		Method:      http.MethodGet,
		Path:        "/v1/blood-request/pet/{id}",
		Summary:     "Получить заявку по ID питомца",
		Description: "Возвращает информацию о конкретной заявке",
		Tags:        []string{"blood-request-v1"},
	}, h.GetBloodRequestByPetID)

	huma.Register(api, huma.Operation{
		OperationID: "get-donor-by-id",
		Method:      http.MethodGet,
		Path:        "/v1/blood-request/donor/{id}",
		Summary:     "Получить информацию о доноре по ID",
		Description: "Возвращает информацию об откликнувшемся на заявку доноре",
		Tags:        []string{"blood-request-v1"},
	}, h.GetDonorByID)

	huma.Register(api, huma.Operation{
		OperationID: "get-donation-by-id",
		Method:      http.MethodGet,
		Path:        "/v1/blood-request/donation/{id}",
		Summary:     "Получить информацию о донации по ID отклика донора",
		Description: "Возвращает детальную информацию о донации по ID отклика",
		Tags:        []string{"blood-request-v1"},
	}, h.GetDonation)

	huma.Register(api, huma.Operation{
		OperationID: "update-blood-request",
		Method:      http.MethodPatch,
		Path:        "/v1/blood-request/{id}",
		Summary:     "Обновить заявку на поиск крови",
		Description: "Частично обновляет информацию о существующей заявке на поиск крови.",
		Tags:        []string{"blood-request-v1"},
	}, h.UpdateBloodRequest)

	huma.Register(api, huma.Operation{
		OperationID: "delete-blood-request",
		Method:      http.MethodDelete,
		Path:        "/v1/blood-request/{id}",
		Summary:     "Удалить заявку",
		Description: "Удаляет заявку на поиск крови (soft delete)",
		Tags:        []string{"blood-request-v1"},
	}, h.DeleteBloodRequest)

	huma.Register(api, huma.Operation{
		OperationID: "apply-donor-response",
		Method:      http.MethodPost,
		Path:        "/v1/blood-request/apply-response/{id}",
		Summary:     "Применить отклик донора",
		Description: "Применяет отклик донора на заявку на поиск крови",
		Tags:        []string{"blood-request-v1"},
	}, h.ApplyResponse)
}

// Handlers

func (h *BloodRequestHandler) AddPetToBloodRequestPool(ctx context.Context, input *dto.CreateBloodRequestInput) (*dto.CreateBloodRequestOutput, error) {
	bloodReq := h.bloodRequestMapper.FromCreate(input.Body)

	result, err := h.createHandler.Handle(ctx, bloodReq)
	if err != nil {
		return nil, err
	}

	return &dto.CreateBloodRequestOutput{Body: dto.CreateBloodRequestResult{
		ID:        result.ID,
		PetID:     result.PetID,
		Status:    string(result.Status),
		CreatedAt: result.CreatedAt,
	}}, nil
}

func (h *BloodRequestHandler) UpdateBloodRequest(ctx context.Context, input *dto.UpdateBloodRequestInput) (*dto.UpdateBloodRequestOutput, error) {
	slog.DebugContext(ctx, "updating blood request", "request_id", input.ID)

	// Получить текущий объект
	existing, err := h.getByIDHandler.Handle(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	// Частично обновить поля
	if input.Body.BloodVolumeNeeded != nil {
		existing.BloodVolumeNeeded = *input.Body.BloodVolumeNeeded
	}
	if input.Body.BloodVolumeReserved != nil {
		existing.BloodVolumeReserved = *input.Body.BloodVolumeReserved
	}
	if len(input.Body.Regions) > 0 {
		existing.Regions = input.Body.Regions
	}
	if input.Body.SmallPetsNotifyAllowed != nil {
		existing.SmallPetsNotifyAllowed = *input.Body.SmallPetsNotifyAllowed
	}
	if input.Body.Description != nil {
		existing.Description = *input.Body.Description
	}
	if len(input.Body.BloodGroupNames) > 0 {
		existing.BloodGroupNames = input.Body.BloodGroupNames
	}
	if len(input.Body.BloodComponentIDs) > 0 {
		existing.BloodComponentIDs = input.Body.BloodComponentIDs
	}
	if len(input.Body.OnBoarding) > 0 {
		existing.OnBoarding = input.Body.OnBoarding
	}
	if input.Body.Status != nil {
		existing.Status = model.BloodRequestStatus(*input.Body.Status)
	}
	if input.Body.PrioritySearch != nil {
		existing.PrioritySearch = *input.Body.PrioritySearch
	}
	if input.Body.IncludeUnknownBloodGroup != nil {
		existing.IncludeUnknownBloodGroup = *input.Body.IncludeUnknownBloodGroup
	}

	result, err := h.updateHandler.Handle(ctx, input.ID, existing)
	if err != nil {
		return nil, err
	}

	return &dto.UpdateBloodRequestOutput{Body: dto.UpdateBloodRequestResult{
		ID:        result.ID,
		UpdatedAt: result.UpdatedAt,
	}}, nil
}

func (h *BloodRequestHandler) GetBloodRequestByID(ctx context.Context, input *dto.BloodRequestIDPath) (*dto.GetBloodRequestByIDOutput, error) {
	bloodReq, err := h.getByIDHandler.Handle(ctx, input.ID)
	if err != nil {
		return nil, err
	}
	zero := 0
	return &dto.GetBloodRequestByIDOutput{Body: h.bloodRequestMapper.ToResponse(bloodReq, &zero)}, nil
}

func (h *BloodRequestHandler) GetBloodRequestByPetID(ctx context.Context, input *dto.PetIDPath) (*dto.GetBloodRequestByPetIDOutput, error) {
	bloodReq, situatableDonors, err := h.getByPetIDHandler.Handle(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	return &dto.GetBloodRequestByPetIDOutput{Body: h.bloodRequestMapper.ToResponse(bloodReq, &situatableDonors)}, nil
}

func (h *BloodRequestHandler) GetDonorByID(ctx context.Context, input *dto.PetIDPath) (*dto.GetDonorByIDOutput, error) {
	donorPet, application, err := h.getDonorByIDHandler.Handle(ctx, input.ID, pet.PetPreloadOptions{
		WithHealth:     true,
		WithTreatments: true,
		WithAnalyses:   true,
		WithBonuses:    true,
		WithOwner:      true,
	})
	if err != nil {
		return nil, err
	}

	output := h.petMapper.ToResponse(*donorPet)

	return &dto.GetDonorByIDOutput{Body: dto.DonorDetail{
		ResponseID: application.ID,
		PetDetail:  output,
		Compensation: dto.Compensation{
			CompensationType: application.CompensationType,
			Taxi:             application.TaxiCompensation,
		},
	}}, nil
}

func (h *BloodRequestHandler) DeleteBloodRequest(ctx context.Context, input *dto.BloodRequestIDPath) (*dto.DeleteBloodRequestOutput, error) {
	slog.DebugContext(ctx, "deleting blood request", "request_id", input.ID)
	if err := h.deleteHandler.Handle(ctx, input.ID); err != nil {
		return nil, err
	}

	return &dto.DeleteBloodRequestOutput{Body: dto.DeleteBloodRequestResult{
		Message: "Заявка удалена",
	}}, nil
}

func (h *BloodRequestHandler) ApplyResponse(ctx context.Context, input *dto.DonorApplicationIDPath) (*dto.ApplyResponseOutput, error) {
	if err := h.applyResponseHandler.Handle(ctx, input.ID); err != nil {
		return nil, err
	}

	return &dto.ApplyResponseOutput{Body: dto.ApplyResponseResult{
		Message: "Отклик донора применен",
	}}, nil
}

func (h *BloodRequestHandler) GetDonation(ctx context.Context, input *dto.DonorApplicationIDPath) (*dto.GetDonationOutput, error) {
	donation, err := h.getDonationHandler.Handle(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	donor := h.petMapper.ToResponse(*donation.DonorPet)

	donorData := dto.PetWithApplication{
		PetDetail: donor,
		Application: dto.CoreApplicationData{
			ID:               donation.Application.ID,
			Amount:           donation.Application.Amount,
			CompensationType: donation.Application.CompensationType,
			TaxiCompensation: donation.Application.TaxiCompensation,
			Status:           string(donation.Application.Status),
		},
	}

	recipientData := dto.RecipientShort{
		PetName:             donation.RecipientPet.Name,
		PetType:             string(donation.RecipientPet.Type),
		BloodVolumeNeeded:   donation.BloodRequest.BloodVolumeNeeded,
		BloodVolumeReserved: donation.BloodRequest.BloodVolumeReserved,
		PhotoURLs:           donation.RecipientPet.PhotoURLs,
	}

	donationCard := dto.DonationCard{
		DonorData:     donorData,
		RecipientData: recipientData,
	}

	return &dto.GetDonationOutput{Body: donationCard}, nil
}
