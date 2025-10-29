package repositories

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/models"
)

// LocationRepository определяет интерфейс для операций с данными локаций
type LocationRepository interface {
	// Create создает новую локацию в базе данных
	Create(ctx context.Context, location *models.Location) error

	// GetByID получает локацию по её ID
	GetByID(ctx context.Context, id int) (*models.Location, error)

	// GetAll получает все локации из базы данных
	GetAll(ctx context.Context) ([]*models.Location, error)

	// Update обновляет существующую локацию в базе данных
	Update(ctx context.Context, location *models.Location) error

	// Delete удаляет локацию по её ID
	Delete(ctx context.Context, id int) error

	// ExistsByName проверяет, существует ли локация с заданным названием
	ExistsByName(ctx context.Context, name string) (bool, error)
}
