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

	// CreateWithTx создает новую заявку на поиск крови в рамках транзакции
	CreateWithTx(ctx context.Context, tx *ent.Tx, request *ent.CreateBloodSearchRequestInput) (*ent.BloodSearchRequest, error)

	// GetByID возвращает заявку по её идентификатору
	GetByID(ctx context.Context, id string) (*ent.BloodSearchRequest, error)

	// GetByPetID возвращает заявку по идентификатору питомца
	GetByPetID(ctx context.Context, petID string) (*ent.BloodSearchRequest, error)

	// Update обновляет информацию о заявке
	Update(ctx context.Context, id string, request *ent.UpdateBloodSearchRequestInput) (*ent.BloodSearchRequest, error)

	// UpdateStatus обновляет статус заявки
	UpdateStatus(ctx context.Context, id string, status string) error

	// UpdateStatusWithTx обновляет статус заявки в рамках транзакции
	UpdateStatusWithTx(ctx context.Context, tx *ent.Tx, id string, status string) error

	// Delete удаляет заявку из хранилища
	Delete(ctx context.Context, id string) error

	// DeleteWithTx удаляет заявку из хранилища в рамках транзакции
	DeleteWithTx(ctx context.Context, tx *ent.Tx, id string) error

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
	bloodRepo BloodRequestRepository
	petRepo   repositories.PetRepository
	storage   repositories.FileStorage
	client    *ent.Client
}

// NewBloodSearchService создает новый экземпляр BloodSearchService
func NewBloodSearchService(repo BloodRequestRepository, petRepo repositories.PetRepository, storage repositories.FileStorage, client *ent.Client) *BloodSearchService {
	return &BloodSearchService{
		bloodRepo: repo,
		petRepo:   petRepo,
		storage:   storage,
		client:    client,
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

	// Устанавливаем статус по умолчанию
	bloodReq.Status = new(bloodsearchrequest.StatusActive)

	// Создаем транзакцию
	tx, err := s.client.Tx(ctx)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to start transaction")
	}

	// Создаем заявку в рамках транзакции
	newReq, err := s.bloodRepo.CreateWithTx(ctx, tx, bloodReq)
	if err != nil {
		tx.Rollback()
		return nil, apperrors.Internal(err, "failed to create blood request")
	}

	// Обновляем статус питомца на "recipient"
	err = s.petRepo.UpdateStatusWithTx(ctx, tx, bloodReq.PetID, "recipient")
	if err != nil {
		tx.Rollback()
		return nil, apperrors.Internal(err, "failed to update pet status")
	}

	// Коммитим транзакцию
	if err := tx.Commit(); err != nil {
		return nil, apperrors.Internal(err, "failed to commit transaction")
	}

	return newReq, nil
}

// GetRequestByID получает заявку по её ID
func (s *BloodSearchService) GetRequestByID(ctx context.Context, id string) (*ent.BloodSearchRequest, error) {
	req, err := s.bloodRepo.GetByID(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, apperrors.ErrBloodRequestNotFound
		}
		return nil, apperrors.Internal(err, "failed to get blood request")
	}

	// Преобразуем пути к фото в полные URL
	req.PhotoUrls = s.BuildFullPhotoURLs(req.PhotoUrls)

	return req, nil
}

// GetRequestByPetID получает активную заявку для конкретного питомца
func (s *BloodSearchService) GetRequestByPetID(ctx context.Context, petID string) (*ent.BloodSearchRequest, error) {
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
func (s *BloodSearchService) UpdateRequest(ctx context.Context, id string, bloodReq *ent.UpdateBloodSearchRequestInput) (*ent.BloodSearchRequest, error) {
	// Проверяем существование
	existing, err := s.bloodRepo.GetByID(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, apperrors.ErrBloodRequestNotFound
		}
		return nil, apperrors.Internal(err, "failed to get blood request")
	}

	// Сохраняем текущий статус, если он не передан
	if bloodReq.Status == nil {
		bloodReq.Status = &existing.Status
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
	req, err := s.bloodRepo.GetByID(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return apperrors.ErrBloodRequestNotFound
		}
		return apperrors.Internal(err, "failed to get blood request")
	}

	// Создаем транзакцию
	tx, err := s.client.Tx(ctx)
	if err != nil {
		return apperrors.Internal(err, "failed to start transaction")
	}

	// Обновляем статус заявки
	err = s.bloodRepo.UpdateStatusWithTx(ctx, tx, id, status)
	if err != nil {
		tx.Rollback()
		if ent.IsNotFound(err) {
			return apperrors.ErrBloodRequestNotFound
		}
		return apperrors.Internal(err, "failed to update status")
	}

	// Если статус закрыт, сбрасываем статус питомца на "none"
	if status == "closed" {
		err = s.petRepo.UpdateStatusWithTx(ctx, tx, req.PetID, "none")
		if err != nil {
			tx.Rollback()
			return apperrors.Internal(err, "failed to update pet status")
		}
	}

	// Коммитим транзакцию
	if err := tx.Commit(); err != nil {
		return apperrors.Internal(err, "failed to commit transaction")
	}

	return nil
}

// DeleteRequest удаляет заявку (soft delete)
func (s *BloodSearchService) DeleteRequest(ctx context.Context, id string) error {
	// Получаем заявку, чтобы узнать PetID
	req, err := s.bloodRepo.GetByID(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return apperrors.ErrBloodRequestNotFound
		}
		return apperrors.Internal(err, "failed to get blood request")
	}

	// Создаем транзакцию
	tx, err := s.client.Tx(ctx)
	if err != nil {
		return apperrors.Internal(err, "failed to start transaction")
	}

	// Обновляем статус питомца на "none"
	err = s.petRepo.UpdateStatusWithTx(ctx, tx, req.PetID, "none")
	if err != nil {
		tx.Rollback()
		return apperrors.Internal(err, "failed to update pet status")
	}

	// Удаляем заявку
	err = s.bloodRepo.DeleteWithTx(ctx, tx, id)
	if err != nil {
		tx.Rollback()
		if ent.IsNotFound(err) {
			return apperrors.ErrBloodRequestNotFound
		}
		return apperrors.Internal(err, "failed to delete blood request")
	}

	// Коммитим транзакцию
	if err := tx.Commit(); err != nil {
		return apperrors.Internal(err, "failed to commit transaction")
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
