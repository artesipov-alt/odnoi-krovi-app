package model

import (
	"errors"
	"time"

	authmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/auth/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/common"
	petmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet/model"
)

// UserRole представляет типы ролей пользователя
type UserRole string

const (
	RoleUser   UserRole = "user"
	RoleAdmin  UserRole = "admin"
	RoleClinic UserRole = "clinic"
)

// User представляет доменную модель пользователя
type User struct {
	ID                  string
	TelegramID          int64
	FullName            string
	Phone               string
	Verified            bool
	Email               string
	PhotoURLs           []string
	OrganizationName    string
	ConsentPd           bool
	OnBoarding          []string
	AllowGeo            bool
	LocationID          *string
	Role                UserRole
	OriginSource        string
	PrioritySearchCount int
	Pets                []*petmodel.Pet
	DonorPreference     *DonorPreference
	Identities          []*authmodel.Identity
	CreatedAt           *time.Time
	UpdatedAt           *time.Time
	DeletedAt           *time.Time
	LastSeenAt          *time.Time
}

// NewUserParams содержит параметры для создания нового пользователя
type NewUserParams struct {
	FullName   string
	Phone      string
	Email      string
	Role       UserRole
	ConsentPd  bool
	LocationID *string
	MetaData   map[string]any
}

// NotificationFrequency определяет, как часто донор хочет получать уведомления
type NotificationFrequency string

const (
	NotifyImmediately NotificationFrequency = "immediately" // Сразу
	NotifyDaily       NotificationFrequency = "daily"       // Раз в день
	NotifyWeekly      NotificationFrequency = "weekly"      // Раз в неделю
	NotifyNever       NotificationFrequency = "never"       // Никогда
)

// DonorPreference представляет предпочтения донора по умолчанию для откликов на донации
type DonorPreference struct {
	ID                    string
	UserID                string
	PreferredLocationIDs  []string
	RecoveryPeriodMonths  int
	CompensationType      common.CompensationType
	TaxiCompensation      bool
	NotificationFrequency NotificationFrequency
	CreatedAt             *time.Time
	UpdatedAt             *time.Time
	DeletedAt             *time.Time
}

// DonorPreferenceParams содержит параметры для создания или обновления DonorPreference
type DonorPreferenceParams struct {
	PreferredLocationIDs  []string
	RecoveryPeriodMonths  int
	CompensationType      common.CompensationType
	TaxiCompensation      bool
	NotificationFrequency NotificationFrequency
}

// NewUser создаёт новый агрегат User с валидацией
func NewUser(userparams NewUserParams) (*User, error) {
	// Валидация
	if userparams.FullName == "" {
		userparams.FullName = "Пользователь портала"
	}
	if len(userparams.FullName) > 100 {
		return nil, errors.New("full name must be less than 100 characters")
	}
	if userparams.Role == "" {
		userparams.Role = RoleUser // роль по умолчанию
	}
	if userparams.Role != RoleUser && userparams.Role != RoleAdmin {
		return nil, errors.New("invalid user role")
	}
	if userparams.Email != "" {
		// Базовая валидация email
		if len(userparams.Email) > 100 {
			return nil, errors.New("email must be less than 100 characters")
		}
	}
	if userparams.Phone != "" && len(userparams.Phone) > 20 {
		return nil, errors.New("phone must be less than 20 characters")
	}

	// Извлечение OriginSource из метаданных (utm_campaign)
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
		Pets:         []*petmodel.Pet{},
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

// MessengerContacts возвращает ID провайдеров для Max и Telegram.
func (u *User) MessengerContacts() (maxID, telegramID string) {
	if u.Identities == nil {
		return "", ""
	}

	for _, identity := range u.Identities {
		switch identity.ProviderName {
		case authmodel.ProviderMax:
			maxID = identity.ProviderUserID
		case authmodel.ProviderTelegram:
			telegramID = identity.ProviderUserID
		}
		if maxID != "" && telegramID != "" {
			return maxID, telegramID
		}
	}
	return maxID, telegramID
}

// HasProviderConflict проверяет, есть ли у двух пользователей пересечение по провайдерам.
// Если хоть один провайдер совпадает — слияние аккаунтов запрещено.
func (u *User) HasProviderConflict(other User) bool {
	if len(other.Identities) == 0 {
		return false
	}
	providers := make(map[authmodel.ProviderName]struct{})
	for _, identity := range u.Identities {
		providers[identity.ProviderName] = struct{}{}
	}
	for _, identity := range other.Identities {
		if _, ok := providers[identity.ProviderName]; ok {
			return true
		}
	}
	return false
}

// DefaultDonorPreference создаёт DonorPreference со значениями по умолчанию
func DefaultDonorPreference() *DonorPreference {
	return &DonorPreference{
		PreferredLocationIDs:  []string{},
		RecoveryPeriodMonths:  2,
		CompensationType:      "",
		TaxiCompensation:      false,
		NotificationFrequency: NotifyImmediately,
	}
}
