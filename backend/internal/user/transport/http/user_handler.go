package user

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/application/user/cmd"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/application/user/query"
	petdto "github.com/artesipov-alt/odnoi-krovi-app/internal/pet/dto"
	petmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/pet/model"
	userdto "github.com/artesipov-alt/odnoi-krovi-app/internal/user/dto"
	usermodel "github.com/artesipov-alt/odnoi-krovi-app/internal/user/model"
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
	userdto.IDPathStr
	userdto.UserPreloadQuery
}) (*userdto.UserResponse, error) {
	slog.DebugContext(ctx, "getting user", "user_id", input.ID)

	usr, pets, err := h.getByIDHandler.Handle(ctx, input.ID, input.WithPets)
	if err != nil {
		return nil, err
	}

	return &userdto.UserResponse{Body: h.toDTO(usr, pets)}, nil
}

func (h *UserHandler) RegisterUserSimple(ctx context.Context, input *struct {
	Body userdto.UserRegistrationSimple
}) (*userdto.UserResponse, error) {
	slog.DebugContext(ctx, "registering user simple", "telegram_id", input.Body.TelegramID)

	u, err := h.createSimpleHandler.Handle(ctx, input.Body.TelegramID, input.Body.FullName, "user")
	if err != nil {
		return nil, err
	}

	if u == nil {
		return nil, apperrors.Internal(nil, "ошибка при создании пользователя")
	}

	return &userdto.UserResponse{Body: h.toDTO(u, nil)}, nil
}

func (h *UserHandler) UpdateUser(ctx context.Context, input *struct {
	userdto.IDPathStr
	Body userdto.UserUpdate
}) (*userdto.MessageResponse, error) {
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

	return &userdto.MessageResponse{
		Body: userdto.MessageBody{
			Message: "Пользователь обновлен",
		},
	}, nil
}

func (h *UserHandler) UserByTelegram(ctx context.Context, input *struct {
	userdto.IDPathInt
	userdto.UserPreloadQuery
}) (*userdto.UserResponse, error) {
	slog.DebugContext(ctx, "getting user by telegram", "telegram_id", input.ID)

	usr, pets, err := h.getByTelegramHandler.Handle(ctx, input.ID, input.WithPets)
	if err != nil {
		return nil, err
	}

	return &userdto.UserResponse{Body: h.toDTO(usr, pets)}, nil
}

func (h *UserHandler) DeleteUser(ctx context.Context, input *userdto.IDPathStr) (*userdto.MessageResponse, error) {
	slog.DebugContext(ctx, "deleting user", "user_id", input.ID)
	if err := h.deleteHandler.Handle(ctx, input.ID); err != nil {
		return nil, err
	}

	return &userdto.MessageResponse{
		Body: userdto.MessageBody{
			Message: "Пользователь удален",
		},
	}, nil
}

func (h *UserHandler) ResetUser(ctx context.Context, input *userdto.IDPathStr) (*userdto.MessageResponse, error) {
	slog.DebugContext(ctx, "resetting user", "user_id", input.ID)
	if err := h.resetHandler.Handle(ctx, input.ID); err != nil {
		return nil, err
	}

	return &userdto.MessageResponse{
		Body: userdto.MessageBody{
			Message: "Пользователь сброшен к заводским настройкам",
		},
	}, nil
}

func (h *UserHandler) RestoreUser(ctx context.Context, input *userdto.IDPathStr) (*userdto.MessageResponse, error) {
	slog.DebugContext(ctx, "restoring user", "user_id", input.ID)
	if err := h.restoreHandler.Handle(ctx, input.ID); err != nil {
		return nil, err
	}

	return &userdto.MessageResponse{
		Body: userdto.MessageBody{
			Message: "Пользователь восстановлен",
		},
	}, nil
}

func (h *UserHandler) DeletedUsers(ctx context.Context, input *struct{}) (*userdto.UsersDeletedResponse, error) {
	slog.DebugContext(ctx, "getting deleted users")
	users, err := h.getDeletedHandler.Handle(ctx)
	if err != nil {
		return nil, err
	}

	userDTOs := make([]userdto.User, len(users))
	for i, u := range users {
		userDTOs[i] = h.toDTO(u, nil)
	}

	return &userdto.UsersDeletedResponse{
		Body: userdto.UsersDeletedBody{
			Message: "Удаленные пользователи получены",
			Users:   userDTOs,
		},
	}, nil
}

// toDTO преобразует модель пользователя в DTO для ответа
func (h *UserHandler) toDTO(u *usermodel.User, pets []*petmodel.Pet) userdto.User {
	if u == nil {
		return userdto.User{}
	}

	userDTO := userdto.User{
		ID:               u.ID,
		TelegramID:       u.TelegramID,
		FullName:         u.FullName,
		Phone:            u.Phone,
		Email:            u.Email,
		PhotoURLs:        u.PhotoURLs,
		OrganizationName: u.OrganizationName,
		ConsentPd:        u.ConsentPd,
		OnBoarding:       u.OnBoarding,
		AllowGeo:         u.AllowGeo,
		LocationID:       "",
		Role:             u.Role,
		CreatedAt:        u.CreatedAt,
		UpdatedAt:        u.UpdatedAt,
		DeletedAt:        u.DeletedAt,
	}

	if u.LocationID != nil {
		userDTO.LocationID = *u.LocationID
	}

	if pets != nil {
		userDTO.Pets = make([]petdto.Pet, len(pets))
		for i, pet := range pets {
			userDTO.Pets[i] = h.petToDTO(pet)
		}
	}

	return userDTO
}

// petToDTO преобразует модель питомца в DTO
func (h *UserHandler) petToDTO(pet *petmodel.Pet) petdto.Pet {
	if pet == nil {
		return petdto.Pet{}
	}

	dtoPet := petdto.Pet{
		ID:         pet.ID,
		Name:       pet.Name,
		ChipNumber: pet.ChipNumber,
		PhotoURLs:  pet.PhotoURLs,
		WeightKg:   pet.WeightKg,
		BirthDate:  pet.BirthDate,
		CreatedAt:  pet.CreatedAt,
		UpdatedAt:  pet.UpdatedAt,
	}

	if pet.BreedRefID != nil {
		dtoPet.BreedID = *pet.BreedRefID
	}

	if pet.BloodGroupName != nil {
		dtoPet.BloodGroup = *pet.BloodGroupName
	}

	if pet.LivingCondition != "" {
		dtoPet.LivingCondition = string(pet.LivingCondition)
	}

	if pet.Gender != "" {
		dtoPet.Gender = string(pet.Gender)
	}

	if pet.Type != "" {
		dtoPet.Type = string(pet.Type)
	}

	return dtoPet
}
