package handlers

import (
	"context"
	"net/http"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/dto"
	repositories "github.com/artesipov-alt/odnoi-krovi-app/internal/repositories"
	"github.com/danielgtaylor/huma/v2"
)

// DevHandler обрабатывает HTTP запросы для инструментов разработки и отладки
type DevHandler struct {
	userRepo repositories.UserRepository
}

// NewDevHandler создает новый обработчик инструментов разработки
func NewDevHandler(userRepo repositories.UserRepository) *DevHandler {
	return &DevHandler{
		userRepo: userRepo,
	}
}

// Register регистрирует маршруты разработки в Huma API
func (h *DevHandler) Register(api huma.API) {
	// Сброс пользователя к начальным настройкам
	huma.Register(api, huma.Operation{
		OperationID: "reset-user",
		Method:      http.MethodPost,
		Path:        "/v1/dev/reset-user/{id}",
		Summary:     "Сброс пользователя к начальным настройкам",
		Description: "Сбрасывает пользователя к заводским настройкам на этапе команды старт от бота",
		Tags:        []string{"dev"},
	}, h.ResetUser)

	// Восстановление удаленного пользователя
	huma.Register(api, huma.Operation{
		OperationID: "restore-user",
		Method:      http.MethodPost,
		Path:        "/v1/dev/restore-user/{id}",
		Summary:     "Восстановление удаленного пользователя",
		Description: "Восстанавливает мягко удаленного пользователя, устанавливая deleted_at в NULL",
		Tags:        []string{"dev"},
	}, h.RestoreUser)

	// Получение всех удаленных пользователей
	huma.Register(api, huma.Operation{
		OperationID: "get-deleted-users",
		Method:      http.MethodGet,
		Path:        "/v1/dev/deleted-users",
		Summary:     "Получение всех удаленных пользователей",
		Description: "Возвращает список всех мягко удаленных пользователей",
		Tags:        []string{"dev"},
	}, h.GetDeletedUsers)
}

// Вспомогательные структуры для Huma

type DevResponseWrapper struct {
	Body dto.DevResponse
}

type GetDeletedUsersResponseWrapper struct {
	Body dto.GetDeletedUsersResponse
}

// Handlers

func (h *DevHandler) ResetUser(ctx context.Context, input *UserIDPath) (*DevResponseWrapper, error) {
	if err := h.userRepo.ResetUser(ctx, input.ID); err != nil {
		return nil, huma.Error500InternalServerError("Ошибка сервера")
	}

	return &DevResponseWrapper{Body: dto.DevResponse{
		Status:  true,
		Message: "Пользователь успешно сброшен к заводским настройкам",
	}}, nil
}

func (h *DevHandler) RestoreUser(ctx context.Context, input *UserIDPath) (*DevResponseWrapper, error) {
	if err := h.userRepo.RestoreUser(ctx, input.ID); err != nil {
		return nil, huma.Error500InternalServerError("Ошибка сервера")
	}

	return &DevResponseWrapper{Body: dto.DevResponse{
		Status:  true,
		Message: "Пользователь успешно восстановлен",
	}}, nil
}

func (h *DevHandler) GetDeletedUsers(ctx context.Context, input *struct{}) (*GetDeletedUsersResponseWrapper, error) {
	users, err := h.userRepo.GetDeletedUsers(ctx)
	if err != nil {
		return nil, huma.Error500InternalServerError("Ошибка сервера")
	}

	return &GetDeletedUsersResponseWrapper{Body: dto.GetDeletedUsersResponse{
		Status:  true,
		Message: "Удаленные пользователи успешно получены",
		Users:   users,
	}}, nil
}
