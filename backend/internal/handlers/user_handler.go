package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/ent"
	"github.com/artesipov-alt/odnoi-krovi-app/ent/user"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
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
	}, h.User)

	// Простая регистрация пользователя
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
	}, h.UserByTelegram)

	// Удаление пользователя по ID
	huma.Register(api, huma.Operation{
		OperationID: "delete-user",
		Method:      http.MethodDelete,
		Path:        "/v1/user/{id}",
		Summary:     "Удаление пользователя по ID",
		Description: "Удаляет пользователя из системы (soft delete)",
		Tags:        []string{"users-v1"},
	}, h.DeleteUser)

	// Сброс пользователя к начальным настройкам
	huma.Register(api, huma.Operation{
		OperationID: "reset-user",
		Method:      http.MethodPost,
		Path:        "/v1/user/reset-user/{id}",
		Summary:     "Сброс пользователя к начальным настройкам",
		Description: "Сбрасывает пользователя к заводским настройкам на этапе команды старт от бота",
		Tags:        []string{"users-v1", "dev"},
	}, h.ResetUser)

	// Восстановление удаленного пользователя
	huma.Register(api, huma.Operation{
		OperationID: "restore-user",
		Method:      http.MethodPost,
		Path:        "/v1/user/restore-user/{id}",
		Summary:     "Восстановление удаленного пользователя",
		Description: "Восстанавливает мягко удаленного пользователя, устанавливая deleted_at в NULL",
		Tags:        []string{"users-v1", "dev"},
	}, h.RestoreUser)

	// Получение всех удаленных пользователей
	huma.Register(api, huma.Operation{
		OperationID: "get-deleted-users",
		Method:      http.MethodGet,
		Path:        "/v1/user/deleted-users",
		Summary:     "Получение всех удаленных пользователей",
		Description: "Возвращает список всех мягко удаленных пользователей",
		Tags:        []string{"users-v1", "dev"},
	}, h.DeletedUsers)
}

// Handlers

func (h *UserHandler) User(ctx context.Context, input *dto.UserIDPath) (*dto.UserResponse, error) {
	u, err := h.userService.GetUserByID(ctx, input.ID)
	if err != nil {
		slog.ErrorContext(ctx, "failed to get user by ID", "user_id", input.ID, "error", err.Error())
		return nil, err
	}

	if u == nil {
		slog.DebugContext(ctx, "user not found", "user_id", input.ID)
		return nil, apperrors.ErrUserNotFound
	}

	return &dto.UserResponse{Body: h.toDTO(u)}, nil
}

func (h *UserHandler) RegisterUserSimple(ctx context.Context, input *struct {
	Body dto.UserRegistrationSimple
}) (*dto.UserResponse, error) {
	fullName := input.Body.FullName
	if fullName == "" {
		fullName = "Пользователь Telegram"
	}

	userData := &ent.User{
		TelegramID: input.Body.TelegramID,
		FullName:   fullName,
		Role:       "user",
	}

	u, err := h.userService.RegisterUserSimple(ctx, userData)
	if err != nil {
		slog.ErrorContext(ctx, "failed to register user simple", "telegram_id", input.Body.TelegramID, "error", err.Error())
		return nil, err
	}

	if u == nil {
		slog.ErrorContext(ctx, "registration returned nil user", "telegram_id", input.Body.TelegramID)
		return nil, apperrors.Internal(errors.New("registration returned nil user"), "ошибка при создании пользователя")
	}

	return &dto.UserResponse{Body: h.toDTO(u)}, nil
}

func (h *UserHandler) UpdateUser(ctx context.Context, input *struct {
	dto.UserIDPath
	Body dto.UserUpdate
}) (*dto.MessageResponse, error) {
	updates := h.toUpdatesMap(input.Body)

	if err := h.userService.UpdateUserProfile(ctx, input.ID, updates); err != nil {
		slog.ErrorContext(ctx, "failed to update user profile", "user_id", input.ID, "error", err.Error())
		return nil, err
	}

	return &dto.MessageResponse{
		Body: dto.MessageBody{
			Message: "Пользователь успешно обновлен",
		},
	}, nil
}

