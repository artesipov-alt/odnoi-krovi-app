package services

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/ent"
	userval "github.com/artesipov-alt/odnoi-krovi-app/ent/user"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
)

// UserRepository определяет интерфейс для операций с данными пользователей
type UserRepository interface {
	// Create создает нового пользователя в базе данных
	Create(ctx context.Context, user *ent.CreateUserInput) (*ent.User, error)

	// GetQuery возвращает query для eager loading
	GetQueryByID(ctx context.Context, id string) *ent.UserQuery

	// GetQueryByTelegram возвращает query для eager loading по Telegram ID
	GetQueryByTelegram(ctx context.Context, telegramID int64) *ent.UserQuery

	// Update обновляет существующего пользователя в базе данных
	Update(ctx context.Context, id string, input *ent.UpdateUserInput) error

	// Delete удаляет пользователя по его ID
	Delete(ctx context.Context, id string) error

	// ExistsByTelegramID проверяет, существует ли пользователь с заданным Telegram ID
	ExistsByTelegramID(ctx context.Context, telegramID int64) (bool, error)

	// ExistsByID проверяет, существует ли пользователь с заданным ID
	ExistsByID(ctx context.Context, id string) (bool, error)

	// ResetUser сбрасывает email и номер телефона пользователя по ID
	ResetUser(ctx context.Context, id string) error

	// RestoreUser восстанавливает пользователя по его ID
	RestoreUser(ctx context.Context, id string) error

	// GetDeletedUsers получает всех удаленных пользователей
	GetDeletedUsers(ctx context.Context) ([]*ent.User, error)

	// AddPhotoURLs добавляет новые пути к фотографиям пользователя
	AddPhotoURLs(ctx context.Context, id string, paths []string) error
}

// LocationRepository определяет интерфейс для операций с данными локаций
type LocationRepository interface {

	// GetByID получает локацию по её ID
	GetByID(ctx context.Context, id string) (*ent.Location, error)

	// GetAll получает все локации из базы данных
	GetAll(ctx context.Context) ([]*ent.Location, error)

	// Exists проверяет, существует ли локация с заданным ID
	Exists(ctx context.Context, id string) (bool, error)
}

type UserOptions struct {
	WithPets bool
}

// UserServiceImpl реализует UserService
type UserServiceImpl struct {
	userRepo     UserRepository
	locationRepo LocationRepository
	storage      FileStorage
}

// NewUserService создает новый сервис пользователей
// NewUserService создает новый экземпляр UserService
func NewUserService(userRepo UserRepository, locationRepo LocationRepository, storage FileStorage) *UserServiceImpl {
	return &UserServiceImpl{
		userRepo:     userRepo,
		locationRepo: locationRepo,
		storage:      storage,
	}
}

// RegisterUser регистрирует нового пользователя в системе
func (s *UserServiceImpl) RegisterUser(ctx context.Context, input *ent.CreateUserInput) (*ent.User, error) {
	// Проверяем, существует ли пользователь уже
	exists, err := s.userRepo.ExistsByTelegramID(ctx, input.TelegramID)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to check user existence")
	}

	if exists {
		return nil, apperrors.ErrUserAlreadyExists
	}

	// Валидируем роль пользователя через ENT-валидатор
	if err := userval.RoleValidator(*input.Role); err != nil {
		return nil, apperrors.ErrUserInvalidRole.WithInternal(err)
	}

	// Проверяем существование локации
	_, err = s.locationRepo.GetByID(ctx, *input.LocationID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, apperrors.ErrLocationNotFound
		}
		return nil, apperrors.Internal(err, "failed to get location")
	}

	newUser, err := s.userRepo.Create(ctx, input)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to create user")
	}

	return newUser, nil
}

// RegisterUserSimple создает нового пользователя с Telegram ID и базовой информацией (для команды Start)
func (s *UserServiceImpl) RegisterUserSimple(ctx context.Context, input *ent.CreateUserInput) (*ent.User, error) {
	// Проверяем, существует ли пользователь уже
	exists, err := s.userRepo.ExistsByTelegramID(ctx, input.TelegramID)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to check user existence")
	}

	if exists {
		return nil, apperrors.ErrUserAlreadyExists
	}

	// 2. Установка дефолтов (Бизнес-логика)
	if input.FullName == nil || *input.FullName == "" {
		input.FullName = new(`Пользователь Telegram`)
	}
	// 3. Установка роли (тоже в сервисе!)
	role := userval.RoleUser
	input.Role = &role

	newUser, err := s.userRepo.Create(ctx, input)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to create user").WithDetails(map[string]any{
			"err": err.Error(),
		})
	}

	return newUser, nil
}

