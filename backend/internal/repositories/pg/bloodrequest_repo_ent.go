package pg

import (
	"context"
	"errors"
	"fmt"

	"github.com/artesipov-alt/odnoi-krovi-app/ent"
	"github.com/artesipov-alt/odnoi-krovi-app/ent/bloodsearchrequest"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
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

// Create создает новую заявку на поиск крови
func (r *EntBloodRequestRepository) Create(ctx context.Context, input *ent.CreateBloodSearchRequestInput) (*ent.BloodSearchRequest, error) {
	newBloodReq, err := r.client(ctx).BloodSearchRequest.
		Create().
		SetInput(*input).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return newBloodReq, nil
}

// GetByID возвращает заявку по её идентификатору
func (r *EntBloodRequestRepository) GetByID(ctx context.Context, id string) (*ent.BloodSearchRequest, error) {
	reqQuery := r.client(ctx).BloodSearchRequest.Query().
		Where(bloodsearchrequest.ID(id))

	req, err := reqQuery.Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) { // This case should ideally be caught by the initial GetByID, but good for defensive programming
			return nil, apperrors.ErrBloodRequestNotFound
		}
		return nil, apperrors.Internal(err, "failed to execute blood request query")
	}
	return req, nil
}

// GetByPetID возвращает заявку по идентификатору питомца
func (r *EntBloodRequestRepository) GetByPetID(ctx context.Context, petID string) (*ent.BloodSearchRequest, error) {
	req, err := r.client(ctx).BloodSearchRequest.Query().
		Where(bloodsearchrequest.PetID(petID)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, apperrors.ErrBloodRequestNotFound
		}
		return nil, apperrors.Internal(err, "failed to execute blood request query by pet ID")
	}
	return req, nil
}

// Update обновляет информацию о заявке
func (r *EntBloodRequestRepository) Update(ctx context.Context, id string, input *ent.UpdateBloodSearchRequestInput) error {
	if input == nil {
		return errors.New("blood request cannot be nil")
	}

	if id == "" {
		return errors.New("invalid blood request ID")
	}

	_, err := r.db.BloodSearchRequest.UpdateOneID(id).
		SetInput(*input).
		Save(ctx)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
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
func (r *EntBloodRequestRepository) List(ctx context.Context, limit, offset int, filters map[string]any) ([]*ent.BloodSearchRequest, error) {
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

	return query.All(ctx)
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

	// Fetch current photo URLs
	req, err := r.client(ctx).BloodSearchRequest.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get blood request for photo update: %w", err)
	}

	// Append new paths
	newPhotoUrls := append(req.PhotoUrls, paths...)

	// Update blood request
	err = r.client(ctx).BloodSearchRequest.UpdateOneID(id).
		SetPhotoUrls(newPhotoUrls).
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to update blood request photo URLs: %w", err)
	}

	return nil
}
