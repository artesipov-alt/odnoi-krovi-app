package handlers

import (
	"context"
	"net/http"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/dto"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/services"
	"github.com/danielgtaylor/huma/v2"
)

// FileHandler обрабатывает HTTP запросы для справочных данных
type FileHandler struct {
	fileService     services.FileService
	petService      services.PetService
	userService     services.UserService
	bloodReqService services.BloodSearchService
}

// NewFileHandler создает новый обработчик справочных данных
func NewFileHandler(fileService services.FileService, petService services.PetService, userService services.UserService, bloodReqService services.BloodSearchService) *FileHandler {
	return &FileHandler{
		fileService:     fileService,
		petService:      petService,
		userService:     userService,
		bloodReqService: bloodReqService,
	}
}

// Register регистрирует маршруты справочников в Huma API
func (h *FileHandler) Register(api huma.API) {
	// Получить ссылку для загрузки фотографии питомца
	huma.Register(api, huma.Operation{
		OperationID: "get-presigned-url",
		Method:      http.MethodPost,
		Path:        "/v1/uploads/presign/{id}",
		Summary:     "Получить ссылку для загрузки фотографии",
		Description: "Возвращает временную ссылку для загрузки фотографии по ID",
		Tags:        []string{"pets-v1", "users-v1"},
	}, h.GetPresignURL)

	// Подтверждение загрузки аватарки питомца
	// huma.Register(api, huma.Operation{
	// 	OperationID: "confirm-upload",
	// 	Method:      http.MethodPost,
	// 	Path:        "/v1/uploads/confirm/{path}",
	// 	Summary:     "Подтверждение загрузки фото",
	// 	Description: "Подтверждает загрузку фотографии, делает её публичной и возвращает публичную ссылку",
	// 	Tags:        []string{"pets-v1", "users-v1"},
	// }, h.ConfirmUpload)
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

// func (h *FileHandler) ConfirmUpload(ctx context.Context, input *dto.AvatarPathParam) (*dto.ConfirmUploadResponse, error) {
// 	publicURL, err := h.petService.UpdatePetAvatar(ctx, input.Path)
// 	if err != nil {
// 		// This error is likely a server-side issue if the path was valid but the update failed.
// 		slog.ErrorContext(ctx, "failed to confirm pet avatar upload", "path", input.Path, "error", err.Error())
// 		return nil, huma.Error500InternalServerError("Внутренняя ошибка сервера")
// 	}

// 	return &dto.ConfirmUploadResponse{Body: struct {
// 		PublicURL string `json:"publicUrl"`
// 	}{PublicURL: publicURL}}, nil
// }
