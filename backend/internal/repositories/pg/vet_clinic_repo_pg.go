package repositories

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/models"
	"gorm.io/gorm"
)

// VetClinicRepositoryImpl реализует VetClinicRepository для PostgreSQL
type VetClinicRepositoryImpl struct {
	db *gorm.DB
}

// NewVetClinicRepository создает новый репозиторий ветеринарных клиник
func NewVetClinicRepository(db *gorm.DB) *VetClinicRepositoryImpl {
	return &VetClinicRepositoryImpl{
		db: db,
	}
}

// Create создает новую ветеринарную клинику в базе данных
func (r *VetClinicRepositoryImpl) Create(ctx context.Context, clinic *models.VetClinic) error {
	return r.db.WithContext(ctx).Create(clinic).Error
}

// GetByID получает клинику по её ID
func (r *VetClinicRepositoryImpl) GetByID(ctx context.Context, id int) (*models.VetClinic, error) {
	var clinic models.VetClinic
	err := r.db.WithContext(ctx).Where("clinic_id = ?", id).First(&clinic).Error
	if err != nil {
		return nil, err
	}
	return &clinic, nil
}

// GetByLocationID получает все клиники по ID локации
func (r *VetClinicRepositoryImpl) GetByLocationID(ctx context.Context, locationID int) ([]*models.VetClinic, error) {
	var clinics []*models.VetClinic
	err := r.db.WithContext(ctx).Where("location_id = ?", locationID).Find(&clinics).Error
	if err != nil {
		return nil, err
	}
	return clinics, nil
}

// Update обновляет существующую клинику в базе данных
func (r *VetClinicRepositoryImpl) Update(ctx context.Context, clinic *models.VetClinic) error {
	return r.db.WithContext(ctx).Save(clinic).Error
}

// Delete удаляет клинику по её ID (soft delete)
func (r *VetClinicRepositoryImpl) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Where("clinic_id = ?", id).Delete(&models.VetClinic{}).Error
}
