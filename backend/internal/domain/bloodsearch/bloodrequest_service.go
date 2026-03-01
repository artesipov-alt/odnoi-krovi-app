package bloodsearch

// import (
// 	"context"
// 	"strconv"
// 	"time"

// 	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
// 	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/filestorage"
// 	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet"
// 	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent"
// 	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent/bloodsearchrequest"
// 	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/presistance"
// )

// // BloodSearchService реализует BloodSearchService
// type BloodSearchService struct {
// 	txManager presistance.TxManager
// 	bloodRepo BloodRequestRepository
// 	petRepo   pet.Repository
// 	donorRepo DonorResponseRepository
// 	storage   filestorage.Repository
// }

// // NewBloodSearchService создает новый экземпляр BloodSearchService
// func NewBloodSearchService(txManager presistance.TxManager, repo BloodRequestRepository, petRepo pet.Repository, donorRepo DonorResponseRepository, storage filestorage.Repository) *BloodSearchService {
// 	return &BloodSearchService{
// 		txManager: txManager,
// 		bloodRepo: repo,
// 		petRepo:   petRepo,
// 		donorRepo: donorRepo,
// 		storage:   storage,
// 	}
// }

// // CreateRequest создает новую заявку на поиск крови
// func (s *BloodSearchService) CreateRequest(ctx context.Context, bloodReq *ent.CreateBloodSearchRequestInput) (*ent.BloodSearchRequest, error) {
// 	// Проверяем существование питомца
// 	exists, err := s.petRepo.ExistsByID(ctx, bloodReq.PetID)
// 	if err != nil {
// 		return nil, apperrors.Internal(err, "failed to check pet existence")
// 	}
// 	if !exists {
// 		return nil, apperrors.ErrPetNotFound
// 	}

// 	// Проверяем, нет ли уже активной заявки для этого питомца
// 	activeExists, err := s.bloodRepo.ExistsByPetID(ctx, bloodReq.PetID)
// 	if err != nil {
// 		return nil, apperrors.Internal(err, "failed to check request existence")
// 	}
// 	if activeExists {
// 		return nil, apperrors.ErrBloodRequestAlreadyExists
// 	}

// 	bloodReq.Status = new(bloodsearchrequest.StatusActive)

// 	newReq, err := s.bloodRepo.Create(ctx, bloodReq)
// 	if err != nil {
// 		return nil, apperrors.Internal(err, "failed to create blood request")
// 	}

// 	newReq.PhotoUrls = s.BuildFullPhotoURLs(newReq.PhotoUrls, newReq.UpdatedAt)

// 	return newReq, nil
// }

// // GetRequestByID получает заявку по её ID
// func (s *BloodSearchService) GetRequestByID(ctx context.Context, id string) (*ent.BloodSearchRequest, error) {
// 	req, err := s.bloodRepo.GetByID(ctx, id)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Преобразуем пути к фото в полные URL
// 	req.PhotoUrls = s.BuildFullPhotoURLs(req.PhotoUrls, req.UpdatedAt)

// 	return req, nil
// }

// // GetRequestByPetID получает активную заявку для конкретного питомца
// func (s *BloodSearchService) GetRequestByPetID(ctx context.Context, petID string) (*ent.BloodSearchRequest, int, error) {
// 	req, err := s.bloodRepo.GetByPetID(ctx, petID)
// 	if err != nil {
// 		return nil, 0, err
// 	}
// 	situatableDonors, err := s.petRepo.CountSuitableDonors(ctx, req.BloodGroupNames)
// 	if err != nil {
// 		return nil, 0, err
// 	}

// 	// Преобразуем пути к фото в полные URL
// 	req.PhotoUrls = s.BuildFullPhotoURLs(req.PhotoUrls, req.UpdatedAt)

// 	return req, situatableDonors, nil
// }

// func (s *BloodSearchService) ApplyForBloodRequest(ctx context.Context, reqID, donorID string, conditions []string) (*ent.DonorResponse, error) {
// 	// Проверяем существование и статус заявки
// 	req, err := s.bloodRepo.GetByID(ctx, reqID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if req.Status != bloodsearchrequest.StatusActive {
// 		return nil, apperrors.ErrInvalidBloodRequestStatus.WithMessage("blood request is not active")
// 	}

// 	// Проверяем существование донора
// 	donorExists, err := s.petRepo.ExistsByID(ctx, donorID)
// 	if err != nil {
// 		return nil, apperrors.Internal(err, "failed to check donor existence")
// 	}
// 	if !donorExists {
// 		return nil, apperrors.ErrPetNotFound
// 	}

