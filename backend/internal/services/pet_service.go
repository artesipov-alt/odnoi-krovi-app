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

// PetService определяет интерфейс для бизнес-логики питомцев
type PetService interface {
	// CreatePet создает нового питомца для пользователя
	CreatePet(ctx context.Context, userID string, pet *ent.Pet) (*ent.Pet, error)

	// GetPetByID получает питомца по ID с preload связей
	GetPetByID(ctx context.Context, petID string, preloads ...string) (*ent.Pet, error)

	// GetUserPets получает всех питомцев пользователя с preload связей
	GetUserPets(ctx context.Context, userID string, preloads ...string) ([]*ent.Pet, error)

	// UpdatePet обновляет информацию о питомце
	UpdatePet(ctx context.Context, petID string, updates map[string]any, health *ent.PetHealth, treatments *ent.PetTreatment, analyses []*ent.PetAnalysis, bonuses *ent.PetBonus) error

	// DeletePet удаляет питомца по ID
	DeletePet(ctx context.Context, petID string) error

	// UpdatePetAvatar обновляет аватар питомца и делает его публичным в хранилище
	UpdatePetAvatar(ctx context.Context, avatarPath string) (string, error)

	// ConfirmPhotos подтверждает загрузку фото для питомца и обновляет PhotoUrls
	ConfirmPhotos(ctx context.Context, petID string, paths []string) error

	// calculateAgeFields вычисляет возраст из даты рождения или дату из возраста
	calculateAgeFields(ageYears, ageMonths *int, birthDate **time.Time)

	// buildFullPhotoURLs преобразует пути к фото в полные публичные URL
	buildFullPhotoURLs(paths []string) []string
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

// CreatePet создает нового питомца для пользователя
func (s *PetServiceImpl) CreatePet(ctx context.Context, userID string, petData *ent.Pet) (*ent.Pet, error) {
	// Проверяем, существует ли пользователь
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, apperrors.ErrUserNotFound
		}
		return nil, apperrors.Internal(err, "failed to get user")
	}

	// Устанавливаем UserID
	petData.UserID = userID

	// Вычисляем возраст или дату рождения
	s.calculateAgeFields(&petData.AgeYears, &petData.AgeMonths, &petData.BirthDate)

	// Валидируем тип животного
	if err := pet.TypeValidator(petData.Type); err != nil {
		return nil, apperrors.Validation("неверный тип питомца", nil).WithInternal(err)
	}

	// Валидируем статус питомца
	if err := pet.PetStatusValidator(petData.PetStatus); err != nil {
		return nil, apperrors.Validation("неверный статус питомца", nil).WithInternal(err)
	}

	// Валидируем пол животного
	if string(petData.Gender) != "" {
		if err := pet.GenderValidator(petData.Gender); err != nil {
			return nil, apperrors.Validation("неверный пол животного", nil).WithInternal(err)
		}
	}

	// Валидируем условия проживания
	if string(petData.LivingCondition) != "" {
		if err := pet.LivingConditionValidator(petData.LivingCondition); err != nil {
			return nil, apperrors.Validation("неверные условия проживания", nil).WithInternal(err)
		}
	}

	// Валидируем вложенные структуры, если они есть
	if petData.Edges.Health != nil {
		if string(petData.Edges.Health.HealthStatus) != "" {
			if err := pethealth.HealthStatusValidator(petData.Edges.Health.HealthStatus); err != nil {
				return nil, apperrors.Validation("неверный статус здоровья", nil).WithInternal(err)
			}
		}
		if string(petData.Edges.Health.ReproductiveStatus) != "" {
			if err := pethealth.ReproductiveStatusValidator(petData.Edges.Health.ReproductiveStatus); err != nil {
				return nil, apperrors.Validation("неверный репродуктивный статус", nil).WithInternal(err)
			}
		}
	}

	if petData.Edges.Analyses != nil {
		for _, a := range petData.Edges.Analyses {
			if string(a.AnalysisName) != "" {
				if err := petanalysis.AnalysisNameValidator(a.AnalysisName); err != nil {
					return nil, apperrors.Validation("неверный тип лейкемии", nil).WithInternal(err)
				}
			}
		}
	}

	newPet, err := s.petRepo.Create(ctx, petData, petData.Edges.Health, petData.Edges.Treatments, petData.Edges.Analyses, petData.Edges.Bonuses)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to create pet")
	}

	// Преобразуем пути к фото в полные URL
	newPet.PhotoUrls = s.buildFullPhotoURLs(newPet.PhotoUrls)

	return newPet, nil
}

