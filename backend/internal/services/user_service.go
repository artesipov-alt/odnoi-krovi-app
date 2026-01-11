package services

import (
	"context"
	"fmt"

	"github.com/artesipov-alt/odnoi-krovi-app/ent"
	userval "github.com/artesipov-alt/odnoi-krovi-app/ent/user"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	repositories "github.com/artesipov-alt/odnoi-krovi-app/internal/repositories"
)

// UserService определяет интерфейс для бизнес-логики пользователей
type UserService interface {
	// RegisterUser регистрирует нового пользователя в системе
	RegisterUser(ctx context.Context, user *ent.User) (*ent.User, error)

	// RegisterUserSimple создает нового пользователя с Telegram ID и базовой информацией (для команды Start)
	RegisterUserSimple(ctx context.Context, user *ent.User) (*ent.User, error)

	// UpdateUserProfile обновляет информацию о пользователе
	UpdateUserProfile(ctx context.Context, userID string, updates map[string]any) error

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
func (s *UserServiceImpl) RegisterUser(ctx context.Context, user *ent.User) (*ent.User, error) {
	// Проверяем, существует ли пользователь уже
	exists, err := s.userRepo.ExistsByTelegramID(ctx, user.TelegramID)
	if err != nil {
		return nil, apperrors.Internal(err, "не удалось проверить существование пользователя")
	}

	if exists {
		return nil, apperrors.NewUserAlreadyExistsError(user.TelegramID)
	}

	// Валидируем роль пользователя через ENT-валидатор
	if err := userval.RoleValidator(user.Role); err != nil {
		return nil, apperrors.ErrUserInvalidRole
	}

	// Проверяем существование локации
	_, err = s.locationRepo.GetByID(ctx, user.LocationID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, apperrors.BadRequest(fmt.Sprintf("локация с ID %d не существует", user.LocationID))
		}
		return nil, apperrors.Internal(err, "не удалось проверить существование локации")
	}

	newUser, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, apperrors.Internal(err, "не удалось создать пользователя")
	}

	return newUser, nil
}

// RegisterUserSimple создает нового пользователя с Telegram ID и базовой информацией (для команды Start)
func (s *UserServiceImpl) RegisterUserSimple(ctx context.Context, user *ent.User) (*ent.User, error) {
	// Проверяем, существует ли пользователь уже
	exists, err := s.userRepo.ExistsByTelegramID(ctx, user.TelegramID)
	if err != nil {
		return nil, apperrors.Internal(err, "не удалось проверить существование пользователя")
	}

	if exists {
		return nil, apperrors.NewUserAlreadyExistsError(user.TelegramID)
	}

	newUser, err := s.userRepo.Create(ctx, user)
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
func (s *UserServiceImpl) UpdateUserProfile(ctx context.Context, userID string, updates map[string]interface{}) error {
	// Получаем существующего пользователя
	u, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if ent.IsNotFound(err) {
			return apperrors.NewUserNotFoundError(userID)
		}
		return apperrors.Internal(err, "не удалось получить пользователя")
	}

	// Применяем обновления
	if val, ok := updates["FullName"]; ok {
		u.FullName = val.(string)
	}
	if val, ok := updates["Phone"]; ok {
		u.Phone = val.(string)
	}
	if val, ok := updates["Email"]; ok {
		u.Email = val.(string)
	}
	if val, ok := updates["AllowGeo"]; ok {
		u.AllowGeo = val.(bool)
	}
	if val, ok := updates["OnBoarding"]; ok {
		u.OnBoarding = val.(bool)
	}
	if val, ok := updates["LocationID"]; ok {
		locationID := val.(int)
		// Проверяем существование локации
		_, err := s.locationRepo.GetByID(ctx, locationID)
		if err != nil {
			if ent.IsNotFound(err) {
				return apperrors.BadRequest(fmt.Sprintf("локация с ID %d не существует", locationID))
			}
			return apperrors.Internal(err, "не удалось проверить существование локации")
		}
		u.LocationID = locationID
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
