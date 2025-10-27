package services

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/models"
	repositories "github.com/artesipov-alt/odnoi-krovi-app/internal/repositories/interfaces"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/utils/validation"
)

// UserService определяет интерфейс для бизнес-логики пользователей
type UserService interface {
	// RegisterUser регистрирует нового пользователя в системе
	RegisterUser(ctx context.Context, telegramID int64, userData UserRegistration) (*models.User, error)

	// RegisterUserSimple создает нового пользователя с Telegram ID и базовой информацией (для команды Start)
	RegisterUserSimple(ctx context.Context, telegramID int64, fullName string) (*models.User, error)

	// GetUserProfile получает полный профиль пользователя с питомцами и клиниками
	GetUserProfile(ctx context.Context, userID int) (*UserProfile, error)

	// UpdateUserProfile обновляет информацию о пользователе
	UpdateUserProfile(ctx context.Context, userID int, updates UserUpdate) error

	// GetUserByTelegramID получает пользователя по Telegram ID
	GetUserByTelegramID(ctx context.Context, telegramID int64) (*models.User, error)

	// DeleteUser удаляет пользователя по ID (soft delete)
	DeleteUser(ctx context.Context, userID int) error
}

// UserRegistration содержит данные для регистрации пользователя
type UserRegistration struct {
	FullName         string          `json:"full_name" validate:"required,min=2,max=255"`
	Phone            string          `json:"phone" validate:"required,e164"`
	Email            string          `json:"email" validate:"omitempty,email"`
	OrganizationName string          `json:"organization_name" validate:"omitempty,max=255"`
	ConsentPD        bool            `json:"consent_pd" validate:"required"`
	LocationID       int             `json:"location_id" validate:"required,min=1"`
	Role             models.UserRole `json:"role" validate:"required,oneof=user clinic_admin"`
}

// UserUpdate содержит поля, которые можно обновить для пользователя
type UserUpdate struct {
	FullName *string `json:"full_name,omitempty" validate:"omitempty,min=2,max=255"`
	Phone    *string `json:"phone,omitempty" validate:"omitempty,e164"`
	Email    *string `json:"email,omitempty" validate:"omitempty,email"`
	// OrganizationName *string `json:"organization_name,omitempty" validate:"omitempty,max=255"`
	AllowGeo   *bool `json:"allow_geo,omitempty" validate:"omitempty"`
	OnBoarding *bool `json:"on_boarding,omitempty" validate:"omitempty"`
	LocationID *int  `json:"location_id,omitempty" validate:"omitempty,min=1"`
}

// UserProfile представляет полный профиль пользователя с связанными данными
type UserProfile struct {
	User   *models.User      `json:"user"`
	Pets   []*models.Pet     `json:"pets,omitempty"`
	Clinic *models.VetClinic `json:"clinic,omitempty"`
}

// UserServiceImpl реализует UserService
type UserServiceImpl struct {
	userRepo repositories.UserRepository
	// petRepo и clinicRepo будут добавлены здесь для полного профиля
}

// NewUserService создает новый сервис пользователей
func NewUserService(userRepo repositories.UserRepository) *UserServiceImpl {
	return &UserServiceImpl{
		userRepo: userRepo,
	}
}

// RegisterUser регистрирует нового пользователя в системе
func (s *UserServiceImpl) RegisterUser(ctx context.Context, telegramID int64, userData UserRegistration) (*models.User, error) {
	// Проверяем, существует ли пользователь уже
	exists, err := s.userRepo.ExistsByTelegramID(ctx, telegramID)
	if err != nil {
		return nil, apperrors.Internal(err, "не удалось проверить существование пользователя")
	}

	if exists {
		return nil, apperrors.NewUserAlreadyExistsError(telegramID)
	}

	// Валидируем роль пользователя
	role, err := validation.ValidateUserRole(string(userData.Role))
	if err != nil {
		return nil, apperrors.BadRequest("неверная роль пользователя")
	}

	// Создаем нового пользователя
	user := &models.User{
		TelegramID:       telegramID,
		FullName:         userData.FullName,
		Phone:            userData.Phone,
		Email:            userData.Email,
		OrganizationName: userData.OrganizationName,
		ConsentPD:        userData.ConsentPD,
		LocationID:       userData.LocationID,
		Role:             role,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, apperrors.Internal(err, "не удалось создать пользователя")
	}

	return user, nil
}

