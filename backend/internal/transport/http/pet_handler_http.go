package http

import (
	"context"
	"net/http"
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/reference"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/transport/http/dto"
	"github.com/danielgtaylor/huma/v2"
)

// PetHandler обрабатывает HTTP запросы для операций с питомцами
type PetHandler struct {
	petService    pet.PetService
	bloodInfoRepo reference.BloodInfoRepository
}

// NewPetHandler создает новый обработчик питомцев
func NewPetHandler(petService pet.PetService, bloodInfoRepo reference.BloodInfoRepository) *PetHandler {
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

	petDomain := new(model.Pet)
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

	petDomain := new(model.Pet)
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
	opts := pet.PetPreloadOptions{
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

	opts := pet.PetPreloadOptions{
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
func ToDTOs(pets []*model.Pet) []dto.Pet {
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

func ToDomain(petDto any, petmodel *model.Pet) {
	switch v := petDto.(type) {
	case dto.PetCreate:
		petmodel.Name = v.Name
		petmodel.Type = model.PetType(v.Type)
		petmodel.WeightKg = v.WeightKg
		petmodel.Gender = model.Gender(v.Gender)
		petmodel.ChipNumber = v.ChipNumber
		petmodel.LivingCondition = model.LivingCondition(v.LivingCondition)
		if v.BreedID != "" {
			petmodel.BreedRefID = &v.BreedID
		}
		if v.BloodGroup != "" {
			petmodel.BloodGroupName = &v.BloodGroup
		}

		if v.BirthDate != nil {
			petmodel.BirthDate = v.BirthDate
		}

		if v.ReproductiveStatus != "" {
			petmodel.ReproductiveStatus = model.ReproductiveStatus(v.ReproductiveStatus)
		}
		if v.AgeMonths != 0 || v.AgeYears != 0 {
			petmodel.BirthDate = calculateBirthDateFromAge(&v.AgeYears, &v.AgeMonths)
		}

		// Handle PetHealth
		if v.Health != nil {
			healthmodel := &model.PetHealth{}
			if v.Health.HealthStatus != nil {
				healthmodel.HealthStatus = model.HealthStatus(*v.Health.HealthStatus)
			}
			if v.Health.LastDonation != nil {
				healthmodel.LastDonation = v.Health.LastDonation
			}
			if v.Health.Transfused != nil {
				healthmodel.Transfused = v.Health.Transfused
			}
			if v.Health.Medications != nil {
				healthmodel.Medications = v.Health.Medications
			}
			if v.Health.SurgicalInterventions != nil {
				healthmodel.SurgicalInterventions = v.Health.SurgicalInterventions
			}
			petmodel.Health = healthmodel
		}

		// Handle PetTreatment
		if v.Treatments != nil {
			treatmentmodel := &model.PetTreatment{}
			if v.Treatments.RabiesVaccinationDate != nil {
				treatmentmodel.RabiesVaccinationDate = v.Treatments.RabiesVaccinationDate
			}
			if v.Treatments.InfectionVaccinationDate != nil {
				treatmentmodel.InfectionVaccinationDate = v.Treatments.InfectionVaccinationDate
			}
			if v.Treatments.EctoparasiteTreatmentDate != nil {
				treatmentmodel.EctoparasiteTreatmentDate = v.Treatments.EctoparasiteTreatmentDate
			}
			if v.Treatments.DewormingDate != nil {
				treatmentmodel.DewormingDate = v.Treatments.DewormingDate
			}
			petmodel.Treatments = treatmentmodel
		}

		// Handle PetAnalysis
		analysesDomain := []*model.PetAnalysis{}
		if v.Analyses != nil {
			processGroup := func(group []*dto.PetAnalysis, name string) {
				for _, a := range group {
					if a != nil && a.AnalysisDate != nil {
						analysis := &model.PetAnalysis{
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
			petmodel.Analyses = analysesDomain
		}
	case *dto.PetUpdate:
		if v.Name != nil {
			petmodel.Name = *v.Name
		}
		if v.Type != nil {
			petmodel.Type = model.PetType(*v.Type)
		}
		if v.WeightKg != nil {
			petmodel.WeightKg = *v.WeightKg
		}
		if v.Gender != nil {
			petmodel.Gender = model.Gender(*v.Gender)
		}
		if v.ChipNumber != nil {
			petmodel.ChipNumber = *v.ChipNumber
		}
		if v.LivingCondition != nil {
			petmodel.LivingCondition = model.LivingCondition(*v.LivingCondition)
		}
		if v.BreedID != nil {
			petmodel.BreedRefID = v.BreedID
		}
		if v.BloodGroup != nil {
			petmodel.BloodGroupName = v.BloodGroup
		}
		if v.BirthDate != nil {
			petmodel.BirthDate = v.BirthDate
		}
		if v.ReproductiveStatus != nil {
			petmodel.ReproductiveStatus = model.ReproductiveStatus(*v.ReproductiveStatus)
		}
		if v.Bonuses != nil {
			petmodel.Bonuses = *v.Bonuses
		}
		if v.AgeMonths != nil || v.AgeYears != nil {
			petmodel.BirthDate = calculateBirthDateFromAge(v.AgeYears, v.AgeMonths)
		}

		// Handle PetHealth
		if v.Health != nil {
			healthmodel := &model.PetHealth{}
			if v.Health.HealthStatus != nil {
				healthmodel.HealthStatus = model.HealthStatus(*v.Health.HealthStatus)
			}
			if v.Health.LastDonation != nil {
				healthmodel.LastDonation = v.Health.LastDonation
			}
			if v.Health.Transfused != nil {
				healthmodel.Transfused = v.Health.Transfused
			}
			if v.Health.Medications != nil {
				healthmodel.Medications = v.Health.Medications
			}
			if v.Health.SurgicalInterventions != nil {
				healthmodel.SurgicalInterventions = v.Health.SurgicalInterventions
			}
			petmodel.Health = healthmodel
		}

		// Handle PetTreatment
		if v.Treatments != nil {
			treatmentmodel := &model.PetTreatment{}
			if v.Treatments.RabiesVaccinationDate != nil {
				treatmentmodel.RabiesVaccinationDate = v.Treatments.RabiesVaccinationDate
			}
			if v.Treatments.InfectionVaccinationDate != nil {
				treatmentmodel.InfectionVaccinationDate = v.Treatments.InfectionVaccinationDate
			}
			if v.Treatments.EctoparasiteTreatmentDate != nil {
				treatmentmodel.EctoparasiteTreatmentDate = v.Treatments.EctoparasiteTreatmentDate
			}
			if v.Treatments.DewormingDate != nil {
				treatmentmodel.DewormingDate = v.Treatments.DewormingDate
			}
			petmodel.Treatments = treatmentmodel
		}

		// Handle PetAnalysis
		analysesDomain := []*model.PetAnalysis{}
		if v.Analyses != nil {
			processGroup := func(group []*dto.PetAnalysis, name string) {
				for _, a := range group {
					if a != nil && a.AnalysisDate != nil {
						analysis := &model.PetAnalysis{
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
			petmodel.Analyses = analysesDomain
		}
	}
}

// ToDTO maps model.Pet to dto.Pet
func ToDTO(petmodel model.Pet) dto.Pet {
	petDTO := dto.Pet{
		ID:                 petmodel.ID,
		Name:               petmodel.Name,
		ChipNumber:         petmodel.ChipNumber,
		PhotoURLs:          petmodel.PhotoURLs,
		WeightKg:           petmodel.WeightKg,
		BirthDate:          petmodel.BirthDate,
		PetStatus:          string(petmodel.PetStatus),
		LivingCondition:    string(petmodel.LivingCondition),
		Gender:             string(petmodel.Gender),
		Type:               string(petmodel.Type),
		ReproductiveStatus: string(petmodel.ReproductiveStatus),
		Bonuses:            petmodel.Bonuses,
		CreatedAt:          petmodel.CreatedAt,
		UpdatedAt:          petmodel.UpdatedAt,
		DeletedAt:          petmodel.DeletedAt,
	}

	// Map StopFactors and WarnFactors into DonorRestrictions
	if len(petmodel.StopFactors) > 0 || len(petmodel.WarnFactors) > 0 {
		var stopFactors []dto.RestrictionFactor
		var warnFactors []dto.RestrictionFactor
		for _, code := range petmodel.StopFactors {
			desc := model.GetFactorDescription(model.FactorCode(code))
			factor := dto.RestrictionFactor{
				Code:           code,
				Description:    desc.Description,
				SubDescription: desc.SubDescription,
			}
			stopFactors = append(stopFactors, factor)
		}
		for _, code := range petmodel.WarnFactors {
			desc := model.GetFactorDescription(model.FactorCode(code))
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

	if petmodel.BreedRefID != nil {
		petDTO.BreedID = *petmodel.BreedRefID
	}
	if petmodel.BloodGroupName != nil {
		petDTO.BloodGroup = *petmodel.BloodGroupName
	}

	if petmodel.Health != nil {
		healthStatus := string(petmodel.Health.HealthStatus)
		petDTO.Health = &dto.PetHealth{
			HealthStatus:          &healthStatus,
			LastDonation:          petmodel.Health.LastDonation,
			Transfused:            petmodel.Health.Transfused,
			Medications:           petmodel.Health.Medications,
			SurgicalInterventions: petmodel.Health.SurgicalInterventions,
		}
	}

	if petmodel.Treatments != nil {
		petDTO.Treatments = &dto.PetTreatment{
			RabiesVaccinationDate:     petmodel.Treatments.RabiesVaccinationDate,
			InfectionVaccinationDate:  petmodel.Treatments.InfectionVaccinationDate,
			EctoparasiteTreatmentDate: petmodel.Treatments.EctoparasiteTreatmentDate,
			DewormingDate:             petmodel.Treatments.DewormingDate,
		}
	}

	// Analyses mapping if needed, but commented in original ToDomain
	if petmodel.Analyses != nil {
		petDTO.Analyses = &dto.PetAnalysisGroup{}
		for _, a := range petmodel.Analyses {
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
