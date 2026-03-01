package http

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/application/user/cmd"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/application/user/query"
	usermodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/mapper"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/transport/http/dto"
	"github.com/danielgtaylor/huma/v2"
)

// UserHandler обрабатывает HTTP запросы для операций с пользователями
type UserHandler struct {
	createSimpleHandler  *cmd.CreateSimpleHandler
	deleteHandler        *cmd.DeleteHandler
	updateHandler        *cmd.UpdateHandler
	resetHandler         *cmd.ResetHandler
	restoreHandler       *cmd.RestoreHandler
	getByIDHandler       *query.GetByIDHandler
	getByTelegramHandler *query.GetByTelegramHandler
	getDeletedHandler    *query.GetDeletedUsersHandler
	userMapper           *mapper.UserMapper
}

// NewUserHandler создает новый обработчик пользователей
func NewUserHandler(
	createSimpleHandler *cmd.CreateSimpleHandler,
	deleteHandler *cmd.DeleteHandler,
	updateHandler *cmd.UpdateHandler,
	resetHandler *cmd.ResetHandler,
	restoreHandler *cmd.RestoreHandler,
	getByIDHandler *query.GetByIDHandler,
	getByTelegramHandler *query.GetByTelegramHandler,
	getDeletedHandler *query.GetDeletedUsersHandler,
) *UserHandler {
	return &UserHandler{
		createSimpleHandler:  createSimpleHandler,
		deleteHandler:        deleteHandler,
		updateHandler:        updateHandler,
		resetHandler:         resetHandler,
		restoreHandler:       restoreHandler,
		getByIDHandler:       getByIDHandler,
		getByTelegramHandler: getByTelegramHandler,
		getDeletedHandler:    getDeletedHandler,
		userMapper:           mapper.NewUserMapper(mapper.NewPetMapper()),
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
		Path:        "/v1/user/telegram/{id}",
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
		Tags:        []string{"dev"},
	}, h.ResetUser)

	// Восстановление удаленного пользователя
	huma.Register(api, huma.Operation{
		OperationID: "restore-user",
		Method:      http.MethodPost,
		Path:        "/v1/user/restore-user/{id}",
		Summary:     "Восстановление удаленного пользователя",
		Description: "Восстанавливает мягко удаленного пользователя, устанавливая deleted_at в NULL",
		Tags:        []string{"dev"},
	}, h.RestoreUser)

	// Получение всех удаленных пользователей
	huma.Register(api, huma.Operation{
		OperationID: "get-deleted-users",
		Method:      http.MethodGet,
		Path:        "/v1/user/deleted-users",
		Summary:     "Получение всех удаленных пользователей",
		Description: "Возвращает список всех мягко удаленных пользователей",
		Tags:        []string{"dev"},
	}, h.DeletedUsers)
}

// Handlers

func (h *UserHandler) GetUser(ctx context.Context, input *struct {
	dto.IDPathStr
	dto.UserPreloadQuery
}) (*dto.UserResponse, error) {
	slog.DebugContext(ctx, "getting user", "user_id", input.ID)

	usr, pets, err := h.getByIDHandler.Handle(ctx, input.ID, input.WithPets)
	if err != nil {
		return nil, err
	}

	return &dto.UserResponse{Body: h.userMapper.ToDTO(usr, pets)}, nil
}

func (h *UserHandler) RegisterUserSimple(ctx context.Context, input *struct {
	Body dto.UserRegistrationSimple
}) (*dto.UserResponse, error) {
	slog.DebugContext(ctx, "registering user simple", "telegram_id", input.Body.TelegramID)

	u, err := h.createSimpleHandler.Handle(ctx, input.Body.TelegramID, input.Body.FullName, "user")
	if err != nil {
		return nil, err
	}

	if u == nil {
		return nil, apperrors.Internal(nil, "ошибка при создании пользователя")
	}

	return &dto.UserResponse{Body: h.userMapper.ToDTO(u, nil)}, nil
}

func (h *UserHandler) UpdateUser(ctx context.Context, input *struct {
	dto.IDPathStr
	Body dto.UserUpdate
}) (*dto.MessageResponse, error) {
	slog.DebugContext(ctx, "updating user", "user_id", input.ID)

	user := &usermodel.User{}

	if input.Body.FullName != nil {
		user.FullName = *input.Body.FullName
	}
	if input.Body.Phone != nil {
		user.Phone = *input.Body.Phone
	}
	if input.Body.Email != nil {
		user.Email = *input.Body.Email
	}
	if input.Body.PhotoURLs != nil {
		user.PhotoURLs = input.Body.PhotoURLs
	}
	if input.Body.AllowGeo != nil {
		user.AllowGeo = *input.Body.AllowGeo
	}
	if input.Body.OnBoarding != nil {
		user.OnBoarding = *input.Body.OnBoarding
	}
	if input.Body.LocationID != nil {
		user.LocationID = input.Body.LocationID
	}

	if err := h.updateHandler.Handle(ctx, input.ID, user); err != nil {
		return nil, err
	}

	return &dto.MessageResponse{
		Body: dto.MessageBody{
			Message: "Пользователь обновлен",
		},
	}, nil
}

func (h *UserHandler) UserByTelegram(ctx context.Context, input *struct {
	dto.IDPathInt
	dto.UserPreloadQuery
}) (*dto.UserResponse, error) {
	slog.DebugContext(ctx, "getting user by telegram", "telegram_id", input.ID)

	usr, pets, err := h.getByTelegramHandler.Handle(ctx, input.ID, input.WithPets)
	if err != nil {
		return nil, err
	}

	return &dto.UserResponse{Body: h.userMapper.ToDTO(usr, pets)}, nil
}

func (h *UserHandler) DeleteUser(ctx context.Context, input *dto.IDPathStr) (*dto.MessageResponse, error) {
	slog.DebugContext(ctx, "deleting user", "user_id", input.ID)
	if err := h.deleteHandler.Handle(ctx, input.ID); err != nil {
		return nil, err
	}

	return &dto.MessageResponse{
		Body: dto.MessageBody{
			Message: "Пользователь удален",
		},
	}, nil
}

func (h *UserHandler) ResetUser(ctx context.Context, input *dto.IDPathStr) (*dto.MessageResponse, error) {
	slog.DebugContext(ctx, "resetting user", "user_id", input.ID)
	if err := h.resetHandler.Handle(ctx, input.ID); err != nil {
		return nil, err
	}

	return &dto.MessageResponse{
		Body: dto.MessageBody{
			Message: "Пользователь сброшен к заводским настройкам",
		},
	}, nil
}

func (h *UserHandler) RestoreUser(ctx context.Context, input *dto.IDPathStr) (*dto.MessageResponse, error) {
	slog.DebugContext(ctx, "restoring user", "user_id", input.ID)
	if err := h.restoreHandler.Handle(ctx, input.ID); err != nil {
		return nil, err
	}

	return &dto.MessageResponse{
		Body: dto.MessageBody{
			Message: "Пользователь восстановлен",
		},
	}, nil
}

func (h *UserHandler) DeletedUsers(ctx context.Context, input *struct{}) (*dto.UsersDeletedResponse, error) {
	slog.DebugContext(ctx, "getting deleted users")
	users, err := h.getDeletedHandler.Handle(ctx)
	if err != nil {
		return nil, err
	}

	return &dto.UsersDeletedResponse{
		Body: dto.UsersDeletedBody{
			Message: "Удаленные пользователи получены",
			Users:   h.userMapper.ToDTOs(users),
		},
	}, nil
}
