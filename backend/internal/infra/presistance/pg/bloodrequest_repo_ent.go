package pg

import (
	"context"
	"errors"
	"fmt"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent/bloodsearchrequest"
)

// EntBloodRequestRepository implements BloodRequestRepository using ENT
type EntBloodRequestRepository struct {
	db *ent.Client
}

// Внутренний хелпер для выбора клиента
func (r *EntBloodRequestRepository) client(ctx context.Context) *ent.Client {
	if tx := ent.TxFromContext(ctx); tx != nil {
		return tx.Client()
	}
	return r.db
}

// NewEntBloodRequestRepository creates a new ENT blood request repository
func NewEntBloodRequestRepository(client *ent.Client) *EntBloodRequestRepository {
	return &EntBloodRequestRepository{
		db: client,
	}
}

// toDomainModel converts ENT BloodSearchRequest to domain BloodRequest
func (r *EntBloodRequestRepository) toDomainModel(entReq *ent.BloodSearchRequest) *model.BloodRequest {
	if entReq == nil {
		return nil
	}

	// Extract response IDs from loaded edges
	var responseIDs []string
	if entReq.Edges.Responses != nil {
		responseIDs = make([]string, len(entReq.Edges.Responses))
		for i, resp := range entReq.Edges.Responses {
			responseIDs[i] = resp.ID
		}
	}

	return &model.BloodRequest{
		ID:                     entReq.ID,
		PetID:                  entReq.PetID,
		BloodVolumeNeeded:      entReq.BloodVolumeNeeded,
		BloodVolumeReserved:    entReq.BloodVolumeReserved,
		Regions:                entReq.Regions,
		SmallPetsNotifyAllowed: entReq.SmallPetsNotifyAllowed,
		Status:                 model.BloodRequestStatus(entReq.Status),
		Description:            entReq.Description,
		PhotoURLs:              entReq.PhotoUrls,
		BloodGroupNames:        entReq.BloodGroupNames,
		BloodComponentIDs:      entReq.BloodComponentIds,
		OnBoarding:             entReq.OnBoarding,
		ResponseIDs:            responseIDs,
		CreatedAt:              entReq.CreatedAt,
		UpdatedAt:              entReq.UpdatedAt,
		DeletedAt:              entReq.DeletedAt,
	}
}

// Create создает новую заявку на поиск крови
func (r *EntBloodRequestRepository) Create(ctx context.Context, req *model.BloodRequest) (*model.BloodRequest, error) {
	newBloodReq, err := r.client(ctx).BloodSearchRequest.
		Create().
		SetPetID(req.PetID).
		SetBloodVolumeNeeded(req.BloodVolumeNeeded).
		SetBloodVolumeReserved(req.BloodVolumeReserved).
		SetRegions(req.Regions).
		SetSmallPetsNotifyAllowed(req.SmallPetsNotifyAllowed).
		SetStatus(bloodsearchrequest.Status(req.Status)).
		SetDescription(req.Description).
		SetPhotoUrls(req.PhotoURLs).
		SetBloodGroupNames(req.BloodGroupNames).
		SetBloodComponentIds(req.BloodComponentIDs).
		SetOnBoarding(req.OnBoarding).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return r.toDomainModel(newBloodReq), nil
}

// GetByID возвращает заявку по её идентификатору
func (r *EntBloodRequestRepository) GetByID(ctx context.Context, id string) (*model.BloodRequest, error) {
	reqQuery := r.client(ctx).BloodSearchRequest.Query().
		Where(bloodsearchrequest.ID(id)).
		WithResponses(func(drq *ent.DonorResponseQuery) {
			drq.WithDonor(func(pq *ent.PetQuery) {
				pq.WithBloodGroupRef()
			})
		})

	req, err := reqQuery.Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, apperrors.ErrBloodRequestNotFound
		}
		return nil, apperrors.Internal(err, "failed to execute blood request query")
	}
	return r.toDomainModel(req), nil
}

