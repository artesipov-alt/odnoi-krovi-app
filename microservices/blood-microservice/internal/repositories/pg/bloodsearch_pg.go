package pg

import (
	"context"
	"encoding/json"

	bloodsearchv1 "github.com/artesipov-alt/odnoi-krovi-app/microservices/blood-microservice/gen/api/bloodsearch/v1"
	"github.com/artesipov-alt/odnoi-krovi-app/microservices/blood-microservice/internal/models"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// BloodSearchRepositoryGorm реализация репозитория на GORM
type BloodSearchRepositoryGorm struct {
	db *gorm.DB
}

// NewBloodSearchRepositoryGorm создает новый экземпляр репозитория
func NewBloodSearchRepositoryGorm(db *gorm.DB) *BloodSearchRepositoryGorm {
	return &BloodSearchRepositoryGorm{db: db}
}

// AddRequest добавляет или обновляет информацию о заявке
func (r *BloodSearchRepositoryGorm) AddRequest(ctx context.Context, request *bloodsearchv1.PetRow) error {
	var petRow models.PetRow
	if err := petRow.FromProto(request); err != nil {
		return err
	}
	return r.db.WithContext(ctx).Save(&petRow).Error
}

// GetRequestByID возвращает заявку по её идентификатору
func (r *BloodSearchRepositoryGorm) GetRequestByID(ctx context.Context, requestID string) (*bloodsearchv1.PetRow, error) {
	var petRow models.PetRow
	err := r.db.WithContext(ctx).Where("pet_id = ?", requestID).First(&petRow).Error
	if err != nil {
		return nil, err
	}
	return petRow.ToProto(), nil
}

// GetRequestsByCriteria возвращает список заявок по заданным критериям
func (r *BloodSearchRepositoryGorm) GetRequestsByCriteria(ctx context.Context, criteria *bloodsearchv1.GetPetRows) ([]*bloodsearchv1.PetRow, error) {
	var petRows []*models.PetRow
	query := r.db.WithContext(ctx).Model(&models.PetRow{})

	if criteria != nil {
		if len(criteria.Regions) > 0 {
			regionsJSON, err := json.Marshal(criteria.Regions)
			if err != nil {
				return nil, err
			}
			query = query.Where("regions @> ?", datatypes.JSON(regionsJSON))
		}
		if criteria.BloodGroup != "" {
			query = query.Where("blood_group = ?", criteria.BloodGroup)
		}
		if criteria.PetType != "" {
			query = query.Where("pet_type = ?", criteria.PetType)
		}
	}

	err := query.Find(&petRows).Error
	if err != nil {
		return nil, err
	}
	return models.PetRowsToProto(petRows), nil
}

// UpdateRequestStatus обновляет статус заявки
func (r *BloodSearchRepositoryGorm) UpdateRequestStatus(ctx context.Context, requestID string, status string) error {
	return r.db.WithContext(ctx).Model(&models.PetRow{}).
		Where("pet_id = ?", requestID).
		Update("status", status).Error
}

// DeleteRequest удаляет заявку из хранилища
func (r *BloodSearchRepositoryGorm) DeleteRequest(ctx context.Context, requestID string) error {
	return r.db.WithContext(ctx).Where("pet_id = ?", requestID).Delete(&models.PetRow{}).Error
}

// GetRequestsByRegion возвращает заявки в указанном регионе
func (r *BloodSearchRepositoryGorm) GetRequestsByRegion(ctx context.Context, region int32) ([]*bloodsearchv1.PetRow, error) {
	var petRows []*models.PetRow
	regionJSON, err := json.Marshal([]int32{region})
	if err != nil {
		return nil, err
	}
	err = r.db.WithContext(ctx).Where("regions @> ?", datatypes.JSON(regionJSON)).Find(&petRows).Error
	if err != nil {
		return nil, err
	}
	return models.PetRowsToProto(petRows), nil
}

// GetRequestsByBloodGroup возвращает заявки с указанной группой крови
func (r *BloodSearchRepositoryGorm) GetRequestsByBloodGroup(ctx context.Context, bloodGroup string) ([]*bloodsearchv1.PetRow, error) {
	var petRows []*models.PetRow
	err := r.db.WithContext(ctx).Where("blood_group = ?", bloodGroup).Find(&petRows).Error
	if err != nil {
		return nil, err
	}
	return models.PetRowsToProto(petRows), nil
}

// GetAllRequests возвращает все заявки (с пагинацией)
func (r *BloodSearchRepositoryGorm) GetAllRequests(ctx context.Context, limit, offset int64) ([]*bloodsearchv1.PetRow, error) {
	var petRows []*models.PetRow
	query := r.db.WithContext(ctx).Model(&models.PetRow{})

	if limit > 0 {
		query = query.Limit(int(limit))
	}
	if offset > 0 {
		query = query.Offset(int(offset))
	}

	err := query.Find(&petRows).Error
	if err != nil {
		return nil, err
	}
	return models.PetRowsToProto(petRows), nil
}

// Exists проверяет существование заявки
func (r *BloodSearchRepositoryGorm) Exists(ctx context.Context, requestID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.PetRow{}).
		Where("pet_id = ?", requestID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Count возвращает общее количество заявок в хранилище
func (r *BloodSearchRepositoryGorm) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.PetRow{}).Count(&count).Error
	return count, err
}
