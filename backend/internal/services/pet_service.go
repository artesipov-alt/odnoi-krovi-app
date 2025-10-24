package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/models"
	repositories "github.com/artesipov-alt/odnoi-krovi-app/internal/repositories/interfaces"
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
	HasChip         bool                   `json:"has_chip"`
	ChipNumber      string                 `json:"chip_number,omitempty" validate:"omitempty,max=50"`
	PhotoURL        string                 `json:"photo_url,omitempty" validate:"omitempty,url,max=255"`
	KnowsBloodGroup bool                   `json:"knows_blood_group"`
	IsGuideDog      bool                   `json:"is_guide_dog"`
	IsTherapist     bool                   `json:"is_therapist"`
	Breed           string                 `json:"breed,omitempty" validate:"omitempty,max=100"`
	WeightKg        float64                `json:"weight_kg,omitempty" validate:"omitempty,min=0"`
	AgeYears        int                    `json:"age_years,omitempty" validate:"omitempty,min=0"`
	AgeMonths       int                    `json:"age_months,omitempty" validate:"omitempty,min=0,max=11"`
	Sterilized      bool                   `json:"sterilized"`
	Latitude        float64                `json:"latitude,omitempty" validate:"omitempty,latitude"`
	Longitude       float64                `json:"longitude,omitempty" validate:"omitempty,longitude"`
	LivingCondition models.LivingCondition `json:"living_condition,omitempty"`
	Gender          models.Gender          `json:"gender,omitempty"`
	Type            models.PetType         `json:"type,omitempty"`
	BloodGroup      string                 `json:"blood_group,omitempty" validate:"omitempty,max=50"`
}

// PetUpdate содержит поля, которые можно обновить для питомца
type PetUpdate struct {
	Name            *string                 `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	HasChip         *bool                   `json:"has_chip,omitempty"`
	ChipNumber      *string                 `json:"chip_number,omitempty" validate:"omitempty,max=50"`
	PhotoURL        *string                 `json:"photo_url,omitempty" validate:"omitempty,url,max=255"`
	KnowsBloodGroup *bool                   `json:"knows_blood_group,omitempty"`
	IsGuideDog      *bool                   `json:"is_guide_dog,omitempty"`
	IsTherapist     *bool                   `json:"is_therapist,omitempty"`
	Breed           *string                 `json:"breed,omitempty" validate:"omitempty,max=100"`
	WeightKg        *float64                `json:"weight_kg,omitempty" validate:"omitempty,min=0"`
	AgeYears        *int                    `json:"age_years,omitempty" validate:"omitempty,min=0"`
	AgeMonths       *int                    `json:"age_months,omitempty" validate:"omitempty,min=0,max=11"`
	Sterilized      *bool                   `json:"sterilized,omitempty"`
	Latitude        *float64                `json:"latitude,omitempty" validate:"omitempty,latitude"`
	Longitude       *float64                `json:"longitude,omitempty" validate:"omitempty,longitude"`
	LivingCondition *models.LivingCondition `json:"living_condition,omitempty"`
	Gender          *models.Gender          `json:"gender,omitempty"`
	Type            *models.PetType         `json:"type,omitempty"`
	BloodGroup      *string                 `json:"blood_group,omitempty" validate:"omitempty,max=50"`
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
		return nil, fmt.Errorf("проверка существования пользователя: %w", err)
	}

	if user == nil {
		return nil, errors.New("пользователь не найден")
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
		LivingCondition: petData.LivingCondition,
		Gender:          petData.Gender,
		Type:            petData.Type,
		BloodGroup:      petData.BloodGroup,
	}

	if err := s.petRepo.Create(ctx, pet); err != nil {
		return nil, fmt.Errorf("создание питомца: %w", err)
	}

	return pet, nil
}

// GetPetByID получает питомца по ID
func (s *PetServiceImpl) GetPetByID(ctx context.Context, petID int) (*models.Pet, error) {
	pet, err := s.petRepo.GetByID(ctx, petID)
	if err != nil {
		return nil, fmt.Errorf("получение питомца: %w", err)
	}

	return pet, nil
}

// GetUserPets получает всех питомцев пользователя
func (s *PetServiceImpl) GetUserPets(ctx context.Context, userID int) ([]*models.Pet, error) {
	// Проверяем, существует ли пользователь
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("проверка существования пользователя: %w", err)
	}

	if user == nil {
		return nil, errors.New("пользователь не найден")
	}

	pets, err := s.petRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("получение питомцев пользователя: %w", err)
	}

	return pets, nil
}

// UpdatePet обновляет информацию о питомце
func (s *PetServiceImpl) UpdatePet(ctx context.Context, petID int, updates PetUpdate) error {
	// Получаем существующего питомца
	pet, err := s.petRepo.GetByID(ctx, petID)
	if err != nil {
		return fmt.Errorf("получение питомца для обновления: %w", err)
	}

	if pet == nil {
		return errors.New("питомец не найден")
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
		pet.LivingCondition = *updates.LivingCondition
	}
	if updates.Gender != nil {
		pet.Gender = *updates.Gender
	}
	if updates.Type != nil {
		pet.Type = *updates.Type
	}
	if updates.BloodGroup != nil {
		pet.BloodGroup = *updates.BloodGroup
	}

	// Сохраняем обновленного питомца
	if err := s.petRepo.Update(ctx, pet); err != nil {
		return fmt.Errorf("обновление питомца: %w", err)
	}

	return nil
}

// DeletePet удаляет питомца по ID
func (s *PetServiceImpl) DeletePet(ctx context.Context, petID int) error {
	// Проверяем, существует ли питомец
	pet, err := s.petRepo.GetByID(ctx, petID)
	if err != nil {
		return fmt.Errorf("получение питомца для удаления: %w", err)
	}

	if pet == nil {
		return errors.New("питомец не найден")
	}

	// Удаляем питомца
	if err := s.petRepo.Delete(ctx, petID); err != nil {
		return fmt.Errorf("удаление питомца: %w", err)
	}

	return nil
}
