package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/artesipov-alt/odnoi-krovi-app/ent"
	"github.com/artesipov-alt/odnoi-krovi-app/ent/user"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/dto"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/services"
	"github.com/danielgtaylor/huma/v2"
	"github.com/jinzhu/copier"
)

// UserService определяет интерфейс для бизнес-логики пользователей
type UserService interface {
	// RegisterUser регистрирует нового пользователя в системе
	RegisterUser(ctx context.Context, user *ent.CreateUserInput) (*ent.User, error)

	// RegisterUserSimple создает нового пользователя с Telegram ID и базовой информацией (для команды Start)
	RegisterUserSimple(ctx context.Context, user *ent.CreateUserInput) (*ent.User, error)

	// GetUserByID получает пользователя по его внутреннему ID
	GetUserByID(ctx context.Context, userID string, opt services.UserPreloadOptions) (*ent.User, error)

	// GetUserByTelegramID получает пользователя по Telegram ID
	GetUserByTelegramID(ctx context.Context, telegramID int64, opts services.UserPreloadOptions) (*ent.User, error)

	// Update обновляет информацию о пользователе
	Update(ctx context.Context, id string, input *ent.UpdateUserInput) error

	// DeleteUser удаляет пользователя по ID (soft delete)
	DeleteUser(ctx context.Context, userID string) error

	// ResetUser сбрасывает пользователя к начальным настройкам
	ResetUser(ctx context.Context, userID string) error

	// RestoreUser восстанавливает удаленного пользователя
	RestoreUser(ctx context.Context, userID string) error

	// GetDeletedUsers получает всех удаленных пользователей
	GetDeletedUsers(ctx context.Context) ([]*ent.User, error)

	// BuildFullPhotoURLs преобразует пути к фото в полные публичные URL
	BuildFullPhotoURLs(paths []string) []string
}

// UserHandler обрабатывает HTTP запросы для операций с пользователями
type UserHandler struct {
	userService UserService
}

// NewUserHandler создает новый обработчик пользователей
func NewUserHandler(userService UserService) *UserHandler {
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

func (h *UserHandler) GetUser(ctx context.Context, input *struct {
	dto.IDPathStr
	dto.UserPreloadQuery
}) (*dto.UserResponse, error) {
	slog.DebugContext(ctx, "getting user", "user_id", input.ID)

	usr, err := h.userService.GetUserByID(ctx, input.ID, services.UserPreloadOptions{
		WithPets: input.WithPets,
	})

	if err != nil {
		return nil, err
	}

	return &dto.UserResponse{Body: h.toDTO(usr)}, nil
}

func (h *UserHandler) RegisterUserSimple(ctx context.Context, input *struct {
	Body dto.UserRegistrationSimple
}) (*dto.UserResponse, error) {
	slog.DebugContext(ctx, "registering user simple", "telegram_id", input.Body.TelegramID)

	var userData ent.CreateUserInput

	if err := copier.Copy(&userData, &input.Body); err != nil {
		return nil, apperrors.Internal(errors.New("failed to copy user data"), "ошибка при копировании данных пользователя")
	}

	u, err := h.userService.RegisterUserSimple(ctx, &userData)
	if err != nil {
		return nil, err
	}

	if u == nil {
		return nil, apperrors.Internal(errors.New("registration returned nil user"), "ошибка при создании пользователя")
	}

	return &dto.UserResponse{Body: h.toDTO(u)}, nil
}

func (h *UserHandler) UpdateUser(ctx context.Context, input *struct {
	dto.IDPathStr
	Body dto.UserUpdate
}) (*dto.MessageResponse, error) {
	slog.DebugContext(ctx, "updating user", "user_id", input.ID)

	var user ent.UpdateUserInput

	if err := copier.Copy(&user, &input.Body); err != nil {
		return nil, apperrors.Internal(err, "failed to copy user update data")
	}

	if err := h.userService.Update(ctx, input.ID, &user); err != nil {
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

	usr, err := h.userService.GetUserByTelegramID(ctx, input.ID, services.UserPreloadOptions{
		WithPets: input.WithPets,
	})

	if err != nil {
		return nil, err
	}

	return &dto.UserResponse{Body: h.toDTO(usr)}, nil
}

func (h *UserHandler) DeleteUser(ctx context.Context, input *dto.IDPathStr) (*dto.MessageResponse, error) {
	slog.DebugContext(ctx, "deleting user", "user_id", input.ID)
	if err := h.userService.DeleteUser(ctx, input.ID); err != nil {
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
	if err := h.userService.ResetUser(ctx, input.ID); err != nil {
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
	if err := h.userService.RestoreUser(ctx, input.ID); err != nil {
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
	users, err := h.userService.GetDeletedUsers(ctx)
	if err != nil {
		return nil, err
	}

	userDTOs := make([]dto.User, len(users))
	for i, u := range users {
		userDTOs[i] = h.toDTO(u)
	}

	return &dto.UsersDeletedResponse{
		Body: dto.UsersDeletedBody{
			Message: "Удаленные пользователи получены",
			Users:   userDTOs,
		},
	}, nil
}

// toDTO преобразует ENT модель пользователя в DTO для ответа
func (h *UserHandler) toDTO(u *ent.User) dto.User {
	if u == nil {
		return dto.User{}
	}

	userDTO := dto.User{
		ID:               u.ID,
		TelegramID:       u.TelegramID,
		FullName:         u.FullName,
		Phone:            u.Phone,
		Email:            u.Email,
		PhotoURLs:        u.PhotoUrls,
		OrganizationName: u.OrganizationName,
		ConsentPd:        u.ConsentPd,
		OnBoarding:       u.OnBoarding,
		AllowGeo:         u.AllowGeo,
		LocationID:       u.LocationID,
		Role:             string(u.Role),
		CreatedAt:        &u.CreatedAt,
		UpdatedAt:        &u.UpdatedAt,
		DeletedAt:        u.DeletedAt,
	}

	if u.Edges.Pets != nil {
		userDTO.Pets = make([]dto.Pet, len(u.Edges.Pets))
		for i, pet := range u.Edges.Pets {
			userDTO.Pets[i] = dto.Pet{
				ID:              pet.ID,
				Name:            pet.Name,
				ChipNumber:      pet.ChipNumber,
				PhotoURLs:       pet.PhotoUrls,
				BreedID:         pet.BreedID,
				WeightKg:        pet.WeightKg,
				BirthDate:       pet.BirthDate,
				LivingCondition: pet.LivingCondition.String(),
				Gender:          pet.Gender.String(),
				Type:            pet.Type.String(),
				BloodGroup:      pet.BloodGroupID,
				// PetStatus:       pet.PetStatus.String(),
				CreatedAt: &pet.CreatedAt,
				UpdatedAt: &pet.UpdatedAt,
			}
		}
	}

	return userDTO
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
		PhotoUrls:        u.PhotoURLs,
		OrganizationName: u.OrganizationName,
		ConsentPd:        u.ConsentPd,
		OnBoarding:       u.OnBoarding,
		AllowGeo:         u.AllowGeo,
		LocationID:       u.LocationID,
		Role:             user.Role(u.Role),
	}
}
