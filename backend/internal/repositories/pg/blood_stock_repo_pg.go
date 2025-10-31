package repositories

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/models"
	repositories "github.com/artesipov-alt/odnoi-krovi-app/internal/repositories/interfaces"
	"gorm.io/gorm"
)

// PostgresBloodStockRepository реализация репозитория для PostgreSQL
type PostgresBloodStockRepository struct {
	db *gorm.DB
}

// NewPostgresBloodStockRepository создает новый экземпляр репозитория
func NewPostgresBloodStockRepository(db *gorm.DB) repositories.BloodStockRepository {
	return &PostgresBloodStockRepository{
		db: db,
	}
}

// GetAll возвращает все запасы крови
func (r *PostgresBloodStockRepository) GetAll(ctx context.Context) ([]models.BloodStock, error) {
	var bloodStocks []models.BloodStock
	result := r.db.WithContext(ctx).Find(&bloodStocks)
	if result.Error != nil {
		return nil, result.Error
	}
	return bloodStocks, nil
}

// GetByID возвращает запас крови по ID
func (r *PostgresBloodStockRepository) GetByID(ctx context.Context, id int) (*models.BloodStock, error) {
	var bloodStock models.BloodStock
	result := r.db.WithContext(ctx).First(&bloodStock, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &bloodStock, nil
}

// GetByClinicID возвращает все запасы крови для конкретной клиники
func (r *PostgresBloodStockRepository) GetByClinicID(ctx context.Context, clinicID int) ([]models.BloodStock, error) {
	var bloodStocks []models.BloodStock
	result := r.db.WithContext(ctx).Where("clinic_id = ?", clinicID).Find(&bloodStocks)
	if result.Error != nil {
		return nil, result.Error
	}
	return bloodStocks, nil
}

// GetByBloodTypeID возвращает все запасы крови по типу крови
func (r *PostgresBloodStockRepository) GetByBloodTypeID(ctx context.Context, bloodTypeID int) ([]models.BloodStock, error) {
	var bloodStocks []models.BloodStock
	result := r.db.WithContext(ctx).Where("blood_type_id = ?", bloodTypeID).Find(&bloodStocks)
	if result.Error != nil {
		return nil, result.Error
	}
	return bloodStocks, nil
}

// Create создает новый запас крови
func (r *PostgresBloodStockRepository) Create(ctx context.Context, bloodStock *models.BloodStock) error {
	result := r.db.WithContext(ctx).Create(bloodStock)
	return result.Error
}

// Update обновляет существующий запас крови
func (r *PostgresBloodStockRepository) Update(ctx context.Context, bloodStock *models.BloodStock) error {
	result := r.db.WithContext(ctx).Save(bloodStock)
	return result.Error
}

// Delete удаляет запас крови по ID
func (r *PostgresBloodStockRepository) Delete(ctx context.Context, id int) error {
	result := r.db.WithContext(ctx).Delete(&models.BloodStock{}, id)
	return result.Error
}

// Search выполняет поиск запасов крови по различным фильтрам
func (r *PostgresBloodStockRepository) Search(ctx context.Context, filters repositories.BloodStockFilters) ([]models.BloodStock, error) {
	var bloodStocks []models.BloodStock
	query := r.db.WithContext(ctx)

	if filters.ClinicID != nil {
		query = query.Where("clinic_id = ?", *filters.ClinicID)
	}

	if filters.PetType != nil {
		query = query.Where("pet_type = ?", *filters.PetType)
	}

	if filters.BloodTypeID != nil {
		query = query.Where("blood_type_id = ?", *filters.BloodTypeID)
	}

	if filters.Status != nil {
		query = query.Where("status = ?", *filters.Status)
	}

	if filters.MinVolume != nil {
		query = query.Where("volume_ml >= ?", *filters.MinVolume)
	}

	if filters.MaxVolume != nil {
		query = query.Where("volume_ml <= ?", *filters.MaxVolume)
	}

	if filters.MinPrice != nil {
		query = query.Where("price_rub >= ?", *filters.MinPrice)
	}

	if filters.MaxPrice != nil {
		query = query.Where("price_rub <= ?", *filters.MaxPrice)
	}

	result := query.Find(&bloodStocks)
	if result.Error != nil {
		return nil, result.Error
	}

	return bloodStocks, nil
}
