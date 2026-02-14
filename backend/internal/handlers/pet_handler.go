package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/ent"

	"github.com/artesipov-alt/odnoi-krovi-app/ent/petanalysis"
	"github.com/artesipov-alt/odnoi-krovi-app/ent/pethealth"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/dto"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/services"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/validator"
	"github.com/danielgtaylor/huma/v2"
	"github.com/jinzhu/copier"
)

// PetHandler обрабатывает HTTP запросы для операций с питомцами
type PetHandler struct {
	petService services.PetService
	validator  validator.DonorValidator
}

// NewPetHandler создает новый обработчик питомцев
func NewPetHandler(petService services.PetService, validator validator.DonorValidator) *PetHandler {
	return &PetHandler{
		petService: petService,
		validator:  validator,
	}
}

// Register регистрирует маршруты питомцев в Huma API
func (h *PetHandler) Register(api huma.API) {
	// Создание нового питомца
	huma.Register(api, huma.Operation{
		OperationID:   "create-pet",
		Method:        http.MethodPost,
		Path:          "/v1/pet/user/{user_id}",
		Summary:       "Создание нового питомца",
		Description:   "Создает нового питомца для пользователя",
		Tags:          []string{"pets-v1"},
		DefaultStatus: http.StatusCreated,
	}, h.CreatePet)

	// Получение питомца по ID
	huma.Register(api, huma.Operation{
		OperationID: "get-pet-by-id",
		Method:      http.MethodGet,
		Path:        "/v1/pet/{id}",
		Summary:     "Получение питомца по ID",
		Description: "Возвращает информацию о питомце по его идентификатору",
		Tags:        []string{"pets-v1"},
	}, h.GetPet)

	// Получение питомцев пользователя
	huma.Register(api, huma.Operation{
		OperationID: "get-user-pets",
		Method:      http.MethodGet,
		Path:        "/v1/pet/user/{user_id}",
		Summary:     "Получение питомцев пользователя",
		Description: "Возвращает всех питомцев конкретного пользователя",
		Tags:        []string{"pets-v1"},
	}, h.GetUserPets)

	// Обновление данных питомца
	huma.Register(api, huma.Operation{
		OperationID: "update-pet",
		Method:      http.MethodPut,
		Path:        "/v1/pet/{id}",
		Summary:     "Обновление данных питомца",
		Description: "Обновляет информацию о питомце",
		Tags:        []string{"pets-v1"},
	}, h.UpdatePet)

	// Удаление питомца по ID
	huma.Register(api, huma.Operation{
		OperationID: "delete-pet",
		Method:      http.MethodDelete,
		Path:        "/v1/pet/{id}",
		Summary:     "Удаление питомца по ID",
		Description: "Удаляет питомца из системы",
		Tags:        []string{"pets-v1"},
	}, h.DeletePet)

	// Валидация донора по ID (изменено на POST)
	huma.Register(api, huma.Operation{
		OperationID:   "validate-donor",
		Method:        http.MethodPost,
		Path:          "/v1/pet/validate-donor/{id}",
		Summary:       "Валидация донора по ID",
		Description:   "Пересчитывает и сохраняет факторы валидации донора для питомца",
		Tags:          []string{"pets-v1"},
		DefaultStatus: http.StatusOK,
	}, h.ValidateDonor)

}

//==========Handlers==============================

