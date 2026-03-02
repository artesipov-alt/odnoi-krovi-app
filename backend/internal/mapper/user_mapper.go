// Package mapper provides conversion functions between domain models and DTOs.
package mapper

import (
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/transport/http/dto"
)

// UserMapper handles conversions between domain User model and DTOs.
type UserMapper struct {
	petMapper *PetMapper
}

// NewUserMapper creates a new UserMapper instance.
func NewUserMapper(petMapper *PetMapper) *UserMapper {
	return &UserMapper{
		petMapper: petMapper,
	}
}

// ToResponse converts a domain User model to a UserDetail DTO.
func (m *UserMapper) ToResponse(u *model.User) dto.UserDetail {
	if u == nil {
		return dto.UserDetail{}
	}

	userDTO := dto.UserDetail{
		ID:               u.ID,
		TelegramID:       u.TelegramID,
		FullName:         u.FullName,
		Phone:            u.Phone,
		Email:            u.Email,
		PhotoURLs:        u.PhotoURLs,
		OrganizationName: u.OrganizationName,
		ConsentPd:        u.ConsentPd,
		OnBoarding:       u.OnBoarding,
		AllowGeo:         u.AllowGeo,
		LocationID:       "",
		Role:             u.Role,
		CreatedAt:        u.CreatedAt,
		UpdatedAt:        u.UpdatedAt,
		DeletedAt:        u.DeletedAt,
	}

	if u.LocationID != nil {
		userDTO.LocationID = *u.LocationID
	}

	return userDTO
}

// ToResponseSlice converts a slice of domain User models to UserDetail DTOs.
func (m *UserMapper) ToResponseSlice(users []*model.User) []dto.UserDetail {
	if users == nil {
		return nil
	}
	dtos := make([]dto.UserDetail, len(users))
	for i, u := range users {
		dtos[i] = m.ToResponse(u)
	}
	return dtos
}

// FromCreate converts a CreateUserBody DTO to a domain User model using the constructor.
func (m *UserMapper) FromCreate(body dto.CreateUserBody) (*model.User, error) {
	params := model.NewUserParams{
		TelegramID: body.TelegramID,
		FullName:   body.FullName,
		Phone:      "", // phone - empty for simple creation
		Email:      "", // email - empty for simple creation
		Role:       model.RoleUser,
		ConsentPd:  false, // consentPd
		LocationID: nil,   // locationID
	}

	prefs := &model.DonorPreferenceParams{
		RecoveryPeriodMonths:  2,
		NotificationFrequency: model.NotifyImmediately,
	}
	return model.NewUser(params, prefs)
}
