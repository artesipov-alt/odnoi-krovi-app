package http

import (
	"context"
	"net/http"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	bloodcmd "github.com/artesipov-alt/odnoi-krovi-app/internal/application/bloodsearch/cmd"
	bloodquery "github.com/artesipov-alt/odnoi-krovi-app/internal/application/bloodsearch/query"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/filestorage"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/transport/http/dto"
	commondto "github.com/artesipov-alt/odnoi-krovi-app/internal/transport/http/dto/common"
	mapper "github.com/artesipov-alt/odnoi-krovi-app/internal/transport/http/dtomapper"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/transport/http/middleware"

	"github.com/danielgtaylor/huma/v2"
)

// BloodRequestHandler обрабатывает HTTP запросы для операций с заявками на поиск крови
type BloodRequestHandler struct {
	createHandler              *bloodcmd.CreateRequestHandler
	updateHandler              *bloodcmd.UpdateRequestHandler
	deleteHandler              *bloodcmd.DeleteRequestHandler
	getByIDHandler             *bloodquery.GetByIDHandler
	getByPetIDHandler          *bloodquery.GetByPetIDHandler
	getDonorByIDHandler        *bloodquery.GetDonorByIDHandler
	getDonationHandler         *bloodquery.GetDonationHandler
	applyResponseHandler       *bloodcmd.ApplyResponseHandler
	confirmDonationHandler     *bloodcmd.ConfirmDonationHandler
	rejectDonationHandler      *bloodcmd.RejectDonationHandler
	closeRequestHandler        *bloodcmd.CloseRequestHandler
	notificationRespondHandler *bloodcmd.NotificationRespondHandler
	selectDonorHandler         *bloodcmd.SelectDonorHandler
	bloodRequestMapper         *mapper.BloodRequestMapper
	petMapper                  *mapper.PetMapper
	storage                    filestorage.Repository
}

