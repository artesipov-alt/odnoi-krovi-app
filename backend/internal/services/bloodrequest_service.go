package services

import (
	"context"
	"errors"

	"github.com/artesipov-alt/odnoi-krovi-app/ent"
	"github.com/artesipov-alt/odnoi-krovi-app/ent/bloodsearchrequest"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/repositories"
)

var (
	ErrBloodRequestNotFound = errors.New("blood search request not found")
	ErrBloodRequestExists   = errors.New("blood search request already exists for this pet")
	ErrInvalidStatus        = errors.New("invalid status")
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

	// DeleteRequest удаляет заявку (soft delete)
	DeleteRequest(ctx context.Context, id string) error

	// ListRequests возвращает список заявок с фильтрацией
	ListRequests(ctx context.Context, limit, offset int, filters map[string]any) ([]*ent.BloodSearchRequest, error)
}

// BloodSearchServiceImpl реализует BloodSearchService
type BloodSearchServiceImpl struct {
	repo    repositories.BloodRequestRepository
	petRepo repositories.PetRepository
}

// NewBloodSearchService создает новый экземпляр BloodSearchService
func NewBloodSearchService(repo repositories.BloodRequestRepository, petRepo repositories.PetRepository) *BloodSearchServiceImpl {
	return &BloodSearchServiceImpl{
		repo:    repo,
		petRepo: petRepo,
	}
}

// CreateRequest создает новую заявку на поиск крови
func (s *BloodSearchServiceImpl) CreateRequest(ctx context.Context, bloodReq *ent.BloodSearchRequest) (*ent.BloodSearchRequest, error) {
	// Проверяем существование питомца
	exists, err := s.petRepo.ExistsByID(ctx, bloodReq.PetID)
	if err != nil {
		return nil, ErrInternal
	}
	if !exists {
		return nil, ErrPetNotFound
	}

	// Проверяем, нет ли уже активной заявки для этого питомца
	activeExists, err := s.repo.ExistsByPetID(ctx, bloodReq.PetID)
	if err != nil {
		return nil, ErrInternal
	}
	if activeExists {
		return nil, ErrBloodRequestExists
	}

	// Устанавливаем статус по умолчанию
	bloodReq.Status = bloodsearchrequest.StatusActive

	// Создаем заявку
	newReq, err := s.repo.Create(ctx, bloodReq)
	if err != nil {
		return nil, ErrInternal
	}

	return newReq, nil
}

// GetRequestByID получает заявку по её ID
func (s *BloodSearchServiceImpl) GetRequestByID(ctx context.Context, id string) (*ent.BloodSearchRequest, error) {
	req, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrBloodRequestNotFound
		}
		return nil, ErrInternal
	}

	return req, nil
}

// GetRequestByPetID получает активную заявку для конкретного питомца
func (s *BloodSearchServiceImpl) GetRequestByPetID(ctx context.Context, petID string) (*ent.BloodSearchRequest, error) {
	req, err := s.repo.GetByPetID(ctx, petID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrBloodRequestNotFound
		}
		return nil, ErrInternal
	}

	return req, nil
}

// UpdateRequest обновляет информацию о заявке
func (s *BloodSearchServiceImpl) UpdateRequest(ctx context.Context, id string, bloodReq *ent.BloodSearchRequest) (*ent.BloodSearchRequest, error) {
	// Проверяем существование
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrBloodRequestNotFound
		}
		return nil, ErrInternal
	}

	// Убеждаемся, что ID совпадает
	bloodReq.ID = id
	// Сохраняем текущий статус, если он не передан
	if bloodReq.Status == "" {
		bloodReq.Status = existing.Status
	}

	updated, err := s.repo.Update(ctx, bloodReq)
	if err != nil {
		return nil, ErrInternal
	}

	return updated, nil
}

// UpdateStatus обновляет статус заявки
func (s *BloodSearchServiceImpl) UpdateStatus(ctx context.Context, id string, status string) error {
	if err := bloodsearchrequest.StatusValidator(bloodsearchrequest.Status(status)); err != nil {
		return ErrInvalidStatus
	}

	err := s.repo.UpdateStatus(ctx, id, status)
	if err != nil {
		if ent.IsNotFound(err) {
			return ErrBloodRequestNotFound
		}
		return ErrInternal
	}

	return nil
}

// DeleteRequest удаляет заявку (soft delete)
func (s *BloodSearchServiceImpl) DeleteRequest(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return ErrBloodRequestNotFound
		}
		return ErrInternal
	}
	return nil
}

// ListRequests возвращает список заявок с фильтрацией
func (s *BloodSearchServiceImpl) ListRequests(ctx context.Context, limit, offset int, filters map[string]any) ([]*ent.BloodSearchRequest, error) {
	requests, err := s.repo.List(ctx, limit, offset, filters)
	if err != nil {
		return nil, ErrInternal
	}

	return requests, nil
}
