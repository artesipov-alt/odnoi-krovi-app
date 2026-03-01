// Package mapper provides conversion functions between domain models and DTOs.
package mapper

import (
	petmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet/model"
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

// ToResponse converts a domain User model to a DTO.
func (m *UserMapper) ToResponse(u *model.User, pets []*petmodel.Pet) dto.User {
	if u == nil {
		return dto.User{}
	}

	userDTO := dto.User{
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

	if pets != nil && m.petMapper != nil {
		userDTO.Pets = m.petMapper.ToSimplifiedResponseSlice(pets)
	}

	return userDTO
}

// ToResponseSlice converts a slice of domain User models to DTOs.
func (m *UserMapper) ToResponseSlice(users []*model.User) []dto.User {
	if users == nil {
		return nil
	}
	dtos := make([]dto.User, len(users))
	for i, u := range users {
		dtos[i] = m.ToResponse(u, nil)
	}
	return dtos
}