// RegisterUserSimple создает нового пользователя с Telegram ID и базовой информацией (для команды Start)
func (s *UserServiceImpl) RegisterUserSimple(ctx context.Context, telegramID int64, fullName string) (*models.User, error) {
	// Проверяем, существует ли пользователь уже
	exists, err := s.userRepo.ExistsByTelegramID(ctx, telegramID)
	if err != nil {
		return nil, apperrors.Internal(err, "не удалось проверить существование пользователя")
	}

	if exists {
		return nil, apperrors.NewUserAlreadyExistsError(telegramID)
	}

	// Создаем нового пользователя с Telegram ID, базовой информацией и значениями по умолчанию для обязательных полей
	user := &models.User{
		TelegramID:       telegramID,
		FullName:         fullName,
		Phone:            "",     // Пустая строка для телефона (будет заполнена позже)
		Email:            "",     // Пустая строка для email (будет заполнена позже)
		OrganizationName: "",     // Пустая строка для организации (будет заполнена позже)
		ConsentPD:        true,   // По умолчанию false, пользователь должен явно согласиться позже
		LocationID:       1,      // По умолчанию 1, пользователь должен установить местоположение позже
		OnBoarding:       false,  // По умолчанию false, пользователь должен пройти опрос
		AllowGeo:         false,  // По умолчанию false, пользователь должен дорегистрацию
		Role:             "user", // Роль по умолчанию
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, apperrors.Internal(err, "не удалось создать пользователя")
	}

	return user, nil
}

// DeleteUser удаляет пользователя по ID (soft delete)
func (s *UserServiceImpl) DeleteUser(ctx context.Context, userID int) error {
	// Проверяем, существует ли пользователь
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return apperrors.Internal(err, "не удалось получить пользователя")
	}

	if user == nil {
		return apperrors.NewUserNotFoundError(userID)
	}

	// Удаляем пользователя
	if err := s.userRepo.Delete(ctx, userID); err != nil {
		return apperrors.Internal(err, "не удалось удалить пользователя")
	}

	return nil
}

// GetUserProfile получает полный профиль пользователя с питомцами и клиниками
func (s *UserServiceImpl) GetUserProfile(ctx context.Context, userID int) (*UserProfile, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, apperrors.Internal(err, "не удалось получить пользователя")
	}

	if user == nil {
		return nil, apperrors.NewUserNotFoundError(userID)
	}

	profile := &UserProfile{
		User: user,
		// TODO: Добавить данные о питомцах и клинике, когда репозитории будут доступны
		Pets:   []*models.Pet{},
		Clinic: nil,
	}

	return profile, nil
}

// UpdateUserProfile обновляет информацию о пользователе
func (s *UserServiceImpl) UpdateUserProfile(ctx context.Context, userID int, updates UserUpdate) error {
	// Получаем существующего пользователя
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return apperrors.Internal(err, "не удалось получить пользователя")
	}

	if user == nil {
		return apperrors.NewUserNotFoundError(userID)
	}

	// Применяем обновления
	if updates.FullName != nil {
		user.FullName = *updates.FullName
	}
	if updates.Phone != nil {
		user.Phone = *updates.Phone
	}
	if updates.Email != nil {
		user.Email = *updates.Email
	}
	// if updates.OrganizationName != nil {
	// 	user.OrganizationName = *updates.OrganizationName
	// }
	if updates.AllowGeo != nil {
		user.AllowGeo = *updates.AllowGeo
	}
	if updates.OnBoarding != nil {
		user.OnBoarding = *updates.OnBoarding
	}
	if updates.LocationID != nil {
		user.LocationID = *updates.LocationID
	}

	// Сохраняем обновленного пользователя
	if err := s.userRepo.Update(ctx, user); err != nil {
		return apperrors.Internal(err, "не удалось обновить пользователя")
	}

	return nil
}

// GetUserByTelegramID получает пользователя по Telegram ID
func (s *UserServiceImpl) GetUserByTelegramID(ctx context.Context, telegramID int64) (*models.User, error) {
	user, err := s.userRepo.GetByTelegramID(ctx, telegramID)
	if err != nil {
		return nil, apperrors.Internal(err, "не удалось получить пользователя по Telegram ID")
	}

	if user == nil {
		return nil, apperrors.NotFound("пользователь с таким Telegram ID не найден").WithDetails(map[string]interface{}{
			"telegram_id": telegramID,
		})
	}

	return user, nil
}
