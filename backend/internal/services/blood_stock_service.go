package services

import (
	"context"
	"errors"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/models"
	repositories "github.com/artesipov-alt/odnoi-krovi-app/internal/repositories/interfaces"
	"gorm.io/gorm"
)

// BloodStockService определяет интерфейс для бизнес-логики запасов крови
type BloodStockService interface {
	// GetAll получает все запасы крови
	GetAll(ctx context.Context) ([]models.BloodStock, error)

	// GetByID получает запас крови по ID
	GetByID(ctx context.Context, id int) (*models.BloodStock, error)

	// GetByClinicID получает все запасы крови для клиники
	GetByClinicID(ctx context.Context, clinicID int) ([]models.BloodStock, error)

	// GetByBloodTypeID получает все запасы крови по типу крови
	GetByBloodTypeID(ctx context.Context, bloodTypeID int) ([]models.BloodStock, error)

	// Search выполняет поиск запасов крови по фильтрам
	Search(ctx context.Context, filters repositories.BloodStockFilters) ([]models.BloodStock, error)

	// CreateBloodStock создает новый запас крови
	CreateBloodStock(ctx context.Context, stockData BloodStockCreate) (*models.BloodStock, error)

	// UpdateBloodStock обновляет запас крови
	UpdateBloodStock(ctx context.Context, id int, updates BloodStockUpdate) error

	// DeleteBloodStock удаляет запас крови
	DeleteBloodStock(ctx context.Context, id int) error
}

// BloodStockCreate содержит данные для создания запаса крови
type BloodStockCreate struct {
	ClinicID       *int                     `json:"clinicId,omitempty" validate:"omitempty,min=1"`
	PetType        models.PetType           `json:"petType" validate:"required,oneof=dog cat"`
	VolumeML       *int                     `json:"volumeMl,omitempty" validate:"omitempty,min=1"`
	ExpirationDate *string                  `json:"expirationDate,omitempty"` // формат: "2024-12-31"
	Status         *models.BloodStockStatus `json:"status,omitempty" validate:"omitempty,oneof=active reserved used expired"`
	BloodTypeID    int                      `json:"bloodTypeId" validate:"required,min=1"`
}

// BloodStockUpdate содержит поля для обновления запаса крови
type BloodStockUpdate struct {
	VolumeML       *int                     `json:"volumeMl,omitempty" validate:"omitempty,min=1"`
	ExpirationDate *string                  `json:"expirationDate,omitempty"`
	Status         *models.BloodStockStatus `json:"status,omitempty" validate:"omitempty,oneof=active reserved used expired"`
}

// BloodStockServiceImpl реализует BloodStockService
type BloodStockServiceImpl struct {
	bloodStockRepo repositories.BloodStockRepository
	bloodTypeRepo  repositories.BloodRepository
	vetClinicRepo  repositories.VetClinicRepository
}

// NewBloodStockService создает новый сервис запасов крови
func NewBloodStockService(
	bloodStockRepo repositories.BloodStockRepository,
	bloodTypeRepo repositories.BloodRepository,
	vetClinicRepo repositories.VetClinicRepository,
) *BloodStockServiceImpl {
	return &BloodStockServiceImpl{
		bloodStockRepo: bloodStockRepo,
		bloodTypeRepo:  bloodTypeRepo,
		vetClinicRepo:  vetClinicRepo,
	}
}

// GetAll получает все запасы крови
func (s *BloodStockServiceImpl) GetAll(ctx context.Context) ([]models.BloodStock, error) {
	stocks, err := s.bloodStockRepo.GetAll(ctx)
	if err != nil {
		return nil, apperrors.Internal(err, "не удалось получить запасы крови")
	}
	return stocks, nil
}

// GetByID получает запас крови по ID
func (s *BloodStockServiceImpl) GetByID(ctx context.Context, id int) (*models.BloodStock, error) {
	stock, err := s.bloodStockRepo.GetByID(ctx, id)
	if err != nil {
		// Если запас крови не найден - возвращаем 404, а не 500
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewBloodStockNotFoundError(id)
		}
		return nil, apperrors.Internal(err, "не удалось получить запас крови")
	}
	if stock == nil {
		return nil, apperrors.NewBloodStockNotFoundError(id)
	}
	return stock, nil
}

// GetByClinicID получает все запасы крови для клиники
func (s *BloodStockServiceImpl) GetByClinicID(ctx context.Context, clinicID int) ([]models.BloodStock, error) {
	// Проверяем существование клиники
	clinic, err := s.vetClinicRepo.GetByID(ctx, clinicID)
	if err != nil {
		// Если клиника не найдена - возвращаем 404, а не 500
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewClinicNotFoundError(clinicID)
		}
		return nil, apperrors.Internal(err, "не удалось проверить существование клиники")
	}
	if clinic == nil {
		return nil, apperrors.NewClinicNotFoundError(clinicID)
	}

	stocks, err := s.bloodStockRepo.GetByClinicID(ctx, clinicID)
	if err != nil {
		return nil, apperrors.Internal(err, "не удалось получить запасы крови клиники")
	}
	return stocks, nil
}

