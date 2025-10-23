package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/models"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/repositories"
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
	FullName         *string `json:"full_name,omitempty" validate:"omitempty,min=2,max=255"`
	Phone            *string `json:"phone,omitempty" validate:"omitempty,e164"`
	Email            *string `json:"email,omitempty" validate:"omitempty,email"`
	OrganizationName *string `json:"organization_name,omitempty" validate:"omitempty,max=255"`
	LocationID       *int    `json:"location_id,omitempty" validate:"omitempty,min=1"`
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
		return nil, fmt.Errorf("проверка существования пользователя: %w", err)
	}

	if exists {
		return nil, errors.New("пользователь с этим telegram ID уже существует")
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
		Role:             userData.Role,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("создание пользователя: %w", err)
	}

	return user, nil
}

// RegisterUserSimple создает нового пользователя с Telegram ID и базовой информацией (для команды Start)
func (s *UserServiceImpl) RegisterUserSimple(ctx context.Context, telegramID int64, fullName string) (*models.User, error) {
	// Проверяем, существует ли пользователь уже
	exists, err := s.userRepo.ExistsByTelegramID(ctx, telegramID)
	if err != nil {
		return nil, fmt.Errorf("проверка существования пользователя: %w", err)
	}

	if exists {
		return nil, errors.New("пользователь с этим telegram ID уже существует")
	}

	// Создаем нового пользователя с Telegram ID, базовой информацией и значениями по умолчанию для обязательных полей
	user := &models.User{
		TelegramID:       telegramID,
		FullName:         fullName,
		Phone:            "",     // Пустая строка для телефона (будет заполнена позже)
		Email:            "",     // Пустая строка для email (будет заполнена позже)
		OrganizationName: "",     // Пустая строка для организации (будет заполнена позже)
		ConsentPD:        true,   // По умолчанию false, пользователь должен явно согласиться позже
		LocationID:       1,      // По умолчанию 0, пользователь должен установить местоположение позже
		Role:             "user", // Роль по умолчанию
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("создание пользователя: %w", err)
	}

	return user, nil
}

// DeleteUser удаляет пользователя по ID (soft delete)
func (s *UserServiceImpl) DeleteUser(ctx context.Context, userID int) error {
	// Проверяем, существует ли пользователь
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("получение пользователя для удаления: %w", err)
	}

	if user == nil {
		return errors.New("пользователь не найден")
	}

	// Удаляем пользователя
	if err := s.userRepo.Delete(ctx, userID); err != nil {
		return fmt.Errorf("удаление пользователя: %w", err)
	}

	return nil
}

// GetUserProfile получает полный профиль пользователя с питомцами и клиниками
func (s *UserServiceImpl) GetUserProfile(ctx context.Context, userID int) (*UserProfile, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("получение пользователя: %w", err)
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
		return fmt.Errorf("получение пользователя для обновления: %w", err)
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
	if updates.OrganizationName != nil {
		user.OrganizationName = *updates.OrganizationName
	}
	if updates.LocationID != nil {
		user.LocationID = *updates.LocationID
	}

	// Сохраняем обновленного пользователя
	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("обновление пользователя: %w", err)
	}

	return nil
}

// GetUserByTelegramID получает пользователя по Telegram ID
func (s *UserServiceImpl) GetUserByTelegramID(ctx context.Context, telegramID int64) (*models.User, error) {
	user, err := s.userRepo.GetByTelegramID(ctx, telegramID)
	if err != nil {
		return nil, fmt.Errorf("получение пользователя по telegram ID: %w", err)
	}
	return user, nil
}