func (h *PetHandler) CreatePet(ctx context.Context, input *struct {
	dto.PetUserIDPath
	Body dto.PetCreate
}) (*dto.PetResponse, error) {
	// Рассчитываем BirthDate, если нужно
	body := input.Body
	if body.BirthDate == nil && (body.AgeYears > 0 || body.AgeMonths > 0) {
		birthDate := time.Now().AddDate(-body.AgeYears, -body.AgeMonths, 0)
		birthDate = time.Date(birthDate.Year(), birthDate.Month(), 1, 0, 0, 0, 0, time.UTC)
		body.BirthDate = &birthDate
	}

	// Копируем в CreatePetInput
	var petInput ent.CreatePetInput
	if err := copier.Copy(&petInput, &body); err != nil {
		slog.ErrorContext(ctx, "failed to copy pet create data", "error", err.Error())
		return nil, apperrors.Internal(err, "failed to copy pet create data")
	}

	// Копируем health
	var healthInput *ent.CreatePetHealthInput
	if body.Health != nil {
		healthInput = &ent.CreatePetHealthInput{}
		if body.Health.HealthStatus != nil {
			status := pethealth.HealthStatus(*body.Health.HealthStatus)
			healthInput.HealthStatus = &status
		}
		healthInput.LastDonation = body.Health.LastDonation
		healthInput.Transfused = body.Health.Transfused
		healthInput.Medications = body.Health.Medications
		healthInput.SurgicalInterventions = body.Health.SurgicalInterventions
	}

	// Копируем treatments
	var treatmentsInput *ent.CreatePetTreatmentInput
	if body.Treatments != nil {
		treatmentsInput = &ent.CreatePetTreatmentInput{}
		if err := copier.Copy(treatmentsInput, body.Treatments); err != nil {
			slog.ErrorContext(ctx, "failed to copy treatments data", "error", err.Error())
			return nil, apperrors.Internal(err, "failed to copy treatments data")
		}
	}

	// Копируем analyses
	var analysesInput []*ent.CreatePetAnalysisInput
	if body.Analyses != nil {
		processGroup := func(group []*dto.PetAnalysis, name petanalysis.AnalysisName) {
			for _, a := range group {
				var ai ent.CreatePetAnalysisInput
				ai.AnalysisName = &name
				if a.AnalysisDate != nil {
					ai.AnalysisDate = *a.AnalysisDate
				}
				if a.AnalysisType != nil {
					at := petanalysis.AnalysisType(*a.AnalysisType)
					ai.AnalysisType = &at
				}
				analysesInput = append(analysesInput, &ai)
			}
		}
		processGroup(body.Analyses.Leukemia, petanalysis.AnalysisNameLeukemia)
		processGroup(body.Analyses.Immunodeficiency, petanalysis.AnalysisNameImmunodeficiency)
		processGroup(body.Analyses.Hemoplasmosis, petanalysis.AnalysisNameHemoplasmosis)
		processGroup(body.Analyses.Bartonellosis, petanalysis.AnalysisNameBartonellosis)
		processGroup(body.Analyses.Babesiosis, petanalysis.AnalysisNameBabesiosis)
		processGroup(body.Analyses.Dirofilaria, petanalysis.AnalysisNameDirofilaria)
		processGroup(body.Analyses.Ehrlichiosis, petanalysis.AnalysisNameEhrlichiosis)
		processGroup(body.Analyses.Anaplasmosis, petanalysis.AnalysisNameAnaplasmosis)
	}

	// Копируем bonuses
	var bonusesInput *ent.CreatePetBonusInput
	if body.Bonuses != nil {
		bonusesInput = &ent.CreatePetBonusInput{}
		if err := copier.Copy(bonusesInput, body.Bonuses); err != nil {
			slog.ErrorContext(ctx, "failed to copy bonuses data", "error", err.Error())
			return nil, apperrors.Internal(err, "failed to copy bonuses data")
		}
	}

	pet, err := h.petService.CreatePet(ctx, input.ID, &petInput, healthInput, treatmentsInput, analysesInput, bonusesInput)
	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			slog.DebugContext(ctx, "user not found for pet creation", "user_id", input.ID, "error", err.Error())
			return nil, huma.Error404NotFound("Пользователь не найден")
		}
		slog.ErrorContext(ctx, "failed to create pet", "user_id", input.ID, "error", err.Error())
		return nil, huma.Error500InternalServerError("Внутренняя ошибка сервера")
	}

	return &dto.PetResponse{Body: h.toDTO(pet)}, nil
}