// GetByBloodTypeID получает все запасы крови по типу крови
func (s *BloodStockServiceImpl) GetByBloodTypeID(ctx context.Context, bloodTypeID int) ([]models.BloodStock, error) {
	// Проверяем существование типа крови
	bloodType, err := s.bloodTypeRepo.GetComponentByID(ctx, bloodTypeID)
	if err != nil {
		// Если тип крови не найден - возвращаем 404, а не 500
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewBloodTypeNotFoundError(bloodTypeID)
		}
		return nil, apperrors.Internal(err, "не удалось проверить существование типа крови")
	}
	if bloodType == nil {
		return nil, apperrors.NewBloodTypeNotFoundError(bloodTypeID)
	}

	stocks, err := s.bloodStockRepo.GetByBloodTypeID(ctx, bloodTypeID)
	if err != nil {
		return nil, apperrors.Internal(err, "не удалось получить запасы крови по типу")
	}
	return stocks, nil
}

// CreateBloodStock создает новый запас крови
func (s *BloodStockServiceImpl) CreateBloodStock(ctx context.Context, stockData BloodStockCreate) (*models.BloodStock, error) {
	// Проверяем существование типа крови
	bloodType, err := s.bloodTypeRepo.GetComponentByID(ctx, stockData.BloodTypeID)
	if err != nil {
		// Если тип крови не найден - возвращаем 404, а не 500
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewBloodTypeNotFoundError(stockData.BloodTypeID)
		}
		return nil, apperrors.Internal(err, "не удалось проверить существование типа крови")
	}
	if bloodType == nil {
		return nil, apperrors.NewBloodTypeNotFoundError(stockData.BloodTypeID)
	}

	// Проверяем существование клиники, если указана
	if stockData.ClinicID != nil {
		clinic, err := s.vetClinicRepo.GetByID(ctx, *stockData.ClinicID)
		if err != nil {
			// Если клиника не найдена - возвращаем 404, а не 500
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, apperrors.NewClinicNotFoundError(*stockData.ClinicID)
			}
			return nil, apperrors.Internal(err, "не удалось проверить существование клиники")
		}
		if clinic == nil {
			return nil, apperrors.NewClinicNotFoundError(*stockData.ClinicID)
		}
	}

	// Создаем новый запас крови
	stock := &models.BloodStock{
		ClinicID:    stockData.ClinicID,
		PetType:     stockData.PetType,
		VolumeML:    stockData.VolumeML,
		BloodTypeID: stockData.BloodTypeID,
	}

	// Устанавливаем статус (по умолчанию active)
	if stockData.Status != nil {
		stock.Status = *stockData.Status
	} else {
		stock.Status = models.BloodStockStatusActive
	}

	// Парсим дату истечения срока, если указана
	if stockData.ExpirationDate != nil && *stockData.ExpirationDate != "" {
		// Здесь можно добавить парсинг даты при необходимости
	}

	if err := s.bloodStockRepo.Create(ctx, stock); err != nil {
		return nil, apperrors.Internal(err, "не удалось создать запас крови")
	}

	return stock, nil
}

// UpdateBloodStock обновляет запас крови
func (s *BloodStockServiceImpl) UpdateBloodStock(ctx context.Context, id int, updates BloodStockUpdate) error {
	// Получаем существующий запас
	stock, err := s.bloodStockRepo.GetByID(ctx, id)
	if err != nil {
		// Если запас крови не найден - возвращаем 404, а не 500
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.NewBloodStockNotFoundError(id)
		}
		return apperrors.Internal(err, "не удалось получить запас крови")
	}
	if stock == nil {
		return apperrors.NewBloodStockNotFoundError(id)
	}

	// Применяем обновления
	if updates.VolumeML != nil {
		stock.VolumeML = updates.VolumeML
	}
	if updates.Status != nil {
		stock.Status = *updates.Status
	}
	if updates.ExpirationDate != nil {
		// Здесь можно добавить парсинг даты при необходимости
	}

	// Сохраняем обновления
	if err := s.bloodStockRepo.Update(ctx, stock); err != nil {
		return apperrors.Internal(err, "не удалось обновить запас крови")
	}

	return nil
}

// DeleteBloodStock удаляет запас крови
func (s *BloodStockServiceImpl) DeleteBloodStock(ctx context.Context, id int) error {
	// Проверяем существование
	stock, err := s.bloodStockRepo.GetByID(ctx, id)
	if err != nil {
		// Если запас крови не найден - возвращаем 404, а не 500
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.NewBloodStockNotFoundError(id)
		}
		return apperrors.Internal(err, "не удалось получить запас крови")
	}
	if stock == nil {
		return apperrors.NewBloodStockNotFoundError(id)
	}

	// Удаляем
	if err := s.bloodStockRepo.Delete(ctx, id); err != nil {
		return apperrors.Internal(err, "не удалось удалить запас крови")
	}

	return nil
}

// Search выполняет поиск запасов крови по фильтрам
func (s *BloodStockServiceImpl) Search(ctx context.Context, filters repositories.BloodStockFilters) ([]models.BloodStock, error) {
	stocks, err := s.bloodStockRepo.Search(ctx, filters)
	if err != nil {
		return nil, apperrors.Internal(err, "не удалось выполнить поиск запасов крови")
	}
	return stocks, nil
}
