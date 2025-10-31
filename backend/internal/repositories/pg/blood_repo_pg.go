package repositories

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/models"
	"gorm.io/gorm"
)

// PostgresBloodRepository реализация репозитория для PostgreSQL
type PostgresBloodRepository struct {
	db *gorm.DB
}

// NewPostgresBloodRepository создает новый экземпляр репозитория
func NewPostgresBloodRepository(db *gorm.DB) *PostgresBloodRepository {
	return &PostgresBloodRepository{
		db: db,
	}
}

// GetAll возвращает все компоненты крови
func (r *PostgresBloodRepository) GetAllComponents(ctx context.Context) ([]models.BloodComponent, error) {
	var BloodComponents []models.BloodComponent
	result := r.db.WithContext(ctx).Find(&BloodComponents)
	if result.Error != nil {
		return nil, result.Error
	}
	return BloodComponents, nil
}

// GetByID возвращает компонент крови по ID
func (r *PostgresBloodRepository) GetComponentByID(ctx context.Context, id int) (*models.BloodComponent, error) {
	var BloodComponent models.BloodComponent
	result := r.db.WithContext(ctx).First(&BloodComponent, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &BloodComponent, nil
}

// GetBloodGroupsByPetType возвращает группы крови по типу животного
func (r *PostgresBloodRepository) GetBloodGroupsByPetType(ctx context.Context, petType models.PetType) ([]*models.BloodGroup, error) {
	var bloodGroups []*models.BloodGroup
	result := r.db.WithContext(ctx).Where("pet_type = ?", petType).Find(&bloodGroups)
	if result.Error != nil {
		return nil, result.Error
	}
	return bloodGroups, nil
}
