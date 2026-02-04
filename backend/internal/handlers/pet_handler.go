package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/ent"
	"github.com/artesipov-alt/odnoi-krovi-app/ent/pet"
	"github.com/artesipov-alt/odnoi-krovi-app/ent/petanalysis"
	"github.com/artesipov-alt/odnoi-krovi-app/ent/pethealth"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/dto"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/services"
	"github.com/danielgtaylor/huma/v2"
)

// PetHandler обрабатывает HTTP запросы для операций с питомцами
type PetHandler struct {
	petService services.PetService
}

// NewPetHandler создает новый обработчик питомцев
func NewPetHandler(petService services.PetService) *PetHandler {
	return &PetHandler{
		petService: petService,
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

}

//==========Handlers==============================

func (h *PetHandler) CreatePet(ctx context.Context, input *struct {
	dto.PetUserIDPath
	Body dto.PetCreate
}) (*dto.PetResponse, error) {
	petData := h.toCreateENT(input.Body)

	pet, err := h.petService.CreatePet(ctx, input.ID, petData)
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
	dto.IDPath
	dto.PetPreloadQuery
}) (*dto.PetResponse, error) {
	preloads := h.getPreloads(input.PetPreloadQuery)

	pet, err := h.petService.GetPetByID(ctx, input.ID, preloads...)
	if err != nil {
		if errors.Is(err, apperrors.ErrPetNotFound) {
			slog.DebugContext(ctx, "pet not found", "pet_id", input.ID, "error", err.Error())
			return nil, huma.Error404NotFound("Питомец не найден")
		}
		slog.ErrorContext(ctx, "failed to get pet by ID", "pet_id", input.ID, "error", err.Error())
		return nil, huma.Error500InternalServerError("Внутренняя ошибка сервера")
	}

	return &dto.PetResponse{Body: h.toDTO(pet)}, nil
}

func (h *PetHandler) GetUserPets(ctx context.Context, input *struct {
	dto.PetUserIDPath
	dto.PetPreloadQuery
}) (*dto.PetsResponse, error) {
	preloads := h.getPreloads(input.PetPreloadQuery)

	pets, err := h.petService.GetUserPets(ctx, input.ID, preloads...)
	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			slog.DebugContext(ctx, "user not found for getting pets", "user_id", input.ID, "error", err.Error())
			return nil, huma.Error404NotFound("Пользователь не найден")
		}
		slog.ErrorContext(ctx, "failed to get user pets", "user_id", input.ID, "error", err.Error())
		return nil, huma.Error500InternalServerError("Внутренняя ошибка сервера")
	}

	var petDTOs []dto.Pet
	for _, p := range pets {
		petDTOs = append(petDTOs, h.toDTO(p))
	}

	return &dto.PetsResponse{Body: petDTOs}, nil
}

