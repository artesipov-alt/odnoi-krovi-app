package handlers

import (
	"context"
	"log/slog"
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
	slog.InfoContext(ctx, "Сброс пользователя к заводским настройкам", "user_id", input.ID)

	if err := h.userRepo.ResetUser(ctx, input.ID); err != nil {
		slog.ErrorContext(ctx, "Ошибка при сбросе пользователя", "error", err, "user_id", input.ID)
		return nil, err
	}

	slog.InfoContext(ctx, "Пользователь успешно сброшен", "user_id", input.ID)
	return &DevResponseWrapper{Body: dto.DevResponse{
		Status:  true,
		Message: "Пользователь успешно сброшен к заводским настройкам",
	}}, nil
}

func (h *DevHandler) RestoreUser(ctx context.Context, input *UserIDPath) (*DevResponseWrapper, error) {
	slog.InfoContext(ctx, "Восстановление удаленного пользователя", "user_id", input.ID)

	if err := h.userRepo.RestoreUser(ctx, input.ID); err != nil {
		slog.ErrorContext(ctx, "Ошибка при восстановлении пользователя", "error", err, "user_id", input.ID)
		return nil, err
	}

	slog.InfoContext(ctx, "Пользователь успешно восстановлен", "user_id", input.ID)
	return &DevResponseWrapper{Body: dto.DevResponse{
		Status:  true,
		Message: "Пользователь успешно восстановлен",
	}}, nil
}

func (h *DevHandler) GetDeletedUsers(ctx context.Context, input *struct{}) (*GetDeletedUsersResponseWrapper, error) {
	slog.InfoContext(ctx, "Получение списка удаленных пользователей")

	users, err := h.userRepo.GetDeletedUsers(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "Ошибка при получении удаленных пользователей", "error", err)
		return nil, err
	}

	slog.InfoContext(ctx, "Удаленные пользователи успешно получены", "count", len(users))
	return &GetDeletedUsersResponseWrapper{Body: dto.GetDeletedUsersResponse{
		Status:  true,
		Message: "Удаленные пользователи успешно получены",
		Users:   users,
	}}, nil
}
