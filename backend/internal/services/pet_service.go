package services

import (
	"context"
	"fmt"
	"log/slog"

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
	CreatePet(ctx context.Context, userID string, input *ent.CreatePetInput, healthInput *ent.CreatePetHealthInput, treatmentsInput *ent.CreatePetTreatmentInput, analysesInput []*ent.CreatePetAnalysisInput, bonusesInput *ent.CreatePetBonusInput) (*ent.Pet, error)

	// GetPetByID получает питомца по ID с preload связей
	GetPetByID(ctx context.Context, petID string, preloads ...string) (*ent.Pet, error)

	// GetPetQuery возвращает query для eager loading
	GetPetQuery(ctx context.Context, petID string) *ent.PetQuery

	// GetUserPets получает всех питомцев пользователя с preload связей
	GetUserPets(ctx context.Context, userID string, preloads ...string) ([]*ent.Pet, error)

	// GetPetsQueryByUser возвращает query для eager loading питомцев пользователя
	GetPetsQueryByUser(ctx context.Context, userID string) *ent.PetQuery

	// Update обновляет питомца и его связанные сущности (здоровье, лечения, анализы, бонусы)
	// Все параметры могут быть nil - тогда соответствующие данные не обновляются
	Update(ctx context.Context, id string, petInput *ent.UpdatePetInput, healthInput *ent.UpdatePetHealthInput, treatmentsInput *ent.UpdatePetTreatmentInput, analysesInput []*ent.UpdatePetAnalysisInput, bonusesInput *ent.UpdatePetBonusInput) (*ent.Pet, error)

	// DeletePet удаляет питомца по ID
	DeletePet(ctx context.Context, petID string) error

	// ApplyValidation применяет валидацию к питомцу, модифицирует объект и сохраняет изменения
	ApplyValidation(ctx context.Context, petID string) ([]validator.FactorCode, []validator.FactorCode, error)

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
func (s *PetServiceImpl) CreatePet(ctx context.Context, userID string, input *ent.CreatePetInput, healthInput *ent.CreatePetHealthInput, treatmentsInput *ent.CreatePetTreatmentInput, analysesInput []*ent.CreatePetAnalysisInput, bonusesInput *ent.CreatePetBonusInput) (*ent.Pet, error) {
	// Проверяем, существует ли пользователь
	exists, err := s.userRepo.ExistsByID(ctx, userID)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to check user existence")
	}
	if !exists {
		return nil, apperrors.ErrUserNotFound
	}

	// Нормализуем опциональные поля: конвертируем пустые строки в nil
	if input.Gender != nil && string(*input.Gender) == "" {
		input.Gender = nil
	}
	if input.LivingCondition != nil && string(*input.LivingCondition) == "" {
		input.LivingCondition = nil
	}
	if input.ReproductiveStatus != nil && string(*input.ReproductiveStatus) == "" {
		input.ReproductiveStatus = nil
	}

	// Валидируем тип животного
	if err := pet.TypeValidator(input.Type); err != nil {
		return nil, apperrors.Validation("неверный тип питомца", nil).WithInternal(err)
	}

	// Валидируем статус питомца
	if err := pet.PetStatusValidator(input.PetStatus); err != nil {
		return nil, apperrors.Validation("неверный статус питомца", nil).WithInternal(err)
	}

	// Валидируем пол животного
	if input.Gender != nil && string(*input.Gender) != "" {
		if err := pet.GenderValidator(*input.Gender); err != nil {
			return nil, apperrors.Validation("неверный пол животного", nil).WithInternal(err)
		}
	}

	// Валидируем условия проживания
	if input.LivingCondition != nil && string(*input.LivingCondition) != "" {
		if err := pet.LivingConditionValidator(*input.LivingCondition); err != nil {
			return nil, apperrors.Validation("неверные условия проживания", nil).WithInternal(err)
		}
	}

	if input.ReproductiveStatus != nil && string(*input.ReproductiveStatus) != "" {
		if err := pet.ReproductiveStatusValidator(*input.ReproductiveStatus); err != nil {
			return nil, apperrors.Validation("неверные условия проживания", nil).WithInternal(err)
		}
	}

	// Валидируем вложенные структуры, если они есть
	if healthInput != nil {
		// Нормализуем HealthStatus: конвертируем пустую строку в nil
		if healthInput.HealthStatus != nil && string(*healthInput.HealthStatus) == "" {
			healthInput.HealthStatus = nil
		}

		if healthInput.HealthStatus != nil && string(*healthInput.HealthStatus) != "" {
			if err := pethealth.HealthStatusValidator(*healthInput.HealthStatus); err != nil {
				return nil, apperrors.Validation("неверный статус здоровья", nil).WithInternal(err)
			}
		}
	}

	for _, a := range analysesInput {
		// Нормализуем AnalysisName и AnalysisType: конвертируем пустые строки в nil
		if a.AnalysisName != nil && string(*a.AnalysisName) == "" {
			a.AnalysisName = nil
		}
		if a.AnalysisType != nil && string(*a.AnalysisType) == "" {
			a.AnalysisType = nil
		}

		if a.AnalysisName != nil && *a.AnalysisName != petanalysis.AnalysisName("") {
			if err := petanalysis.AnalysisNameValidator(*a.AnalysisName); err != nil {
				return nil, apperrors.Validation("неверный тип анализа", nil).WithInternal(err)
			}
		}
	}

	// Set the owner ID for the pet
	input.OwnerID = &userID

	newPet, err := s.petRepo.Create(ctx, input, healthInput, treatmentsInput, analysesInput, bonusesInput)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to create pet")
	}

	// Применяем валидацию
	_, _, err = s.ApplyValidation(ctx, newPet.ID)
	if err != nil {
		return nil, err
	}

	// Получаем обновленного питомца
	updatedPet, err := s.petRepo.GetPetQuery(ctx, newPet.ID).Only(ctx)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to get updated pet")
	}

	// Преобразуем пути к фото в полные URL
	s.BuildFullPhotoURLs(updatedPet)

	return updatedPet, nil
}

