package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/artesipov-alt/odnoi-krovi-app/ent"
	"github.com/artesipov-alt/odnoi-krovi-app/ent/pet"
	"github.com/artesipov-alt/odnoi-krovi-app/ent/petanalysis"
	"github.com/artesipov-alt/odnoi-krovi-app/ent/pethealth"
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

	// Получить ссылку для загрузки фотографии питомца
	huma.Register(api, huma.Operation{
		OperationID: "get-avatar-upload-url",
		Method:      http.MethodGet,
		Path:        "/v1/pet/upload/avatar/{id}",
		Summary:     "Получить ссылку для загрузки фотографии питомца",
		Description: "Возвращает временную ссылку для загрузки фотографии питомца по ID",
		Tags:        []string{"pets-v1"},
	}, h.GetAvatarUploadURL)

	// Подтверждение загрузки аватарки питомца
	huma.Register(api, huma.Operation{
		OperationID: "confirm-pet-avatar-upload",
		Method:      http.MethodPost,
		Path:        "/v1/pet/upload/avatar/confirm/{path}",
		Summary:     "Подтверждение загрузки аватарки питомца",
		Description: "Подтверждает загрузку аватарки питомца, делает её публичной и возвращает публичную ссылку",
		Tags:        []string{"pets-v1"},
	}, h.ConfirmPetAvatarUpload)
}

// Вспомогательные структуры для Huma

type PetIDPath struct {
	ID string `path:"id" doc:"ID питомца" minLength:"1" example:"PET-25-000001"`
}

type PetUserIDPath struct {
	ID string `path:"user_id" doc:"ID пользователя" minLength:"1" example:"1"`
}

type PetPreloadQuery struct {
	WithHealth     bool `query:"with_health" doc:"Включить данные о здоровье"`
	WithTreatments bool `query:"with_treatments" doc:"Включить данные о ветеринарных обработках"`
	WithAnalysis   bool `query:"with_analysis" doc:"Включить данные об анализах"`
	WithBonuses    bool `query:"with_bonuses" doc:"Включить данные о бонусах"`
	WithAll        bool `query:"with_all" doc:"Включить все связанные данные"`
}

type AvatarPathParam struct {
	Path string `path:"path" doc:"Путь к аватарке питомца" example:"pets/PET-25-000001/avatar.jpg"`
}

type PetResponse struct {
	Body dto.PetResponseDTO
}

type PetsResponse struct {
	Body []dto.PetResponseDTO
}

type UploadURLResponse struct {
	Body struct {
		URL  string `json:"url"`
		Path string `json:"path"`
	}
}

type ConfirmUploadResponse struct {
	Body struct {
		PublicURL string `json:"publicUrl"`
	}
}

type MessageResponse struct {
	Body struct {
		Message string `json:"message" example:"Успешно"`
	}
}

// getPreloads извлекает список связей для предзагрузки из query-параметров
func (h *PetHandler) getPreloads(pq PetPreloadQuery) []string {
	var preloads []string
	if pq.WithHealth {
		preloads = append(preloads, "Health")
	}
	if pq.WithTreatments {
		preloads = append(preloads, "Treatments")
	}
	if pq.WithAnalysis {
		preloads = append(preloads, "Analysis")
	}
	if pq.WithBonuses {
		preloads = append(preloads, "Bonuses")
	}
	if pq.WithAll {
		return []string{"Health", "Treatments", "Analysis", "Bonuses"}
	}
	return preloads
}

