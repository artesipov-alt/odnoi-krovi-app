package services

import (
	"context"
	"fmt"
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/ent"
	"github.com/artesipov-alt/odnoi-krovi-app/ent/bloodsearchrequest"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/repositories"
)

// DTOs for Blood Request operations
type (
	// BloodSearchPetRequest представляет запрос на добавление питомца в пул поиска крови
	BloodSearchPetRequest struct {
		PetID                  string   `json:"petId" validate:"required"`
		BloodVolumeNeeded      int32    `json:"bloodVolumeNeeded" validate:"required,gt=0"`
		BloodVolumeReserved    int32    `json:"bloodVolumeReserved"`
		Regions                []int32  `json:"regions" validate:"required,min=1"`
		SmallPetsNotifyAllowed bool     `json:"smallPetsNotifyAllowed"`
		Description            string   `json:"description"`
		PhotoUrls              []string `json:"photoUrls"`
		BloodGroupIds          []string `json:"bloodGroupIds"`
		BloodComponentIds      []int    `json:"bloodComponentIds"`
	}

	// BloodSearchPetResponse представляет ответ после создания заявки
	BloodSearchPetResponse struct {
		ID     string `json:"id"`
		PetID  string `json:"petId"`
		Status string `json:"status"`
	}

	// BloodSearchFilterRequest представляет фильтры для поиска заявок
	BloodSearchFilterRequest struct {
		PetID  string `json:"petId"`
		Status string `json:"status"`
		Limit  int    `json:"limit"`
		Offset int    `json:"offset"`
	}

	// BloodSearchRequestDTO представляет DTO для заявки на поиск крови
	BloodSearchRequestDTO struct {
		ID                     string    `json:"id"`
		PetID                  string    `json:"petId"`
		BloodVolumeNeeded      int32     `json:"bloodVolumeNeeded"`
		BloodVolumeReserved    int32     `json:"bloodVolumeReserved"`
		Regions                []int32   `json:"regions"`
		SmallPetsNotifyAllowed bool      `json:"smallPetsNotifyAllowed"`
		Description            string    `json:"description"`
		PhotoUrls              []string  `json:"photoUrls"`
		BloodGroupIds          []string  `json:"bloodGroupIds"`
		BloodComponentIds      []int     `json:"bloodComponentIds"`
		Status                 string    `json:"status"`
		CreatedAt              time.Time `json:"createdAt"`
		UpdatedAt              time.Time `json:"updatedAt"`
	}

	// BloodSearchPetsResponse представляет список заявок
	BloodSearchPetsResponse struct {
		Requests []BloodSearchRequestDTO `json:"requests"`
	}
)

