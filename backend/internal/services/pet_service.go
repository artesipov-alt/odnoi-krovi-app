package services

import (
	"context"
	"strings"
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/ent"
	"github.com/artesipov-alt/odnoi-krovi-app/ent/pet"
	"github.com/artesipov-alt/odnoi-krovi-app/ent/petanalysis"
	"github.com/artesipov-alt/odnoi-krovi-app/ent/pethealth"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	repositories "github.com/artesipov-alt/odnoi-krovi-app/internal/repositories"
)

// calculateAgeFields вычисляет возраст из даты рождения или дату из возраста
func calculateAgeFields(ageYears, ageMonths *int, birthDate **time.Time) {
	now := time.Now()
	if *birthDate != nil && **birthDate != (time.Time{}) {
		// Вычисляем возраст из даты рождения
		birth := **birthDate
		years := now.Year() - birth.Year()
		months := int(now.Month()) - int(birth.Month())
		if now.Day() < birth.Day() {
			months--
		}
		if months < 0 {
			years--
			months += 12
		}
		if ageYears != nil {
			*ageYears = years
		}
		if ageMonths != nil {
			*ageMonths = months
		}
	} else if ageYears != nil && ageMonths != nil && (*ageYears > 0 || *ageMonths > 0) {
		// Вычисляем дату рождения из возраста
		*birthDate = new(time.Time)
		**birthDate = now.AddDate(-*ageYears, -*ageMonths, 0)
	}
}

// PetService определяет интерфейс для бизнес-логики питомцев
type PetService interface {
	// CreatePet создает нового питомца для пользователя
	CreatePet(ctx context.Context, userID string, petData PetCreate) (*ent.Pet, error)

	// GetPetByID получает питомца по ID
	GetPetByID(ctx context.Context, petID string) (*ent.Pet, error)

	// GetUserPets получает всех питомцев пользователя
	GetUserPets(ctx context.Context, userID string) ([]*ent.Pet, error)

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
	Name            string              `json:"name" validate:"required,min=1,max=100"`
	ChipNumber      string              `json:"chipNumber,omitempty" validate:"omitempty,len=15"`
	PhotoURL        string              `json:"photoUrl,omitempty" validate:"omitempty,url,max=255"`
	BreedID         int                 `json:"breedId,omitempty" validate:"omitempty,min=1"`
	WeightKg        float64             `json:"weightKg,omitempty" validate:"omitempty,min=0"`
	AgeYears        int                 `json:"ageYears,omitempty" validate:"omitempty,min=0"`
	AgeMonths       int                 `json:"ageMonths,omitempty" validate:"omitempty,min=0,max=11"`
	BirthDate       *time.Time          `json:"birthDate,omitempty"`
	LivingCondition pet.LivingCondition `json:"livingCondition,omitempty"`
	Gender          pet.Gender          `json:"gender,omitempty"`
	Type            pet.Type            `json:"type" validate:"required"`
	BloodGroup      string              `json:"bloodGroup,omitempty" validate:"omitempty,max=50"`
	PetStatus       pet.PetStatus       `json:"petStatus" validate:"required"`

	// Вложенные структуры (DTO)
	Health     *ent.PetHealth     `json:"health,omitempty"`
	Treatments *ent.PetTreatment  `json:"treatments,omitempty"`
	Analyses   []*ent.PetAnalysis `json:"analyses,omitempty"`
	Bonuses    *ent.PetBonus      `json:"bonuses,omitempty"`
}

