package http

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/application/user/cmd"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/application/user/query"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/filestorage"
	usermodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/mapper"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/transport/http/dto"
	"github.com/danielgtaylor/huma/v2"
)

// UserHandler обрабатывает HTTP запросы для операций с пользователями
type UserHandler struct {
	deleteHandler        *cmd.DeleteHandler
	updateHandler        *cmd.UpdateHandler
	resetHandler         *cmd.ResetHandler
	restoreHandler       *cmd.RestoreHandler
	getByIDHandler       *query.GetByIDHandler
	getByTelegramHandler *query.GetByTelegramHandler
	getDeletedHandler    *query.GetDeletedUsersHandler
	userMapper           *mapper.UserMapper
	storage              filestorage.Repository
}

// NewUserHandler создает новый обработчик пользователей
func NewUserHandler(
	deleteHandler *cmd.DeleteHandler,
	updateHandler *cmd.UpdateHandler,
	resetHandler *cmd.ResetHandler,
	restoreHandler *cmd.RestoreHandler,
	getByIDHandler *query.GetByIDHandler,
	getByTelegramHandler *query.GetByTelegramHandler,
	getDeletedHandler *query.GetDeletedUsersHandler,
	storage filestorage.Repository,
) *UserHandler {
	return &UserHandler{
		deleteHandler:        deleteHandler,
		updateHandler:        updateHandler,
		resetHandler:         resetHandler,
		restoreHandler:       restoreHandler,
		getByIDHandler:       getByIDHandler,
		getByTelegramHandler: getByTelegramHandler,
		getDeletedHandler:    getDeletedHandler,
		userMapper:           mapper.NewUserMapper(mapper.NewPetMapper(storage), storage),
		storage:              storage,
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
		Deprecated:  true,
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

func (h *UserHandler) GetUser(ctx context.Context, input *dto.GetUserByIDInput) (*dto.GetUserByIDOutput, error) {
	slog.DebugContext(ctx, "getting user", "user_id", input.ID)

	usr, err := h.getByIDHandler.Handle(ctx, input.ID, input.WithPets, input.WithDonorPreference)
	if err != nil {
		return nil, err
	}

	return &dto.GetUserByIDOutput{Body: h.userMapper.ToResponse(usr)}, nil
}

func (h *UserHandler) UpdateUser(ctx context.Context, input *dto.UpdateUserInput) (*dto.UpdateUserOutput, error) {
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
	if input.Body.AllowGeo != nil {
		user.AllowGeo = *input.Body.AllowGeo
	}
	if input.Body.OnBoarding != nil {
		user.OnBoarding = *input.Body.OnBoarding
	}
	if input.Body.LocationID != nil {
		user.LocationID = input.Body.LocationID
	}
	if input.Body.DonorPreference != nil {
		user.DonorPreference = &usermodel.DonorPreference{
			PreferredLocationIDs: input.Body.DonorPreference.PreferredLocationIDs,
		}
		if input.Body.DonorPreference.RecoveryPeriodMonths != nil {
			user.DonorPreference.RecoveryPeriodMonths = *input.Body.DonorPreference.RecoveryPeriodMonths
		}
		if input.Body.DonorPreference.CompensationType != nil {
			user.DonorPreference.CompensationType = usermodel.CompensationType(*input.Body.DonorPreference.CompensationType)
		}
		if input.Body.DonorPreference.TaxiCompensation != nil {
			user.DonorPreference.TaxiCompensation = *input.Body.DonorPreference.TaxiCompensation
		}
		if input.Body.DonorPreference.NotificationFrequency != nil {
			user.DonorPreference.NotificationFrequency = usermodel.NotificationFrequency(*input.Body.DonorPreference.NotificationFrequency)
		}
	}

	user, err := h.updateHandler.Handle(ctx, input.ID, user)
	if err != nil {
		return nil, err
	}

	return &dto.UpdateUserOutput{Body: dto.UpdateUserResult{
		ID:        user.ID,
		UpdatedAt: user.UpdatedAt,
	}}, nil
}

func (h *UserHandler) UserByTelegram(ctx context.Context, input *dto.GetUserByTelegramInput) (*dto.GetUserByTelegramOutput, error) {
	slog.DebugContext(ctx, "getting user by telegram", "telegram_id", input.ID)

	usr, err := h.getByTelegramHandler.Handle(ctx, input.ID, input.WithPets, input.WithDonorPreference)
	if err != nil {
		return nil, err
	}

	return &dto.GetUserByTelegramOutput{Body: h.userMapper.ToResponse(usr)}, nil
}

func (h *UserHandler) DeleteUser(ctx context.Context, input *dto.DeleteUserInput) (*dto.DeleteUserOutput, error) {
	slog.DebugContext(ctx, "deleting user", "user_id", input.ID)
	if err := h.deleteHandler.Handle(ctx, input.ID); err != nil {
		return nil, err
	}

	return &dto.DeleteUserOutput{Body: dto.DeleteUserResult{
		Message: "Пользователь удален",
	}}, nil
}

func (h *UserHandler) ResetUser(ctx context.Context, input *dto.ResetUserInput) (*dto.ResetUserOutput, error) {
	slog.DebugContext(ctx, "resetting user", "user_id", input.ID)
	if err := h.resetHandler.Handle(ctx, input.ID); err != nil {
		return nil, err
	}

	return &dto.ResetUserOutput{Body: dto.ResetUserResult{
		Message: "Пользователь сброшен к заводским настройкам",
	}}, nil
}

func (h *UserHandler) RestoreUser(ctx context.Context, input *dto.RestoreUserInput) (*dto.RestoreUserOutput, error) {
	slog.DebugContext(ctx, "restoring user", "user_id", input.ID)
	if err := h.restoreHandler.Handle(ctx, input.ID); err != nil {
		return nil, err
	}

	return &dto.RestoreUserOutput{Body: dto.RestoreUserResult{
		Message: "Пользователь восстановлен",
	}}, nil
}

func (h *UserHandler) DeletedUsers(ctx context.Context, input *dto.GetDeletedUsersInput) (*dto.GetDeletedUsersOutput, error) {
	slog.DebugContext(ctx, "getting deleted users")
	users, err := h.getDeletedHandler.Handle(ctx)
	if err != nil {
		return nil, err
	}

	return &dto.GetDeletedUsersOutput{Body: dto.DeletedUsersList{
		Message: "Удаленные пользователи получены",
		Users:   h.userMapper.ToResponseSlice(users),
	}}, nil
}