// BloodSearchService определяет интерфейс для бизнес-логики заявок на поиск крови
type BloodSearchService interface {
	// CreateRequest создает новую заявку на поиск крови
	CreateRequest(ctx context.Context, req BloodSearchPetRequest) (BloodSearchPetResponse, error)

	// GetRequestByID получает заявку по её ID
	GetRequestByID(ctx context.Context, id string) (BloodSearchRequestDTO, error)

	// GetRequestByPetID получает активную заявку для конкретного питомца
	GetRequestByPetID(ctx context.Context, petID string) (BloodSearchRequestDTO, error)

	// UpdateRequest обновляет информацию о заявке
	UpdateRequest(ctx context.Context, id string, req BloodSearchPetRequest) (BloodSearchRequestDTO, error)

	// UpdateStatus обновляет статус заявки
	UpdateStatus(ctx context.Context, id string, status string) error

	// DeleteRequest удаляет заявку (soft delete)
	DeleteRequest(ctx context.Context, id string) error

	// ListRequests возвращает список заявок с фильтрацией
	ListRequests(ctx context.Context, limit, offset int, filters map[string]interface{}) ([]BloodSearchRequestDTO, error)
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
func (s *BloodSearchServiceImpl) CreateRequest(ctx context.Context, req BloodSearchPetRequest) (BloodSearchPetResponse, error) {
	// Проверяем существование питомца
	exists, err := s.petRepo.ExistsByID(ctx, req.PetID)
	if err != nil {
		return BloodSearchPetResponse{}, apperrors.Internal(err, "не удалось проверить существование питомца")
	}
	if !exists {
		return BloodSearchPetResponse{}, apperrors.BadRequest(fmt.Sprintf("питомец с ID %s не найден", req.PetID))
	}

	// Проверяем, нет ли уже активной заявки для этого питомца
	activeExists, err := s.repo.ExistsByPetID(ctx, req.PetID)
	if err != nil {
		return BloodSearchPetResponse{}, apperrors.Internal(err, "не удалось проверить наличие активных заявок")
	}
	if activeExists {
		return BloodSearchPetResponse{}, apperrors.BadRequest("для этого питомца уже есть активная заявка")
	}

	// Маппинг DTO в Ent модель
	bloodReq := &ent.BloodSearchRequest{
		PetID:                  req.PetID,
		BloodVolumeNeeded:      req.BloodVolumeNeeded,
		BloodVolumeReserved:    req.BloodVolumeReserved,
		Regions:                req.Regions,
		SmallPetsNotifyAllowed: req.SmallPetsNotifyAllowed,
		Description:            req.Description,
		PhotoUrls:              req.PhotoUrls,
		BloodGroupIds:          req.BloodGroupIds,
		BloodComponentIds:      req.BloodComponentIds,
		Status:                 bloodsearchrequest.StatusActive,
	}

	// Создаем заявку
	newReq, err := s.repo.Create(ctx, bloodReq)
	if err != nil {
		return BloodSearchPetResponse{}, apperrors.Internal(err, "не удалось создать заявку")
	}

	// Маппинг в DTO ответа
	resp := BloodSearchPetResponse{
		ID:     newReq.ID,
		PetID:  newReq.PetID,
		Status: string(newReq.Status),
	}

	return resp, nil
}

// GetRequestByID получает заявку по её ID
func (s *BloodSearchServiceImpl) GetRequestByID(ctx context.Context, id string) (BloodSearchRequestDTO, error) {
	req, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return BloodSearchRequestDTO{}, apperrors.NotFound(fmt.Sprintf("заявка с ID %s не найдена", id))
		}
		return BloodSearchRequestDTO{}, apperrors.Internal(err, "не удалось получить заявку")
	}

	// Маппинг в DTO
	dto := BloodSearchRequestDTO{
		ID:                     req.ID,
		PetID:                  req.PetID,
		BloodVolumeNeeded:      req.BloodVolumeNeeded,
		BloodVolumeReserved:    req.BloodVolumeReserved,
		Regions:                req.Regions,
		SmallPetsNotifyAllowed: req.SmallPetsNotifyAllowed,
		Description:            req.Description,
		PhotoUrls:              req.PhotoUrls,
		BloodGroupIds:          req.BloodGroupIds,
		BloodComponentIds:      req.BloodComponentIds,
		Status:                 string(req.Status),
		CreatedAt:              req.CreatedAt,
		UpdatedAt:              req.UpdatedAt,
	}

	return dto, nil
}

// GetRequestByPetID получает активную заявку для конкретного питомца
func (s *BloodSearchServiceImpl) GetRequestByPetID(ctx context.Context, petID string) (BloodSearchRequestDTO, error) {
	req, err := s.repo.GetByPetID(ctx, petID)
	if err != nil {
		if ent.IsNotFound(err) {
			return BloodSearchRequestDTO{}, apperrors.NotFound(fmt.Sprintf("активная заявка для питомца %s не найдена", petID))
		}
		return BloodSearchRequestDTO{}, apperrors.Internal(err, "не удалось получить заявку")
	}

	// Маппинг в DTO
	dto := BloodSearchRequestDTO{
		ID:                     req.ID,
		PetID:                  req.PetID,
		BloodVolumeNeeded:      req.BloodVolumeNeeded,
		BloodVolumeReserved:    req.BloodVolumeReserved,
		Regions:                req.Regions,
		SmallPetsNotifyAllowed: req.SmallPetsNotifyAllowed,
		Description:            req.Description,
		PhotoUrls:              req.PhotoUrls,
		BloodGroupIds:          req.BloodGroupIds,
		BloodComponentIds:      req.BloodComponentIds,
		Status:                 string(req.Status),
		CreatedAt:              req.CreatedAt,
		UpdatedAt:              req.UpdatedAt,
	}

	return dto, nil
}

