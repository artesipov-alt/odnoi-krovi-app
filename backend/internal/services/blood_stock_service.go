package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/models"
	repositories "github.com/artesipov-alt/odnoi-krovi-app/internal/repositories/interfaces"
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
	ClinicID       *int                     `json:"clinic_id,omitempty" validate:"omitempty,min=1"`
	PetType        models.PetType           `json:"pet_type" validate:"required,oneof=dog cat"`
	VolumeML       *int                     `json:"volume_ml,omitempty" validate:"omitempty,min=1"`
	ExpirationDate *string                  `json:"expiration_date,omitempty"` // формат: "2024-12-31"
	Status         *models.BloodStockStatus `json:"status,omitempty" validate:"omitempty,oneof=active reserved used expired"`
	BloodTypeID    int                      `json:"blood_type_id" validate:"required,min=1"`
}

// BloodStockUpdate содержит поля для обновления запаса крови
type BloodStockUpdate struct {
	VolumeML       *int                     `json:"volume_ml,omitempty" validate:"omitempty,min=1"`
	ExpirationDate *string                  `json:"expiration_date,omitempty"`
	Status         *models.BloodStockStatus `json:"status,omitempty" validate:"omitempty,oneof=active reserved used expired"`
}

// BloodStockServiceImpl реализует BloodStockService
type BloodStockServiceImpl struct {
	bloodStockRepo repositories.BloodStockRepository
	bloodTypeRepo  repositories.BloodTypeRepository
	vetClinicRepo  repositories.VetClinicRepository
}

// NewBloodStockService создает новый сервис запасов крови
func NewBloodStockService(
	bloodStockRepo repositories.BloodStockRepository,
	bloodTypeRepo repositories.BloodTypeRepository,
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
		return nil, fmt.Errorf("получение всех запасов крови: %w", err)
	}
	return stocks, nil
}

// GetByID получает запас крови по ID
func (s *BloodStockServiceImpl) GetByID(ctx context.Context, id int) (*models.BloodStock, error) {
	stock, err := s.bloodStockRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("получение запаса крови: %w", err)
	}
	if stock == nil {
		return nil, errors.New("запас крови не найден")
	}
	return stock, nil
}

// GetByClinicID получает все запасы крови для клиники
func (s *BloodStockServiceImpl) GetByClinicID(ctx context.Context, clinicID int) ([]models.BloodStock, error) {
	// Проверяем существование клиники
	_, err := s.vetClinicRepo.GetByID(ctx, clinicID)
	if err != nil {
		return nil, fmt.Errorf("клиника не найдена: %w", err)
	}

	stocks, err := s.bloodStockRepo.GetByClinicID(ctx, clinicID)
	if err != nil {
		return nil, fmt.Errorf("получение запасов крови клиники: %w", err)
	}
	return stocks, nil
}

// GetByBloodTypeID получает все запасы крови по типу крови
func (s *BloodStockServiceImpl) GetByBloodTypeID(ctx context.Context, bloodTypeID int) ([]models.BloodStock, error) {
	// Проверяем существование типа крови
	_, err := s.bloodTypeRepo.GetByID(ctx, bloodTypeID)
	if err != nil {
		return nil, fmt.Errorf("тип крови не найден: %w", err)
	}

	stocks, err := s.bloodStockRepo.GetByBloodTypeID(ctx, bloodTypeID)
	if err != nil {
		return nil, fmt.Errorf("получение запасов крови по типу: %w", err)
	}
	return stocks, nil
}

// CreateBloodStock создает новый запас крови
func (s *BloodStockServiceImpl) CreateBloodStock(ctx context.Context, stockData BloodStockCreate) (*models.BloodStock, error) {
	// Проверяем существование типа крови
	_, err := s.bloodTypeRepo.GetByID(ctx, stockData.BloodTypeID)
	if err != nil {
		return nil, fmt.Errorf("тип крови не найден: %w", err)
	}

	// Проверяем существование клиники, если указана
	if stockData.ClinicID != nil {
		_, err := s.vetClinicRepo.GetByID(ctx, *stockData.ClinicID)
		if err != nil {
			return nil, fmt.Errorf("клиника не найдена: %w", err)
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
		return nil, fmt.Errorf("создание запаса крови: %w", err)
	}

	return stock, nil
}

// UpdateBloodStock обновляет запас крови
func (s *BloodStockServiceImpl) UpdateBloodStock(ctx context.Context, id int, updates BloodStockUpdate) error {
	// Получаем существующий запас
	stock, err := s.bloodStockRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("получение запаса крови для обновления: %w", err)
	}
	if stock == nil {
		return errors.New("запас крови не найден")
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
		return fmt.Errorf("обновление запаса крови: %w", err)
	}

	return nil
}

// DeleteBloodStock удаляет запас крови
func (s *BloodStockServiceImpl) DeleteBloodStock(ctx context.Context, id int) error {
	// Проверяем существование
	stock, err := s.bloodStockRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("получение запаса крови для удаления: %w", err)
	}
	if stock == nil {
		return errors.New("запас крови не найден")
	}

	// Удаляем
	if err := s.bloodStockRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("удаление запаса крови: %w", err)
	}

	return nil
}

// Search выполняет поиск запасов крови по фильтрам
func (s *BloodStockServiceImpl) Search(ctx context.Context, filters repositories.BloodStockFilters) ([]models.BloodStock, error) {
	stocks, err := s.bloodStockRepo.Search(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("поиск запасов крови: %w", err)
	}
	return stocks, nil
}