// GetPetByID получает питомца по ID с preload связей
func (s *PetServiceImpl) GetPetByID(ctx context.Context, petID string, preloads ...string) (*ent.Pet, error) {
	if petID == "" {
		return nil, apperrors.BadRequest("неверный ID питомца")
	}

	p, err := s.petRepo.GetPetQuery(ctx, petID).Only(ctx)
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

// GetPetQuery возвращает query для eager loading
func (s *PetServiceImpl) GetPetQuery(ctx context.Context, petID string) *ent.PetQuery {
	return s.petRepo.GetPetQuery(ctx, petID)
}

// GetPetsQueryByUser возвращает query для eager loading питомцев пользователя
func (s *PetServiceImpl) GetPetsQueryByUser(ctx context.Context, userID string) *ent.PetQuery {
	return s.petRepo.GetPetsQueryByUser(ctx, userID)
}

// GetUserPets получает всех питомцев пользователя с preload связей
func (s *PetServiceImpl) GetUserPets(ctx context.Context, userID string, preloads ...string) ([]*ent.Pet, error) {
	// Проверяем, существует ли пользователь
	query := s.userRepo.GetQueryByID(ctx, userID)
	_, err := query.Only(ctx)
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
	for _, p := range pets {
		s.BuildFullPhotoURLs(p)
	}

	return pets, nil
}

// Update обновляет питомца и его связанные сущности
func (s *PetServiceImpl) Update(ctx context.Context, id string, petInput *ent.UpdatePetInput, healthInput *ent.UpdatePetHealthInput, treatmentsInput *ent.UpdatePetTreatmentInput, analysesInput []*ent.UpdatePetAnalysisInput, bonusesInput *ent.UpdatePetBonusInput) (*ent.Pet, error) {
	// Валидируем основные данные питомца
	if petInput != nil {
		if petInput.Type != nil && *petInput.Type != "" {
			if err := pet.TypeValidator(*petInput.Type); err != nil {
				return nil, apperrors.Validation("неверный тип питомца", nil).WithInternal(err)
			}
		}
		if petInput.PetStatus != nil && *petInput.PetStatus != "" {
			if err := pet.PetStatusValidator(*petInput.PetStatus); err != nil {
				return nil, apperrors.Validation("неверный статус питомца", nil).WithInternal(err)
			}
		}
		if petInput.Gender != nil && string(*petInput.Gender) != "" {
			if err := pet.GenderValidator(*petInput.Gender); err != nil {
				return nil, apperrors.Validation("неверный пол животного", nil).WithInternal(err)
			}
		}
		if petInput.LivingCondition != nil && string(*petInput.LivingCondition) != "" {
			if err := pet.LivingConditionValidator(*petInput.LivingCondition); err != nil {
				return nil, apperrors.Validation("неверные условия проживания", nil).WithInternal(err)
			}
		}
		if petInput.ReproductiveStatus != nil && string(*petInput.ReproductiveStatus) != "" {
			if err := pet.ReproductiveStatusValidator(*petInput.ReproductiveStatus); err != nil {
				return nil, apperrors.Validation("неверный репродуктивный статус", nil).WithInternal(err)
			}
		}

	}

	// Валидируем связанные данные
	if healthInput != nil {
		if healthInput.HealthStatus != nil && string(*healthInput.HealthStatus) != "" {
			if err := pethealth.HealthStatusValidator(*healthInput.HealthStatus); err != nil {
				return nil, apperrors.Validation("неверный статус здоровья", nil).WithInternal(err)
			}
		}
	}

	for _, a := range analysesInput {
		if a.AnalysisName != nil && *a.AnalysisName != petanalysis.AnalysisName("") {
			if err := petanalysis.AnalysisNameValidator(*a.AnalysisName); err != nil {
				return nil, apperrors.Validation("неверный тип анализа", nil).WithInternal(err)
			}
		}
	}

	// Выполняем обновление через репозиторий
	_, err := s.petRepo.Update(ctx, id, petInput, healthInput, treatmentsInput, analysesInput, bonusesInput)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, apperrors.ErrPetNotFound
		}
		return nil, apperrors.Internal(err, "failed to update pet")
	}

	// Применяем валидацию
	_, _, err = s.ApplyValidation(ctx, id)
	if err != nil {
		return nil, err
	}

	// Получаем обновленного питомца
	updatedPet, err := s.petRepo.GetPetQuery(ctx, id).Only(ctx)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to get updated pet")
	}

	// Преобразуем пути к фото в полные URL
	s.BuildFullPhotoURLs(updatedPet)

	return updatedPet, nil
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
	bloodRequests, err := s.bloodRepo.List(ctx, 0, 0, map[string]any{"pet_id": petID})
	if err != nil {
		return apperrors.Internal(err, "failed to list blood requests for pet")
	}

	// Удаляем каждую связанную заявку
	for _, req := range bloodRequests {
		if err := s.bloodRepo.Delete(ctx, req.ID); err != nil {
			slog.WarnContext(ctx, "Failed to delete blood request for pet", "blood_request_id", req.ID, "pet_id", petID, "error", err)
		}
	}

	if err := s.petRepo.Delete(ctx, petID); err != nil {
		return apperrors.Internal(err, "failed to delete pet")
	}

	return nil
}

