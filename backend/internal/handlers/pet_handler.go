package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/dto"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/services"
	"github.com/danielgtaylor/huma/v2"
)

// PetHandler обрабатывает HTTP запросы для операций с питомцами
type PetHandler struct {
	petService    services.PetService
	bloodInfoRepo services.BloodInfoRepository
}

// NewPetHandler создает новый обработчик питомцев
func NewPetHandler(petService services.PetService, bloodInfoRepo services.BloodInfoRepository) *PetHandler {
	return &PetHandler{
		petService:    petService,
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
}) (*dto.BodyPetCreateResponse, error) {
	body := &input.Body

	petDomain := new(domain.Pet)
	ToDomain(*body, petDomain)

	createdPet, err := h.petService.CreatePet(ctx, input.ID, petDomain)
	if err != nil {
		return nil, err
	}

	return &dto.BodyPetCreateResponse{Body: dto.PetCreateResponse{ID: createdPet.ID, CreatedAt: createdPet.CreatedAt}}, nil
}

func (h *PetHandler) UpdatePet(ctx context.Context, input *struct {
	dto.IDPathStr
	Body dto.PetUpdate
}) (*dto.PetResponse, error) {
	body := &input.Body

	petDomain := new(domain.Pet)
	ToDomain(body, petDomain)

	updatedPet, err := h.petService.Update(ctx, input.ID, petDomain)
	if err != nil {
		return nil, err
	}

	return &dto.PetResponse{Body: ToDTO(*updatedPet)}, nil
}

func (h *PetHandler) GetPet(ctx context.Context,
	input *struct {
		dto.IDPathStr
		dto.PetPreloadQuery
	}) (*dto.PetResponse, error) {
	// Добавляем опции к запросу.
	opts := services.PetPreloadOptions{
		WithHealth:     input.WithHealth,
		WithTreatments: input.WithTreatments,
		WithAnalyses:   input.WithAnalysis,
		WithBonuses:    input.WithBonuses,
		WithAll:        input.WithAll,
	}

	pet, err := h.petService.GetPet(ctx, input.ID, opts)
	if err != nil {
		return nil, err
	}

	return &dto.PetResponse{Body: ToDTO(*pet)}, nil
}

func (h *PetHandler) GetUserPets(ctx context.Context,
	input *struct {
		dto.PetUserIDPath
		dto.PetPreloadQuery
	}) (*dto.PetsResponse, error) {

	opts := services.PetPreloadOptions{
		WithHealth:     input.WithHealth,
		WithTreatments: input.WithTreatments,
		WithAnalyses:   input.WithAnalysis,
		WithBonuses:    input.WithBonuses,
		WithAll:        input.WithAll,
	}

	pets, err := h.petService.GetUserPets(ctx, input.ID, opts)
	if err != nil {
		return nil, err
	}

	return &dto.PetsResponse{Body: ToDTOs(pets)}, nil
}

func (h *PetHandler) DeletePet(ctx context.Context, input *dto.IDPathStr) (*dto.MessageResponse, error) {
	if err := h.petService.DeletePet(ctx, input.ID); err != nil {
		return nil, err
	}

	resp := &dto.MessageResponse{}
	resp.Body.Message = "Питомец удален"
	return resp, nil
}

