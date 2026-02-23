package services

import (
	"context"
	"log/slog"

	"github.com/artesipov-alt/odnoi-krovi-app/ent"
	"github.com/artesipov-alt/odnoi-krovi-app/ent/breed"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/validator"
)

// PetRepository определяет интерфейс для операций с данными питомцев
type PetRepository interface {
	// Create создает нового питомца в базе данных
	Create(ctx context.Context, petDomain *domain.Pet) (*domain.Pet, error)

	// GetPetQuery возвращает query для eager loading
	GetPet(ctx context.Context, id string, opts PetPreloadOptions) (*domain.Pet, error)

	// GetPetsQueryByUser возвращает query для eager loading питомцев пользователя
	GetPetsByUser(ctx context.Context, userID string, opts PetPreloadOptions) ([]*domain.Pet, error)

	// Update обновляет питомца и его связанные сущности в одной транзакции
	// Все параметры (кроме id и pet) могут быть nil - тогда соответствующие данные не обновляются
	Update(ctx context.Context, id string, petDomain *domain.Pet) (*domain.Pet, error)

	// Delete удаляет питомца по его ID
	Delete(ctx context.Context, id string) error

	// ExistsByID проверяет, существует ли питомец с заданным ID
	ExistsByID(ctx context.Context, id string) (bool, error)

	// // UpdateStatus обновляет статус питомца по его ID
	// UpdateStatus(ctx context.Context, id string, status string) error

	// AddPhotoURLs добавляет новые пути к фотографиям питомца
	AddPhotoURLs(ctx context.Context, id string, paths []string) error
}

// BreedRepository определяет интерфейс для операций с данными пород
type BreedRepository interface {
	// GetAll возвращает все породы из базы данных
	GetAll(ctx context.Context) ([]*ent.Breed, error)

	// GetByID получает породу по её ID
	GetByID(ctx context.Context, id string) (*ent.Breed, error)

	// GetByPetType получает породы по типу животного
	GetByPetType(ctx context.Context, petType breed.Type) ([]*ent.Breed, error)

	// Create создает новую породу в базе данных
	Create(ctx context.Context, b *ent.Breed) (*ent.Breed, error)

	// Update обновляет существующую породу в базе данных
	Update(ctx context.Context, b *ent.Breed) (*ent.Breed, error)

	// Delete удаляет породу по её ID
	Delete(ctx context.Context, id string) error

	// ExistsByName проверяет, существует ли порода с заданным названием
	ExistsByName(ctx context.Context, name string) (bool, error)
}

// PetPreloadOptions определяет опции для preload связанных данных питомца
type PetPreloadOptions struct {
	WithHealth     bool
	WithTreatments bool
	WithAnalyses   bool
	WithBonuses    bool
	WithAll        bool
}

// PetService реализует PetService
type PetService struct {
	petRepo      PetRepository
	userRepo     UserRepository
	storage      FileStorage
	bloodReqRepo BloodRequestRepository
	blooInfoRepo BloodInfoRepository
}

// NewPetService создает новый сервис питомцев
func NewPetService(petRepo PetRepository, userRepo UserRepository, bloodReqRepo BloodRequestRepository, bloodInfoRepo BloodInfoRepository, storage FileStorage, validator validator.DonorValidator) *PetService {
	return &PetService{
		petRepo:      petRepo,
		userRepo:     userRepo,
		storage:      storage,
		bloodReqRepo: bloodReqRepo,
		blooInfoRepo: bloodInfoRepo,
	}
}

// CreatePet создает нового питомца для пользователя
func (s *PetService) CreatePet(ctx context.Context, userID string, pet *domain.Pet) (*domain.Pet, error) {
	// Проверяем, существует ли пользователь
	exists, err := s.userRepo.ExistsByID(ctx, userID)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to check user existence")
	}
	if !exists {
		return nil, apperrors.ErrUserNotFound
	}

	// Проверяем, существует ли группа крови, если указана
	if pet.BloodGroupRefID != nil {
		bg, err := s.blooInfoRepo.FindByName(ctx, *pet.BloodGroupRefID)
		if err != nil {
			return nil, apperrors.Internal(err, "failed to check blood group existence")
		}
		pet.BloodGroupRefID = &bg.ID
	}

	// Set the owner ID for the pet
	pet.OwnerID = userID

	newPet, err := s.petRepo.Create(ctx, pet)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to create pet")
	}

	return newPet, nil
}