func NewBloodRequestHandler(
	createHandler *bloodcmd.CreateRequestHandler,
	updateHandler *bloodcmd.UpdateRequestHandler,
	deleteHandler *bloodcmd.DeleteRequestHandler,
	getByIDHandler *bloodquery.GetByIDHandler,
	getByPetIDHandler *bloodquery.GetByPetIDHandler,
	getDonorByIDHandler *bloodquery.GetDonorByIDHandler,
	getDonationHandler *bloodquery.GetDonationHandler,
	applyResponseHandler *bloodcmd.ApplyResponseHandler,
	confirmDonationHandler *bloodcmd.ConfirmDonationHandler,
	rejectDonationHandler *bloodcmd.RejectDonationHandler,
	closeRequestHandler *bloodcmd.CloseRequestHandler,
	notificationRespondHandler *bloodcmd.NotificationRespondHandler,
	selectDonorHandler *bloodcmd.SelectDonorHandler,
	storage filestorage.Repository,
) *BloodRequestHandler {
	return &BloodRequestHandler{
		createHandler:              createHandler,
		updateHandler:              updateHandler,
		deleteHandler:              deleteHandler,
		getByIDHandler:             getByIDHandler,
		getByPetIDHandler:          getByPetIDHandler,
		getDonorByIDHandler:        getDonorByIDHandler,
		getDonationHandler:         getDonationHandler,
		applyResponseHandler:       applyResponseHandler,
		confirmDonationHandler:     confirmDonationHandler,
		rejectDonationHandler:      rejectDonationHandler,
		closeRequestHandler:        closeRequestHandler,
		notificationRespondHandler: notificationRespondHandler,
		selectDonorHandler:         selectDonorHandler,
		bloodRequestMapper:         mapper.NewBloodRequestMapper(storage),
		petMapper:                  mapper.NewPetMapper(storage),
		storage:                    storage,
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
		Path:        "/v1/blood-request/{req_id}",
		Summary:     "Получить заявку по ID",
		Description: "Возвращает информацию о конкретной заявке",
		Tags:        []string{"blood-request-v1"},
	}, h.GetBloodRequestByID)

	huma.Register(api, huma.Operation{
		OperationID: "get-blood-request-by-pet-id",
		Method:      http.MethodGet,
		Path:        "/v1/blood-request/pet/{pet_id}",
		Summary:     "Получить заявку по ID питомца",
		Description: "Возвращает информацию о конкретной заявке",
		Tags:        []string{"blood-request-v1"},
	}, h.GetBloodRequestByPetID)

	huma.Register(api, huma.Operation{
		OperationID: "get-donor-by-id",
		Method:      http.MethodGet,
		Path:        "/v1/blood-request/donor/{pet_id}",
		Summary:     "Получить информацию о доноре по ID",
		Description: "Возвращает информацию об откликнувшемся на заявку доноре",
		Tags:        []string{"blood-request-v1"},
	}, h.GetDonorByID)

	huma.Register(api, huma.Operation{
		OperationID: "get-donation-by-id",
		Method:      http.MethodGet,
		Path:        "/v1/blood-request/donation/{res_id}",
		Summary:     "Получить информацию о донации по ID",
		Description: "Возвращает детальную информацию о донации по ID",
		Tags:        []string{"blood-request-v1"},
	}, h.GetDonation)

	huma.Register(api, huma.Operation{
		OperationID: "update-blood-request",
		Method:      http.MethodPatch,
		Path:        "/v1/blood-request/{req_id}",
		Summary:     "Обновить заявку на поиск крови",
		Description: "Частично обновляет информацию о существующей заявке на поиск крови.",
		Tags:        []string{"blood-request-v1"},
	}, h.UpdateBloodRequest)

	huma.Register(api, huma.Operation{
		OperationID: "delete-blood-request",
		Method:      http.MethodDelete,
		Path:        "/v1/blood-request/{req_id}",
		Summary:     "Удалить заявку",
		Description: "Удаляет заявку на поиск крови (soft delete)",
		Tags:        []string{"blood-request-v1"},
	}, h.DeleteBloodRequest)

	huma.Register(api, huma.Operation{
		OperationID: "apply-donor-response",
		Method:      http.MethodPost,
		Path:        "/v1/blood-request/apply-response/{res_id}",
		Summary:     "Принять отклик донора",
		Description: "Принимает отклик донора на заявку на поиск крови",
		Tags:        []string{"blood-request-v1"},
	}, h.ApplyResponse)

	huma.Register(api, huma.Operation{
		OperationID: "confirm-donation-by-id",
		Method:      http.MethodPost,
		Path:        "/v1/blood-request/donation/{res_id}/confirm",
		Summary:     "Подтвердить донацию по ID",
		Description: "Подтверждает факт проведения донации по ID отклика донора",
		Tags:        []string{"blood-request-v1"},
	}, h.ConfirmDonation)

	huma.Register(api, huma.Operation{
		OperationID: "reject-donation-by-id",
		Method:      http.MethodPost,
		Path:        "/v1/blood-request/donation/{res_id}/reject",
		Summary:     "Отклонить донацию по ID",
		Description: "Отклоняет факт проведения донации по ID отклика донора",
		Tags:        []string{"blood-request-v1"},
	}, h.RejectDonation)

	huma.Register(api, huma.Operation{
		OperationID: "close-blood-request-by-id",
		Method:      http.MethodPost,
		Path:        "/v1/blood-request/close/{req_id}",
		Summary:     "Закрыть заявку на поиск крови по ID заявки",
		Description: "Закрывает заявку на поиск крови по ID заявки",
		Tags:        []string{"blood-request-v1"},
	}, h.CloseBloodSearch)

	huma.Register(api, huma.Operation{
		OperationID: "respond-to-notification",
		Method:      http.MethodPost,
		Path:        "/v1/blood-request/notification/respond/{req_id}",
		Summary:     "Ответить на уведомление",
		Description: "Отвечает на уведомление",
		Tags:        []string{"blood-request-v1"},
	}, h.NotificationRespond)

	huma.Register(api, huma.Operation{
		OperationID: "select-donor",
		Method:      http.MethodPost,
		Path:        "/v1/blood-request/{req_id}/donor/select",
		Summary:     "Выбрать донора из списка потенциальных",
		Description: "Реципиент выбирает конкретного донора из списка потенциальных (open for contact). Создаёт DonorResponse со статусом accepted.",
		Tags:        []string{"blood-request-v1"},
	}, h.SelectDonor)

}

// Handlers

func (h *BloodRequestHandler) AddPetToBloodRequestPool(ctx context.Context, input *dto.CreateBloodRequestInput) (*dto.CreateBloodRequestOutput, error) {
	userID := middleware.GetUserID(ctx)
	if userID == "" {
		return nil, apperrors.Unauthorized("user ID is missing in context")
	}
	bloodReq := h.bloodRequestMapper.FromCreate(input.Body)
	result, err := h.createHandler.Handle(ctx, bloodReq)
	if err != nil {
		return nil, err
	}

	return &dto.CreateBloodRequestOutput{Body: dto.CreateBloodRequestResult{
		ID:        result.BloodRequest.ID,
		PetID:     result.BloodRequest.PetID,
		Status:    string(result.BloodRequest.Status),
		CreatedAt: result.BloodRequest.CreatedAt,
	}}, nil
}