// DeleteUser удаляет пользователя по ID (soft delete)
func (s *UserServiceImpl) DeleteUser(ctx context.Context, userID string) error {
	// Проверяем, существует ли пользователь
	query := s.userRepo.GetQueryByID(ctx, userID)
	_, err := query.Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return apperrors.ErrUserNotFound
		}
		return apperrors.Internal(err, "failed to get user")
	}

	// Удаляем пользователя
	if err := s.userRepo.Delete(ctx, userID); err != nil {
		return apperrors.Internal(err, "failed to delete user")
	}

	return nil
}

func (s *UserServiceImpl) Update(ctx context.Context, id string, input *ent.UpdateUserInput) error {
	// 1. Если пришел LocationID, проверяем его прямо здесь (или в репо)
	if input.LocationID != nil {
		exists, err := s.locationRepo.Exists(ctx, *input.LocationID)
		if err != nil || !exists {
			return apperrors.ErrLocationNotFound
		}
	}

	// 2. Просто обновляем. SetInput сам проигнорирует nil поля.
	err := s.userRepo.Update(ctx, id, input)
	if err != nil {
		// Используем ent.IsNotFound - это идиоматический способ для Ent
		if ent.IsNotFound(err) {
			return apperrors.ErrUserNotFound
		}
		return apperrors.Internal(err, "failed to update user")
	}

	return nil
}

// GetUserByID получает пользователя по ID
func (s *UserServiceImpl) GetUserByID(ctx context.Context, userID string, opts UserOptions) (*ent.User, error) {
	query := s.userRepo.GetQueryByID(ctx, userID)

	if opts.WithPets {
		query = query.WithPets()
	}

	u, err := query.Only(ctx)
	if err != nil {
		return nil, err
	}

	// Преобразуем пути к фото в полные URL
	u.PhotoUrls = s.BuildFullPhotoURLs(u.PhotoUrls)

	return u, nil
}

// GetUserByTelegramID получает пользователя по Telegram ID
func (s *UserServiceImpl) GetUserByTelegramID(ctx context.Context, telegramID int64, opts UserOptions) (*ent.User, error) {
	query := s.userRepo.GetQueryByTelegram(ctx, telegramID)

	if opts.WithPets {
		query = query.WithPets()
	}

	u, err := query.First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, apperrors.ErrUserNotFound
		}
		return nil, apperrors.Internal(err, "failed to get user")
	}

	// Преобразуем пути к фото в полные URL
	u.PhotoUrls = s.BuildFullPhotoURLs(u.PhotoUrls)

	return u, nil
}

// ResetUser сбрасывает пользователя к начальным настройкам
func (s *UserServiceImpl) ResetUser(ctx context.Context, userID string) error {
	if err := s.userRepo.ResetUser(ctx, userID); err != nil {
		return apperrors.Internal(err, "failed to reset user")
	}
	return nil
}

// RestoreUser восстанавливает удаленного пользователя
func (s *UserServiceImpl) RestoreUser(ctx context.Context, userID string) error {
	if err := s.userRepo.RestoreUser(ctx, userID); err != nil {
		return apperrors.Internal(err, "failed to restore user")
	}
	return nil
}

// GetDeletedUsers получает всех удаленных пользователей
func (s *UserServiceImpl) GetDeletedUsers(ctx context.Context) ([]*ent.User, error) {
	users, err := s.userRepo.GetDeletedUsers(ctx)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to get deleted users")
	}
	return users, nil
}

// BuildFullPhotoURLs преобразует пути к фото в полные публичные URL
func (s *UserServiceImpl) BuildFullPhotoURLs(paths []string) []string {
	if len(paths) == 0 {
		return []string{}
	}
	result := make([]string, len(paths))
	for i, path := range paths {
		if path == "" {
			result[i] = ""
		} else {
			result[i] = s.storage.GetPublicURLFromPath(path)
		}
	}
	return result
}
