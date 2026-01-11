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

// UserHandler обрабатывает HTTP запросы для операций с пользователями
type UserHandler struct {
	userService services.UserService
}

// NewUserHandler создает новый обработчик пользователей
func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// Register регистрирует маршруты пользователя в Huma API
func (h *UserHandler) Register(api huma.API) {
	// Получение пользователя по ID
	huma.Register(api, huma.Operation{
		OperationID: "get-user-by-id",
		Method:      http.MethodGet,
		Path:        "/v1/user/{id}",
		Summary:     "Получение пользователя по ID",
		Description: "Возвращает информацию о пользователе по его идентификатору",
		Tags:        []string{"users-v1"},
	}, h.GetUser)

	// Простая регистрация пользователя
	huma.Register(api, huma.Operation{
		OperationID:   "register-user-simple",
		Method:        http.MethodPost,
		Path:          "/v1/user/register/simple",
		Summary:       "Простая регистрация пользователя",
		Description:   "Создает пользователя с Telegram ID и именем (для команды Start)",
		Tags:          []string{"users-v1"},
		DefaultStatus: http.StatusCreated,
	}, h.RegisterUserSimple)

	// Регистрация нового пользователя
	huma.Register(api, huma.Operation{
		OperationID:   "register-user",
		Method:        http.MethodPost,
		Path:          "/v1/user/register",
		Summary:       "Регистрация нового пользователя",
		Description:   "Регистрирует нового пользователя в системе",
		Tags:          []string{"users-v1"},
		DefaultStatus: http.StatusCreated,
		Deprecated:    true,
	}, h.RegisterUser)

	// Обновление данных пользователя
	huma.Register(api, huma.Operation{
		OperationID: "update-user",
		Method:      http.MethodPut,
		Path:        "/v1/user/{id}",
		Summary:     "Обновление данных пользователя",
		Description: "Обновляет информацию о пользователе",
		Tags:        []string{"users-v1"},
	}, h.UpdateUser)

	// Получение пользователя по Telegram ID
	huma.Register(api, huma.Operation{
		OperationID: "get-user-by-telegram",
		Method:      http.MethodGet,
		Path:        "/v1/user/telegram",
		Summary:     "Получение пользователя по Telegram ID",
		Description: "Возвращает информацию о пользователе по его Telegram ID",
		Tags:        []string{"users-v1"},
	}, h.GetUserByTelegram)

	// Удаление пользователя по ID
	huma.Register(api, huma.Operation{
		OperationID: "delete-user",
		Method:      http.MethodDelete,
		Path:        "/v1/user/{id}",
		Summary:     "Удаление пользователя по ID",
		Description: "Удаляет пользователя из системы (soft delete)",
		Tags:        []string{"users-v1"},
	}, h.DeleteUser)
}

// Вспомогательные структуры для Huma

type UserIDPath struct {
	ID string `path:"id" doc:"ID пользователя" minLength:"1" example:"1"`
}

type TelegramIDQuery struct {
	TelegramID int64 `query:"telegram_id" doc:"Telegram ID пользователя" minimum:"1" example:"123456789"`
}

type UserResponse struct {
	Body dto.UserResponseDTO
}

type MessageResponse struct {
	Body struct {
		Message string `json:"message" example:"Успешно"`
	}
}

