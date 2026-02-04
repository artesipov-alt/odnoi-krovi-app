package services

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/ent"
	"github.com/artesipov-alt/odnoi-krovi-app/ent/bloodsearchrequest"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/repositories"
)

// BloodSearchService определяет интерфейс для бизнес-логики заявок на поиск крови
type BloodSearchService interface {
	// CreateRequest создает новую заявку на поиск крови
	CreateRequest(ctx context.Context, bloodReq *ent.BloodSearchRequest) (*ent.BloodSearchRequest, error)

	// GetRequestByID получает заявку по её ID
	GetRequestByID(ctx context.Context, id string) (*ent.BloodSearchRequest, error)

	// GetRequestByPetID получает активную заявку для конкретного питомца
	GetRequestByPetID(ctx context.Context, petID string) (*ent.BloodSearchRequest, error)

	// UpdateRequest обновляет информацию о заявке
	UpdateRequest(ctx context.Context, id string, bloodReq *ent.BloodSearchRequest) (*ent.BloodSearchRequest, error)

	// UpdateStatus обновляет статус заявки
	UpdateStatus(ctx context.Context, id string, status string) error

	// ExistsByID проверяет существование заявки по её ID
	ExistsByID(ctx context.Context, id string) (bool, error)

	// DeleteRequest удаляет заявку (soft delete)
	DeleteRequest(ctx context.Context, id string) error

	// ListRequests возвращает список заявок с фильтрацией
	ListRequests(ctx context.Context, limit, offset int, filters map[string]any) ([]*ent.BloodSearchRequest, error)

	// buildFullPhotoURLs преобразует пути к фото в полные публичные URL
	buildFullPhotoURLs(paths []string) []string
}

// BloodSearchServiceImpl реализует BloodSearchService
type BloodSearchServiceImpl struct {
	bloodRepo repositories.BloodRequestRepository
	petRepo   repositories.PetRepository
	storage   repositories.FileStorage
}

// NewBloodSearchService создает новый экземпляр BloodSearchService
func NewBloodSearchService(repo repositories.BloodRequestRepository, petRepo repositories.PetRepository, storage repositories.FileStorage) *BloodSearchServiceImpl {
	return &BloodSearchServiceImpl{
		bloodRepo: repo,
		petRepo:   petRepo,
		storage:   storage,
	}
}

// CreateRequest создает новую заявку на поиск крови
func (s *BloodSearchServiceImpl) CreateRequest(ctx context.Context, bloodReq *ent.BloodSearchRequest) (*ent.BloodSearchRequest, error) {
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

	// Устанавливаем статус по умолчанию
	bloodReq.Status = bloodsearchrequest.StatusActive

	// Создаем заявку
	newReq, err := s.bloodRepo.Create(ctx, bloodReq)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to create blood request")
	}

	return newReq, nil
}

// GetRequestByID получает заявку по её ID
func (s *BloodSearchServiceImpl) GetRequestByID(ctx context.Context, id string) (*ent.BloodSearchRequest, error) {
	req, err := s.bloodRepo.GetByID(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, apperrors.ErrBloodRequestNotFound
		}
		return nil, apperrors.Internal(err, "failed to get blood request")
	}

	// Преобразуем пути к фото в полные URL
	req.PhotoUrls = s.buildFullPhotoURLs(req.PhotoUrls)

	return req, nil
}

// GetRequestByPetID получает активную заявку для конкретного питомца
func (s *BloodSearchServiceImpl) GetRequestByPetID(ctx context.Context, petID string) (*ent.BloodSearchRequest, error) {
	req, err := s.bloodRepo.GetByPetID(ctx, petID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, apperrors.ErrBloodRequestNotFound
		}
		return nil, apperrors.Internal(err, "failed to get blood request by pet ID")
	}

	return req, nil
}

// UpdateRequest обновляет информацию о заявке
func (s *BloodSearchServiceImpl) UpdateRequest(ctx context.Context, id string, bloodReq *ent.BloodSearchRequest) (*ent.BloodSearchRequest, error) {
	// Проверяем существование
	existing, err := s.bloodRepo.GetByID(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, apperrors.ErrBloodRequestNotFound
		}
		return nil, apperrors.Internal(err, "failed to get blood request")
	}

	// Убеждаемся, что ID совпадает
	bloodReq.ID = id
	// Сохраняем текущий статус, если он не передан
	if bloodReq.Status == "" {
		bloodReq.Status = existing.Status
	}

	updated, err := s.bloodRepo.Update(ctx, bloodReq)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to update blood request")
	}

	return updated, nil
}

// UpdateStatus обновляет статус заявки
func (s *BloodSearchServiceImpl) UpdateStatus(ctx context.Context, id string, status string) error {
	if err := bloodsearchrequest.StatusValidator(bloodsearchrequest.Status(status)); err != nil {
		return apperrors.ErrInvalidBloodRequestStatus.WithInternal(err)
	}

	err := s.bloodRepo.UpdateStatus(ctx, id, status)
	if err != nil {
		if ent.IsNotFound(err) {
			return apperrors.ErrBloodRequestNotFound
		}
		return apperrors.Internal(err, "failed to update status")
	}

	return nil
}

// DeleteRequest удаляет заявку (soft delete)
func (s *BloodSearchServiceImpl) DeleteRequest(ctx context.Context, id string) error {
	err := s.bloodRepo.Delete(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return apperrors.ErrBloodRequestNotFound
		}
		return apperrors.Internal(err, "failed to delete blood request")
	}
	return nil
}

// ListRequests возвращает список заявок с фильтрацией
func (s *BloodSearchServiceImpl) ListRequests(ctx context.Context, limit, offset int, filters map[string]any) ([]*ent.BloodSearchRequest, error) {
	requests, err := s.bloodRepo.List(ctx, limit, offset, filters)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to list blood requests")
	}

	return requests, nil
}

// ExistsByID проверяет существование заявки по её ID
func (s *BloodSearchServiceImpl) ExistsByID(ctx context.Context, id string) (bool, error) {
	return s.bloodRepo.ExistsByID(ctx, id)
}

// buildFullPhotoURLs преобразует пути к фото в полные публичные URL
func (s *BloodSearchServiceImpl) buildFullPhotoURLs(paths []string) []string {
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
