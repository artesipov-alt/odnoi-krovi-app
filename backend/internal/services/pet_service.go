package services

import (
	"context"
	"errors"
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
	CreatePet(ctx context.Context, userID string, pet *ent.Pet) (*ent.Pet, error)

	// GetPetByID получает питомца по ID с preload связей
	GetPetByID(ctx context.Context, petID string, preloads ...string) (*ent.Pet, error)

	// GetUserPets получает всех питомцев пользователя с preload связей
	GetUserPets(ctx context.Context, userID string, preloads ...string) ([]*ent.Pet, error)

	// UpdatePet обновляет информацию о питомце
	UpdatePet(ctx context.Context, petID string, updates map[string]interface{}, health *ent.PetHealth, treatments *ent.PetTreatment, analyses []*ent.PetAnalysis, bonuses *ent.PetBonus) error

	// DeletePet удаляет питомца по ID
	DeletePet(ctx context.Context, petID string) error

	// GetAvatarUploadURL Возвращает ссылку для загрузки аватарки питомца.
	GetAvatarUploadURL(ctx context.Context, petID string) (string, string, error)

	// UpdatePetAvatar обновляет аватар питомца и делает его публичным в хранилище
	UpdatePetAvatar(ctx context.Context, avatarPath string) (string, error)
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
func (s *PetServiceImpl) CreatePet(ctx context.Context, userID string, petData *ent.Pet) (*ent.Pet, error) {
	// Проверяем, существует ли пользователь
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, apperrors.NewUserNotFoundError(userID)
		}
		return nil, apperrors.Internal(err, "не удалось проверить существование пользователя")
	}

	// Устанавливаем UserID
	petData.UserID = userID

	// Валидируем тип животного
	if err := pet.TypeValidator(petData.Type); err != nil {
		return nil, apperrors.ErrInvalidPetType
	}

	// Валидируем статус питомца
	if err := pet.PetStatusValidator(petData.PetStatus); err != nil {
		return nil, apperrors.ErrPetInvalidStatus
	}

	// Валидируем пол животного
	if string(petData.Gender) != "" {
		if err := pet.GenderValidator(petData.Gender); err != nil {
			return nil, apperrors.ErrInvalidGender
		}
	}

	// Валидируем условия проживания
	if string(petData.LivingCondition) != "" {
		if err := pet.LivingConditionValidator(petData.LivingCondition); err != nil {
			return nil, apperrors.ErrInvalidLivingCondition
		}
	}

	// Валидируем вложенные структуры, если они есть
	if petData.Edges.Health != nil {
		if string(petData.Edges.Health.HealthStatus) != "" {
			if err := pethealth.HealthStatusValidator(petData.Edges.Health.HealthStatus); err != nil {
				return nil, apperrors.ErrInvalidHealthStatus
			}
		}
		if string(petData.Edges.Health.ReproductiveStatus) != "" {
			if err := pethealth.ReproductiveStatusValidator(petData.Edges.Health.ReproductiveStatus); err != nil {
				return nil, apperrors.ErrInvalidReproductiveStatus
			}
		}
	}

	if petData.Edges.Analyses != nil {
		for _, a := range petData.Edges.Analyses {
			if string(a.LeukemiaType) != "" {
				if err := petanalysis.LeukemiaTypeValidator(a.LeukemiaType); err != nil {
					return nil, apperrors.ErrInvalidAnalysisType
				}
			}
			if string(a.ImmunodeficiencyType) != "" {
				if err := petanalysis.ImmunodeficiencyTypeValidator(a.ImmunodeficiencyType); err != nil {
					return nil, apperrors.ErrInvalidAnalysisType
				}
			}
			if string(a.HemoplasmosisType) != "" {
				if err := petanalysis.HemoplasmosisTypeValidator(a.HemoplasmosisType); err != nil {
					return nil, apperrors.ErrInvalidAnalysisType
				}
			}
			if string(a.BartonellosisType) != "" {
				if err := petanalysis.BartonellosisTypeValidator(a.BartonellosisType); err != nil {
					return nil, apperrors.ErrInvalidAnalysisType
				}
			}
			if string(a.BabesiosisType) != "" {
				if err := petanalysis.BabesiosisTypeValidator(a.BabesiosisType); err != nil {
					return nil, apperrors.ErrInvalidAnalysisType
				}
			}
			if string(a.DirofilariaType) != "" {
				if err := petanalysis.DirofilariaTypeValidator(a.DirofilariaType); err != nil {
					return nil, apperrors.ErrInvalidAnalysisType
				}
			}
			if string(a.EhrlichiosisType) != "" {
				if err := petanalysis.EhrlichiosisTypeValidator(a.EhrlichiosisType); err != nil {
					return nil, apperrors.ErrInvalidAnalysisType
				}
			}
			if string(a.AnaplasmosisType) != "" {
				if err := petanalysis.AnaplasmosisTypeValidator(a.AnaplasmosisType); err != nil {
					return nil, apperrors.ErrInvalidAnalysisType
				}
			}
		}
	}

	newPet, err := s.petRepo.Create(ctx, petData, petData.Edges.Health, petData.Edges.Treatments, petData.Edges.Analyses, petData.Edges.Bonuses)
	if err != nil {
		return nil, apperrors.Internal(err, "не удалось создать питомца")
	}

	// Преобразуем путь к фото в полный URL
	newPet.PhotoURL = s.buildFullPhotoURL(newPet.PhotoURL)

	return newPet, nil
}

