package pet

// import (
// 	"context"
// 	"errors"
// 	"log/slog"
// 	"strconv"
// 	"time"

// 	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
// 	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"
// 	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/filestorage"
// 	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet/model"
// 	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user"
// 	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent"
// )

// // PetService реализует PetService
// type PetService struct {
// 	petRepo      Repository
// 	userRepo     user.Repository
// 	bloodReqRepo bloodsearch.BloodRequestRepository
// 	storage      filestorage.Repository
// }

// // NewPetService создает новый сервис питомцев
// func NewPetService(petRepo Repository, userRepo user.Repository, bloodReqRepo bloodsearch.BloodRequestRepository, storage filestorage.Repository) *PetService {
// 	return &PetService{
// 		petRepo:      petRepo,
// 		userRepo:     userRepo,
// 		storage:      storage,
// 		bloodReqRepo: bloodReqRepo,
// 	}
// }

// // CreatePet создает нового питомца для пользователя
// func (s *PetService) CreatePet(ctx context.Context, userID string, pet *model.Pet) (*model.Pet, error) {
// 	// Проверяем, существует ли пользователь
// 	exists, err := s.userRepo.ExistsByID(ctx, userID)
// 	if err != nil {
// 		return nil, apperrors.Internal(err, "failed to check user existence")
// 	}
// 	if !exists {
// 		return nil, apperrors.ErrUserNotFound
// 	}

// 	// Set the owner ID for the pet
// 	pet.OwnerID = userID

// 	// Calculate stop and warn factors and set DonorRestrictions
// 	stopFactors := pet.GetStopFactors(time.Now())
// 	pet.StopFactors = make([]string, len(stopFactors))
// 	for i, f := range stopFactors {
// 		pet.StopFactors[i] = string(f)
// 	}

// 	warnFactors := pet.GetWarnFactors(time.Now())
// 	pet.WarnFactors = make([]string, len(warnFactors))
// 	for i, f := range warnFactors {
// 		pet.WarnFactors[i] = string(f)
// 	}

// 	newPet, err := s.petRepo.Create(ctx, pet)
// 	if err != nil {
// 		return nil, apperrors.Internal(err, "failed to create pet")
// 	}

// 	return newPet, nil
// }

// // GetPet получает питомца с preload связанных данных
// func (s *PetService) GetPet(ctx context.Context, petID string, opts PetPreloadOptions) (*model.Pet, error) {
// 	pet, err := s.petRepo.GetPet(ctx, petID, opts)
// 	if err != nil {
// 		return nil, err
// 	}
// 	bloodReq, err := s.bloodReqRepo.GetByPetID(ctx, petID)
// 	if err != nil && !errors.Is(err, apperrors.ErrBloodRequestNotFound) {
// 		return nil, err
// 	}

// 	pet.PhotoURLs = s.BuildFullPhotoURLs(pet.PhotoURLs, *pet.UpdatedAt)

// 	// Set status based on blood request and stored DonorRestrictions
// 	if bloodReq != nil {
// 		if len(bloodReq.Edges.Responses) > 0 {
// 			pet.PetStatus = model.PetStatusBloodFound
// 		} else {
// 			pet.PetStatus = model.PetStatusRecipient
// 		}
// 	} else {
// 		if len(pet.StopFactors) > 0 {
// 			pet.PetStatus = model.PetStatusNone
// 		} else {
// 			pet.PetStatus = model.PetStatusDonor
// 		}
// 	}

// 	return pet, nil
// }

// // GetUserPets получает всех питомцев пользователя с preload связей
// func (s *PetService) GetUserPets(ctx context.Context, userID string, opts PetPreloadOptions) ([]*model.Pet, error) {
// 	_, _, err := s.userRepo.GetByID(ctx, userID, user.UserPreloadOptions{})
// 	if err != nil {
// 		if ent.IsNotFound(err) {
// 			return nil, apperrors.ErrUserNotFound
// 		}
// 		return nil, apperrors.Internal(err, "failed to get user")
// 	}

// 	pets, err := s.petRepo.GetPetsByUser(ctx, userID, opts)
// 	if err != nil {
// 		return nil, apperrors.Internal(err, "failed to get pets")
// 	}

// 	for i := range pets {
// 		pets[i].PhotoURLs = s.BuildFullPhotoURLs(pets[i].PhotoURLs, *pets[i].UpdatedAt)

// 		bloodReq, err := s.bloodReqRepo.GetByPetID(ctx, pets[i].ID)
// 		if err != nil && !errors.Is(err, apperrors.ErrBloodRequestNotFound) {
// 			return nil, err
// 		}

