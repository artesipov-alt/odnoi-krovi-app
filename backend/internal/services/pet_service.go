package services

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/models"
	repositories "github.com/artesipov-alt/odnoi-krovi-app/internal/repositories/interfaces"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/utils/validation"
)

// PetService определяет интерфейс для бизнес-логики питомцев
type PetService interface {
	// CreatePet создает нового питомца для пользователя
	CreatePet(ctx context.Context, userID int, petData PetCreate) (*models.Pet, error)

	// GetPetByID получает питомца по ID
	GetPetByID(ctx context.Context, petID int) (*models.Pet, error)

	// GetUserPets получает всех питомцев пользователя
	GetUserPets(ctx context.Context, userID int) ([]*models.Pet, error)

	// UpdatePet обновляет информацию о питомце
	UpdatePet(ctx context.Context, petID int, updates PetUpdate) error

	// DeletePet удаляет питомца по ID
	DeletePet(ctx context.Context, petID int) error
}

// PetCreate содержит данные для создания питомца
type PetCreate struct {
	Name            string                 `json:"name" validate:"required,min=1,max=100"`
	HasChip         bool                   `json:"hasChip"`
	ChipNumber      string                 `json:"chipNumber,omitempty" validate:"omitempty,max=50"`
	PhotoURL        string                 `json:"photoUrl,omitempty" validate:"omitempty,url,max=255"`
	KnowsBloodGroup bool                   `json:"knowsBloodGroup"`
	IsGuideDog      bool                   `json:"isGuideDog"`
	IsTherapist     bool                   `json:"isTherapist"`
	Breed           string                 `json:"breed,omitempty" validate:"omitempty,max=100"`
	WeightKg        float64                `json:"weightKg,omitempty" validate:"omitempty,min=0"`
	AgeYears        int                    `json:"ageYears,omitempty" validate:"omitempty,min=0"`
	AgeMonths       int                    `json:"ageMonths,omitempty" validate:"omitempty,min=0,max=11"`
	Sterilized      bool                   `json:"sterilized"`
	Latitude        float64                `json:"latitude,omitempty" validate:"omitempty,latitude"`
	Longitude       float64                `json:"longitude,omitempty" validate:"omitempty,longitude"`
	LivingCondition models.LivingCondition `json:"livingCondition,omitempty"`
	Gender          models.Gender          `json:"gender,omitempty"`
	Type            models.PetType         `json:"type,omitempty"`
	BloodGroup      string                 `json:"bloodGroup,omitempty" validate:"omitempty,max=50"`
}