// GetPetByID получает питомца по ID с preload связей
func (s *PetServiceImpl) GetPetByID(ctx context.Context, petID string, preloads ...string) (*ent.Pet, error) {
	if petID == "" {
		return nil, apperrors.BadRequest("неверный ID питомца")
	}

	p, err := s.petRepo.GetByID(ctx, petID, preloads...)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, apperrors.ErrPetNotFound
		}
		return nil, apperrors.Internal(err, "failed to get pet")
	}

	// Преобразуем пути к фото в полные URL
	p.PhotoUrls = s.buildFullPhotoURLs(p.PhotoUrls)

	return p, nil
}

// GetUserPets получает всех питомцев пользователя с preload связей
func (s *PetServiceImpl) GetUserPets(ctx context.Context, userID string, preloads ...string) ([]*ent.Pet, error) {
	// Проверяем, существует ли пользователь
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, apperrors.ErrUserNotFound
		}
		return nil, apperrors.Internal(err, "failed to get user")
	}

	pets, err := s.petRepo.GetByUserID(ctx, userID, preloads...)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to get user pets")
	}

	// Преобразуем пути к фото в полные URL
	for _, pet := range pets {
		pet.PhotoUrls = s.buildFullPhotoURLs(pet.PhotoUrls)
	}

	return pets, nil
}

// UpdatePet обновляет информацию о питомце
func (s *PetServiceImpl) UpdatePet(ctx context.Context, petID string, updates map[string]interface{}, health *ent.PetHealth, treatments *ent.PetTreatment, analyses []*ent.PetAnalysis, bonuses *ent.PetBonus) error {
	// Получаем существующего питомца
	p, err := s.petRepo.GetByID(ctx, petID)
	if err != nil {
		if ent.IsNotFound(err) {
			return apperrors.ErrPetNotFound
		}
		return apperrors.Internal(err, "failed to get pet")
	}

	// Применяем обновления и валидируем
	if val, ok := updates["Name"]; ok {
		p.Name = val.(string)
	}
	if val, ok := updates["ChipNumber"]; ok {
		p.ChipNumber = val.(string)
	}
	if val, ok := updates["PhotoUrls"]; ok {
		p.PhotoUrls = val.([]string)
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
			return apperrors.Validation("неверные условия проживания", nil).WithInternal(err)
		}
		p.LivingCondition = pet.LivingCondition(lc)
	}
	if val, ok := updates["Gender"]; ok {
		g := val.(string)
		if err := pet.GenderValidator(pet.Gender(g)); err != nil {
			return apperrors.Validation("неверный пол животного", nil).WithInternal(err)
		}
		p.Gender = pet.Gender(g)
	}
	if val, ok := updates["Type"]; ok {
		t := val.(string)
		if err := pet.TypeValidator(pet.Type(t)); err != nil {
			return apperrors.Validation("неверный тип питомца", nil).WithInternal(err)
		}
		p.Type = pet.Type(t)
	}
	if val, ok := updates["BloodGroup"]; ok {
		p.BloodGroup = val.(string)
	}
	if val, ok := updates["PetStatus"]; ok {
		ps := val.(string)
		if err := pet.PetStatusValidator(pet.PetStatus(ps)); err != nil {
			return apperrors.Validation("неверный статус питомца", nil).WithInternal(err)
		}
		p.PetStatus = pet.PetStatus(ps)
	}

	// Вычисляем возраст или дату рождения при обновлении
	s.calculateAgeFields(&p.AgeYears, &p.AgeMonths, &p.BirthDate)

	// Валидируем вложенные структуры
	if health != nil {
		if health.HealthStatus != "" {
			if err := pethealth.HealthStatusValidator(health.HealthStatus); err != nil {
				return apperrors.Validation("неверный статус здоровья", nil).WithInternal(err)
			}
		}
		if health.ReproductiveStatus != "" {
			if err := pethealth.ReproductiveStatusValidator(health.ReproductiveStatus); err != nil {
				return apperrors.Validation("неверный репродуктивный статус", nil).WithInternal(err)
			}
		}
	}
	for _, a := range analyses {
		if string(a.AnalysisName) != "" {
			if err := petanalysis.AnalysisNameValidator(a.AnalysisName); err != nil {
				return apperrors.Validation("неверный тип лейкемии", nil).WithInternal(err)
			}
		}
	}

	// Сохраняем обновленного питомца
	if _, err := s.petRepo.Update(ctx, p, health, treatments, analyses, bonuses); err != nil {
		return apperrors.Internal(err, "failed to update pet")
	}

	return nil
}