func (h *PetHandler) UpdatePet(ctx context.Context, input *struct {
	dto.IDPath
	Body dto.PetUpdate
}) (*dto.MessageResponse, error) {
	updates, health, treatments, analyses, bonuses := h.toUpdateENT(input.Body)

	if err := h.petService.UpdatePet(ctx, input.ID, updates, health, treatments, analyses, bonuses); err != nil {
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

func (h *PetHandler) DeletePet(ctx context.Context, input *dto.IDPath) (*dto.MessageResponse, error) {
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

// getPreloads извлекает список связей для предзагрузки из query-параметров
func (h *PetHandler) getPreloads(pq dto.PetPreloadQuery) []string {
	var preloads []string
	if pq.WithHealth {
		preloads = append(preloads, "Health")
	}
	if pq.WithTreatments {
		preloads = append(preloads, "Treatments")
	}
	if pq.WithAnalysis {
		preloads = append(preloads, "Analyses") // Changed from "Analysis" to "Analyses" to match the DTO struct
	}
	if pq.WithBonuses {
		preloads = append(preloads, "Bonuses")
	}
	if pq.WithAll {
		return []string{"Health", "Treatments", "Analyses", "Bonuses"} // Changed from "Analysis" to "Analyses"
	}
	return preloads
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
	return petDTO
}

// mapDTOToPet преобразует DTO создания питомца в ENT модель
func (h *PetHandler) toCreateENT(d dto.PetCreate) *ent.Pet {
	if d.BirthDate == nil && (d.AgeYears > 0 || d.AgeMonths > 0) {
		birthDate := time.Now().AddDate(-d.AgeYears, -d.AgeMonths, 0)
		birthDate = time.Date(birthDate.Year(), birthDate.Month(), 1, 0, 0, 0, 0, time.UTC)
		d.BirthDate = &birthDate
	}
	p := &ent.Pet{
		Name:               d.Name,
		ChipNumber:         d.ChipNumber,
		BreedID:            d.BreedID,
		WeightKg:           d.WeightKg,
		BirthDate:          d.BirthDate,
		LivingCondition:    pet.LivingCondition(d.LivingCondition),
		Gender:             pet.Gender(d.Gender),
		Type:               pet.Type(d.Type),
		BloodGroup:         d.BloodGroup,
		ReproductiveStatus: pet.ReproductiveStatus(d.ReproductiveStatus),
		PetStatus:          pet.PetStatus(d.PetStatus),
	}

	if d.Health != nil {
		p.Edges.Health = &ent.PetHealth{
			LastDonation: d.Health.LastDonation,
		}
		if d.Health.HealthStatus != nil {
			p.Edges.Health.HealthStatus = pethealth.HealthStatus(*d.Health.HealthStatus)
		}
		if d.Health.Transfused != nil {
			p.Edges.Health.Transfused = *d.Health.Transfused
		}
		if d.Health.Medications != nil {
			p.Edges.Health.Medications = *d.Health.Medications
		}
		if d.Health.SurgicalInterventions != nil {
			p.Edges.Health.SurgicalInterventions = *d.Health.SurgicalInterventions
		}
	}

	if d.Treatments != nil {
		p.Edges.Treatments = &ent.PetTreatment{
			RabiesVaccinationDate:     d.Treatments.RabiesVaccinationDate,
			InfectionVaccinationDate:  d.Treatments.InfectionVaccinationDate,
			EctoparasiteTreatmentDate: d.Treatments.EctoparasiteTreatmentDate,
			DewormingDate:             d.Treatments.DewormingDate,
		}
	}

	if d.Analyses != nil {
		var analyses []*ent.PetAnalysis
		processGroup := func(group []*dto.PetAnalysis, name petanalysis.AnalysisName) {
			for _, a := range group {
				entA := &ent.PetAnalysis{
					AnalysisName: name,
					AnalysisDate: a.AnalysisDate,
				}
				if a.AnalysisType != nil {
					entA.AnalysisType = petanalysis.AnalysisType(*a.AnalysisType)
				}
				analyses = append(analyses, entA)
			}
		}

		processGroup(d.Analyses.Leukemia, petanalysis.AnalysisNameLeukemia)
		processGroup(d.Analyses.Immunodeficiency, petanalysis.AnalysisNameImmunodeficiency)
		processGroup(d.Analyses.Hemoplasmosis, petanalysis.AnalysisNameHemoplasmosis)
		processGroup(d.Analyses.Bartonellosis, petanalysis.AnalysisNameBartonellosis)
		processGroup(d.Analyses.Babesiosis, petanalysis.AnalysisNameBabesiosis)
		processGroup(d.Analyses.Dirofilaria, petanalysis.AnalysisNameDirofilaria)
		processGroup(d.Analyses.Ehrlichiosis, petanalysis.AnalysisNameEhrlichiosis)
		processGroup(d.Analyses.Anaplasmosis, petanalysis.AnalysisNameAnaplasmosis)

		p.Edges.Analyses = analyses
	}

	if d.Bonuses != nil {
		p.Edges.Bonuses = &ent.PetBonus{
			IsArtist:      d.Bonuses.IsArtist,
			IsTherapist:   d.Bonuses.IsTherapist,
			IsFormerDonor: d.Bonuses.IsFormerDonor,
			IsGuideDog:    d.Bonuses.IsGuideDog,
		}
	}

	return p
}

// mapDTOToPetUpdates преобразует DTO обновления питомца в аргументы для сервиса
func (h *PetHandler) toUpdateENT(d dto.PetUpdate) (map[string]any, *ent.PetHealth, *ent.PetTreatment, []*ent.PetAnalysis, *ent.PetBonus) {
	updates := make(map[string]any)
	if d.Name != nil {
		updates["Name"] = *d.Name
	}
	if d.ChipNumber != nil {
		updates["ChipNumber"] = *d.ChipNumber
	}
	if d.PhotoURLs != nil {
		updates["PhotoUrls"] = d.PhotoURLs
	}
	if d.BreedID != nil {
		updates["BreedID"] = *d.BreedID
	}
	if d.WeightKg != nil {
		updates["WeightKg"] = *d.WeightKg
	}
	if d.BirthDate != nil {
		updates["BirthDate"] = d.BirthDate
	} else if d.AgeYears != nil || d.AgeMonths != nil {
		ageYears := 0
		if d.AgeYears != nil {
			ageYears = *d.AgeYears
		}
		ageMonths := 0
		if d.AgeMonths != nil {
			ageMonths = *d.AgeMonths
		}
		if ageYears > 0 || ageMonths > 0 {
			birthDate := time.Now().AddDate(-ageYears, -ageMonths, 0)
			birthDate = time.Date(birthDate.Year(), birthDate.Month(), 1, 0, 0, 0, 0, time.UTC)
			updates["BirthDate"] = &birthDate
		}
	}
	if d.LivingCondition != nil {
		updates["LivingCondition"] = *d.LivingCondition
	}
	if d.Gender != nil {
		updates["Gender"] = *d.Gender
	}
	if d.Type != nil {
		updates["Type"] = *d.Type
	}
	if d.BloodGroup != nil {
		updates["BloodGroup"] = *d.BloodGroup
	}
	if d.ReproductiveStatus != nil {
		updates["ReproductiveStatus"] = *d.ReproductiveStatus
	}
	if d.PetStatus != nil {
		updates["PetStatus"] = *d.PetStatus
	}

	var health *ent.PetHealth
	if d.Health != nil {
		health = &ent.PetHealth{
			LastDonation: d.Health.LastDonation,
		}
		if d.Health.HealthStatus != nil {
			health.HealthStatus = pethealth.HealthStatus(*d.Health.HealthStatus)
		}
		if d.Health.Transfused != nil {
			health.Transfused = *d.Health.Transfused
		}
		if d.Health.Medications != nil {
			health.Medications = *d.Health.Medications
		}
		if d.Health.SurgicalInterventions != nil {
			health.SurgicalInterventions = *d.Health.SurgicalInterventions
		}
	}

	var treatments *ent.PetTreatment
	if d.Treatments != nil {
		treatments = &ent.PetTreatment{
			RabiesVaccinationDate:     d.Treatments.RabiesVaccinationDate,
			InfectionVaccinationDate:  d.Treatments.InfectionVaccinationDate,
			EctoparasiteTreatmentDate: d.Treatments.EctoparasiteTreatmentDate,
			DewormingDate:             d.Treatments.DewormingDate,
		}
	}

	var analyses []*ent.PetAnalysis
	if d.Analyses != nil {
		processGroup := func(group []*dto.PetAnalysis, name petanalysis.AnalysisName) {
			for _, a := range group {
				entA := &ent.PetAnalysis{
					AnalysisName: name,
					AnalysisDate: a.AnalysisDate,
				}
				if a.AnalysisType != nil {
					entA.AnalysisType = petanalysis.AnalysisType(*a.AnalysisType)
				}
				analyses = append(analyses, entA)
			}
		}

		processGroup(d.Analyses.Leukemia, petanalysis.AnalysisNameLeukemia)
		processGroup(d.Analyses.Immunodeficiency, petanalysis.AnalysisNameImmunodeficiency)
		processGroup(d.Analyses.Hemoplasmosis, petanalysis.AnalysisNameHemoplasmosis)
		processGroup(d.Analyses.Bartonellosis, petanalysis.AnalysisNameBartonellosis)
		processGroup(d.Analyses.Babesiosis, petanalysis.AnalysisNameBabesiosis)
		processGroup(d.Analyses.Dirofilaria, petanalysis.AnalysisNameDirofilaria)
		processGroup(d.Analyses.Ehrlichiosis, petanalysis.AnalysisNameEhrlichiosis)
		processGroup(d.Analyses.Anaplasmosis, petanalysis.AnalysisNameAnaplasmosis)
	}

	var bonuses *ent.PetBonus
	if d.Bonuses != nil {
		bonuses = &ent.PetBonus{
			IsArtist:      d.Bonuses.IsArtist,
			IsTherapist:   d.Bonuses.IsTherapist,
			IsFormerDonor: d.Bonuses.IsFormerDonor,
			IsGuideDog:    d.Bonuses.IsGuideDog,
		}
	}

	return updates, health, treatments, analyses, bonuses
}
