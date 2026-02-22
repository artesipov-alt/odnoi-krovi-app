package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/dto"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/services"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/validator"
	"github.com/danielgtaylor/huma/v2"
)

// PetHandler обрабатывает HTTP запросы для операций с питомцами
type PetHandler struct {
	petService    services.PetService
	validator     validator.DonorValidator
	bloodInfoRepo services.BloodInfoRepository
}

// NewPetHandler создает новый обработчик питомцев
func NewPetHandler(petService services.PetService, validator validator.DonorValidator, bloodInfoRepo services.BloodInfoRepository) *PetHandler {
	return &PetHandler{
		petService:    petService,
		validator:     validator,
		bloodInfoRepo: bloodInfoRepo,
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

	// // Получение питомца по ID
	// huma.Register(api, huma.Operation{
	// 	OperationID: "get-pet-by-id",
	// 	Method:      http.MethodGet,
	// 	Path:        "/v1/pet/{id}",
	// 	Summary:     "Получение питомца по ID",
	// 	Description: "Возвращает информацию о питомце по его идентификатору",
	// 	Tags:        []string{"pets-v1"},
	// }, h.GetPet)

	// // Получение питомцев пользователя
	// huma.Register(api, huma.Operation{
	// 	OperationID: "get-user-pets",
	// 	Method:      http.MethodGet,
	// 	Path:        "/v1/pet/user/{user_id}",
	// 	Summary:     "Получение питомцев пользователя",
	// 	Description: "Возвращает всех питомцев конкретного пользователя",
	// 	Tags:        []string{"pets-v1"},
	// }, h.GetUserPets)

	// // Обновление данных питомца
	// huma.Register(api, huma.Operation{
	// 	OperationID: "update-pet",
	// 	Method:      http.MethodPut,
	// 	Path:        "/v1/pet/{id}",
	// 	Summary:     "Обновление данных питомца",
	// 	Description: "Обновляет информацию о питомце",
	// 	Tags:        []string{"pets-v1"},
	// }, h.UpdatePet)

	// // Удаление питомца по ID
	// huma.Register(api, huma.Operation{
	// 	OperationID: "delete-pet",
	// 	Method:      http.MethodDelete,
	// 	Path:        "/v1/pet/{id}",
	// 	Summary:     "Удаление питомца по ID",
	// 	Description: "Удаляет питомца из системы",
	// 	Tags:        []string{"pets-v1"},
	// }, h.DeletePet)

	// // Валидация донора по ID (изменено на POST)
	// huma.Register(api, huma.Operation{
	// 	OperationID:   "validate-donor",
	// 	Method:        http.MethodPost,
	// 	Path:          "/v1/pet/validate-donor/{id}",
	// 	Summary:       "Валидация донора по ID",
	// 	Description:   "Пересчитывает и сохраняет факторы валидации донора для питомца",
	// 	Tags:          []string{"pets-v1"},
	// 	DefaultStatus: http.StatusOK,
	// }, h.ValidateDonor)

}

//==========Handlers==============================

func (h *PetHandler) CreatePet(ctx context.Context, input *struct {
	dto.PetUserIDPath
	Body dto.PetCreate
}) (*dto.PetResponse, error) {
	body := &input.Body

	petDomain := new(domain.Pet)
	ToDomain(*body, petDomain)

	createdPet, err := h.petService.CreatePet(ctx, input.ID, petDomain)
	if err != nil {
		return nil, err
	}

	return &dto.PetResponse{Body: ToDTO(*createdPet)}, nil
}

// func (h *PetHandler) UpdatePet(ctx context.Context, input *struct {
// 	dto.IDPathStr
// 	Body dto.PetUpdate
// }) (*dto.MessageResponse, error) {

// 	body := input.Body
// 	body.BirthDate = calculateBirthDateFromAge(body.AgeYears, body.AgeMonths)

// 	// Копируем в UpdatePetInput
// 	var petInput ent.UpdatePetInput
// 	if err := copier.Copy(&petInput, &body); err != nil {
// 		return nil, apperrors.Internal(err, "failed to copy pet update data")
// 	}

// 	// Устанавливаем BreedRefID, так как copier не копирует поле с другим именем
// 	if body.BreedID != nil {
// 		petInput.BreedRefID = body.BreedID
// 	}

// 	// Устанавливаем BloodGroupRefID, так как copier не копирует поле с другим именем
// 	if body.BloodGroup != nil && *body.BloodGroup != "" {
// 		bg, err := h.bloodInfoRepo.FindByBloodGroup(ctx, *body.BloodGroup)
// 		if err != nil {
// 			return nil, apperrors.Internal(err, "failed to find blood group")
// 		}
// 		petInput.BloodGroupRefID = &bg.ID
// 	}

// 	// Копируем health
// 	var healthInput *ent.UpdatePetHealthInput
// 	if body.Health != nil {
// 		healthInput = &ent.UpdatePetHealthInput{}
// 		if body.Health.HealthStatus != nil {
// 			status := pethealth.HealthStatus(*body.Health.HealthStatus)
// 			healthInput.HealthStatus = &status
// 		}
// 		healthInput.LastDonation = body.Health.LastDonation
// 		healthInput.Transfused = body.Health.Transfused
// 		healthInput.Medications = body.Health.Medications
// 		healthInput.SurgicalInterventions = body.Health.SurgicalInterventions
// 	}

// 	// Копируем treatments
// 	var treatmentsInput *ent.UpdatePetTreatmentInput
// 	if body.Treatments != nil {
// 		treatmentsInput = &ent.UpdatePetTreatmentInput{}
// 		if err := copier.Copy(treatmentsInput, body.Treatments); err != nil {
// 			return nil, apperrors.Internal(err, "failed to copy treatments update data")
// 		}
// 	}

// 	// Копируем analyses
// 	var analysesInput []*ent.UpdatePetAnalysisInput
// 	if body.Analyses != nil {
// 		analysesInput = []*ent.UpdatePetAnalysisInput{}
// 		processGroup := func(group []*dto.PetAnalysis, name petanalysis.AnalysisName) {
// 			for _, a := range group {
// 				var ai ent.UpdatePetAnalysisInput
// 				ai.AnalysisName = &name
// 				if a.AnalysisDate != nil {
// 					ai.AnalysisDate = a.AnalysisDate
// 				}
// 				if a.AnalysisType != nil {
// 					at := petanalysis.AnalysisType(*a.AnalysisType)
// 					ai.AnalysisType = &at
// 				}
// 				analysesInput = append(analysesInput, &ai)
// 			}
// 		}
// 		processGroup(body.Analyses.Leukemia, petanalysis.AnalysisNameLeukemia)
// 		processGroup(body.Analyses.Immunodeficiency, petanalysis.AnalysisNameImmunodeficiency)
// 		processGroup(body.Analyses.Hemoplasmosis, petanalysis.AnalysisNameHemoplasmosis)
// 		processGroup(body.Analyses.Bartonellosis, petanalysis.AnalysisNameBartonellosis)
// 		processGroup(body.Analyses.Babesiosis, petanalysis.AnalysisNameBabesiosis)
// 		processGroup(body.Analyses.Dirofilaria, petanalysis.AnalysisNameDirofilaria)
// 		processGroup(body.Analyses.Ehrlichiosis, petanalysis.AnalysisNameEhrlichiosis)
// 		processGroup(body.Analyses.Anaplasmosis, petanalysis.AnalysisNameAnaplasmosis)
// 	}

// 	// Копируем bonuses
// 	var bonusesInput *ent.UpdatePetBonusInput
// 	if body.Bonuses != nil {
// 		bonusesInput = &ent.UpdatePetBonusInput{}
// 		if err := copier.Copy(bonusesInput, body.Bonuses); err != nil {
// 			return nil, apperrors.Internal(err, "failed to copy bonuses update data")
// 		}
// 	}

// 	if err := h.petService.Update(ctx, input.ID, &petInput, healthInput, treatmentsInput, analysesInput, bonusesInput); err != nil {
// 		return nil, err
// 	}

// 	return &dto.MessageResponse{Body: dto.MessageBody{Message: "Питомец обновлен"}}, nil
// }

// func (h *PetHandler) GetPet(ctx context.Context,
// 	input *struct {
// 		dto.IDPathStr
// 		dto.PetPreloadQuery
// 	}) (*dto.PetResponse, error) {
// 	// Добавляем опции к запросу.
// 	opts := services.PetPreloadOptions{
// 		WithHealth:     input.WithHealth,
// 		WithTreatments: input.WithTreatments,
// 		WithAnalyses:   input.WithAnalysis,
// 		WithBonuses:    input.WithBonuses,
// 		WithAll:        input.WithAll,
// 	}

// 	pet, err := h.petService.GetPet(ctx, input.ID, opts)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &dto.PetResponse{Body: h.toDTO(pet)}, nil
// }

// func (h *PetHandler) GetUserPets(ctx context.Context,
// 	input *struct {
// 		dto.PetUserIDPath
// 		dto.PetPreloadQuery
// 	}) (*dto.PetsResponse, error) {

// 	opts := services.PetPreloadOptions{
// 		WithHealth:     input.WithHealth,
// 		WithTreatments: input.WithTreatments,
// 		WithAnalyses:   input.WithAnalysis,
// 		WithBonuses:    input.WithBonuses,
// 		WithAll:        input.WithAll,
// 	}

// 	pets, err := h.petService.GetUserPets(ctx, input.ID, opts)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &dto.PetsResponse{Body: h.toPetsDTO(pets)}, nil
// }

// func (h *PetHandler) DeletePet(ctx context.Context, input *dto.IDPathStr) (*dto.MessageResponse, error) {
// 	if err := h.petService.DeletePet(ctx, input.ID); err != nil {
// 		return nil, err
// 	}

// 	resp := &dto.MessageResponse{}
// 	resp.Body.Message = "Питомец удален"
// 	return resp, nil
// }

// func (h *PetHandler) ValidateDonor(ctx context.Context, input *dto.IDPathStr) (*dto.PetResponse, error) {
// 	// Загружаем все связанные данные для полной валидации
// 	p, err := h.petService.GetPet(ctx, input.ID, services.PetPreloadOptions{WithAll: true})
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Применяем валидацию: пересчитываем и сохраняем факторы
// 	_, _, err = h.petService.ApplyValidation(ctx, p.ID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Возвращаем обновлённый DTO (факторы теперь сохранены и будут включены)
// 	return &dto.PetResponse{Body: h.toDTO(p)}, nil
// }

// mapPetToDTO преобразует ENT модель питомца в DTO для ответа
// func (h *PetHandler) toDTO(p *ent.Pet) dto.Pet {
// 	petDTO := dto.Pet{
// 		ID:                 p.ID,
// 		Name:               p.Name,
// 		ChipNumber:         p.ChipNumber,
// 		PhotoURLs:          p.PhotoUrls,
// 		BreedID:            p.BreedID,
// 		WeightKg:           p.WeightKg,
// 		BirthDate:          p.BirthDate,
// 		LivingCondition:    string(p.LivingCondition),
// 		Gender:             string(p.Gender),
// 		Type:               string(p.Type),
// 		BloodGroup:         "",
// 		ReproductiveStatus: string(p.ReproductiveStatus),
// 		// PetStatus:          string(p.PetStatus),
// 		CreatedAt: &p.CreatedAt,
// 		UpdatedAt: &p.UpdatedAt,
// 		DeletedAt: p.DeletedAt,
// 	}
// 	if p.Edges.BloodGroupRef != nil {
// 		petDTO.BloodGroup = p.Edges.BloodGroupRef.BloodGroup
// 	}
// 	if p.Edges.Health != nil {
// 		healthStatus := string(p.Edges.Health.HealthStatus)
// 		medications := p.Edges.Health.Medications
// 		surgical := p.Edges.Health.SurgicalInterventions
// 		petDTO.Health = &dto.PetHealth{
// 			HealthStatus:          &healthStatus,
// 			LastDonation:          p.Edges.Health.LastDonation,
// 			Transfused:            &p.Edges.Health.Transfused,
// 			Medications:           &medications,
// 			SurgicalInterventions: &surgical,
// 		}
// 	}
// 	if p.Edges.Treatments != nil {
// 		petDTO.Treatments = &dto.PetTreatment{
// 			RabiesVaccinationDate:     p.Edges.Treatments.RabiesVaccinationDate,
// 			InfectionVaccinationDate:  p.Edges.Treatments.InfectionVaccinationDate,
// 			EctoparasiteTreatmentDate: p.Edges.Treatments.EctoparasiteTreatmentDate,
// 			DewormingDate:             p.Edges.Treatments.DewormingDate,
// 		}
// 	}
// 	if p.Edges.Analyses != nil {
// 		petDTO.Analyses = &dto.PetAnalysisGroup{}
// 		for _, a := range p.Edges.Analyses {
// 			analysisName := string(a.AnalysisName)
// 			analysisType := string(a.AnalysisType)
// 			dtoAnalysis := &dto.PetAnalysis{
// 				ID:           &a.ID,
// 				AnalysisName: &analysisName,
// 				AnalysisType: &analysisType,
// 				AnalysisDate: a.AnalysisDate,
// 			}

// 			switch a.AnalysisName {
// 			case petanalysis.AnalysisNameLeukemia:
// 				petDTO.Analyses.Leukemia = append(petDTO.Analyses.Leukemia, dtoAnalysis)
// 			case petanalysis.AnalysisNameImmunodeficiency:
// 				petDTO.Analyses.Immunodeficiency = append(petDTO.Analyses.Immunodeficiency, dtoAnalysis)
// 			case petanalysis.AnalysisNameHemoplasmosis:
// 				petDTO.Analyses.Hemoplasmosis = append(petDTO.Analyses.Hemoplasmosis, dtoAnalysis)
// 			case petanalysis.AnalysisNameBartonellosis:
// 				petDTO.Analyses.Bartonellosis = append(petDTO.Analyses.Bartonellosis, dtoAnalysis)
// 			case petanalysis.AnalysisNameBabesiosis:
// 				petDTO.Analyses.Babesiosis = append(petDTO.Analyses.Babesiosis, dtoAnalysis)
// 			case petanalysis.AnalysisNameDirofilaria:
// 				petDTO.Analyses.Dirofilaria = append(petDTO.Analyses.Dirofilaria, dtoAnalysis)
// 			case petanalysis.AnalysisNameEhrlichiosis:
// 				petDTO.Analyses.Ehrlichiosis = append(petDTO.Analyses.Ehrlichiosis, dtoAnalysis)
// 			case petanalysis.AnalysisNameAnaplasmosis:
// 				petDTO.Analyses.Anaplasmosis = append(petDTO.Analyses.Anaplasmosis, dtoAnalysis)
// 			}
// 		}
// 	}
// 	if p.Edges.Bonuses != nil {
// 		petDTO.Bonuses = &dto.PetBonus{
// 			IsArtist:      p.Edges.Bonuses.IsArtist,
// 			IsTherapist:   p.Edges.Bonuses.IsTherapist,
// 			IsFormerDonor: p.Edges.Bonuses.IsFormerDonor,
// 			IsGuideDog:    p.Edges.Bonuses.IsGuideDog,
// 		}
// 	}

// 	// Заполняем DonorRestrictions из сохранённых данных
// 	// if p.PetStatus != "" && len(p.DonorRestrictions) > 0 {
// 	// 	var stopRestrictionFactors []dto.RestrictionFactor
// 	// 	var warnRestrictionFactors []dto.RestrictionFactor

// 	// 	for _, codeStr := range p.DonorRestrictions {
// 	// 		code := validator.FactorCode(codeStr)
// 	// 		desc := validator.GetFactorDescription(code)

// 	// 		factor := dto.RestrictionFactor{
// 	// 			Code:           codeStr,
// 	// 			Description:    desc.Description,
// 	// 			SubDescription: desc.SubDescription,
// 	// 		}

// 	// 		// Разделяем на стопы и предупреждения по префиксу
// 	// 		if strings.HasPrefix(codeStr, "STOP_") {
// 	// 			stopRestrictionFactors = append(stopRestrictionFactors, factor)
// 	// 		} else if strings.HasPrefix(codeStr, "WARN_") {
// 	// 			warnRestrictionFactors = append(warnRestrictionFactors, factor)
// 	// 		}
// 	// 	}

// 	// 	petDTO.DonorRestrictions = &dto.DonorRestrictions{
// 	// 		StopFactors: stopRestrictionFactors,
// 	// 		WarnFactors: warnRestrictionFactors,
// 	// 	}
// 	// } else if p.PetStatus != "" {
// 	// 	// Если факторы ещё не рассчитаны, возвращаем пустые (не рассчитываем на лету)
// 	// 	petDTO.DonorRestrictions = &dto.DonorRestrictions{
// 	// 		StopFactors: []dto.RestrictionFactor{},
// 	// 		WarnFactors: []dto.RestrictionFactor{},
// 	// 	}
// 	// }

// 	return petDTO
// }

// toPetsDTO преобразует слайс ENT питомцев в слайс DTO
// func (h *PetHandler) toPetsDTO(pets []*ent.Pet) []dto.Pet {
// 	petDTOs := make([]dto.Pet, len(pets))
// 	for i, p := range pets {
// 		petDTOs[i] = h.toDTO(p)
// 	}
// 	return petDTOs
// }

// calculateBirthDateFromAge вычисляет дату рождения на основе возраста в годах и месяцах
func calculateBirthDateFromAge(ageYears, ageMonths *int) *time.Time {
	if ageYears == nil && ageMonths == nil {
		return nil
	}

	years := 0
	if ageYears != nil {
		years = *ageYears
	}

	months := 0
	if ageMonths != nil {
		months = *ageMonths
	}

	if years > 0 || months > 0 {
		birthDate := time.Now().AddDate(-years, -months, 0)
		birthDate = time.Date(birthDate.Year(), birthDate.Month(), 1, 0, 0, 0, 0, time.UTC)
		return &birthDate
	}
	return nil
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

func ToDomain(dto dto.PetCreate, petDomain *domain.Pet) {
	petDomain.Name = dto.Name
	petDomain.Type = domain.PetType(dto.Type)
	petDomain.WeightKg = dto.WeightKg
	petDomain.Gender = domain.Gender(dto.Gender)
	petDomain.ChipNumber = dto.ChipNumber
	petDomain.LivingCondition = domain.LivingCondition(dto.LivingCondition)
	petDomain.BreedRefID = dto.BreedID
	petDomain.BloodGroupRefID = dto.BloodGroup

	if dto.BirthDate != nil {
		petDomain.BirthDate = dto.BirthDate
	}

	if dto.ReproductiveStatus != "" {
		petDomain.ReproductiveStatus = domain.ReproductiveStatus(dto.ReproductiveStatus)
	}

	// Handle PetHealth
	if dto.Health != nil {
		healthDomain := &domain.PetHealth{}
		healthDomain.HealthStatus = domain.HealthStatus(*dto.Health.HealthStatus)
		healthDomain.LastDonation = *dto.Health.LastDonation
		healthDomain.Transfused = *dto.Health.Transfused
		healthDomain.Medications = *dto.Health.Medications
		healthDomain.SurgicalInterventions = *dto.Health.SurgicalInterventions
		petDomain.Health = healthDomain
	}

	// Handle PetTreatment
	if dto.Treatments != nil {
		treatmentDomain := &domain.PetTreatment{}
		treatmentDomain.RabiesVaccinationDate = *dto.Treatments.RabiesVaccinationDate
		treatmentDomain.InfectionVaccinationDate = *dto.Treatments.InfectionVaccinationDate
		treatmentDomain.EctoparasiteTreatmentDate = *dto.Treatments.EctoparasiteTreatmentDate
		treatmentDomain.DewormingDate = *dto.Treatments.DewormingDate
		petDomain.Treatments = treatmentDomain
	}

	// Handle PetAnalysis
	// var analysesDomain []*domain.PetAnalysis
	// if dto.Analyses != nil {
	// 	processGroup := func(group []*dto.PetAnalysis, name string) {
	// 		for _, a := range group {
	// 			if a != nil && a.AnalysisDate != nil {
	// 				analysis := &domain.PetAnalysis{
	// 					AnalysisName: name,
	// 					AnalysisDate: *a.AnalysisDate,
	// 				}
	// 				if a.AnalysisType != nil {
	// 					analysis.AnalysisType = *a.AnalysisType
	// 				}
	// 				analysesDomain = append(analysesDomain, analysis)
	// 			}
	// 		}
	// 	}

	// 	processGroup(dto.Analyses.Leukemia, "leukemia")
	// 	processGroup(dto.Analyses.Immunodeficiency, "immunodeficiency")
	// 	processGroup(dto.Analyses.Hemoplasmosis, "hemoplasmosis")
	// 	processGroup(dto.Analyses.Bartonellosis, "bartonellosis")
	// 	processGroup(dto.Analyses.Babesiosis, "babesiosis")
	// 	processGroup(dto.Analyses.Dirofilaria, "dirofilaria")
	// 	processGroup(dto.Analyses.Ehrlichiosis, "ehrlichiosis")
	// 	processGroup(dto.Analyses.Anaplasmosis, "anaplasmosis")
	// 	petDomain.Analyses = analysesDomain
	// }

}

// ToDTO maps domain.Pet to dto.Pet
func ToDTO(petDomain domain.Pet) dto.Pet {
	petDTO := dto.Pet{
		ID:                 petDomain.ID,
		Name:               petDomain.Name,
		ChipNumber:         petDomain.ChipNumber,
		PhotoURLs:          petDomain.PhotoURLs,
		BreedID:            petDomain.BreedRefID,
		WeightKg:           petDomain.WeightKg,
		BirthDate:          petDomain.BirthDate,
		LivingCondition:    string(petDomain.LivingCondition),
		Gender:             string(petDomain.Gender),
		Type:               string(petDomain.Type),
		BloodGroup:         petDomain.BloodGroupRefID,
		ReproductiveStatus: string(petDomain.ReproductiveStatus),
		CreatedAt:          petDomain.CreatedAt,
		UpdatedAt:          petDomain.UpdatedAt,
		DeletedAt:          petDomain.DeletedAt,
	}

	if petDomain.Health != nil {
		healthStatus := string(petDomain.Health.HealthStatus)
		petDTO.Health = &dto.PetHealth{
			HealthStatus:          &healthStatus,
			LastDonation:          nilable(petDomain.Health.LastDonation),
			Transfused:            nilable(petDomain.Health.Transfused),
			Medications:           nilable(petDomain.Health.Medications),
			SurgicalInterventions: nilable(petDomain.Health.SurgicalInterventions),
		}
	}

	if petDomain.Treatments != nil {
		petDTO.Treatments = &dto.PetTreatment{
			RabiesVaccinationDate:     nilable(petDomain.Treatments.RabiesVaccinationDate),
			InfectionVaccinationDate:  nilable(petDomain.Treatments.InfectionVaccinationDate),
			EctoparasiteTreatmentDate: nilable(petDomain.Treatments.EctoparasiteTreatmentDate),
			DewormingDate:             nilable(petDomain.Treatments.DewormingDate),
		}
	}

	// Analyses mapping if needed, but commented in original ToDomain
	// if petDomain.Analyses != nil {
	// 	petDTO.Analyses = &dto.PetAnalysisGroup{}
	// 	for _, a := range petDomain.Analyses {
	// 		analysisName := string(a.AnalysisName)
	// 		analysisType := string(a.AnalysisType)
	// 		dtoAnalysis := &dto.PetAnalysis{
	// 			ID:           nilable(a.ID),
	// 			AnalysisName: &analysisName,
	// 			AnalysisType: &analysisType,
	// 			AnalysisDate: nilable(a.AnalysisDate),
	// 		}
	// 		switch a.AnalysisName {
	// 		case "leukemia":
	// 			petDTO.Analyses.Leukemia = append(petDTO.Analyses.Leukemia, dtoAnalysis)
	// 		// Add other cases similarly
	// 		}
	// 	}
	// }

	return petDTO
}

func nilable[T any](val T) *T {
	// Проверяем, является ли значение "нулевым" для своего типа
	// Это работает для большинства встроенных типов (числа, строки, булевы)
	// и для структур, если все их поля нулевые.
	// Для более сложных случаев (например, пользовательские типы с неэкспортированными полями)
	// может потребоваться более сложная логика или использование reflect.
	var zero T
	if any(val) == any(zero) {
		return nil
	}
	return &val
}