// DeletePet удаляет питомца по ID
func (s *PetServiceImpl) DeletePet(ctx context.Context, petID string) error {
	exists, err := s.petRepo.ExistsByID(ctx, petID)
	if err != nil {
		return apperrors.Internal(err, "failed to check pet existence")
	}
	if !exists {
		return apperrors.ErrPetNotFound
	}

	if err := s.petRepo.Delete(ctx, petID); err != nil {
		return apperrors.Internal(err, "failed to delete pet")
	}

	return nil
}

// UpdatePetAvatar обновляет аватар питомца и делает его публичным в хранилище
func (s *PetServiceImpl) UpdatePetAvatar(ctx context.Context, avatarPath string) (string, error) {
	decodedPath := strings.ReplaceAll(avatarPath, "%2F", "/")
	parts := strings.Split(decodedPath, "/")
	if len(parts) < 2 {
		return "", apperrors.BadRequest("неверный формат пути аватара")
	}
	petID := parts[1]

	p, err := s.petRepo.GetByID(ctx, petID)
	if err != nil {
		if ent.IsNotFound(err) {
			return "", apperrors.ErrPetNotFound
		}
		return "", apperrors.Internal(err, "failed to get pet")
	}

	exists, err := s.storage.CheckObjectExists(ctx, decodedPath)
	if err != nil {
		return "", apperrors.Internal(err, "failed to check file existence")
	}
	if !exists {
		return "", apperrors.NotFound("файл не найден")
	}

	err = s.storage.SetObjectPublicACL(ctx, decodedPath)
	if err != nil {
		return "", apperrors.Internal(err, "failed to set public ACL")
	}

	publicURL := s.storage.GetAvatarPublicURL(petID)
	if publicURL == "" {
		return "", apperrors.Internal(errors.New("failed to generate public URL"), "failed to get public URL")
	}

	p.PhotoUrls = []string{decodedPath}
	if _, err := s.petRepo.Update(ctx, p, nil, nil, nil, nil); err != nil {
		return "", apperrors.Internal(err, "failed to update pet avatar")
	}

	return publicURL, nil
}

// ConfirmPhotos подтверждает загрузку фото для питомца и обновляет PhotoUrls
func (s *PetServiceImpl) ConfirmPhotos(ctx context.Context, petID string, paths []string) error {
	// Получить питомца
	p, err := s.petRepo.GetByID(ctx, petID)
	if err != nil {
		if ent.IsNotFound(err) {
			return apperrors.ErrPetNotFound
		}
		return apperrors.Internal(err, "failed to get pet")
	}

	// Обновить PhotoUrls: добавить новые пути к существующим
	p.PhotoUrls = append(p.PhotoUrls, paths...)

	// Сохранить обновленного питомца
	if _, err := s.petRepo.Update(ctx, p, nil, nil, nil, nil); err != nil {
		return apperrors.Internal(err, "failed to update pet photos")
	}

	return nil
}

//===================HELPERS===============================================

// calculateAgeFields вычисляет возраст из даты рождения или дату из возраста
func (s *PetServiceImpl) calculateAgeFields(ageYears, ageMonths *int, birthDate **time.Time) {
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

// buildFullPhotoURLs преобразует пути к фото в полные публичные URL
func (s *PetServiceImpl) buildFullPhotoURLs(paths []string) []string {
	if len(paths) == 0 {
		return []string{}
	}
	result := make([]string, len(paths))
	for i, path := range paths {
		if path == "" {
			result[i] = ""
		} else {
			result[i] = s.storage.GetPublicURLFromPath(path)
		}
	}
	return result
}