// ApplyValidation применяет валидацию к питомцу, модифицирует объект и сохраняет изменения
func (s *PetServiceImpl) ApplyValidation(ctx context.Context, petID string) ([]validator.FactorCode, []validator.FactorCode, error) {
	// Получаем питомца для валидации
	p, err := s.petRepo.GetPetQuery(ctx, petID).
		WithHealth().
		WithTreatments().
		WithAnalyses().
		WithBonuses().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil, apperrors.ErrPetNotFound
		}
		return nil, nil, apperrors.Internal(err, "failed to get pet for validation")
	}

	stopFactors := s.validator.GetStopFactors(p)
	warnFactors := s.validator.GetWarnFactors(p)

	// Дедуплицируем факторы
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

	// Определяем новый статус
	newStatus := ""
	if len(stopFactors) == 0 {
		newStatus = "donor"
	}

	// Создаем UpdatePetInput для сохранения результатов валидации
	updateInput := &ent.UpdatePetInput{
		DonorRestrictions: allFactors,
	}
	if newStatus != "" {
		petStatus := pet.PetStatus(newStatus)
		updateInput.PetStatus = &petStatus
	}

	// Обновляем только поля валидации
	if _, err := s.petRepo.Update(ctx, petID, updateInput, nil, nil, nil, nil); err != nil {
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
