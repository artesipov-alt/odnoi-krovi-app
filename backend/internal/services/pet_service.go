package services

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/ent"
	"github.com/artesipov-alt/odnoi-krovi-app/ent/pet"
	"github.com/artesipov-alt/odnoi-krovi-app/ent/petanalysis"
	"github.com/artesipov-alt/odnoi-krovi-app/ent/pethealth"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	repositories "github.com/artesipov-alt/odnoi-krovi-app/internal/repositories"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/validator"
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

	// buildFullPhotoURLs преобразует пути к фото в полные публичные URL
	buildFullPhotoURLs(paths []string) []string
}

// PetServiceImpl реализует PetService
type PetServiceImpl struct {
	petRepo            repositories.PetRepository
	userRepo           repositories.UserRepository
	storage            repositories.FileStorage
	bloodSearchService BloodSearchService
	validator          validator.PetValidator
}

// NewPetService создает новый сервис питомцев
func NewPetService(petRepo repositories.PetRepository, userRepo repositories.UserRepository, storage repositories.FileStorage, bloodSearchService BloodSearchService) *PetServiceImpl {
	return &PetServiceImpl{
		petRepo:            petRepo,
		userRepo:           userRepo,
		storage:            storage,
		bloodSearchService: bloodSearchService,
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

	if string(petData.ReproductiveStatus) != "" {
		if err := pet.ReproductiveStatusValidator(petData.ReproductiveStatus); err != nil {
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
	if val, ok := updates["ReproductiveStatus"]; ok {
		rs := val.(string)
		if err := pet.ReproductiveStatusValidator(pet.ReproductiveStatus(rs)); err != nil {
			return apperrors.Validation("неверный репродуктивный статус", nil).WithInternal(err)
		}
		p.ReproductiveStatus = pet.ReproductiveStatus(rs)
	}
	if val, ok := updates["PetStatus"]; ok {
		ps := val.(string)
		if err := pet.PetStatusValidator(pet.PetStatus(ps)); err != nil {
			return apperrors.Validation("неверный статус питомца", nil).WithInternal(err)
		}
		p.PetStatus = pet.PetStatus(ps)
	}

	// Валидируем вложенные структуры
	if health != nil {
		if health.HealthStatus != "" {
			if err := pethealth.HealthStatusValidator(health.HealthStatus); err != nil {
				return apperrors.Validation("неверный статус здоровья", nil).WithInternal(err)
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

	// Получаем все заявки на поиск крови, связанные с этим питомцем
	bloodRequests, err := s.bloodSearchService.ListRequests(ctx, 0, 0, map[string]any{"pet_id": petID})
	if err != nil {
		return apperrors.Internal(err, "failed to list blood requests for pet")
	}

	// Удаляем каждую связанную заявку
	for _, req := range bloodRequests {
		if err := s.bloodSearchService.DeleteRequest(ctx, req.ID); err != nil {
			// Логируем ошибку, но продолжаем удаление питомца, чтобы не блокировать операцию
			// В реальном приложении здесь может быть более сложная логика обработки ошибок
			// например, попытка повтора или уведомление администратора.
			// Для данного кейса, просто логируем и продолжаем.
			slog.WarnContext(ctx, "Failed to delete blood request for pet", "blood_request_id", req.ID, "pet_id", petID, "error", err)
		}
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

//===================HELPERS===============================================

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