func (h *PetHandler) GetPet(ctx context.Context, input *struct {
	dto.IDPathStr
	dto.PetPreloadQuery
}) (*dto.PetResponse, error) {
	query := h.petService.GetPetQuery(ctx, input.ID)
	if input.WithHealth {
		query = query.WithHealth()
	}
	if input.WithTreatments {
		query = query.WithTreatments()
	}
	if input.WithAnalysis {
		query = query.WithAnalyses()
	}
	if input.WithBonuses {
		query = query.WithBonuses()
	}
	if input.WithAll {
		query = query.WithHealth().WithTreatments().WithAnalyses().WithBonuses()
	}

	pet, err := query.Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, apperrors.ErrPetNotFound
		}
		return nil, apperrors.Internal(err, "failed to get pet")
	}

	return &dto.PetResponse{Body: h.toDTO(pet)}, nil
}

func (h *PetHandler) GetUserPets(ctx context.Context, input *struct {
	dto.PetUserIDPath
	dto.PetPreloadQuery
}) (*dto.PetsResponse, error) {
	query := h.petService.GetPetsQueryByUser(ctx, input.ID)
	if input.WithHealth {
		query = query.WithHealth()
	}
	if input.WithTreatments {
		query = query.WithTreatments()
	}
	if input.WithAnalysis {
		query = query.WithAnalyses()
	}
	if input.WithBonuses {
		query = query.WithBonuses()
	}
	if input.WithAll {
		query = query.WithHealth().WithTreatments().WithAnalyses().WithBonuses()
	}

	pets, err := query.All(ctx)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to get user pets")
	}

	// Преобразуем пути к фото в полные URL для каждого питомца
	for _, pet := range pets {
		h.petService.BuildFullPhotoURLs(pet)
	}

	return &dto.PetsResponse{Body: h.toPetsDTO(pets)}, nil
}