func (h *BloodRequestHandler) UpdateBloodRequest(ctx context.Context, input *dto.UpdateBloodRequestInput) (*dto.UpdateBloodRequestOutput, error) {
	// Получить текущий объект
	existing, _, err := h.getByIDHandler.Handle(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	// Частично обновить поля
	if input.Body.BloodVolumeNeeded != nil {
		existing.BloodRequest.BloodVolumeNeeded = *input.Body.BloodVolumeNeeded
	}
	if input.Body.BloodVolumeReserved != nil {
		existing.BloodRequest.BloodVolumeReserved = *input.Body.BloodVolumeReserved
	}
	if len(input.Body.Regions) > 0 {
		existing.BloodRequest.Regions = input.Body.Regions
	}
	if input.Body.SmallPetsNotifyAllowed != nil {
		existing.BloodRequest.SmallPetsNotifyAllowed = *input.Body.SmallPetsNotifyAllowed
	}
	if input.Body.Description != nil {
		existing.BloodRequest.AdvancedInfo.Description = *input.Body.Description
	}
	if len(input.Body.BloodGroupNames) > 0 {
		existing.BloodRequest.BloodGroupNames = input.Body.BloodGroupNames
	}
	if len(input.Body.BloodComponentIDs) > 0 {
		existing.BloodRequest.BloodComponentIDs = input.Body.BloodComponentIDs
	}
	if len(input.Body.OnBoarding) > 0 {
		existing.BloodRequest.OnBoarding = input.Body.OnBoarding
	}
	if input.Body.Status != nil {
		existing.BloodRequest.Status = model.BloodRequestStatus(*input.Body.Status)
	}
	if input.Body.PrioritySearch != nil {
		existing.BloodRequest.PrioritySearch = *input.Body.PrioritySearch
	}
	if input.Body.IncludeUnknownBloodGroup != nil {
		existing.BloodRequest.IncludeUnknownBloodGroup = *input.Body.IncludeUnknownBloodGroup
	}

	result, err := h.updateHandler.Handle(ctx, input.ID, existing)
	if err != nil {
		return nil, err
	}

	return &dto.UpdateBloodRequestOutput{Body: dto.UpdateBloodRequestResult{
		ID:        result.BloodRequest.ID,
		UpdatedAt: result.BloodRequest.UpdatedAt,
	}}, nil
}

func (h *BloodRequestHandler) GetBloodRequestByID(ctx context.Context, input *commondto.BloodRequestIDPath) (*dto.GetBloodRequestByIDOutput, error) {
	bloodReq, situatableDonors, err := h.getByIDHandler.Handle(ctx, input.ID)
	if err != nil {
		return nil, err
	}
	return &dto.GetBloodRequestByIDOutput{Body: h.bloodRequestMapper.ToResponse(bloodReq, &situatableDonors)}, nil
}

func (h *BloodRequestHandler) GetBloodRequestByPetID(ctx context.Context, input *commondto.PetIDPath) (*dto.GetBloodRequestByPetIDOutput, error) {
	userID := middleware.GetUserID(ctx)
	if userID == "" {
		return nil, apperrors.Unauthorized("user ID is missing in context")
	}

	result, err := h.getByPetIDHandler.Handle(ctx, userID, input.ID)
	if err != nil {
		return nil, err
	}

	return &dto.GetBloodRequestByPetIDOutput{Body: h.bloodRequestMapper.ToResponseWithPotential(result)}, nil
}

func (h *BloodRequestHandler) GetDonorByID(ctx context.Context, input *commondto.PetIDPath) (*dto.GetDonorByIDOutput, error) {
	donorPet, application, err := h.getDonorByIDHandler.Handle(ctx, input.ID, pet.PetPreloadOptions{
		WithAll: true,
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

func (h *BloodRequestHandler) DeleteBloodRequest(ctx context.Context, input *commondto.BloodRequestIDPath) (*commondto.DefaultMessageOutput, error) {
	if err := h.deleteHandler.Handle(ctx, input.ID); err != nil {
		return nil, err
	}

	return &commondto.DefaultMessageOutput{Body: commondto.ResultMessage{
		Message: "Заявка удалена",
	}}, nil
}

func (h *BloodRequestHandler) ApplyResponse(ctx context.Context, input *commondto.DonorApplicationIDPath) (*commondto.DefaultMessageOutput, error) {
	if err := h.applyResponseHandler.Handle(ctx, input.ID); err != nil {
		return nil, err
	}

	return &commondto.DefaultMessageOutput{Body: commondto.ResultMessage{
		Message: "Отклик донора применен",
	}}, nil
}

func (h *BloodRequestHandler) GetDonation(ctx context.Context, input *commondto.DonorApplicationIDPath) (*dto.GetDonationOutput, error) {
	donation, err := h.getDonationHandler.Handle(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	donor := h.petMapper.ToResponse(*donation.DonorPet)

	identities := make([]dto.Identity, 0, len(donation.DonorOwnerData.Identities))
	for _, identity := range donation.DonorOwnerData.Identities {
		identities = append(identities, dto.Identity{
			ProviderName: string(identity.ProviderName),
			ProviderID:   identity.ProviderUserID,
		})
	}
	donorData := dto.PetWithApplication{
		PetDetail:  donor,
		OwnerPhone: donation.DonorOwnerData.Phone,
		Identities: identities,
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
		BloodVolumeDonated:  donation.BloodRequest.BloodRequest.BloodVolumeDonated,
		BloodVolumeNeeded:   donation.BloodRequest.BloodRequest.BloodVolumeNeeded,
		BloodVolumeReserved: donation.BloodRequest.BloodRequest.BloodVolumeReserved,
		PhotoURLs:           h.storage.BuildPhotoURLs(donation.RecipientPet.PhotoURLs, *donation.Application.UpdatedAt),
	}

	donationCard := dto.DonationCard{
		DonorData:     donorData,
		RecipientData: recipientData,
	}

	return &dto.GetDonationOutput{Body: donationCard}, nil
}

func (h *BloodRequestHandler) ConfirmDonation(ctx context.Context, input *dto.ConfirmDonorApplicationInput) (*commondto.DefaultMessageOutput, error) {
	if err := h.confirmDonationHandler.Handle(ctx, input.ID, input.Body.Amount); err != nil {
		return nil, err
	}
	return &commondto.DefaultMessageOutput{Body: commondto.ResultMessage{Message: "Донация успешно подтверждена"}}, nil
}

func (h *BloodRequestHandler) RejectDonation(ctx context.Context, input *dto.RejectDonorApplicationInput) (*commondto.DefaultMessageOutput, error) {
	var reason string
	if input.Body.Reason != nil {
		reason = *input.Body.Reason
	}

	if err := h.rejectDonationHandler.Handle(ctx, input.ID, reason); err != nil {
		return nil, err
	}
	return &commondto.DefaultMessageOutput{Body: commondto.ResultMessage{Message: "Донация отменена"}}, nil
}

func (h *BloodRequestHandler) CloseBloodSearch(ctx context.Context, input *commondto.BloodRequestIDPath) (*commondto.DefaultMessageOutput, error) {
	if err := h.closeRequestHandler.Handle(ctx, input.ID); err != nil {
		return nil, err
	}
	return &commondto.DefaultMessageOutput{Body: commondto.ResultMessage{Message: "Заявка успешно закрыта, все невыполненные донации отменены"}}, nil
}

func (h *BloodRequestHandler) NotificationRespond(ctx context.Context, input *dto.NotificationRespondInput) (*commondto.DefaultMessageOutput, error) {
	if err := h.notificationRespondHandler.Handle(ctx, bloodcmd.NotificationRespondInput{
		BloodRequestID: input.ID,
		Action:         bloodcmd.NotificationAction(input.Body.Action),
	}); err != nil {
		return nil, err
	}
	return &commondto.DefaultMessageOutput{Body: commondto.ResultMessage{
		Message: "Ответ принят",
	}}, nil
}

// SelectDonor выбирает донора из списка потенциальных (recipient-initiated).
// Действие адресовано конкретной заявке (req_id в URL), донор передаётся в теле.
// Условия донации (компенсация, такси) берутся из DonorPreference владельца донора,
// а не из тела запроса.
func (h *BloodRequestHandler) SelectDonor(ctx context.Context, input *dto.SelectDonorInput) (*dto.SelectDonorOutput, error) {
	userID := middleware.GetUserID(ctx)
	if userID == "" {
		return nil, apperrors.Unauthorized("user ID is missing in context")
	}

	resp, err := h.selectDonorHandler.Handle(
		ctx,
		userID,
		input.ID,
		input.Body.DonorID,
	)
	if err != nil {
		return nil, err
	}

	return &dto.SelectDonorOutput{Body: h.bloodRequestMapper.DonorResponseToSelected(resp)}, nil
}
