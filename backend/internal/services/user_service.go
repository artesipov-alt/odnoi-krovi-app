package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/models"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/repositories"
)

// UserService defines the interface for user business logic
type UserService interface {
	// RegisterUser registers a new user in the system
	RegisterUser(ctx context.Context, telegramID int64, userData UserRegistration) (*models.User, error)

	// GetUserProfile retrieves complete user profile with pets and clinics
	GetUserProfile(ctx context.Context, userID int) (*UserProfile, error)

	// UpdateUserProfile updates user information
	UpdateUserProfile(ctx context.Context, userID int, updates UserUpdate) error

	// GetUserByTelegramID retrieves user by Telegram ID
	GetUserByTelegramID(ctx context.Context, telegramID int64) (*models.User, error)
}

// UserRegistration contains data for user registration
type UserRegistration struct {
	FullName         string `json:"full_name" validate:"required,min=2,max=255"`
	Phone            string `json:"phone" validate:"required,e164"`
	Email            string `json:"email" validate:"omitempty,email"`
	OrganizationName string `json:"organization_name" validate:"omitempty,max=255"`
	ConsentPD        bool   `json:"consent_pd" validate:"required"`
	LocationID       int    `json:"location_id" validate:"required,min=1"`
	Role             string `json:"role" validate:"required,oneof=user clinic_admin"`
}

// UserUpdate contains fields that can be updated for a user
type UserUpdate struct {
	FullName         *string `json:"full_name,omitempty" validate:"omitempty,min=2,max=255"`
	Phone            *string `json:"phone,omitempty" validate:"omitempty,e164"`
	Email            *string `json:"email,omitempty" validate:"omitempty,email"`
	OrganizationName *string `json:"organization_name,omitempty" validate:"omitempty,max=255"`
	LocationID       *int    `json:"location_id,omitempty" validate:"omitempty,min=1"`
}

// UserProfile represents a complete user profile with related data
type UserProfile struct {
	User   *models.User      `json:"user"`
	Pets   []*models.Pet     `json:"pets,omitempty"`
	Clinic *models.VetClinic `json:"clinic,omitempty"`
}

// UserServiceImpl implements UserService
type UserServiceImpl struct {
	userRepo repositories.UserRepository
	// petRepo and clinicRepo would be added here for complete profile
}

// NewUserService creates a new user service
func NewUserService(userRepo repositories.UserRepository) *UserServiceImpl {
	return &UserServiceImpl{
		userRepo: userRepo,
	}
}

// RegisterUser registers a new user in the system
func (s *UserServiceImpl) RegisterUser(ctx context.Context, telegramID int64, userData UserRegistration) (*models.User, error) {
	// Check if user already exists
	exists, err := s.userRepo.ExistsByTelegramID(ctx, telegramID)
	if err != nil {
		return nil, fmt.Errorf("check user existence: %w", err)
	}

	if exists {
		return nil, errors.New("user with this telegram ID already exists")
	}

	// Create new user
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
		return nil, fmt.Errorf("create user: %w", err)
	}

	return user, nil
}

// GetUserProfile retrieves complete user profile with pets and clinics
func (s *UserServiceImpl) GetUserProfile(ctx context.Context, userID int) (*UserProfile, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}

	profile := &UserProfile{
		User: user,
		// TODO: Add pets and clinic data when repositories are available
		Pets:   []*models.Pet{},
		Clinic: nil,
	}

	return profile, nil
}

// UpdateUserProfile updates user information
func (s *UserServiceImpl) UpdateUserProfile(ctx context.Context, userID int, updates UserUpdate) error {
	// Get existing user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("get user for update: %w", err)
	}

	// Apply updates
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

	// Save updated user
	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("update user: %w", err)
	}

	return nil
}

// GetUserByTelegramID retrieves user by Telegram ID
func (s *UserServiceImpl) GetUserByTelegramID(ctx context.Context, telegramID int64) (*models.User, error) {
	user, err := s.userRepo.GetByTelegramID(ctx, telegramID)
	if err != nil {
		return nil, fmt.Errorf("get user by telegram ID: %w", err)
	}
	return user, nil
}
