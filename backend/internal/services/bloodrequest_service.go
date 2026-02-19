package services

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/ent"
	"github.com/artesipov-alt/odnoi-krovi-app/ent/bloodsearchrequest"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/repositories"
)

// BloodRequestRepository определяет интерфейс для работы с данными заявок на поиск крови питомцев
type BloodRequestRepository interface {
	// Create создает новую заявку на поиск крови
	Create(ctx context.Context, request *ent.CreateBloodSearchRequestInput) (*ent.BloodSearchRequest, error)

	// GetByID возвращает заявку по её идентификатору
	GetByID(ctx context.Context, id string) *ent.BloodSearchRequestQuery

	// GetByPetID возвращает заявку по идентификатору питомца
	GetByPetID(ctx context.Context, petID string) *ent.BloodSearchRequestQuery

	// Update обновляет информацию о заявке
	Update(ctx context.Context, id string, request *ent.UpdateBloodSearchRequestInput) (*ent.BloodSearchRequest, error)

	// UpdateStatus обновляет статус заявки
	UpdateStatus(ctx context.Context, id string, status string) error

	// Delete удаляет заявку из хранилища
	Delete(ctx context.Context, id string) error

	// List возвращает список заявок с фильтрацией и пагинацией
	List(ctx context.Context, limit, offset int, filters map[string]any) ([]*ent.BloodSearchRequest, error)

	// ExistsByPetID проверяет существование активной заявки для питомца
	ExistsByPetID(ctx context.Context, petID string) (bool, error)

	// ExistsByID проверяет существование заявки по её идентификатору
	ExistsByID(ctx context.Context, id string) (bool, error)

	// Count возвращает общее количество заявок в хранилище
	Count(ctx context.Context) (int, error)

	// AddPhotoURLs добавляет новые пути к фотографиям заявки
	AddPhotoURLs(ctx context.Context, id string, paths []string) error
}

// BloodSearchService реализует BloodSearchService
type BloodSearchService struct {
	txManager repositories.TxManager
	bloodRepo BloodRequestRepository
	petRepo   repositories.PetRepository
	storage   repositories.FileStorage
}

// NewBloodSearchService создает новый экземпляр BloodSearchService
func NewBloodSearchService(txManager repositories.TxManager, repo BloodRequestRepository, petRepo repositories.PetRepository, storage repositories.FileStorage) *BloodSearchService {
	return &BloodSearchService{
		txManager: txManager,
		bloodRepo: repo,
		petRepo:   petRepo,
		storage:   storage,
	}
}

// CreateRequest создает новую заявку на поиск крови
func (s *BloodSearchService) CreateRequest(ctx context.Context, bloodReq *ent.CreateBloodSearchRequestInput) (*ent.BloodSearchRequest, error) {
	// Проверяем существование питомца
	exists, err := s.petRepo.ExistsByID(ctx, bloodReq.PetID)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to check pet existence")
	}
	if !exists {
		return nil, apperrors.ErrPetNotFound
	}

	// Проверяем, нет ли уже активной заявки для этого питомца
	activeExists, err := s.bloodRepo.ExistsByPetID(ctx, bloodReq.PetID)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to check request existence")
	}
	if activeExists {
		return nil, apperrors.ErrBloodRequestAlreadyExists
	}

	bloodReq.Status = new(bloodsearchrequest.StatusActive)

	var newReq *ent.BloodSearchRequest
	err = s.txManager.WithTx(ctx, func(txCtx context.Context) error {
		entTx := ent.TxFromContext(txCtx)
		if entTx == nil {
			return apperrors.Internal(nil, "ent.TxFromContext returned nil")
		}

		newReq, err = s.bloodRepo.Create(txCtx, bloodReq)
		if err != nil {
			return apperrors.Internal(err, "failed to create blood request")
		}

		err = s.petRepo.UpdateStatus(txCtx, bloodReq.PetID, "recipient")
		if err != nil {
			return apperrors.Internal(err, "failed to update pet status")
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return newReq, nil
}

// GetRequestByID получает заявку по её ID
func (s *BloodSearchService) GetRequestByID(ctx context.Context, id string) (*ent.BloodSearchRequest, error) {
	reqQuery := s.bloodRepo.GetByID(ctx, id)

	req, err := reqQuery.Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) { // This case should ideally be caught by the initial GetByID, but good for defensive programming
			return nil, apperrors.ErrBloodRequestNotFound
		}
		return nil, apperrors.Internal(err, "failed to execute blood request query")
	}

	// Преобразуем пути к фото в полные URL
	req.PhotoUrls = s.BuildFullPhotoURLs(req.PhotoUrls)

	return req, nil
}