// UpdateRequest обновляет информацию о заявке
func (s *BloodSearchServiceImpl) UpdateRequest(ctx context.Context, id string, req BloodSearchPetRequest) (BloodSearchRequestDTO, error) {
	// Проверяем существование
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return BloodSearchRequestDTO{}, apperrors.NotFound(fmt.Sprintf("заявка с ID %s не найдена", id))
		}
		return BloodSearchRequestDTO{}, apperrors.Internal(err, "не удалось получить заявку для обновления")
	}

	// Маппинг DTO в Ent модель
	updateReq := &ent.BloodSearchRequest{
		ID:                     id,
		PetID:                  req.PetID,
		BloodVolumeNeeded:      req.BloodVolumeNeeded,
		BloodVolumeReserved:    req.BloodVolumeReserved,
		Regions:                req.Regions,
		SmallPetsNotifyAllowed: req.SmallPetsNotifyAllowed,
		Description:            req.Description,
		PhotoUrls:              req.PhotoUrls,
		BloodGroupIds:          req.BloodGroupIds,
		BloodComponentIds:      req.BloodComponentIds,
		Status:                 existing.Status, // Сохраняем текущий статус
	}

	updated, err := s.repo.Update(ctx, updateReq)
	if err != nil {
		return BloodSearchRequestDTO{}, apperrors.Internal(err, "не удалось обновить заявку")
	}

	// Маппинг в DTO
	dto := BloodSearchRequestDTO{
		ID:                     updated.ID,
		PetID:                  updated.PetID,
		BloodVolumeNeeded:      updated.BloodVolumeNeeded,
		BloodVolumeReserved:    updated.BloodVolumeReserved,
		Regions:                updated.Regions,
		SmallPetsNotifyAllowed: updated.SmallPetsNotifyAllowed,
		Description:            updated.Description,
		PhotoUrls:              updated.PhotoUrls,
		BloodGroupIds:          updated.BloodGroupIds,
		BloodComponentIds:      updated.BloodComponentIds,
		Status:                 string(updated.Status),
		CreatedAt:              updated.CreatedAt,
		UpdatedAt:              updated.UpdatedAt,
	}

	return dto, nil
}

// UpdateStatus обновляет статус заявки
func (s *BloodSearchServiceImpl) UpdateStatus(ctx context.Context, id string, status string) error {
	if err := bloodsearchrequest.StatusValidator(bloodsearchrequest.Status(status)); err != nil {
		return apperrors.BadRequest("недопустимый статус заявки")
	}

	err := s.repo.UpdateStatus(ctx, id, status)
	if err != nil {
		if ent.IsNotFound(err) {
			return apperrors.NotFound(fmt.Sprintf("заявка с ID %s не найдена", id))
		}
		return apperrors.Internal(err, "не удалось обновить статус заявки")
	}

	return nil
}

// DeleteRequest удаляет заявку (soft delete)
func (s *BloodSearchServiceImpl) DeleteRequest(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return apperrors.NotFound(fmt.Sprintf("заявка с ID %s не найдена", id))
		}
		return apperrors.Internal(err, "не удалось удалить заявку")
	}
	return nil
}

// ListRequests возвращает список заявок с фильтрацией
func (s *BloodSearchServiceImpl) ListRequests(ctx context.Context, limit, offset int, filters map[string]interface{}) ([]BloodSearchRequestDTO, error) {
	requests, err := s.repo.List(ctx, limit, offset, filters)
	if err != nil {
		return nil, apperrors.Internal(err, "не удалось получить список заявок")
	}

	// Маппинг в DTO
	dtos := make([]BloodSearchRequestDTO, len(requests))
	for i, req := range requests {
		dtos[i] = BloodSearchRequestDTO{
			ID:                     req.ID,
			PetID:                  req.PetID,
			BloodVolumeNeeded:      req.BloodVolumeNeeded,
			BloodVolumeReserved:    req.BloodVolumeReserved,
			Regions:                req.Regions,
			SmallPetsNotifyAllowed: req.SmallPetsNotifyAllowed,
			Description:            req.Description,
			PhotoUrls:              req.PhotoUrls,
			BloodGroupIds:          req.BloodGroupIds,
			BloodComponentIds:      req.BloodComponentIds,
			Status:                 string(req.Status),
			CreatedAt:              req.CreatedAt,
			UpdatedAt:              req.UpdatedAt,
		}
	}

	return dtos, nil
}