func (h *PetHandler) UpdatePet(ctx context.Context, input *struct {
	dto.IDPathStr
	Body dto.PetUpdate
}) (*dto.MessageResponse, error) {
	// Рассчитываем BirthDate, если нужно
	if input.Body.BirthDate == nil && (input.Body.AgeYears != nil || input.Body.AgeMonths != nil) {
		ageYears := 0
		if input.Body.AgeYears != nil {
			ageYears = *input.Body.AgeYears
		}
		ageMonths := 0
		if input.Body.AgeMonths != nil {
			ageMonths = *input.Body.AgeMonths
		}
		if ageYears > 0 || ageMonths > 0 {
			birthDate := time.Now().AddDate(-ageYears, -ageMonths, 0)
			birthDate = time.Date(birthDate.Year(), birthDate.Month(), 1, 0, 0, 0, 0, time.UTC)
			input.Body.BirthDate = &birthDate
		}
	}

	// Копируем данные в UpdatePetInput
	var petInput ent.UpdatePetInput
	if err := copier.Copy(&petInput, &input.Body); err != nil {
		slog.ErrorContext(ctx, "failed to copy pet update data", "error", err.Error())
		return nil, apperrors.Internal(err, "failed to copy pet update data")
	}

	// Подготавливаем связанные данные
	var health *ent.UpdatePetHealthInput
	if input.Body.Health != nil {
		health = &ent.UpdatePetHealthInput{}
		if err := copier.Copy(health, &input.Body.Health); err != nil {
			return nil, apperrors.Internal(err, "failed to copy pet health data")
		}
	}

	var treatments *ent.UpdatePetTreatmentInput
	if input.Body.Treatments != nil {
		treatments = &ent.UpdatePetTreatmentInput{}
		if err := copier.Copy(treatments, &input.Body.Treatments); err != nil {
			return nil, apperrors.Internal(err, "failed to copy treatments data")
		}
	}

	var analyses []*ent.UpdatePetAnalysisInput
	if input.Body.Analyses != nil {
		analyses = []*ent.UpdatePetAnalysisInput{}
		processGroup := func(group []*dto.PetAnalysis, name petanalysis.AnalysisName) {
			for _, a := range group {
				entA := &ent.UpdatePetAnalysisInput{
					AnalysisName: &name,
					AnalysisDate: a.AnalysisDate,
				}
				if a.AnalysisType != nil {
					at := petanalysis.AnalysisType(*a.AnalysisType)
					entA.AnalysisType = &at
				}
				analyses = append(analyses, entA)
			}
		}

		processGroup(input.Body.Analyses.Leukemia, petanalysis.AnalysisNameLeukemia)
		processGroup(input.Body.Analyses.Immunodeficiency, petanalysis.AnalysisNameImmunodeficiency)
		processGroup(input.Body.Analyses.Hemoplasmosis, petanalysis.AnalysisNameHemoplasmosis)
		processGroup(input.Body.Analyses.Bartonellosis, petanalysis.AnalysisNameBartonellosis)
		processGroup(input.Body.Analyses.Babesiosis, petanalysis.AnalysisNameBabesiosis)
		processGroup(input.Body.Analyses.Dirofilaria, petanalysis.AnalysisNameDirofilaria)
		processGroup(input.Body.Analyses.Ehrlichiosis, petanalysis.AnalysisNameEhrlichiosis)
		processGroup(input.Body.Analyses.Anaplasmosis, petanalysis.AnalysisNameAnaplasmosis)
	}

	var bonuses *ent.UpdatePetBonusInput
	if input.Body.Bonuses != nil {
		bonuses = &ent.UpdatePetBonusInput{}
		if err := copier.Copy(bonuses, &input.Body.Bonuses); err != nil {
			return nil, apperrors.Internal(err, "failed to copy bonuses data")
		}
	}

	// Обновляем питомца и все связанные данные одним вызовом
	_, err := h.petService.Update(ctx, input.ID, &petInput, health, treatments, analyses, bonuses)
	if err != nil {
		if errors.Is(err, apperrors.ErrPetNotFound) {
			slog.DebugContext(ctx, "pet not found for update", "pet_id", input.ID, "error", err.Error())
			return nil, huma.Error404NotFound("Питомец не найден")
		}
		slog.ErrorContext(ctx, "failed to update pet", "pet_id", input.ID, "error", err.Error())
		return nil, huma.Error500InternalServerError("Внутренняя ошибка сервера")
	}

	resp := &dto.MessageResponse{}
	resp.Body.Message = "Питомец успешно обновлен"
	return resp, nil
}

func (h *PetHandler) DeletePet(ctx context.Context, input *dto.IDPathStr) (*dto.MessageResponse, error) {
	if err := h.petService.DeletePet(ctx, input.ID); err != nil {
		if errors.Is(err, apperrors.ErrPetNotFound) {
			slog.DebugContext(ctx, "pet not found for deletion", "pet_id", input.ID, "error", err.Error())
			return nil, huma.Error404NotFound("Питомец не найден")
		}
		slog.ErrorContext(ctx, "failed to delete pet", "pet_id", input.ID, "error", err.Error())
		return nil, huma.Error500InternalServerError("Внутренняя ошибка сервера")
	}

	resp := &dto.MessageResponse{}
	resp.Body.Message = "Питомец успешно удален"
	return resp, nil
}