// GetRequestByPetID получает активную заявку для конкретного питомца
func (s *BloodSearchService) GetRequestByPetID(ctx context.Context, petID string) (*ent.BloodSearchRequest, error) {
	req, err := s.bloodRepo.GetByPetID(ctx, petID).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, apperrors.ErrBloodRequestNotFound
		}
		return nil, apperrors.Internal(err, "failed to get blood request by pet ID")
	}

	// Преобразуем пути к фото в полные URL
	req.PhotoUrls = s.BuildFullPhotoURLs(req.PhotoUrls)

	return req, nil
}

// UpdateRequest обновляет информацию о заявке
func (s *BloodSearchService) UpdateRequest(ctx context.Context, id string, bloodReq *ent.UpdateBloodSearchRequestInput) (*ent.BloodSearchRequest, error) {
	// Проверяем существование
	req, err := s.bloodRepo.GetByID(ctx, id).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, apperrors.ErrBloodRequestNotFound
		}
		return nil, apperrors.Internal(err, "failed to get blood request")
	}

	// Сохраняем текущий статус, если он не передан
	if bloodReq.Status == nil {
		bloodReq.Status = &req.Status
	}

	updated, err := s.bloodRepo.Update(ctx, id, bloodReq)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to update blood request")
	}

	return updated, nil
}

// UpdateStatus обновляет статус заявки
func (s *BloodSearchService) UpdateStatus(ctx context.Context, id string, status string) error {
	if err := bloodsearchrequest.StatusValidator(bloodsearchrequest.Status(status)); err != nil {
		return apperrors.ErrInvalidBloodRequestStatus.WithInternal(err)
	}

	// Получаем заявку, чтобы узнать PetID
	req, err := s.bloodRepo.GetByID(ctx, id).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return apperrors.ErrBloodRequestNotFound
		}
		return apperrors.Internal(err, "failed to get blood request")
	}

	err = s.txManager.WithTx(ctx, func(txCtx context.Context) error {
		entTx := ent.TxFromContext(txCtx)
		if entTx == nil {
			return apperrors.Internal(nil, "ent.TxFromContext returned nil")
		}

		// Обновляем статус заявки
		err = s.bloodRepo.UpdateStatus(txCtx, id, status)
		if err != nil {
			if ent.IsNotFound(err) {
				return apperrors.ErrBloodRequestNotFound
			}
			return apperrors.Internal(err, "failed to update status")
		}

		// Если статус закрыт, сбрасываем статус питомца на "none"
		if status == "closed" {
			err = s.petRepo.UpdateStatus(txCtx, req.PetID, "none")
			if err != nil {
				return apperrors.Internal(err, "failed to update pet status")
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// DeleteRequest удаляет заявку (soft delete)
func (s *BloodSearchService) DeleteRequest(ctx context.Context, id string) error {
	// Получаем заявку, чтобы узнать PetID
	req, err := s.bloodRepo.GetByID(ctx, id).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return apperrors.ErrBloodRequestNotFound
		}
		return apperrors.Internal(err, "failed to get blood request")
	}

	err = s.txManager.WithTx(ctx, func(txCtx context.Context) error {
		entTx := ent.TxFromContext(txCtx)
		if entTx == nil {
			return apperrors.Internal(nil, "ent.TxFromContext returned nil")
		}

		// Обновляем статус питомца на "none"
		err = s.petRepo.UpdateStatus(txCtx, req.PetID, "none")
		if err != nil {
			return apperrors.Internal(err, "failed to update pet status")
		}

		// Удаляем заявку
		err = s.bloodRepo.Delete(txCtx, id)
		if err != nil {
			if ent.IsNotFound(err) {
				return apperrors.ErrBloodRequestNotFound
			}
			return apperrors.Internal(err, "failed to delete blood request")
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

// ListRequests возвращает список заявок с фильтрацией
func (s *BloodSearchService) ListRequests(ctx context.Context, limit, offset int, filters map[string]any) ([]*ent.BloodSearchRequest, error) {
	requests, err := s.bloodRepo.List(ctx, limit, offset, filters)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to list blood requests")
	}

	return requests, nil
}

// ExistsByID проверяет существование заявки по её ID
func (s *BloodSearchService) ExistsByID(ctx context.Context, id string) (bool, error) {
	return s.bloodRepo.ExistsByID(ctx, id)
}

// buildFullPhotoURLs преобразует пути к фото в полные публичные URL
func (s *BloodSearchService) BuildFullPhotoURLs(paths []string) []string {
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
