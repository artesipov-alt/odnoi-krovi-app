package repositories

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/models"
)

// BloodStockRepository определяет интерфейс для работы с запасами крови
type BloodStockRepository interface {
	// GetAll возвращает все запасы крови
	GetAll(ctx context.Context) ([]models.BloodStock, error)

	// GetByID возвращает запас крови по ID
	GetByID(ctx context.Context, id int) (*models.BloodStock, error)

	// GetByClinicID возвращает все запасы крови для конкретной клиники
	GetByClinicID(ctx context.Context, clinicID int) ([]models.BloodStock, error)

	// GetByBloodTypeID возвращает все запасы крови по типу крови
	GetByBloodTypeID(ctx context.Context, bloodTypeID int) ([]models.BloodStock, error)

	// Create создает новый запас крови
	Create(ctx context.Context, bloodStock *models.BloodStock) error

	// Update обновляет существующий запас крови
	Update(ctx context.Context, bloodStock *models.BloodStock) error

	// Delete удаляет запас крови по ID
	Delete(ctx context.Context, id int) error

	// Search выполняет поиск запасов крови по различным фильтрам
	Search(ctx context.Context, filters BloodStockFilters) ([]models.BloodStock, error)
}

// BloodStockFilters представляет фильтры для поиска запасов крови
type BloodStockFilters struct {
	ClinicID    *int                     `json:"clinic_id,omitempty"`
	PetType     *models.PetType          `json:"pet_type,omitempty"`
	BloodTypeID *int                     `json:"blood_type_id,omitempty"`
	Status      *models.BloodStockStatus `json:"status,omitempty"`
	MinVolume   *int                     `json:"min_volume,omitempty"`
	MaxVolume   *int                     `json:"max_volume,omitempty"`
	MinPrice    *float64                 `json:"min_price,omitempty"`
	MaxPrice    *float64                 `json:"max_price,omitempty"`
}