func (h *PetHandler) ValidateDonor(ctx context.Context, input *dto.IDPathStr) (*dto.PetResponse, error) {
	// Загружаем все связанные данные для полной валидации
	preloads := []string{"Health", "Treatments", "Analyses", "Bonuses"}
	p, err := h.petService.GetPetByID(ctx, input.ID, preloads...)
	if err != nil {
		if errors.Is(err, apperrors.ErrPetNotFound) {
			slog.DebugContext(ctx, "pet not found for donor validation", "pet_id", input.ID, "error", err.Error())
			return nil, huma.Error404NotFound("Питомец не найден")
		}
		slog.ErrorContext(ctx, "failed to get pet for donor validation", "pet_id", input.ID, "error", err.Error())
		return nil, huma.Error500InternalServerError("Внутренняя ошибка сервера")
	}

	// Применяем валидацию: пересчитываем и сохраняем факторы
	_, _, err = h.petService.ApplyValidation(ctx, p.ID)
	if err != nil {
		slog.ErrorContext(ctx, "failed to apply validation", "pet_id", input.ID, "error", err.Error())
		return nil, huma.Error500InternalServerError("Ошибка валидации")
	}

	// Возвращаем обновлённый DTO (факторы теперь сохранены и будут включены)
	resp := &dto.PetResponse{}
	resp.Body = h.toDTO(p)
	return resp, nil
}

// mapPetToDTO преобразует ENT модель питомца в DTO для ответа
func (h *PetHandler) toDTO(p *ent.Pet) dto.Pet {
	petDTO := dto.Pet{
		ID:                 p.ID,
		Name:               p.Name,
		ChipNumber:         p.ChipNumber,
		PhotoURLs:          p.PhotoUrls,
		BreedID:            p.BreedID,
		WeightKg:           p.WeightKg,
		BirthDate:          p.BirthDate,
		LivingCondition:    string(p.LivingCondition),
		Gender:             string(p.Gender),
		Type:               string(p.Type),
		BloodGroup:         p.BloodGroup,
		ReproductiveStatus: string(p.ReproductiveStatus),
		PetStatus:          string(p.PetStatus),
		CreatedAt:          &p.CreatedAt,
		UpdatedAt:          &p.UpdatedAt,
		DeletedAt:          p.DeletedAt,
	}
	if p.Edges.Health != nil {
		healthStatus := string(p.Edges.Health.HealthStatus)
		medications := p.Edges.Health.Medications
		surgical := p.Edges.Health.SurgicalInterventions
		petDTO.Health = &dto.PetHealth{
			HealthStatus:          &healthStatus,
			LastDonation:          p.Edges.Health.LastDonation,
			Transfused:            &p.Edges.Health.Transfused,
			Medications:           &medications,
			SurgicalInterventions: &surgical,
		}
	}
	if p.Edges.Treatments != nil {
		petDTO.Treatments = &dto.PetTreatment{
			RabiesVaccinationDate:     p.Edges.Treatments.RabiesVaccinationDate,
			InfectionVaccinationDate:  p.Edges.Treatments.InfectionVaccinationDate,
			EctoparasiteTreatmentDate: p.Edges.Treatments.EctoparasiteTreatmentDate,
			DewormingDate:             p.Edges.Treatments.DewormingDate,
		}
	}
	if p.Edges.Analyses != nil {
		petDTO.Analyses = &dto.PetAnalysisGroup{}
		for _, a := range p.Edges.Analyses {
			analysisName := string(a.AnalysisName)
			analysisType := string(a.AnalysisType)
			dtoAnalysis := &dto.PetAnalysis{
				ID:           &a.ID,
				AnalysisName: &analysisName,
				AnalysisType: &analysisType,
				AnalysisDate: a.AnalysisDate,
			}

			switch a.AnalysisName {
			case petanalysis.AnalysisNameLeukemia:
				petDTO.Analyses.Leukemia = append(petDTO.Analyses.Leukemia, dtoAnalysis)
			case petanalysis.AnalysisNameImmunodeficiency:
				petDTO.Analyses.Immunodeficiency = append(petDTO.Analyses.Immunodeficiency, dtoAnalysis)
			case petanalysis.AnalysisNameHemoplasmosis:
				petDTO.Analyses.Hemoplasmosis = append(petDTO.Analyses.Hemoplasmosis, dtoAnalysis)
			case petanalysis.AnalysisNameBartonellosis:
				petDTO.Analyses.Bartonellosis = append(petDTO.Analyses.Bartonellosis, dtoAnalysis)
			case petanalysis.AnalysisNameBabesiosis:
				petDTO.Analyses.Babesiosis = append(petDTO.Analyses.Babesiosis, dtoAnalysis)
			case petanalysis.AnalysisNameDirofilaria:
				petDTO.Analyses.Dirofilaria = append(petDTO.Analyses.Dirofilaria, dtoAnalysis)
			case petanalysis.AnalysisNameEhrlichiosis:
				petDTO.Analyses.Ehrlichiosis = append(petDTO.Analyses.Ehrlichiosis, dtoAnalysis)
			case petanalysis.AnalysisNameAnaplasmosis:
				petDTO.Analyses.Anaplasmosis = append(petDTO.Analyses.Anaplasmosis, dtoAnalysis)
			}
		}
	}
	if p.Edges.Bonuses != nil {
		petDTO.Bonuses = &dto.PetBonus{
			IsArtist:      p.Edges.Bonuses.IsArtist,
			IsTherapist:   p.Edges.Bonuses.IsTherapist,
			IsFormerDonor: p.Edges.Bonuses.IsFormerDonor,
			IsGuideDog:    p.Edges.Bonuses.IsGuideDog,
		}
	}

	// Заполняем DonorRestrictions из сохранённых данных
	if p.PetStatus != "" && len(p.DonorRestrictions) > 0 {
		var stopRestrictionFactors []dto.RestrictionFactor
		var warnRestrictionFactors []dto.RestrictionFactor

		for _, codeStr := range p.DonorRestrictions {
			code := validator.FactorCode(codeStr)
			desc := validator.GetFactorDescription(code)

			factor := dto.RestrictionFactor{
				Code:           codeStr,
				Description:    desc.Description,
				SubDescription: desc.SubDescription,
			}

			// Разделяем на стопы и предупреждения по префиксу
			if strings.HasPrefix(codeStr, "STOP_") {
				stopRestrictionFactors = append(stopRestrictionFactors, factor)
			} else if strings.HasPrefix(codeStr, "WARN_") {
				warnRestrictionFactors = append(warnRestrictionFactors, factor)
			}
		}

		petDTO.DonorRestrictions = &dto.DonorRestrictions{
			StopFactors: stopRestrictionFactors,
			WarnFactors: warnRestrictionFactors,
		}
	} else if p.PetStatus != "" {
		// Если факторы ещё не рассчитаны, возвращаем пустые (не рассчитываем на лету)
		petDTO.DonorRestrictions = &dto.DonorRestrictions{
			StopFactors: []dto.RestrictionFactor{},
			WarnFactors: []dto.RestrictionFactor{},
		}
	}

	return petDTO
}