func (h *PetHandler) ValidateDonor(ctx context.Context, input *dto.IDPathStr) (*dto.PetResponse, error) {
	pet, err := h.petService.RevalidateDonor(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	return &dto.PetResponse{Body: ToDTO(*pet)}, nil
}

// toPetsDTO преобразует слайс ENT питомцев в слайс DTO
func ToDTOs(pets []*domain.Pet) []dto.Pet {
	petDTOs := make([]dto.Pet, len(pets))
	for i, p := range pets {
		petDTOs[i] = ToDTO(*p)
	}
	return petDTOs
}

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

func ToDomain(petDto any, petDomain *domain.Pet) {
	switch v := petDto.(type) {
	case dto.PetCreate:
		petDomain.Name = v.Name
		petDomain.Type = domain.PetType(v.Type)
		petDomain.WeightKg = v.WeightKg
		petDomain.Gender = domain.Gender(v.Gender)
		petDomain.ChipNumber = v.ChipNumber
		petDomain.LivingCondition = domain.LivingCondition(v.LivingCondition)
		if v.BreedID != "" {
			petDomain.BreedRefID = &v.BreedID
		}
		if v.BloodGroup != "" {
			petDomain.BloodGroupName = &v.BloodGroup
		}

		if v.BirthDate != nil {
			petDomain.BirthDate = v.BirthDate
		}

		if v.ReproductiveStatus != "" {
			petDomain.ReproductiveStatus = domain.ReproductiveStatus(v.ReproductiveStatus)
		}
		if v.AgeMonths != 0 || v.AgeYears != 0 {
			petDomain.BirthDate = calculateBirthDateFromAge(&v.AgeYears, &v.AgeMonths)
		}

		// Handle PetHealth
		if v.Health != nil {
			healthDomain := &domain.PetHealth{}
			if v.Health.HealthStatus != nil {
				healthDomain.HealthStatus = domain.HealthStatus(*v.Health.HealthStatus)
			}
			if v.Health.LastDonation != nil {
				healthDomain.LastDonation = v.Health.LastDonation
			}
			if v.Health.Transfused != nil {
				healthDomain.Transfused = v.Health.Transfused
			}
			if v.Health.Medications != nil {
				healthDomain.Medications = v.Health.Medications
			}
			if v.Health.SurgicalInterventions != nil {
				healthDomain.SurgicalInterventions = v.Health.SurgicalInterventions
			}
			petDomain.Health = healthDomain
		}

		// Handle PetTreatment
		if v.Treatments != nil {
			treatmentDomain := &domain.PetTreatment{}
			if v.Treatments.RabiesVaccinationDate != nil {
				treatmentDomain.RabiesVaccinationDate = v.Treatments.RabiesVaccinationDate
			}
			if v.Treatments.InfectionVaccinationDate != nil {
				treatmentDomain.InfectionVaccinationDate = v.Treatments.InfectionVaccinationDate
			}
			if v.Treatments.EctoparasiteTreatmentDate != nil {
				treatmentDomain.EctoparasiteTreatmentDate = v.Treatments.EctoparasiteTreatmentDate
			}
			if v.Treatments.DewormingDate != nil {
				treatmentDomain.DewormingDate = v.Treatments.DewormingDate
			}
			petDomain.Treatments = treatmentDomain
		}

		// Handle PetAnalysis
		analysesDomain := []*domain.PetAnalysis{}
		if v.Analyses != nil {
			processGroup := func(group []*dto.PetAnalysis, name string) {
				for _, a := range group {
					if a != nil && a.AnalysisDate != nil {
						analysis := &domain.PetAnalysis{
							AnalysisName: name,
							AnalysisDate: a.AnalysisDate,
						}
						if a.AnalysisType != nil {
							analysis.AnalysisType = *a.AnalysisType
						}
						analysesDomain = append(analysesDomain, analysis)
					}
				}
			}

			processGroup(v.Analyses.Leukemia, "leukemia")
			processGroup(v.Analyses.Immunodeficiency, "immunodeficiency")
			processGroup(v.Analyses.Hemoplasmosis, "hemoplasmosis")
			processGroup(v.Analyses.Bartonellosis, "bartonellosis")
			processGroup(v.Analyses.Babesiosis, "babesiosis")
			processGroup(v.Analyses.Dirofilaria, "dirofilaria")
			processGroup(v.Analyses.Ehrlichiosis, "ehrlichiosis")
			processGroup(v.Analyses.Anaplasmosis, "anaplasmosis")
			petDomain.Analyses = analysesDomain
		}
	case *dto.PetUpdate:
		if v.Name != nil {
			petDomain.Name = *v.Name
		}
		if v.Type != nil {
			petDomain.Type = domain.PetType(*v.Type)
		}
		if v.WeightKg != nil {
			petDomain.WeightKg = *v.WeightKg
		}
		if v.Gender != nil {
			petDomain.Gender = domain.Gender(*v.Gender)
		}
		if v.ChipNumber != nil {
			petDomain.ChipNumber = *v.ChipNumber
		}
		if v.LivingCondition != nil {
			petDomain.LivingCondition = domain.LivingCondition(*v.LivingCondition)
		}
		if v.BreedID != nil {
			petDomain.BreedRefID = v.BreedID
		}
		if v.BloodGroup != nil {
			petDomain.BloodGroupName = v.BloodGroup
		}
		if v.BirthDate != nil {
			petDomain.BirthDate = v.BirthDate
		}
		if v.ReproductiveStatus != nil {
			petDomain.ReproductiveStatus = domain.ReproductiveStatus(*v.ReproductiveStatus)
		}
		if v.Bonuses != nil {
			petDomain.Bonuses = *v.Bonuses
		}
		if v.AgeMonths != nil || v.AgeYears != nil {
			petDomain.BirthDate = calculateBirthDateFromAge(v.AgeYears, v.AgeMonths)
		}

		// Handle PetHealth
		if v.Health != nil {
			healthDomain := &domain.PetHealth{}
			if v.Health.HealthStatus != nil {
				healthDomain.HealthStatus = domain.HealthStatus(*v.Health.HealthStatus)
			}
			if v.Health.LastDonation != nil {
				healthDomain.LastDonation = v.Health.LastDonation
			}
			if v.Health.Transfused != nil {
				healthDomain.Transfused = v.Health.Transfused
			}
			if v.Health.Medications != nil {
				healthDomain.Medications = v.Health.Medications
			}
			if v.Health.SurgicalInterventions != nil {
				healthDomain.SurgicalInterventions = v.Health.SurgicalInterventions
			}
			petDomain.Health = healthDomain
		}

		// Handle PetTreatment
		if v.Treatments != nil {
			treatmentDomain := &domain.PetTreatment{}
			if v.Treatments.RabiesVaccinationDate != nil {
				treatmentDomain.RabiesVaccinationDate = v.Treatments.RabiesVaccinationDate
			}
			if v.Treatments.InfectionVaccinationDate != nil {
				treatmentDomain.InfectionVaccinationDate = v.Treatments.InfectionVaccinationDate
			}
			if v.Treatments.EctoparasiteTreatmentDate != nil {
				treatmentDomain.EctoparasiteTreatmentDate = v.Treatments.EctoparasiteTreatmentDate
			}
			if v.Treatments.DewormingDate != nil {
				treatmentDomain.DewormingDate = v.Treatments.DewormingDate
			}
			petDomain.Treatments = treatmentDomain
		}

		// Handle PetAnalysis
		analysesDomain := []*domain.PetAnalysis{}
		if v.Analyses != nil {
			processGroup := func(group []*dto.PetAnalysis, name string) {
				for _, a := range group {
					if a != nil && a.AnalysisDate != nil {
						analysis := &domain.PetAnalysis{
							AnalysisName: name,
							AnalysisDate: a.AnalysisDate,
						}
						if a.AnalysisType != nil {
							analysis.AnalysisType = *a.AnalysisType
						}
						analysesDomain = append(analysesDomain, analysis)
					}
				}
			}

			processGroup(v.Analyses.Leukemia, "leukemia")
			processGroup(v.Analyses.Immunodeficiency, "immunodeficiency")
			processGroup(v.Analyses.Hemoplasmosis, "hemoplasmosis")
			processGroup(v.Analyses.Bartonellosis, "bartonellosis")
			processGroup(v.Analyses.Babesiosis, "babesiosis")
			processGroup(v.Analyses.Dirofilaria, "dirofilaria")
			processGroup(v.Analyses.Ehrlichiosis, "ehrlichiosis")
			processGroup(v.Analyses.Anaplasmosis, "anaplasmosis")
			petDomain.Analyses = analysesDomain
		}
	}
}

// ToDTO maps domain.Pet to dto.Pet
func ToDTO(petDomain domain.Pet) dto.Pet {
	petDTO := dto.Pet{
		ID:                 petDomain.ID,
		Name:               petDomain.Name,
		ChipNumber:         petDomain.ChipNumber,
		PhotoURLs:          petDomain.PhotoURLs,
		WeightKg:           petDomain.WeightKg,
		BirthDate:          petDomain.BirthDate,
		PetStatus:          string(petDomain.PetStatus),
		LivingCondition:    string(petDomain.LivingCondition),
		Gender:             string(petDomain.Gender),
		Type:               string(petDomain.Type),
		ReproductiveStatus: string(petDomain.ReproductiveStatus),
		Bonuses:            petDomain.Bonuses,
		CreatedAt:          petDomain.CreatedAt,
		UpdatedAt:          petDomain.UpdatedAt,
		DeletedAt:          petDomain.DeletedAt,
	}

	// Map StopFactors and WarnFactors into DonorRestrictions
	if len(petDomain.StopFactors) > 0 || len(petDomain.WarnFactors) > 0 {
		var stopFactors []dto.RestrictionFactor
		var warnFactors []dto.RestrictionFactor
		for _, code := range petDomain.StopFactors {
			desc := domain.GetFactorDescription(domain.FactorCode(code))
			factor := dto.RestrictionFactor{
				Code:           code,
				Description:    desc.Description,
				SubDescription: desc.SubDescription,
			}
			stopFactors = append(stopFactors, factor)
		}
		for _, code := range petDomain.WarnFactors {
			desc := domain.GetFactorDescription(domain.FactorCode(code))
			factor := dto.RestrictionFactor{
				Code:           code,
				Description:    desc.Description,
				SubDescription: desc.SubDescription,
			}
			warnFactors = append(warnFactors, factor)
		}
		petDTO.DonorRestrictions = &dto.DonorRestrictions{
			StopFactors: stopFactors,
			WarnFactors: warnFactors,
		}
	}

	if petDomain.BreedRefID != nil {
		petDTO.BreedID = *petDomain.BreedRefID
	}
	if petDomain.BloodGroupName != nil {
		petDTO.BloodGroup = *petDomain.BloodGroupName
	}

	if petDomain.Health != nil {
		healthStatus := string(petDomain.Health.HealthStatus)
		petDTO.Health = &dto.PetHealth{
			HealthStatus:          &healthStatus,
			LastDonation:          petDomain.Health.LastDonation,
			Transfused:            petDomain.Health.Transfused,
			Medications:           petDomain.Health.Medications,
			SurgicalInterventions: petDomain.Health.SurgicalInterventions,
		}
	}

	if petDomain.Treatments != nil {
		petDTO.Treatments = &dto.PetTreatment{
			RabiesVaccinationDate:     petDomain.Treatments.RabiesVaccinationDate,
			InfectionVaccinationDate:  petDomain.Treatments.InfectionVaccinationDate,
			EctoparasiteTreatmentDate: petDomain.Treatments.EctoparasiteTreatmentDate,
			DewormingDate:             petDomain.Treatments.DewormingDate,
		}
	}

	// Analyses mapping if needed, but commented in original ToDomain
	if petDomain.Analyses != nil {
		petDTO.Analyses = &dto.PetAnalysisGroup{}
		for _, a := range petDomain.Analyses {
			analysisName := string(a.AnalysisName)
			analysisType := string(a.AnalysisType)
			dtoAnalysis := &dto.PetAnalysis{
				ID:           nilable(a.ID),
				AnalysisName: &analysisName,
				AnalysisType: &analysisType,
				AnalysisDate: a.AnalysisDate,
			}
			switch a.AnalysisName {
			case "leukemia":
				petDTO.Analyses.Leukemia = append(petDTO.Analyses.Leukemia, dtoAnalysis)
			case "immunodeficiency":
				petDTO.Analyses.Immunodeficiency = append(petDTO.Analyses.Immunodeficiency, dtoAnalysis)
			case "hemoplasmosis":
				petDTO.Analyses.Hemoplasmosis = append(petDTO.Analyses.Hemoplasmosis, dtoAnalysis)
			case "bartonellosis":
				petDTO.Analyses.Bartonellosis = append(petDTO.Analyses.Bartonellosis, dtoAnalysis)
			case "babesiosis":
				petDTO.Analyses.Babesiosis = append(petDTO.Analyses.Babesiosis, dtoAnalysis)
			case "dirofilaria":
				petDTO.Analyses.Dirofilaria = append(petDTO.Analyses.Dirofilaria, dtoAnalysis)
			case "ehrlichiosis":
				petDTO.Analyses.Ehrlichiosis = append(petDTO.Analyses.Ehrlichiosis, dtoAnalysis)
			case "anaplasmosis":
				petDTO.Analyses.Anaplasmosis = append(petDTO.Analyses.Anaplasmosis, dtoAnalysis)
			}
		}
	}

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
