package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/artesipov-alt/odnoi-krovi-app/ent"
	"github.com/artesipov-alt/odnoi-krovi-app/ent/user"
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

	// Сброс пользователя к начальным настройкам
	huma.Register(api, huma.Operation{
		OperationID:   "reset-user",
		Method:        http.MethodPost,
		Path:          "/v1/user/reset-user/{id}",
		Summary:       "Сброс пользователя к начальным настройкам",
		Description:   "Сбрасывает пользователя к заводским настройкам на этапе команды старт от бота",
		Tags:          []string{"users-v1", "dev"},
		DefaultStatus: http.StatusOK,
	}, h.ResetUser)

	// Восстановление удаленного пользователя
	huma.Register(api, huma.Operation{
		OperationID:   "restore-user",
		Method:        http.MethodPost,
		Path:          "/v1/user/restore-user/{id}",
		Summary:       "Восстановление удаленного пользователя",
		Description:   "Восстанавливает мягко удаленного пользователя, устанавливая deleted_at в NULL",
		Tags:          []string{"users-v1", "dev"},
		DefaultStatus: http.StatusOK,
	}, h.RestoreUser)

	// Получение всех удаленных пользователей
	huma.Register(api, huma.Operation{
		OperationID:   "get-deleted-users",
		Method:        http.MethodGet,
		Path:          "/v1/user/deleted-users",
		Summary:       "Получение всех удаленных пользователей",
		Description:   "Возвращает список всех мягко удаленных пользователей",
		Tags:          []string{"users-v1", "dev"},
		DefaultStatus: http.StatusOK,
	}, h.GetDeletedUsers)
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

type GetDeletedUsersResponse struct {
	Message string                `json:"message"`
	Users   []dto.UserResponseDTO `json:"users"`
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

// mapDTOToUser преобразует DTO регистрации в ENT пользователя
func mapDTOToUser(dto dto.UserRegistrationFull, telegramID int64) *ent.User {
	return &ent.User{
		TelegramID: telegramID,
		FullName:   dto.FullName,
		Phone:      dto.Phone,
		Email:      dto.Email,
		ConsentPd:  dto.ConsentPD,
		LocationID: dto.LocationID,
		Role:       user.Role(dto.Role),
	}
}

// mapDTOToUserSimple преобразует DTO простой регистрации в ENT пользователя
func mapDTOToUserSimple(dto dto.UserRegistrationSimple) *ent.User {
	return &ent.User{
		TelegramID: dto.TelegramID,
		FullName:   dto.FullName,
		Phone:      "",
		Email:      "",
		ConsentPd:  true,
		OnBoarding: false,
		AllowGeo:   false,
		Role:       user.RoleUser,
	}
}

// mapDTOToUpdates преобразует DTO обновления в map для сервиса
func mapDTOToUpdates(dto dto.UserUpdate) map[string]any {
	updates := make(map[string]any)
	if dto.FullName != nil {
		updates["FullName"] = *dto.FullName
	}
	if dto.Phone != nil {
		updates["Phone"] = *dto.Phone
	}
	if dto.Email != nil {
		updates["Email"] = *dto.Email
	}
	if dto.AllowGeo != nil {
		updates["AllowGeo"] = *dto.AllowGeo
	}
	if dto.OnBoarding != nil {
		updates["OnBoarding"] = *dto.OnBoarding
	}
	if dto.LocationID != nil {
		updates["LocationID"] = *dto.LocationID
	}
	return updates
}

// Handlers

func (h *UserHandler) GetUser(ctx context.Context, input *UserIDPath) (*UserResponse, error) {
	user, err := h.userService.GetUserByID(ctx, input.ID)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			return nil, huma.Error404NotFound("Пользователь не найден")
		}
		return nil, huma.Error500InternalServerError("Ошибка сервера")
	}
	return &UserResponse{Body: mapUserToDTO(user)}, nil
}

func (h *UserHandler) RegisterUserSimple(ctx context.Context, input *struct {
	Body dto.UserRegistrationSimple
}) (*UserResponse, error) {
	fullName := input.Body.FullName
	if fullName == "" {
		fullName = "Пользователь Telegram"
	}

	userData := dto.UserRegistrationSimple{
		TelegramID: input.Body.TelegramID,
		FullName:   fullName,
	}

	user := mapDTOToUserSimple(userData)

	user, err := h.userService.RegisterUserSimple(ctx, user)
	if err != nil {
		if errors.Is(err, services.ErrUserAlreadyExists) {
			return nil, huma.Error409Conflict("Пользователь уже существует")
		}
		return nil, huma.Error500InternalServerError("Ошибка сервера")
	}

	return &UserResponse{Body: mapUserToDTO(user)}, nil
}

