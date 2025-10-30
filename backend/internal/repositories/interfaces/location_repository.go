package repositories

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/models"
)

// LocationRepository определяет интерфейс для операций с данными локаций
type LocationRepository interface {

	// GetByID получает локацию по её ID
	GetByID(ctx context.Context, id int) (*models.Location, error)

	// GetAll получает все локации из базы данных
	GetAll(ctx context.Context) ([]*models.Location, error)
}
