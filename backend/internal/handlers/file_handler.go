package handlers

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/dto"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/services"
	"github.com/danielgtaylor/huma/v2"
)

// FileHandler обрабатывает HTTP запросы для загрузки и подтверждения файлов
type FileHandler struct {
	fileService services.FileService
}

// NewFileHandler создает новый обработчик файлов
func NewFileHandler(fileService services.FileService) *FileHandler {
	return &FileHandler{
		fileService: fileService,
	}
}

// Register регистрирует маршруты для работы с файлами в Huma API
func (h *FileHandler) Register(api huma.API) {
	// Получить ссылку для загрузки фотографии
	huma.Register(api, huma.Operation{
		OperationID: "get-presigned-url",
		Method:      http.MethodPost,
		Path:        "/v1/uploads/presign/{id}",
		Summary:     "Получить ссылку для загрузки фотографии",
		Description: "Возвращает временную ссылку для загрузки фотографии по ID",
		Tags:        []string{"pets-v1", "users-v1", "blood-request-v1"},
	}, h.GetPresignURL)

	// Подтверждение загрузки фото
	huma.Register(api, huma.Operation{
		OperationID: "confirm-upload",
		Method:      http.MethodPost,
		Path:        "/v1/uploads/confirm",
		Summary:     "Подтверждение загрузки фото",
		Description: "Подтверждает загрузку массива фотографий, делает их публичными и обновляет сущность",
		Tags:        []string{"pets-v1", "users-v1", "blood-request-v1"},
	}, h.ConfirmUpload)
}

func (h *FileHandler) GetPresignURL(ctx context.Context, input *struct {
	dto.IDPath
	dto.PhotoPreloadQuery
}) (*dto.UploadURLResponse, error) {
	var preloads []string
	if input.ForPetAvatar {
		preloads = append(preloads, "pet_avatar")
	}
	if input.ForUserAvatar {
		preloads = append(preloads, "user_avatar")
	}
	if input.ForPetBlood {
		preloads = append(preloads, "blood_req")
	}

	if len(preloads) != 1 {
		return nil, huma.Error400BadRequest("Неверный запрос, Выбрать можно только один query запрос.")
	}
	// Если количество фотографий не передано, устанавливаем дефолтное значение 1
	if input.PhotosCount == 0 {
		input.PhotosCount = 1
	}
	if input.PhotosCount > 5 {
		return nil, huma.Error400BadRequest("Максимальное количество фотографий - 5")
	}

	uploadInfos, err := h.fileService.GetPresignURLs(ctx, input.ID, input.PhotosCount, preloads...)
	if err != nil {
		return nil, err
	}

	if len(uploadInfos) == 0 {
		return nil, huma.Error500InternalServerError("Не удалось получить ссылки для загрузки.")
	}

	items := make([]dto.UploadItem, len(uploadInfos))
	for i, info := range uploadInfos {
		items[i] = dto.UploadItem{
			URL:  info.UploadURL,
			Path: info.ObjectPath,
		}
	}

	return &dto.UploadURLResponse{
		Body: struct {
			Items []dto.UploadItem `json:"items" doc:"Список ссылок для загрузки"`
		}{
			Items: items,
		},
	}, nil
}

// getEntityType определяет тип сущности по префиксу ID
func getEntityType(id string) string {
	switch {
	case strings.HasPrefix(id, "USR"):
		return "user"
	case strings.HasPrefix(id, "PET"):
		return "pet"
	case strings.HasPrefix(id, "BLS"):
		return "blood_req"
	default:
		return ""
	}
}

func (h *FileHandler) ConfirmUpload(ctx context.Context, input *struct {
	Body dto.ConfirmUploadRequest
}) (*dto.MessageResponse, error) {
	entityType := getEntityType(input.Body.EntityID)
	if entityType == "" {
		return nil, huma.Error400BadRequest("Неверный ID сущности")
	}

	var preload string
	switch entityType {
	case "pet":
		preload = "pet_avatar"
	case "user":
		preload = "user_avatar"
	case "blood_req":
		preload = "blood_req"
	}

	err := h.fileService.ConfirmUploads(ctx, input.Body.EntityID, input.Body.Paths, preload)
	if err != nil {
		slog.ErrorContext(ctx, "failed to confirm uploads", "entity_id", input.Body.EntityID, "type", entityType, "error", err.Error())
		return nil, err
	}

	return &dto.MessageResponse{
		Body: dto.MessageBody{
			Message: "Фото успешно подтверждены и добавлены",
		},
	}, nil
}
