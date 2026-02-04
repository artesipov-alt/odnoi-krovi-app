package pg

import (
	"context"
	"errors"
	"fmt"

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
		SetBloodGroupNames(req.BloodGroupNames).
		SetBloodComponentIds(req.BloodComponentIds).
		Save(ctx)
}

// CreateWithTx создает новую заявку на поиск крови в рамках транзакции
func (r *EntBloodRequestRepository) CreateWithTx(ctx context.Context, tx *ent.Tx, req *ent.BloodSearchRequest) (*ent.BloodSearchRequest, error) {
	return tx.BloodSearchRequest.Create().
		SetPetID(req.PetID).
		SetBloodVolumeNeeded(req.BloodVolumeNeeded).
		SetBloodVolumeReserved(req.BloodVolumeReserved).
		SetRegions(req.Regions).
		SetSmallPetsNotifyAllowed(req.SmallPetsNotifyAllowed).
		SetStatus(bloodsearchrequest.Status(req.Status)).
		SetNillableDescription(&req.Description).
		SetPhotoUrls(req.PhotoUrls).
		SetBloodGroupNames(req.BloodGroupNames).
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
		SetBloodGroupNames(req.BloodGroupNames).
		SetBloodComponentIds(req.BloodComponentIds).
		Save(ctx)
}

// UpdateStatus обновляет статус заявки
func (r *EntBloodRequestRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	return r.client.BloodSearchRequest.UpdateOneID(id).
		SetStatus(bloodsearchrequest.Status(status)).
		Exec(ctx)
}

// UpdateStatusWithTx обновляет статус заявки в рамках транзакции
func (r *EntBloodRequestRepository) UpdateStatusWithTx(ctx context.Context, tx *ent.Tx, id string, status string) error {
	return tx.BloodSearchRequest.UpdateOneID(id).
		SetStatus(bloodsearchrequest.Status(status)).
		Exec(ctx)
}

// Delete удаляет заявку из хранилища (soft delete)
func (r *EntBloodRequestRepository) Delete(ctx context.Context, id string) error {
	return r.client.BloodSearchRequest.DeleteOneID(id).Exec(ctx)
}

// DeleteWithTx удаляет заявку из хранилища в рамках транзакции (soft delete)
func (r *EntBloodRequestRepository) DeleteWithTx(ctx context.Context, tx *ent.Tx, id string) error {
	return tx.BloodSearchRequest.DeleteOneID(id).Exec(ctx)
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

// ExistsByID проверяет существование заявки по её идентификатору
func (r *EntBloodRequestRepository) ExistsByID(ctx context.Context, id string) (bool, error) {
	return r.client.BloodSearchRequest.Query().
		Where(bloodsearchrequest.ID(id)).
		Exist(ctx)
}

// Count возвращает общее количество заявок в хранилище
func (r *EntBloodRequestRepository) Count(ctx context.Context) (int, error) {
	return r.client.BloodSearchRequest.Query().Count(ctx)
}

// AddPhotoURLs adds new photo paths to the blood request's PhotoUrls array
func (r *EntBloodRequestRepository) AddPhotoURLs(ctx context.Context, id string, paths []string) error {
	if id == "" {
		return errors.New("invalid blood request ID")
	}

	// Fetch current photo URLs
	req, err := r.client.BloodSearchRequest.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get blood request for photo update: %w", err)
	}

	// Append new paths
	newPhotoUrls := append(req.PhotoUrls, paths...)

	// Update blood request
	err = r.client.BloodSearchRequest.UpdateOneID(id).
		SetPhotoUrls(newPhotoUrls).
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to update blood request photo URLs: %w", err)
	}

	return nil
}