// GetPetByID получает питомца по ID с preload связей
func (s *PetServiceImpl) GetPetByID(ctx context.Context, petID string, preloads ...string) (*ent.Pet, error) {
	if petID == "" {
		return nil, errors.New("invalid pet ID")
	}

	p, err := s.petRepo.GetByID(ctx, petID, preloads...)
	if err != nil {
		return nil, apperrors.Internal(err, "не удалось получить питомца")
	}

	// Преобразуем путь к фото в полный URL
	p.PhotoURL = s.buildFullPhotoURL(p.PhotoURL)

	return p, nil
}

// GetUserPets получает всех питомцев пользователя с preload связей
func (s *PetServiceImpl) GetUserPets(ctx context.Context, userID string, preloads ...string) ([]*ent.Pet, error) {
	// Проверяем, существует ли пользователь
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, apperrors.NewUserNotFoundError(userID)
		}
		return nil, apperrors.Internal(err, "не удалось проверить существование пользователя")
	}

	pets, err := s.petRepo.GetByUserID(ctx, userID, preloads...)
	if err != nil {
		return nil, apperrors.Internal(err, "не удалось получить питомцев пользователя")
	}

	// Преобразуем пути к фото в полные URL
	for _, pet := range pets {
		pet.PhotoURL = s.buildFullPhotoURL(pet.PhotoURL)
	}

	return pets, nil
}

// UpdatePet обновляет информацию о питомце
func (s *PetServiceImpl) UpdatePet(ctx context.Context, petID string, updates map[string]interface{}, health *ent.PetHealth, treatments *ent.PetTreatment, analyses []*ent.PetAnalysis, bonuses *ent.PetBonus) error {
	// Получаем существующего питомца
	p, err := s.petRepo.GetByID(ctx, petID)
	if err != nil {
		if ent.IsNotFound(err) {
			return apperrors.NewPetNotFoundError(petID)
		}
		return apperrors.Internal(err, "не удалось получить питомца")
	}

	// Применяем обновления и валидируем
	if val, ok := updates["Name"]; ok {
		p.Name = val.(string)
	}
	if val, ok := updates["ChipNumber"]; ok {
		p.ChipNumber = val.(string)
	}
	if val, ok := updates["PhotoURL"]; ok {
		p.PhotoURL = val.(string)
	}
	if val, ok := updates["BreedID"]; ok {
		p.BreedID = val.(int)
	}
	if val, ok := updates["WeightKg"]; ok {
		p.WeightKg = val.(float64)
	}
	if val, ok := updates["AgeYears"]; ok {
		p.AgeYears = val.(int)
	}
	if val, ok := updates["AgeMonths"]; ok {
		p.AgeMonths = val.(int)
	}
	if val, ok := updates["BirthDate"]; ok {
		p.BirthDate = val.(*time.Time)
	}
	if val, ok := updates["LivingCondition"]; ok {
		lc := val.(string)
		if err := pet.LivingConditionValidator(pet.LivingCondition(lc)); err != nil {
			return apperrors.ErrInvalidLivingCondition
		}
		p.LivingCondition = pet.LivingCondition(lc)
	}
	if val, ok := updates["Gender"]; ok {
		g := val.(string)
		if err := pet.GenderValidator(pet.Gender(g)); err != nil {
			return apperrors.ErrInvalidGender
		}
		p.Gender = pet.Gender(g)
	}
	if val, ok := updates["Type"]; ok {
		t := val.(string)
		if err := pet.TypeValidator(pet.Type(t)); err != nil {
			return apperrors.ErrInvalidPetType
		}
		p.Type = pet.Type(t)
	}
	if val, ok := updates["BloodGroup"]; ok {
		p.BloodGroup = val.(string)
	}
	if val, ok := updates["PetStatus"]; ok {
		ps := val.(string)
		if err := pet.PetStatusValidator(pet.PetStatus(ps)); err != nil {
			return apperrors.ErrPetInvalidStatus
		}
		p.PetStatus = pet.PetStatus(ps)
	}

	// Вычисляем возраст или дату рождения при обновлении
	calculateAgeFields(&p.AgeYears, &p.AgeMonths, &p.BirthDate)

	// Валидируем вложенные структуры
	if health != nil {
		if health.HealthStatus != "" {
			if err := pethealth.HealthStatusValidator(health.HealthStatus); err != nil {
				return apperrors.ErrInvalidHealthStatus
			}
		}
		if health.ReproductiveStatus != "" {
			if err := pethealth.ReproductiveStatusValidator(health.ReproductiveStatus); err != nil {
				return apperrors.ErrInvalidReproductiveStatus
			}
		}
	}
	for _, a := range analyses {
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

	// Сохраняем обновленного питомца
	if _, err := s.petRepo.Update(ctx, p, health, treatments, analyses, bonuses); err != nil {
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
