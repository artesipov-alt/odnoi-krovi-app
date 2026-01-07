package pg

import (
	"context"
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/ent"
	"github.com/artesipov-alt/odnoi-krovi-app/ent/bloodsearchrequest"
)

// EntBloodRequestRepository implements BloodRequestRepository using ENT
type EntBloodRequestRepository struct {
	client *ent.Client
}

// NewEntBloodRequestRepository creates a new ENT blood request repository
func NewEntBloodRequestRepository(client *ent.Client) *EntBloodRequestRepository {
	return &EntBloodRequestRepository{
		client: client,
	}
}

// Create создает новую заявку на поиск крови
func (r *EntBloodRequestRepository) Create(ctx context.Context, req *ent.BloodSearchRequest) (*ent.BloodSearchRequest, error) {
	return r.client.BloodSearchRequest.Create().
		SetPetID(req.PetID).
		SetBloodVolumeNeeded(req.BloodVolumeNeeded).
		SetBloodVolumeReserved(req.BloodVolumeReserved).
		SetRegions(req.Regions).
		SetSmallPetsNotifyAllowed(req.SmallPetsNotifyAllowed).
		SetStatus(bloodsearchrequest.Status(req.Status)).
		SetNillableDescription(&req.Description).
		SetPhotoUrls(req.PhotoUrls).
		SetBloodGroupIds(req.BloodGroupIds).
		SetBloodComponentIds(req.BloodComponentIds).
		Save(ctx)
}

// GetByID возвращает заявку по её идентификатору
func (r *EntBloodRequestRepository) GetByID(ctx context.Context, id string) (*ent.BloodSearchRequest, error) {
	return r.client.BloodSearchRequest.Get(ctx, id)
}

// GetByPetID возвращает заявку по идентификатору питомца
func (r *EntBloodRequestRepository) GetByPetID(ctx context.Context, petID string) (*ent.BloodSearchRequest, error) {
	return r.client.BloodSearchRequest.Query().
		Where(bloodsearchrequest.PetID(petID)).
		Only(ctx)
}

// Update обновляет информацию о заявке
func (r *EntBloodRequestRepository) Update(ctx context.Context, req *ent.BloodSearchRequest) (*ent.BloodSearchRequest, error) {
	return r.client.BloodSearchRequest.UpdateOneID(req.ID).
		SetBloodVolumeNeeded(req.BloodVolumeNeeded).
		SetBloodVolumeReserved(req.BloodVolumeReserved).
		SetRegions(req.Regions).
		SetSmallPetsNotifyAllowed(req.SmallPetsNotifyAllowed).
		SetStatus(bloodsearchrequest.Status(req.Status)).
		SetDescription(req.Description).
		SetPhotoUrls(req.PhotoUrls).
		SetBloodGroupIds(req.BloodGroupIds).
		SetBloodComponentIds(req.BloodComponentIds).
		Save(ctx)
}

// UpdateStatus обновляет статус заявки
func (r *EntBloodRequestRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	return r.client.BloodSearchRequest.UpdateOneID(id).
		SetStatus(bloodsearchrequest.Status(status)).
		Exec(ctx)
}

// Delete удаляет заявку из хранилища (soft delete)
func (r *EntBloodRequestRepository) Delete(ctx context.Context, id string) error {
	return r.client.BloodSearchRequest.UpdateOneID(id).
		SetDeletedAt(time.Now()).
		Exec(ctx)
}

// List возвращает список заявок с фильтрацией и пагинацией
func (r *EntBloodRequestRepository) List(ctx context.Context, limit, offset int, filters map[string]interface{}) ([]*ent.BloodSearchRequest, error) {
	query := r.client.BloodSearchRequest.Query()

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

	return query.All(ctx)
}

// ExistsByPetID проверяет существование активной заявки для питомца
func (r *EntBloodRequestRepository) ExistsByPetID(ctx context.Context, petID string) (bool, error) {
	return r.client.BloodSearchRequest.Query().
		Where(
			bloodsearchrequest.PetID(petID),
			bloodsearchrequest.StatusEQ(bloodsearchrequest.StatusActive),
		).
		Exist(ctx)
}

// Count возвращает общее количество заявок в хранилище
func (r *EntBloodRequestRepository) Count(ctx context.Context) (int, error) {
	return r.client.BloodSearchRequest.Query().Count(ctx)
}
