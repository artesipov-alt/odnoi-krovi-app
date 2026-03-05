package model

import (
	"errors"
	"time"

	pet "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet/model"
)

// UserRole represents user role types
type UserRole string

const (
	RoleUser  UserRole = "user"
	RoleAdmin UserRole = "admin"
)

// User представляет доменную модель пользователя
type User struct {
	ID               string
	TelegramID       int64
	ProviderID       int64
	ProviderName     string
	FullName         string
	Phone            string
	Email            string
	PhotoURLs        []string
	OrganizationName string
	ConsentPd        bool
	OnBoarding       []string
	AllowGeo         bool
	LocationID       *string
	Role             string
	OriginSource     string
	Pets             []*pet.Pet
	DonorPreference  *DonorPreference
	CreatedAt        *time.Time
	UpdatedAt        *time.Time
	DeletedAt        *time.Time
}

// NewUserParams holds the parameters for creating a new User
type NewUserParams struct {
	ProviderID   int64
	ProviderName string
	FullName     string
	Phone        string
	Email        string
	Role         UserRole
	ConsentPd    bool
	LocationID   *string
	MetaData     map[string]string
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
func NewUser(userparams NewUserParams, donorparams *DonorPreferenceParams) (*User, error) {
	// Validation
	if userparams.FullName == "" {
		return nil, errors.New("full name is required")
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
		originSource = userparams.MetaData["utm_campaign"]
	}

	user := &User{
		ProviderID:   userparams.ProviderID,
		ProviderName: userparams.ProviderName,
		FullName:     userparams.FullName,
		Phone:        userparams.Phone,
		Email:        userparams.Email,
		Role:         string(userparams.Role),
		ConsentPd:    userparams.ConsentPd,
		LocationID:   userparams.LocationID,
		OriginSource: originSource,
		PhotoURLs:    []string{},
		OnBoarding:   []string{},
		Pets:         []*pet.Pet{},
	}

	return user, nil
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

// Metadata represents user metadata with UTM and other fields
type UTM struct {
	Source   string
	Medium   string
	Campaign string
	Content  string
	Term     string
}

type Metadata struct {
	UTMData *UTM
}

func NewUserMetadata(metadata map[string]string) *Metadata {
	source, medium, campaign, content, term := extractUTMFromMetadata(metadata)
	return &Metadata{
		UTMData: &UTM{
			Source:   source,
			Medium:   medium,
			Campaign: campaign,
			Content:  content,
			Term:     term,
		},
	}
}

func extractUTMFromMetadata(metadata map[string]string) (string, string, string, string, string) {
	return metadata["utm_source"], metadata["utm_medium"], metadata["utm_campaign"], metadata["utm_content"], metadata["utm_term"]
}
