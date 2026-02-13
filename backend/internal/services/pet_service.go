package services

import (
	"context"
	"fmt"
	"log/slog"
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

	// ApplyValidation применяет валидацию к питомцу, модифицирует объект и сохраняет изменения
	ApplyValidation(ctx context.Context, pet *ent.Pet) ([]validator.FactorCode, []validator.FactorCode, error)

	// BuildFullPhotoURLs преобразует пути к фото в полные публичные URL
	BuildFullPhotoURLs(pet *ent.Pet)
}

// PetServiceImpl реализует PetService
type PetServiceImpl struct {
	petRepo   repositories.PetRepository
	userRepo  repositories.UserRepository
	storage   repositories.FileStorage
	bloodRepo repositories.BloodRequestRepository
	validator validator.DonorValidator
}

// NewPetService создает новый сервис питомцев
func NewPetService(petRepo repositories.PetRepository, userRepo repositories.UserRepository, bloodRepo repositories.BloodRequestRepository, storage repositories.FileStorage, validator validator.DonorValidator) *PetServiceImpl {
	return &PetServiceImpl{
		petRepo:   petRepo,
		userRepo:  userRepo,
		storage:   storage,
		bloodRepo: bloodRepo,
		validator: validator,
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

	// Применяем валидацию, модифицируем объект и сохраняем
	_, _, err = s.ApplyValidation(ctx, newPet)
	if err != nil {
		return nil, err
	}

	// Преобразуем пути к фото в полные URL
	s.BuildFullPhotoURLs(newPet)

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
	s.BuildFullPhotoURLs(p)

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
		s.BuildFullPhotoURLs(pet)
	}

	return pets, nil
}

// UpdatePet обновляет информацию о питомце
func (s *PetServiceImpl) UpdatePet(ctx context.Context, petID string, updates map[string]any, health *ent.PetHealth, treatments *ent.PetTreatment, analyses []*ent.PetAnalysis, bonuses *ent.PetBonus) error {
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
	p, err = s.petRepo.Update(ctx, p, health, treatments, analyses, bonuses)
	if err != nil {
		return apperrors.Internal(err, "failed to update pet")
	}

	// Применяем валидацию, модифицируем объект и сохраняем
	_, _, err = s.ApplyValidation(ctx, p)
	if err != nil {
		return err
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

	// Получаем все заявки на поиск крови, связанные с этим питомцем напрямую через репозиторий
	bloodRequests, err := s.bloodRepo.List(ctx, 0, 0, map[string]any{"pet_id": petID})
	if err != nil {
		return apperrors.Internal(err, "failed to list blood requests for pet")
	}

	// Удаляем каждую связанную заявку
	for _, req := range bloodRequests {
		if err := s.bloodRepo.Delete(ctx, req.ID); err != nil {
			// Логируем ошибку, но продолжаем удаление питомца, чтобы не блокировать операцию
			slog.WarnContext(ctx, "Failed to delete blood request for pet", "blood_request_id", req.ID, "pet_id", petID, "error", err)
		}
	}

	if err := s.petRepo.Delete(ctx, petID); err != nil {
		return apperrors.Internal(err, "failed to delete pet")
	}

	return nil
}

// ValidatePet валидирует питомца и возвращает стоп-факторы и предупреждения
func (s *PetServiceImpl) ApplyValidation(ctx context.Context, p *ent.Pet) ([]validator.FactorCode, []validator.FactorCode, error) {
	stopFactors := s.validator.GetStopFactors(p)
	warnFactors := s.validator.GetWarnFactors(p)

	// Присваиваем результаты валидации объекту питомца (дедуплицируем)
	factorSet := make(map[string]bool)
	var allFactors []string
	for _, f := range stopFactors {
		code := string(f)
		if !factorSet[code] {
			factorSet[code] = true
			allFactors = append(allFactors, code)
		}
	}
	for _, f := range warnFactors {
		code := string(f)
		if !factorSet[code] {
			factorSet[code] = true
			allFactors = append(allFactors, code)
		}
	}
	p.DonorRestrictions = allFactors
	if len(stopFactors) == 0 {
		p.PetStatus = "donor"
	}

	// Сохраняем изменения
	if _, err := s.petRepo.Update(ctx, p, nil, nil, nil, nil); err != nil {
		return nil, nil, apperrors.Internal(err, "failed to save validation results")
	}

	return stopFactors, warnFactors, nil
}

//===================HELPERS===============================================

// BuildFullPhotoURLs преобразует пути к фото в полные публичные URL
func (s *PetServiceImpl) BuildFullPhotoURLs(pet *ent.Pet) {
	if len(pet.PhotoUrls) == 0 {
		return
	}
	for i, path := range pet.PhotoUrls {
		if path == "" {
			continue
		}
		url := s.storage.GetPublicURLFromPath(path)
		timestamp := pet.UpdatedAt.Unix()
		pet.PhotoUrls[i] = fmt.Sprintf("%s?t=%d", url, timestamp)
	}
}
