package services

import (
	"context"
	"fmt"

	"github.com/artesipov-alt/odnoi-krovi-app/ent"
	"github.com/artesipov-alt/odnoi-krovi-app/ent/user"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/dto"
	repositories "github.com/artesipov-alt/odnoi-krovi-app/internal/repositories"
)

// UserService определяет интерфейс для бизнес-логики пользователей
type UserService interface {
	// RegisterUser регистрирует нового пользователя в системе
	RegisterUser(ctx context.Context, telegramID int64, userData dto.UserRegistrationFull) (*ent.User, error)

	// RegisterUserSimple создает нового пользователя с Telegram ID и базовой информацией (для команды Start)
	RegisterUserSimple(ctx context.Context, userData dto.UserRegistrationSimple) (*ent.User, error)

	// UpdateUserProfile обновляет информацию о пользователе
	UpdateUserProfile(ctx context.Context, userID string, updates dto.UserUpdate) error

	// GetUserByID получает пользователя по его внутреннему ID
	GetUserByID(ctx context.Context, userID string) (*ent.User, error)

	// GetUserByTelegramID получает пользователя по Telegram ID
	GetUserByTelegramID(ctx context.Context, telegramID int64) (*ent.User, error)

	// DeleteUser удаляет пользователя по ID (soft delete)
	DeleteUser(ctx context.Context, userID string) error
}

// UserServiceImpl реализует UserService
type UserServiceImpl struct {
	userRepo     repositories.UserRepository
	locationRepo repositories.LocationRepository
}

// NewUserService создает новый сервис пользователей
func NewUserService(userRepo repositories.UserRepository, locationRepo repositories.LocationRepository) *UserServiceImpl {
	return &UserServiceImpl{
		userRepo:     userRepo,
		locationRepo: locationRepo,
	}
}

// RegisterUser регистрирует нового пользователя в системе
func (s *UserServiceImpl) RegisterUser(ctx context.Context, telegramID int64, userData dto.UserRegistrationFull) (*ent.User, error) {
	// Проверяем, существует ли пользователь уже
	exists, err := s.userRepo.ExistsByTelegramID(ctx, telegramID)
	if err != nil {
		return nil, apperrors.Internal(err, "не удалось проверить существование пользователя")
	}

	if exists {
		return nil, apperrors.NewUserAlreadyExistsError(telegramID)
	}

	// Валидируем роль пользователя через ENT-валидатор
	if err := user.RoleValidator(user.Role(userData.Role)); err != nil {
		return nil, apperrors.ErrUserInvalidRole
	}

	// Проверяем существование локации
	_, err = s.locationRepo.GetByID(ctx, userData.LocationID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, apperrors.BadRequest(fmt.Sprintf("локация с ID %d не существует", userData.LocationID))
		}
		return nil, apperrors.Internal(err, "не удалось проверить существование локации")
	}

	// Создаем нового пользователя
	u := &ent.User{
		TelegramID: telegramID,
		FullName:   userData.FullName,
		Phone:      userData.Phone,
		Email:      userData.Email,
		ConsentPd:  userData.ConsentPD,
		LocationID: userData.LocationID,
		Role:       user.Role(userData.Role),
	}

	newUser, err := s.userRepo.Create(ctx, u)
	if err != nil {
		return nil, apperrors.Internal(err, "не удалось создать пользователя")
	}

	return newUser, nil
}

// RegisterUserSimple создает нового пользователя с Telegram ID и базовой информацией (для команды Start)
func (s *UserServiceImpl) RegisterUserSimple(ctx context.Context, userdata dto.UserRegistrationSimple) (*ent.User, error) {
	// Проверяем, существует ли пользователь уже
	exists, err := s.userRepo.ExistsByTelegramID(ctx, userdata.TelegramID)
	if err != nil {
		return nil, apperrors.Internal(err, "не удалось проверить существование пользователя")
	}

	if exists {
		return nil, apperrors.NewUserAlreadyExistsError(userdata.TelegramID)
	}

	// Создаем нового пользователя с Telegram ID, базовой информацией и значениями по умолчанию
	u := &ent.User{
		TelegramID: userdata.TelegramID,
		FullName:   userdata.FullName,
		Phone:      "",
		Email:      "",
		ConsentPd:  true,
		OnBoarding: false,
		AllowGeo:   false,
		Role:       user.RoleUser,
	}

	newUser, err := s.userRepo.Create(ctx, u)
	if err != nil {
		return nil, apperrors.Internal(err, "не удалось создать пользователя")
	}

	return newUser, nil
}

// DeleteUser удаляет пользователя по ID (soft delete)
func (s *UserServiceImpl) DeleteUser(ctx context.Context, userID string) error {
	// Проверяем, существует ли пользователь
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if ent.IsNotFound(err) {
			return apperrors.NewUserNotFoundError(userID)
		}
		return apperrors.Internal(err, "не удалось получить пользователя")
	}

	// Удаляем пользователя
	if err := s.userRepo.Delete(ctx, userID); err != nil {
		return apperrors.Internal(err, "не удалось удалить пользователя")
	}

	return nil
}

// GetUserByID получает пользователя по ID
func (s *UserServiceImpl) GetUserByID(ctx context.Context, userID string) (*ent.User, error) {
	u, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, apperrors.NewUserNotFoundError(userID)
		}
		return nil, apperrors.Internal(err, "не удалось получить пользователя")
	}

	return u, nil
}

// UpdateUserProfile обновляет информацию о пользователе
func (s *UserServiceImpl) UpdateUserProfile(ctx context.Context, userID string, updates dto.UserUpdate) error {
	// Получаем существующего пользователя
	u, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if ent.IsNotFound(err) {
			return apperrors.NewUserNotFoundError(userID)
		}
		return apperrors.Internal(err, "не удалось получить пользователя")
	}

	// Применяем обновления
	if updates.FullName != nil {
		u.FullName = *updates.FullName
	}
	if updates.Phone != nil {
		u.Phone = *updates.Phone
	}
	if updates.Email != nil {
		u.Email = *updates.Email
	}
	if updates.AllowGeo != nil {
		u.AllowGeo = *updates.AllowGeo
	}
	if updates.OnBoarding != nil {
		u.OnBoarding = *updates.OnBoarding
	}
	if updates.LocationID != nil {
		// Проверяем существование локации
		_, err := s.locationRepo.GetByID(ctx, *updates.LocationID)
		if err != nil {
			if ent.IsNotFound(err) {
				return apperrors.BadRequest(fmt.Sprintf("локация с ID %d не существует", *updates.LocationID))
			}
			return apperrors.Internal(err, "не удалось проверить существование локации")
		}
		u.LocationID = *updates.LocationID
	}

	// Сохраняем обновленного пользователя
	if _, err := s.userRepo.Update(ctx, u); err != nil {
		return apperrors.Internal(err, "не удалось обновить пользователя")
	}

	return nil
}

// GetUserByTelegramID получает пользователя по Telegram ID
func (s *UserServiceImpl) GetUserByTelegramID(ctx context.Context, telegramID int64) (*ent.User, error) {
	u, err := s.userRepo.GetByTelegramID(ctx, telegramID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, apperrors.NotFound("пользователь с таким Telegram ID не найден").WithDetails(map[string]any{
				"telegram_id": telegramID,
			})
		}
		return nil, apperrors.Internal(err, "не удалось получить пользователя по Telegram ID")
	}

	return u, nil
}