// 		// Set status based on blood request and stored DonorRestrictions
// 		if bloodReq != nil {
// 			if len(bloodReq.Edges.Responses) > 0 {
// 				pets[i].PetStatus = model.PetStatusBloodFound
// 			} else {
// 				pets[i].PetStatus = model.PetStatusRecipient
// 			}
// 		} else {
// 			if len(pets[i].StopFactors) > 0 {
// 				pets[i].PetStatus = model.PetStatusNone
// 			} else {
// 				pets[i].PetStatus = model.PetStatusDonor
// 			}
// 		}
// 	}

// 	return pets, nil
// }

// func (s *PetService) RevalidateDonor(ctx context.Context, petID string) (*model.Pet, error) {
// 	pet, err := s.petRepo.GetPet(ctx, petID, PetPreloadOptions{
// 		WithAll: true,
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	stopFactors := pet.GetStopFactors(time.Now())
// 	warnFactors := pet.GetWarnFactors(time.Now())

// 	// Create a new empty pet structure to update only StopFactors and WarnFactors
// 	updatePet := &model.Pet{}
// 	updatePet.StopFactors = make([]string, len(stopFactors))
// 	for i, f := range stopFactors {
// 		updatePet.StopFactors[i] = string(f)
// 	}
// 	updatePet.WarnFactors = make([]string, len(warnFactors))
// 	for i, f := range warnFactors {
// 		updatePet.WarnFactors[i] = string(f)
// 	}

// 	updatedPet, err := s.petRepo.Update(ctx, petID, updatePet)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return updatedPet, nil
// }

// // Update обновляет питомца и его связанные сущности
// func (s *PetService) Update(ctx context.Context, id string, petInput *model.Pet) (*model.Pet, error) {
// 	exists, err := s.petRepo.ExistsByID(ctx, id)
// 	if err != nil {
// 		return nil, apperrors.Internal(err, "failed to check pet existence")
// 	}
// 	if !exists {
// 		return nil, apperrors.ErrPetNotFound
// 	}

// 	// Calculate stop and warn factors and set DonorRestrictions
// 	stopFactors := petInput.GetStopFactors(time.Now())
// 	petInput.StopFactors = make([]string, len(stopFactors))
// 	for i, f := range stopFactors {
// 		petInput.StopFactors[i] = string(f)
// 	}

// 	warnFactors := petInput.GetWarnFactors(time.Now())
// 	petInput.WarnFactors = make([]string, len(warnFactors))
// 	for i, f := range warnFactors {
// 		petInput.WarnFactors[i] = string(f)
// 	}

// 	updatedPet, err := s.petRepo.Update(ctx, id, petInput)
// 	if err != nil {
// 		return nil, apperrors.Internal(err, "failed to update pet")
// 	}

// 	return updatedPet, nil
// }

// // DeletePet удаляет питомца по ID
// func (s *PetService) DeletePet(ctx context.Context, petID string) error {
// 	exists, err := s.petRepo.ExistsByID(ctx, petID)
// 	if err != nil {
// 		return apperrors.Internal(err, "failed to check pet existence")
// 	}
// 	if !exists {
// 		return apperrors.ErrPetNotFound
// 	}

// 	// Получаем все заявки на поиск крови, связанные с этим питомцем
// 	bloodRequests, err := s.bloodReqRepo.List(ctx, 0, 0, map[string]any{"pet_id": petID})
// 	if err != nil {
// 		return apperrors.Internal(err, "failed to list blood requests for pet")
// 	}

// 	// Удаляем каждую связанную заявку
// 	for _, req := range bloodRequests {
// 		if err := s.bloodReqRepo.Delete(ctx, req.ID); err != nil {
// 			slog.WarnContext(ctx, "Failed to delete blood request for pet", "blood_request_id", req.ID, "pet_id", petID, "error", err)
// 		}
// 	}

// 	if err := s.petRepo.Delete(ctx, petID); err != nil {
// 		return apperrors.Internal(err, "failed to delete pet")
// 	}

// 	return nil
// }

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
// func (s *PetService) BuildFullPhotoURLs(paths []string, updatedAt time.Time) []string {
// 	if len(paths) == 0 {
// 		return []string{}
// 	}
// 	result := make([]string, len(paths))
// 	for i, path := range paths {
// 		if path == "" {
// 			result[i] = ""
// 		} else {
// 			url := s.storage.GetPublicURLFromPath(path)
// 			result[i] = url + "?t=" + strconv.FormatInt(updatedAt.Unix(), 10)
// 		}
// 	}
// 	return result
// }
