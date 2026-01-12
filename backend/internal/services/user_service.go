package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/artesipov-alt/odnoi-krovi-app/ent"
	userval "github.com/artesipov-alt/odnoi-krovi-app/ent/user"
	repositories "github.com/artesipov-alt/odnoi-krovi-app/internal/repositories"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrLocationNotFound  = errors.New("location not found")
	ErrInvalidRole       = errors.New("invalid role")
	ErrInternal          = errors.New("internal error")
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

	// ResetUser сбрасывает пользователя к начальным настройкам
	ResetUser(ctx context.Context, userID string) error

	// RestoreUser восстанавливает удаленного пользователя
	RestoreUser(ctx context.Context, userID string) error

	// GetDeletedUsers получает всех удаленных пользователей
	GetDeletedUsers(ctx context.Context) ([]*ent.User, error)
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
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}

	if exists {
		return nil, ErrUserAlreadyExists
	}

	// Валидируем роль пользователя через ENT-валидатор
	if err := userval.RoleValidator(user.Role); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidRole, err)
	}

	// Проверяем существование локации
	_, err = s.locationRepo.GetByID(ctx, user.LocationID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrLocationNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}

	newUser, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}

	return newUser, nil
}

// RegisterUserSimple создает нового пользователя с Telegram ID и базовой информацией (для команды Start)
func (s *UserServiceImpl) RegisterUserSimple(ctx context.Context, user *ent.User) (*ent.User, error) {
	// Проверяем, существует ли пользователь уже
	exists, err := s.userRepo.ExistsByTelegramID(ctx, user.TelegramID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}

	if exists {
		return nil, ErrUserAlreadyExists
	}

	newUser, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}

	return newUser, nil
}

// DeleteUser удаляет пользователя по ID (soft delete)
func (s *UserServiceImpl) DeleteUser(ctx context.Context, userID string) error {
	// Проверяем, существует ли пользователь
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if ent.IsNotFound(err) {
			return ErrUserNotFound
		}
		return fmt.Errorf("%w: %v", ErrInternal, err)
	}

	// Удаляем пользователя
	if err := s.userRepo.Delete(ctx, userID); err != nil {
		return fmt.Errorf("%w: %v", ErrInternal, err)
	}

	return nil
}

// GetUserByID получает пользователя по ID
func (s *UserServiceImpl) GetUserByID(ctx context.Context, userID string) (*ent.User, error) {
	u, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}

	return u, nil
}

// UpdateUserProfile обновляет информацию о пользователе
func (s *UserServiceImpl) UpdateUserProfile(ctx context.Context, userID string, updates map[string]interface{}) error {
	// Получаем существующего пользователя
	u, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if ent.IsNotFound(err) {
			return ErrUserNotFound
		}
		return fmt.Errorf("%w: %v", ErrInternal, err)
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
				return ErrLocationNotFound
			}
			return fmt.Errorf("%w: %v", ErrInternal, err)
		}
		u.LocationID = locationID
	}

	// Сохраняем обновленного пользователя
	if _, err := s.userRepo.Update(ctx, u); err != nil {
		return fmt.Errorf("%w: %v", ErrInternal, err)
	}

	return nil
}

// GetUserByTelegramID получает пользователя по Telegram ID
func (s *UserServiceImpl) GetUserByTelegramID(ctx context.Context, telegramID int64) (*ent.User, error) {
	u, err := s.userRepo.GetByTelegramID(ctx, telegramID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}

	return u, nil
}

// ResetUser сбрасывает пользователя к начальным настройкам
func (s *UserServiceImpl) ResetUser(ctx context.Context, userID string) error {
	if err := s.userRepo.ResetUser(ctx, userID); err != nil {
		return fmt.Errorf("%w: %v", ErrInternal, err)
	}
	return nil
}

// RestoreUser восстанавливает удаленного пользователя
func (s *UserServiceImpl) RestoreUser(ctx context.Context, userID string) error {
	if err := s.userRepo.RestoreUser(ctx, userID); err != nil {
		return fmt.Errorf("%w: %v", ErrInternal, err)
	}
	return nil
}

// GetDeletedUsers получает всех удаленных пользователей
func (s *UserServiceImpl) GetDeletedUsers(ctx context.Context) ([]*ent.User, error) {
	users, err := s.userRepo.GetDeletedUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}
	return users, nil
}
