package repositories

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/models"
	"gorm.io/gorm"
)

// PostgresBloodInfoRepo реализация репозитория для PostgreSQL
type PostgresBloodInfoRepo struct {
	db *gorm.DB
}

// NewPostgresBloodInfoRepo создает новый экземпляр репозитория
func NewPostgresBloodInfoRepo(db *gorm.DB) *PostgresBloodInfoRepo {
	return &PostgresBloodInfoRepo{
		db: db,
	}
}

// All возвращает все компоненты крови
func (r *PostgresBloodInfoRepo) AllComponents(ctx context.Context) ([]models.BloodComponent, error) {
	var BloodComponents []models.BloodComponent
	result := r.db.WithContext(ctx).Find(&BloodComponents)
	if result.Error != nil {
		return nil, result.Error
	}
	return BloodComponents, nil
}

// ByID возвращает компонент крови по ID
func (r *PostgresBloodInfoRepo) ComponentByID(ctx context.Context, id int) (*models.BloodComponent, error) {
	var BloodComponent models.BloodComponent
	result := r.db.WithContext(ctx).First(&BloodComponent, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &BloodComponent, nil
}

// BloodGroupsByPetType возвращает группы крови по типу животного
func (r *PostgresBloodInfoRepo) BloodGroupsByPetType(ctx context.Context, petType models.PetType) ([]*models.BloodGroup, error) {
	var bloodGroups []*models.BloodGroup
	result := r.db.WithContext(ctx).Where("pet_type = ?", petType).Find(&bloodGroups)
	if result.Error != nil {
		return nil, result.Error
	}
	return bloodGroups, nil
}