// toPetsDTO преобразует слайс ENT питомцев в слайс DTO
func (h *PetHandler) toPetsDTO(pets []*ent.Pet) []dto.Pet {
	petDTOs := make([]dto.Pet, len(pets))
	for i, p := range pets {
		petDTOs[i] = h.toDTO(p)
	}
	return petDTOs
}

// // mapDTOToPet преобразует DTO создания питомца в ENT модель
// func (h *PetHandler) toCreateENT(d dto.PetCreate) *ent.Pet {
// 	if d.BirthDate == nil && (d.AgeYears > 0 || d.AgeMonths > 0) {
// 		birthDate := time.Now().AddDate(-d.AgeYears, -d.AgeMonths, 0)
// 		birthDate = time.Date(birthDate.Year(), birthDate.Month(), 1, 0, 0, 0, 0, time.UTC)
// 		d.BirthDate = &birthDate
// 	}
// 	p := &ent.Pet{
// 		Name:               d.Name,
// 		ChipNumber:         d.ChipNumber,
// 		BreedID:            d.BreedID,
// 		WeightKg:           d.WeightKg,
// 		BirthDate:          d.BirthDate,
// 		LivingCondition:    pet.LivingCondition(d.LivingCondition),
// 		Gender:             pet.Gender(d.Gender),
// 		Type:               pet.Type(d.Type),
// 		BloodGroup:         d.BloodGroup,
// 		ReproductiveStatus: pet.ReproductiveStatus(d.ReproductiveStatus),
// 		PetStatus:          pet.PetStatus(d.PetStatus),
// 	}

