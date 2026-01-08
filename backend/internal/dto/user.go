package dto

// SimpleRegistrationRequest represents a simple registration request for users
type SimpleRegistrationRequest struct {
	TelegramID int64  `json:"telegramId" validate:"required,min=1" example:"123456789"`
	FullName   string `json:"fullName,omitempty" validate:"omitempty,min=1,max=255" example:"Иван Иванов"`
}

// UserRegistration represents the full user registration structure
type UserRegistration struct {
	FullName   string `json:"fullName" validate:"required,min=2,max=255"`
	Phone      string `json:"phone" validate:"required,e164"`
	Email      string `json:"email" validate:"omitempty,email"`
	ConsentPD  bool   `json:"consentPd" validate:"required"`
	LocationID int    `json:"locationId" validate:"required,min=1"`
	Role       string `json:"role" validate:"required,oneof=user admin"`
}

// UserUpdate represents the structure for updating user information
type UserUpdate struct {
	FullName   *string `json:"fullName,omitempty" validate:"omitempty,min=2,max=255"`
	Phone      *string `json:"phone,omitempty" validate:"omitempty,e164"`
	Email      *string `json:"email,omitempty" validate:"omitempty,email"`
	AllowGeo   *bool   `json:"allowGeo,omitempty" validate:"omitempty"`
	OnBoarding *bool   `json:"onBoarding,omitempty" validate:"omitempty"`
	LocationID *int    `json:"locationId,omitempty" validate:"omitempty,min=1"`
}
