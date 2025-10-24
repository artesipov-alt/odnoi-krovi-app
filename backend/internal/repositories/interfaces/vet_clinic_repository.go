package repositories

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/models"
)

// VetClinicRepository определяет интерфейс для операций с данными ветеринарных клиник
type VetClinicRepository interface {
	// Create создает новую ветеринарную клинику в базе данных
	Create(ctx context.Context, clinic *models.VetClinic) error

	// GetByID получает ветеринарную клинику по её ID
	GetByID(ctx context.Context, id int) (*models.VetClinic, error)

	// GetByLocationID получает все ветеринарные клиники по ID локации
	GetByLocationID(ctx context.Context, locationID int) ([]*models.VetClinic, error)

	// Update обновляет существующую ветеринарную клинику в базе данных
	Update(ctx context.Context, clinic *models.VetClinic) error

	// Delete удаляет ветеринарную клинику по её ID (soft delete)
	Delete(ctx context.Context, id int) error
}
