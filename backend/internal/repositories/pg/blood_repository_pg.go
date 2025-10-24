package repositories

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/models"
	"gorm.io/gorm"
)

// PostgresBloodTypeRepository реализация репозитория для PostgreSQL
type PostgresBloodTypeRepository struct {
	db *gorm.DB
}

// NewPostgresBloodTypeRepository создает новый экземпляр репозитория
func NewPostgresBloodTypeRepository(db *gorm.DB) *PostgresBloodTypeRepository {
	return &PostgresBloodTypeRepository{
		db: db,
	}
}

// GetAll возвращает все типы крови
func (r *PostgresBloodTypeRepository) GetAll(ctx context.Context) ([]models.BloodType, error) {
	var bloodTypes []models.BloodType
	result := r.db.WithContext(ctx).Find(&bloodTypes)
	if result.Error != nil {
		return nil, result.Error
	}
	return bloodTypes, nil
}

// GetByID возвращает тип крови по ID
func (r *PostgresBloodTypeRepository) GetByID(ctx context.Context, id int) (*models.BloodType, error) {
	var bloodType models.BloodType
	result := r.db.WithContext(ctx).First(&bloodType, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &bloodType, nil
}

// Create создает новый тип крови
func (r *PostgresBloodTypeRepository) Create(ctx context.Context, bloodType *models.BloodType) error {
	result := r.db.WithContext(ctx).Create(bloodType)
	return result.Error
}

// Update обновляет существующий тип крови
func (r *PostgresBloodTypeRepository) Update(ctx context.Context, bloodType *models.BloodType) error {
	result := r.db.WithContext(ctx).Save(bloodType)
	return result.Error
}

// Delete удаляет тип крови по ID
func (r *PostgresBloodTypeRepository) Delete(ctx context.Context, id int) error {
	result := r.db.WithContext(ctx).Delete(&models.BloodType{}, id)
	return result.Error
}
