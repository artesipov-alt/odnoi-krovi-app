package pg

import (
	"context"
	"encoding/json"

	bloodrequestv1 "github.com/artesipov-alt/odnoi-krovi-app/microservices/blood-microservice/gen/api/bloodrequest/v1"
	"github.com/artesipov-alt/odnoi-krovi-app/microservices/blood-microservice/internal/models"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// BloodRequestRepositoryGorm реализация репозитория на GORM
type BloodRequestRepositoryGorm struct {
	db *gorm.DB
}

// NewBloodRequestRepositoryGorm создает новый экземпляр репозитория
func NewBloodRequestRepositoryGorm(db *gorm.DB) *BloodRequestRepositoryGorm {
	return &BloodRequestRepositoryGorm{db: db}
}

// AddRequest добавляет или обновляет информацию о заявке
func (r *BloodRequestRepositoryGorm) AddRequest(ctx context.Context, request *bloodrequestv1.BloodRequest) error {
	var bloodRequest models.BloodRequest
	if err := bloodRequest.FromProto(request); err != nil {
		return err
	}
	return r.db.WithContext(ctx).Save(&bloodRequest).Error
}

// GetRequestByID возвращает заявку по её идентификатору
func (r *BloodRequestRepositoryGorm) GetRequestByID(ctx context.Context, requestID string) (*bloodrequestv1.BloodRequest, error) {
	var bloodRequest models.BloodRequest
	err := r.db.WithContext(ctx).Where("pet_id = ?", requestID).First(&bloodRequest).Error
	if err != nil {
		return nil, err
	}
	return bloodRequest.ToProto(), nil
}

// GetRequestsByCriteria возвращает список заявок по заданным критериям
func (r *BloodRequestRepositoryGorm) GetRequestsByCriteria(ctx context.Context, criteria *bloodrequestv1.GetBloodRequests) ([]*bloodrequestv1.BloodRequest, error) {
	var bloodRequests []*models.BloodRequest
	query := r.db.WithContext(ctx).Model(&models.BloodRequest{})

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

	err := query.Find(&bloodRequests).Error
	if err != nil {
		return nil, err
	}
	return models.BloodRequestsToProto(bloodRequests), nil
}

// UpdateRequestStatus обновляет статус заявки
func (r *BloodRequestRepositoryGorm) UpdateRequestStatus(ctx context.Context, requestID string, status string) error {
	return r.db.WithContext(ctx).Model(&models.BloodRequest{}).
		Where("pet_id = ?", requestID).
		Update("status", status).Error
}

// DeleteRequest удаляет заявку из хранилища
func (r *BloodRequestRepositoryGorm) DeleteRequest(ctx context.Context, requestID string) error {
	return r.db.WithContext(ctx).Where("pet_id = ?", requestID).Delete(&models.BloodRequest{}).Error
}

// GetRequestsByRegion возвращает заявки в указанном регионе
func (r *BloodRequestRepositoryGorm) GetRequestsByRegion(ctx context.Context, region int32) ([]*bloodrequestv1.BloodRequest, error) {
	var bloodRequests []*models.BloodRequest
	regionJSON, err := json.Marshal([]int32{region})
	if err != nil {
		return nil, err
	}
	err = r.db.WithContext(ctx).Where("regions @> ?", datatypes.JSON(regionJSON)).Find(&bloodRequests).Error
	if err != nil {
		return nil, err
	}
	return models.BloodRequestsToProto(bloodRequests), nil
}

// GetRequestsByBloodGroup возвращает заявки с указанной группой крови
func (r *BloodRequestRepositoryGorm) GetRequestsByBloodGroup(ctx context.Context, bloodGroup string) ([]*bloodrequestv1.BloodRequest, error) {
	var bloodRequests []*models.BloodRequest
	err := r.db.WithContext(ctx).Where("blood_group = ?", bloodGroup).Find(&bloodRequests).Error
	if err != nil {
		return nil, err
	}
	return models.BloodRequestsToProto(bloodRequests), nil
}

// GetAllRequests возвращает все заявки (с пагинацией)
func (r *BloodRequestRepositoryGorm) GetAllRequests(ctx context.Context, limit, offset int64) ([]*bloodrequestv1.BloodRequest, error) {
	var bloodRequests []*models.BloodRequest
	query := r.db.WithContext(ctx).Model(&models.BloodRequest{})

	if limit > 0 {
		query = query.Limit(int(limit))
	}
	if offset > 0 {
		query = query.Offset(int(offset))
	}

	err := query.Find(&bloodRequests).Error
	if err != nil {
		return nil, err
	}
	return models.BloodRequestsToProto(bloodRequests), nil
}

// Exists проверяет существование заявки
func (r *BloodRequestRepositoryGorm) Exists(ctx context.Context, requestID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.BloodRequest{}).
		Where("pet_id = ?", requestID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Count возвращает общее количество заявок в хранилище
func (r *BloodRequestRepositoryGorm) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.BloodRequest{}).Count(&count).Error
	return count, err
}
