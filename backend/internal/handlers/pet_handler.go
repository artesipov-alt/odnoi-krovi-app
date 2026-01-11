package handlers

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/artesipov-alt/odnoi-krovi-app/ent"
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

// Handlers

func (h *PetHandler) CreatePet(ctx context.Context, input *struct {
	PetUserIDPath
	Body dto.PetCreate
}) (*PetResponse, error) {
	slog.InfoContext(ctx, "Начало создания питомца", "user_id", input.ID, "pet_name", input.Body.Name)

	pet, err := h.petService.CreatePet(ctx, input.ID, input.Body)
	if err != nil {
		slog.ErrorContext(ctx, "Ошибка создания питомца", "user_id", input.ID, "error", err)
		return nil, err
	}

	slog.InfoContext(ctx, "Питомец успешно создан", "user_id", input.ID, "pet_id", pet.ID)
	return &PetResponse{Body: mapPetToDTO(pet)}, nil
}

func (h *PetHandler) GetPet(ctx context.Context, input *struct {
	PetIDPath
	PetPreloadQuery
}) (*PetResponse, error) {
	preloads := h.getPreloads(input.PetPreloadQuery)
	slog.InfoContext(ctx, "Начало получения питомца по ID", "pet_id", input.ID, "preloads", preloads)

	pet, err := h.petService.GetPetByID(ctx, input.ID, preloads...)
	if err != nil {
		slog.ErrorContext(ctx, "Ошибка получения питомца по ID", "pet_id", input.ID, "error", err)
		return nil, err
	}

	slog.InfoContext(ctx, "Питомец успешно получен по ID", "pet_id", input.ID)
	return &PetResponse{Body: mapPetToDTO(pet)}, nil
}

func (h *PetHandler) GetUserPets(ctx context.Context, input *struct {
	PetUserIDPath
	PetPreloadQuery
}) (*PetsResponse, error) {
	preloads := h.getPreloads(input.PetPreloadQuery)
	slog.InfoContext(ctx, "Начало получения питомцев пользователя", "user_id", input.ID, "preloads", preloads)

	pets, err := h.petService.GetUserPets(ctx, input.ID, preloads...)
	if err != nil {
		slog.ErrorContext(ctx, "Ошибка получения питомцев пользователя", "user_id", input.ID, "error", err)
		return nil, err
	}

	var petDTOs []dto.PetResponseDTO
	for _, p := range pets {
		petDTOs = append(petDTOs, mapPetToDTO(p))
	}

	slog.InfoContext(ctx, "Питомцы пользователя успешно получены", "user_id", input.ID, "count", len(petDTOs))
	return &PetsResponse{Body: petDTOs}, nil
}

func (h *PetHandler) UpdatePet(ctx context.Context, input *struct {
	PetIDPath
	Body dto.PetUpdate
}) (*MessageResponse, error) {
	slog.InfoContext(ctx, "Начало обновления данных питомца", "pet_id", input.ID)

	if err := h.petService.UpdatePet(ctx, input.ID, input.Body); err != nil {
		slog.ErrorContext(ctx, "Ошибка обновления данных питомца", "pet_id", input.ID, "error", err)
		return nil, err
	}

	slog.InfoContext(ctx, "Данные питомца успешно обновлены", "pet_id", input.ID)
	resp := &MessageResponse{}
	resp.Body.Message = "Питомец успешно обновлен"
	return resp, nil
}

func (h *PetHandler) DeletePet(ctx context.Context, input *PetIDPath) (*MessageResponse, error) {
	slog.InfoContext(ctx, "Начало удаления питомца по ID", "pet_id", input.ID)

	if err := h.petService.DeletePet(ctx, input.ID); err != nil {
		slog.ErrorContext(ctx, "Ошибка удаления питомца по ID", "pet_id", input.ID, "error", err)
		return nil, err
	}

	slog.InfoContext(ctx, "Питомец успешно удален по ID", "pet_id", input.ID)
	resp := &MessageResponse{}
	resp.Body.Message = "Питомец успешно удален"
	return resp, nil
}

func (h *PetHandler) GetAvatarUploadURL(ctx context.Context, input *PetIDPath) (*UploadURLResponse, error) {
	slog.InfoContext(ctx, "Начало получения ссылки для загрузки фотографии питомца", "pet_id", input.ID)

	url, path, err := h.petService.GetAvatarUploadURL(ctx, input.ID)
	if err != nil {
		slog.ErrorContext(ctx, "Ошибка получения ссылки для загрузки фотографии питомца", "pet_id", input.ID, "error", err)
		return nil, err
	}

	slog.InfoContext(ctx, "Ссылка для загрузки фотографии питомца успешно получена", "pet_id", input.ID)
	return &UploadURLResponse{Body: struct {
		URL  string `json:"url"`
		Path string `json:"path"`
	}{URL: url, Path: path}}, nil
}

func (h *PetHandler) ConfirmPetAvatarUpload(ctx context.Context, input *AvatarPathParam) (*ConfirmUploadResponse, error) {
	slog.InfoContext(ctx, "Начало подтверждения загрузки аватарки питомца", "avatar_path", input.Path)

	publicURL, err := h.petService.UpdatePetAvatar(ctx, input.Path)
	if err != nil {
		slog.ErrorContext(ctx, "Ошибка подтверждения загрузки аватарки питомца", "avatar_path", input.Path, "error", err)
		return nil, err
	}

	slog.InfoContext(ctx, "Аватарка питомца успешно подтверждена", "avatar_path", input.Path)
	return &ConfirmUploadResponse{Body: struct {
		PublicURL string `json:"publicUrl"`
	}{PublicURL: publicURL}}, nil
}
