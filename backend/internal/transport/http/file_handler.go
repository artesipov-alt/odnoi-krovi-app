package http

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	bonuscmd "github.com/artesipov-alt/odnoi-krovi-app/internal/application/bonus/cmd"
	filecmd "github.com/artesipov-alt/odnoi-krovi-app/internal/application/file/cmd"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/transport/http/dto"
	"github.com/danielgtaylor/huma/v2"
)

// FileHandler обрабатывает HTTP запросы для загрузки и подтверждения файлов
type FileHandler struct {
	getPresignedHandler  *filecmd.GetPresignedURLsHandler
	confirmUploadHandler *filecmd.ConfirmUploadHandler
	importHandler        *bonuscmd.ImportBonusesHandler
}

// NewFileHandler создает новый обработчик файлов
func NewFileHandler(
	getPresignedHandler *filecmd.GetPresignedURLsHandler,
	confirmUploadHandler *filecmd.ConfirmUploadHandler,
	importHandler *bonuscmd.ImportBonusesHandler,
) *FileHandler {
	return &FileHandler{
		getPresignedHandler:  getPresignedHandler,
		confirmUploadHandler: confirmUploadHandler,
		importHandler:        importHandler,
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

	huma.Register(api, huma.Operation{
		OperationID: "import-bonuses",
		Method:      http.MethodPost,
		Path:        "/v1/admin/bonuses/import",
		Summary:     "Импорт бонусов из Excel",
		Description: "Загружает бонусы из Excel файла. Требуются права администратора.",
		Tags:        []string{"admin-v1"},
	}, h.ImportBonuses)
}

func (h *FileHandler) GetPresignURL(ctx context.Context, input *dto.GetUploadURLsInput) (*dto.GetUploadURLsOutput, error) {
	var preloads []string
	if input.ForPetAvatar {
		preloads = append(preloads, "pet_avatar")
	}
	if input.ForUserAvatar {
		preloads = append(preloads, "user_avatar")
	}
	if input.ForBloodReq {
		preloads = append(preloads, "blood_req")
	}

	if len(preloads) != 1 {
		return nil, apperrors.BadRequest("Неверный запрос, Выбрать можно только один query запрос.")
	}
	// Если количество фотографий не передано, устанавливаем дефолтное значение 1
	if input.PhotosCount == 0 {
		input.PhotosCount = 1
	}
	if input.PhotosCount > 5 {
		return nil, apperrors.BadRequest("Максимальное количество фотографий - 5")
	}

	uploadInfos, err := h.getPresignedHandler.Handle(ctx, input.ID, input.PhotosCount, preloads[0])
	if err != nil {
		return nil, err
	}

	if len(uploadInfos) == 0 {
		return nil, apperrors.Internal(nil, "Не удалось получить ссылки для загрузки.")
	}

	items := make([]dto.UploadItem, len(uploadInfos))
	for i, info := range uploadInfos {
		items[i] = dto.UploadItem{
			URL:  info.UploadURL,
			Path: info.ObjectPath,
		}
	}

	return &dto.GetUploadURLsOutput{
			Body: dto.UploadURLsResult{
				Items: items,
			},
		},
		nil
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

func (h *FileHandler) ConfirmUpload(ctx context.Context, input *dto.ConfirmUploadInput) (*dto.ConfirmUploadOutput, error) {
	entityType := getEntityType(input.Body.EntityID)
	if entityType == "" {
		return nil, apperrors.BadRequest("Неверный ID сущности")
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

	err := h.confirmUploadHandler.Handle(ctx, input.Body.EntityID, input.Body.Paths, preload)
	if err != nil {
		return nil, err
	}

	return &dto.ConfirmUploadOutput{
			Body: dto.ConfirmUploadResult{
				Message: "Фото подтверждены и добавлены",
			},
		},
		nil
}

func (h *FileHandler) ImportBonuses(ctx context.Context, input *dto.ImportBonusesInput) (*dto.ImportBonusesOutput, error) {
	formData := input.RawBody.Data()

	fileContent, err := io.ReadAll(formData.File)
	if err != nil {
		return nil, huma.Error400BadRequest("Не удалось прочитать файл")
	}

	result, err := h.importHandler.Handle(ctx, bytes.NewReader(fileContent))
	if err != nil {
		return nil, err
	}

	return &dto.ImportBonusesOutput{
		Body: dto.ImportBonusesResult{
			TotalRows: result.TotalRows,
			Imported:  result.Imported,
			Skipped:   result.Skipped,
			Errors:    result.Errors,
		},
	}, nil
}
