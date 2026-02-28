package reference

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/reference/model"
)

// LocationRepository определяет интерфейс для операций с данными локаций
type LocationRepository interface {
	// GetByID получает локацию по её ID
	GetByID(ctx context.Context, id string) (*model.Location, error)

	// GetAll получает все локации из базы данных
	GetAll(ctx context.Context) ([]*model.Location, error)

	// Exists проверяет, существует ли локация с заданным ID
	Exists(ctx context.Context, id string) (bool, error)
}