func (h *UserHandler) UserByTelegram(ctx context.Context, input *dto.TelegramIDQuery) (*dto.UserResponse, error) {
	u, err := h.userService.GetUserByTelegramID(ctx, input.TelegramID)
	if err != nil {
		slog.ErrorContext(ctx, "failed to get user by Telegram ID", "telegram_id", input.TelegramID, "error", err.Error())
		return nil, err
	}

	if u == nil {
		slog.DebugContext(ctx, "user not found", "telegram_id", input.TelegramID)
		return nil, apperrors.ErrUserNotFound
	}

	return &dto.UserResponse{Body: h.toDTO(u)}, nil
}

func (h *UserHandler) DeleteUser(ctx context.Context, input *dto.UserIDPath) (*dto.MessageResponse, error) {
	if err := h.userService.DeleteUser(ctx, input.ID); err != nil {
		slog.ErrorContext(ctx, "failed to delete user", "user_id", input.ID, "error", err.Error())
		return nil, err
	}

	return &dto.MessageResponse{
		Body: dto.MessageBody{
			Message: "Пользователь успешно удален",
		},
	}, nil
}

func (h *UserHandler) ResetUser(ctx context.Context, input *dto.UserIDPath) (*dto.MessageResponse, error) {
	if err := h.userService.ResetUser(ctx, input.ID); err != nil {
		slog.ErrorContext(ctx, "failed to reset user", "user_id", input.ID, "error", err.Error())
		return nil, err
	}

	return &dto.MessageResponse{
		Body: dto.MessageBody{
			Message: "Пользователь успешно сброшен к заводским настройкам",
		},
	}, nil
}

func (h *UserHandler) RestoreUser(ctx context.Context, input *dto.UserIDPath) (*dto.MessageResponse, error) {
	if err := h.userService.RestoreUser(ctx, input.ID); err != nil {
		slog.ErrorContext(ctx, "failed to restore user", "user_id", input.ID, "error", err.Error())
		return nil, err
	}

	return &dto.MessageResponse{
		Body: dto.MessageBody{
			Message: "Пользователь успешно восстановлен",
		},
	}, nil
}

func (h *UserHandler) DeletedUsers(ctx context.Context, input *struct{}) (*dto.UsersDeletedResponse, error) {
	users, err := h.userService.GetDeletedUsers(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "failed to get deleted users", "error", err.Error())
		return nil, err
	}

	userDTOs := make([]dto.User, len(users))
	for i, u := range users {
		userDTOs[i] = h.toDTO(u)
	}

	return &dto.UsersDeletedResponse{
		Body: dto.UsersDeletedBody{
			Message: "Удаленные пользователи успешно получены",
			Users:   userDTOs,
		},
	}, nil
}

// Helpers

// toDTO преобразует ENT модель пользователя в DTO для ответа
func (h *UserHandler) toDTO(u *ent.User) dto.User {
	if u == nil {
		return dto.User{}
	}

	formatDate := func(t any) string {
		if t == nil {
			return ""
		}
		switch v := t.(type) {
		case *time.Time:
			if v == nil || v.IsZero() || v.Unix() <= 0 {
				return ""
			}
			return v.Format(time.RFC3339)
		case time.Time:
			if v.IsZero() || v.Unix() <= 0 {
				return ""
			}
			return v.Format(time.RFC3339)
		default:
			return ""
		}
	}

	return dto.User{
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
		DeletedAt:        formatDate(u.DeletedAt),
	}
}

// toENT преобразует DTO пользователя в ENT модель
func (h *UserHandler) toENT(u *dto.User) ent.User {
	if u == nil {
		return ent.User{}
	}
	return ent.User{
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
		Role:             user.Role(u.Role),
	}
}

// toUpdatesMap преобразует DTO обновления в карту для сервиса
func (h *UserHandler) toUpdatesMap(d dto.UserUpdate) map[string]any {
	updates := make(map[string]any)
	if d.FullName != nil {
		updates["FullName"] = *d.FullName
	}
	if d.Phone != nil {
		updates["Phone"] = *d.Phone
	}
	if d.Email != nil {
		updates["Email"] = *d.Email
	}
	if d.AllowGeo != nil {
		updates["AllowGeo"] = *d.AllowGeo
	}
	if d.OnBoarding != nil {
		updates["OnBoarding"] = *d.OnBoarding
	}
	if d.LocationID != nil {
		updates["LocationID"] = *d.LocationID
	}
	return updates
}
