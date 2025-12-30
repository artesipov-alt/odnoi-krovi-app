package services

import (
	"context"
	"errors"
	"strings"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/models"
	repositories "github.com/artesipov-alt/odnoi-krovi-app/internal/repositories"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/storage"
	validation "github.com/artesipov-alt/odnoi-krovi-app/internal/utils/enums"
	"github.com/artesipov-alt/odnoi-krovi-app/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// PetService определяет интерфейс для бизнес-логики питомцев
type PetService interface {
	// CreatePet создает нового питомца для пользователя
	CreatePet(ctx context.Context, userID string, petData PetCreate) (*models.Pet, error)

	// GetPetByID получает питомца по ID
	GetPetByID(ctx context.Context, petID string) (*models.Pet, error)

	// GetUserPets получает всех питомцев пользователя
	GetUserPets(ctx context.Context, userID string) ([]*models.Pet, error)

	// UpdatePet обновляет информацию о питомце
	UpdatePet(ctx context.Context, petID string, updates PetUpdate) error

	// DeletePet удаляет питомца по ID
	DeletePet(ctx context.Context, petID string) error

	// GetAvatarUploadURL Возвращает ссылку для загрузки аватарки питомца.
	GetAvatarUploadURL(ctx context.Context, petID string) (string, string, error)

	// UpdatePetAvatar обновляет аватар питомца и делает его публичным в хранилище
	UpdatePetAvatar(ctx context.Context, avatarPath string) (string, error)
}

// PetCreate содержит данные для создания питомца
type PetCreate struct {
	Name            string                 `json:"name" validate:"required,min=1,max=100"`
	ChipNumber      string                 `json:"chipNumber,omitempty" validate:"omitempty,len=15"`
	PhotoURL        string                 `json:"photoUrl,omitempty" validate:"omitempty,url,max=255"`
	Breed           string                 `json:"breed,omitempty" validate:"omitempty,max=100"`
	WeightKg        float64                `json:"weightKg,omitempty" validate:"omitempty,min=0"`
	AgeYears        int                    `json:"ageYears,omitempty" validate:"omitempty,min=0"`
	AgeMonths       int                    `json:"ageMonths,omitempty" validate:"omitempty,min=0,max=11"`
	LivingCondition models.LivingCondition `json:"livingCondition,omitempty"`
	Gender          models.Gender          `json:"gender,omitempty"`
	Type            models.PetType         `json:"type,omitempty"`
	BloodGroup      string                 `json:"bloodGroup,omitempty" validate:"omitempty,max=50"`
	PetStatus       models.PetRole         `json:"petStatus,omitempty" validate:"omitempty,max=50"`
}

// PetUpdate содержит поля, которые можно обновить для питомца
type PetUpdate struct {
	Name            *string                 `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	ChipNumber      *string                 `json:"chipNumber,omitempty" validate:"omitempty,len=15"`
	PhotoURL        *string                 `json:"photoUrl,omitempty" validate:"omitempty,url,max=255"`
	Breed           *string                 `json:"breed,omitempty" validate:"omitempty,max=100"`
	WeightKg        *float64                `json:"weightKg,omitempty" validate:"omitempty,min=0"`
	AgeYears        *int                    `json:"ageYears,omitempty" validate:"omitempty,min=0"`
	AgeMonths       *int                    `json:"ageMonths,omitempty" validate:"omitempty,min=0,max=11"`
	LivingCondition *models.LivingCondition `json:"livingCondition,omitempty"`
	Gender          *models.Gender          `json:"gender,omitempty"`
	Type            *models.PetType         `json:"type,omitempty"`
	BloodGroup      *string                 `json:"bloodGroup,omitempty" validate:"omitempty,max=50"`
	PetStatus       *models.PetRole         `json:"petStatus,omitempty" validate:"omitempty,max=50"`
}

// PetServiceImpl реализует PetService
type PetServiceImpl struct {
	petRepo  repositories.PetRepository
	userRepo repositories.UserRepository
	storage  storage.FileStorage
}

// NewPetService создает новый сервис питомцев
func NewPetService(petRepo repositories.PetRepository, userRepo repositories.UserRepository, storage storage.FileStorage) *PetServiceImpl {
	return &PetServiceImpl{
		petRepo:  petRepo,
		userRepo: userRepo,
		storage:  storage,
	}
}

// buildFullPhotoURL преобразует путь к фото в полный публичный URL
func (s *PetServiceImpl) buildFullPhotoURL(path string) string {
	if path == "" {
		return ""
	}
	return s.storage.GetPublicURLFromPath(path)
}

// CreatePet создает нового питомца для пользователя
func (s *PetServiceImpl) CreatePet(ctx context.Context, userID string, petData PetCreate) (*models.Pet, error) {
	// Проверяем, существует ли пользователь
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		// Если пользователь не найден - возвращаем 404, а не 500
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewUserNotFoundError(userID)
		}
		return nil, apperrors.Internal(err, "не удалось проверить существование пользователя")
	}

	if user == nil {
		return nil, apperrors.NewUserNotFoundError(userID)
	}

	// Валидируем тип животного
	_, err = validation.LocalizePetType(string(petData.Type))
	if err != nil {
		return nil, apperrors.ErrInvalidPetType
	}

	// Валидируем пол животного
	_, err = validation.LocalizeGender(string(petData.Gender))
	if err != nil {
		return nil, apperrors.ErrInvalidGender
	}

	// Валидируем условия проживания
	_, err = validation.LocalizeLivingCondition(string(petData.LivingCondition))
	if err != nil {
		return nil, apperrors.ErrInvalidLivingCondition
	}

	// Создаем нового питомца
	pet := &models.Pet{
		OwnerID:         userID,
		Name:            petData.Name,
		ChipNumber:      petData.ChipNumber,
		PhotoURL:        petData.PhotoURL,
		Breed:           petData.Breed,
		WeightKg:        petData.WeightKg,
		AgeYears:        petData.AgeYears,
		AgeMonths:       petData.AgeMonths,
		LivingCondition: petData.LivingCondition,
		Gender:          petData.Gender,
		Type:            petData.Type,
		BloodGroup:      petData.BloodGroup,
		PetStatus:       petData.PetStatus,
	}

	if err := s.petRepo.Create(ctx, pet); err != nil {
		return nil, apperrors.Internal(err, "не удалось создать питомца")
	}

	// Преобразуем путь к фото в полный URL
	pet.PhotoURL = s.buildFullPhotoURL(pet.PhotoURL)

	return pet, nil
}

// GetPetByID получает питомца по ID
func (s *PetServiceImpl) GetPetByID(ctx context.Context, petID string) (*models.Pet, error) {
	pet, err := s.petRepo.GetByID(ctx, petID)
	if err != nil {
		// Если питомец не найден - возвращаем 404, а не 500
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewPetNotFoundError(petID)
		}
		return nil, apperrors.Internal(err, "не удалось получить питомца")
	}

	if pet == nil {
		return nil, apperrors.NewPetNotFoundError(petID)
	}

	// Преобразуем путь к фото в полный URL
	pet.PhotoURL = s.buildFullPhotoURL(pet.PhotoURL)

	return pet, nil
}

// GetUserPets получает всех питомцев пользователя
func (s *PetServiceImpl) GetUserPets(ctx context.Context, userID string) ([]*models.Pet, error) {
	// Проверяем, существует ли пользователь
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		// Если пользователь не найден - возвращаем 404, а не 500
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewUserNotFoundError(userID)
		}
		return nil, apperrors.Internal(err, "не удалось проверить существование пользователя")
	}

	if user == nil {
		return nil, apperrors.NewUserNotFoundError(userID)
	}

	pets, err := s.petRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, apperrors.Internal(err, "не удалось получить питомцев пользователя")
	}

	// Преобразуем пути к фото в полные URL для каждого питомца
	for i := range pets {
		if pets[i] != nil {
			pets[i].PhotoURL = s.buildFullPhotoURL(pets[i].PhotoURL)
		}
	}

	return pets, nil
}

// UpdatePet обновляет информацию о питомце
func (s *PetServiceImpl) UpdatePet(ctx context.Context, petID string, updates PetUpdate) error {
	// Получаем существующего питомца
	pet, err := s.petRepo.GetByID(ctx, petID)
	if err != nil {
		// Если питомец не найден - возвращаем 404, а не 500
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.NewPetNotFoundError(petID)
		}
		return apperrors.Internal(err, "не удалось получить питомца")
	}

	if pet == nil {
		return apperrors.NewPetNotFoundError(petID)
	}

	// Применяем обновления
	if updates.Name != nil {
		pet.Name = *updates.Name
	}
	if updates.ChipNumber != nil {
		pet.ChipNumber = *updates.ChipNumber
	}
	if updates.PhotoURL != nil {
		pet.PhotoURL = *updates.PhotoURL
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
	if updates.LivingCondition != nil {
		_, err := validation.LocalizeLivingCondition(string(*updates.LivingCondition))
		if err != nil {
			return apperrors.ErrInvalidLivingCondition
		}
		pet.LivingCondition = *updates.LivingCondition
	}
	if updates.Gender != nil {
		_, err := validation.LocalizeGender(string(*updates.Gender))
		if err != nil {
			return apperrors.ErrInvalidGender
		}
		pet.Gender = *updates.Gender
	}
	if updates.Type != nil {
		_, err := validation.LocalizePetType(string(*updates.Type))
		if err != nil {
			return apperrors.ErrInvalidPetType
		}
		pet.Type = *updates.Type
	}
	if updates.BloodGroup != nil {
		pet.BloodGroup = *updates.BloodGroup
	}
	if updates.PetStatus != nil {
		pet.PetStatus = *updates.PetStatus
	}

	// Сохраняем обновленного питомца
	if err := s.petRepo.Update(ctx, pet); err != nil {
		return apperrors.Internal(err, "не удалось обновить питомца")
	}

	return nil
}

// DeletePet удаляет питомца по ID
func (s *PetServiceImpl) DeletePet(ctx context.Context, petID string) error {
	// Проверяем, существует ли питомец
	pet, err := s.petRepo.GetByID(ctx, petID)
	if err != nil {
		// Если питомец не найден - возвращаем 404, а не 500
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.NewPetNotFoundError(petID)
		}
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

// GetAvatarUploadURL Возвращает ссылку для загрузки аватарки питомца.
func (s *PetServiceImpl) GetAvatarUploadURL(ctx context.Context, petID string) (string, string, error) {
	// Проверяем, существует ли питомец
	pet, err := s.petRepo.GetByID(ctx, petID)
	if err != nil {
		// Если питомец не найден - возвращаем 404, а не 500
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", "", apperrors.NewPetNotFoundError(petID)
		}
		return "", "", apperrors.Internal(err, "не удалось получить питомца")
	}

	if pet == nil {
		return "", "", apperrors.NewPetNotFoundError(petID)
	}

	// Генерируем ссылку питомца
	url, path, err := s.storage.GetAvatarUploadInfo(ctx, petID)
	if err != nil {
		return "", "", apperrors.Internal(err, "не удалось сгенерировать URL для загрузки аватара")
	}

	return url, path, nil
}

// UpdatePetAvatar обновляет аватар питомца и делает его публичным в хранилище
func (s *PetServiceImpl) UpdatePetAvatar(ctx context.Context, avatarPath string) (string, error) {
	// Извлекаем petID из пути. Пример: "pets%2FPET-25-000006%2Favatar.jpg"
	// Сначала декодируем URL, затем разбиваем по "/"
	decodedPath := strings.ReplaceAll(avatarPath, "%2F", "/")

	parts := strings.Split(decodedPath, "/")
	logger.Log.Debug("Parts:", zap.Strings("p", parts))
	if len(parts) < 2 {
		return "", apperrors.Internal(nil, "некорректный формат пути аватара")
	}
	petID := parts[1]

	// Проверяем, существует ли питомец
	pet, err := s.petRepo.GetByID(ctx, petID)
	if err != nil {
		// Если питомец не найден - возвращаем 404, а не 500
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", apperrors.NewPetNotFoundError(petID)
		}
		return "", apperrors.Internal(err, "не удалось получить питомца")
	}

	if pet == nil {
		return "", apperrors.NewPetNotFoundError(petID)
	}

	// Проверяем, что файл действительно загружен
	exists, err := s.storage.CheckObjectExists(ctx, decodedPath)
	if err != nil {
		return "", apperrors.Internal(err, "не удалось проверить существование файла")
	}
	if !exists {
		return "", apperrors.Internal(nil, "файл не найден")
	}

	// Устанавливаем публичный ACL для объекта
	err = s.storage.SetObjectPublicACL(ctx, decodedPath)
	if err != nil {
		return "", apperrors.Internal(err, "не удалось установить публичный ACL")
	}

	// Получаем публичный URL для аватарки
	publicURL := s.storage.GetAvatarPublicURL(petID)
	if publicURL == "" {
		return "", apperrors.Internal(nil, "не удалось получить публичный URL для аватарки")
	}

	// Обновляем photoURL в базе данных
	pet.PhotoURL = decodedPath
	if err := s.petRepo.Update(ctx, pet); err != nil {
		return "", apperrors.Internal(err, "не удалось обновить photoURL питомца")
	}

	return publicURL, nil
}