// GetByPetID возвращает заявку по идентификатору питомца
func (r *EntBloodRequestRepository) GetByPetID(ctx context.Context, petID string) (*model.BloodRequest, error) {
	req, err := r.client(ctx).BloodSearchRequest.Query().
		Where(bloodsearchrequest.PetID(petID)).
		WithResponses(func(drq *ent.DonorResponseQuery) {
			drq.WithDonor(func(pq *ent.PetQuery) {
				pq.WithBloodGroupRef()
			})
		}).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, apperrors.ErrBloodRequestNotFound
		}
		return nil, apperrors.Internal(err, "failed to execute blood request query by pet ID")
	}
	return r.toDomainModel(req), nil
}

// Update обновляет информацию о заявке
func (r *EntBloodRequestRepository) Update(ctx context.Context, id string, req *model.BloodRequest) (*model.BloodRequest, error) {
	if req == nil {
		return nil, errors.New("blood request cannot be nil")
	}

	if id == "" {
		return nil, errors.New("invalid blood request ID")
	}

	updater := r.client(ctx).BloodSearchRequest.UpdateOneID(id).
		SetBloodVolumeNeeded(req.BloodVolumeNeeded).
		SetBloodVolumeReserved(req.BloodVolumeReserved).
		SetRegions(req.Regions).
		SetSmallPetsNotifyAllowed(req.SmallPetsNotifyAllowed).
		SetStatus(bloodsearchrequest.Status(req.Status)).
		SetDescription(req.Description).
		SetPhotoUrls(req.PhotoURLs).
		SetBloodGroupNames(req.BloodGroupNames).
		SetBloodComponentIds(req.BloodComponentIDs).
		SetOnBoarding(req.OnBoarding)

	updatedBloodReq, err := updater.Save(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, apperrors.ErrBloodRequestNotFound
		}
		return nil, apperrors.Internal(err, "failed to update blood request")
	}

	return r.toDomainModel(updatedBloodReq), nil
}

// UpdateStatus обновляет статус заявки
func (r *EntBloodRequestRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	return r.client(ctx).BloodSearchRequest.UpdateOneID(id).
		SetStatus(bloodsearchrequest.Status(status)).
		Exec(ctx)
}

// Delete удаляет заявку из хранилища (soft delete)
func (r *EntBloodRequestRepository) Delete(ctx context.Context, id string) error {
	return r.client(ctx).BloodSearchRequest.DeleteOneID(id).Exec(ctx)
}

// List возвращает список заявок с фильтрацией и пагинацией
func (r *EntBloodRequestRepository) List(ctx context.Context, limit, offset int, filters map[string]any) ([]*model.BloodRequest, error) {
	query := r.client(ctx).BloodSearchRequest.Query()

	if status, ok := filters["status"].(string); ok {
		query = query.Where(bloodsearchrequest.StatusEQ(bloodsearchrequest.Status(status)))
	}

	if petID, ok := filters["pet_id"].(string); ok {
		query = query.Where(bloodsearchrequest.PetID(petID))
	}

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	entReqs, err := query.All(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]*model.BloodRequest, len(entReqs))
	for i, entReq := range entReqs {
		result[i] = r.toDomainModel(entReq)
	}

	return result, nil
}

// ExistsByPetID проверяет существование активной заявки для питомца
func (r *EntBloodRequestRepository) ExistsByPetID(ctx context.Context, petID string) (bool, error) {
	return r.client(ctx).BloodSearchRequest.Query().
		Where(
			bloodsearchrequest.PetID(petID),
			bloodsearchrequest.StatusEQ(bloodsearchrequest.StatusActive),
		).
		Exist(ctx)
}

// ExistsByID проверяет существование заявки по её идентификатору
func (r *EntBloodRequestRepository) ExistsByID(ctx context.Context, id string) (bool, error) {
	return r.client(ctx).BloodSearchRequest.Query().
		Where(bloodsearchrequest.ID(id)).
		Exist(ctx)
}

// Count возвращает общее количество заявок в хранилище
func (r *EntBloodRequestRepository) Count(ctx context.Context) (int, error) {
	return r.client(ctx).BloodSearchRequest.Query().Count(ctx)
}

// AddPhotoURLs adds new photo paths to the blood request's PhotoUrls array
func (r *EntBloodRequestRepository) AddPhotoURLs(ctx context.Context, id string, paths []string) error {
	if id == "" {
		return errors.New("invalid blood request ID")
	}

	err := r.client(ctx).BloodSearchRequest.UpdateOneID(id).
		SetPhotoUrls(paths).
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to update blood request photo URLs: %w", err)
	}

	return nil
}