func (h *UserHandler) RegisterUser(ctx context.Context, input *struct {
	Body dto.UserRegistrationFull
}) (*UserResponse, error) {
	// Извлекаем telegram_id из контекста (устанавливается middleware)
	telegramID, _ := ctx.Value("telegram_id").(int64)

	user := mapDTOToUser(input.Body, telegramID)

	user, err := h.userService.RegisterUser(ctx, user)
	if err != nil {
		if errors.Is(err, services.ErrUserAlreadyExists) {
			return nil, huma.Error409Conflict("Пользователь уже существует")
		}
		if errors.Is(err, services.ErrInvalidRole) {
			return nil, huma.Error400BadRequest("Неверная роль")
		}
		if errors.Is(err, services.ErrLocationNotFound) {
			return nil, huma.Error400BadRequest("Неверная локация")
		}
		return nil, huma.Error500InternalServerError("Ошибка сервера")
	}

	return &UserResponse{Body: mapUserToDTO(user)}, nil
}

func (h *UserHandler) UpdateUser(ctx context.Context, input *struct {
	UserIDPath
	Body dto.UserUpdate
}) (*MessageResponse, error) {
	updates := mapDTOToUpdates(input.Body)

	if err := h.userService.UpdateUserProfile(ctx, input.ID, updates); err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			return nil, huma.Error404NotFound("Пользователь не найден")
		}
		if errors.Is(err, services.ErrLocationNotFound) {
			return nil, huma.Error400BadRequest("Неверная локация")
		}
		return nil, huma.Error500InternalServerError("Ошибка сервера")
	}

	resp := &MessageResponse{}
	resp.Body.Message = "Пользователь успешно обновлен"
	return resp, nil
}

func (h *UserHandler) GetUserByTelegram(ctx context.Context, input *TelegramIDQuery) (*UserResponse, error) {
	user, err := h.userService.GetUserByTelegramID(ctx, input.TelegramID)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			return nil, huma.Error404NotFound("Пользователь не найден")
		}
		return nil, huma.Error500InternalServerError("Ошибка сервера")
	}

	return &UserResponse{Body: mapUserToDTO(user)}, nil
}

func (h *UserHandler) DeleteUser(ctx context.Context, input *UserIDPath) (*MessageResponse, error) {
	if err := h.userService.DeleteUser(ctx, input.ID); err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			return nil, huma.Error404NotFound("Пользователь не найден")
		}
		return nil, huma.Error500InternalServerError("Ошибка сервера")
	}

	resp := &MessageResponse{}
	resp.Body.Message = "Пользователь успешно удален"
	return resp, nil
}

func (h *UserHandler) ResetUser(ctx context.Context, input *UserIDPath) (*dto.DevResponse, error) {
	if err := h.userService.ResetUser(ctx, input.ID); err != nil {
		return nil, huma.Error500InternalServerError("Ошибка сервера")
	}

	return &dto.DevResponse{
		Message: "Пользователь успешно сброшен к заводским настройкам",
	}, nil
}

func (h *UserHandler) RestoreUser(ctx context.Context, input *UserIDPath) (*dto.DevResponse, error) {
	if err := h.userService.RestoreUser(ctx, input.ID); err != nil {
		return nil, huma.Error500InternalServerError("Ошибка сервера")
	}

	return &dto.DevResponse{
		Message: "Пользователь успешно восстановлен",
	}, nil
}

func (h *UserHandler) GetDeletedUsers(ctx context.Context, input *struct{}) (*GetDeletedUsersResponse, error) {
	users, err := h.userService.GetDeletedUsers(ctx)
	if err != nil {
		return nil, huma.Error500InternalServerError("Ошибка сервера")
	}

	userDTOs := make([]dto.UserResponseDTO, len(users))
	for i, user := range users {
		userDTOs[i] = mapUserToDTO(user)
	}

	return &GetDeletedUsersResponse{
		Message: "Удаленные пользователи успешно получены",
		Users:   userDTOs,
	}, nil
}