// GetPet получает питомца с preload связанных данных
func (s *PetService) GetPet(ctx context.Context, petID string, opts PetPreloadOptions) (*domain.Pet, error) {
	pet, err := s.petRepo.GetPet(ctx, petID, opts)
	if err != nil {
		return nil, err
	}

	pet.PhotoURLs = s.BuildFullPhotoURLs(pet.PhotoURLs)

	return pet, nil
}

// GetUserPets получает всех питомцев пользователя с preload связей
func (s *PetService) GetUserPets(ctx context.Context, userID string, opts PetPreloadOptions) ([]*domain.Pet, error) {
	_, err := s.userRepo.GetByID(ctx, userID, UserPreloadOptions{})
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, apperrors.ErrUserNotFound
		}
		return nil, apperrors.Internal(err, "failed to get user")
	}

	pets, err := s.petRepo.GetPetsByUser(ctx, userID, opts)
	if err != nil {
		return nil, err
	}

	for i, pet := range pets {
		pets[i].PhotoURLs = s.BuildFullPhotoURLs(pet.PhotoURLs)
	}

	return pets, nil
}

// Update обновляет питомца и его связанные сущности
func (s *PetService) Update(ctx context.Context, id string, petInput *domain.Pet) (*domain.Pet, error) {
	exists, err := s.petRepo.ExistsByID(ctx, id)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to check user existence")
	}
	if !exists {
		return nil, apperrors.ErrPetNotFound
	}

	updatedPet, err := s.petRepo.Update(ctx, id, petInput)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to create pet")
	}

	return updatedPet, nil

}

// DeletePet удаляет питомца по ID
func (s *PetService) DeletePet(ctx context.Context, petID string) error {
	exists, err := s.petRepo.ExistsByID(ctx, petID)
	if err != nil {
		return apperrors.Internal(err, "failed to check pet existence")
	}
	if !exists {
		return apperrors.ErrPetNotFound
	}

	// Получаем все заявки на поиск крови, связанные с этим питомцем
	bloodRequests, err := s.bloodReqRepo.List(ctx, 0, 0, map[string]any{"pet_id": petID})
	if err != nil {
		return apperrors.Internal(err, "failed to list blood requests for pet")
	}

	// Удаляем каждую связанную заявку
	for _, req := range bloodRequests {
		if err := s.bloodReqRepo.Delete(ctx, req.ID); err != nil {
			slog.WarnContext(ctx, "Failed to delete blood request for pet", "blood_request_id", req.ID, "pet_id", petID, "error", err)
		}
	}

	if err := s.petRepo.Delete(ctx, petID); err != nil {
		return apperrors.Internal(err, "failed to delete pet")
	}

	return nil
}

// TODO Продумать как сделать валидацию более правильно
// ApplyValidation применяет валидацию к питомцу, модифицирует объект и сохраняет изменения
// func (s *PetService) ApplyValidation(ctx context.Context, petID string) ([]validator.FactorCode, []validator.FactorCode, error) {
// 	// Получаем питомца для валидации
// 	p, err := s.GetPet(ctx, petID, PetPreloadOptions{WithAll: true})
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	stopFactors := s.validator.GetStopFactors(p)
// 	warnFactors := s.validator.GetWarnFactors(p)

// 	// Дедуплицируем факторы
// 	factorSet := make(map[string]bool)
// 	var allFactors []string
// 	for _, f := range stopFactors {
// 		code := string(f)
// 		if !factorSet[code] {
// 			factorSet[code] = true
// 			allFactors = append(allFactors, code)
// 		}
// 	}
// 	for _, f := range warnFactors {
// 		code := string(f)
// 		if !factorSet[code] {
// 			factorSet[code] = true
// 			allFactors = append(allFactors, code)
// 		}
// 	}

// 	// Определяем новый статус
// 	// newStatus := ""
// 	// if len(stopFactors) == 0 {
// 	// 	newStatus = "donor"
// 	// }

// 	// Создаем UpdatePetInput для сохранения результатов валидации
// 	updateInput := &ent.UpdatePetInput{
// 		DonorRestrictions: allFactors,
// 	}
// 	// if newStatus != "" {
// 	// 	petStatus := pet.PetStatus(newStatus)
// 	// 	updateInput.PetStatus = &petStatus
// 	// }

// 	// Обновляем только поля валидации
// 	if _, err := s.petRepo.Update(ctx, petID, updateInput, nil, nil, nil, nil); err != nil {
// 		return nil, nil, apperrors.Internal(err, "failed to save validation results")
// 	}

// 	return stopFactors, warnFactors, nil
// }

//===================HELPERS===============================================

// BuildFullPhotoURLs преобразует пути к фото в полные публичные URL
func (s *PetService) BuildFullPhotoURLs(paths []string) []string {
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