// mapPetToDTO преобразует ENT модель питомца в DTO для ответа
func mapPetToDTO(p *ent.Pet) dto.PetResponseDTO {
	petDTO := dto.PetResponseDTO{
		ID:              p.ID,
		Name:            p.Name,
		ChipNumber:      p.ChipNumber,
		PhotoURL:        p.PhotoURL,
		BreedID:         p.BreedID,
		WeightKg:        p.WeightKg,
		AgeYears:        p.AgeYears,
		AgeMonths:       p.AgeMonths,
		BirthDate:       p.BirthDate,
		LivingCondition: string(p.LivingCondition),
		Gender:          string(p.Gender),
		Type:            string(p.Type),
		BloodGroup:      p.BloodGroup,
		PetStatus:       string(p.PetStatus),
		CreatedAt:       p.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:       p.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	if p.Edges.Health != nil {
		reproStatus := string(p.Edges.Health.ReproductiveStatus)
		healthStatus := string(p.Edges.Health.HealthStatus)
		medications := p.Edges.Health.Medications
		surgical := p.Edges.Health.SurgicalInterventions
		petDTO.Health = &dto.PetHealthDTO{
			ReproductiveStatus:    &reproStatus,
			HealthStatus:          &healthStatus,
			LastDonation:          p.Edges.Health.LastDonation,
			Transfused:            &p.Edges.Health.Transfused,
			Medications:           &medications,
			SurgicalInterventions: &surgical,
		}
	}
	if p.Edges.Treatments != nil {
		petDTO.Treatments = &dto.PetTreatmentDTO{
			RabiesVaccinationDate:     p.Edges.Treatments.RabiesVaccinationDate,
			InfectionVaccinationDate:  p.Edges.Treatments.InfectionVaccinationDate,
			EctoparasiteTreatmentDate: p.Edges.Treatments.EctoparasiteTreatmentDate,
			DewormingDate:             p.Edges.Treatments.DewormingDate,
		}
	}
	if p.Edges.Analyses != nil {
		petDTO.Analyses = make([]*dto.PetAnalysisDTO, len(p.Edges.Analyses))
		for i, a := range p.Edges.Analyses {
			leukemiaType := string(a.LeukemiaType)
			immunoType := string(a.ImmunodeficiencyType)
			hemoType := string(a.HemoplasmosisType)
			bartType := string(a.BartonellosisType)
			babeType := string(a.BabesiosisType)
			diroType := string(a.DirofilariaType)
			ehriType := string(a.EhrlichiosisType)
			anaType := string(a.AnaplasmosisType)
			petDTO.Analyses[i] = &dto.PetAnalysisDTO{
				LeukemiaDate:         a.LeukemiaDate,
				LeukemiaType:         &leukemiaType,
				ImmunodeficiencyDate: a.ImmunodeficiencyDate,
				ImmunodeficiencyType: &immunoType,
				HemoplasmosisDate:    a.HemoplasmosisDate,
				HemoplasmosisType:    &hemoType,
				BartonellosisDate:    a.BartonellosisDate,
				BartonellosisType:    &bartType,
				BabesiosisDate:       a.BabesiosisDate,
				BabesiosisType:       &babeType,
				DirofilariaDate:      a.DirofilariaDate,
				DirofilariaType:      &diroType,
				EhrlichiosisDate:     a.EhrlichiosisDate,
				EhrlichiosisType:     &ehriType,
				AnaplasmosisDate:     a.AnaplasmosisDate,
				AnaplasmosisType:     &anaType,
			}
		}
	}
	if p.Edges.Bonuses != nil {
		petDTO.Bonuses = &dto.PetBonusDTO{
			IsArtist:      p.Edges.Bonuses.IsArtist,
			IsTherapist:   p.Edges.Bonuses.IsTherapist,
			IsFormerDonor: p.Edges.Bonuses.IsFormerDonor,
			IsGuideDog:    p.Edges.Bonuses.IsGuideDog,
		}
	}
	return petDTO
}

// mapDTOToPet преобразует DTO создания питомца в ENT модель
func mapDTOToPet(d dto.PetCreate) *ent.Pet {
	p := &ent.Pet{
		Name:            d.Name,
		ChipNumber:      d.ChipNumber,
		PhotoURL:        d.PhotoURL,
		BreedID:         d.BreedID,
		WeightKg:        d.WeightKg,
		AgeYears:        d.AgeYears,
		AgeMonths:       d.AgeMonths,
		BirthDate:       d.BirthDate,
		LivingCondition: pet.LivingCondition(d.LivingCondition),
		Gender:          pet.Gender(d.Gender),
		Type:            pet.Type(d.Type),
		BloodGroup:      d.BloodGroup,
		PetStatus:       pet.PetStatus(d.PetStatus),
	}

	if d.Health != nil {
		p.Edges.Health = &ent.PetHealth{
			LastDonation: d.Health.LastDonation,
		}
		if d.Health.ReproductiveStatus != nil {
			p.Edges.Health.ReproductiveStatus = pethealth.ReproductiveStatus(*d.Health.ReproductiveStatus)
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
		p.Edges.Analyses = make([]*ent.PetAnalysis, len(d.Analyses))
		for i, a := range d.Analyses {
			p.Edges.Analyses[i] = &ent.PetAnalysis{
				LeukemiaDate:         a.LeukemiaDate,
				ImmunodeficiencyDate: a.ImmunodeficiencyDate,
				HemoplasmosisDate:    a.HemoplasmosisDate,
				BartonellosisDate:    a.BartonellosisDate,
				BabesiosisDate:       a.BabesiosisDate,
				DirofilariaDate:      a.DirofilariaDate,
				EhrlichiosisDate:     a.EhrlichiosisDate,
				AnaplasmosisDate:     a.AnaplasmosisDate,
			}
			if a.LeukemiaType != nil {
				p.Edges.Analyses[i].LeukemiaType = petanalysis.LeukemiaType(*a.LeukemiaType)
			}
			if a.ImmunodeficiencyType != nil {
				p.Edges.Analyses[i].ImmunodeficiencyType = petanalysis.ImmunodeficiencyType(*a.ImmunodeficiencyType)
			}
			if a.HemoplasmosisType != nil {
				p.Edges.Analyses[i].HemoplasmosisType = petanalysis.HemoplasmosisType(*a.HemoplasmosisType)
			}
			if a.BartonellosisType != nil {
				p.Edges.Analyses[i].BartonellosisType = petanalysis.BartonellosisType(*a.BartonellosisType)
			}
			if a.BabesiosisType != nil {
				p.Edges.Analyses[i].BabesiosisType = petanalysis.BabesiosisType(*a.BabesiosisType)
			}
			if a.DirofilariaType != nil {
				p.Edges.Analyses[i].DirofilariaType = petanalysis.DirofilariaType(*a.DirofilariaType)
			}
			if a.EhrlichiosisType != nil {
				p.Edges.Analyses[i].EhrlichiosisType = petanalysis.EhrlichiosisType(*a.EhrlichiosisType)
			}
			if a.AnaplasmosisType != nil {
				p.Edges.Analyses[i].AnaplasmosisType = petanalysis.AnaplasmosisType(*a.AnaplasmosisType)
			}
		}
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
func mapDTOToPetUpdates(d dto.PetUpdate) (map[string]any, *ent.PetHealth, *ent.PetTreatment, []*ent.PetAnalysis, *ent.PetBonus) {
	updates := make(map[string]any)
	if d.Name != nil {
		updates["Name"] = *d.Name
	}
	if d.ChipNumber != nil {
		updates["ChipNumber"] = *d.ChipNumber
	}
	if d.PhotoURL != nil {
		updates["PhotoURL"] = *d.PhotoURL
	}
	if d.BreedID != nil {
		updates["BreedID"] = *d.BreedID
	}
	if d.WeightKg != nil {
		updates["WeightKg"] = *d.WeightKg
	}
	if d.AgeYears != nil {
		updates["AgeYears"] = *d.AgeYears
	}
	if d.AgeMonths != nil {
		updates["AgeMonths"] = *d.AgeMonths
	}
	if d.BirthDate != nil {
		updates["BirthDate"] = d.BirthDate
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
	if d.PetStatus != nil {
		updates["PetStatus"] = *d.PetStatus
	}

	var health *ent.PetHealth
	if d.Health != nil {
		health = &ent.PetHealth{
			LastDonation: d.Health.LastDonation,
		}
		if d.Health.ReproductiveStatus != nil {
			health.ReproductiveStatus = pethealth.ReproductiveStatus(*d.Health.ReproductiveStatus)
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
		analyses = make([]*ent.PetAnalysis, len(d.Analyses))
		for i, a := range d.Analyses {
			analyses[i] = &ent.PetAnalysis{
				LeukemiaDate:         a.LeukemiaDate,
				ImmunodeficiencyDate: a.ImmunodeficiencyDate,
				HemoplasmosisDate:    a.HemoplasmosisDate,
				BartonellosisDate:    a.BartonellosisDate,
				BabesiosisDate:       a.BabesiosisDate,
				DirofilariaDate:      a.DirofilariaDate,
				EhrlichiosisDate:     a.EhrlichiosisDate,
				AnaplasmosisDate:     a.AnaplasmosisDate,
			}
			if a.LeukemiaType != nil {
				analyses[i].LeukemiaType = petanalysis.LeukemiaType(*a.LeukemiaType)
			}
			if a.ImmunodeficiencyType != nil {
				analyses[i].ImmunodeficiencyType = petanalysis.ImmunodeficiencyType(*a.ImmunodeficiencyType)
			}
			if a.HemoplasmosisType != nil {
				analyses[i].HemoplasmosisType = petanalysis.HemoplasmosisType(*a.HemoplasmosisType)
			}
			if a.BartonellosisType != nil {
				analyses[i].BartonellosisType = petanalysis.BartonellosisType(*a.BartonellosisType)
			}
			if a.BabesiosisType != nil {
				analyses[i].BabesiosisType = petanalysis.BabesiosisType(*a.BabesiosisType)
			}
			if a.DirofilariaType != nil {
				analyses[i].DirofilariaType = petanalysis.DirofilariaType(*a.DirofilariaType)
			}
			if a.EhrlichiosisType != nil {
				analyses[i].EhrlichiosisType = petanalysis.EhrlichiosisType(*a.EhrlichiosisType)
			}
			if a.AnaplasmosisType != nil {
				analyses[i].AnaplasmosisType = petanalysis.AnaplasmosisType(*a.AnaplasmosisType)
			}
		}
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

// Handlers

func (h *PetHandler) CreatePet(ctx context.Context, input *struct {
	PetUserIDPath
	Body dto.PetCreate
}) (*PetResponse, error) {
	petData := mapDTOToPet(input.Body)

	pet, err := h.petService.CreatePet(ctx, input.ID, petData)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			return nil, huma.Error404NotFound("Пользователь не найден")
		}
		return nil, huma.Error500InternalServerError("Внутренняя ошибка сервера")
	}

	return &PetResponse{Body: mapPetToDTO(pet)}, nil
}

func (h *PetHandler) GetPet(ctx context.Context, input *struct {
	PetIDPath
	PetPreloadQuery
}) (*PetResponse, error) {
	preloads := h.getPreloads(input.PetPreloadQuery)

	pet, err := h.petService.GetPetByID(ctx, input.ID, preloads...)
	if err != nil {
		if errors.Is(err, services.ErrPetNotFound) {
			return nil, huma.Error404NotFound("Питомец не найден")
		}
		return nil, huma.Error500InternalServerError("Внутренняя ошибка сервера")
	}

	return &PetResponse{Body: mapPetToDTO(pet)}, nil
}

func (h *PetHandler) GetUserPets(ctx context.Context, input *struct {
	PetUserIDPath
	PetPreloadQuery
}) (*PetsResponse, error) {
	preloads := h.getPreloads(input.PetPreloadQuery)

	pets, err := h.petService.GetUserPets(ctx, input.ID, preloads...)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			return nil, huma.Error404NotFound("Пользователь не найден")
		}
		return nil, huma.Error500InternalServerError("Внутренняя ошибка сервера")
	}

	var petDTOs []dto.PetResponseDTO
	for _, p := range pets {
		petDTOs = append(petDTOs, mapPetToDTO(p))
	}

	return &PetsResponse{Body: petDTOs}, nil
}

func (h *PetHandler) UpdatePet(ctx context.Context, input *struct {
	PetIDPath
	Body dto.PetUpdate
}) (*MessageResponse, error) {
	updates, health, treatments, analyses, bonuses := mapDTOToPetUpdates(input.Body)

	if err := h.petService.UpdatePet(ctx, input.ID, updates, health, treatments, analyses, bonuses); err != nil {
		if errors.Is(err, services.ErrPetNotFound) {
			return nil, huma.Error404NotFound("Питомец не найден")
		}
		return nil, huma.Error500InternalServerError("Внутренняя ошибка сервера")
	}

	resp := &MessageResponse{}
	resp.Body.Message = "Питомец успешно обновлен"
	return resp, nil
}

func (h *PetHandler) DeletePet(ctx context.Context, input *PetIDPath) (*MessageResponse, error) {
	if err := h.petService.DeletePet(ctx, input.ID); err != nil {
		if errors.Is(err, services.ErrPetNotFound) {
			return nil, huma.Error404NotFound("Питомец не найден")
		}
		return nil, huma.Error500InternalServerError("Внутренняя ошибка сервера")
	}

	resp := &MessageResponse{}
	resp.Body.Message = "Питомец успешно удален"
	return resp, nil
}

func (h *PetHandler) GetAvatarUploadURL(ctx context.Context, input *PetIDPath) (*UploadURLResponse, error) {
	url, path, err := h.petService.GetAvatarUploadURL(ctx, input.ID)
	if err != nil {
		if errors.Is(err, services.ErrPetNotFound) {
			return nil, huma.Error404NotFound("Питомец не найден")
		}
		return nil, huma.Error500InternalServerError("Внутренняя ошибка сервера")
	}

	return &UploadURLResponse{Body: struct {
		URL  string `json:"url"`
		Path string `json:"path"`
	}{URL: url, Path: path}}, nil
}

func (h *PetHandler) ConfirmPetAvatarUpload(ctx context.Context, input *AvatarPathParam) (*ConfirmUploadResponse, error) {
	publicURL, err := h.petService.UpdatePetAvatar(ctx, input.Path)
	if err != nil {
		return nil, huma.Error500InternalServerError("Внутренняя ошибка сервера")
	}

	return &ConfirmUploadResponse{Body: struct {
		PublicURL string `json:"publicUrl"`
	}{PublicURL: publicURL}}, nil
}