// mapUserToDTO преобразует ENT модель пользователя в DTO для ответа
func mapUserToDTO(u *ent.User) dto.UserResponseDTO {
	return dto.UserResponseDTO{
		ID:               u.ID,
		TelegramID:       u.TelegramID,
		FullName:         u.FullName,
		Phone:            u.Phone,
		Email:            u.Email,
		OrganizationName: u.OrganizationName,
		ConsentPd:        u.ConsentPd,
		OnBoarding:       u.OnBoarding,
		AllowGeo:         u.AllowGeo,
		LocationID:       u.LocationID,
		Role:             string(u.Role),
		CreatedAt:        u.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:        u.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// Handlers

func (h *UserHandler) GetUser(ctx context.Context, input *UserIDPath) (*UserResponse, error) {
	slog.InfoContext(ctx, "Начало получения пользователя по ID", "user_id", input.ID)
	user, err := h.userService.GetUserByID(ctx, input.ID)
	if err != nil {
		slog.ErrorContext(ctx, "Ошибка получения пользователя по ID", "user_id", input.ID, "error", err)
		return nil, err
	}

	slog.InfoContext(ctx, "Пользователь успешно получен по ID", "user_id", input.ID)
	return &UserResponse{Body: mapUserToDTO(user)}, nil
}

func (h *UserHandler) RegisterUserSimple(ctx context.Context, input *struct {
	Body dto.SimpleRegistrationRequest
}) (*UserResponse, error) {
	slog.InfoContext(ctx, "Начало простой регистрации пользователя", "telegram_id", input.Body.TelegramID, "full_name", input.Body.FullName)

	fullName := input.Body.FullName
	if fullName == "" {
		fullName = "Пользователь Telegram"
	}

	user, err := h.userService.RegisterUserSimple(ctx, input.Body.TelegramID, fullName)
	if err != nil {
		slog.ErrorContext(ctx, "Ошибка простой регистрации пользователя", "telegram_id", input.Body.TelegramID, "error", err)
		return nil, err
	}

	slog.InfoContext(ctx, "Пользователь успешно зарегистрирован просто", "telegram_id", input.Body.TelegramID, "user_id", user.ID)
	return &UserResponse{Body: mapUserToDTO(user)}, nil
}

func (h *UserHandler) RegisterUser(ctx context.Context, input *struct {
	Body dto.UserRegistration
}) (*UserResponse, error) {
	// Извлекаем telegram_id из контекста (устанавливается middleware)
	telegramID, _ := ctx.Value("telegram_id").(int64)
	slog.InfoContext(ctx, "Начало регистрации пользователя", "telegram_id", telegramID)

	user, err := h.userService.RegisterUser(ctx, telegramID, input.Body)
	if err != nil {
		slog.ErrorContext(ctx, "Ошибка регистрации пользователя", "telegram_id", telegramID, "error", err)
		return nil, err
	}

	slog.InfoContext(ctx, "Пользователь успешно зарегистрирован", "telegram_id", telegramID, "user_id", user.ID)
	return &UserResponse{Body: mapUserToDTO(user)}, nil
}

func (h *UserHandler) UpdateUser(ctx context.Context, input *struct {
	UserIDPath
	Body dto.UserUpdate
}) (*MessageResponse, error) {
	slog.InfoContext(ctx, "Начало обновления данных пользователя", "user_id", input.ID)

	if err := h.userService.UpdateUserProfile(ctx, input.ID, input.Body); err != nil {
		slog.ErrorContext(ctx, "Ошибка обновления данных пользователя", "user_id", input.ID, "error", err)
		return nil, err
	}

	slog.InfoContext(ctx, "Данные пользователя успешно обновлены", "user_id", input.ID)
	resp := &MessageResponse{}
	resp.Body.Message = "Пользователь успешно обновлен"
	return resp, nil
}

func (h *UserHandler) GetUserByTelegram(ctx context.Context, input *TelegramIDQuery) (*UserResponse, error) {
	slog.InfoContext(ctx, "Начало получения пользователя по Telegram ID", "telegram_id", input.TelegramID)

	user, err := h.userService.GetUserByTelegramID(ctx, input.TelegramID)
	if err != nil {
		slog.ErrorContext(ctx, "Ошибка получения пользователя по Telegram ID", "telegram_id", input.TelegramID, "error", err)
		return nil, err
	}

	slog.InfoContext(ctx, "Пользователь успешно получен по Telegram ID", "telegram_id", input.TelegramID)
	return &UserResponse{Body: mapUserToDTO(user)}, nil
}

func (h *UserHandler) DeleteUser(ctx context.Context, input *UserIDPath) (*MessageResponse, error) {
	slog.InfoContext(ctx, "Начало удаления пользователя по ID", "user_id", input.ID)

	if err := h.userService.DeleteUser(ctx, input.ID); err != nil {
		slog.ErrorContext(ctx, "Ошибка удаления пользователя по ID", "user_id", input.ID, "error", err)
		return nil, err
	}

	slog.InfoContext(ctx, "Пользователь успешно удален по ID", "user_id", input.ID)
	resp := &MessageResponse{}
	resp.Body.Message = "Пользователь успешно удален"
	return resp, nil
}