// PetUpdate содержит поля, которые можно обновить для питомца
type PetUpdate struct {
	Name            *string              `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	ChipNumber      *string              `json:"chipNumber,omitempty" validate:"omitempty,len=15"`
	PhotoURL        *string              `json:"photoUrl,omitempty" validate:"omitempty,url,max=255"`
	BreedID         *int                 `json:"breedId,omitempty" validate:"omitempty,min=1"`
	WeightKg        *float64             `json:"weightKg,omitempty" validate:"omitempty,min=0"`
	AgeYears        *int                 `json:"ageYears,omitempty" validate:"omitempty,min=0"`
	AgeMonths       *int                 `json:"ageMonths,omitempty" validate:"omitempty,min=0,max=11"`
	BirthDate       *time.Time           `json:"birthDate,omitempty"`
	LivingCondition *pet.LivingCondition `json:"livingCondition,omitempty"`
	Gender          *pet.Gender          `json:"gender,omitempty"`
	Type            *pet.Type            `json:"type,omitempty"`
	BloodGroup      *string              `json:"bloodGroup,omitempty" validate:"omitempty,max=50"`
	PetStatus       *pet.PetStatus       `json:"petStatus,omitempty" validate:"omitempty,max=50"`

	// Вложенные структуры
	Health     *ent.PetHealth     `json:"health,omitempty"`
	Treatments *ent.PetTreatment  `json:"treatments,omitempty"`
	Analyses   []*ent.PetAnalysis `json:"analyses,omitempty"`
	Bonuses    *ent.PetBonus      `json:"bonuses,omitempty"`
}

// PetServiceImpl реализует PetService
type PetServiceImpl struct {
	petRepo  repositories.PetRepository
	userRepo repositories.UserRepository
	storage  repositories.FileStorage
}

// NewPetService создает новый сервис питомцев
func NewPetService(petRepo repositories.PetRepository, userRepo repositories.UserRepository, storage repositories.FileStorage) *PetServiceImpl {
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
func (s *PetServiceImpl) CreatePet(ctx context.Context, userID string, petData PetCreate) (*ent.Pet, error) {
	// Проверяем, существует ли пользователь
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, apperrors.NewUserNotFoundError(userID)
		}
		return nil, apperrors.Internal(err, "не удалось проверить существование пользователя")
	}

	// Валидируем тип животного
	if err := pet.TypeValidator(petData.Type); err != nil {
		return nil, apperrors.ErrInvalidPetType
	}

	// Валидируем статус питомца
	if err := pet.PetStatusValidator(petData.PetStatus); err != nil {
		return nil, apperrors.ErrPetInvalidStatus
	}

	// Валидируем пол животного
	if petData.Gender != "" {
		if err := pet.GenderValidator(petData.Gender); err != nil {
			return nil, apperrors.ErrInvalidGender
		}
	}

	// Валидируем условия проживания
	if petData.LivingCondition != "" {
		if err := pet.LivingConditionValidator(petData.LivingCondition); err != nil {
			return nil, apperrors.ErrInvalidLivingCondition
		}
	}

	// Валидируем вложенные структуры, если они есть
	if petData.Health != nil {
		if petData.Health.HealthStatus != "" {
			if err := pethealth.HealthStatusValidator(petData.Health.HealthStatus); err != nil {
				return nil, apperrors.ErrInvalidHealthStatus
			}
		}
		if petData.Health.ReproductiveStatus != "" {
			if err := pethealth.ReproductiveStatusValidator(petData.Health.ReproductiveStatus); err != nil {
				return nil, apperrors.ErrInvalidReproductiveStatus
			}
		}
	}

	if petData.Analyses != nil {
		for _, a := range petData.Analyses {
			if a.LeukemiaType != "" {
				if err := petanalysis.LeukemiaTypeValidator(a.LeukemiaType); err != nil {
					return nil, apperrors.ErrInvalidAnalysisType
				}
			}
			if a.ImmunodeficiencyType != "" {
				if err := petanalysis.ImmunodeficiencyTypeValidator(a.ImmunodeficiencyType); err != nil {
					return nil, apperrors.ErrInvalidAnalysisType
				}
			}
			if a.HemoplasmosisType != "" {
				if err := petanalysis.HemoplasmosisTypeValidator(a.HemoplasmosisType); err != nil {
					return nil, apperrors.ErrInvalidAnalysisType
				}
			}
			if a.BartonellosisType != "" {
				if err := petanalysis.BartonellosisTypeValidator(a.BartonellosisType); err != nil {
					return nil, apperrors.ErrInvalidAnalysisType
				}
			}
			if a.BabesiosisType != "" {
				if err := petanalysis.BabesiosisTypeValidator(a.BabesiosisType); err != nil {
					return nil, apperrors.ErrInvalidAnalysisType
				}
			}
			if a.DirofilariaType != "" {
				if err := petanalysis.DirofilariaTypeValidator(a.DirofilariaType); err != nil {
					return nil, apperrors.ErrInvalidAnalysisType
				}
			}
			if a.EhrlichiosisType != "" {
				if err := petanalysis.EhrlichiosisTypeValidator(a.EhrlichiosisType); err != nil {
					return nil, apperrors.ErrInvalidAnalysisType
				}
			}
			if a.AnaplasmosisType != "" {
				if err := petanalysis.AnaplasmosisTypeValidator(a.AnaplasmosisType); err != nil {
					return nil, apperrors.ErrInvalidAnalysisType
				}
			}
		}
	}

	// Вычисляем возраст или дату рождения
	ageYears := petData.AgeYears
	ageMonths := petData.AgeMonths
	birthDate := petData.BirthDate
	calculateAgeFields(&ageYears, &ageMonths, &birthDate)

	// Создаем нового питомца (ENT entity)
	p := &ent.Pet{
		UserID:          userID,
		Name:            petData.Name,
		ChipNumber:      petData.ChipNumber,
		PhotoURL:        petData.PhotoURL,
		BreedID:         petData.BreedID,
		WeightKg:        petData.WeightKg,
		AgeYears:        ageYears,
		AgeMonths:       ageMonths,
		BirthDate:       birthDate,
		LivingCondition: petData.LivingCondition,
		Gender:          petData.Gender,
		Type:            petData.Type,
		BloodGroup:      petData.BloodGroup,
		PetStatus:       petData.PetStatus,
	}

	newPet, err := s.petRepo.Create(ctx, p, petData.Health, petData.Treatments, petData.Analyses, petData.Bonuses)
	if err != nil {
		return nil, apperrors.Internal(err, "не удалось создать питомца")
	}

	// Преобразуем путь к фото в полный URL
	newPet.PhotoURL = s.buildFullPhotoURL(newPet.PhotoURL)

	return newPet, nil
}

// GetPetByID получает питомца по ID
func (s *PetServiceImpl) GetPetByID(ctx context.Context, petID string) (*ent.Pet, error) {
	p, err := s.petRepo.GetByID(ctx, petID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, apperrors.NewPetNotFoundError(petID)
		}
		return nil, apperrors.Internal(err, "не удалось получить питомца")
	}

	// Преобразуем путь к фото в полный URL
	p.PhotoURL = s.buildFullPhotoURL(p.PhotoURL)

	return p, nil
}

// GetUserPets получает всех питомцев пользователя
func (s *PetServiceImpl) GetUserPets(ctx context.Context, userID string) ([]*ent.Pet, error) {
	// Проверяем, существует ли пользователь
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, apperrors.NewUserNotFoundError(userID)
		}
		return nil, apperrors.Internal(err, "не удалось проверить существование пользователя")
	}

	pets, err := s.petRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, apperrors.Internal(err, "не удалось получить питомцев пользователя")
	}

	// Преобразуем пути к фото в полные URL
	for _, p := range pets {
		if p != nil {
			p.PhotoURL = s.buildFullPhotoURL(p.PhotoURL)
		}
	}

	return pets, nil
}

// UpdatePet обновляет информацию о питомце
func (s *PetServiceImpl) UpdatePet(ctx context.Context, petID string, updates PetUpdate) error {
	// Получаем существующего питомца
	p, err := s.petRepo.GetByID(ctx, petID)
	if err != nil {
		if ent.IsNotFound(err) {
			return apperrors.NewPetNotFoundError(petID)
		}
		return apperrors.Internal(err, "не удалось получить питомца")
	}

	// Применяем обновления и валидируем
	if updates.Name != nil {
		p.Name = *updates.Name
	}
	if updates.ChipNumber != nil {
		p.ChipNumber = *updates.ChipNumber
	}
	if updates.PhotoURL != nil {
		p.PhotoURL = *updates.PhotoURL
	}
	if updates.BreedID != nil {
		p.BreedID = *updates.BreedID
	}
	if updates.WeightKg != nil {
		p.WeightKg = *updates.WeightKg
	}
	if updates.AgeYears != nil {
		p.AgeYears = *updates.AgeYears
	}
	if updates.AgeMonths != nil {
		p.AgeMonths = *updates.AgeMonths
	}
	if updates.BirthDate != nil {
		p.BirthDate = updates.BirthDate
	}
	if updates.LivingCondition != nil {
		if err := pet.LivingConditionValidator(*updates.LivingCondition); err != nil {
			return apperrors.ErrInvalidLivingCondition
		}
		p.LivingCondition = *updates.LivingCondition
	}
	if updates.Gender != nil {
		if err := pet.GenderValidator(*updates.Gender); err != nil {
			return apperrors.ErrInvalidGender
		}
		p.Gender = *updates.Gender
	}
	if updates.Type != nil {
		if err := pet.TypeValidator(*updates.Type); err != nil {
			return apperrors.ErrInvalidPetType
		}
		p.Type = *updates.Type
	}
	if updates.BloodGroup != nil {
		p.BloodGroup = *updates.BloodGroup
	}
	if updates.PetStatus != nil {
		if err := pet.PetStatusValidator(*updates.PetStatus); err != nil {
			return apperrors.ErrPetInvalidStatus
		}
		p.PetStatus = *updates.PetStatus
	}

	// Вычисляем возраст или дату рождения при обновлении
	if updates.BirthDate != nil || (updates.AgeYears != nil && updates.AgeMonths != nil) {
		calculateAgeFields(updates.AgeYears, updates.AgeMonths, &updates.BirthDate)
		if updates.BirthDate != nil {
			p.BirthDate = updates.BirthDate
		}
	}

	// Валидируем вложенные структуры
	if updates.Health != nil {
		h := updates.Health
		if h.HealthStatus != "" {
			if err := pethealth.HealthStatusValidator(h.HealthStatus); err != nil {
				return apperrors.ErrInvalidHealthStatus
			}
		}
		if h.ReproductiveStatus != "" {
			if err := pethealth.ReproductiveStatusValidator(h.ReproductiveStatus); err != nil {
				return apperrors.ErrInvalidReproductiveStatus
			}
		}
	}
	if updates.Analyses != nil {
		for _, a := range updates.Analyses {
			if a.LeukemiaType != "" {
				if err := petanalysis.LeukemiaTypeValidator(a.LeukemiaType); err != nil {
					return apperrors.ErrInvalidAnalysisType
				}
			}
			if a.ImmunodeficiencyType != "" {
				if err := petanalysis.ImmunodeficiencyTypeValidator(a.ImmunodeficiencyType); err != nil {
					return apperrors.ErrInvalidAnalysisType
				}
			}
			if a.HemoplasmosisType != "" {
				if err := petanalysis.HemoplasmosisTypeValidator(a.HemoplasmosisType); err != nil {
					return apperrors.ErrInvalidAnalysisType
				}
			}
			if a.BartonellosisType != "" {
				if err := petanalysis.BartonellosisTypeValidator(a.BartonellosisType); err != nil {
					return apperrors.ErrInvalidAnalysisType
				}
			}
			if a.BabesiosisType != "" {
				if err := petanalysis.BabesiosisTypeValidator(a.BabesiosisType); err != nil {
					return apperrors.ErrInvalidAnalysisType
				}
			}
			if a.DirofilariaType != "" {
				if err := petanalysis.DirofilariaTypeValidator(a.DirofilariaType); err != nil {
					return apperrors.ErrInvalidAnalysisType
				}
			}
			if a.EhrlichiosisType != "" {
				if err := petanalysis.EhrlichiosisTypeValidator(a.EhrlichiosisType); err != nil {
					return apperrors.ErrInvalidAnalysisType
				}
			}
			if a.AnaplasmosisType != "" {
				if err := petanalysis.AnaplasmosisTypeValidator(a.AnaplasmosisType); err != nil {
					return apperrors.ErrInvalidAnalysisType
				}
			}
		}
	}

	// Сохраняем обновленного питомца
	if _, err := s.petRepo.Update(ctx, p, updates.Health, updates.Treatments, updates.Analyses, updates.Bonuses); err != nil {
		return apperrors.Internal(err, "не удалось обновить питомца")
	}

	return nil
}

// DeletePet удаляет питомца по ID
func (s *PetServiceImpl) DeletePet(ctx context.Context, petID string) error {
	exists, err := s.petRepo.ExistsByID(ctx, petID)
	if err != nil {
		return apperrors.Internal(err, "не удалось проверить существование питомца")
	}
	if !exists {
		return apperrors.NewPetNotFoundError(petID)
	}

	if err := s.petRepo.Delete(ctx, petID); err != nil {
		return apperrors.Internal(err, "не удалось удалить питомца")
	}

	return nil
}

// GetAvatarUploadURL Возвращает ссылку для загрузки аватарки питомца.
func (s *PetServiceImpl) GetAvatarUploadURL(ctx context.Context, petID string) (string, string, error) {
	exists, err := s.petRepo.ExistsByID(ctx, petID)
	if err != nil {
		return "", "", apperrors.Internal(err, "не удалось проверить существование питомца")
	}
	if !exists {
		return "", "", apperrors.NewPetNotFoundError(petID)
	}

	url, path, err := s.storage.GetAvatarUploadInfo(ctx, petID)
	if err != nil {
		return "", "", apperrors.Internal(err, "не удалось сгенерировать URL для загрузки аватара")
	}

	return url, path, nil
}

// UpdatePetAvatar обновляет аватар питомца и делает его публичным в хранилище
func (s *PetServiceImpl) UpdatePetAvatar(ctx context.Context, avatarPath string) (string, error) {
	decodedPath := strings.ReplaceAll(avatarPath, "%2F", "/")
	parts := strings.Split(decodedPath, "/")
	if len(parts) < 2 {
		return "", apperrors.Internal(nil, "некорректный формат пути аватара")
	}
	petID := parts[1]

	p, err := s.petRepo.GetByID(ctx, petID)
	if err != nil {
		if ent.IsNotFound(err) {
			return "", apperrors.NewPetNotFoundError(petID)
		}
		return "", apperrors.Internal(err, "не удалось получить питомца")
	}

	exists, err := s.storage.CheckObjectExists(ctx, decodedPath)
	if err != nil {
		return "", apperrors.Internal(err, "не удалось проверить существование файла")
	}
	if !exists {
		return "", apperrors.Internal(nil, "файл не найден")
	}

	err = s.storage.SetObjectPublicACL(ctx, decodedPath)
	if err != nil {
		return "", apperrors.Internal(err, "не удалось установить публичный ACL")
	}

	publicURL := s.storage.GetAvatarPublicURL(petID)
	if publicURL == "" {
		return "", apperrors.Internal(nil, "не удалось получить публичный URL для аватарки")
	}

	p.PhotoURL = decodedPath
	if _, err := s.petRepo.Update(ctx, p, nil, nil, nil, nil); err != nil {
		return "", apperrors.Internal(err, "не удалось обновить photoURL питомца")
	}

	return publicURL, nil
}
