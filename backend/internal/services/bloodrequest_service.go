package services

import (
	"context"
	"strings"

	"github.com/artesipov-alt/odnoi-krovi-app/ent"
	"github.com/artesipov-alt/odnoi-krovi-app/ent/bloodgroup"
	"github.com/artesipov-alt/odnoi-krovi-app/ent/bloodsearchrequest"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/repositories"
)

// BloodRequestRepository определяет интерфейс для работы с данными заявок на поиск крови питомцев
type BloodRequestRepository interface {
	// Create создает новую заявку на поиск крови
	Create(ctx context.Context, request *ent.CreateBloodSearchRequestInput) (*ent.BloodSearchRequest, error)

	// GetByID возвращает заявку по её идентификатору
	GetByID(ctx context.Context, id string) (*ent.BloodSearchRequest, error)

	// GetByPetID возвращает заявку по идентификатору питомца
	GetByPetID(ctx context.Context, petID string) (*ent.BloodSearchRequest, error)

	// Update обновляет информацию о заявке
	Update(ctx context.Context, id string, request *ent.UpdateBloodSearchRequestInput) error

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

// BloodInfoRepository определяет интерфейс для работы с типами крови
type BloodInfoRepository interface {
	AllComponents(ctx context.Context) ([]*ent.BloodComponent, error)
	ComponentByID(ctx context.Context, id string) (*ent.BloodComponent, error)
	BloodGroupsByPetType(ctx context.Context, petType bloodgroup.PetType) ([]*ent.BloodGroup, error)
	FindByTypeAndBloodGroup(ctx context.Context, petType bloodgroup.PetType, bloodGroup string) (*ent.BloodGroup, error)
	FindByBloodGroup(ctx context.Context, bloodGroup string) (*ent.BloodGroup, error)
}

// DonorResponseRepository определяет интерфейс для работы с откликами доноров
type DonorResponseRepository interface {
	CreateDonorResponse(ctx context.Context, reqID, donorID string, amountML int32) error
	GetDonorResponseByID(ctx context.Context, id string) (*ent.DonorResponse, error)
	UpdateDonorResponseStatus(ctx context.Context, id, status string) error
	DeleteDonorResponse(ctx context.Context, id string) error
	GetDonorResponsesByRequestID(ctx context.Context, reqID string) ([]*ent.DonorResponse, error)
	GetDonorResponsesByDonorID(ctx context.Context, donorID string) ([]*ent.DonorResponse, error)
	ExistsByID(ctx context.Context, id string) (bool, error)
	ExistsByRequestID(ctx context.Context, reqID string) (bool, error)
	ExistsByDonorID(ctx context.Context, donorID string) (bool, error)
	Count(ctx context.Context) (int, error)
}

// BloodSearchService реализует BloodSearchService
type BloodSearchService struct {
	txManager repositories.TxManager
	bloodRepo BloodRequestRepository
	petRepo   PetRepository
	donorRepo DonorResponseRepository
	storage   FileStorage
}

// NewBloodSearchService создает новый экземпляр BloodSearchService
func NewBloodSearchService(txManager repositories.TxManager, repo BloodRequestRepository, petRepo PetRepository, donorRepo DonorResponseRepository, storage FileStorage) *BloodSearchService {
	return &BloodSearchService{
		txManager: txManager,
		bloodRepo: repo,
		petRepo:   petRepo,
		donorRepo: donorRepo,
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
	req, err := s.bloodRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Преобразуем пути к фото в полные URL
	req.PhotoUrls = s.BuildFullPhotoURLs(req.PhotoUrls)

	return req, nil
}

// GetRequestByPetID получает активную заявку для конкретного питомца
func (s *BloodSearchService) GetRequestByPetID(ctx context.Context, petID string) (*ent.BloodSearchRequest, error) {
	req, err := s.bloodRepo.GetByPetID(ctx, petID)
	if err != nil {
		return nil, err
	}

	// Преобразуем пути к фото в полные URL
	req.PhotoUrls = s.BuildFullPhotoURLs(req.PhotoUrls)

	return req, nil
}

func (s *BloodSearchService) ApplyForBloodRequest(ctx context.Context, reqID, donorID string, amountML int32) error {
	// Проверяем существование и статус заявки
	req, err := s.bloodRepo.GetByID(ctx, reqID)
	if err != nil {
		return err
	}
	if req.Status != bloodsearchrequest.StatusActive {
		return apperrors.ErrInvalidBloodRequestStatus.WithMessage("blood request is not active")
	}

	// Проверяем существование донора
	donorExists, err := s.petRepo.ExistsByID(ctx, donorID)
	if err != nil {
		return apperrors.Internal(err, "failed to check donor existence")
	}
	if !donorExists {
		return apperrors.ErrPetNotFound
	}

	// Проверяем статус донора
	donor, err := s.petRepo.GetPet(ctx, donorID, PetPreloadOptions{})
	if err != nil {
		return err
	}
	if len(donor.DonorRestrictions) > 0 {
		for _, restriction := range donor.DonorRestrictions {
			if strings.HasPrefix(restriction, "STOP") {
				return apperrors.ErrInvalidPetStatus.WithMessage("pet is not a donor")
			}
		}
	}

	// Проверяем, нет ли уже отклика от этого донора на эту заявку
	responses, err := s.donorRepo.GetDonorResponsesByDonorID(ctx, donorID)
	if err != nil {
		return apperrors.Internal(err, "failed to check existing responses")
	}
	for _, resp := range responses {
		if resp.Edges.Request.ID == reqID {
			return apperrors.ErrDonorResponseAlreadyExists
		}
	}

	// Создаем отклик
	return s.donorRepo.CreateDonorResponse(ctx, reqID, donorID, amountML)
}

// UpdateRequest обновляет информацию о заявке
func (s *BloodSearchService) UpdateRequest(ctx context.Context, id string, bloodReq *ent.UpdateBloodSearchRequestInput) error {
	// Проверяем существование
	req, err := s.bloodRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Сохраняем текущий статус, если он не передан
	if bloodReq.Status == nil {
		bloodReq.Status = &req.Status
	}

	if err := s.bloodRepo.Update(ctx, id, bloodReq); err != nil {
		return apperrors.Internal(err, "failed to update blood request")
	}

	return nil
}

// UpdateStatus обновляет статус заявки
func (s *BloodSearchService) UpdateStatus(ctx context.Context, id string, status string) error {
	if err := bloodsearchrequest.StatusValidator(bloodsearchrequest.Status(status)); err != nil {
		return apperrors.ErrInvalidBloodRequestStatus.WithInternal(err)
	}

	// Получаем заявку, чтобы узнать PetID
	req, err := s.bloodRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	err = s.txManager.WithTx(ctx, func(txCtx context.Context) error {
		entTx := ent.TxFromContext(txCtx)
		if entTx == nil {
			return apperrors.Internal(nil, "ent.TxFromContext returned nil")
		}

		// Обновляем статус заявки
		err = s.bloodRepo.UpdateStatus(txCtx, id, status)
		if err != nil {
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
	req, err := s.bloodRepo.GetByID(ctx, id)
	if err != nil {
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