// 	if d.Health != nil {
// 		p.Edges.Health = &ent.PetHealth{
// 			LastDonation: d.Health.LastDonation,
// 		}
// 		if d.Health.HealthStatus != nil {
// 			p.Edges.Health.HealthStatus = pethealth.HealthStatus(*d.Health.HealthStatus)
// 		}
// 		if d.Health.Transfused != nil {
// 			p.Edges.Health.Transfused = *d.Health.Transfused
// 		}
// 		if d.Health.Medications != nil {
// 			p.Edges.Health.Medications = *d.Health.Medications
// 		}
// 		if d.Health.SurgicalInterventions != nil {
// 			p.Edges.Health.SurgicalInterventions = *d.Health.SurgicalInterventions
// 		}
// 	}

// 	if d.Treatments != nil {
// 		p.Edges.Treatments = &ent.PetTreatment{
// 			RabiesVaccinationDate:     d.Treatments.RabiesVaccinationDate,
// 			InfectionVaccinationDate:  d.Treatments.InfectionVaccinationDate,
// 			EctoparasiteTreatmentDate: d.Treatments.EctoparasiteTreatmentDate,
// 			DewormingDate:             d.Treatments.DewormingDate,
// 		}
// 	}

// 	if d.Analyses != nil {
// 		var analyses []*ent.PetAnalysis
// 		processGroup := func(group []*dto.PetAnalysis, name petanalysis.AnalysisName) {
// 			for _, a := range group {
// 				entA := &ent.PetAnalysis{
// 					AnalysisName: name,
// 					AnalysisDate: a.AnalysisDate,
// 				}
// 				if a.AnalysisType != nil {
// 					entA.AnalysisType = petanalysis.AnalysisType(*a.AnalysisType)
// 				}
// 				analyses = append(analyses, entA)
// 			}
// 		}

// 		processGroup(d.Analyses.Leukemia, petanalysis.AnalysisNameLeukemia)
// 		processGroup(d.Analyses.Immunodeficiency, petanalysis.AnalysisNameImmunodeficiency)
// 		processGroup(d.Analyses.Hemoplasmosis, petanalysis.AnalysisNameHemoplasmosis)
// 		processGroup(d.Analyses.Bartonellosis, petanalysis.AnalysisNameBartonellosis)
// 		processGroup(d.Analyses.Babesiosis, petanalysis.AnalysisNameBabesiosis)
// 		processGroup(d.Analyses.Dirofilaria, petanalysis.AnalysisNameDirofilaria)
// 		processGroup(d.Analyses.Ehrlichiosis, petanalysis.AnalysisNameEhrlichiosis)
// 		processGroup(d.Analyses.Anaplasmosis, petanalysis.AnalysisNameAnaplasmosis)

// 		p.Edges.Analyses = analyses
// 	}

// 	if d.Bonuses != nil {
// 		p.Edges.Bonuses = &ent.PetBonus{
// 			IsArtist:      d.Bonuses.IsArtist,
// 			IsTherapist:   d.Bonuses.IsTherapist,
// 			IsFormerDonor: d.Bonuses.IsFormerDonor,
// 			IsGuideDog:    d.Bonuses.IsGuideDog,
// 		}
// 	}

// 	return p
// }

// mapDTOToPetUpdates преобразует DTO обновления питомца в аргументы для сервиса