// 	// Проверяем, нет ли уже отклика от этого донора на эту заявку
// 	responses, err := s.donorRepo.GetDonorResponsesByDonorID(ctx, donorID)
// 	if err != nil {
// 		return nil, apperrors.Internal(err, "failed to check existing responses")
// 	}
// 	for _, resp := range responses {
// 		if resp.Edges.Request.ID == reqID {
// 			return nil, apperrors.ErrDonorResponseAlreadyExists
// 		}
// 	}

// 	resp, err := s.donorRepo.CreateDonorResponse(ctx, reqID, donorID, conditions)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return resp, nil
// }

// // UpdateRequest обновляет информацию о заявке
// func (s *BloodSearchService) UpdateRequest(ctx context.Context, id string, bloodReq *ent.UpdateBloodSearchRequestInput) (*ent.BloodSearchRequest, error) {
// 	// Проверяем существование
// 	req, err := s.bloodRepo.GetByID(ctx, id)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Сохраняем текущий статус, если он не передан
// 	if bloodReq.Status == nil {
// 		bloodReq.Status = &req.Status
// 	}

// 	updatedReq, err := s.bloodRepo.Update(ctx, id, bloodReq)
// 	if err != nil {
// 		return nil, apperrors.Internal(err, "failed to update blood request")
// 	}

// 	updatedReq.PhotoUrls = s.BuildFullPhotoURLs(updatedReq.PhotoUrls, updatedReq.UpdatedAt)

// 	return updatedReq, nil
// }

// // UpdateStatus обновляет статус заявки
// func (s *BloodSearchService) UpdateStatus(ctx context.Context, id string, status string) error {
// 	if err := bloodsearchrequest.StatusValidator(bloodsearchrequest.Status(status)); err != nil {
// 		return apperrors.ErrInvalidBloodRequestStatus.WithInternal(err)
// 	}

// 	// Получаем заявку, чтобы узнать PetID
// 	_, err := s.bloodRepo.GetByID(ctx, id)
// 	if err != nil {
// 		return err
// 	}

// 	err = s.txManager.WithTx(ctx, func(txCtx context.Context) error {
// 		entTx := ent.TxFromContext(txCtx)
// 		if entTx == nil {
// 			return apperrors.Internal(nil, "ent.TxFromContext returned nil")
// 		}

// 		// Обновляем статус заявки
// 		err = s.bloodRepo.UpdateStatus(txCtx, id, status)
// 		if err != nil {
// 			return apperrors.Internal(err, "failed to update status")
// 		}

// 		// Если статус закрыт, сбрасываем статус питомца на "none"
// 		// if status == "closed" {
// 		// 	err = s.petRepo.UpdateStatus(txCtx, req.PetID, "none")
// 		// 	if err != nil {
// 		// 		return apperrors.Internal(err, "failed to update pet status")
// 		// 	}
// 		// }
// 		return nil
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// // DeleteRequest удаляет заявку (soft delete)
// func (s *BloodSearchService) DeleteRequest(ctx context.Context, id string) error {
// 	// Получаем заявку, чтобы узнать PetID
// 	_, err := s.bloodRepo.GetByID(ctx, id)
// 	if err != nil {
// 		return apperrors.Internal(err, "failed to get blood request")
// 	}

// 	err = s.txManager.WithTx(ctx, func(txCtx context.Context) error {
// 		entTx := ent.TxFromContext(txCtx)
// 		if entTx == nil {
// 			return apperrors.Internal(nil, "ent.TxFromContext returned nil")
// 		}

// 		// // Обновляем статус питомца на "none"
// 		// err = s.petRepo.UpdateStatus(txCtx, req.PetID, "none")
// 		// if err != nil {
// 		// 	return apperrors.Internal(err, "failed to update pet status")
// 		// }

// 		// Удаляем заявку
// 		err = s.bloodRepo.Delete(txCtx, id)
// 		if err != nil {
// 			if ent.IsNotFound(err) {
// 				return apperrors.ErrBloodRequestNotFound
// 			}
// 			return apperrors.Internal(err, "failed to delete blood request")
// 		}
// 		return nil
// 	})

// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// // ListRequests возвращает список заявок с фильтрацией
// func (s *BloodSearchService) ListRequests(ctx context.Context, limit, offset int, filters map[string]any) ([]*ent.BloodSearchRequest, error) {
// 	requests, err := s.bloodRepo.List(ctx, limit, offset, filters)
// 	if err != nil {
// 		return nil, apperrors.Internal(err, "failed to list blood requests")
// 	}

// 	return requests, nil
// }

// // ExistsByID проверяет существование заявки по её ID
// func (s *BloodSearchService) ExistsByID(ctx context.Context, id string) (bool, error) {
// 	return s.bloodRepo.ExistsByID(ctx, id)
// }

// // BuildFullPhotoURLs преобразует пути к фото в полные публичные URL
// func (s *BloodSearchService) BuildFullPhotoURLs(paths []string, updatedAt time.Time) []string {
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