// PetUpdate содержит поля, которые можно обновить для питомца
type PetUpdate struct {
	Name            *string                 `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	HasChip         *bool                   `json:"hasChip,omitempty"`
	ChipNumber      *string                 `json:"chipNumber,omitempty" validate:"omitempty,max=50"`
	PhotoURL        *string                 `json:"photoUrl,omitempty" validate:"omitempty,url,max=255"`
	KnowsBloodGroup *bool                   `json:"knowsBloodGroup,omitempty"`
	IsGuideDog      *bool                   `json:"isGuideDog,omitempty"`
	IsTherapist     *bool                   `json:"isTherapist,omitempty"`
	Breed           *string                 `json:"breed,omitempty" validate:"omitempty,max=100"`
	WeightKg        *float64                `json:"weightKg,omitempty" validate:"omitempty,min=0"`
	AgeYears        *int                    `json:"ageYears,omitempty" validate:"omitempty,min=0"`
	AgeMonths       *int                    `json:"ageMonths,omitempty" validate:"omitempty,min=0,max=11"`
	Sterilized      *bool                   `json:"sterilized,omitempty"`
	Latitude        *float64                `json:"latitude,omitempty" validate:"omitempty,latitude"`
	Longitude       *float64                `json:"longitude,omitempty" validate:"omitempty,longitude"`
	LivingCondition *models.LivingCondition `json:"livingCondition,omitempty"`
	Gender          *models.Gender          `json:"gender,omitempty"`
	Type            *models.PetType         `json:"type,omitempty"`
	BloodGroup      *string                 `json:"bloodGroup,omitempty" validate:"omitempty,max=50"`
}

// PetServiceImpl реализует PetService
type PetServiceImpl struct {
	petRepo  repositories.PetRepository
	userRepo repositories.UserRepository
}

// NewPetService создает новый сервис питомцев
func NewPetService(petRepo repositories.PetRepository, userRepo repositories.UserRepository) *PetServiceImpl {
	return &PetServiceImpl{
		petRepo:  petRepo,
		userRepo: userRepo,
	}
}

// CreatePet создает нового питомца для пользователя
func (s *PetServiceImpl) CreatePet(ctx context.Context, userID int, petData PetCreate) (*models.Pet, error) {
	// Проверяем, существует ли пользователь
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, apperrors.Internal(err, "не удалось проверить существование пользователя")
	}

	if user == nil {
		return nil, apperrors.NewUserNotFoundError(userID)
	}

	// Валидируем тип животного
	validatedType, err := validation.ValidatePetType(string(petData.Type))
	if err != nil {
		return nil, apperrors.ErrInvalidPetType
	}

	// Валидируем пол животного
	validatedGender, err := validation.ValidateGender(string(petData.Gender))
	if err != nil {
		return nil, apperrors.ErrInvalidGender
	}

	// Валидируем условия проживания
	validatedLivingCondition, err := validation.ValidateLivingCondition(string(petData.LivingCondition))
	if err != nil {
		return nil, apperrors.ErrInvalidLivingCondition
	}

	// Создаем нового питомца
	pet := &models.Pet{
		OwnerID:         userID,
		Name:            petData.Name,
		HasChip:         petData.HasChip,
		ChipNumber:      petData.ChipNumber,
		PhotoURL:        petData.PhotoURL,
		KnowsBloodGroup: petData.KnowsBloodGroup,
		IsGuideDog:      petData.IsGuideDog,
		IsTherapist:     petData.IsTherapist,
		Breed:           petData.Breed,
		WeightKg:        petData.WeightKg,
		AgeYears:        petData.AgeYears,
		AgeMonths:       petData.AgeMonths,
		Sterilized:      petData.Sterilized,
		Latitude:        petData.Latitude,
		Longitude:       petData.Longitude,
		LivingCondition: validatedLivingCondition,
		Gender:          validatedGender,
		Type:            validatedType,
		BloodGroup:      petData.BloodGroup,
	}

	if err := s.petRepo.Create(ctx, pet); err != nil {
		return nil, apperrors.Internal(err, "не удалось создать питомца")
	}

	return pet, nil
}

// GetPetByID получает питомца по ID
func (s *PetServiceImpl) GetPetByID(ctx context.Context, petID int) (*models.Pet, error) {
	pet, err := s.petRepo.GetByID(ctx, petID)
	if err != nil {
		return nil, apperrors.Internal(err, "не удалось получить питомца")
	}

	if pet == nil {
		return nil, apperrors.NewPetNotFoundError(petID)
	}

	return pet, nil
}

// GetUserPets получает всех питомцев пользователя
func (s *PetServiceImpl) GetUserPets(ctx context.Context, userID int) ([]*models.Pet, error) {
	// Проверяем, существует ли пользователь
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, apperrors.Internal(err, "не удалось проверить существование пользователя")
	}

	if user == nil {
		return nil, apperrors.NewUserNotFoundError(userID)
	}

	pets, err := s.petRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, apperrors.Internal(err, "не удалось получить питомцев пользователя")
	}

	return pets, nil
}

// UpdatePet обновляет информацию о питомце
func (s *PetServiceImpl) UpdatePet(ctx context.Context, petID int, updates PetUpdate) error {
	// Получаем существующего питомца
	pet, err := s.petRepo.GetByID(ctx, petID)
	if err != nil {
		return apperrors.Internal(err, "не удалось получить питомца")
	}

	if pet == nil {
		return apperrors.NewPetNotFoundError(petID)
	}

	// Применяем обновления
	if updates.Name != nil {
		pet.Name = *updates.Name
	}
	if updates.HasChip != nil {
		pet.HasChip = *updates.HasChip
	}
	if updates.ChipNumber != nil {
		pet.ChipNumber = *updates.ChipNumber
	}
	if updates.PhotoURL != nil {
		pet.PhotoURL = *updates.PhotoURL
	}
	if updates.KnowsBloodGroup != nil {
		pet.KnowsBloodGroup = *updates.KnowsBloodGroup
	}
	if updates.IsGuideDog != nil {
		pet.IsGuideDog = *updates.IsGuideDog
	}
	if updates.IsTherapist != nil {
		pet.IsTherapist = *updates.IsTherapist
	}
	if updates.Breed != nil {
		pet.Breed = *updates.Breed
	}
	if updates.WeightKg != nil {
		pet.WeightKg = *updates.WeightKg
	}
	if updates.AgeYears != nil {
		pet.AgeYears = *updates.AgeYears
	}
	if updates.AgeMonths != nil {
		pet.AgeMonths = *updates.AgeMonths
	}
	if updates.Sterilized != nil {
		pet.Sterilized = *updates.Sterilized
	}
	if updates.Latitude != nil {
		pet.Latitude = *updates.Latitude
	}
	if updates.Longitude != nil {
		pet.Longitude = *updates.Longitude
	}
	if updates.LivingCondition != nil {
		validatedLivingCondition, err := validation.ValidateLivingCondition(string(*updates.LivingCondition))
		if err != nil {
			return apperrors.ErrInvalidLivingCondition
		}
		pet.LivingCondition = validatedLivingCondition
	}
	if updates.Gender != nil {
		validatedGender, err := validation.ValidateGender(string(*updates.Gender))
		if err != nil {
			return apperrors.ErrInvalidGender
		}
		pet.Gender = validatedGender
	}
	if updates.Type != nil {
		validatedType, err := validation.ValidatePetType(string(*updates.Type))
		if err != nil {
			return apperrors.ErrInvalidPetType
		}
		pet.Type = validatedType
	}
	if updates.BloodGroup != nil {
		pet.BloodGroup = *updates.BloodGroup
	}

	// Сохраняем обновленного питомца
	if err := s.petRepo.Update(ctx, pet); err != nil {
		return apperrors.Internal(err, "не удалось обновить питомца")
	}

	return nil
}

// DeletePet удаляет питомца по ID
func (s *PetServiceImpl) DeletePet(ctx context.Context, petID int) error {
	// Проверяем, существует ли питомец
	pet, err := s.petRepo.GetByID(ctx, petID)
	if err != nil {
		return apperrors.Internal(err, "не удалось получить питомца")
	}

	if pet == nil {
		return apperrors.NewPetNotFoundError(petID)
	}

	// Удаляем питомца
	if err := s.petRepo.Delete(ctx, petID); err != nil {
		return apperrors.Internal(err, "не удалось удалить питомца")
	}

	return nil
}
