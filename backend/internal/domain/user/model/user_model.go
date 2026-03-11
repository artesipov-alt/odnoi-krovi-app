package model

import (
	"errors"
	"time"

	pet "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet/model"
)

// UserRole represents user role types
type UserRole string

const (
	RoleUser   UserRole = "user"
	RoleAdmin  UserRole = "admin"
	RoleClinic UserRole = "clinic"
)

// User представляет доменную модель пользователя
type User struct {
	ID               string
	TelegramID       int64
	FullName         string
	Phone            string
	Email            string
	PhotoURLs        []string
	OrganizationName string
	ConsentPd        bool
	OnBoarding       []string
	AllowGeo         bool
	LocationID       *string
	Role             UserRole
	OriginSource     string
	Pets             []*pet.Pet
	DonorPreference  *DonorPreference
	CreatedAt        *time.Time
	UpdatedAt        *time.Time
	DeletedAt        *time.Time
}

// NewUserParams holds the parameters for creating a new User
type NewUserParams struct {
	FullName   string
	Phone      string
	Email      string
	Role       UserRole
	ConsentPd  bool
	LocationID *string
	MetaData   map[string]any
}

// CompensationType represents donor's compensation preference
type CompensationType string

const (
	CompensationFree CompensationType = "free" // Готов помочь безвозмездно
	CompensationPaid CompensationType = "paid" // Не готов помочь бесплатно
	CompensationFood CompensationType = "food" // Готов помочь за корм
)

// NotificationFrequency represents how often donor wants to be notified
type NotificationFrequency string

const (
	NotifyImmediately NotificationFrequency = "immediately" // Сразу
	NotifyDaily       NotificationFrequency = "daily"       // Раз в день
	NotifyWeekly      NotificationFrequency = "weekly"      // Раз в неделю
	NotifyNever       NotificationFrequency = "never"       // Никогда
)

// DonorPreference represents donor's default preferences for blood donation responses
type DonorPreference struct {
	ID                    string
	UserID                string
	PreferredLocationIDs  []string
	RecoveryPeriodMonths  int
	CompensationType      CompensationType
	TaxiCompensation      bool
	NotificationFrequency NotificationFrequency
	CreatedAt             *time.Time
	UpdatedAt             *time.Time
	DeletedAt             *time.Time
}

// DonorPreferenceParams holds the parameters for creating or updating a DonorPreference.
type DonorPreferenceParams struct {
	PreferredLocationIDs  []string
	RecoveryPeriodMonths  int
	CompensationType      CompensationType
	TaxiCompensation      bool
	NotificationFrequency NotificationFrequency
}

// NewUser creates a new User aggregate with validation
func NewUser(userparams NewUserParams) (*User, error) {
	// Validation
	if userparams.FullName == "" {
		userparams.FullName = "Пользователь портала"
	}
	if len(userparams.FullName) > 100 {
		return nil, errors.New("full name must be less than 100 characters")
	}
	if userparams.Role == "" {
		userparams.Role = RoleUser // default role
	}
	if userparams.Role != RoleUser && userparams.Role != RoleAdmin {
		return nil, errors.New("invalid user role")
	}
	if userparams.Email != "" {
		// Basic email validation could be added here
		if len(userparams.Email) > 100 {
			return nil, errors.New("email must be less than 100 characters")
		}
	}
	if userparams.Phone != "" && len(userparams.Phone) > 20 {
		return nil, errors.New("phone must be less than 20 characters")
	}

	// Extract OriginSource from metadata (utm_campaign)
	originSource := ""
	if userparams.MetaData != nil {
		if val, ok := userparams.MetaData["utm_campaign"]; ok {
			originSource, _ = val.(string)
		}
	}

	user := &User{
		FullName:     userparams.FullName,
		Phone:        userparams.Phone,
		Email:        userparams.Email,
		Role:         userparams.Role,
		ConsentPd:    userparams.ConsentPd,
		LocationID:   userparams.LocationID,
		OriginSource: originSource,
		PhotoURLs:    []string{},
		OnBoarding:   []string{},
		Pets:         []*pet.Pet{},
	}

	return user, nil
}

func NewDefaultUser(fullName string, originSource string) (*User, error) {

	if fullName == "" {
		fullName = "Пользователь портала"
	}
	if originSource == "" {
		originSource = "self"
	}
	return &User{
		FullName:     fullName,
		OriginSource: originSource,
	}, nil
}

func (u *User) SetRole(role string) {
	u.Role = UserRole(role)
}

// NewDonorPreferenceParams creates a new DonorPreferenceParams with default values
func DefaultDonorPreference() *DonorPreference {
	return &DonorPreference{
		PreferredLocationIDs:  []string{},
		RecoveryPeriodMonths:  2,
		CompensationType:      "",
		TaxiCompensation:      false,
		NotificationFrequency: NotifyImmediately,
	}
}
